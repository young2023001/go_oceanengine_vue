//go:build !cloudstudio
// +build !cloudstudio

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"oceanengine-backend/config"
	accountModel "oceanengine-backend/internal/app/account/model"
	adModel "oceanengine-backend/internal/app/ad/model"
	adminModel "oceanengine-backend/internal/app/admin/model"
	advertiserModel "oceanengine-backend/internal/app/advertiser/model"
	analyticsModel "oceanengine-backend/internal/app/analytics/model"
	audienceModel "oceanengine-backend/internal/app/audience/model"
	batchModel "oceanengine-backend/internal/app/batch/model"
	campaignModel "oceanengine-backend/internal/app/campaign/model"
	creativeModel "oceanengine-backend/internal/app/creative/model"
	enterpriseModel "oceanengine-backend/internal/app/enterprise/model"
	groupModel "oceanengine-backend/internal/app/group/model"
	mediaModel "oceanengine-backend/internal/app/media/model"
	projectModel "oceanengine-backend/internal/app/project/model"
	reportModel "oceanengine-backend/internal/app/report/model"
	scopeModel "oceanengine-backend/internal/app/scope/model"
	templateModel "oceanengine-backend/internal/app/template/model"
	tenantModel "oceanengine-backend/internal/app/tenant/model"
	batchService "oceanengine-backend/internal/app/batch/service"
	batchRepository "oceanengine-backend/internal/app/batch/repository"
	tenantRepository "oceanengine-backend/internal/app/tenant/repository"
	tenantService "oceanengine-backend/internal/app/tenant/service"
	"github.com/redis/go-redis/v9"
	"oceanengine-backend/internal/router"
	"oceanengine-backend/internal/scheduler"
	"oceanengine-backend/pkg/auth"
	"oceanengine-backend/pkg/database"
	"oceanengine-backend/pkg/logger"
	"oceanengine-backend/pkg/ocean"
	"oceanengine-backend/pkg/oceanengine"
	"oceanengine-backend/pkg/ratelimiter"
)

func main() {
	// 加载配置
	cfg, err := config.Load("config/settings.yml")
	if err != nil {
		fmt.Printf("加载配置失败: %v\n", err)
		os.Exit(1)
	}

	// 初始化日志
	log, err := logger.Init(&cfg.Logger)
	if err != nil {
		fmt.Printf("初始化日志失败: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	log.Info("启动 OceanEngine 后台管理系统")

	// 初始化数据库
	db, err := database.Init(&cfg.Database, log)
	if err != nil {
		log.Fatal(fmt.Sprintf("初始化数据库失败: %v", err))
	}
	log.Info("数据库连接成功")

	// 自动迁移所有表
	if err := db.AutoMigrate(
		// 系统管理
		&adminModel.User{},
		&adminModel.Role{},
		&adminModel.Menu{},
		&adminModel.RoleMenu{},
		&adminModel.OperationLog{},
		&adminModel.UserSetting{},
		&adminModel.Notification{},
		&adminModel.DictType{},
		&adminModel.DictData{},
		// 租户
		&tenantModel.Tenant{},
		// 账户与分组
		&accountModel.LocalAccount{},
		&accountModel.Store{},
		&groupModel.AccountGroup{},
		&groupModel.AccountGroupMember{},
		&scopeModel.UserAccountScope{},
		// 广告主
		&advertiserModel.Advertiser{},
		&advertiserModel.AdvertiserFund{},
		// 广告系列、广告组、创意
		&campaignModel.Campaign{},
		&adModel.Ad{},
		&creativeModel.Creative{},
		// 项目与计划
		&projectModel.LocalProject{},
		&projectModel.LocalPromotion{},
		// 批量任务
		&batchModel.BatchTask{},
		&batchModel.BatchTaskItem{},
		// 模板
		&templateModel.ProjectTemplate{},
		&templateModel.PromotionTemplate{},
		// 报表
		&analyticsModel.ReportDaily{},
		&analyticsModel.ExportTask{},
		&reportModel.AdvertiserReport{},
		&reportModel.CampaignReport{},
		&reportModel.AdReport{},
		&reportModel.ExportTask{},
		// 素材
		&mediaModel.MaterialImage{},
		&mediaModel.MaterialVideo{},
		// 人群定向
		&audienceModel.AudiencePackage{},
		&audienceModel.CustomAudience{},
		// 企业号
		&enterpriseModel.ReplyTemplate{},
	); err != nil {
		log.Warn(fmt.Sprintf("auto migrate warning (non-fatal): %v", err))
	}
	log.Info("数据库迁移完成")

	// 初始化 Redis (可选)
	var redisClient *redis.Client
	if cfg.Redis.Addr != "" {
		redisClient, err = database.InitRedis(&cfg.Redis, log)
		if err != nil {
			log.Warn(fmt.Sprintf("初始化 Redis 失败: %v", err))
		} else {
			log.Info("Redis 连接成功")
		}
	}

	// 初始化 JWT 管理器
	jwtManager := auth.NewJWTManager(&cfg.JWT)

	// 设置路由
	r := router.NewRouter(db, log, jwtManager, &cfg.Ocean)
	engine := r.Setup(cfg.Server.Mode)

	// Token 自动续期
	refreshCtx, refreshCancel := context.WithCancel(context.Background())
	defer refreshCancel()

	tRepo := tenantRepository.NewTenantRepository(db)
	oauthClient := tenantService.NewOAuthClient()
	tokenRefresher := tenantService.NewTokenRefresher(tRepo, oauthClient, log)
	go tokenRefresher.Start(refreshCtx)

	// 批量任务断点续传
	recovery := batchService.NewTaskRecovery(db, log)
	recovery.RecoverOnStartup(refreshCtx)

	// 启动批量任务 Worker
	batchRepo := batchRepository.NewBatchRepository(db)
	oceanEngineClient := oceanengine.NewClient(cfg.Ocean.AppID, cfg.Ocean.Secret)
	tokenStore := tenantService.NewDBTokenStore(db)
	if redisClient != nil {
		workerLimiter := ratelimiter.NewTenantRateLimiter(redisClient, 10)
		worker := batchService.NewWorker(batchRepo, oceanEngineClient, workerLimiter, tokenStore, log)
		go worker.Start(refreshCtx)
	}

	// 项目同步定时任务
	if redisClient != nil {
		rateLimiter := ratelimiter.NewTenantRateLimiter(redisClient, 10)
		oceanClient := ocean.NewClient()
		syncer := scheduler.NewProjectSyncer(db, oceanClient, rateLimiter, log)
		go syncer.Start(refreshCtx)
	}

	// 创建 HTTP 服务器
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      engine,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// 启动服务器（非阻塞）
	go func() {
		log.Info(fmt.Sprintf("服务器启动于 %s", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(fmt.Sprintf("服务器启动失败: %v", err))
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("正在关闭服务器...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(fmt.Sprintf("服务器关闭失败: %v", err))
	}

	log.Info("服务器已关闭")
}
