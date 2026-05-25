<template>
  <div class="p-6">
    <Breadcrumb :items="[{ name: '本地推', path: '/local' }, { name: '项目管理' }, { name: '项目列表' }]" />

    <div class="flex justify-between items-center mb-6">
      <div>
        <h1 class="text-2xl font-bold text-gray-900">项目列表</h1>
        <p class="text-gray-600 mt-1">管理本地推广项目</p>
      </div>
      <router-link to="/local/project/create" class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700">
        创建项目
      </router-link>
    </div>

    <!-- 筛选 -->
    <div class="bg-white rounded-lg shadow p-4 mb-6">
      <div class="flex flex-wrap gap-4">
        <select v-model="filters.status" class="border border-gray-300 rounded-lg px-3 py-2">
          <option value="">全部状态</option>
          <option value="running">投放中</option>
          <option value="pause">已暂停</option>
        </select>
        <button @click="handleSearch" class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700">搜索</button>
      </div>
    </div>

    <!-- 加载/错误状态 -->
    <div v-if="error" class="bg-red-50 text-red-600 p-4 rounded-lg mb-6">{{ error }}</div>

    <!-- 项目列表 -->
    <div class="bg-white rounded-lg shadow">
      <div v-if="loading" class="p-12 text-center text-gray-400">加载中...</div>
      <template v-else>
        <table class="min-w-full">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-4 py-3 text-left text-sm font-medium text-gray-500">项目信息</th>
              <th class="px-4 py-3 text-left text-sm font-medium text-gray-500">状态</th>
              <th class="px-4 py-3 text-left text-sm font-medium text-gray-500">创建时间</th>
              <th class="px-4 py-3 text-left text-sm font-medium text-gray-500">操作</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr v-for="project in projects" :key="project.id">
              <td class="px-4 py-3">
                <div class="font-medium">{{ project.name }}</div>
                <div class="text-sm text-gray-500">ID: {{ project.id }}</div>
              </td>
              <td class="px-4 py-3">
                <span :class="project.status === 'enable' ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'"
                  class="px-2 py-1 text-xs rounded-full">
                  {{ project.status === 'enable' ? '投放中' : '已暂停' }}
                </span>
              </td>
              <td class="px-4 py-3 text-sm text-gray-500">{{ project.created_at }}</td>
              <td class="px-4 py-3">
                <div class="flex space-x-2">
                  <router-link :to="`/local/project/${project.id}`" class="text-blue-600 hover:text-blue-800 text-sm">详情</router-link>
                  <button v-if="project.status === 'enable'" @click="handlePause(project)" class="text-orange-600 hover:text-orange-800 text-sm">暂停</button>
                  <button v-else @click="handleStart(project)" class="text-green-600 hover:text-green-800 text-sm">启动</button>
                </div>
              </td>
            </tr>
            <tr v-if="projects.length === 0">
              <td colspan="4" class="px-4 py-12 text-center text-gray-400">暂无数据</td>
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
import Breadcrumb from '@/components/common/Breadcrumb.vue'
import Pagination from '@/components/common/Pagination.vue'
import { projectApi, type Project } from '@/api/project'

const loading = ref(false)
const error = ref('')

const filters = ref({
  status: ''
})

const pagination = ref({
  page: 1,
  pageSize: 20,
  total: 0
})

const projects = ref<Project[]>([])

const fetchProjects = async () => {
  loading.value = true
  error.value = ''

  try {
    const res = await projectApi.getList({
      status: filters.value.status || undefined,
      page: pagination.value.page,
      page_size: pagination.value.pageSize
    })

    const data = res as any
    projects.value = data.list ?? []
    pagination.value.total = data.total ?? 0
  } catch (e: unknown) {
    error.value = e instanceof Error ? e.message : '获取项目列表失败'
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.value.page = 1
  fetchProjects()
}

const handlePageChange = (page: number) => {
  pagination.value.page = page
  fetchProjects()
}

const handlePause = async (project: Project) => {
  if (confirm(`确定暂停项目「${project.name}」吗？`)) {
    try {
      await projectApi.updateStatus(project.id, 'disable')
      project.status = 'disable'
    } catch (e: unknown) {
      error.value = e instanceof Error ? e.message : '操作失败'
    }
  }
}

const handleStart = async (project: Project) => {
  try {
    await projectApi.updateStatus(project.id, 'enable')
    project.status = 'enable'
  } catch (e: unknown) {
    error.value = e instanceof Error ? e.message : '操作失败'
  }
}

onMounted(() => {
  fetchProjects()
})
</script>
