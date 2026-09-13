<template>
  <BaseDialog
    :show="show"
    :title="t('keys.postmanModal.title')"
    width="wide"
    @close="emit('close')"
  >
    <div class="space-y-4">
      <!-- Subtitle & Context Card -->
      <div class="rounded-xl border border-gray-200 bg-gray-50/70 p-4 dark:border-dark-700 dark:bg-dark-800/60">
        <div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <p class="text-sm text-gray-600 dark:text-gray-300">
              {{ t('keys.postmanModal.subtitle') }}
            </p>
            <div class="mt-2 flex flex-wrap items-center gap-2 text-xs">
              <span v-if="keyName" class="inline-flex items-center gap-1 rounded-md bg-primary-50 px-2 py-0.5 font-medium text-primary-700 dark:bg-primary-950/40 dark:text-primary-300">
                <Icon name="key" size="xs" />
                {{ keyName }}
              </span>
              <span v-if="groupName" class="inline-flex items-center gap-1 rounded-md bg-gray-100 px-2 py-0.5 font-medium text-gray-700 dark:bg-dark-700 dark:text-gray-300">
                <Icon name="cube" size="xs" />
                {{ groupName }}
              </span>
              <span class="inline-flex items-center gap-1 font-mono text-gray-500 dark:text-gray-400">
                {{ maskedKey }}
              </span>
            </div>
          </div>
          <div class="flex items-center gap-2">
            <button
              type="button"
              @click="copyApiKey"
              class="inline-flex items-center gap-1.5 rounded-lg border border-gray-300 bg-white px-2.5 py-1.5 text-xs font-medium text-gray-700 shadow-sm transition hover:bg-gray-50 dark:border-dark-600 dark:bg-dark-700 dark:text-gray-200 dark:hover:bg-dark-600"
            >
              <Icon :name="copiedKey ? 'check' : 'copy'" size="xs" />
              <span>{{ copiedKey ? t('keys.postmanModal.copied') : t('keys.useKeyModal.copy') }} Key</span>
            </button>
          </div>
        </div>

        <!-- Base URL and Model Settings Row -->
        <div class="mt-4 grid grid-cols-1 gap-3 sm:grid-cols-2">
          <!-- Base URL Input -->
          <div>
            <label class="block text-xs font-medium text-gray-700 dark:text-gray-300">
              Base URL
            </label>
            <div class="mt-1 flex rounded-lg shadow-sm">
              <input
                v-model="customBaseUrl"
                type="text"
                class="block w-full rounded-lg border border-gray-300 bg-white px-3 py-1.5 text-xs text-gray-900 focus:border-primary-500 focus:outline-none focus:ring-1 focus:ring-primary-500 dark:border-dark-600 dark:bg-dark-900 dark:text-gray-100"
                :placeholder="defaultBaseUrl"
              />
            </div>
          </div>

          <!-- Model Selector / Input -->
          <div>
            <label class="block text-xs font-medium text-gray-700 dark:text-gray-300">
              {{ t('keys.postmanModal.selectModel') }}
            </label>
            <div class="mt-1 flex items-center gap-1.5">
              <select
                v-if="availableModels.length > 0"
                v-model="selectedModel"
                class="block w-full rounded-lg border border-gray-300 bg-white px-3 py-1.5 text-xs text-gray-900 focus:border-primary-500 focus:outline-none focus:ring-1 focus:ring-primary-500 dark:border-dark-600 dark:bg-dark-900 dark:text-gray-100"
              >
                <option v-for="m in availableModels" :key="m" :value="m">
                  {{ m }}
                </option>
                <option value="__custom__">
                  + {{ t('keys.postmanModal.customModelOption') }}
                </option>
              </select>
              <input
                v-if="availableModels.length === 0 || selectedModel === '__custom__'"
                v-model="customModelInput"
                type="text"
                class="block w-full rounded-lg border border-gray-300 bg-white px-3 py-1.5 text-xs text-gray-900 focus:border-primary-500 focus:outline-none focus:ring-1 focus:ring-primary-500 dark:border-dark-600 dark:bg-dark-900 dark:text-gray-100"
                :placeholder="t('keys.postmanModal.customModelPlaceholder')"
              />
            </div>
          </div>
        </div>
      </div>

      <!-- Protocol Selection Tabs -->
      <div class="border-b border-gray-200 dark:border-dark-700">
        <nav class="-mb-px flex space-x-2 sm:space-x-4 overflow-x-auto" aria-label="Protocol Tabs">
          <button
            v-for="proto in protocols"
            :key="proto.id"
            type="button"
            @click="activeProtocol = proto.id"
            :class="[
              'group inline-flex items-center gap-2 border-b-2 px-3 py-2 text-xs font-medium whitespace-nowrap transition-colors',
              activeProtocol === proto.id
                ? 'border-primary-500 text-primary-600 dark:text-primary-400'
                : 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700 dark:text-gray-400 dark:hover:border-dark-600 dark:hover:text-gray-300'
            ]"
          >
            <span
              :class="[
                'rounded px-1.5 py-0.5 text-[10px] font-bold uppercase tracking-wider',
                activeProtocol === proto.id
                  ? 'bg-primary-100 text-primary-700 dark:bg-primary-950/60 dark:text-primary-300'
                  : 'bg-gray-100 text-gray-600 dark:bg-dark-700 dark:text-gray-400'
              ]"
            >
              {{ proto.badge }}
            </span>
            <span>{{ proto.label }}</span>
          </button>
        </nav>
      </div>

      <!-- Active Protocol Info Bar -->
      <div class="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between rounded-lg bg-gray-100/60 px-3 py-2 text-xs dark:bg-dark-800/80">
        <div class="flex items-center gap-2 font-mono text-gray-700 dark:text-gray-300">
          <span class="rounded bg-emerald-100 px-1.5 py-0.5 font-bold text-emerald-800 dark:bg-emerald-950/60 dark:text-emerald-300">POST</span>
          <span class="break-all">{{ currentProtocolConfig.endpoint }}</span>
        </div>
        <div class="flex items-center gap-2">
          <span class="text-gray-500 dark:text-gray-400">{{ currentProtocolConfig.note }}</span>
        </div>
      </div>

      <!-- cURL Code Block -->
      <div class="relative rounded-xl border border-gray-800 bg-gray-900 text-gray-100 shadow-inner">
        <div class="flex items-center justify-between border-b border-gray-800 bg-gray-950/70 px-4 py-2 text-xs">
          <div class="flex items-center gap-2 text-gray-400">
            <span class="h-2.5 w-2.5 rounded-full bg-red-500/80"></span>
            <span class="h-2.5 w-2.5 rounded-full bg-yellow-500/80"></span>
            <span class="h-2.5 w-2.5 rounded-full bg-green-500/80"></span>
            <span class="ml-2 font-mono text-[11px] text-gray-400">cURL / Postman Raw</span>
          </div>
          <div class="flex items-center gap-2">
            <button
              type="button"
              @click="copyBody"
              class="inline-flex items-center gap-1 rounded-md bg-gray-800 px-2 py-1 text-[11px] text-gray-300 transition hover:bg-gray-700 hover:text-white"
            >
              <Icon :name="copiedBodyState ? 'check' : 'clipboard'" size="xs" />
              <span>{{ copiedBodyState ? t('keys.postmanModal.copied') : t('keys.postmanModal.copyBody') }}</span>
            </button>
            <button
              type="button"
              @click="copyCurl"
              class="inline-flex items-center gap-1 rounded-md bg-primary-600 px-2.5 py-1 text-[11px] font-medium text-white shadow-sm transition hover:bg-primary-500"
            >
              <Icon :name="copiedCurlState ? 'check' : 'copy'" size="xs" />
              <span>{{ copiedCurlState ? t('keys.postmanModal.copied') : t('keys.postmanModal.copyCurl') }}</span>
            </button>
          </div>
        </div>

        <!-- Code Content -->
        <pre class="overflow-x-auto p-4 font-mono text-xs leading-relaxed text-gray-200"><code>{{ currentCurlCommand }}</code></pre>
      </div>

      <!-- Postman Import Step-by-Step Guide -->
      <div class="rounded-xl border border-amber-200 bg-amber-50/70 p-4 dark:border-amber-900/40 dark:bg-amber-950/20">
        <div class="flex items-start gap-3">
          <div class="mt-0.5 flex-shrink-0 text-amber-600 dark:text-amber-400">
            <Icon name="lightbulb" size="md" />
          </div>
          <div class="space-y-1.5 text-xs text-amber-900 dark:text-amber-200">
            <p class="font-semibold text-amber-950 dark:text-amber-100">
              {{ t('keys.postmanModal.postmanTipTitle') }}
            </p>
            <ol class="list-decimal space-y-1 pl-4 text-amber-800 dark:text-amber-300">
              <li>{{ t('keys.postmanModal.postmanStep1') }}</li>
              <li>{{ t('keys.postmanModal.postmanStep2') }}</li>
              <li>{{ t('keys.postmanModal.postmanStep3') }}</li>
              <li>{{ t('keys.postmanModal.postmanStep4') }}</li>
            </ol>
          </div>
        </div>
      </div>
    </div>
  </BaseDialog>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import BaseDialog from '@/components/common/BaseDialog.vue'
