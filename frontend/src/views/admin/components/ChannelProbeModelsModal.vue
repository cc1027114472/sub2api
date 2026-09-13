<template>
  <BaseDialog
    :show="visible"
    :title="t('admin.channels.probeModal.title', '从上游探测并选择模型')"
    width="wide"
    @close="handleClose"
  >
    <div class="space-y-4">
      <!-- 1. 连接配置区域 -->
      <div class="rounded-xl border border-gray-200 bg-gray-50/70 p-4 dark:border-dark-700 dark:bg-dark-800/60">
        <div class="flex items-center justify-between mb-3">
          <div class="flex items-center gap-2">
            <Icon name="bolt" size="sm" class="text-primary-600 dark:text-primary-400" />
            <span class="text-sm font-semibold text-gray-800 dark:text-gray-200">
              {{ t('admin.channels.probeModal.connectionTitle', '上游节点连接配置') }}
            </span>
          </div>
          <span
            v-if="probed"
            class="inline-flex items-center gap-1 rounded-md bg-green-50 px-2 py-0.5 text-xs font-medium text-green-700 dark:bg-green-900/20 dark:text-green-400"
          >
            <Icon name="check" size="xs" />
            {{ t('admin.channels.probeModal.probedSuccess', '探测成功') }}
          </span>
        </div>

        <div class="grid grid-cols-1 gap-3 sm:grid-cols-2">
          <!-- Base URL -->
          <div class="sm:col-span-2">
            <label class="input-label text-xs">
              {{ t('admin.channels.quickSync.baseUrl', '基础接口地址 (Base URL)') }}
              <span class="text-red-500">*</span>
            </label>
            <input
              v-model="form.base_url"
              type="text"
              class="input font-mono text-xs"
              placeholder="http://154.36.173.146 或 https://api.openai.com"
              :disabled="loading"
            />
          </div>

          <!-- API Key -->
          <div>
            <label class="input-label text-xs">
              {{ t('admin.channels.quickSync.apiKey', 'API 密钥 (API Key)') }}
              <span class="text-red-500">*</span>
            </label>
            <div class="relative">
              <input
                v-model="form.api_key"
                :type="showApiKey ? 'text' : 'password'"
                class="input pr-9 font-mono text-xs"
                placeholder="sk-... 或上游 Token"
                :disabled="loading"
              />
              <button
                type="button"
                class="absolute right-2.5 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
                @click="showApiKey = !showApiKey"
              >
                <Icon :name="showApiKey ? 'eyeOff' : 'eye'" size="xs" />
              </button>
            </div>
          </div>

          <!-- 平台类型 -->
          <div>
            <label class="input-label text-xs">{{ t('admin.channels.quickSync.platform', '平台类型') }}</label>
            <select v-model="form.platform" class="input text-xs" :disabled="loading">
              <option value="antigravity">Antigravity (多协议桥接)</option>
              <option value="gemini">Gemini</option>
              <option value="openai">OpenAI</option>
              <option value="anthropic">Anthropic</option>
            </select>
          </div>
        </div>

        <div class="mt-3 flex items-center justify-between">
          <p class="text-[11px] text-gray-400 dark:text-dark-400">
            {{ t('admin.channels.probeModal.autoSaveHint', '连接信息将自动安全保存在本地浏览器中') }}
          </p>
          <button
            type="button"
            class="btn btn-primary py-1.5 px-3 text-xs inline-flex items-center gap-1.5"
            :disabled="loading || !form.base_url.trim() || !form.api_key.trim()"
            @click="probeUpstream"
          >
            <Icon v-if="loading" name="refresh" size="xs" class="animate-spin" />
            <Icon v-else name="bolt" size="xs" />
            <span>{{ loading ? t('admin.channels.quickSync.probing', '正在探测...') : t('admin.channels.probeModal.startProbe', '连接并探测模型') }}</span>
          </button>
        </div>

        <div v-if="probeError" class="mt-3 rounded-lg bg-red-50 p-2.5 text-xs text-red-600 dark:bg-red-900/20 dark:text-red-400">
          {{ probeError }}
        </div>
      </div>

      <!-- 2. 模型列表与多选区域 -->
      <div v-if="probed" class="space-y-3">
        <!-- 工具栏：搜索 + 统计 + 批量勾选 -->
        <div class="flex flex-wrap items-center justify-between gap-2 border-b border-gray-100 pb-2 dark:border-dark-700">
          <div class="flex items-center gap-2">
            <span class="text-xs font-semibold text-gray-700 dark:text-gray-300">
              {{ t('admin.channels.probeModal.modelsListTitle', '探测到的模型列表') }}
            </span>
            <span class="text-xs text-primary-600 dark:text-primary-400 font-medium">
              {{ t('admin.channels.probeModal.selectedCount', { count: selectedCount }, `已选 ${selectedCount} 个`) }}
            </span>
          </div>

          <div class="flex items-center gap-3">
            <div class="relative w-44 sm:w-56">
              <Icon name="search" size="xs" class="absolute left-2.5 top-1/2 -translate-y-1/2 text-gray-400" />
              <input
                v-model="searchQuery"
                type="text"
                class="input py-1 pl-8 pr-3 text-xs"
                :placeholder="t('admin.channels.probeModal.searchPlaceholder', '搜索模型名称...')"
              />
            </div>
            <div class="flex items-center gap-1.5 text-xs text-gray-500">
              <button
                type="button"
                class="hover:text-primary-600 dark:hover:text-primary-400 font-medium"
                @click="selectAllAvailable(true)"
              >
                {{ t('admin.channels.probeModal.selectAllNew', '全选未添加') }}
              </button>
              <span>/</span>
              <button
                type="button"
                class="hover:text-primary-600 dark:hover:text-primary-400"
                @click="selectAllAvailable(false)"
              >
                {{ t('admin.channels.probeModal.clearAll', '全不选') }}
              </button>
            </div>
          </div>
        </div>

        <!-- 列表容器 -->
        <div class="max-h-72 overflow-y-auto rounded-lg border border-gray-200 divide-y divide-gray-100 dark:border-dark-700 dark:divide-dark-700/60 bg-white dark:bg-dark-800">
          <div
            v-for="m in filteredModels"
            :key="m.id"
            class="flex items-center justify-between p-2.5 transition hover:bg-gray-50 dark:hover:bg-dark-700/40"
            :class="{ 'opacity-60 bg-gray-50/50 dark:bg-dark-700/20': isAlreadyInChannel(m.id) }"
          >
            <label class="flex items-center gap-2.5 cursor-pointer min-w-0 flex-1">
              <input
                type="checkbox"
                class="rounded border-gray-300 text-primary-600 focus:ring-primary-500 dark:border-dark-600 dark:bg-dark-700"
                :disabled="isAlreadyInChannel(m.id)"
                :checked="isAlreadyInChannel(m.id) || selectedMap[m.id]"
                @change="toggleSelect(m.id)"
              />
              <span class="font-mono text-xs font-medium text-gray-800 dark:text-gray-200 truncate select-all">
                {{ m.id }}
              </span>
            </label>

            <!-- 状态标签 -->
            <div class="flex items-center gap-2 shrink-0 ml-3">
              <span
                v-if="isAlreadyInChannel(m.id)"
                class="inline-flex items-center rounded-md bg-gray-100 px-2 py-0.5 text-[11px] font-medium text-gray-500 dark:bg-dark-700 dark:text-dark-300"
              >
                {{ t('admin.channels.probeModal.alreadyInChannel', '渠道中已有') }}
              </span>
              <span
                v-else
                class="inline-flex items-center rounded-md bg-blue-50 px-2 py-0.5 text-[11px] font-medium text-blue-700 dark:bg-blue-900/20 dark:text-blue-300"
              >
                {{ t('admin.channels.probeModal.canAdd', '未添加') }}
              </span>
            </div>
          </div>

          <div v-if="filteredModels.length === 0" class="p-6 text-center text-xs text-gray-400">
            {{ t('admin.channels.probeModal.noMatchingModels', '没有找到匹配的模型') }}
          </div>
        </div>

        <!-- 目标定价归属选择 -->
        <div class="flex flex-wrap items-center justify-between gap-2 pt-1 text-xs text-gray-600 dark:text-gray-300">
          <div class="flex items-center gap-2">
            <span class="font-medium">{{ t('admin.channels.probeModal.targetEntryLabel', '导入目标:') }}</span>
            <select v-model="targetEntryIndex" class="input py-1 text-xs w-56">
              <option value="new">
                ⭐ {{ t('admin.channels.probeModal.targetNewEntry', '创建为新的定价配置卡片') }}
              </option>
              <option
                v-for="idx in pricingEntriesCount"
                :key="idx"
                :value="idx - 1"
              >
                {{ t('admin.channels.probeModal.targetExistingEntry', { index: idx }, `追加到现有第 ${idx} 组配置`) }}
              </option>
            </select>
          </div>
        </div>
      </div>
    </div>

    <!-- Footer -->
    <template #footer>
      <div class="flex items-center justify-between w-full">
        <div class="text-xs text-gray-500 dark:text-dark-400">
          <template v-if="probed">
            {{ t('admin.channels.probeModal.summaryTotal', { total: probedModels.length, available: availableModels.length }, `共 ${probedModels.length} 个模型，${availableModels.length} 个未添加`) }}
          </template>
        </div>
        <div class="flex items-center gap-2">
          <button type="button" class="btn btn-secondary py-1.5 px-3 text-xs" @click="handleClose">
            {{ t('common.cancel', '取消') }}
          </button>
          <button
            type="button"
            class="btn btn-primary py-1.5 px-3 text-xs inline-flex items-center gap-1.5"
            :disabled="!probed || selectedCount === 0"
            @click="confirmImport"
          >
            <Icon name="check" size="xs" />
            <span>{{ t('admin.channels.probeModal.confirmImport', { count: selectedCount }, `导入已选 (${selectedCount}) 个模型`) }}</span>
          </button>
        </div>
      </div>
    </template>
  </BaseDialog>
