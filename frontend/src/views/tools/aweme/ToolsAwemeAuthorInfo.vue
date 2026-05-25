<script setup lang="ts">
import { ref } from 'vue'
import Breadcrumb from '@/components/common/Breadcrumb.vue'

const authorInfo = ref({
  name: '科技测评达人',
  avatar: '👨',
  douyinId: 'tech_review_123',
  followers: 2560000,
  following: 156,
  likes: 8900000,
  videos: 234,
  description: '专注数码科技评测，每周更新3-5条内容',
  category: '数码科技',
  region: '北京',
  verified: true
})

const recentVideos = ref([
  { id: 'V001', title: 'iPhone 15 Pro深度体验', plays: 1250000, likes: 89000, comments: 5600 },
  { id: 'V002', title: '小米14开箱评测', plays: 980000, likes: 67000, comments: 4200 },
  { id: 'V003', title: '华为Mate60实测', plays: 1560000, likes: 112000, comments: 8900 }
])

const formatNumber = (num: number) => {
  if (num >= 10000) return (num / 10000).toFixed(1) + '万'
  return num.toLocaleString()
}

const handleAddTarget = () => {
  // TODO: implement
}
</script>

<template>
  <div class="space-y-6">
    <div>
      <Breadcrumb :items="[{ name: '抖音工具' }, { name: '达人信息' }]" />
    </div>

    <div class="bg-white rounded-lg border border-gray-200 p-6">
      <div class="flex items-start gap-6">
        <div class="w-24 h-24 bg-gray-100 rounded-full flex items-center justify-center text-5xl">
          {{ authorInfo.avatar }}
        </div>
        <div class="flex-1">
          <div class="flex items-center gap-3">
            <h1 class="text-2xl font-bold text-gray-900">{{ authorInfo.name }}</h1>
            <span v-if="authorInfo.verified" class="px-2 py-0.5 bg-blue-100 text-blue-700 rounded text-xs">已认证</span>
            <span class="px-2 py-0.5 bg-gray-100 text-gray-700 rounded text-xs">{{ authorInfo.category }}</span>
          </div>
          <p class="text-gray-500 mt-1">抖音号: {{ authorInfo.douyinId }}</p>
          <p class="text-gray-600 mt-2">{{ authorInfo.description }}</p>
          <div class="flex items-center gap-6 mt-4">
            <div><span class="text-xl font-bold text-gray-900">{{ formatNumber(authorInfo.followers) }}</span><span class="text-gray-500 text-sm ml-1">粉丝</span></div>
            <div><span class="text-xl font-bold text-gray-900">{{ authorInfo.following }}</span><span class="text-gray-500 text-sm ml-1">关注</span></div>
            <div><span class="text-xl font-bold text-pink-600">{{ formatNumber(authorInfo.likes) }}</span><span class="text-gray-500 text-sm ml-1">获赞</span></div>
            <div><span class="text-xl font-bold text-blue-600">{{ authorInfo.videos }}</span><span class="text-gray-500 text-sm ml-1">作品</span></div>
          </div>
        </div>
<button @click="handleAddTarget" class="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700">
          添加定向
        </button>
      </div>
    </div>

    <div class="bg-white rounded-lg border border-gray-200 p-6">
      <h3 class="font-semibold text-gray-900 mb-4">近期热门作品</h3>
      <div class="space-y-3">
        <div v-for="video in recentVideos" :key="video.id"
             class="flex items-center gap-4 p-3 border border-gray-100 rounded-lg hover:bg-gray-50">
          <div class="w-16 h-16 bg-gray-200 rounded flex items-center justify-center text-2xl">🎬</div>
          <div class="flex-1">
            <h4 class="font-medium text-gray-900">{{ video.title }}</h4>
            <div class="flex gap-4 text-sm text-gray-500 mt-1">
              <span>播放 {{ formatNumber(video.plays) }}</span>
              <span>点赞 {{ formatNumber(video.likes) }}</span>
              <span>评论 {{ formatNumber(video.comments) }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