import Icon from '@/components/icons/Icon.vue'
import { useClipboard } from '@/composables/useClipboard'

const props = withDefaults(
  defineProps<{
    show: boolean
    apiKey: string
    keyName?: string
    baseUrl: string
    platform?: string | null
    groupId?: number | null
    groupName?: string
    allowedModels?: string[]
  }>(),
  {
    keyName: '',
    platform: null,
    groupId: null,
    groupName: '',
    allowedModels: () => []
  }
)

const emit = defineEmits<{
  (e: 'close'): void
}>()

const { t } = useI18n()
const { copyToClipboard } = useClipboard()

// Effective Base URL
const defaultBaseUrl = computed(() => {
  if (props.baseUrl && props.baseUrl.trim() !== '') {
    return props.baseUrl.replace(/\/+$/, '')
  }
  return window.location.origin
})

const customBaseUrl = ref('')
const activeBaseUrl = computed(() => {
  if (customBaseUrl.value && customBaseUrl.value.trim() !== '') {
    return customBaseUrl.value.trim().replace(/\/+$/, '')
  }
  return defaultBaseUrl.value
})

// Models
const defaultFallbackModels = [
  'gemini-3.8-flash-high',
  'gemini-3.7-flash-high',
  'claude-sonnet-4-6',
  'gpt-4o-mini',
  'gpt-4o'
]

