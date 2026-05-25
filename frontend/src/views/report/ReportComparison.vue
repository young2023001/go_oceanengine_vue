<script setup lang="ts">
import { ref } from 'vue'
import Breadcrumb from '@/components/common/Breadcrumb.vue'

const compareType = ref('period')
const period1 = ref('2025-11-21')
const period2 = ref('2025-11-28')

const comparisonData = ref([
  { metric: '消耗', current: 125680, previous: 98560, change: 27.5 },
  { metric: '展示', current: 2568000, previous: 2156000, change: 19.1 },
  { metric: '点击', current: 89560, previous: 72300, change: 23.9 },
  { metric: '转化', current: 3256, previous: 2890, change: 12.7 },
  { metric: 'CTR', current: 3.49, previous: 3.35, change: 4.2, isPercent: true },
  { metric: 'CVR', current: 3.64, previous: 4.00, change: -9.0, isPercent: true },
  { metric: 'CPC', current: 1.40, previous: 1.36, change: 2.9, isCurrency: true },
  { metric: 'CPA', current: 38.6, previous: 34.1, change: 13.2, isCurrency: true }
])

const formatValue = (value: number, isPercent?: boolean, isCurrency?: boolean) => {
  if (isPercent) return value.toFixed(2) + '%'
  if (isCurrency) return '¥' + value.toFixed(2)
  if (value >= 10000) return (value / 10000).toFixed(1) + '万'
  return value.toLocaleString()
}

const handleCompare = () => {
  // TODO: implement
}
</script>

<template>
  <div class="space-y-6">
    <div>
      <Breadcrumb :items="[{ name: '报表中心' }, { name: '数据对比' }]" />
      <h1 class="text-3xl font-bold text-gray-900">数据对比分析</h1>
      <p class="mt-2 text-gray-600">对比不同时段的投放效果</p>
    </div>

    <div class="bg-white rounded-lg border border-gray-200 p-4">
      <div class="flex flex-wrap items-center gap-4">
        <select v-model="compareType" class="px-4 py-2 border border-gray-300 rounded-lg">
          <option value="period">时间对比</option>
          <option value="campaign">计划对比</option>
          <option value="creative">素材对比</option>
        </select>
        <div class="flex items-center gap-2">
          <input v-model="period1" type="date" class="px-4 py-2 border border-gray-300 rounded-lg">
          <span class="text-gray-500">VS</span>
          <input v-model="period2" type="date" class="px-4 py-2 border border-gray-300 rounded-lg">
        </div>
<button class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700" @click="handleCompare">
          开始对比
        </button>
      </div>
    </div>

    <div class="bg-white rounded-lg border border-gray-200">
      <table class="w-full">
        <thead class="bg-gray-50 border-b border-gray-200">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-semibold text-gray-600 uppercase">指标</th>
            <th class="px-6 py-3 text-left text-xs font-semibold text-gray-600 uppercase">当前周期</th>
            <th class="px-6 py-3 text-left text-xs font-semibold text-gray-600 uppercase">对比周期</th>
            <th class="px-6 py-3 text-left text-xs font-semibold text-gray-600 uppercase">变化</th>
            <th class="px-6 py-3 text-left text-xs font-semibold text-gray-600 uppercase">趋势</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-200">
          <tr v-for="row in comparisonData" :key="row.metric" class="hover:bg-gray-50">
            <td class="px-6 py-4 font-medium text-gray-900">{{ row.metric }}</td>
            <td class="px-6 py-4 text-gray-900">{{ formatValue(row.current, row.isPercent, row.isCurrency) }}</td>
            <td class="px-6 py-4 text-gray-600">{{ formatValue(row.previous, row.isPercent, row.isCurrency) }}</td>
            <td class="px-6 py-4">
              <span :class="['font-medium', row.change >= 0 ? 'text-green-600' : 'text-red-600']">
                {{ row.change >= 0 ? '+' : '' }}{{ row.change.toFixed(1) }}%
              </span>
            </td>
            <td class="px-6 py-4 text-2xl">
              {{ row.change >= 0 ? '📈' : '📉' }}
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
