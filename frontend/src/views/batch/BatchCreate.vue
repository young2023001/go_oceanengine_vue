<template>
  <div class="max-w-4xl mx-auto p-6">
    <h1 class="text-2xl font-bold mb-8">批量操作向导</h1>

    <!-- 步骤指示器 -->
    <div class="flex items-center mb-10">
      <template v-for="(step, index) in steps" :key="index">
        <div class="flex items-center">
          <div
            class="w-8 h-8 rounded-full flex items-center justify-center text-sm font-medium"
            :class="getStepClass(index)"
          >
            {{ index + 1 }}
          </div>
          <span
            class="ml-2 text-sm"
            :class="currentStep >= index ? 'text-gray-900 font-medium' : 'text-gray-400'"
          >
            {{ step }}
          </span>
        </div>
        <div
          v-if="index < steps.length - 1"
          class="flex-1 h-px mx-4"
          :class="currentStep > index ? 'bg-blue-500' : 'bg-gray-200'"
        />
      </template>
    </div>

    <!-- 步骤 1: 选择操作类型 -->
    <div v-if="currentStep === 0" class="space-y-4">
      <h2 class="text-lg font-semibold mb-4">选择操作类型</h2>
      <div class="grid grid-cols-2 gap-4">
        <label
          v-for="option in taskTypeOptions"
          :key="option.value"
          class="flex items-center p-4 border rounded-lg cursor-pointer transition-colors"
          :class="
            formState.taskType === option.value
              ? 'border-blue-500 bg-blue-50'
              : 'border-gray-200 hover:border-gray-300'
          "
        >
          <input
            v-model="formState.taskType"
            type="radio"
            :value="option.value"
            class="mr-3"
          />
          <div>
            <div class="font-medium">{{ option.label }}</div>
            <div class="text-sm text-gray-500">{{ option.description }}</div>
          </div>
        </label>
      </div>
    </div>

    <!-- 步骤 2: 选择目标账户 -->
    <div v-if="currentStep === 1" class="space-y-4">
      <h2 class="text-lg font-semibold mb-4">选择目标账户</h2>
      <div v-if="accountsLoading" class="text-gray-500">加载中...</div>
      <div v-else-if="accounts.length === 0" class="text-gray-500">暂无可用账户</div>
      <div v-else class="space-y-2 max-h-96 overflow-y-auto">
        <label
          v-for="account in accounts"
          :key="account.id"
          class="flex items-center p-3 border rounded-lg cursor-pointer transition-colors"
          :class="
            formState.accountIds.includes(account.id)
              ? 'border-blue-500 bg-blue-50'
              : 'border-gray-200 hover:border-gray-300'
          "
        >
          <input
            v-model="formState.accountIds"
            type="checkbox"
            :value="account.id"
            class="mr-3"
          />
          <div>
            <div class="font-medium">{{ account.name }}</div>
            <div class="text-sm text-gray-400">ID: {{ account.local_account_id }}</div>
          </div>
        </label>
      </div>
    </div>

    <!-- 步骤 3: 配置参数 -->
    <div v-if="currentStep === 2" class="space-y-4">
      <h2 class="text-lg font-semibold mb-4">配置参数</h2>

      <!-- 批量创建项目 -->
      <div v-if="formState.taskType === 'create_project'" class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">项目名称</label>
          <input
            v-model="formState.config.projectName"
            type="text"
            placeholder="请输入项目名称"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">日预算 (元)</label>
          <input
            v-model.number="formState.config.dailyBudget"
            type="number"
            placeholder="请输入日预算"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
      </div>

      <!-- 调预算 -->
      <div v-if="formState.taskType === 'adjust_budget'" class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">调整方式</label>
          <select
            v-model="formState.config.adjustMode"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="set">设置为固定值</option>
            <option value="increase">增加</option>
            <option value="decrease">减少</option>
          </select>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">金额 (元)</label>
          <input
            v-model.number="formState.config.budgetAmount"
            type="number"
            placeholder="请输入金额"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
      </div>

      <!-- 调出价 -->
      <div v-if="formState.taskType === 'adjust_bid'" class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">调整方式</label>
          <select
            v-model="formState.config.adjustMode"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="set">设置为固定值</option>
            <option value="increase_percent">按比例增加</option>
            <option value="decrease_percent">按比例减少</option>
          </select>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">
            {{ formState.config.adjustMode === 'set' ? '出价 (元)' : '比例 (%)' }}
          </label>
          <input
            v-model.number="formState.config.bidValue"
            type="number"
            placeholder="请输入数值"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
      </div>

      <!-- 启停 -->
      <div v-if="formState.taskType === 'toggle_status'" class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">操作</label>
          <select
            v-model="formState.config.targetStatus"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="enable">启用</option>
            <option value="disable">暂停</option>
          </select>
        </div>
      </div>
    </div>

    <!-- 步骤 4: 确认提交 -->
    <div v-if="currentStep === 3" class="space-y-4">
      <h2 class="text-lg font-semibold mb-4">确认提交</h2>
      <div class="bg-gray-50 rounded-lg p-6 space-y-3">
        <div class="flex">
          <span class="text-gray-500 w-24">操作类型:</span>
          <span class="font-medium">{{ getTaskTypeLabel(formState.taskType) }}</span>
        </div>
        <div class="flex">
          <span class="text-gray-500 w-24">目标账户:</span>
          <span class="font-medium">已选择 {{ formState.accountIds.length }} 个账户</span>
        </div>
        <div class="flex">
          <span class="text-gray-500 w-24">配置参数:</span>
          <span class="font-medium">{{ getConfigSummary() }}</span>
        </div>
      </div>
    </div>

    <!-- 底部按钮 -->
    <div class="flex justify-between mt-10 pt-6 border-t">
      <button
        v-if="currentStep > 0"
        class="px-6 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50 transition-colors"
        @click="prevStep"
      >
        上一步
      </button>
      <div v-else />
      <button
        v-if="currentStep < steps.length - 1"
        :disabled="!canProceed"
        class="px-6 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
        @click="nextStep"
      >
        下一步
      </button>
      <button
        v-else
        :disabled="submitting"
        class="px-6 py-2 bg-green-500 text-white rounded-lg hover:bg-green-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
        @click="handleSubmit"
      >
        {{ submitting ? '提交中...' : '确认提交' }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { batchApi, type CreateBatchTaskParams } from '@/api/batch'
import { accountApi, type Account } from '@/api/account'

const router = useRouter()

const steps = ['选择操作类型', '选择目标账户', '配置参数', '确认提交']
const currentStep = ref(0)
const submitting = ref(false)
const accountsLoading = ref(false)
const accounts = ref<Account[]>([])

interface TaskTypeOption {
  value: string
  label: string
  description: string
}

const taskTypeOptions: TaskTypeOption[] = [
  { value: 'create_project', label: '批量创建项目', description: '为多个账户批量创建广告项目' },
  { value: 'adjust_budget', label: '批量调预算', description: '批量调整账户或计划的预算' },
  { value: 'adjust_bid', label: '批量调出价', description: '批量调整广告出价' },
  { value: 'toggle_status', label: '批量启停', description: '批量启用或暂停广告计划' },
]

interface FormConfig {
  projectName: string
  dailyBudget: number | null
  adjustMode: string
  budgetAmount: number | null
  bidValue: number | null
  targetStatus: string
}

interface FormState {
  taskType: string
  accountIds: number[]
  config: FormConfig
}

const formState: FormState = reactive({
  taskType: '',
  accountIds: [],
  config: {
    projectName: '',
    dailyBudget: null,
    adjustMode: 'set',
    budgetAmount: null,
    bidValue: null,
    targetStatus: 'enable',
  },
})

const canProceed = computed((): boolean => {
  switch (currentStep.value) {
    case 0:
      return formState.taskType !== ''
    case 1:
      return formState.accountIds.length > 0
    case 2:
      return isConfigValid()
    default:
      return true
  }
})

function isConfigValid(): boolean {
  switch (formState.taskType) {
    case 'create_project':
      return formState.config.projectName.trim() !== '' && formState.config.dailyBudget !== null && formState.config.dailyBudget > 0
    case 'adjust_budget':
      return formState.config.budgetAmount !== null && formState.config.budgetAmount > 0
    case 'adjust_bid':
      return formState.config.bidValue !== null && formState.config.bidValue > 0
    case 'toggle_status':
      return formState.config.targetStatus !== ''
    default:
      return false
  }
}

function getStepClass(index: number): string {
  if (currentStep.value > index) {
    return 'bg-blue-500 text-white'
  }
  if (currentStep.value === index) {
    return 'bg-blue-500 text-white'
  }
  return 'bg-gray-200 text-gray-500'
}

function getTaskTypeLabel(value: string): string {
  const option = taskTypeOptions.find((o) => o.value === value)
  return option?.label ?? ''
}

function getConfigSummary(): string {
  switch (formState.taskType) {
    case 'create_project':
      return `项目名称: ${formState.config.projectName}, 日预算: ${formState.config.dailyBudget}元`
    case 'adjust_budget':
      return `${getAdjustModeLabel(formState.config.adjustMode)} ${formState.config.budgetAmount}元`
    case 'adjust_bid':
      return `${getAdjustModeLabel(formState.config.adjustMode)} ${formState.config.bidValue}${formState.config.adjustMode === 'set' ? '元' : '%'}`
    case 'toggle_status':
      return formState.config.targetStatus === 'enable' ? '启用' : '暂停'
    default:
      return ''
  }
}

function getAdjustModeLabel(mode: string): string {
  const labels: Record<string, string> = {
    set: '设置为',
    increase: '增加',
    decrease: '减少',
    increase_percent: '按比例增加',
    decrease_percent: '按比例减少',
  }
  return labels[mode] ?? mode
}

function nextStep(): void {
  if (currentStep.value < steps.length - 1) {
    currentStep.value++
  }
}

function prevStep(): void {
  if (currentStep.value > 0) {
    currentStep.value--
  }
}

function buildTaskConfig(): object {
  switch (formState.taskType) {
    case 'create_project':
      return { projectName: formState.config.projectName, dailyBudget: formState.config.dailyBudget }
    case 'adjust_budget':
      return { adjustMode: formState.config.adjustMode, amount: formState.config.budgetAmount }
    case 'adjust_bid':
      return { adjustMode: formState.config.adjustMode, value: formState.config.bidValue }
    case 'toggle_status':
      return { targetStatus: formState.config.targetStatus }
    default:
      return {}
  }
}

async function handleSubmit(): Promise<void> {
  submitting.value = true
  try {
    const params: CreateBatchTaskParams = {
      task_type: formState.taskType,
      title: `${getTaskTypeLabel(formState.taskType)} - ${formState.accountIds.length}个账户`,
      config: buildTaskConfig(),
      account_ids: formState.accountIds,
    }
    await batchApi.createTask(params)
    router.push('/batch/tasks')
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : '操作失败，请重试'
    alert(message)
  } finally {
    submitting.value = false
  }
}

async function fetchAccounts(): Promise<void> {
  accountsLoading.value = true
  try {
    const res = await accountApi.getList({ page: 1, page_size: 200 })
    accounts.value = res.list ?? []
  } finally {
    accountsLoading.value = false
  }
}

onMounted(() => {
  fetchAccounts()
})
</script>
