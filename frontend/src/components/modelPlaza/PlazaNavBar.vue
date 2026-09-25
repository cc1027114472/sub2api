<template>
  <header
    class="sticky top-0 z-40 border-b border-gray-200/80 bg-white/75 backdrop-blur-xl transition-colors dark:border-white/10 dark:bg-dark-900/75"
  >
    <div class="mx-auto flex max-w-7xl items-center justify-between gap-4 px-4 py-3 sm:px-6">
      <!-- 左:站点 logo + 名称(点击返回首页) -->
      <RouterLink to="/home" class="group flex min-w-0 items-center gap-3">
        <template v-if="settings">
          <div
            class="flex h-9 w-9 shrink-0 items-center justify-center overflow-hidden rounded-xl border border-gray-200/80 bg-white p-1 shadow-sm transition-transform duration-200 group-hover:scale-105 dark:border-white/10 dark:bg-dark-800"
          >
            <img :src="siteLogo || '/logo.svg?v=2.1.0'" alt="Logo" class="h-full w-full object-contain" />
          </div>
          <div class="flex flex-col">
            <span class="truncate text-sm font-bold tracking-tight text-gray-900 dark:text-white">
              {{ siteName }}
            </span>
            <span class="text-[10px] font-mono tracking-widest uppercase text-primary-500 font-semibold">MODEL PLAZA</span>
          </div>
        </template>
        <template v-else>
          <span class="h-9 w-9 shrink-0 animate-pulse rounded-xl bg-gray-200 dark:bg-dark-700" aria-hidden="true"></span>
          <span class="h-5 w-24 animate-pulse rounded bg-gray-200 dark:bg-dark-700" aria-hidden="true"></span>
        </template>
      </RouterLink>

      <!-- 中间快速链接 (桌面端) -->
      <nav class="hidden md:flex items-center gap-1 text-xs font-medium text-gray-600 dark:text-dark-300">
        <RouterLink
          to="/home"
          class="rounded-lg px-3 py-1.5 transition-colors hover:bg-gray-100 hover:text-gray-900 dark:hover:bg-dark-800 dark:hover:text-white"
        >
          {{ t('nav.home') || '首页门户' }}
        </RouterLink>
        <RouterLink
          v-if="showAgentDownload"
          to="/download"
          class="rounded-lg px-3 py-1.5 transition-colors hover:bg-gray-100 hover:text-gray-900 dark:hover:bg-dark-800 dark:hover:text-white"
        >
          {{ t('nav.agentDownload') || '客户端下载' }}
        </RouterLink>
        <RouterLink
          to="/key-usage"
          class="rounded-lg px-3 py-1.5 transition-colors hover:bg-gray-100 hover:text-gray-900 dark:hover:bg-dark-800 dark:hover:text-white"
        >
          {{ t('keyUsage.title') || 'API Key 用量查询' }}
        </RouterLink>
        <a
          v-if="docUrl"
          :href="docUrl"
          target="_blank"
          rel="noopener noreferrer"
          class="rounded-lg px-3 py-1.5 transition-colors hover:bg-gray-100 hover:text-gray-900 dark:hover:bg-dark-800 dark:hover:text-white"
        >
          {{ t('home.viewDocs') || '开发文档' }}
        </a>
      </nav>

      <!-- 右:主题切换 + 登录 / 回到后台 -->
      <div class="flex items-center gap-2.5">
        <!-- 主题切换 -->
        <button
          type="button"
          class="flex h-9 w-9 items-center justify-center rounded-xl border border-gray-200/80 bg-white/60 text-gray-500 shadow-sm transition hover:bg-gray-100 dark:border-white/10 dark:bg-dark-800/60 dark:text-dark-300 dark:hover:bg-dark-700"
          :title="isDark ? t('home.switchToLight') : t('home.switchToDark')"
          @click="toggleTheme"
        >
          <Icon :name="isDark ? 'sun' : 'moon'" size="xs" />
        </button>

        <!-- 登录 / 控制台按钮 -->
        <RouterLink
          v-if="isAuthenticated"
          :to="backTarget"
          class="inline-flex shrink-0 items-center justify-center gap-1.5 rounded-xl bg-gradient-to-r from-primary-600 to-indigo-600 px-4 py-2 text-xs font-semibold text-white shadow-md shadow-primary-500/20 transition-all duration-200 hover:from-primary-500 hover:to-indigo-500 hover:shadow-lg hover:shadow-primary-500/30 active:scale-[0.98]"
        >
          <Icon name="cog" size="xs" />
          <span>{{ t('modelPlaza.nav.backToDashboard') }}</span>
        </RouterLink>
        <RouterLink
          v-else
          :to="{ path: '/login', query: { redirect: '/model-plaza' } }"
          class="inline-flex shrink-0 items-center justify-center gap-1.5 rounded-xl bg-gradient-to-r from-primary-600 to-indigo-600 px-4 py-2 text-xs font-semibold text-white shadow-md shadow-primary-500/20 transition-all duration-200 hover:from-primary-500 hover:to-indigo-500 hover:shadow-lg hover:shadow-primary-500/30 active:scale-[0.98]"
        >
          <span>{{ t('modelPlaza.nav.login') }}</span>
          <svg class="h-3.5 w-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3" />
          </svg>
        </RouterLink>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import Icon from '@/components/icons/Icon.vue'
import { sanitizeUrl } from '@/utils/url'
import { useAppStore } from '@/stores/app'
import { useAuthStore } from '@/stores/auth'

const { t } = useI18n()
const appStore = useAppStore()
const authStore = useAuthStore()

const settings = computed(() => appStore.cachedPublicSettings)
const siteName = computed(() => appStore.siteName || '英国api.cc')
const siteLogo = computed(() =>
  sanitizeUrl(settings.value?.site_logo || '', { allowRelative: true, allowDataUrl: true })
)
const docUrl = computed(() => appStore.docUrl)
const showAgentDownload = computed(() => {
  const val = (settings.value as Record<string, unknown> | null)?.agent_download_enabled
  return val === undefined || Boolean(val)
})

const isAuthenticated = computed(() => authStore.isAuthenticated)
const backTarget = computed(() => (authStore.isAdmin ? '/admin/dashboard' : '/dashboard'))

// 主题状态
const isDark = ref(false)
const initTheme = () => {
  isDark.value =
    localStorage.theme === 'dark' ||
    (!('theme' in localStorage) && window.matchMedia('(prefers-color-scheme: dark)').matches)
}
const toggleTheme = () => {
  if (isDark.value) {
    document.documentElement.classList.remove('dark')
    localStorage.theme = 'light'
    isDark.value = false
  } else {
    document.documentElement.classList.add('dark')
    localStorage.theme = 'dark'
    isDark.value = true
  }
}

onMounted(() => {
  initTheme()
})
</script>