const availableModels = computed(() => {
  if (props.allowedModels && props.allowedModels.length > 0) {
    return props.allowedModels
  }
  return defaultFallbackModels
})

// Protocol tabs
type ProtocolId = 'openaiChat' | 'claudeMessages' | 'geminiNative' | 'openaiResponses'
const activeProtocol = ref<ProtocolId>('openaiChat')

const pickRecommendedModelForProtocol = (proto: ProtocolId, current: string): string => {
  const models = availableModels.value
  if (!models || models.length === 0) return current

  const findMatch = (keywords: string[]) => {
    for (const kw of keywords) {
      const match = models.find(m => m.toLowerCase().includes(kw))
      if (match) return match
    }
    return null
  }

  if (proto === 'claudeMessages') {
    if (current && current.toLowerCase().includes('claude')) {
      return current
    }
    return findMatch(['claude', 'sonnet', 'opus', 'haiku']) || 'claude-sonnet-4-6'
  }

  if (proto === 'geminiNative') {
    if (current && current.toLowerCase().includes('gemini') && current !== 'gemini-3-flash') {
      return current
    }
    return findMatch(['gemini-3.8-flash-high', 'gemini-3.7-flash-high', 'gemini-3.6-flash-high', 'gemini-3-pro', 'gemini']) || 'gemini-3.8-flash-high'
  }

  if (proto === 'openaiChat' || proto === 'openaiResponses') {
    if (current && models.includes(current) && current !== 'gemini-3-flash') {
      return current
    }
    return findMatch(['gemini-3.8-flash-high', 'gemini-3.7-flash-high', 'gpt-4o', 'gemini', 'gpt']) || models[0] || 'gemini-3.8-flash-high'
  }

  return current
}

const selectedModel = ref('')
const customModelInput = ref('')

// Initialize model selection
watch(
  () => props.allowedModels,
  (models) => {
    if (models && models.length > 0) {
      selectedModel.value = pickRecommendedModelForProtocol(activeProtocol.value, models[0])
    } else {
      selectedModel.value = pickRecommendedModelForProtocol(activeProtocol.value, defaultFallbackModels[0])
    }
  },
  { immediate: true }
)

// Auto switch recommended model when tab changes
watch(activeProtocol, (newProto) => {
  if (selectedModel.value !== '__custom__') {
    selectedModel.value = pickRecommendedModelForProtocol(newProto, selectedModel.value)
  }
})

const effectiveModel = computed(() => {
  if (selectedModel.value === '__custom__') {
    return customModelInput.value.trim() || 'gemini-3.8-flash-high'
  }
  return selectedModel.value || 'gemini-3.8-flash-high'
})

