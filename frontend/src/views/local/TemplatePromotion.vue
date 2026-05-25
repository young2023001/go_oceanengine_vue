<script setup lang="ts">
import { ref, onMounted } from 'vue'
import Breadcrumb from '@/components/common/Breadcrumb.vue'
import { templateApi, type PromotionTemplate, type CreateTemplateParams } from '@/api/template'

const templates = ref<PromotionTemplate[]>([])
const loading = ref(false)
const showDialog = ref(false)
const form = ref<CreateTemplateParams>({ name: '', config: {} })
const configText = ref('')

const fetchList = async () => {
  loading.value = true
  try {
    const res = await templateApi.listPromotions()
    templates.value = Array.isArray(res) ? res : []
  } finally {
    loading.value = false
  }
}

const openCreate = () => {
  form.value = { name: '', config: {} }
  configText.value = '{}'
  showDialog.value = true
}

const handleSave = async () => {
  try {
    const config = JSON.parse(configText.value)
    const data: CreateTemplateParams = { name: form.value.name, config }
    await templateApi.createPromotion(data)
    showDialog.value = false
    await fetchList()
  } catch (e: unknown) {
    if (e instanceof SyntaxError) {
      alert('JSON 配置格式错误，请检查')
    }
  }
}

const handleDelete = async (tpl: PromotionTemplate) => {
  if (!confirm(`确定删除模板「${tpl.name}」吗？`)) return
  await templateApi.deletePromotion(tpl.id)
  await fetchList()
}

onMounted(fetchList)
</script>

<template>
  <div class="space-y-6">
    <div>
      <Breadcrumb :items="[{ name: '本地投放' }, { name: '单元模板' }]" />
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-3xl font-bold text-gray-900">单元模板管理</h1>
          <p class="mt-2 text-gray-600">管理批量创建广告单元的配置模板</p>
        </div>
        <button @click="openCreate" class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 flex items-center gap-2">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
          </svg>
          新建模板
        </button>
      </div>
    </div>

    <div class="bg-white rounded-lg border border-gray-200">
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead class="bg-gray-50 border-b border-gray-200">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-semibold text-gray-600 uppercase">模板名称</th>
              <th class="px-6 py-3 text-left text-xs font-semibold text-gray-600 uppercase">使用次数</th>
              <th class="px-6 py-3 text-left text-xs font-semibold text-gray-600 uppercase">创建时间</th>
              <th class="px-6 py-3 text-left text-xs font-semibold text-gray-600 uppercase">操作</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr v-for="tpl in templates" :key="tpl.id" class="hover:bg-gray-50">
              <td class="px-6 py-4 text-sm font-medium text-gray-900">{{ tpl.name }}</td>
              <td class="px-6 py-4 text-sm text-gray-700">{{ tpl.use_count }}</td>
              <td class="px-6 py-4 text-sm text-gray-500">{{ tpl.created_at }}</td>
              <td class="px-6 py-4 text-sm">
                <button @click="handleDelete(tpl)" class="text-red-600 hover:text-red-800">删除</button>
              </td>
            </tr>
            <tr v-if="!loading && templates.length === 0">
              <td colspan="4" class="px-6 py-8 text-center text-gray-500">暂无模板数据</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- 新建弹窗 -->
    <div v-if="showDialog" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
      <div class="bg-white rounded-lg shadow-xl w-full max-w-lg mx-4">
        <div class="px-6 py-4 border-b border-gray-200">
          <h3 class="text-lg font-semibold text-gray-900">新建模板</h3>
        </div>
        <div class="px-6 py-4 space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">模板名称</label>
            <input v-model="form.name" type="text" placeholder="请输入模板名称" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">JSON 配置</label>
            <textarea v-model="configText" rows="10" placeholder="{}" class="w-full px-3 py-2 border border-gray-300 rounded-lg font-mono text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"></textarea>
          </div>
        </div>
        <div class="px-6 py-4 border-t border-gray-200 flex justify-end gap-3">
          <button @click="showDialog = false" class="px-4 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50">取消</button>
          <button @click="handleSave" class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700">保存</button>
        </div>
      </div>
    </div>
  </div>
</template>
