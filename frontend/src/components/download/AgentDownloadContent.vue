<template>
  <div class="space-y-12">
    <!-- Hero 区域 -->
    <div class="relative overflow-hidden rounded-3xl border border-gray-200/80 bg-gradient-to-b from-primary-50/60 via-white to-white p-6 shadow-sm dark:border-dark-700/60 dark:from-primary-950/20 dark:via-dark-900 dark:to-dark-900 sm:p-12">
      <!-- 柔和背景光晕 -->
      <div class="pointer-events-none absolute -top-24 left-1/2 h-96 w-96 -translate-x-1/2 rounded-full bg-primary-400/15 blur-3xl dark:bg-primary-500/10"></div>

      <div class="relative z-10 flex flex-col items-center text-center">
        <!-- 品牌 Logo 徽章 -->
        <div class="relative mb-6">
          <div class="flex h-24 w-24 items-center justify-center rounded-2xl bg-white p-3 shadow-xl ring-1 ring-gray-900/5 transition-transform duration-300 hover:scale-105 dark:bg-dark-800 dark:ring-white/10 sm:h-28 sm:w-28">
            <img
              src="/assets/agent/mowan-logo-128.png"
              alt="mowan-harness Logo"
              class="h-full w-full object-contain"
            />
          </div>
          <span class="absolute -bottom-2.5 left-1/2 -translate-x-1/2 whitespace-nowrap rounded-full bg-primary-600 px-3 py-0.5 text-xs font-semibold text-white shadow-sm dark:bg-primary-500">
            {{ t('agentDownload.badge') }}
          </span>
        </div>

        <!-- 标题与简介 -->
        <h1 class="text-3xl font-extrabold tracking-tight text-gray-900 dark:text-white sm:text-4xl lg:text-5xl">
          {{ t('agentDownload.title') }}
        </h1>
        <p class="mx-auto mt-4 max-w-2xl text-base text-gray-600 dark:text-dark-300 sm:text-lg">
          {{ t('agentDownload.subtitle') }}
        </p>

        <!-- 系统匹配提示 -->
        <div class="mt-6 flex flex-wrap items-center justify-center gap-2 text-xs text-gray-500 dark:text-dark-400">
          <span>{{ t('agentDownload.detectedPlatform') }}</span>
          <span class="inline-flex items-center gap-1 rounded-md bg-gray-100 px-2 py-0.5 font-medium text-gray-800 dark:bg-dark-800 dark:text-dark-200">
            <Icon name="check" size="xs" class="text-emerald-500" />
            {{ detectedPlatform.platformLabel }}
          </span>
          <span class="hidden sm:inline">·</span>
          <span>{{ t('agentDownload.version') }} {{ downloadConfig.version }}</span>
          <span class="hidden sm:inline">·</span>
          <span>{{ t('agentDownload.fileSize') }} {{ downloadConfig.windows.size }}</span>
        </div>

        <!-- 下载行动主按钮组 -->
        <div class="mt-8 flex flex-col items-center justify-center gap-4 sm:flex-row">
          <a
            :href="detectedPlatform.primaryUrl"
            download
            class="group inline-flex w-full items-center justify-center gap-2.5 rounded-2xl bg-gradient-to-r from-primary-600 via-primary-500 to-indigo-600 px-8 py-4 text-base font-bold text-white shadow-lg shadow-primary-500/25 transition-all duration-200 hover:from-primary-500 hover:to-indigo-500 hover:shadow-xl hover:shadow-primary-500/30 active:scale-[0.98] sm:w-auto"
          >
            <Icon name="download" size="md" class="transition-transform group-hover:-translate-y-0.5" />
            <span>{{ t('agentDownload.downloadPrimary') }}</span>
          </a>

          <a
            v-if="downloadConfig.windows.portableUrl"
            :href="downloadConfig.windows.portableUrl"
            download
            class="inline-flex w-full items-center justify-center gap-2 rounded-2xl border border-gray-300 bg-white px-6 py-4 text-base font-semibold text-gray-700 shadow-sm transition-colors hover:bg-gray-50 hover:text-gray-900 dark:border-dark-700 dark:bg-dark-800 dark:text-dark-200 dark:hover:bg-dark-700 dark:hover:text-white sm:w-auto"
          >
            <Icon name="folder" size="md" />
            <span>{{ t('agentDownload.downloadPortable') }}</span>
          </a>
        </div>

        <!-- 环境说明提醒 -->
        <div class="mt-6 flex items-center justify-center gap-1.5 text-xs text-emerald-600 dark:text-emerald-400">
          <Icon name="checkCircle" size="xs" />
          <span>{{ t('agentDownload.envIncluded') }}</span>
        </div>
      </div>
    </div>

    <!-- 客户端功能预览与上手指南 (方案 A: 图文全景流) -->
    <AgentGuideTour />

    <!-- 快速接入配置助手 -->
    <div class="rounded-3xl border border-gray-200/80 bg-white p-6 shadow-sm dark:border-dark-700/60 dark:bg-dark-900 sm:p-8">
      <div class="flex flex-col justify-between gap-4 border-b border-gray-100 pb-6 dark:border-dark-800 sm:flex-row sm:items-center">
        <div>
          <div class="flex items-center gap-2">
            <span class="flex h-8 w-8 items-center justify-center rounded-lg bg-primary-100 text-primary-600 dark:bg-primary-900/40 dark:text-primary-400">
              <Icon name="cog" size="sm" />
            </span>
            <h2 class="text-xl font-bold text-gray-900 dark:text-white">
              {{ t('agentDownload.quickStartTitle') }}
            </h2>
          </div>
          <p class="mt-1 text-sm text-gray-500 dark:text-dark-400">
            {{ t('agentDownload.quickStartDesc') }}
          </p>
        </div>

        <!-- 3步极简流程标签 -->
        <div class="flex flex-wrap items-center gap-2 text-xs text-gray-600 dark:text-dark-400">
          <span class="rounded-md bg-gray-100 px-2 py-1 dark:bg-dark-800">{{ t('agentDownload.step1') }}</span>
          <span>→</span>
          <span class="rounded-md bg-primary-50 px-2 py-1 font-medium text-primary-700 dark:bg-primary-950/40 dark:text-primary-300">{{ t('agentDownload.step2') }}</span>
          <span>→</span>
          <span class="rounded-md bg-gray-100 px-2 py-1 dark:bg-dark-800">{{ t('agentDownload.step3') }}</span>
        </div>
      </div>

      <!-- 配置操作区 -->
      <div class="mt-6 space-y-6">
        <!-- 基础接口地址 -->
        <div>
          <label class="block text-xs font-semibold uppercase tracking-wider text-gray-500 dark:text-dark-400">
            {{ t('agentDownload.apiBaseUrl') }}
          </label>
          <div class="mt-2 flex items-center gap-2">
            <div class="relative flex-1">
              <input
                type="text"
                readonly
                :value="currentApiBaseUrl"
                class="w-full rounded-xl border border-gray-200 bg-gray-50 px-4 py-2.5 font-mono text-sm text-gray-800 focus:outline-none dark:border-dark-700 dark:bg-dark-800 dark:text-dark-200"
              />
            </div>
            <button
              type="button"
              class="inline-flex items-center gap-1.5 rounded-xl border border-gray-200 bg-white px-3.5 py-2.5 text-xs font-semibold text-gray-700 shadow-sm transition-colors hover:bg-gray-50 dark:border-dark-700 dark:bg-dark-800 dark:text-dark-200 dark:hover:bg-dark-700"
              @click="copyText(currentApiBaseUrl)"
            >
              <Icon :name="copiedKey === currentApiBaseUrl ? 'check' : 'clipboard'" size="xs" />
              <span>{{ copiedKey === currentApiBaseUrl ? t('common.copied') : t('agentDownload.copyUrlBtn') }}</span>
            </button>
          </div>
        </div>

        <!-- API Key 选择与联动 -->
        <div v-if="isAuthenticated">
          <label class="block text-xs font-semibold uppercase tracking-wider text-gray-500 dark:text-dark-400">
            {{ t('agentDownload.apiKeySelect') }}
          </label>
          <div class="mt-2 flex flex-col gap-2 sm:flex-row sm:items-center">
            <div v-if="userKeys.length > 0" class="flex-1">
              <select
                v-model="selectedKeyId"
                class="w-full rounded-xl border border-gray-200 bg-white px-4 py-2.5 text-sm text-gray-900 focus:border-primary-500 focus:outline-none focus:ring-1 focus:ring-primary-500 dark:border-dark-700 dark:bg-dark-800 dark:text-dark-200"
              >
                <option v-for="key in userKeys" :key="key.id" :value="key.id">
                  {{ key.name }} ({{ maskKey(key.key) }})
                </option>
              </select>
            </div>
            <div v-else class="flex flex-1 items-center justify-between rounded-xl border border-amber-200 bg-amber-50 px-4 py-2.5 text-xs text-amber-800 dark:border-amber-900/50 dark:bg-amber-950/30 dark:text-amber-200">
              <span>{{ t('agentDownload.noKeyHint') }}</span>
              <router-link
                to="/keys"
                class="font-semibold underline hover:text-amber-900 dark:hover:text-amber-100"
              >
                {{ t('agentDownload.createKeyBtn') }}
              </router-link>
            </div>

            <button
              v-if="selectedKeySecret"
              type="button"
              class="inline-flex items-center gap-1.5 rounded-xl border border-gray-200 bg-white px-3.5 py-2.5 text-xs font-semibold text-gray-700 shadow-sm transition-colors hover:bg-gray-50 dark:border-dark-700 dark:bg-dark-800 dark:text-dark-200 dark:hover:bg-dark-700"
              @click="copyText(selectedKeySecret)"
            >
              <Icon :name="copiedKey === selectedKeySecret ? 'check' : 'clipboard'" size="xs" />
              <span>{{ copiedKey === selectedKeySecret ? t('common.copied') : t('agentDownload.copyKeyBtn') }}</span>
            </button>
          </div>
        </div>

        <!-- 未登录引导 -->
        <div v-else class="flex flex-col items-center justify-between gap-4 rounded-2xl border border-primary-100 bg-primary-50/50 p-4 dark:border-primary-900/30 dark:bg-primary-950/20 sm:flex-row">
          <div class="flex items-center gap-2.5 text-sm text-primary-900 dark:text-primary-200">
            <Icon name="infoCircle" size="sm" class="text-primary-600 dark:text-primary-400" />
            <span>{{ t('agentDownload.loginPrompt') }}</span>
          </div>
          <router-link
            :to="{ path: '/login', query: { redirect: '/download' } }"
            class="inline-flex shrink-0 items-center justify-center rounded-xl bg-primary-600 px-4 py-2 text-xs font-semibold text-white shadow-sm hover:bg-primary-700"
          >
            {{ t('agentDownload.loginNow') }}
          </router-link>
        </div>

        <!-- 整合配置代码预览 -->
        <div class="relative">
          <div class="flex items-center justify-between rounded-t-xl bg-gray-900 px-4 py-2 text-xs text-gray-400">
            <span class="font-mono">config.yaml (Agent Connection)</span>
            <button
              type="button"
              class="flex items-center gap-1 text-xs text-gray-300 transition-colors hover:text-white"
              @click="copyText(generatedSnippet)"
            >
              <Icon :name="copiedKey === generatedSnippet ? 'check' : 'clipboard'" size="xs" />
              <span>{{ copiedKey === generatedSnippet ? t('common.copied') : t('agentDownload.copyBtn') }}</span>
            </button>
          </div>
          <pre class="overflow-x-auto rounded-b-xl bg-gray-950 p-4 font-mono text-xs leading-relaxed text-gray-200">{{ generatedSnippet }}</pre>
        </div>
      </div>
    </div>

    <!-- 多平台下载矩阵 -->
    <div class="space-y-4">
      <h2 class="text-xl font-bold text-gray-900 dark:text-white">
        {{ t('agentDownload.manualSelect') }}
      </h2>

      <div class="grid grid-cols-1 gap-4 md:grid-cols-3">
        <!-- Windows -->
        <div class="relative flex flex-col justify-between rounded-2xl border-2 border-primary-500/80 bg-white p-6 shadow-sm dark:bg-dark-900">
          <div class="absolute -top-3 right-4 rounded-full bg-primary-500 px-2.5 py-0.5 text-[10px] font-bold uppercase tracking-wider text-white">
            推荐 / 主推
          </div>
          <div>
            <div class="flex h-12 w-12 items-center justify-center rounded-xl bg-primary-50 text-primary-600 dark:bg-primary-950/40 dark:text-primary-400">
              <Icon name="desktop" size="md" />
            </div>
            <h3 class="mt-4 text-base font-bold text-gray-900 dark:text-white">
              {{ t('agentDownload.windowsTitle') }}
            </h3>
            <p class="mt-1 text-xs text-gray-500 dark:text-dark-400">
              {{ t('agentDownload.windowsDesc') }}
            </p>
          </div>
          <div class="mt-6 flex flex-col gap-2">
            <a
              :href="downloadConfig.windows.installerUrl"
              download
              class="inline-flex items-center justify-center gap-2 rounded-xl bg-primary-600 px-4 py-2.5 text-xs font-semibold text-white shadow-sm transition-colors hover:bg-primary-700"
            >
              <Icon name="download" size="xs" />
              {{ t('agentDownload.windowsInstaller') }}
            </a>
            <a
              v-if="downloadConfig.windows.portableUrl"
              :href="downloadConfig.windows.portableUrl"
              download
              class="inline-flex items-center justify-center gap-2 rounded-xl border border-gray-200 bg-white px-4 py-2 text-xs font-medium text-gray-700 hover:bg-gray-50 dark:border-dark-700 dark:bg-dark-800 dark:text-dark-200"
            >
              {{ t('agentDownload.windowsPortable') }}
            </a>
          </div>
        </div>

        <!-- macOS -->
        <div class="flex flex-col justify-between rounded-2xl border border-gray-200/80 bg-white p-6 shadow-sm dark:border-dark-700/60 dark:bg-dark-900">
          <div>
            <div class="flex h-12 w-12 items-center justify-center rounded-xl bg-gray-100 text-gray-600 dark:bg-dark-800 dark:text-dark-300">
              <Icon name="deviceMobile" size="md" />
            </div>
            <h3 class="mt-4 text-base font-bold text-gray-900 dark:text-white">
              {{ t('agentDownload.macTitle') }}
            </h3>
            <p class="mt-1 text-xs text-gray-500 dark:text-dark-400">
              {{ t('agentDownload.macDesc') }}
            </p>
          </div>
          <div class="mt-6">
            <span class="inline-flex w-full items-center justify-center rounded-xl bg-gray-100 px-4 py-2.5 text-xs font-medium text-gray-500 dark:bg-dark-800 dark:text-dark-400">
              {{ t('agentDownload.macComingSoon') }}
            </span>
          </div>
        </div>

        <!-- Linux / Docker -->
        <div class="flex flex-col justify-between rounded-2xl border border-gray-200/80 bg-white p-6 shadow-sm dark:border-dark-700/60 dark:bg-dark-900">
          <div>
            <div class="flex h-12 w-12 items-center justify-center rounded-xl bg-gray-100 text-gray-600 dark:bg-dark-800 dark:text-dark-300">
              <Icon name="terminal" size="md" />
            </div>
            <h3 class="mt-4 text-base font-bold text-gray-900 dark:text-white">
              {{ t('agentDownload.linuxTitle') }}
            </h3>
            <p class="mt-1 text-xs text-gray-500 dark:text-dark-400">
              {{ t('agentDownload.linuxDesc') }}
            </p>
          </div>
          <div class="mt-6">
            <a
              :href="downloadConfig.linux.cliUrl"
              target="_blank"
              rel="noopener noreferrer"
              class="inline-flex w-full items-center justify-center gap-1.5 rounded-xl border border-gray-200 bg-white px-4 py-2.5 text-xs font-medium text-gray-700 transition-colors hover:bg-gray-50 dark:border-dark-700 dark:bg-dark-800 dark:text-dark-200"
            >
              <Icon name="codeBracket" size="xs" />
              <span>GitHub / Docker</span>
            </a>
          </div>
        </div>
      </div>
    </div>

    <!-- 4大核心特性介绍 -->
    <div class="space-y-6">
      <div class="text-center">
        <h2 class="text-2xl font-bold tracking-tight text-gray-900 dark:text-white sm:text-3xl">
          {{ t('agentDownload.features.title') }}
        </h2>
      </div>

      <div class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4">
        <!-- Feature 1: 全明星模型与 Gemini 旗舰深度调优 -->
        <div class="rounded-2xl border border-gray-200/70 bg-white p-6 shadow-sm transition-all hover:shadow-md dark:border-dark-800 dark:bg-dark-900">
          <div class="flex h-11 w-11 items-center justify-center rounded-xl bg-indigo-100 text-indigo-600 dark:bg-indigo-950/40 dark:text-indigo-400">
            <Icon name="bolt" size="sm" />
          </div>
          <h3 class="mt-4 text-base font-bold text-gray-900 dark:text-white">
            {{ t('agentDownload.features.modelsTitle') }}
          </h3>
          <p class="mt-2 text-xs leading-relaxed text-gray-500 dark:text-dark-400">
            {{ t('agentDownload.features.modelsDesc') }}
          </p>
        </div>

        <!-- Feature 2: 自愈 Ralph 循环 -->
        <div class="rounded-2xl border border-gray-200/70 bg-white p-6 shadow-sm transition-all hover:shadow-md dark:border-dark-800 dark:bg-dark-900">
          <div class="flex h-11 w-11 items-center justify-center rounded-xl bg-primary-100 text-primary-600 dark:bg-primary-950/40 dark:text-primary-400">
            <Icon name="sync" size="sm" />
          </div>
          <h3 class="mt-4 text-base font-bold text-gray-900 dark:text-white">
            {{ t('agentDownload.features.loopTitle') }}
          </h3>
          <p class="mt-2 text-xs leading-relaxed text-gray-500 dark:text-dark-400">
            {{ t('agentDownload.features.loopDesc') }}
          </p>
        </div>

        <!-- Feature 3: 安全沙箱与门禁 -->
        <div class="rounded-2xl border border-gray-200/70 bg-white p-6 shadow-sm transition-all hover:shadow-md dark:border-dark-800 dark:bg-dark-900">
          <div class="flex h-11 w-11 items-center justify-center rounded-xl bg-emerald-100 text-emerald-600 dark:bg-emerald-950/40 dark:text-emerald-400">
            <Icon name="shield" size="sm" />
          </div>
          <h3 class="mt-4 text-base font-bold text-gray-900 dark:text-white">
            {{ t('agentDownload.features.sandboxTitle') }}
          </h3>
          <p class="mt-2 text-xs leading-relaxed text-gray-500 dark:text-dark-400">
            {{ t('agentDownload.features.sandboxDesc') }}
          </p>
        </div>

        <!-- Feature 4: 跨会话状态持久化 -->
        <div class="rounded-2xl border border-gray-200/70 bg-white p-6 shadow-sm transition-all hover:shadow-md dark:border-dark-800 dark:bg-dark-900">
          <div class="flex h-11 w-11 items-center justify-center rounded-xl bg-amber-100 text-amber-600 dark:bg-amber-950/40 dark:text-amber-400">
            <Icon name="database" size="sm" />
          </div>
          <h3 class="mt-4 text-base font-bold text-gray-900 dark:text-white">
            {{ t('agentDownload.features.memoryTitle') }}
          </h3>
          <p class="mt-2 text-xs leading-relaxed text-gray-500 dark:text-dark-400">
            {{ t('agentDownload.features.memoryDesc') }}
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import Icon from '@/components/icons/Icon.vue'
import AgentGuideTour from './AgentGuideTour.vue'
import { defaultAgentDownloadConfig, resolvePlatformDownload } from '@/constants/agentDownload'
import { useAuthStore } from '@/stores/auth'
import * as keysAPI from '@/api/keys'
import type { ApiKey } from '@/types'