const protocols = computed(() => [
  {
    id: 'openaiChat' as ProtocolId,
    badge: 'OpenAI',
    label: t('keys.postmanModal.protocolTabs.openaiChat')
  },
  {
    id: 'claudeMessages' as ProtocolId,
    badge: 'Claude',
    label: t('keys.postmanModal.protocolTabs.claudeMessages')
  },
  {
    id: 'geminiNative' as ProtocolId,
    badge: 'Gemini',
    label: t('keys.postmanModal.protocolTabs.geminiNative')
  },
  {
    id: 'openaiResponses' as ProtocolId,
    badge: 'Responses',
    label: t('keys.postmanModal.protocolTabs.openaiResponses')
  }
])

// Masked key display
const maskedKey = computed(() => {
  if (!props.apiKey) return ''
  if (props.apiKey.length <= 12) return props.apiKey
  return `${props.apiKey.slice(0, 7)}...${props.apiKey.slice(-4)}`
})

// Protocol configs & cURL generation
const currentProtocolConfig = computed(() => {
  const base = activeBaseUrl.value
  const model = effectiveModel.value
  const key = props.apiKey

  switch (activeProtocol.value) {
    case 'openaiChat':
      return {
        endpoint: `${base}/v1/chat/completions`,
        note: t('keys.postmanModal.notes.chat'),
        bodyObj: {
          model,
          messages: [
            {
              role: 'user',
              content: 'Hello, how are you?'
            }
          ],
          stream: false
        }
      }
    case 'claudeMessages':
      return {
        endpoint: `${base}/v1/messages`,
        note: t('keys.postmanModal.notes.messages'),
        bodyObj: {
          model,
          max_tokens: 1024,
          messages: [
            {
              role: 'user',
              content: 'Hello, how are you?'
            }
          ]
        }
      }
    case 'geminiNative':
      return {
        endpoint: `${base}/v1beta/models/${model}:generateContent?key=${key}`,
        note: t('keys.postmanModal.notes.gemini'),
        bodyObj: {
          contents: [
            {
              role: 'user',
              parts: [
                {
                  text: 'Hello, how are you?'
                }
              ]
            }
          ]
        }
      }
    case 'openaiResponses':
      return {
        endpoint: `${base}/v1/responses`,
        note: t('keys.postmanModal.notes.responses'),
        bodyObj: {
          model,
          input: 'Hello, how are you?'
        }
      }
  }
})

const currentCurlCommand = computed(() => {
  const base = activeBaseUrl.value
  const model = effectiveModel.value
  const key = props.apiKey

  switch (activeProtocol.value) {
    case 'openaiChat':
      return `curl -X POST "${base}/v1/chat/completions" \\
  -H "Authorization: Bearer ${key}" \\
  -H "Content-Type: application/json" \\
  -d '{
    "model": "${model}",
    "messages": [
      {
        "role": "user",
        "content": "Hello, how are you?"
      }
    ],
    "stream": false
  }'`

    case 'claudeMessages':
      return `curl -X POST "${base}/v1/messages" \\
  -H "x-api-key: ${key}" \\
  -H "anthropic-version: 2023-06-01" \\
  -H "Content-Type: application/json" \\
  -d '{
    "model": "${model}",
    "max_tokens": 1024,
    "messages": [
      {
        "role": "user",
        "content": "Hello, how are you?"
      }
    ]
  }'`

    case 'geminiNative':
      return `curl -X POST "${base}/v1beta/models/${model}:generateContent?key=${key}" \\
  -H "Content-Type: application/json" \\
  -d '{
    "contents": [
      {
        "role": "user",
        "parts": [
          {
            "text": "Hello, how are you?"
          }
        ]
      }
    ]
  }'`

    case 'openaiResponses':
      return `curl -X POST "${base}/v1/responses" \\
  -H "Authorization: Bearer ${key}" \\
  -H "Content-Type: application/json" \\
  -d '{
    "model": "${model}",
    "input": "Hello, how are you?"
  }'`
  }
})

// Copy states
const copiedKey = ref(false)
const copiedCurlState = ref(false)
const copiedBodyState = ref(false)

async function copyApiKey() {
  if (!props.apiKey) return
  const success = await copyToClipboard(props.apiKey)
  if (success) {
    copiedKey.value = true
    setTimeout(() => {
      copiedKey.value = false
    }, 2000)
  }
}

async function copyCurl() {
  const success = await copyToClipboard(currentCurlCommand.value)
  if (success) {
    copiedCurlState.value = true
    setTimeout(() => {
      copiedCurlState.value = false
    }, 2000)
  }
}

async function copyBody() {
  const jsonStr = JSON.stringify(currentProtocolConfig.value.bodyObj, null, 2)
  const success = await copyToClipboard(jsonStr)
  if (success) {
    copiedBodyState.value = true
    setTimeout(() => {
      copiedBodyState.value = false
    }, 2000)
  }
}
</script>
