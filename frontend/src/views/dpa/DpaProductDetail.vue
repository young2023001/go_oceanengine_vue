<script setup lang="ts">
import { ref } from 'vue'
import { useRoute } from 'vue-router'
import Breadcrumb from '@/components/common/Breadcrumb.vue'

const route = useRoute()
const productId = route.params.id || 'DPA001'

const product = ref({
  id: productId,
  name: '智能手表Pro Max',
  category: '数码3C/智能穿戴',
  price: 1299,
  originalPrice: 1599,
  status: 'active',
  stock: 1520,
  salesCount: 3256,
  images: ['📱', '⌚', '🎁'],
  description: '全新一代智能手表，支持血氧检测、心率监测、睡眠分析等功能...',
  attributes: [
    { name: '品牌', value: 'TechPro' },
    { name: '型号', value: 'PM-2025' },
    { name: '颜色', value: '深空灰/银白/玫瑰金' },
    { name: '屏幕', value: '1.75英寸 AMOLED' },
    { name: '电池', value: '410mAh' }
  ],
  createdAt: '2025-10-01',
  updatedAt: '2025-11-20'
})

const adStats = ref({
  impressions: 125600,
  clicks: 4523,
  ctr: 3.6,
  conversions: 328,
  cost: 15680
})

const handleCreateAd = () => {
  // TODO: implement
}
</script>

<template>
  <div class="space-y-6">
    <div>
      <Breadcrumb :items="[{ name: 'DPA商品' }, { name: '商品详情' }]" />
      <div class="flex items-center justify-between">
        <h1 class="text-3xl font-bold text-gray-900">商品详情</h1>
        <div class="flex items-center gap-4">
          <router-link :to="`/dpa/product/edit/${product.id}`" class="px-4 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50">
            编辑商品
          </router-link>
          <button class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700" @click="handleCreateAd">
            创建广告
          </button>
        </div>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <div class="lg:col-span-2 space-y-6">
        <div class="bg-white rounded-lg border border-gray-200 p-6">
          <div class="flex gap-6">
            <div class="w-48 h-48 bg-gray-100 rounded-lg flex items-center justify-center text-7xl">
              {{ product.images[0] }}
            </div>
            <div class="flex-1">
              <div class="flex items-start justify-between">
                <div>
                  <h2 class="text-xl font-bold text-gray-900">{{ product.name }}</h2>
                  <p class="text-sm text-gray-500 mt-1">{{ product.category }}</p>
                </div>
                <span :class="['px-3 py-1 rounded-full text-sm font-medium',
                  product.status === 'active' ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800']">
                  {{ product.status === 'active' ? '在售' : '下架' }}
                </span>
              </div>
              <div class="mt-4">
                <span class="text-3xl font-bold text-red-600">¥{{ product.price }}</span>
                <span class="text-sm text-gray-400 line-through ml-2">¥{{ product.originalPrice }}</span>
              </div>
              <div class="flex items-center gap-6 mt-4 text-sm text-gray-600">
                <span>库存: {{ product.stock }}</span>
                <span>销量: {{ product.salesCount }}</span>
                <span>ID: {{ product.id }}</span>
              </div>
            </div>
          </div>
        </div>

        <div class="bg-white rounded-lg border border-gray-200 p-6">
          <h3 class="font-semibold text-gray-900 mb-4">商品属性</h3>
          <div class="grid grid-cols-2 gap-4">
            <div v-for="attr in product.attributes" :key="attr.name" 
                 class="flex items-center gap-2 py-2 border-b border-gray-100">
              <span class="text-sm text-gray-500 w-20">{{ attr.name }}</span>
              <span class="text-sm text-gray-900">{{ attr.value }}</span>
            </div>
          </div>
        </div>

        <div class="bg-white rounded-lg border border-gray-200 p-6">
          <h3 class="font-semibold text-gray-900 mb-4">商品描述</h3>
          <p class="text-gray-600">{{ product.description }}</p>
        </div>
      </div>

      <div class="space-y-6">
        <div class="bg-white rounded-lg border border-gray-200 p-6">
          <h3 class="font-semibold text-gray-900 mb-4">投放数据</h3>
          <div class="space-y-4">
            <div class="flex justify-between items-center py-2 border-b border-gray-100">
              <span class="text-gray-500">展示次数</span>
              <span class="font-semibold text-gray-900">{{ adStats.impressions.toLocaleString() }}</span>
            </div>
            <div class="flex justify-between items-center py-2 border-b border-gray-100">
              <span class="text-gray-500">点击次数</span>
              <span class="font-semibold text-gray-900">{{ adStats.clicks.toLocaleString() }}</span>
            </div>
            <div class="flex justify-between items-center py-2 border-b border-gray-100">
              <span class="text-gray-500">点击率</span>
              <span class="font-semibold text-blue-600">{{ adStats.ctr }}%</span>
            </div>
            <div class="flex justify-between items-center py-2 border-b border-gray-100">
              <span class="text-gray-500">转化数</span>
              <span class="font-semibold text-green-600">{{ adStats.conversions }}</span>
            </div>
            <div class="flex justify-between items-center py-2">
              <span class="text-gray-500">消耗</span>
              <span class="font-semibold text-gray-900">¥{{ adStats.cost.toLocaleString() }}</span>
            </div>
          </div>
        </div>

        <div class="bg-white rounded-lg border border-gray-200 p-6">
          <h3 class="font-semibold text-gray-900 mb-4">更新记录</h3>
          <div class="space-y-2 text-sm">
            <div class="flex justify-between text-gray-600">
              <span>创建时间</span>
              <span>{{ product.createdAt }}</span>
            </div>
            <div class="flex justify-between text-gray-600">
              <span>最后更新</span>
              <span>{{ product.updatedAt }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
