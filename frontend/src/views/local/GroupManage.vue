<template>
  <div class="p-6">
    <Breadcrumb :items="[{ name: '本地推', path: '/local' }, { name: '分组管理' }]" />

    <div class="flex justify-between items-center mb-6">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">分组管理</h1>
        <p class="text-gray-600 mt-1">管理账户分组与成员</p>
      </div>
      <button @click="showCreateDialog = true" class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700">
        新建分组
      </button>
    </div>

    <div class="flex gap-6">
      <!-- 左侧：分组列表 -->
      <div class="w-1/2 bg-white rounded-lg shadow">
        <div class="p-4 border-b">
          <select v-model="typeFilter" class="border border-gray-300 rounded-lg px-3 py-2 w-full">
            <option value="">全部类型</option>
            <option value="franchisee">加盟商</option>
            <option value="region">区域</option>
            <option value="custom">自定义</option>
          </select>
        </div>
        <table class="min-w-full">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-4 py-3 text-left text-sm font-medium text-gray-500">名称</th>
              <th class="px-4 py-3 text-left text-sm font-medium text-gray-500">类型</th>
              <th class="px-4 py-3 text-left text-sm font-medium text-gray-500">操作</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr v-for="group in filteredGroups" :key="group.id"
              :class="selectedGroup?.id === group.id ? 'bg-blue-50' : ''"
              class="cursor-pointer hover:bg-gray-50"
              @click="selectGroup(group)">
              <td class="px-4 py-3 text-sm font-medium">{{ group.name }}</td>
              <td class="px-4 py-3 text-sm">
                <span class="px-2 py-1 text-xs rounded-full"
                  :class="typeTagClass(group.type)">
                  {{ typeLabel(group.type) }}
                </span>
              </td>
              <td class="px-4 py-3">
                <div class="flex space-x-2">
                  <button @click.stop="handleEdit(group)" class="text-blue-600 hover:text-blue-800 text-sm">编辑</button>
                  <button @click.stop="handleDelete(group)" class="text-red-600 hover:text-red-800 text-sm">删除</button>
                </div>
              </td>
            </tr>
            <tr v-if="filteredGroups.length === 0">
              <td colspan="3" class="px-4 py-8 text-center text-gray-400">暂无分组数据</td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- 右侧：成员列表 -->
      <div class="w-1/2 bg-white rounded-lg shadow">
        <div class="p-4 border-b flex justify-between items-center">
          <h3 class="font-medium text-gray-900">
            {{ selectedGroup ? `「${selectedGroup.name}」的成员` : '请选择分组' }}
          </h3>
          <button v-if="selectedGroup" @click="showAddMemberDialog = true"
            class="px-3 py-1.5 bg-blue-600 text-white text-sm rounded-lg hover:bg-blue-700">
            添加成员
          </button>
        </div>
        <table v-if="selectedGroup" class="min-w-full">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-4 py-3 text-left text-sm font-medium text-gray-500">账户名称</th>
              <th class="px-4 py-3 text-left text-sm font-medium text-gray-500">账户ID</th>
              <th class="px-4 py-3 text-left text-sm font-medium text-gray-500">操作</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr v-for="member in members" :key="member.id">
              <td class="px-4 py-3 text-sm">{{ member.name }}</td>
              <td class="px-4 py-3 text-sm text-gray-500">{{ member.local_account_id }}</td>
              <td class="px-4 py-3">
                <button @click="handleRemoveMember(member)" class="text-red-600 hover:text-red-800 text-sm">移除</button>
              </td>
            </tr>
            <tr v-if="members.length === 0">
              <td colspan="3" class="px-4 py-8 text-center text-gray-400">暂无成员</td>
            </tr>
          </tbody>
        </table>
        <div v-else class="p-8 text-center text-gray-400">
          请在左侧选择一个分组查看成员
        </div>
      </div>
    </div>

    <!-- 新建分组弹窗 -->
    <div v-if="showCreateDialog" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg shadow-xl w-96 p-6">
        <h3 class="text-lg font-medium mb-4">新建分组</h3>
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">分组名称</label>
            <input v-model="createForm.name" type="text" placeholder="请输入分组名称"
              class="w-full border border-gray-300 rounded-lg px-3 py-2">
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">分组类型</label>
            <select v-model="createForm.type" class="w-full border border-gray-300 rounded-lg px-3 py-2">
              <option value="franchisee">加盟商</option>
              <option value="region">区域</option>
              <option value="custom">自定义</option>
            </select>
          </div>
        </div>
        <div class="flex justify-end space-x-3 mt-6">
          <button @click="showCreateDialog = false" class="px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50">取消</button>
          <button @click="handleCreate" class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700">确定</button>
        </div>
      </div>
    </div>

    <!-- 添加成员弹窗 -->
    <div v-if="showAddMemberDialog" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg shadow-xl w-96 p-6">
        <h3 class="text-lg font-medium mb-4">添加成员</h3>
        <div class="max-h-64 overflow-y-auto border rounded-lg">
          <label v-for="account in availableAccounts" :key="account.id"
            class="flex items-center px-4 py-2 hover:bg-gray-50 cursor-pointer">
            <input type="checkbox" v-model="selectedAccountIds" :value="account.id"
              class="mr-3 rounded border-gray-300">
            <span class="text-sm">{{ account.name }}</span>
            <span class="text-xs text-gray-400 ml-2">ID: {{ account.local_account_id }}</span>
          </label>
          <div v-if="availableAccounts.length === 0" class="px-4 py-8 text-center text-gray-400">
            暂无可添加的账户
          </div>
        </div>
        <div class="flex justify-end space-x-3 mt-6">
          <button @click="showAddMemberDialog = false" class="px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50">取消</button>
          <button @click="handleAddMembers" :disabled="selectedAccountIds.length === 0"
            class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed">
            确定
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import Breadcrumb from '@/components/common/Breadcrumb.vue'
import { groupApi, type Group } from '@/api/group'
import { accountApi, type Account } from '@/api/account'

