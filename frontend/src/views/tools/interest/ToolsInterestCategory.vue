<script setup lang="ts">
import { ref, computed } from 'vue'
import Breadcrumb from '@/components/common/Breadcrumb.vue'

const searchKeyword = ref('')
const expandedCategories = ref<string[]>(['1', '2'])

const categories = ref([
  { id: '1', name: '汽车', icon: '🚗', children: [
    { id: '101', name: '新能源汽车' },
    { id: '102', name: '豪华车' },
    { id: '103', name: 'SUV' },
    { id: '104', name: '汽车改装' }
  ]},
  { id: '2', name: '数码科技', icon: '📱', children: [
    { id: '201', name: '手机' },
    { id: '202', name: '电脑' },
    { id: '203', name: '智能穿戴' },
    { id: '204', name: '摄影器材' }
  ]},
  { id: '3', name: '美妆护肤', icon: '💄', children: [
    { id: '301', name: '彩妆' },
    { id: '302', name: '护肤' },
    { id: '303', name: '香水' }
  ]},
  { id: '4', name: '母婴育儿', icon: '👶', children: [
    { id: '401', name: '孕产' },
    { id: '402', name: '婴儿用品' },
    { id: '403', name: '早教' }
  ]},
  { id: '5', name: '旅游出行', icon: '✈️', children: [
    { id: '501', name: '国内游' },
    { id: '502', name: '出境游' },
    { id: '503', name: '酒店住宿' }
  ]}
])

const filteredCategories = computed(() => {
  if (!searchKeyword.value) return categories.value
  const keyword = searchKeyword.value.toLowerCase()
  return categories.value.filter(cat =>
    cat.name.toLowerCase().includes(keyword) ||
    cat.children.some(c => c.name.toLowerCase().includes(keyword))
  )
})

const toggleExpand = (id: string) => {
  const index = expandedCategories.value.indexOf(id)
  if (index > -1) {
    expandedCategories.value.splice(index, 1)
  } else {
    expandedCategories.value.push(id)
  }
}

const handleCopyId = (id: string) => {
  navigator.clipboard.writeText(id)
}
</script>

<template>
  <div class="space-y-6">
    <div>
      <Breadcrumb :items="[{ name: '兴趣定向' }, { name: '兴趣类目' }]" />
      <h1 class="text-3xl font-bold text-gray-900">兴趣类目查询</h1>
      <p class="mt-2 text-gray-600">查询可用的兴趣定向类目</p>
    </div>

    <div class="bg-white rounded-lg border border-gray-200 p-4">
      <input v-model="searchKeyword" type="text" placeholder="搜索兴趣类目..."
             class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500">
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      <div v-for="category in filteredCategories" :key="category.id"
           class="bg-white rounded-lg border border-gray-200 overflow-hidden">
        <div class="p-4 border-b border-gray-100 bg-gray-50 cursor-pointer hover:bg-gray-100"
             @click="toggleExpand(category.id)">
          <div class="flex items-center gap-3">
            <span class="text-2xl">{{ category.icon }}</span>
            <div class="flex-1">
              <h3 class="font-semibold text-gray-900">{{ category.name }}</h3>
              <p class="text-xs text-gray-500">{{ category.children.length }} 个子类目</p>
            </div>
            <svg :class="['w-5 h-5 text-gray-400 transition-transform',
                         expandedCategories.includes(category.id) ? 'rotate-180' : '']"
                 fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"/>
            </svg>
          </div>
        </div>
        <div v-if="expandedCategories.includes(category.id)" class="p-2">
          <div v-for="child in category.children" :key="child.id"
               class="flex items-center justify-between px-3 py-2 rounded hover:bg-gray-50 cursor-pointer group">
            <span class="text-sm text-gray-700">{{ child.name }}</span>
            <div class="flex items-center gap-2">
              <span class="text-xs text-gray-400">{{ child.id }}</span>
<button @click.stop="handleCopyId(child.id)" class="opacity-0 group-hover:opacity-100 text-blue-500 text-xs">复制ID</button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
