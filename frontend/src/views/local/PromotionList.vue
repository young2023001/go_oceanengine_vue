<template>
  <div class="p-6">
    <Breadcrumb :items="[{ name: '本地推', path: '/local' }, { name: '推广管理' }]" />

    <div class="flex justify-between items-center mb-6">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">推广管理</h1>
        <p class="text-gray-600 mt-1">管理本地推广计划</p>
      </div>
      <button @click="handleCreate" class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700">
        新建推广
      </button>
    </div>

    <!-- 筛选 -->
    <div class="bg-white rounded-lg shadow p-4 mb-6">
      <div class="flex flex-wrap gap-4">
        <select v-model="filters.projectId" class="border border-gray-300 rounded-lg px-3 py-2">
          <option :value="undefined">全部项目</option>
          <option v-for="p in projectOptions" :key="p.id" :value="p.id">{{ p.name }}</option>
        </select>
        <button @click="handleSearch" class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700">搜索</button>
      </div>
    </div>

    <!-- 加载/错误状态 -->
    <div v-if="error" class="bg-red-50 text-red-600 p-4 rounded-lg mb-6">{{ error }}</div>

    <!-- 推广列表 -->
    <div class="bg-white rounded-lg shadow">
      <div v-if="loading" class="p-12 text-center text-gray-400">加载中...</div>
      <template v-else>
        <table class="min-w-full">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-4 py-3 text-left text-sm font-medium text-gray-500">推广信息</th>
              <th class="px-4 py-3 text-left text-sm font-medium text-gray-500">所属项目</th>
              <th class="px-4 py-3 text-left text-sm font-medium text-gray-500">状态</th>
              <th class="px-4 py-3 text-left text-sm font-medium text-gray-500">创建时间</th>
              <th class="px-4 py-3 text-left text-sm font-medium text-gray-500">操作</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr v-for="promo in promotions" :key="promo.id">
              <td class="px-4 py-3">
                <div class="font-medium">{{ promo.name }}</div>
                <div class="text-sm text-gray-500">ID: {{ promo.id }}</div>
              </td>
              <td class="px-4 py-3 text-sm text-gray-500">{{ promo.project_id }}</td>
              <td class="px-4 py-3">
                <span :class="getStatusClass(promo.status)" class="px-2 py-1 text-xs rounded-full">
                  {{ getStatusText(promo.status) }}
                </span>
              </td>
              <td class="px-4 py-3 text-sm text-gray-500">{{ promo.created_at }}</td>
              <td class="px-4 py-3">
                <div class="flex space-x-2">
                  <router-link :to="`/local/promotion/${promo.id}`" class="text-blue-600 hover:text-blue-800 text-sm">详情</router-link>
                </div>
              </td>
            </tr>
            <tr v-if="promotions.length === 0">
              <td colspan="5" class="px-4 py-12 text-center text-gray-400">暂无数据</td>
            </tr>
          </tbody>
        </table>
        <div class="p-4 border-t">
          <Pagination :current="pagination.page" :total="pagination.total" :page-size="pagination.pageSize" @change="handlePageChange" />
        </div>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Breadcrumb from '@/components/common/Breadcrumb.vue'
import Pagination from '@/components/common/Pagination.vue'
import { projectApi, type Promotion, type Project } from '@/api/project'

const router = useRouter()
const loading = ref(false)
const error = ref('')

const filters = ref<{ projectId: number | undefined }>({
  projectId: undefined
})

const pagination = ref({
  page: 1,
  pageSize: 20,
  total: 0
})

const promotions = ref<Promotion[]>([])
const projectOptions = ref<Project[]>([])

const getStatusClass = (status: string): string => {
  const classes: Record<string, string> = {
    enable: 'bg-green-100 text-green-800',
    disable: 'bg-orange-100 text-orange-800',
    ended: 'bg-gray-100 text-gray-800'
  }
  return classes[status] ?? 'bg-gray-100 text-gray-800'
}

const getStatusText = (status: string): string => {
  const texts: Record<string, string> = {
    enable: '投放中',
    disable: '已暂停',
    ended: '已结束'
  }
  return texts[status] ?? status
}

const fetchPromotions = async () => {
  loading.value = true
  error.value = ''

  try {
    const res = await projectApi.getPromotions({
      project_id: filters.value.projectId,
      page: pagination.value.page,
      page_size: pagination.value.pageSize
    })

    const data = res as any
    promotions.value = data.list ?? []
    pagination.value.total = data.total ?? 0
  } catch (e: unknown) {
    error.value = e instanceof Error ? e.message : '获取推广列表失败'
  } finally {
    loading.value = false
  }
}

const fetchProjectOptions = async () => {
  try {
    const res = await projectApi.getList({ page: 1, page_size: 200 })
    const data = res as any
    projectOptions.value = data.list ?? []
  } catch {
    // 静默处理，不影响主流程
  }
}

const handleSearch = () => {
  pagination.value.page = 1
  fetchPromotions()
}

const handlePageChange = (page: number) => {
  pagination.value.page = page
  fetchPromotions()
}

const handleCreate = () => {
  router.push('/local/promotion/create')
}

onMounted(() => {
  fetchProjectOptions()
  fetchPromotions()
})
</script>
