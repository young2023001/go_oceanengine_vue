<script setup lang="ts">
import { ref, onMounted } from 'vue'
import Breadcrumb from '@/components/common/Breadcrumb.vue'
import request from '@/api/request'
import { scopeApi } from '@/api/scope'
import { accountApi, type Account } from '@/api/account'

interface User {
  id: number
  username: string
  nickname?: string
}

const users = ref<User[]>([])
const accounts = ref<Account[]>([])
const selectedUserId = ref<number | null>(null)
const checkedAccountIds = ref<number[]>([])
const loading = ref(false)
const saving = ref(false)

const fetchUsers = async () => {
  const res = await request.get<User[]>('/system/users')
  users.value = Array.isArray(res) ? res : []
}

const fetchAccounts = async () => {
  const res = await accountApi.getList({ page: 1, page_size: 1000 })
  accounts.value = Array.isArray(res) ? res : (res as any)?.list ?? []
}

const selectUser = async (user: User) => {
  selectedUserId.value = user.id
  loading.value = true
  try {
    const res = await scopeApi.getScope(user.id)
    checkedAccountIds.value = (res as any)?.account_ids ?? []
  } finally {
    loading.value = false
  }
}

const toggleAccount = (accountId: number) => {
  const idx = checkedAccountIds.value.indexOf(accountId)
  if (idx >= 0) {
    checkedAccountIds.value = [
      ...checkedAccountIds.value.slice(0, idx),
      ...checkedAccountIds.value.slice(idx + 1)
    ]
  } else {
    checkedAccountIds.value = [...checkedAccountIds.value, accountId]
  }
}

const handleSave = async () => {
  if (!selectedUserId.value) return
  saving.value = true
  try {
    await scopeApi.setScope(selectedUserId.value, { account_ids: checkedAccountIds.value })
    alert('权限保存成功')
  } finally {
    saving.value = false
  }
}

onMounted(async () => {
  await Promise.all([fetchUsers(), fetchAccounts()])
})
</script>

<template>
  <div class="space-y-6">
    <div>
      <Breadcrumb :items="[{ name: '系统管理' }, { name: '权限分配' }]" />
      <h1 class="text-3xl font-bold text-gray-900">投手权限分配</h1>
      <p class="mt-2 text-gray-600">为投手分配可操作的广告账户权限</p>
    </div>

    <div class="flex gap-6 min-h-[500px]">
      <!-- 左侧：用户列表 -->
      <div class="w-64 shrink-0 bg-white rounded-lg border border-gray-200">
        <div class="px-4 py-3 border-b border-gray-200">
          <h3 class="text-sm font-semibold text-gray-700">用户列表</h3>
        </div>
        <div class="overflow-y-auto max-h-[450px]">
          <button
            v-for="user in users"
            :key="user.id"
            @click="selectUser(user)"
            :class="[
              'w-full px-4 py-3 text-left text-sm border-b border-gray-100 hover:bg-blue-50 transition-colors',
              selectedUserId === user.id ? 'bg-blue-50 text-blue-700 font-medium' : 'text-gray-700'
            ]"
          >
            {{ user.nickname || user.username }}
          </button>
          <div v-if="users.length === 0" class="px-4 py-8 text-center text-gray-500 text-sm">
            暂无用户
          </div>
        </div>
      </div>

      <!-- 右侧：账户权限 -->
      <div class="flex-1 bg-white rounded-lg border border-gray-200">
        <div class="px-6 py-3 border-b border-gray-200 flex items-center justify-between">
          <h3 class="text-sm font-semibold text-gray-700">
            账户权限范围
            <span v-if="selectedUserId" class="text-gray-400 font-normal ml-2">
              已选 {{ checkedAccountIds.length }} / {{ accounts.length }}
            </span>
          </h3>
          <button
            v-if="selectedUserId"
            @click="handleSave"
            :disabled="saving"
            class="px-4 py-1.5 bg-blue-600 text-white text-sm rounded-lg hover:bg-blue-700 disabled:opacity-50"
          >
            {{ saving ? '保存中...' : '保存' }}
          </button>
        </div>

        <div v-if="!selectedUserId" class="flex items-center justify-center h-80 text-gray-400">
          请从左侧选择用户
        </div>

        <div v-else-if="loading" class="flex items-center justify-center h-80 text-gray-400">
          加载中...
        </div>

        <div v-else class="p-6 overflow-y-auto max-h-[450px]">
          <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
            <label
              v-for="account in accounts"
              :key="account.id"
              class="flex items-center gap-3 px-3 py-2 rounded-lg border border-gray-200 hover:border-blue-300 cursor-pointer transition-colors"
              :class="checkedAccountIds.includes(account.id) ? 'bg-blue-50 border-blue-300' : ''"
            >
              <input
                type="checkbox"
                :checked="checkedAccountIds.includes(account.id)"
                @change="toggleAccount(account.id)"
                class="w-4 h-4 text-blue-600 rounded border-gray-300 focus:ring-blue-500"
              />
              <div class="min-w-0">
                <div class="text-sm font-medium text-gray-900 truncate">{{ account.name }}</div>
                <div class="text-xs text-gray-500">ID: {{ account.local_account_id }}</div>
              </div>
            </label>
          </div>
          <div v-if="accounts.length === 0" class="text-center text-gray-500 py-8">
            暂无账户数据
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
