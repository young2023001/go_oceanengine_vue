<script setup lang="ts">
import { ref } from 'vue'
import Breadcrumb from '@/components/common/Breadcrumb.vue'

const expandedCategories = ref<string[]>(['1', '3'])

const categories = ref([
  { id: '1', name: '数码3C', icon: '📱', count: 1256, children: [
    { id: '101', name: '手机', count: 456 },
    { id: '102', name: '电脑', count: 328 },
    { id: '103', name: '平板', count: 245 },
    { id: '104', name: '智能穿戴', count: 227 }
  ]},
  { id: '2', name: '服装鞋包', icon: '👔', count: 2890, children: [
    { id: '201', name: '女装', count: 1256 },
    { id: '202', name: '男装', count: 892 },
    { id: '203', name: '鞋靴', count: 456 },
    { id: '204', name: '箱包', count: 286 }
  ]},
  { id: '3', name: '美妆个护', icon: '💄', count: 1680, children: [
    { id: '301', name: '护肤', count: 756 },
    { id: '302', name: '彩妆', count: 524 },
    { id: '303', name: '香水', count: 200 },
    { id: '304', name: '个护', count: 200 }
  ]},
  { id: '4', name: '家居家装', icon: '🏠', count: 980, children: [
    { id: '401', name: '家具', count: 456 },
    { id: '402', name: '家纺', count: 324 },
    { id: '403', name: '灯具', count: 200 }
  ]},
  { id: '5', name: '食品饮料', icon: '🍎', count: 756, children: [
    { id: '501', name: '零食', count: 356 },
    { id: '502', name: '饮料', count: 200 },
    { id: '503', name: '生鲜', count: 200 }
  ]}
])

const toggleExpand = (id: string) => {
  const index = expandedCategories.value.indexOf(id)
  if (index > -1) {
    expandedCategories.value.splice(index, 1)
  } else {
    expandedCategories.value.push(id)
  }
}

const handleCreateCategory = () => {
  // TODO: implement
}

const handleEditCategory = (_childId: string, _childName: string) => {
  // TODO: implement
}
</script>

<template>
  <div class="space-y-6">
    <div>
      <Breadcrumb :items="[{ name: 'DPA商品' }, { name: '商品分类' }]" />
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-3xl font-bold text-gray-900">商品分类管理</h1>
          <p class="mt-2 text-gray-600">管理DPA商品库分类目录</p>
        </div>
        <button class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700" @click="handleCreateCategory">
          新建分类
        </button>
      </div>
    </div>

    <div class="grid grid-cols-1 sm:grid-cols-5 gap-4">
      <div class="bg-white rounded-lg border border-gray-200 p-4">
        <p class="text-sm text-gray-500">一级分类</p>
        <p class="text-2xl font-bold text-gray-900 mt-1">{{ categories.length }}</p>
      </div>
      <div class="bg-white rounded-lg border border-gray-200 p-4">
        <p class="text-sm text-gray-500">二级分类</p>
        <p class="text-2xl font-bold text-blue-600 mt-1">{{ categories.reduce((acc, c) => acc + c.children.length, 0) }}</p>
      </div>
      <div class="bg-white rounded-lg border border-gray-200 p-4">
        <p class="text-sm text-gray-500">总商品数</p>
        <p class="text-2xl font-bold text-green-600 mt-1">{{ categories.reduce((acc, c) => acc + c.count, 0).toLocaleString() }}</p>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <div v-for="category in categories" :key="category.id"
           class="bg-white rounded-lg border border-gray-200 overflow-hidden">
        <div class="p-4 bg-gray-50 border-b border-gray-200 cursor-pointer hover:bg-gray-100"
             @click="toggleExpand(category.id)">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-3">
              <span class="text-2xl">{{ category.icon }}</span>
              <div>
                <h3 class="font-semibold text-gray-900">{{ category.name }}</h3>
                <p class="text-xs text-gray-500">{{ category.count }} 件商品</p>
              </div>
            </div>
            <div class="flex items-center gap-4">
              <span class="text-sm text-gray-500">{{ category.children.length }} 个子类</span>
              <svg :class="['w-5 h-5 text-gray-400 transition-transform',
                           expandedCategories.includes(category.id) ? 'rotate-180' : '']"
                   fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"/>
              </svg>
            </div>
          </div>
        </div>
        <div v-if="expandedCategories.includes(category.id)" class="divide-y divide-gray-100">
          <div v-for="child in category.children" :key="child.id"
               class="px-4 py-3 flex items-center justify-between hover:bg-gray-50">
            <div class="flex items-center gap-3">
              <span class="w-6 h-6 bg-gray-100 rounded flex items-center justify-center text-xs text-gray-500">{{ child.id.slice(-2) }}</span>
              <span class="text-sm text-gray-700">{{ child.name }}</span>
            </div>
            <div class="flex items-center gap-4">
              <span class="text-sm text-gray-500">{{ child.count }} 件</span>
              <button class="text-blue-600 text-sm hover:text-blue-800" @click="handleEditCategory(child.id, child.name)">编辑</button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
