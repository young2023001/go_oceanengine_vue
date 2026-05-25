<script setup lang="ts">
import { ref, reactive } from 'vue'
import Breadcrumb from '@/components/common/Breadcrumb.vue'
import Pagination from '@/components/common/Pagination.vue'

const pagination = reactive({ page: 1, pageSize: 12, total: 36 })
const filterCategory = ref('')
const filterDuration = ref('')
const searchKeyword = ref('')

const templates = ref([
  { id: 'VT001', name: '商品展示模板', category: '电商', duration: '15s', resolution: '1080x1920', usedCount: 256, rating: 4.8, preview: '🛍️' },
  { id: 'VT002', name: '促销倒计时', category: '促销', duration: '10s', resolution: '1080x1920', usedCount: 189, rating: 4.7, preview: '⏰' },
  { id: 'VT003', name: '产品对比', category: '电商', duration: '20s', resolution: '1080x1920', usedCount: 145, rating: 4.6, preview: '⚖️' },
  { id: 'VT004', name: '用户评价', category: '口碑', duration: '15s', resolution: '1080x1920', usedCount: 98, rating: 4.5, preview: '⭐' },
  { id: 'VT005', name: '开箱体验', category: '电商', duration: '30s', resolution: '1080x1920', usedCount: 76, rating: 4.9, preview: '📦' },
  { id: 'VT006', name: '节日促销', category: '促销', duration: '15s', resolution: '1080x1920', usedCount: 312, rating: 4.8, preview: '🎉' }
])

const handlePageChange = (page: number) => {
  pagination.page = page
}

const handleUploadTemplate = () => {
  // TODO: implement
}

const handleUseTemplate = (_tpl: typeof templates.value[0]) => {
  // TODO: implement
}

const handlePreviewTemplate = (_tpl: typeof templates.value[0]) => {
  // TODO: implement
}
</script>

<template>
  <div class="space-y-6">
    <div>
      <Breadcrumb :items="[{ name: 'DPA商品' }, { name: '视频模板' }]" />
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-3xl font-bold text-gray-900">DPA视频模板</h1>
          <p class="mt-2 text-gray-600">选择模板快速生成商品视频广告</p>
        </div>
        <button class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700" @click="handleUploadTemplate">
          上传模板
        </button>
      </div>
    </div>

    <div class="bg-white rounded-lg border border-gray-200 p-4">
      <div class="flex flex-wrap gap-4">
        <select v-model="filterCategory" class="px-4 py-2 border border-gray-300 rounded-lg">
          <option value="">全部分类</option>
          <option value="ecommerce">电商</option>
          <option value="promotion">促销</option>
          <option value="review">口碑</option>
        </select>
        <select v-model="filterDuration" class="px-4 py-2 border border-gray-300 rounded-lg">
          <option value="">时长</option>
          <option value="10">10秒</option>
          <option value="15">15秒</option>
          <option value="30">30秒</option>
        </select>
        <input v-model="searchKeyword" type="text" placeholder="搜索模板..." class="flex-1 min-w-[200px] px-4 py-2 border border-gray-300 rounded-lg">
      </div>
    </div>

    <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-6 gap-4">
      <div v-for="tpl in templates" :key="tpl.id"
           class="bg-white rounded-lg border border-gray-200 overflow-hidden hover:shadow-lg transition-shadow group cursor-pointer"
           @click="handlePreviewTemplate(tpl)">
        <div class="aspect-[9/16] bg-gradient-to-br from-purple-100 to-blue-100 flex items-center justify-center text-6xl relative">
          {{ tpl.preview }}
          <div class="absolute inset-0 bg-black/50 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center">
            <button class="px-4 py-2 bg-white text-gray-900 rounded-lg text-sm font-medium" @click.stop="handleUseTemplate(tpl)">
              使用模板
            </button>
          </div>
          <span class="absolute bottom-2 right-2 px-2 py-0.5 bg-black/60 text-white text-xs rounded">{{ tpl.duration }}</span>
        </div>
        <div class="p-3">
          <h4 class="font-medium text-sm text-gray-900 truncate">{{ tpl.name }}</h4>
          <div class="flex items-center justify-between mt-2">
            <span class="px-2 py-0.5 bg-blue-100 text-blue-700 rounded text-xs">{{ tpl.category }}</span>
            <div class="flex items-center gap-1">
              <span class="text-yellow-400 text-xs">★</span>
              <span class="text-xs text-gray-600">{{ tpl.rating }}</span>
            </div>
          </div>
          <p class="text-xs text-gray-500 mt-1">{{ tpl.usedCount }}次使用</p>
        </div>
      </div>
    </div>

    <div class="flex justify-center">
      <Pagination :current="pagination.page" :total="pagination.total" :page-size="pagination.pageSize" @change="handlePageChange" />
    </div>
  </div>
</template>
