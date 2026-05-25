<template>
  <div class="p-6">
    <h1 class="text-2xl font-bold mb-6">批量任务列表</h1>

    <div class="overflow-x-auto bg-white rounded-lg shadow">
      <table class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">标题</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">类型</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">状态</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">进度</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">操作</th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200">
          <tr v-for="task in tasks" :key="task.id" class="hover:bg-gray-50">
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ task.title }}</td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ task.task_type }}</td>
            <td class="px-6 py-4 whitespace-nowrap">
              <span :class="statusClass(task.status)" class="px-2 py-1 text-xs font-semibold rounded-full">
                {{ statusLabel(task.status) }}
              </span>
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
              {{ task.account_ids?.length ?? 0 }} 个账户
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm">
              <button
                v-if="task.status === 'running'"
                class="px-3 py-1 text-xs font-medium text-white bg-red-500 rounded hover:bg-red-600 transition-colors"
                @click="handleCancel(task.id)"
              >
                取消
              </button>
              <button
                v-if="task.status === 'failed' || task.status === 'partial_success'"
                class="px-3 py-1 text-xs font-medium text-white bg-blue-500 rounded hover:bg-blue-600 transition-colors"
                @click="handleRetry(task.id)"
              >
                重试
              </button>
            </td>
          </tr>
          <tr v-if="tasks.length === 0">
            <td colspan="5" class="px-6 py-8 text-center text-sm text-gray-400">暂无任务</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { batchApi, type BatchTask } from '@/api/batch'

const tasks = ref<BatchTask[]>([])

const POLL_INTERVAL_MS = 5000

let pollTimer: ReturnType<typeof setInterval> | null = null

async function fetchTasks(): Promise<void> {
  try {
    const res = await batchApi.getList()
    tasks.value = res.list ?? []
  } catch {
    // 静默处理轮询错误，避免干扰用户
  }
}

function startPolling(): void {
  pollTimer = setInterval(fetchTasks, POLL_INTERVAL_MS)
}

function stopPolling(): void {
  if (pollTimer !== null) {
    clearInterval(pollTimer)
    pollTimer = null
  }
}

async function handleCancel(id: number): Promise<void> {
  try {
    await batchApi.cancel(id)
    await fetchTasks()
  } catch {
    // 错误处理可根据项目统一方案补充
  }
}

async function handleRetry(id: number): Promise<void> {
  try {
    await batchApi.retry(id)
    await fetchTasks()
  } catch {
    // 错误处理可根据项目统一方案补充
  }
}

const STATUS_CONFIG: Record<string, { label: string; classes: string }> = {
  running: { label: '运行中', classes: 'bg-blue-100 text-blue-800' },
  completed: { label: '已完成', classes: 'bg-green-100 text-green-800' },
  failed: { label: '失败', classes: 'bg-red-100 text-red-800' },
  partial_success: { label: '部分成功', classes: 'bg-yellow-100 text-yellow-800' },
  cancelled: { label: '已取消', classes: 'bg-gray-100 text-gray-800' },
}

function statusClass(status: string): string {
  return STATUS_CONFIG[status]?.classes ?? 'bg-gray-100 text-gray-800'
}

function statusLabel(status: string): string {
  return STATUS_CONFIG[status]?.label ?? status
}

onMounted(() => {
  fetchTasks()
  startPolling()
})

onUnmounted(() => {
  stopPolling()
})
</script>
