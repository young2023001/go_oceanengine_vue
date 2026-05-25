<template>
  <div class="p-6">
    <Breadcrumb :items="[{ name: '服务市场', path: '/servemarket' }, { name: '订阅管理' }]" />
    
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-gray-900">订阅管理</h1>
      <p class="text-gray-600 mt-1">管理已订阅的服务</p>
    </div>

    <div class="mb-6">
      <div class="flex border-b">
        <button v-for="tab in tabs" :key="tab.key" @click="currentTab = tab.key"
          :class="currentTab === tab.key ? 'border-b-2 border-blue-600 text-blue-600' : 'text-gray-500'"
          class="px-4 py-2 font-medium">{{ tab.label }}</button>
      </div>
    </div>

    <div v-if="currentTab === 'active'" class="space-y-4">
      <div v-for="sub in activeSubscriptions" :key="sub.id" class="bg-white rounded-lg shadow p-4">
        <div class="flex items-center justify-between">
          <div class="flex items-center">
            <div class="w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center mr-4">
              <span class="text-2xl">{{ sub.icon }}</span>
            </div>
            <div>
              <h4 class="font-medium">{{ sub.name }}</h4>
              <p class="text-sm text-gray-500">{{ sub.desc }}</p>
            </div>
          </div>
          <div class="text-right">
            <div class="text-sm text-gray-500">到期时间: {{ sub.expireDate }}</div>
            <div class="flex space-x-2 mt-2">
              <button @click="handleRenew(sub)" class="px-3 py-1 text-sm text-blue-600 border border-blue-600 rounded hover:bg-blue-50">续费</button>
              <button @click="handleCancelSub(sub)" class="px-3 py-1 text-sm text-gray-600 border border-gray-300 rounded hover:bg-gray-50">取消</button>
            </div>
          </div>
        </div>
        <div class="mt-4 flex items-center text-sm text-gray-500">
          <span>订阅时间: {{ sub.startDate }}</span>
          <span class="mx-4">|</span>
          <span>价格: ¥{{ sub.price }}/月</span>
          <span class="mx-4">|</span>
          <span>自动续费: {{ sub.autoRenew ? '已开启' : '已关闭' }}</span>
        </div>
      </div>
    </div>

    <div v-if="currentTab === 'expired'" class="space-y-4">
      <div v-for="sub in expiredSubscriptions" :key="sub.id" class="bg-white rounded-lg shadow p-4 opacity-75">
        <div class="flex items-center justify-between">
          <div class="flex items-center">
            <div class="w-12 h-12 bg-gray-100 rounded-lg flex items-center justify-center mr-4">
              <span class="text-2xl grayscale">{{ sub.icon }}</span>
            </div>
            <div>
              <h4 class="font-medium text-gray-500">{{ sub.name }}</h4>
              <p class="text-sm text-gray-400">已于 {{ sub.expireDate }} 到期</p>
            </div>
          </div>
          <button @click="handleResubscribe(sub)" class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700">重新订阅</button>
        </div>
      </div>
    </div>

    <div v-if="currentTab === 'all'" class="bg-white rounded-lg shadow">
      <table class="min-w-full">
        <thead class="bg-gray-50">
          <tr>
            <th class="px-4 py-3 text-left text-sm font-medium text-gray-500">服务</th>
            <th class="px-4 py-3 text-left text-sm font-medium text-gray-500">状态</th>
            <th class="px-4 py-3 text-left text-sm font-medium text-gray-500">到期时间</th>
            <th class="px-4 py-3 text-right text-sm font-medium text-gray-500">价格</th>
            <th class="px-4 py-3 text-left text-sm font-medium text-gray-500">操作</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-200">
          <tr v-for="sub in [...activeSubscriptions, ...expiredSubscriptions]" :key="sub.id">
            <td class="px-4 py-3 font-medium">{{ sub.name }}</td>
            <td class="px-4 py-3">
              <span :class="sub.status === 'active' ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'" class="px-2 py-1 text-xs rounded">
                {{ sub.status === 'active' ? '使用中' : '已过期' }}
              </span>
            </td>
            <td class="px-4 py-3 text-sm text-gray-500">{{ sub.expireDate }}</td>
            <td class="px-4 py-3 text-sm text-right">¥{{ sub.price }}/月</td>
            <td class="px-4 py-3">
              <button @click="handleTableAction(sub)" class="text-blue-600 hover:text-blue-800 text-sm">{{ sub.status === 'active' ? '续费' : '订阅' }}</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import Breadcrumb from '@/components/common/Breadcrumb.vue'

const currentTab = ref('active')

const handleRenew = (_sub: typeof activeSubscriptions.value[0]) => {
  // TODO: implement
}

const handleCancelSub = (sub: typeof activeSubscriptions.value[0]) => {
  if (confirm(`确定取消订阅 ${sub.name} 吗？`)) {
  }
}

const handleResubscribe = (_sub: typeof expiredSubscriptions.value[0]) => {
  // TODO: implement
}

const handleTableAction = (_sub: any) => {
  // TODO: implement
}
const tabs = [
  { key: 'active', label: '使用中' },
  { key: 'expired', label: '已过期' },
  { key: 'all', label: '全部' }
]

const activeSubscriptions = ref([
  { id: 1, name: '智能投放助手', icon: '🤖', desc: 'AI驱动的智能投放优化工具', startDate: '2024-01-01', expireDate: '2024-12-31', price: 299, autoRenew: true, status: 'active' },
  { id: 2, name: '创意素材生成', icon: '🎨', desc: '一键生成高质量创意素材', startDate: '2024-03-15', expireDate: '2024-09-15', price: 199, autoRenew: false, status: 'active' },
  { id: 3, name: '数据分析Pro', icon: '📊', desc: '深度数据洞察分析工具', startDate: '2024-02-20', expireDate: '2024-08-20', price: 399, autoRenew: true, status: 'active' }
])

const expiredSubscriptions = ref([
  { id: 4, name: '竞品分析工具', icon: '🔍', desc: '实时竞品数据监控', expireDate: '2024-05-01', price: 199, status: 'expired' },
  { id: 5, name: '自动报表生成', icon: '📋', desc: '定时自动生成数据报表', expireDate: '2024-04-15', price: 99, status: 'expired' }
])
</script>
