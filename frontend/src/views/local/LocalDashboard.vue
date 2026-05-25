<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Line } from 'vue-chartjs'
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler
} from 'chart.js'
import { analyticsApi, type CompareData, type TrendData, type RankItem } from '@/api/analytics'
import Breadcrumb from '@/components/common/Breadcrumb.vue'

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend, Filler)

const loading = ref(false)
const compareData = ref<CompareData | null>(null)
const trendData = ref<TrendData[]>([])
const rankData = ref<RankItem[]>([])

function getDateString(daysAgo: number): string {
  const date = new Date()
  date.setDate(date.getDate() - daysAgo)
  return date.toISOString().slice(0, 10)
}

function calcChangePercent(current: number, previous: number): string {
  if (previous === 0) return '+0.0'
  const change = ((current - previous) / previous) * 100
  const sign = change >= 0 ? '+' : ''
  return `${sign}${change.toFixed(1)}`
}

const summaryCards = computed(() => {
  if (!compareData.value) return []
  const { current, previous } = compareData.value
  return [
    {
      title: '今日消耗',
      value: `¥${current.cost.toLocaleString()}`,
      change: calcChangePercent(current.cost, previous.cost),
      color: 'blue'
    },
    {
      title: '今日展示',
      value: current.impressions.toLocaleString(),
      change: calcChangePercent(current.impressions, previous.impressions),
      color: 'green'
    },
    {
      title: '今日点击',
      value: current.clicks.toLocaleString(),
      change: calcChangePercent(current.clicks, previous.clicks),
      color: 'purple'
    },
    {
      title: '今日转化',
      value: current.conversions.toLocaleString(),
      change: calcChangePercent(current.conversions, previous.conversions),
      color: 'orange'
    }
  ]
})

const chartData = computed(() => ({
  labels: trendData.value.map(item => item.date.slice(5)),
  datasets: [
    {
      label: '消耗',
      data: trendData.value.map(item => item.value),
      borderColor: 'rgb(59, 130, 246)',
      backgroundColor: 'rgba(59, 130, 246, 0.1)',
      fill: true,
      tension: 0.4
    }
  ]
}))

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: { position: 'top' as const }
  },
  scales: {
    y: {
      beginAtZero: true,
      ticks: {
        callback: (value: string | number) => `¥${Number(value).toLocaleString()}`
      }
    }
  }
}

const maxRankValue = computed(() => {
  if (rankData.value.length === 0) return 1
  return rankData.value[0]?.value ?? 1
})

async function fetchData(): Promise<void> {
  loading.value = true
  try {
    const today = getDateString(0)
    const sevenDaysAgo = getDateString(6)

    const [compareRes, trendRes, rankRes] = await Promise.all([
      analyticsApi.getCompare({ date: today }),
      analyticsApi.getTrend({ start_date: sevenDaysAgo, end_date: today }),
      analyticsApi.getRank({ start_date: sevenDaysAgo, end_date: today, order_by: 'cost', limit: 10 })
    ])

    compareData.value = compareRes as unknown as CompareData
    trendData.value = trendRes as unknown as TrendData[]
    rankData.value = rankRes as unknown as RankItem[]
  } finally {
    loading.value = false
  }
}

function getChangeClass(change: string): string {
  if (change.startsWith('+')) return 'text-red-500'
  if (change.startsWith('-')) return 'text-green-500'
  return 'text-gray-500'
}

function getBarColor(index: number): string {
  const colors = ['bg-blue-500', 'bg-blue-400', 'bg-blue-300']
  if (index < 3) return colors[index]
  return 'bg-blue-200'
}

onMounted(() => {
  fetchData()
})
</script>

<template>
  <div class="p-6 space-y-6">
    <Breadcrumb :items="[{ name: '本地推', path: '/local' }, { name: '工作台' }]" />

    <div class="mb-2">
      <h1 class="text-2xl font-bold text-gray-900">本地推工作台</h1>
      <p class="text-gray-600 mt-1">本地生活服务推广数据概览</p>
    </div>

    <!-- 汇总卡片 -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <div
        v-for="card in summaryCards"
        :key="card.title"
        class="bg-white rounded-lg border border-gray-200 p-5"
      >
        <p class="text-sm text-gray-500">{{ card.title }}</p>
        <p class="mt-2 text-2xl font-bold text-gray-900">{{ card.value }}</p>
        <p class="mt-1 text-sm" :class="getChangeClass(card.change)">
          环比 {{ card.change }}%
        </p>
      </div>

      <!-- 加载占位 -->
      <template v-if="summaryCards.length === 0 && loading">
        <div
          v-for="i in 4"
          :key="i"
          class="bg-white rounded-lg border border-gray-200 p-5 animate-pulse"
        >
          <div class="h-4 bg-gray-200 rounded w-20 mb-3"></div>
          <div class="h-7 bg-gray-200 rounded w-32 mb-2"></div>
          <div class="h-4 bg-gray-200 rounded w-16"></div>
        </div>
      </template>
    </div>

    <!-- 趋势图 -->
    <div class="bg-white rounded-lg border border-gray-200 p-6">
      <h3 class="text-lg font-semibold text-gray-900 mb-4">最近 7 天消耗趋势</h3>
      <div v-if="trendData.length > 0" style="height: 300px;">
        <Line :data="chartData" :options="chartOptions" />
      </div>
      <div v-else-if="loading" class="h-[300px] flex items-center justify-center">
        <p class="text-gray-400">加载中...</p>
      </div>
      <div v-else class="h-[300px] flex items-center justify-center">
        <p class="text-gray-400">暂无数据</p>
      </div>
    </div>

    <!-- 排行榜表格 -->
    <div class="bg-white rounded-lg border border-gray-200 p-6">
      <h3 class="text-lg font-semibold text-gray-900 mb-4">TOP 10 门店消耗排行</h3>
      <div v-if="rankData.length > 0" class="overflow-x-auto">
        <table class="w-full text-sm">
          <thead>
            <tr class="border-b border-gray-200">
              <th class="text-left py-3 px-2 font-medium text-gray-500 w-12">排名</th>
              <th class="text-left py-3 px-2 font-medium text-gray-500">门店名称</th>
              <th class="text-left py-3 px-2 font-medium text-gray-500 w-48">消耗</th>
              <th class="text-right py-3 px-2 font-medium text-gray-500 w-24">金额</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="(item, index) in rankData"
              :key="item.id"
              class="border-b border-gray-100 hover:bg-gray-50"
            >
              <td class="py-3 px-2">
                <span
                  class="inline-flex items-center justify-center w-6 h-6 rounded-full text-xs font-bold"
                  :class="index < 3 ? 'bg-blue-100 text-blue-700' : 'bg-gray-100 text-gray-600'"
                >
                  {{ index + 1 }}
                </span>
              </td>
              <td class="py-3 px-2 text-gray-900 font-medium">{{ item.name }}</td>
              <td class="py-3 px-2">
                <div class="flex items-center gap-2">
                  <div class="flex-1 h-2 bg-gray-100 rounded-full overflow-hidden">
                    <div
                      class="h-full rounded-full transition-all"
                      :class="getBarColor(index)"
                      :style="{ width: `${(item.value / maxRankValue) * 100}%` }"
                    ></div>
                  </div>
                </div>
              </td>
              <td class="py-3 px-2 text-right text-gray-900">
                ¥{{ item.value.toLocaleString() }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <div v-else-if="loading" class="py-12 text-center">
        <p class="text-gray-400">加载中...</p>
      </div>
      <div v-else class="py-12 text-center">
        <p class="text-gray-400">暂无数据</p>
      </div>
    </div>
  </div>
</template>
