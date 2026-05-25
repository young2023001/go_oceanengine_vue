<template>
  <div class="p-6">
    <Breadcrumb :items="[{ name: '本地推', path: '/local' }, { name: '门店列表' }]" />

    <div class="mb-6">
      <h1 class="text-2xl font-bold text-gray-900">门店列表</h1>
      <p class="text-gray-600 mt-1">查看已关联的门店信息</p>
    </div>

    <div class="bg-white rounded-lg shadow">
      <table class="min-w-full">
        <thead class="bg-gray-50">
          <tr>
            <th class="px-4 py-3 text-left text-sm font-medium text-gray-500">门店名称</th>
            <th class="px-4 py-3 text-left text-sm font-medium text-gray-500">所属账户</th>
            <th class="px-4 py-3 text-left text-sm font-medium text-gray-500">城市</th>
            <th class="px-4 py-3 text-left text-sm font-medium text-gray-500">区域</th>
            <th class="px-4 py-3 text-left text-sm font-medium text-gray-500">地址</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-200">
          <tr v-for="store in stores" :key="store.id">
            <td class="px-4 py-3 text-sm font-medium">{{ store.name }}</td>
            <td class="px-4 py-3 text-sm">{{ store.accountName }}</td>
            <td class="px-4 py-3 text-sm">{{ store.city }}</td>
            <td class="px-4 py-3 text-sm">{{ store.region }}</td>
            <td class="px-4 py-3 text-sm text-gray-500">{{ store.address }}</td>
          </tr>
        </tbody>
      </table>

      <!-- 空状态 -->
      <div v-if="stores.length === 0" class="p-12 text-center">
        <div class="text-gray-400 text-lg mb-2">暂无门店数据</div>
        <p class="text-gray-400 text-sm">门店数据将在账户导入时自动关联</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import Breadcrumb from '@/components/common/Breadcrumb.vue'

interface StoreItem {
  id: number
  name: string
  accountName: string
  city: string
  region: string
  address: string
}

const stores = ref<StoreItem[]>([])
</script>
