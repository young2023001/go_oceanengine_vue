<template>
  <div class="p-6">
    <div class="flex justify-between items-center mb-6">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">租户管理</h1>
        <p class="text-gray-600 mt-1">管理巨量引擎 OAuth 租户授权</p>
      </div>
      <button @click="showCreateModal = true" class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700">
        新建租户
      </button>
    </div>

    <!-- 租户列表 -->
    <div class="bg-white rounded-lg shadow">
      <table class="min-w-full">
        <thead class="bg-gray-50">
          <tr>
            <th class="px-4 py-3 text-left text-sm font-medium text-gray-500">租户名称</th>
            <th class="px-4 py-3 text-left text-sm font-medium text-gray-500">AppID</th>
            <th class="px-4 py-3 text-left text-sm font-medium text-gray-500">Token状态</th>
            <th class="px-4 py-3 text-left text-sm font-medium text-gray-500">操作</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-200">
          <tr v-for="tenant in tenants" :key="tenant.id">
            <td class="px-4 py-3">
              <div class="font-medium text-gray-900">{{ tenant.name }}</div>
            </td>
            <td class="px-4 py-3 text-sm text-gray-600">{{ tenant.oauth_app_id }}</td>
            <td class="px-4 py-3">
              <span :class="tokenStatusClass(tenant)" class="px-2 py-1 text-xs rounded-full">
                {{ tokenStatusLabel(tenant) }}
              </span>
            </td>
            <td class="px-4 py-3">
              <button @click="handleAuthorize(tenant)" class="text-blue-600 hover:text-blue-800 text-sm">
                {{ tenant.token_status === 'active' ? '重新授权' : '发起授权' }}
              </button>
            </td>
          </tr>
          <tr v-if="tenants.length === 0 && !loading">
            <td colspan="4" class="px-4 py-8 text-center text-gray-500">暂无租户数据</td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- 新建租户弹窗 -->
    <div v-if="showCreateModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg shadow-xl w-full max-w-md p-6">
        <h2 class="text-lg font-bold text-gray-900 mb-4">新建租户</h2>
        <form @submit.prevent="handleCreate">
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">租户名称</label>
              <input
                v-model="createForm.name"
                type="text"
                required
                placeholder="请输入租户名称"
                class="w-full border border-gray-300 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">AppID</label>
              <input
                v-model="createForm.oauth_app_id"
                type="text"
                required
                placeholder="请输入 OAuth AppID"
                class="w-full border border-gray-300 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Secret</label>
              <input
                v-model="createForm.oauth_secret"
                type="text"
                required
                placeholder="请输入 OAuth Secret"
                class="w-full border border-gray-300 rounded-lg px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
          </div>
          <div class="flex justify-end space-x-3 mt-6">
            <button
              type="button"
              @click="showCreateModal = false"
              class="px-4 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50"
            >
              取消
            </button>
            <button
              type="submit"
              :disabled="creating"
              class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50"
            >
              {{ creating ? '创建中...' : '确认创建' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { tenantApi, type Tenant } from '@/api/tenant'

interface TenantWithStatus extends Tenant {
  token_status?: 'active' | 'expiring' | 'expired' | 'failed'
}

const loading = ref(false)
const tenants = ref<TenantWithStatus[]>([])
const showCreateModal = ref(false)
const creating = ref(false)

const createForm = ref({
  name: '',
  oauth_app_id: '',
  oauth_secret: ''
})

const tokenStatusClass = (tenant: TenantWithStatus): string => {
  switch (tenant.token_status) {
    case 'active':
      return 'bg-green-100 text-green-800'
    case 'expiring':
      return 'bg-yellow-100 text-yellow-800'
    case 'expired':
    case 'failed':
      return 'bg-red-100 text-red-800'
    default:
      return 'bg-gray-100 text-gray-800'
  }
}

const tokenStatusLabel = (tenant: TenantWithStatus): string => {
  switch (tenant.token_status) {
    case 'active':
      return '正常'
    case 'expiring':
      return '即将过期'
    case 'expired':
      return '已过期'
    case 'failed':
      return '授权失败'
    default:
      return '未授权'
  }
}

const fetchTenants = async () => {
  loading.value = true
  try {
    const res = await tenantApi.getList()
    tenants.value = (res as TenantWithStatus[]) || []
  } catch (e: unknown) {
    const message = e instanceof Error ? e.message : '获取租户列表失败'
    alert(message)
  } finally {
    loading.value = false
  }
}

const handleCreate = async () => {
  creating.value = true
  try {
    await tenantApi.create({
      name: createForm.value.name,
      oauth_app_id: createForm.value.oauth_app_id,
      oauth_secret: createForm.value.oauth_secret
    })
    showCreateModal.value = false
    createForm.value = { name: '', oauth_app_id: '', oauth_secret: '' }
    await fetchTenants()
  } catch (e: unknown) {
    const message = e instanceof Error ? e.message : '创建租户失败'
    alert(message)
  } finally {
    creating.value = false
  }
}

const handleAuthorize = async (tenant: TenantWithStatus) => {
  try {
    const res = await tenantApi.getOAuthURL(tenant.id)
    if (res && res.url) {
      window.open(res.url, '_blank')
    }
  } catch (e: unknown) {
    const message = e instanceof Error ? e.message : '获取授权链接失败'
    alert(message)
  }
}

onMounted(() => {
  fetchTenants()
})
</script>
