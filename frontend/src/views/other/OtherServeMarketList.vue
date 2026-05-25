<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import Breadcrumb from '@/components/common/Breadcrumb.vue'

const router = useRouter()

const filterCategory = ref('')
const sortBy = ref('')
const searchKeyword = ref('')

const categories = ref([
  { id: 'creative', name: '创意服务', icon: '🎨', count: 45 },
  { id: 'data', name: '数据服务', icon: '📊', count: 32 },
  { id: 'tech', name: '技术服务', icon: '⚙️', count: 28 },
  { id: 'marketing', name: '营销服务', icon: '📢', count: 38 }
])

const services = ref([
  { id: 'SM001', name: '专业视频制作', provider: '优创科技', category: '创意服务', price: 2000, rating: 4.8, orders: 256 },
  { id: 'SM002', name: '数据分析报告', provider: '数据魔方', category: '数据服务', price: 500, rating: 4.6, orders: 189 },
  { id: 'SM003', name: 'API对接服务', provider: '技术先锋', category: '技术服务', price: 3000, rating: 4.9, orders: 78 },
  { id: 'SM004', name: '品牌策划', provider: '品牌大师', category: '营销服务', price: 5000, rating: 4.7, orders: 145 }
])

const handleConsult = (service: typeof services.value[0]) => {
  router.push(`/other/servemarket/detail/${service.id}`)
}

const handleViewCategory = (_cat: typeof categories.value[0]) => {
  // TODO: implement
}
</script>

<template>
  <div class="space-y-6">
    <div>
      <Breadcrumb :items="[{ name: '其他' }, { name: '服务市场' }]" />
      <h1 class="text-3xl font-bold text-gray-900">服务市场</h1>
      <p class="mt-2 text-gray-600">发现优质的广告服务供应商</p>
    </div>

    <div class="grid grid-cols-4 gap-4">
      <div v-for="cat in categories" :key="cat.id"
           class="bg-white rounded-lg border border-gray-200 p-4 hover:shadow-md transition-shadow cursor-pointer"
           @click="handleViewCategory(cat)">
        <div class="flex items-center gap-3">
          <span class="text-3xl">{{ cat.icon }}</span>
          <div>
            <h4 class="font-medium text-gray-900">{{ cat.name }}</h4>
            <p class="text-sm text-gray-500">{{ cat.count }} 个服务</p>
          </div>
        </div>
      </div>
    </div>

    <div class="bg-white rounded-lg border border-gray-200 p-4">
      <div class="flex gap-4">
        <select v-model="filterCategory" class="px-4 py-2 border border-gray-300 rounded-lg">
          <option value="">全部分类</option>
          <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
        </select>
        <select v-model="sortBy" class="px-4 py-2 border border-gray-300 rounded-lg">
          <option value="">按评分排序</option>
          <option value="rating_desc">评分从高到低</option>
          <option value="price_asc">价格从低到高</option>
          <option value="orders_desc">订单量从高到低</option>
        </select>
        <input v-model="searchKeyword" type="text" placeholder="搜索服务..." class="flex-1 px-4 py-2 border border-gray-300 rounded-lg">
      </div>
    </div>

    <div class="grid grid-cols-2 gap-4">
      <div v-for="service in services" :key="service.id"
           class="bg-white rounded-lg border border-gray-200 p-6 hover:shadow-md transition-shadow">
        <div class="flex items-start justify-between">
          <div>
            <h4 class="font-semibold text-gray-900 text-lg">{{ service.name }}</h4>
            <p class="text-sm text-gray-500 mt-1">{{ service.provider }}</p>
          </div>
          <span class="px-2 py-1 bg-blue-100 text-blue-700 rounded text-xs">{{ service.category }}</span>
        </div>
        <div class="mt-4 flex items-center gap-4">
          <div class="flex items-center gap-1">
            <span class="text-yellow-400">★</span>
            <span class="text-sm font-medium">{{ service.rating }}</span>
          </div>
          <span class="text-sm text-gray-500">{{ service.orders }} 单</span>
        </div>
        <div class="mt-4 flex items-center justify-between">
          <p class="text-xl font-bold text-gray-900">¥{{ service.price.toLocaleString() }}<span class="text-sm font-normal text-gray-500">/次</span></p>
          <button class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 text-sm" @click="handleConsult(service)">
            立即咨询
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
