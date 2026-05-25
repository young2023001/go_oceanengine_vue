<template>
  <div class="p-6">
    <Breadcrumb :items="[{ name: '本地推', path: '/local' }, { name: '素材管理' }, { name: '视频管理' }]" />
    
    <div class="flex justify-between items-center mb-6">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">视频管理</h1>
        <p class="text-gray-600 mt-1">管理本地推广告素材视频</p>
      </div>
      <button @click="showUploadModal = true" class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 flex items-center">
        <span class="mr-2">+</span> 上传视频
      </button>
    </div>

    <!-- 筛选条件 -->
    <div class="bg-white rounded-lg shadow p-4 mb-6">
      <div class="flex flex-wrap gap-4 items-end">
        <div>
          <label class="block text-sm text-gray-600 mb-1">视频名称</label>
          <input type="text" v-model="filters.keyword" placeholder="请输入视频名称" class="border border-gray-300 rounded px-3 py-2 w-48">
        </div>
        <div>
          <label class="block text-sm text-gray-600 mb-1">审核状态</label>
          <select v-model="filters.auditStatus" class="border border-gray-300 rounded px-3 py-2 w-32">
            <option value="">全部状态</option>
            <option value="pending">待审核</option>
            <option value="approved">已通过</option>
            <option value="rejected">已拒绝</option>
          </select>
        </div>
        <div>
          <label class="block text-sm text-gray-600 mb-1">上传时间</label>
          <div class="flex items-center space-x-2">
            <input type="date" v-model="filters.startDate" class="border border-gray-300 rounded px-3 py-2">
            <span>至</span>
            <input type="date" v-model="filters.endDate" class="border border-gray-300 rounded px-3 py-2">
          </div>
        </div>
        <button @click="handleSearch" class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700">查询</button>
        <button @click="handleReset" class="px-4 py-2 border border-gray-300 rounded hover:bg-gray-50">重置</button>
      </div>
    </div>

    <!-- 视频列表 -->
    <div class="bg-white rounded-lg shadow">
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 p-4">
        <div v-for="video in videos" :key="video.id" class="border rounded-lg overflow-hidden hover:shadow-md transition-shadow">
          <div class="aspect-video bg-gray-100 relative group cursor-pointer">
            <div class="absolute inset-0 flex items-center justify-center">
              <span class="text-5xl">🎬</span>
            </div>
            <div class="absolute inset-0 bg-black bg-opacity-0 group-hover:bg-opacity-30 transition-all flex items-center justify-center">
              <span class="text-white opacity-0 group-hover:opacity-100 text-3xl">▶</span>
            </div>
            <div class="absolute bottom-2 right-2 bg-black bg-opacity-70 text-white text-xs px-2 py-1 rounded">
              {{ video.duration }}
            </div>
          </div>
          <div class="p-3">
            <div class="font-medium text-sm line-clamp-2 mb-2" :title="video.name">{{ video.name }}</div>
            <div class="flex items-center justify-between text-xs text-gray-500 mb-2">
              <span>{{ video.size }}</span>
              <span>{{ video.resolution }}</span>
            </div>
            <div class="flex items-center justify-between">
              <span :class="getStatusClass(video.auditStatus)" class="text-xs px-2 py-1 rounded">
                {{ getStatusText(video.auditStatus) }}
              </span>
              <span class="text-xs text-gray-400">{{ video.uploadTime }}</span>
            </div>
            <div class="mt-3 pt-3 border-t flex justify-end space-x-3">
              <button @click="handlePreview(video)" class="text-blue-600 text-xs hover:underline">预览</button>
              <button @click="handleUse(video)" class="text-blue-600 text-xs hover:underline">使用</button>
              <button @click="handleDelete(video)" class="text-red-600 text-xs hover:underline">删除</button>
            </div>
          </div>
        </div>
      </div>

      <!-- 空状态 -->
      <div v-if="videos.length === 0" class="py-16 text-center">
        <div class="text-5xl mb-4">📹</div>
        <div class="text-gray-500 mb-4">暂无视频素材</div>
        <button @click="showUploadModal = true" class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700">上传视频</button>
      </div>

      <div class="p-4 border-t">
        <Pagination :current="1" :total="24" :page-size="8" />
      </div>
    </div>

    <!-- 上传弹窗 -->
    <div v-if="showUploadModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg w-[600px] max-h-[80vh] overflow-y-auto">
        <div class="flex justify-between items-center p-4 border-b">
          <h3 class="text-lg font-medium">上传视频</h3>
          <button @click="showUploadModal = false" class="text-gray-400 hover:text-gray-600">✕</button>
        </div>
        <div class="p-6">