defineProps<{
  embedded?: boolean
}>()

const { t } = useI18n()
const authStore = useAuthStore()

const downloadConfig = defaultAgentDownloadConfig
const detectedPlatform = computed(() => resolvePlatformDownload())
const isAuthenticated = computed(() => authStore.isAuthenticated)

const currentApiBaseUrl = computed(() => {
  if (typeof window !== 'undefined') {
    return `${window.location.origin}/v1`
  }
  return 'https://英国api.cc/v1'
})

const userKeys = ref<ApiKey[]>([])
const selectedKeyId = ref<number | null>(null)
const copiedKey = ref<string | null>(null)

const selectedKeySecret = computed(() => {
  const found = userKeys.value.find((k) => k.id === selectedKeyId.value)
  return found?.key || ''
})

const generatedSnippet = computed(() => {
  const keyStr = selectedKeySecret.value || 'sk-your-api-key-here'
  return `# mowan-harness (魔丸) Connection Config
api_url: "${currentApiBaseUrl.value}"
api_key: "${keyStr}"
default_model: "gemini-2.5-pro"
provider: "sub2api"`
})

function maskKey(key: string): string {
  if (!key || key.length < 10) return 'sk-***'
  return `${key.slice(0, 6)}...${key.slice(-4)}`
}

async function copyText(text: string) {
  if (!text) return
  try {
    await navigator.clipboard.writeText(text)
    copiedKey.value = text
    setTimeout(() => {
      if (copiedKey.value === text) {
        copiedKey.value = null
      }
    }, 2500)
  } catch {
    // Fallback for clipboard if needed
    const textarea = document.createElement('textarea')
    textarea.value = text
    document.body.appendChild(textarea)
    textarea.select()
    document.execCommand('copy')
    document.body.removeChild(textarea)
    copiedKey.value = text
    setTimeout(() => {
      copiedKey.value = null
    }, 2500)
  }
}

onMounted(async () => {
  if (isAuthenticated.value) {
    try {
      const res = await keysAPI.list(1, 20, { status: 'active' })
      if (res && res.items && res.items.length > 0) {
        userKeys.value = res.items
        selectedKeyId.value = res.items[0].id
      }
    } catch {
      // Non-blocking if keys cannot be loaded
    }
  }
})
</script>
