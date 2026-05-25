import http from 'node:http'

const port = process.env.PORT || 8080

const cors = (res) => {
  res.setHeader('Access-Control-Allow-Origin', '*')
  res.setHeader('Access-Control-Allow-Methods', 'GET, POST, PUT, DELETE, OPTIONS')
  res.setHeader('Access-Control-Allow-Headers', 'Content-Type, Authorization, X-Access-Token, Access-Token')
}

const json = (res, data) => {
  res.setHeader('Content-Type', 'application/json')
  res.end(JSON.stringify(data))
}

const routes = {
  '/health': (req, res) => {
    json(res, { status: 'ok', timestamp: new Date().toISOString() })
  },

  '/api/v1/info': (req, res) => {
    json(res, { app: 'OceanEngine Backend', version: '1.0.0', environment: 'local-mock', timestamp: new Date().toISOString() })
  },

  '/api/v1/auth/login': (req, res) => {
    json(res, {
      code: 0,
      message: '登录成功',
      data: {
        token: 'mock_token_for_local_dev',
        user: { id: 1, username: 'admin', nickname: '管理员', email: 'admin@example.com', role: 'admin' }
      }
    })
  },

  '/api/v1/auth/me': (req, res) => {
    json(res, {
      code: 0,
      message: 'success',
      data: { id: 1, username: 'admin', nickname: '管理员', email: 'admin@example.com', role: 'admin', avatar: '' }
    })
  },

  '/api/v1/advertisers': (req, res) => {
    json(res, {
      code: 0,
      message: 'success',
      data: {
        list: [
          { id: 1, name: '测试广告主1', status: 'active', created_at: '2023-01-01T00:00:00Z' },
          { id: 2, name: '测试广告主2', status: 'active', created_at: '2023-01-02T00:00:00Z' }
        ],
        total: 2, page: 1, page_size: 10
      }
    })
  },

  '/api/v1/local/projects': (req, res) => {
    json(res, {
      code: 0,
      message: 'success',
      data: {
        list: [
          { project_id: 1001, project_name: '本地推-餐饮门店', advertiser_id: 1, status: 'enable', budget: 500, budget_mode: 'day', delivery_range: '5km', create_time: '2024-06-01' },
          { project_id: 1002, project_name: '本地推-美容美发', advertiser_id: 1, status: 'enable', budget: 300, budget_mode: 'day', delivery_range: '3km', create_time: '2024-06-10' }
        ],
        total: 2, page: 1, page_size: 10
      }
    })
  },

  '/api/v1/local/promotions': (req, res) => {
    json(res, {
      code: 0,
      message: 'success',
      data: {
        list: [
          { promotion_id: 2001, promotion_name: '午餐优惠券推广', project_id: 1001, status: 'enable', opt_status: 'enable', budget: 200, create_time: '2024-06-02' },
          { promotion_id: 2002, promotion_name: '新店开业活动', project_id: 1001, status: 'enable', opt_status: 'enable', budget: 300, create_time: '2024-06-05' }
        ],
        total: 2, page: 1, page_size: 10
      }
    })
  },

  '/api/v1/local/clues': (req, res) => {
    json(res, {
      code: 0,
      message: 'success',
      data: {
        list: [
          { clue_id: 3001, promotion_id: 2001, clue_type: 1, clue_source: 1, create_time: '2024-06-15 10:30:00', phone: '138****1234', name: '张先生', address: '朝阳区望京', follow_status: 0 },
          { clue_id: 3002, promotion_id: 2001, clue_type: 2, clue_source: 2, create_time: '2024-06-15 14:20:00', phone: '139****5678', name: '李女士', address: '海淀区中关村', follow_status: 1 }
        ],
        total: 2, page: 1, page_size: 10
      }
    })
  },

  '/api/v1/local/stores': (req, res) => {
    json(res, {
      code: 0,
      message: 'success',
      data: {
        list: [
          { poi_id: 'poi_001', poi_name: '好味道餐厅(望京店)', address: '北京市朝阳区望京街道', latitude: 39.99, longitude: 116.48, phone: '010-12345678', status: 1 },
          { poi_id: 'poi_002', poi_name: '好味道餐厅(中关村店)', address: '北京市海淀区中关村大街', latitude: 39.98, longitude: 116.31, phone: '010-87654321', status: 1 }
        ],
        total: 2, page: 1, page_size: 10
      }
    })
  },

  '/api/v1/local/report': (req, res) => {
    json(res, {
      code: 0,
      message: 'success',
      data: {
        list: [
          { stat_datetime: '2024-06-14', cost: 156.8, show_cnt: 12500, click_cnt: 380, convert_cnt: 25, ctr: 3.04, cpm: 12.54, cpc: 0.41, cpa_convert: 6.27 },
          { stat_datetime: '2024-06-15', cost: 203.5, show_cnt: 15800, click_cnt: 520, convert_cnt: 38, ctr: 3.29, cpm: 12.88, cpc: 0.39, cpa_convert: 5.36 }
        ],
        total: 2, page: 1, page_size: 10
      }
    })
  },

  '/api/v1/local/materials': (req, res) => {
    json(res, {
      code: 0,
      message: 'success',
      data: {
        list: [
          { video_id: 'v001', video_url: '', cover_url: '', width: 1080, height: 1920, duration: 15.0, status: 1, create_time: '2024-06-01' },
          { video_id: 'v002', video_url: '', cover_url: '', width: 1080, height: 1920, duration: 30.0, status: 1, create_time: '2024-06-05' }
        ],
        total: 2, page: 1, page_size: 10
      }
    })
  }
}

const server = http.createServer((req, res) => {
  cors(res)
  if (req.method === 'OPTIONS') { res.writeHead(200); res.end(); return }

  const url = new URL(req.url, `http://localhost:${port}`)
  const handler = routes[url.pathname]

  if (handler) {
    handler(req, res)
  } else if (url.pathname.startsWith('/api/')) {
    json(res, { code: 0, message: 'Mock API', path: url.pathname, timestamp: new Date().toISOString() })
  } else {
    res.writeHead(200, { 'Content-Type': 'text/html' })
    res.end('<h1>OceanEngine Mock Backend Running</h1>')
  }
})

server.listen(port, () => {
  console.log(`🚀 Mock Backend 启动成功 http://localhost:${port}`)
  console.log(`🏥 健康检查: http://localhost:${port}/health`)
  console.log(`📍 本地推API: http://localhost:${port}/api/v1/local/projects`)
})