</template>

<script setup lang="ts">
import { ref, computed, watch, reactive } from 'vue'
import { useI18n } from 'vue-i18n'
import BaseDialog from '@/components/common/BaseDialog.vue'
import Icon from '@/components/icons/Icon.vue'
import adminAPI from '@/api/admin'
import type { QuickSyncProbeModel } from '@/api/admin/channels'

const { t } = useI18n()

const props = defineProps<{
  visible: boolean
  platform: string
  existingModels: string[]
  pricingEntriesCount: number
}>()

const emit = defineEmits<{
  'update:visible': [value: boolean]
  'import-models': [payload: {
    models: string[]
    targetEntryIndex: number | 'new'
  }]
}>()

const STORAGE_KEY_BASE_URL = 'sub2api_upstream_base_url'
const STORAGE_KEY_API_KEY = 'sub2api_upstream_api_key'

const form = reactive({
  base_url: '',
  api_key: '',
  platform: props.platform || 'antigravity'
})

const showApiKey = ref(false)
const loading = ref(false)
const probed = ref(false)
const probeError = ref('')
const probedModels = ref<QuickSyncProbeModel[]>([])
const selectedMap = reactive<Record<string, boolean>>({})
const searchQuery = ref('')
const targetEntryIndex = ref<number | 'new'>('new')

