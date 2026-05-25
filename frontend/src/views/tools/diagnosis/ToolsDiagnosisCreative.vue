<script setup lang="ts">
import { ref, reactive } from 'vue'
import Breadcrumb from '@/components/common/Breadcrumb.vue'
import Pagination from '@/components/common/Pagination.vue'

const pagination = reactive({ page: 1, pageSize: 10, total: 68 })

const creatives = ref([
  { id: 'CR001', name: '促销视频A', type: 'video', score: 88, ctr: 3.5, cvr: 2.1, issues: ['画质建议提升'], status: 'good' },
  { id: 'CR002', name: '产品图片B', type: 'image', score: 72, ctr: 2.8, cvr: 1.5, issues: ['文字占比过高', '色彩对比度低'], status: 'warning' },
  { id: 'CR003', name: '品牌Banner', type: 'image', score: 55, ctr: 1.2, cvr: 0.8, issues: ['点击率过低', '吸引力不足', '需要优化文案'], status: 'critical' }
])

const handlePageChange = (page: number) => {
  pagination.page = page
}

const handleUpload = () => {
  // TODO: implement
}

const handleDetail = (_creative: typeof creatives.value[0]) => {
  // TODO: implement
}

const handleSuggestion = (_creative: typeof creatives.value[0]) => {
  // TODO: implement
}
</script>

<template>
  <div class="space-y-6">
    <div>
      <Breadcrumb :items="[{ name: '诊断工具' }, { name: '创意诊断' }]" />
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-3xl font-bold text-gray-900">创意诊断</h1>
          <p class="mt-2 text-gray-600">分析广告创意效果并提供优化建议</p>
        </div>
<button class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700" @click="handleUpload">
          上传创意分析
        </button>
      </div>
    </div>

    <div class="grid grid-cols-5 gap-4">
      <div class="bg-white rounded-lg border border-gray-200 p-4">
        <p class="text-sm text-gray-500">创意总数</p>
        <p class="text-2xl font-bold text-gray-900 mt-1">{{ pagination.total }}</p>
      </div>
      <div class="bg-white rounded-lg border border-gray-200 p-4">
        <p class="text-sm text-gray-500">优秀创意</p>
        <p class="text-2xl font-bold text-green-600 mt-1">42</p>
      </div>
      <div class="bg-white rounded-lg border border-gray-200 p-4">
        <p class="text-sm text-gray-500">待优化</p>
        <p class="text-2xl font-bold text-yellow-600 mt-1">18</p>
      </div>
      <div class="bg-white rounded-lg border border-gray-200 p-4">
        <p class="text-sm text-gray-500">效果差</p>
        <p class="text-2xl font-bold text-red-600 mt-1">8</p>
      </div>
      <div class="bg-white rounded-lg border border-gray-200 p-4">
        <p class="text-sm text-gray-500">平均CTR</p>
        <p class="text-2xl font-bold text-blue-600 mt-1">2.5%</p>
      </div>
    </div>

    <div class="bg-white rounded-lg border border-gray-200">
      <table class="w-full">
        <thead class="bg-gray-50 border-b border-gray-200">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-semibold text-gray-600 uppercase">创意信息</th>
            <th class="px-6 py-3 text-left text-xs font-semibold text-gray-600 uppercase">类型</th>
            <th class="px-6 py-3 text-left text-xs font-semibold text-gray-600 uppercase">评分</th>
            <th class="px-6 py-3 text-left text-xs font-semibold text-gray-600 uppercase">CTR</th>
            <th class="px-6 py-3 text-left text-xs font-semibold text-gray-600 uppercase">CVR</th>
            <th class="px-6 py-3 text-left text-xs font-semibold text-gray-600 uppercase">问题标签</th>
            <th class="px-6 py-3 text-left text-xs font-semibold text-gray-600 uppercase">操作</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-200">
          <tr v-for="creative in creatives" :key="creative.id" class="hover:bg-gray-50">
            <td class="px-6 py-4">
              <div class="flex items-center">
                <div class="w-12 h-12 bg-gray-200 rounded flex items-center justify-center mr-3">
                  <span class="text-xl">{{ creative.type === 'video' ? '🎬' : '🖼️' }}</span>
                </div>
                <div>
                  <div class="text-sm font-medium text-gray-900">{{ creative.name }}</div>
                  <div class="text-xs text-gray-500">{{ creative.id }}</div>
                </div>
              </div>
            </td>
            <td class="px-6 py-4">
              <span :class="['px-2 py-1 rounded text-xs',
                     creative.type === 'video' ? 'bg-purple-100 text-purple-700' : 'bg-blue-100 text-blue-700']">
                {{ creative.type === 'video' ? '视频' : '图片' }}
              </span>
            </td>
            <td class="px-6 py-4">
              <span :class="['text-sm font-medium',
                     creative.score >= 80 ? 'text-green-600' : creative.score >= 60 ? 'text-yellow-600' : 'text-red-600']">
                {{ creative.score }}分
              </span>
            </td>
            <td class="px-6 py-4 text-sm text-gray-600">{{ creative.ctr }}%</td>
            <td class="px-6 py-4 text-sm text-gray-600">{{ creative.cvr }}%</td>
            <td class="px-6 py-4">
              <div class="flex flex-wrap gap-1">
                <span v-for="(issue, index) in creative.issues.slice(0, 2)" :key="index"
                      class="px-2 py-0.5 bg-red-50 text-red-600 rounded text-xs">{{ issue }}</span>
                <span v-if="creative.issues.length > 2" class="px-2 py-0.5 bg-gray-100 text-gray-600 rounded text-xs">
                  +{{ creative.issues.length - 2 }}
                </span>
              </div>
            </td>
            <td class="px-6 py-4 text-sm">
<button class="text-blue-600 hover:text-blue-800 mr-3" @click="handleDetail(creative)">详情</button>
              <button class="text-green-600 hover:text-green-800" @click="handleSuggestion(creative)">优化建议</button>
            </td>
          </tr>
        </tbody>
      </table>
      <div class="px-6 py-4 border-t border-gray-200">
        <Pagination :current="pagination.page" :total="pagination.total" :page-size="pagination.pageSize" @change="handlePageChange" />
      </div>
    </div>
  </div>
</template>
