<script setup lang="ts">
import { ref, reactive } from 'vue'
import Breadcrumb from '@/components/common/Breadcrumb.vue'
import Pagination from '@/components/common/Pagination.vue'

const pagination = reactive({ page: 1, pageSize: 20, total: 245 })
const selectedImages = ref<string[]>([])

const searchKeyword = ref('')
const filterStatus = ref('')
const filterFormat = ref('')

const images = ref([
  { id: 'IMG001', name: '产品主图', thumbnail: '🖼️', size: '1.2MB', dimensions: '800x800', format: 'JPG', status: 'active', usedCount: 45, createdAt: '2025-10-01' },
  { id: 'IMG002', name: '促销Banner', thumbnail: '🎨', size: '856KB', dimensions: '1200x628', format: 'PNG', status: 'active', usedCount: 32, createdAt: '2025-10-03' },
  { id: 'IMG003', name: '品牌Logo', thumbnail: '✨', size: '128KB', dimensions: '500x500', format: 'PNG', status: 'active', usedCount: 89, createdAt: '2025-10-05' },
  { id: 'IMG004', name: '活动海报', thumbnail: '📸', size: '2.1MB', dimensions: '1080x1920', format: 'JPG', status: 'reviewing', usedCount: 0, createdAt: '2025-10-08' },
  { id: 'IMG005', name: '用户晒单', thumbnail: '🏞️', size: '1.5MB', dimensions: '750x1000', format: 'JPG', status: 'active', usedCount: 18, createdAt: '2025-10-10' },
  { id: 'IMG006', name: '功能介绍图', thumbnail: '📊', size: '680KB', dimensions: '1200x800', format: 'PNG', status: 'rejected', usedCount: 0, createdAt: '2025-10-12' },
  { id: 'IMG007', name: '节日素材', thumbnail: '🎄', size: '920KB', dimensions: '800x800', format: 'PNG', status: 'active', usedCount: 56, createdAt: '2025-10-15' },
  { id: 'IMG008', name: '背景图', thumbnail: '🌄', size: '1.8MB', dimensions: '1920x1080', format: 'JPG', status: 'active', usedCount: 23, createdAt: '2025-10-18' }
])

const toggleSelect = (id: string) => {
  const idx = selectedImages.value.indexOf(id)
  if (idx > -1) {
    selectedImages.value.splice(idx, 1)
  } else {
    selectedImages.value.push(id)
  }
}

const getStatusConfig = (status: string) => {
  switch (status) {
    case 'active': return { label: '可用', class: 'bg-green-100 text-green-800' }
    case 'reviewing': return { label: '审核中', class: 'bg-yellow-100 text-yellow-800' }
    case 'rejected': return { label: '已拒绝', class: 'bg-red-100 text-red-800' }
    default: return { label: status, class: 'bg-gray-100 text-gray-800' }
  }
}

const handlePageChange = (page: number) => {
  pagination.page = page
}

const handleBatchDelete = () => {
  if (confirm(`确定删除选中的 ${selectedImages.value.length} 张图片?`)) {
    selectedImages.value = []
  }
}
</script>

<template>
  <div class="space-y-6">
    <div>
      <Breadcrumb :items="[{ name: '素材管理' }, { name: '图片素材' }]" />
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-3xl font-bold text-gray-900">图片素材</h1>
          <p class="mt-2 text-gray-600">管理广告投放使用的图片素材</p>
        </div>
        <div class="flex items-center gap-4">
          <span v-if="selectedImages.length > 0" class="text-sm text-gray-600">
            已选择 {{ selectedImages.length }} 张
          </span>
          <button v-if="selectedImages.length > 0" class="px-4 py-2 border border-red-300 text-red-600 rounded-lg hover:bg-red-50" @click="handleBatchDelete">
            批量删除
          </button>
          <router-link to="/material/image/upload" class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700">
            上传图片
          </router-link>
        </div>
      </div>
    </div>

    <div class="grid grid-cols-1 sm:grid-cols-4 gap-4">
      <div class="bg-white rounded-lg border border-gray-200 p-4">
        <p class="text-sm text-gray-500">总图片</p>
        <p class="text-2xl font-bold text-gray-900 mt-1">{{ pagination.total }}</p>
      </div>
      <div class="bg-white rounded-lg border border-gray-200 p-4">
        <p class="text-sm text-gray-500">可用</p>
        <p class="text-2xl font-bold text-green-600 mt-1">218</p>
      </div>
      <div class="bg-white rounded-lg border border-gray-200 p-4">
        <p class="text-sm text-gray-500">存储占用</p>
        <p class="text-2xl font-bold text-blue-600 mt-1">1.2GB</p>
      </div>
      <div class="bg-white rounded-lg border border-gray-200 p-4">
        <p class="text-sm text-gray-500">总使用次数</p>
        <p class="text-2xl font-bold text-purple-600 mt-1">2,456</p>
      </div>
    </div>

    <div class="bg-white rounded-lg border border-gray-200 p-4">
      <div class="flex flex-wrap gap-4">
        <input v-model="searchKeyword" type="text" placeholder="搜索图片名称..."
               class="flex-1 min-w-[200px] px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500">
        <select v-model="filterStatus" class="px-4 py-2 border border-gray-300 rounded-lg">
          <option value="">全部状态</option>
          <option value="active">可用</option>
          <option value="reviewing">审核中</option>
          <option value="rejected">已拒绝</option>
        </select>
        <select v-model="filterFormat" class="px-4 py-2 border border-gray-300 rounded-lg">
          <option value="">全部格式</option>
          <option value="jpg">JPG</option>
          <option value="png">PNG</option>
          <option value="gif">GIF</option>
        </select>
      </div>
    </div>

    <div class="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-8 gap-4">
      <div v-for="img in images" :key="img.id" 
           :class="['bg-white rounded-lg border-2 overflow-hidden cursor-pointer transition-all group',
             selectedImages.includes(img.id) ? 'border-blue-500' : 'border-gray-200 hover:border-gray-300']"
           @click="toggleSelect(img.id)">
        <div class="aspect-square bg-gray-100 flex items-center justify-center text-4xl relative">
          {{ img.thumbnail }}
          <div v-if="selectedImages.includes(img.id)" 
               class="absolute top-2 right-2 w-6 h-6 bg-blue-500 rounded-full flex items-center justify-center">
            <svg class="w-4 h-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
            </svg>
          </div>
          <span :class="['absolute bottom-2 left-2 px-1.5 py-0.5 rounded text-xs', getStatusConfig(img.status).class]">
            {{ getStatusConfig(img.status).label }}
          </span>
        </div>
        <div class="p-2">
          <h4 class="text-xs font-medium text-gray-900 truncate">{{ img.name }}</h4>
          <p class="text-xs text-gray-500 mt-1">{{ img.dimensions }}</p>
        </div>
      </div>
    </div>

    <div class="flex justify-center">
      <Pagination :current="pagination.page" :total="pagination.total" :page-size="pagination.pageSize" @change="handlePageChange" />
    </div>
  </div>
</template>