<div class="border-2 border-dashed border-gray-300 rounded-lg p-8 text-center hover:border-blue-500 cursor-pointer" @click="handleSelectFile">
            <div class="text-4xl mb-4">📤</div>
            <div class="text-gray-600 mb-2">点击或拖拽视频到此处上传</div>
            <div class="text-xs text-gray-400">支持 MP4、MOV 格式，最大 500MB</div>
          </div>
          <div class="mt-4">
            <label class="block text-sm text-gray-700 mb-2">视频名称</label>
            <input v-model="uploadForm.name" type="text" class="w-full border border-gray-300 rounded px-3 py-2" placeholder="请输入视频名称">
          </div>
        </div>
        <div class="flex justify-end space-x-3 p-4 border-t bg-gray-50">
          <button @click="showUploadModal = false" class="px-4 py-2 border border-gray-300 rounded hover:bg-gray-50">取消</button>
          <button @click="handleUpload" :disabled="uploading" class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50">{{ uploading ? '上传中...' : '确认上传' }}</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import Breadcrumb from '@/components/common/Breadcrumb.vue'
import Pagination from '@/components/common/Pagination.vue'

const filters = ref({
  keyword: '',
  auditStatus: '',
  startDate: '',
  endDate: ''
})

const showUploadModal = ref(false)
const uploading = ref(false)
const uploadForm = ref({
  name: '',
  file: null as File | null
})

const videos = ref([
  { id: 'V001', name: '618大促限时特惠活动宣传', duration: '00:30', size: '15.6MB', resolution: '1080x1920', auditStatus: 'approved', uploadTime: '2024-06-01' },
  { id: 'V002', name: '新品上市发布会精彩回顾', duration: '01:15', size: '45.2MB', resolution: '1080x1920', auditStatus: 'approved', uploadTime: '2024-05-28' },
  { id: 'V003', name: '店铺环境展示视频', duration: '00:45', size: '28.3MB', resolution: '1080x1920', auditStatus: 'pending', uploadTime: '2024-05-25' },
  { id: 'V004', name: '招牌菜品制作过程', duration: '02:00', size: '68.5MB', resolution: '1080x1920', auditStatus: 'approved', uploadTime: '2024-05-22' },
  { id: 'V005', name: '会员专享福利介绍', duration: '00:20', size: '12.1MB', resolution: '1080x1920', auditStatus: 'rejected', uploadTime: '2024-05-20' },
  { id: 'V006', name: '周末特惠活动预告', duration: '00:25', size: '14.8MB', resolution: '1080x1920', auditStatus: 'approved', uploadTime: '2024-05-18' },
  { id: 'V007', name: '品牌故事纪录片', duration: '03:30', size: '125.6MB', resolution: '1920x1080', auditStatus: 'pending', uploadTime: '2024-05-15' },
  { id: 'V008', name: '用户好评合集', duration: '01:00', size: '38.2MB', resolution: '1080x1920', auditStatus: 'approved', uploadTime: '2024-05-12' }
])

const getStatusClass = (status: string) => {
  const classes: Record<string, string> = {
    pending: 'bg-yellow-100 text-yellow-800',
    approved: 'bg-green-100 text-green-800',
    rejected: 'bg-red-100 text-red-800'
  }
  return classes[status] || 'bg-gray-100 text-gray-800'
}

const getStatusText = (status: string) => {
  const texts: Record<string, string> = {
    pending: '待审核',
    approved: '已通过',
    rejected: '已拒绝'
  }
  return texts[status] || '未知'
}

const handleSearch = () => {
  // TODO: implement
}

const handleReset = () => {
  filters.value = { keyword: '', auditStatus: '', startDate: '', endDate: '' }
}

const handlePreview = (_video: typeof videos.value[0]) => {
  // TODO: implement
}

const handleUse = (_video: typeof videos.value[0]) => {
  // TODO: implement
}

const handleDelete = (video: typeof videos.value[0]) => {
  if (confirm(`确定删除视频「${video.name}」吗？`)) {
    const idx = videos.value.findIndex(v => v.id === video.id)
    if (idx > -1) videos.value.splice(idx, 1)
  }
}

const handleSelectFile = () => {
  // TODO: implement
}

const handleUpload = async () => {
  uploading.value = true
  try {
    await new Promise(r => setTimeout(r, 1000))
    showUploadModal.value = false
    uploadForm.value = { name: '', file: null }
  } finally {
    uploading.value = false
  }
}
</script>
