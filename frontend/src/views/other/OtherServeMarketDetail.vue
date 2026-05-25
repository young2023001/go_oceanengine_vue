<script setup lang="ts">
import { ref } from 'vue'
import Breadcrumb from '@/components/common/Breadcrumb.vue'

const service = ref({
  id: 'SV001',
  name: '智能投放助手',
  provider: '巨量引擎官方',
  category: '投放优化',
  rating: 4.8,
  reviews: 256,
  users: 12500,
  price: '免费',
  description: '基于AI的智能投放优化工具，自动调整出价和定向，提升广告效果。',
  features: ['自动出价优化', '智能定向调整', '实时效果监控', '异常预警通知']
})

const reviews = ref([
  { id: 'R001', user: '广告主A', rating: 5, content: '效果非常好，ROI提升了30%', date: '2025-11-25' },
  { id: 'R002', user: '广告主B', rating: 4, content: '功能很实用，希望增加更多分析维度', date: '2025-11-20' },
  { id: 'R003', user: '广告主C', rating: 5, content: '省时省力，强烈推荐', date: '2025-11-15' }
])

const handleUseService = () => {
  // TODO: implement
}
</script>

<template>
  <div class="space-y-6">
    <div>
      <Breadcrumb :items="[{ name: '服务市场' }, { name: '服务详情' }]" />
    </div>

    <div class="bg-white rounded-lg border border-gray-200 p-6">
      <div class="flex gap-6">
        <div class="w-24 h-24 bg-gradient-to-br from-blue-500 to-purple-500 rounded-xl flex items-center justify-center text-white text-4xl">
          🤖
        </div>
        <div class="flex-1">
          <div class="flex items-start justify-between">
            <div>
              <h1 class="text-2xl font-bold text-gray-900">{{ service.name }}</h1>
              <p class="text-gray-500 mt-1">{{ service.provider }} · {{ service.category }}</p>
            </div>
<button class="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700" @click="handleUseService">
              立即使用
            </button>
          </div>
          <div class="flex items-center gap-6 mt-4">
            <div class="flex items-center gap-1">
              <span class="text-yellow-500">⭐</span>
              <span class="font-medium">{{ service.rating }}</span>
              <span class="text-gray-500">({{ service.reviews }} 评价)</span>
            </div>
            <div class="text-gray-500">{{ service.users.toLocaleString() }} 人在用</div>
            <div class="text-green-600 font-medium">{{ service.price }}</div>
          </div>
        </div>
      </div>
    </div>

    <div class="grid grid-cols-3 gap-6">
      <div class="col-span-2 space-y-6">
        <div class="bg-white rounded-lg border border-gray-200 p-6">
          <h3 class="font-semibold text-gray-900 mb-4">服务介绍</h3>
          <p class="text-gray-600">{{ service.description }}</p>
          <div class="mt-4">
            <h4 class="text-sm font-medium text-gray-700 mb-2">主要功能：</h4>
            <div class="flex flex-wrap gap-2">
              <span v-for="feature in service.features" :key="feature"
                    class="px-3 py-1 bg-blue-50 text-blue-700 rounded-full text-sm">
                {{ feature }}
              </span>
            </div>
          </div>
        </div>

        <div class="bg-white rounded-lg border border-gray-200 p-6">
          <h3 class="font-semibold text-gray-900 mb-4">用户评价</h3>
          <div class="space-y-4">
            <div v-for="review in reviews" :key="review.id"
                 class="p-4 bg-gray-50 rounded-lg">
              <div class="flex items-center justify-between">
                <div class="flex items-center gap-2">
                  <span class="font-medium text-gray-900">{{ review.user }}</span>
                  <span class="text-yellow-500">{{ '⭐'.repeat(review.rating) }}</span>
                </div>
                <span class="text-sm text-gray-400">{{ review.date }}</span>
              </div>
              <p class="text-gray-600 mt-2">{{ review.content }}</p>
            </div>
          </div>
        </div>
      </div>

      <div class="space-y-4">
        <div class="bg-white rounded-lg border border-gray-200 p-4">
          <h4 class="font-medium text-gray-900 mb-3">服务信息</h4>
          <div class="space-y-3 text-sm">
            <div class="flex justify-between">
              <span class="text-gray-500">服务ID</span>
              <span class="text-gray-900">{{ service.id }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-500">服务商</span>
              <span class="text-gray-900">{{ service.provider }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-500">分类</span>
              <span class="text-gray-900">{{ service.category }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-500">价格</span>
              <span class="text-green-600 font-medium">{{ service.price }}</span>
            </div>
          </div>
        </div>

        <div class="bg-blue-50 border border-blue-200 rounded-lg p-4">
          <h4 class="font-medium text-blue-800 mb-2">使用须知</h4>
          <ul class="text-sm text-blue-700 space-y-1">
            <li>• 需要授权账户访问权限</li>
            <li>• 支持随时开启/关闭</li>
            <li>• 数据仅用于优化投放</li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</template>
