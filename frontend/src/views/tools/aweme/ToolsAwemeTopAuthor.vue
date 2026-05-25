<script setup lang="ts">
import { ref } from 'vue'
import Breadcrumb from '@/components/common/Breadcrumb.vue'

const category = ref('all')

const topAuthors = ref([
  { rank: 1, name: '李子柒', avatar: '👩', followers: 56800000, likes: 12500000, category: '生活', trend: 'up' },
  { rank: 2, name: '刘畊宏', avatar: '🏋️', followers: 48900000, likes: 9800000, category: '健身', trend: 'up' },
  { rank: 3, name: '疯狂小杨哥', avatar: '👨', followers: 45600000, likes: 8900000, category: '娱乐', trend: 'same' },
  { rank: 4, name: '东方甄选', avatar: '📚', followers: 42300000, likes: 7600000, category: '知识', trend: 'down' },
  { rank: 5, name: '张同学', avatar: '🧑', followers: 38900000, likes: 6500000, category: '生活', trend: 'up' }
])

const formatNumber = (num: number) => {
  if (num >= 100000000) return (num / 100000000).toFixed(1) + '亿'
  if (num >= 10000) return (num / 10000).toFixed(0) + '万'
  return num.toLocaleString()
}

const getTrendIcon = (trend: string) => {
  if (trend === 'up') return '📈'
  if (trend === 'down') return '📉'
  return '➡️'
}

const handleViewDetail = (_author: typeof topAuthors.value[0]) => {
  // TODO: implement
}
</script>

<template>
  <div class="space-y-6">
    <div>
      <Breadcrumb :items="[{ name: '抖音工具' }, { name: '热门达人' }]" />
      <h1 class="text-3xl font-bold text-gray-900">热门达人榜</h1>
      <p class="mt-2 text-gray-600">发现抖音平台热门创作者</p>
    </div>

    <div class="bg-white rounded-lg border border-gray-200 p-4">
      <div class="flex gap-2">
        <button :class="['px-4 py-2 rounded-lg', category === 'all' ? 'bg-blue-600 text-white' : 'bg-gray-100']"
                @click="category = 'all'">全部</button>
        <button :class="['px-4 py-2 rounded-lg', category === 'life' ? 'bg-blue-600 text-white' : 'bg-gray-100']"
                @click="category = 'life'">生活</button>
        <button :class="['px-4 py-2 rounded-lg', category === 'entertainment' ? 'bg-blue-600 text-white' : 'bg-gray-100']"
                @click="category = 'entertainment'">娱乐</button>
        <button :class="['px-4 py-2 rounded-lg', category === 'knowledge' ? 'bg-blue-600 text-white' : 'bg-gray-100']"
                @click="category = 'knowledge'">知识</button>
      </div>
    </div>

    <div class="space-y-4">
      <div v-for="author in topAuthors" :key="author.rank"
           class="bg-white rounded-lg border border-gray-200 p-4 flex items-center gap-4">
        <div :class="['w-12 h-12 rounded-full flex items-center justify-center font-bold text-xl',
                      author.rank === 1 ? 'bg-yellow-100 text-yellow-700' :
                      author.rank === 2 ? 'bg-gray-100 text-gray-700' :
                      author.rank === 3 ? 'bg-orange-100 text-orange-700' : 'bg-gray-50 text-gray-600']">
          {{ author.rank }}
        </div>
        <div class="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center text-3xl">
          {{ author.avatar }}
        </div>
        <div class="flex-1">
          <h4 class="font-semibold text-gray-900 text-lg">{{ author.name }}</h4>
          <span class="px-2 py-0.5 bg-blue-100 text-blue-700 rounded text-xs">{{ author.category }}</span>
        </div>
        <div class="text-center">
          <p class="text-xl font-bold text-gray-900">{{ formatNumber(author.followers) }}</p>
          <p class="text-xs text-gray-500">粉丝数</p>
        </div>
        <div class="text-center">
          <p class="text-xl font-bold text-pink-600">{{ formatNumber(author.likes) }}</p>
          <p class="text-xs text-gray-500">获赞数</p>
        </div>
        <div class="text-2xl">{{ getTrendIcon(author.trend) }}</div>
<button @click="handleViewDetail(author)" class="px-4 py-2 text-sm text-blue-600 border border-blue-300 rounded hover:bg-blue-50">
          查看详情
        </button>
      </div>
    </div>
  </div>
</template>