const groups = ref<Group[]>([])
const members = ref<Account[]>([])
const allAccounts = ref<Account[]>([])
const selectedGroup = ref<Group | null>(null)
const typeFilter = ref('')
const showCreateDialog = ref(false)
const showAddMemberDialog = ref(false)
const selectedAccountIds = ref<number[]>([])

const createForm = ref({
  name: '',
  type: 'custom' as 'franchisee' | 'region' | 'custom'
})

const filteredGroups = computed(() => {
  if (!typeFilter.value) return groups.value
  return groups.value.filter(g => g.type === typeFilter.value)
})

const availableAccounts = computed(() => {
  const memberIds = new Set(members.value.map(m => m.id))
  return allAccounts.value.filter(a => !memberIds.has(a.id))
})

const typeLabel = (type: string): string => {
  const map: Record<string, string> = { franchisee: '加盟商', region: '区域', custom: '自定义' }
  return map[type] || type
}

const typeTagClass = (type: string): string => {
  const map: Record<string, string> = {
    franchisee: 'bg-purple-100 text-purple-800',
    region: 'bg-green-100 text-green-800',
    custom: 'bg-gray-100 text-gray-800'
  }
  return map[type] || 'bg-gray-100 text-gray-800'
}

const fetchGroups = async () => {
  try {
    const res = await groupApi.getList(typeFilter.value ? { type: typeFilter.value } : undefined)
    groups.value = (res as unknown as Group[]) || []
  } catch {
    groups.value = []
  }
}

const fetchAccounts = async () => {
  try {
    const res = await accountApi.getList({ page: 1, page_size: 1000 })
    allAccounts.value = (res as any)?.list || []
  } catch {
    allAccounts.value = []
  }
}

const selectGroup = async (group: Group) => {
  selectedGroup.value = group
  // 成员列表暂时从 allAccounts 中模拟，实际应有独立接口
  members.value = []
}

const handleCreate = async () => {
  if (!createForm.value.name.trim()) return
  try {
    await groupApi.create({
      name: createForm.value.name,
      type: createForm.value.type
    })
    showCreateDialog.value = false
    createForm.value = { name: '', type: 'custom' }
    await fetchGroups()
  } catch {
    // 错误处理
  }
}

const handleEdit = async (group: Group) => {
  const newName = prompt('请输入新名称', group.name)
  if (newName && newName !== group.name) {
    try {
      await groupApi.update(group.id, { name: newName })
      await fetchGroups()
    } catch {
      // 错误处理
    }
  }
}

const handleDelete = async (group: Group) => {
  if (!confirm(`确定删除分组「${group.name}」吗？`)) return
  try {
    await groupApi.remove(group.id)
    if (selectedGroup.value?.id === group.id) {
      selectedGroup.value = null
      members.value = []
    }
    await fetchGroups()
  } catch {
    // 错误处理
  }
}

const handleAddMembers = async () => {
  if (!selectedGroup.value || selectedAccountIds.value.length === 0) return
  try {
    await groupApi.addMembers(selectedGroup.value.id, { account_ids: selectedAccountIds.value })
    showAddMemberDialog.value = false
    selectedAccountIds.value = []
    await selectGroup(selectedGroup.value)
  } catch {
    // 错误处理
  }
}

const handleRemoveMember = async (member: Account) => {
  if (!selectedGroup.value) return
  if (!confirm(`确定移除成员「${member.name}」吗？`)) return
  try {
    await groupApi.removeMembers(selectedGroup.value.id, { account_ids: [member.id] })
    await selectGroup(selectedGroup.value)
  } catch {
    // 错误处理
  }
}

onMounted(() => {
  fetchGroups()
  fetchAccounts()
})
</script>
