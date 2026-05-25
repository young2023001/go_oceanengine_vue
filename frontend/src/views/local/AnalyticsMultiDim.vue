<template>
  <div class="p-6">
    <Breadcrumb :items="[{ name: '本地推', path: '/local' }, { name: '数据分析' }, { name: '多维分析' }]" />

    <div class="mb-6">
      <h1 class="text-2xl font-bold text-gray-900">多维分析</h1>
      <p class="text-gray-600 mt-1">按不同维度查看投放数据排行</p>
    </div>

    <div class="bg-white rounded-lg shadow p-4 mb-6">
      <div class="flex flex-wrap gap-4 items-center">
        <input type="date" v-model="filters.startDate" class="border border-gray-300 rounded px-3 py-2">
        <span>至</span>
        <input type="date" v-model="filters.endDate" class="border border-gray-300 rounded px-3 py-2">
        <select v-model="filters.dimension" class="border border-gray-300 rounded px-3 py-2">
          <option value="store">按门店</option>
          <option value="project" disabled>按项目（即将上线）</option>
          <option value="promotion" disabled>按计划（即将上线）</option>
        </select>
        <button @click="fetchData" class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700">查询</button>
      </div>
    </div>

    <div class="bg-white rounded-lg shadow">
      <table class="min-w-full">
        <thead class="bg-gray-50">
          <tr>
            <th class="px-4 py-3 text-left text-sm font-medium text-gray-500">名称</th>
            <th class="px-4 py-3 text-right text-sm font-medium text-gray-500">消耗</th>
            <th class="px-4 py-3 text-right text-sm font-medium text-gray-500">展示</th>
            <th class="px-4 py-3 text-right text-sm font-medium text-gray-500">点击</th>
            <th class="px-4 py-3 text-right text-sm font-medium text-gray-500">转化</th>
            <th class="px-4 py-3 text-right text-sm font-medium text-gray-500">转化成本</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-200">
          <tr v-if="loading">
            <td colspan="6" class="px-4 py-8 text-center text-gray-500">加载中...</td>
          </tr>
          <tr v-else-if="rankList.length === 0">
            <td colspan="6" class="px-4 py-8 text-center text-gray-500">暂无数据</td>
          </tr>
          <tr v-for="item in rankList" :key="item.id">
            <td class="px-4 py-3 text-sm font-medium">{{ item.name }}</td>
            <td class="px-4 py-3 text-sm text-right">¥{{ item.cost.toLocaleString() }}</td>
            <td class="px-4 py-3 text-sm text-right">{{ item.impressions.toLocaleString() }}</td>
            <td class="px-4 py-3 text-sm text-right">{{ item.clicks.toLocaleString() }}</td>
            <td class="px-4 py-3 text-sm text-right">{{ item.conversions.toLocaleString() }}</td>
            <td class="px-4 py-3 text-sm text-right">¥{{ item.conversionCost }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import Breadcrumb from '@/components/common/Breadcrumb.vue'
import { analyticsApi } from '@/api/analytics'

interface RankRow {
  id: number
  name: string
  cost: number
  impressions: number
  clicks: number
  conversions: number
  conversionCost: string
}

function formatDate(date: Date): string {
  const y = date.getFullYear()
  const m = String(date.getMonth() + 1).padStart(2, '0')
  const d = String(date.getDate()).padStart(2, '0')
  return `${y}-${m}-${d}`
}

function getDefaultDateRange(): { startDate: string; endDate: string } {
  const end = new Date()
  const start = new Date()
  start.setDate(start.getDate() - 6)
  return { startDate: formatDate(start), endDate: formatDate(end) }
}

const { startDate, endDate } = getDefaultDateRange()

const filters = ref({
  startDate,
  endDate,
  dimension: 'store'
})

const loading = ref(false)
const rankList = ref<RankRow[]>([])

async function fetchData(): Promise<void> {
  if (filters.value.dimension !== 'store') {
    return
  }

  loading.value = true
  try {
    const res = await analyticsApi.getRank({
      start_date: filters.value.startDate,
      end_date: filters.value.endDate,
      order_by: 'cost',
      limit: 50
    })
    const items = res.data ?? []
    rankList.value = items.map((item) => ({
      id: item.id,
      name: item.name,
      cost: item.value,
      impressions: 0,
      clicks: 0,
      conversions: 0,
      conversionCost: item.value > 0 ? '—' : '—'
    }))
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchData()
})
</script>
