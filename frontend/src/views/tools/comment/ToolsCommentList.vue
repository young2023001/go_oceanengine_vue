<script setup lang="ts">
import { ref, reactive } from 'vue'
import Breadcrumb from '@/components/common/Breadcrumb.vue'
import Pagination from '@/components/common/Pagination.vue'

const pagination = reactive({ page: 1, pageSize: 10, total: 156 })
const activeTab = ref('all')

const tabs = [
  { key: 'all', label: '全部', count: 156 },
  { key: 'pending', label: '待处理', count: 23 },
  { key: 'replied', label: '已回复', count: 98 },
  { key: 'hidden', label: '已隐藏', count: 35 }
]

const comments = ref([
  { id: 'C001', content: '这个产品太好用了，强烈推荐！', adName: '品牌推广A', user: '用户***8', time: '10分钟前', status: 'pending', sentiment: 'positive' },
  { id: 'C002', content: '价格有点贵啊，能便宜点吗？', adName: '促销活动B', user: '用户***2', time: '30分钟前', status: 'pending', sentiment: 'neutral' },
  { id: 'C003', content: '已经收到货了，质量不错', adName: '品牌推广A', user: '用户***5', time: '1小时前', status: 'replied', sentiment: 'positive' },
  { id: 'C004', content: '发货太慢了，差评', adName: '双11活动', user: '用户***9', time: '2小时前', status: 'pending', sentiment: 'negative' }
])

const getSentimentConfig = (sentiment: string) => {
  switch (sentiment) {
    case 'positive': return { icon: '😊', class: 'text-green-600' }
    case 'negative': return { icon: '😞', class: 'text-red-600' }
    default: return { icon: '😐', class: 'text-gray-600' }
  }
}

const handlePageChange = (page: number) => {
  pagination.page = page
}

const handleReply = (_comment: typeof comments.value[0]) => {
  // TODO: implement
}

const handleHide = (comment: typeof comments.value[0]) => {
  comment.status = 'hidden'
}
</script>

<template>
  <div class="space-y-6">
    <div>
      <Breadcrumb :items="[{ name: '评论工具' }, { name: '评论管理' }]" />
      <h1 class="text-3xl font-bold text-gray-900">评论管理</h1>
      <p class="mt-2 text-gray-600">管理广告评论，及时回复用户反馈</p>
    </div>

    <div class="grid grid-cols-1 sm:grid-cols-4 gap-4">
      <div class="bg-white rounded-lg border border-gray-200 p-4">
        <p class="text-sm text-gray-500">今日评论</p>
        <p class="text-2xl font-bold text-gray-900 mt-1">89</p>
      </div>
      <div class="bg-white rounded-lg border border-gray-200 p-4">
        <p class="text-sm text-gray-500">待处理</p>
        <p class="text-2xl font-bold text-yellow-600 mt-1">23</p>
      </div>
      <div class="bg-white rounded-lg border border-gray-200 p-4">
        <p class="text-sm text-gray-500">好评率</p>
        <p class="text-2xl font-bold text-green-600 mt-1">78%</p>
      </div>
      <div class="bg-white rounded-lg border border-gray-200 p-4">
        <p class="text-sm text-gray-500">回复率</p>
        <p class="text-2xl font-bold text-blue-600 mt-1">92%</p>
      </div>
    </div>

    <div class="bg-white rounded-lg border border-gray-200">
      <div class="border-b border-gray-200">
        <nav class="flex">
          <button v-for="tab in tabs" :key="tab.key" @click="activeTab = tab.key"
                  :class="['px-6 py-3 text-sm font-medium border-b-2 -mb-px',
                    activeTab === tab.key ? 'border-blue-500 text-blue-600' : 'border-transparent text-gray-500 hover:text-gray-700']">
            {{ tab.label }} ({{ tab.count }})
          </button>
        </nav>
      </div>
      
      <div class="divide-y divide-gray-200">
        <div v-for="comment in comments" :key="comment.id" class="p-4 hover:bg-gray-50">
          <div class="flex items-start gap-4">
            <div class="w-10 h-10 bg-gray-200 rounded-full flex items-center justify-center text-gray-500">
              {{ comment.user.slice(-2, -1) }}
            </div>
            <div class="flex-1">
              <div class="flex items-center gap-2 mb-1">
                <span class="font-medium text-gray-900">{{ comment.user }}</span>
                <span class="text-xs text-gray-400">{{ comment.time }}</span>
                <span :class="['text-sm', getSentimentConfig(comment.sentiment).class]">
                  {{ getSentimentConfig(comment.sentiment).icon }}
                </span>
              </div>
              <p class="text-gray-700">{{ comment.content }}</p>
              <div class="flex items-center gap-4 mt-2">
                <span class="text-xs text-gray-500">来自: {{ comment.adName }}</span>
<button @click="handleReply(comment)" class="text-blue-600 text-xs hover:text-blue-800">回复</button>
                <button @click="handleHide(comment)" class="text-gray-500 text-xs hover:text-gray-700">隐藏</button>
              </div>
            </div>
            <span :class="['px-2 py-1 rounded text-xs',
              comment.status === 'pending' ? 'bg-yellow-100 text-yellow-700' :
              comment.status === 'replied' ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-700']">
              {{ comment.status === 'pending' ? '待处理' : comment.status === 'replied' ? '已回复' : '已隐藏' }}
            </span>
          </div>
        </div>
      </div>
      
      <div class="px-6 py-4 border-t border-gray-200">
        <Pagination :current="pagination.page" :total="pagination.total" :page-size="pagination.pageSize" @change="handlePageChange" />
      </div>
    </div>
  </div>
</template>
