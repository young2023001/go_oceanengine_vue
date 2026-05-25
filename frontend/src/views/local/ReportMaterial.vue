<template>
  <div class="p-6">
    <Breadcrumb :items="[{ name: '本地推', path: '/local' }, { name: '数据报表' }, { name: '素材报表' }]" />
    
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-gray-900">素材报表</h1>
      <p class="text-gray-600 mt-1">查看素材投放效果数据</p>
    </div>

    <!-- 筛选条件 -->
    <div class="bg-white rounded-lg shadow p-4 mb-6">
      <div class="flex flex-wrap gap-4 items-end">
        <div>
          <label class="block text-sm text-gray-600 mb-1">日期范围</label>
          <div class="flex items-center space-x-2">
            <input type="date" v-model="filters.startDate" class="border border-gray-300 rounded px-3 py-2">
            <span>至</span>
            <input type="date" v-model="filters.endDate" class="border border-gray-300 rounded px-3 py-2">
          </div>
        </div>
        <div>
          <label class="block text-sm text-gray-600 mb-1">素材类型</label>
          <select v-model="filters.type" class="border border-gray-300 rounded px-3 py-2 w-32">
            <option value="">全部类型</option>
            <option value="video">视频</option>
            <option value="image">图片</option>
          </select>
        </div>
        <button @click="handleSearch" class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700">查询</button>
        <button @click="handleExport" class="px-4 py-2 border border-gray-300 rounded hover:bg-gray-50">导出</button>
      </div>
    </div>

    <!-- 汇总数据 -->
    <div class="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-4 mb-6">
      <div class="bg-white rounded-lg shadow p-4">
        <div class="text-sm text-gray-500">素材数</div>
        <div class="text-xl font-bold text-gray-900 mt-1">{{ summary.count }}</div>
      </div>
      <div class="bg-white rounded-lg shadow p-4">
        <div class="text-sm text-gray-500">消耗</div>
        <div class="text-xl font-bold text-gray-900 mt-1">¥{{ summary.cost.toLocaleString() }}</div>
      </div>
      <div class="bg-white rounded-lg shadow p-4">
        <div class="text-sm text-gray-500">曝光量</div>
        <div class="text-xl font-bold text-gray-900 mt-1">{{ formatNumber(summary.impressions) }}</div>
      </div>
      <div class="bg-white rounded-lg shadow p-4">
        <div class="text-sm text-gray-500">点击量</div>
        <div class="text-xl font-bold text-gray-900 mt-1">{{ formatNumber(summary.clicks) }}</div>
      </div>
      <div class="bg-white rounded-lg shadow p-4">
        <div class="text-sm text-gray-500">平均点击率</div>
        <div class="text-xl font-bold text-gray-900 mt-1">{{ summary.avgCtr }}%</div>
      </div>
      <div class="bg-white rounded-lg shadow p-4">
        <div class="text-sm text-gray-500">平均转化率</div>
        <div class="text-xl font-bold text-gray-900 mt-1">{{ summary.avgCvr }}%</div>
      </div>
    </div>

    <!-- 素材列表 -->
    <div class="bg-white rounded-lg shadow">
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 p-4">
        <div v-for="item in materials" :key="item.id" class="border rounded-lg overflow-hidden">
          <div class="aspect-video bg-gray-100 flex items-center justify-center relative">
            <span class="text-4xl">{{ item.type === 'video' ? '🎬' : '🖼️' }}</span>
            <span class="absolute top-2 right-2 px-2 py-1 bg-black bg-opacity-50 text-white text-xs rounded">
              {{ item.type === 'video' ? '视频' : '图片' }}
            </span>
          </div>
          <div class="p-4">
            <div class="font-medium text-sm line-clamp-2 mb-2">{{ item.name }}</div>
            <div class="grid grid-cols-2 gap-2 text-xs">
              <div>
                <span class="text-gray-500">消耗:</span>
                <span class="font-medium ml-1">¥{{ item.cost.toLocaleString() }}</span>
              </div>
              <div>
                <span class="text-gray-500">曝光:</span>
                <span class="font-medium ml-1">{{ formatNumber(item.impressions) }}</span>
              </div>
              <div>
                <span class="text-gray-500">点击率:</span>
                <span class="font-medium ml-1">{{ item.ctr }}%</span>
              </div>
              <div>
                <span class="text-gray-500">转化率:</span>
                <span class="font-medium ml-1">{{ item.cvr }}%</span>
              </div>
            </div>
            <div class="mt-3 pt-3 border-t flex justify-between items-center">
              <span class="text-xs text-gray-400">{{ item.createdAt }}</span>
              <button @click="handleMaterialDetail(item)" class="text-blue-600 text-xs hover:underline">查看详情</button>
            </div>
          </div>
        </div>
      </div>
      <div class="p-4 border-t">
        <Pagination :current="1" :total="50" :page-size="9" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import Breadcrumb from '@/components/common/Breadcrumb.vue'
import Pagination from '@/components/common/Pagination.vue'

const handleSearch = () => {
  // TODO: implement
}

const handleExport = () => {
  // TODO: implement
}

const handleMaterialDetail = (_item: typeof materials.value[0]) => {
  // TODO: implement
}

const filters = ref({
  startDate: '',
  endDate: '',
  type: ''
})

const summary = ref({
  count: 28,
  cost: 45680,
  impressions: 1256800,
  clicks: 85600,
  avgCtr: 6.81,
  avgCvr: 0.66
})

const materials = ref([
  { id: 'M001', name: '618大促限时特惠活动视频', type: 'video', cost: 8560, impressions: 256800, ctr: 7.2, cvr: 0.85, createdAt: '2024-06-01' },
  { id: 'M002', name: '新品上市宣传视频', type: 'video', cost: 6580, impressions: 186500, ctr: 6.8, cvr: 0.72, createdAt: '2024-05-28' },
  { id: 'M003', name: '店铺环境展示图', type: 'image', cost: 3560, impressions: 125600, ctr: 5.6, cvr: 0.58, createdAt: '2024-05-25' },
  { id: 'M004', name: '招牌菜品展示', type: 'image', cost: 4280, impressions: 156800, ctr: 6.2, cvr: 0.65, createdAt: '2024-05-22' },
  { id: 'M005', name: '会员福利介绍视频', type: 'video', cost: 5680, impressions: 168500, ctr: 7.5, cvr: 0.78, createdAt: '2024-05-20' },
  { id: 'M006', name: '周末特惠活动海报', type: 'image', cost: 2860, impressions: 98600, ctr: 5.8, cvr: 0.55, createdAt: '2024-05-18' }
])

const formatNumber = (num: number) => {
  if (num >= 10000) {
    return (num / 10000).toFixed(1) + 'w'
  }
  return num.toLocaleString()
}
</script>
