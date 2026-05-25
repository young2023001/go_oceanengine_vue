<script setup lang="ts">
import { ref, reactive } from 'vue'
import Breadcrumb from '@/components/common/Breadcrumb.vue'
import Pagination from '@/components/common/Pagination.vue'

const pagination = reactive({ page: 1, pageSize: 12, total: 86 })

const assets = ref([
  { id: 'AST001', name: '商品主图模板', type: 'image', format: 'PNG', size: '256KB', usedCount: 45, status: 'active', preview: '🖼️' },
  { id: 'AST002', name: '促销视频模板', type: 'video', format: 'MP4', size: '12MB', usedCount: 32, status: 'active', preview: '🎬' },
  { id: 'AST003', name: '轮播图组件', type: 'component', format: '-', size: '-', usedCount: 89, status: 'active', preview: '📱' },
  { id: 'AST004', name: '价格标签组件', type: 'component', format: '-', size: '-', usedCount: 156, status: 'active', preview: '🏷️' },
  { id: 'AST005', name: '倒计时组件', type: 'component', format: '-', size: '-', usedCount: 78, status: 'active', preview: '⏱️' },
  { id: 'AST006', name: '商品详情模板', type: 'template', format: 'HTML', size: '45KB', usedCount: 23, status: 'active', preview: '📋' }
])

const getTypeConfig = (type: string) => {
  switch (type) {
    case 'image': return { label: '图片', class: 'bg-blue-100 text-blue-700' }
    case 'video': return { label: '视频', class: 'bg-purple-100 text-purple-700' }
    case 'component': return { label: '组件', class: 'bg-green-100 text-green-700' }
    default: return { label: '模板', class: 'bg-orange-100 text-orange-700' }
  }
}

const handlePageChange = (page: number) => {
  pagination.page = page
}

const handleUploadAsset = () => {
  // TODO: implement
}

const handleViewAsset = (_asset: typeof assets.value[0]) => {
  // TODO: implement
}
</script>

<template>
  <div class="space-y-6">
    <div>
      <Breadcrumb :items="[{ name: 'DPA商品' }, { name: '素材库' }]" />
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-3xl font-bold text-gray-900">DPA素材库</h1>
          <p class="mt-2 text-gray-600">管理DPA广告素材和组件</p>
        </div>
        <button class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700" @click="handleUploadAsset">
          上传素材
        </button>
      </div>
    </div>

    <div class="grid grid-cols-1 sm:grid-cols-4 gap-4">
      <div class="bg-white rounded-lg border border-gray-200 p-4">
        <p class="text-sm text-gray-500">总素材</p>
        <p class="text-2xl font-bold text-gray-900 mt-1">{{ pagination.total }}</p>
      </div>
      <div class="bg-white rounded-lg border border-gray-200 p-4">
        <p class="text-sm text-gray-500">图片/视频</p>
        <p class="text-2xl font-bold text-blue-600 mt-1">45</p>
      </div>
      <div class="bg-white rounded-lg border border-gray-200 p-4">
        <p class="text-sm text-gray-500">组件</p>
        <p class="text-2xl font-bold text-green-600 mt-1">28</p>
      </div>
      <div class="bg-white rounded-lg border border-gray-200 p-4">
        <p class="text-sm text-gray-500">模板</p>
        <p class="text-2xl font-bold text-purple-600 mt-1">13</p>
      </div>
    </div>

    <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-6 gap-4">
      <div v-for="asset in assets" :key="asset.id"
           class="bg-white rounded-lg border border-gray-200 overflow-hidden hover:shadow-lg transition-shadow cursor-pointer group"
           @click="handleViewAsset(asset)">
        <div class="aspect-square bg-gray-100 flex items-center justify-center text-5xl relative">
          {{ asset.preview }}
          <span :class="['absolute top-2 left-2 px-2 py-0.5 rounded text-xs', getTypeConfig(asset.type).class]">
            {{ getTypeConfig(asset.type).label }}
          </span>
        </div>
        <div class="p-3">
          <h4 class="font-medium text-sm text-gray-900 truncate">{{ asset.name }}</h4>
          <div class="flex items-center justify-between mt-2 text-xs text-gray-500">
            <span>{{ asset.size !== '-' ? asset.size : asset.format }}</span>
            <span>{{ asset.usedCount }}次使用</span>
          </div>
        </div>
      </div>
    </div>

    <div class="flex justify-center">
      <Pagination :current="pagination.page" :total="pagination.total" :page-size="pagination.pageSize" @change="handlePageChange" />
    </div>
  </div>
</template>
