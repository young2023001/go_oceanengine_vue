<script setup lang="ts">
import { ref, reactive } from 'vue'
import Breadcrumb from '@/components/common/Breadcrumb.vue'
import Pagination from '@/components/common/Pagination.vue'

const pagination = reactive({ page: 1, pageSize: 10, total: 25 })

const alerts = ref([
  { id: 'AL001', title: '预算即将耗尽', desc: '账户预算剩余不足10%', level: 'high', target: '主账户', time: '2025-11-28 10:00', status: 'unread' },
  { id: 'AL002', title: '广告效果下降', desc: 'CTR较昨日下降30%', level: 'medium', target: '促销广告', time: '2025-11-28 09:30', status: 'unread' },
  { id: 'AL003', title: '创意审核失败', desc: '2个创意未通过审核', level: 'low', target: '新品素材', time: '2025-11-28 08:45', status: 'read' }
])

const handlePageChange = (page: number) => {
  pagination.page = page
}

const handleMarkAllRead = () => {
  // TODO: implement
}

const handleProcessAlert = (_alert: any) => {
  // TODO: implement
}
</script>

<template>
  <div class="space-y-6">
    <div>
      <Breadcrumb :items="[{ name: '工作台' }, { name: '预警中心' }]" />
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-3xl font-bold text-gray-900">预警中心</h1>
          <p class="mt-2 text-gray-600">查看系统预警和通知</p>
        </div>
<button class="px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50" @click="handleMarkAllRead">
          全部已读
        </button>
      </div>
    </div>

    <div class="grid grid-cols-4 gap-4">
      <div class="bg-white rounded-lg border border-gray-200 p-4">
        <p class="text-sm text-gray-500">预警总数</p>
        <p class="text-2xl font-bold text-gray-900 mt-1">{{ pagination.total }}</p>
      </div>
      <div class="bg-white rounded-lg border border-gray-200 p-4">
        <p class="text-sm text-gray-500">紧急</p>
        <p class="text-2xl font-bold text-red-600 mt-1">3</p>
      </div>
      <div class="bg-white rounded-lg border border-gray-200 p-4">
        <p class="text-sm text-gray-500">重要</p>
        <p class="text-2xl font-bold text-yellow-600 mt-1">8</p>
      </div>
      <div class="bg-white rounded-lg border border-gray-200 p-4">
        <p class="text-sm text-gray-500">未读</p>
        <p class="text-2xl font-bold text-blue-600 mt-1">14</p>
      </div>
    </div>

    <div class="bg-white rounded-lg border border-gray-200 divide-y divide-gray-200">
      <div v-for="alert in alerts" :key="alert.id" 
           :class="['px-6 py-4 flex items-center hover:bg-gray-50 cursor-pointer', alert.status === 'unread' ? 'bg-blue-50' : '']">
        <div :class="['w-10 h-10 rounded-full flex items-center justify-center mr-4',
               alert.level === 'high' ? 'bg-red-100' : alert.level === 'medium' ? 'bg-yellow-100' : 'bg-gray-100']">
          <span :class="['text-lg',
                 alert.level === 'high' ? 'text-red-600' : alert.level === 'medium' ? 'text-yellow-600' : 'text-gray-600']">
            {{ alert.level === 'high' ? '🚨' : alert.level === 'medium' ? '⚠️' : 'ℹ️' }}
          </span>
        </div>
        <div class="flex-1">
          <div class="flex items-center">
            <h4 class="text-sm font-medium text-gray-900">{{ alert.title }}</h4>
            <span :class="['ml-2 px-2 py-0.5 rounded text-xs',
                   alert.level === 'high' ? 'bg-red-100 text-red-700' :
                   alert.level === 'medium' ? 'bg-yellow-100 text-yellow-700' : 'bg-gray-100 text-gray-700']">
              {{ alert.level === 'high' ? '紧急' : alert.level === 'medium' ? '重要' : '一般' }}
            </span>
            <span v-if="alert.status === 'unread'" class="ml-2 w-2 h-2 bg-blue-600 rounded-full"></span>
          </div>
          <p class="text-sm text-gray-500 mt-1">{{ alert.desc }}</p>
          <div class="flex items-center mt-1 text-xs text-gray-400">
            <span>{{ alert.target }}</span>
            <span class="mx-2">·</span>
            <span>{{ alert.time }}</span>
          </div>
        </div>
<button class="px-3 py-1 text-sm text-blue-600 border border-blue-300 rounded hover:bg-blue-50" @click="handleProcessAlert(alert)">
          处理
        </button>
      </div>
    </div>

    <div class="flex justify-center">
      <Pagination :current="pagination.page" :total="pagination.total" :page-size="pagination.pageSize" @change="handlePageChange" />
    </div>
  </div>
</template>
