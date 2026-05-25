<script setup lang="ts">
import { ref } from 'vue'
import Breadcrumb from '@/components/common/Breadcrumb.vue'

const seedWord = ref('')
const expandedWords = ref<{ word: string; heat: number; related: string[] }[]>([])

const mockExpand = () => {
  expandedWords.value = [
    { word: '智能手机推荐', heat: 95, related: ['性价比手机', '旗舰手机', '千元机'] },
    { word: '手机测评', heat: 92, related: ['开箱测评', '对比测试', '深度体验'] },
    { word: '新款手机', heat: 88, related: ['新品发布', '首发体验', '限时优惠'] },
    { word: '手机配件', heat: 75, related: ['手机壳', '充电器', '耳机'] },
    { word: '手机维修', heat: 65, related: ['换屏幕', '换电池', '维修店'] }
  ]
}

const handleAddToLibrary = (_item: typeof expandedWords.value[0]) => {
  // TODO: implement
}

const handleUse = (_item: typeof expandedWords.value[0]) => {
  // TODO: implement
}
</script>

<template>
  <div class="space-y-6">
    <div>
      <Breadcrumb :items="[{ name: '创意工具' }, { name: '词汇扩展' }]" />
      <h1 class="text-3xl font-bold text-gray-900">创意词汇扩展</h1>
      <p class="mt-2 text-gray-600">基于种子词扩展相关创意词汇</p>
    </div>

    <div class="bg-white rounded-lg border border-gray-200 p-4">
      <div class="flex gap-4">
        <input v-model="seedWord" type="text" placeholder="输入种子词，如：手机"
               class="flex-1 px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500">
        <button class="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700" @click="mockExpand">
          开始扩展
        </button>
      </div>
    </div>

    <div v-if="expandedWords.length > 0" class="space-y-4">
      <div class="flex items-center justify-between">
        <h3 class="font-semibold text-gray-900">扩展结果</h3>
        <span class="text-sm text-gray-500">共 {{ expandedWords.length }} 个词组</span>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <div v-for="item in expandedWords" :key="item.word"
             class="bg-white rounded-lg border border-gray-200 p-4 hover:shadow-md transition-shadow">
          <div class="flex items-start justify-between">
            <h4 class="font-semibold text-gray-900">{{ item.word }}</h4>
            <div class="flex items-center gap-1">
              <span class="text-orange-500">🔥</span>
              <span class="text-sm font-medium text-orange-600">{{ item.heat }}</span>
            </div>
          </div>
          <div class="mt-3">
            <p class="text-xs text-gray-500 mb-2">相关词汇：</p>
            <div class="flex flex-wrap gap-1">
              <span v-for="rel in item.related" :key="rel"
                    class="px-2 py-0.5 bg-blue-50 text-blue-700 rounded text-xs">
                {{ rel }}
              </span>
            </div>
          </div>
          <div class="mt-3 flex gap-2">
<button class="flex-1 py-1.5 text-xs text-blue-600 border border-blue-300 rounded hover:bg-blue-50" @click="handleAddToLibrary(item)">
              添加到词库
            </button>
            <button class="flex-1 py-1.5 text-xs text-white bg-blue-600 rounded hover:bg-blue-700" @click="handleUse(item)">
              使用
            </button>
          </div>
        </div>
      </div>
    </div>

    <div v-else class="bg-gray-50 rounded-lg p-12 text-center">
      <div class="text-5xl mb-4">💡</div>
      <p class="text-gray-600">输入种子词开始扩展创意词汇</p>
      <p class="text-sm text-gray-400 mt-2">支持商品名、行业词、品牌词等</p>
    </div>
  </div>
</template>
