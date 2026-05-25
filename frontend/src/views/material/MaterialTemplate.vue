<script setup lang="ts">
import { ref, reactive } from 'vue'
import Breadcrumb from '@/components/common/Breadcrumb.vue'
import Pagination from '@/components/common/Pagination.vue'

const pagination = reactive({ page: 1, pageSize: 12, total: 86 })
const activeCategory = ref('all')

const categories = [
  { key: 'all', label: '全部' },
  { key: 'ecommerce', label: '电商促销' },
  { key: 'brand', label: '品牌宣传' },
  { key: 'app', label: '应用推广' },
  { key: 'game', label: '游戏推广' }
]

const templates = ref([
  { id: 'TPL001', name: '电商爆款模板', category: 'ecommerce', preview: '🛍️', type: 'video', duration: '15s', usedCount: 1256, rating: 4.8 },
  { id: 'TPL002', name: '品牌故事模板', category: 'brand', preview: '✨', type: 'video', duration: '30s', usedCount: 892, rating: 4.7 },
  { id: 'TPL003', name: 'APP下载模板', category: 'app', preview: '📲', type: 'video', duration: '15s', usedCount: 2350, rating: 4.9 },
  { id: 'TPL004', name: '游戏预告模板', category: 'game', preview: '🎮', type: 'video', duration: '20s', usedCount: 756, rating: 4.6 },
  { id: 'TPL005', name: '促销海报模板', category: 'ecommerce', preview: '🎯', type: 'image', duration: '-', usedCount: 3450, rating: 4.8 },
  { id: 'TPL006', name: '产品展示模板', category: 'brand', preview: '📦', type: 'video', duration: '10s', usedCount: 1890, rating: 4.5 }
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
</script>

<template>
  <div class="space-y-6">
    <div>
      <Breadcrumb :items="[{ name: '素材管理' }, { name: '素材模板' }]" />
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-3xl font-bold text-gray-900">素材模板</h1>
          <p class="mt-2 text-gray-600">选择模板快速制作广告素材</p>
        </div>
        <button class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700" @click="handleUploadTemplate">
          上传模板
        </button>
      </div>
    </div>

    <div class="bg-white rounded-lg border border-gray-200 p-4">
      <div class="flex flex-wrap gap-4">
        <button v-for="cat in categories" :key="cat.key"
                :class="['px-4 py-2 rounded-lg text-sm transition-colors',
                  activeCategory === cat.key ? 'bg-blue-600 text-white' : 'bg-gray-100 text-gray-700 hover:bg-gray-200']"
                @click="activeCategory = cat.key">
          {{ cat.label }}
        </button>
      </div>
    </div>

    <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-6 gap-4">
      <div v-for="tpl in templates" :key="tpl.id"
           class="bg-white rounded-lg border border-gray-200 overflow-hidden hover:shadow-lg transition-shadow group cursor-pointer">
        <div class="aspect-[9/16] bg-gradient-to-br from-blue-50 to-purple-50 flex items-center justify-center text-6xl relative">
          {{ tpl.preview }}
          <div class="absolute inset-0 bg-black/50 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center">
            <button class="px-4 py-2 bg-white text-gray-900 rounded-lg text-sm font-medium" @click.stop="handleUseTemplate(tpl)">
              使用模板
            </button>
          </div>
          <span v-if="tpl.duration !== '-'" class="absolute bottom-2 right-2 px-2 py-0.5 bg-black/60 text-white text-xs rounded">
            {{ tpl.duration }}
          </span>
        </div>
        <div class="p-3">
          <h4 class="font-medium text-sm text-gray-900 truncate">{{ tpl.name }}</h4>
          <div class="flex items-center justify-between mt-2">
            <div class="flex items-center gap-1">
              <span class="text-yellow-400 text-xs">★</span>
              <span class="text-xs text-gray-600">{{ tpl.rating }}</span>
            </div>
            <span class="text-xs text-gray-500">{{ tpl.usedCount }}次使用</span>
          </div>
        </div>
      </div>
    </div>

    <div class="flex justify-center">
      <Pagination :current="pagination.page" :total="pagination.total" :page-size="pagination.pageSize" @change="handlePageChange" />
    </div>
  </div>
</template>
