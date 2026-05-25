<template>
  <div class="p-6">
    <Breadcrumb :items="[{ name: '本地推', path: '/local' }, { name: '数据分析' }, { name: '报表导出' }]" />

    <div class="mb-6">
      <h1 class="text-2xl font-bold text-gray-900">报表导出</h1>
      <p class="text-gray-600 mt-1">创建并下载数据报表</p>
    </div>

    <div class="bg-white rounded-lg shadow p-6 mb-6">
      <h2 class="text-lg font-medium text-gray-900 mb-4">创建导出任务</h2>
      <div class="flex flex-wrap gap-4 items-end">
        <div>
          <label class="block text-sm text-gray-600 mb-1">报表标题</label>
          <input
            v-model="form.title"
            type="text"
            placeholder="请输入报表标题"
            class="border border-gray-300 rounded px-3 py-2 w-64"
          >
        </div>
        <div>
          <label class="block text-sm text-gray-600 mb-1">开始日期</label>
          <input type="date" v-model="form.startDate" class="border border-gray-300 rounded px-3 py-2">
        </div>
        <div>
          <label class="block text-sm text-gray-600 mb-1">结束日期</label>
          <input type="date" v-model="form.endDate" class="border border-gray-300 rounded px-3 py-2">
        </div>
        <button
          @click="handleCreateExport"
          :disabled="submitting"
          class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50"
        >
          {{ submitting ? '提交中...' : '提交导出' }}
        </button>
      </div>
    </div>

    <div class="bg-white rounded-lg shadow">
      <div class="px-6 py-4 border-b">
        <h2 class="text-lg font-medium text-gray-900">导出任务列表</h2>
      </div>
      <table class="min-w-full">
        <thead class="bg-gray-50">
          <tr>
            <th class="px-4 py-3 text-left text-sm font-medium text-gray-500">标题</th>
            <th class="px-4 py-3 text-left text-sm font-medium text-gray-500">状态</th>
            <th class="px-4 py-3 text-right text-sm font-medium text-gray-500">行数</th>
            <th class="px-4 py-3 text-left text-sm font-medium text-gray-500">创建时间</th>
            <th class="px-4 py-3 text-left text-sm font-medium text-gray-500">操作</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-200">
          <tr v-if="loadingList">
            <td colspan="5" class="px-4 py-8 text-center text-gray-500">加载中...</td>
          </tr>
          <tr v-else-if="exportList.length === 0">
            <td colspan="5" class="px-4 py-8 text-center text-gray-500">暂无导出任务</td>
          </tr>
          <tr v-for="task in exportList" :key="task.id">
            <td class="px-4 py-3 text-sm font-medium">{{ task.title }}</td>
            <td class="px-4 py-3 text-sm">
              <span :class="statusClass(task.status)">{{ statusLabel(task.status) }}</span>
            </td>
            <td class="px-4 py-3 text-sm text-right">{{ task.rows ?? '—' }}</td>
            <td class="px-4 py-3 text-sm text-gray-500">{{ task.created_at }}</td>
            <td class="px-4 py-3 text-sm">
              <button
                v-if="task.status === 'completed'"
                @click="handleDownload(task.id)"
                class="text-blue-600 hover:text-blue-800"
              >
                下载
              </button>
              <span v-else class="text-gray-400">—</span>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import Breadcrumb from '@/components/common/Breadcrumb.vue'
import { analyticsApi, type ExportTask } from '@/api/analytics'

interface ExportRow extends ExportTask {
  rows?: number
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

const form = ref({
  title: '',
  startDate,
  endDate
})

const submitting = ref(false)
const loadingList = ref(false)
const exportList = ref<ExportRow[]>([])

function statusLabel(status: string): string {
  const map: Record<string, string> = {
    pending: '等待中',
    processing: '处理中',
    completed: '已完成',
    failed: '失败'
  }
  return map[status] ?? status
}

function statusClass(status: string): string {
  const map: Record<string, string> = {
    pending: 'text-yellow-600',
    processing: 'text-blue-600',
    completed: 'text-green-600',
    failed: 'text-red-600'
  }
  return map[status] ?? 'text-gray-600'
}

async function fetchExports(): Promise<void> {
  loadingList.value = true
  try {
    const res = await analyticsApi.getExports()
    exportList.value = Array.isArray(res) ? res : []
  } finally {
    loadingList.value = false
  }
}

async function handleCreateExport(): Promise<void> {
  if (!form.value.title.trim()) {
    alert('请输入报表标题')
    return
  }

  submitting.value = true
  try {
    await analyticsApi.createExport({
      title: form.value.title,
      start_date: form.value.startDate,
      end_date: form.value.endDate
    })
    form.value.title = ''
    await fetchExports()
  } finally {
    submitting.value = false
  }
}

function handleDownload(id: number): void {
  window.open(`/api/v1/analytics/exports/${id}/download`)
}

onMounted(() => {
  fetchExports()
})
</script>
