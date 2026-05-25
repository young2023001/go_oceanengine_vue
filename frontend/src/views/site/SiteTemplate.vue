<script setup lang="ts">
import { ref } from 'vue'
import Breadcrumb from '@/components/common/Breadcrumb.vue'

const activeCategory = ref('all')
const categories = [
  { id: 'all', name: '全部', count: 68 },
  { id: 'ecommerce', name: '电商', count: 25 },
  { id: 'lead', name: '线索', count: 18 },
  { id: 'app', name: '应用', count: 15 },
  { id: 'brand', name: '品牌', count: 10 }
]

const templates = ref([
  { id: 'T001', name: '电商促销模板', category: 'ecommerce', uses: 2560, preview: '🛒' },
  { id: 'T002', name: '线索收集模板', category: 'lead', uses: 1890, preview: '📋' },
  { id: 'T003', name: '应用下载模板', category: 'app', uses: 1256, preview: '📱' },
  { id: 'T004', name: '品牌展示模板', category: 'brand', uses: 856, preview: '✨' },
  { id: 'T005', name: '活动推广模板', category: 'ecommerce', uses: 768, preview: '🎉' },
  { id: 'T006', name: '产品详情模板', category: 'ecommerce', uses: 656, preview: '📦' }
])

const handleUseTemplate = (_template: any) => {
  // TODO: implement
}

const handleLoadMore = () => {
  // TODO: implement
}

const handlePreviewTemplate = (_template: any) => {
  // TODO: implement
}
</script>

<template>
  <div class="space-y-6">
    <div>
      <Breadcrumb :items="[{ name: '落地页' }, { name: '页面模板' }]" />
      <h1 class="text-3xl font-bold text-gray-900">落地页模板</h1>
      <p class="mt-2 text-gray-600">使用模板快速创建落地页</p>
    </div>

    <div class="flex gap-2 overflow-x-auto pb-2">
      <button v-for="cat in categories" :key="cat.id"
              :class="['px-4 py-2 rounded-full text-sm whitespace-nowrap transition-colors',
                       activeCategory === cat.id ? 'bg-blue-600 text-white' : 'bg-gray-100 text-gray-700 hover:bg-gray-200']"
              @click="activeCategory = cat.id">
        {{ cat.name }} ({{ cat.count }})
      </button>
    </div>

    <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
      <div v-for="template in templates" :key="template.id"
           class="bg-white rounded-lg border border-gray-200 overflow-hidden hover:shadow-lg transition-shadow cursor-pointer"
           @click="handlePreviewTemplate(template)">
        <div class="aspect-[9/16] bg-gradient-to-br from-gray-100 to-gray-200 flex items-center justify-center">
          <span class="text-6xl">{{ template.preview }}</span>
        </div>
        <div class="p-4">
          <span class="px-2 py-0.5 bg-blue-100 text-blue-700 rounded text-xs mb-2 inline-block">{{ template.category }}</span>
          <h4 class="font-medium text-gray-900">{{ template.name }}</h4>
          <div class="flex items-center justify-between mt-3">
            <span class="text-xs text-gray-500">{{ template.uses.toLocaleString() }} 次使用</span>
<button class="px-3 py-1 text-sm text-blue-600 border border-blue-300 rounded hover:bg-blue-50" @click.stop="handleUseTemplate(template)">使用</button>
          </div>
        </div>
      </div>
    </div>

<div class="flex justify-center">
      <button class="px-6 py-2 border border-gray-300 rounded-lg hover:bg-gray-50" @click="handleLoadMore">
        加载更多
      </button>
    </div>
  </div>
</template>
