<template>
  <div class="p-6">
    <Breadcrumb :items="[{ name: '本地推', path: '/local' }, { name: '线索管理', path: '/local/clue' }, { name: '线索详情' }]" />

    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold text-gray-900">线索详情</h1>
      <button @click="$router.back()" class="px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50">返回</button>
    </div>

    <!-- 基本信息 -->
    <div class="bg-white rounded-lg shadow p-6 mb-6">
      <h2 class="text-lg font-semibold text-gray-900 mb-4">基本信息</h2>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <span class="text-sm text-gray-500">客户姓名</span>
          <p class="font-medium">{{ clue.name }}</p>
        </div>
        <div>
          <span class="text-sm text-gray-500">联系电话</span>
          <p class="font-medium">{{ clue.phone }}</p>
        </div>
        <div>
          <span class="text-sm text-gray-500">线索来源</span>
          <p class="font-medium">{{ clue.source }}</p>
        </div>
        <div>
          <span class="text-sm text-gray-500">所属项目</span>
          <p class="font-medium">{{ clue.project }}</p>
        </div>
        <div>
          <span class="text-sm text-gray-500">当前状态</span>
          <p>
            <span :class="getStatusClass(clue.status)" class="px-2 py-1 text-xs rounded-full">
              {{ getStatusText(clue.status) }}
            </span>
          </p>
        </div>
        <div>
          <span class="text-sm text-gray-500">获取时间</span>
          <p class="font-medium">{{ clue.createTime }}</p>
        </div>
        <div>
          <span class="text-sm text-gray-500">跟进人</span>
          <p class="font-medium">{{ clue.follower || '未分配' }}</p>
        </div>
        <div>
          <span class="text-sm text-gray-500">线索ID</span>
          <p class="font-medium">{{ clue.id }}</p>
        </div>
      </div>
    </div>

    <!-- 跟进记录 -->
    <div class="bg-white rounded-lg shadow p-6 mb-6">
      <h2 class="text-lg font-semibold text-gray-900 mb-4">跟进记录</h2>
      <div v-if="followRecords.length === 0" class="text-center text-gray-400 py-8">暂无跟进记录</div>
      <div v-else class="space-y-4">
        <div v-for="record in followRecords" :key="record.id" class="border-l-2 border-blue-400 pl-4 py-2">
          <div class="flex justify-between items-center">
            <span class="font-medium text-gray-900">{{ record.operator }}</span>
            <span class="text-sm text-gray-500">{{ record.time }}</span>
          </div>
          <p class="text-sm text-gray-600 mt-1">{{ record.content }}</p>
        </div>
      </div>
    </div>

    <!-- 操作按钮 -->
    <div class="bg-white rounded-lg shadow p-6">
      <div class="flex space-x-4">
        <button @click="handleFollow" class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700">添加跟进</button>
        <button @click="handleAssign" class="px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50">分配跟进人</button>
        <button @click="handleMarkConverted" class="px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700">标记成交</button>
        <button @click="handleMarkInvalid" class="px-4 py-2 border border-red-300 text-red-600 rounded-lg hover:bg-red-50">标记无效</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRoute } from 'vue-router'
import Breadcrumb from '@/components/common/Breadcrumb.vue'

const route = useRoute()
const clueId = route.params.id

const clue = ref({
  id: clueId,
  name: '陈女士',
  phone: '137****7890',
  source: '表单提交',
  project: '健身房获客',
  status: 'new',
  createTime: '2024-03-15 09:30',
  follower: ''
})

const followRecords = ref([
  { id: 1, operator: '销售A', time: '2024-03-15 10:00', content: '首次电话联系，客户表示有兴趣，约定明天下午回访' },
  { id: 2, operator: '销售A', time: '2024-03-16 14:30', content: '二次回访，客户询问了价格方案，已发送报价单' }
])

const getStatusClass = (status: string) => {
  const classes: Record<string, string> = {
    new: 'bg-blue-100 text-blue-800',
    following: 'bg-orange-100 text-orange-800',
    converted: 'bg-green-100 text-green-800',
    invalid: 'bg-gray-100 text-gray-800'
  }
  return classes[status] || 'bg-gray-100 text-gray-800'
}

const getStatusText = (status: string) => {
  const texts: Record<string, string> = {
    new: '新线索',
    following: '跟进中',
    converted: '已成交',
    invalid: '无效'
  }
  return texts[status] || status
}

const handleFollow = () => {
  // TODO: implement
}

const handleAssign = () => {
  // TODO: implement
}

const handleMarkConverted = () => {
  if (confirm('确定将此线索标记为已成交吗？')) {
    clue.value.status = 'converted'
  }
}

const handleMarkInvalid = () => {
  if (confirm('确定将此线索标记为无效吗？')) {
    clue.value.status = 'invalid'
  }
}
</script>