// 监听弹窗打开，初始化配置
watch(
  () => props.visible,
  (val) => {
    if (val) {
      probeError.value = ''
      // 从本地存储或默认值载入
      form.base_url = localStorage.getItem(STORAGE_KEY_BASE_URL) || form.base_url || 'http://154.36.173.146'
      form.api_key = localStorage.getItem(STORAGE_KEY_API_KEY) || form.api_key || ''
      form.platform = props.platform || 'antigravity'
      targetEntryIndex.value = 'new'
    }
  },
  { immediate: true }
)

const existingModelSet = computed(() => new Set(props.existingModels.map(m => m.trim().toLowerCase())))

function isAlreadyInChannel(modelId: string): boolean {
  return existingModelSet.value.has(modelId.trim().toLowerCase())
}

const availableModels = computed(() => {
  return probedModels.value.filter((m: QuickSyncProbeModel) => !isAlreadyInChannel(m.id))
})

const filteredModels = computed(() => {
  const q = searchQuery.value.trim().toLowerCase()
  if (!q) return probedModels.value
  const keywords = q.split(/\s+/).filter(Boolean)
  return probedModels.value.filter((m: QuickSyncProbeModel) => {
    const name = m.id.toLowerCase()
    return keywords.every(kw => name.includes(kw))
  })
})

const selectedCount = computed(() => {
  return availableModels.value.filter((m: QuickSyncProbeModel) => selectedMap[m.id]).length
})

function toggleSelect(modelId: string) {
  if (isAlreadyInChannel(modelId)) return
  selectedMap[modelId] = !selectedMap[modelId]
}

function selectAllAvailable(selected: boolean) {
  for (const m of availableModels.value) {
    selectedMap[m.id] = selected
  }
}

async function probeUpstream() {
  if (!form.base_url.trim() || !form.api_key.trim()) return
  loading.value = true
  probeError.value = ''

  try {
    const result = await adminAPI.channels.quickSyncProbe({
      base_url: form.base_url.trim(),
      api_key: form.api_key.trim(),
      platform: form.platform
    })

    probedModels.value = result.models || []
    probed.value = true

    // 默认全选所有未添加的新模型
    for (const key of Object.keys(selectedMap)) {
      delete selectedMap[key]
    }
    for (const m of probedModels.value) {
      if (!isAlreadyInChannel(m.id)) {
        selectedMap[m.id] = true
      }
    }

    // 保存到本地存储供下次记忆
    localStorage.setItem(STORAGE_KEY_BASE_URL, form.base_url.trim())
    localStorage.setItem(STORAGE_KEY_API_KEY, form.api_key.trim())
  } catch (err: any) {
    probeError.value = err.message || t('admin.channels.quickSync.probeFailed', '探测失败，请检查 Base URL 与 API Key 是否正确')
  } finally {
    loading.value = false
  }
}

function confirmImport() {
  const toImport = availableModels.value
    .filter((m: QuickSyncProbeModel) => selectedMap[m.id])
    .map((m: QuickSyncProbeModel) => m.id)

  if (toImport.length === 0) return

  emit('import-models', {
    models: toImport,
    targetEntryIndex: targetEntryIndex.value
  })

  handleClose()
}

function handleClose() {
  emit('update:visible', false)
}
</script>
