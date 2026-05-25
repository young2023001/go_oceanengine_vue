<template>
  <div class="p-6">
    <div class="flex justify-between items-center mb-6">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">账户列表</h1>
        <p class="text-gray-600 mt-1">管理巨量引擎广告账户</p>
      </div>
      <button @click="showImportModal = true" class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700">
        批量导入
      </button>
    </div>

    <!-- 账户列表 -->
    <div class="bg-white rounded-lg shadow">
      <table class="min-w-full">
        <thead class="bg-gray-50">
          <tr>
            <th class="px-4 py-3 text-left text-sm font-medium text-gray-500">ID</th>
            <th class="px-4 py-3 text-left text-sm font-medium text-gray-500">巨量账户ID</th>
            <th class="px-4 py-3 text-left text-sm font-medium text-gray-500">名称</th>
            <th class="px-4 py-3 text-left text-sm font-medium text-gray-500">状态</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-200">
          <tr v-for="account in accounts" :key="account.id">
            <td class="px-4 py-3 text-sm text-gray-600">{{ account.id }}</td>
            <td class="px-4 py-3 text-sm text-gray-600">{{ account.local_account_id }}</td>
            <td class="px-4 py-3 font-medium text-gray-900">{{ account.name }}</td>
            <td class="px-4 py-3">
              <span class="px-2 py-1 text-xs rounded-full bg-green-100 text-green-800">正常</span>
            </td>
          </tr>
          <tr v-if="accounts.length === 0 && !loading">
            <td colspan="4" class="px-4 py-8 text-center text-gray-500">暂无账户数据</td>
          </tr>
        </tbody>
      </table>

      <!-- 分页 -->
      <div v-if="pagination.total > 0" class="p-4 border-t flex items-center justify-between">
        <span class="text-sm text-gray-500">共 {{ pagination.total }} 条</span>
        <div class="flex space-x-2">
          <button
            :disabled="pagination.page <= 1"
            @click="handlePageChange(pagination.page - 1)"
            class="px-3 py-1 border border-gray-300 rounded text-sm disabled:opacity-50 hover:bg-gray-50"
          >
            上一页
          </button>
          <span class="px-3 py-1 text-sm text-gray-700">{{ pagination.page }} / {{ totalPages }}</span>
          <button
            :disabled="pagination.page >= totalPages"
            @click="handlePageChange(pagination.page + 1)"
            class="px-3 py-1 border border-gray-300 rounded text-sm disabled:opacity-50 hover:bg-gray-50"
          >
            下一页
          </button>
        </div>
      </div>
    </div>

    <!-- 批量导入弹窗 -->
    <div v-if="showImportModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg shadow-xl w-full max-w-lg p-6">
        <h2 class="text-lg font-bold text-gray-900 mb-4">批量导入账户</h2>
        <div class="space-y-3">
          <label class="block text-sm font-medium text-gray-700">
            请输入 JSON 数据
          </label>
          <textarea
            v-model="importJson"
            rows="8"
            placeholder='[{"local_account_id":100001,"name":"门店A"}]'
            class="w-full border border-gray-300 rounded-lg px-3 py-2 font-mono text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
          ></textarea>
          <p class="text-xs text-gray-500">
            格式：[{"local_account_id": 数字, "name": "名称"}, ...]
          </p>
          <p v-if="importError" class="text-xs text-red-600">{{ importError }}</p>
        </div>
        <div class="flex justify-end space-x-3 mt-6">
          <button
            type="button"
            @click="closeImportModal"
            class="px-4 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50"
          >
            取消
          </button>
          <button
            @click="handleImport"
            :disabled="importing"
            class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50"
          >
            {{ importing ? '导入中...' : '确认导入' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { accountApi, type Account, type ImportAccountItem } from '@/api/account'

const loading = ref(false)
const accounts = ref<Account[]>([])
const showImportModal = ref(false)
const importing = ref(false)
const importJson = ref('')
const importError = ref('')

const pagination = ref({
  page: 1,
  pageSize: 20,
  total: 0
})

const totalPages = computed(() => Math.ceil(pagination.value.total / pagination.value.pageSize))

const fetchAccounts = async () => {
  loading.value = true
  try {
    const res = await accountApi.getList({
      page: pagination.value.page,
      page_size: pagination.value.pageSize
    })
    if (res) {
      const data = res as any
      accounts.value = Array.isArray(data) ? data : (data.list || [])
      pagination.value.total = Array.isArray(data) ? data.length : (data.total || 0)
    }
  } catch (e: unknown) {
    const message = e instanceof Error ? e.message : '获取账户列表失败'
    alert(message)
  } finally {
    loading.value = false
  }
}

const handlePageChange = (page: number) => {
  pagination.value.page = page
  fetchAccounts()
}

const closeImportModal = () => {
  showImportModal.value = false
  importJson.value = ''
  importError.value = ''
}

const handleImport = async () => {
  importError.value = ''

  let items: ImportAccountItem[]
  try {
    items = JSON.parse(importJson.value)
  } catch {
    importError.value = 'JSON 格式错误，请检查输入'
    return
  }

  if (!Array.isArray(items) || items.length === 0) {
    importError.value = '数据必须是非空数组'
    return
  }

  const isValid = items.every(
    (item) => typeof item.local_account_id === 'number' && typeof item.name === 'string'
  )
  if (!isValid) {
    importError.value = '每条数据必须包含 local_account_id (数字) 和 name (字符串)'
    return
  }

  importing.value = true
  try {
    await accountApi.importAccounts({ items })
    closeImportModal()
    await fetchAccounts()
  } catch (e: unknown) {
    const message = e instanceof Error ? e.message : '导入失败'
    importError.value = message
  } finally {
    importing.value = false
  }
}

onMounted(() => {
  fetchAccounts()
})
</script>
