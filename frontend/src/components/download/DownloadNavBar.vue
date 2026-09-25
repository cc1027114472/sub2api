<template>
  <header
    class="sticky top-0 z-40 border-b border-slate-200/80 bg-white/80 backdrop-blur-md transition-colors dark:border-white/10 dark:bg-dark-900/80"
  >
    <div class="mx-auto flex max-w-7xl items-center justify-between gap-4 px-4 py-3 sm:px-6">
      <!-- 左侧: 魔丸 Logo + 名称 + 版本徽章 -->
      <router-link to="/download" class="group flex min-w-0 items-center gap-3">
        <img
          src="/logo.svg?v=2.1.0"
          alt="Logo"
          class="h-9 w-9 shrink-0 object-contain drop-shadow-sm transition-transform duration-200 group-hover:scale-105"
        />
        <div class="flex items-center gap-2">
          <span class="truncate text-lg font-black tracking-tight text-slate-950 dark:text-white">
            {{ siteName }}
          </span>
          <span class="rounded-full border border-orange-200 bg-orange-100 px-2 py-0.5 text-[11px] font-bold text-orange-700 dark:border-orange-500/30 dark:bg-orange-950/50 dark:text-orange-300">
            v2.1.0 官方版
          </span>
        </div>
      </router-link>

      <!-- 中间快速锚点导航 (桌面端) -->
      <nav class="hidden lg:flex items-center gap-6 text-sm font-semibold text-slate-600 dark:text-dark-300">
        <a
          href="#models"
          class="text-indigo-600 transition hover:text-orange-600 dark:text-indigo-400 dark:hover:text-orange-400"
        >
          全模型矩阵·专精Gemini
        </a>
        <a
          href="#lay-flat"
          class="transition hover:text-orange-600 dark:hover:text-orange-400"
        >
          手机躺平操控
        </a>
        <a
          href="#hardcore-proof"
          class="transition hover:text-orange-600 dark:hover:text-orange-400"
        >
          长任务实测
        </a>
        <a
          href="#tour"
          class="transition hover:text-orange-600 dark:hover:text-orange-400"
        >
          功能全景
        </a>
        <a
          href="#guide"
          class="transition hover:text-orange-600 dark:hover:text-orange-400"
        >
          上手指南
        </a>
      </nav>

      <!-- 右侧: 主题切换 + 登录/控制台 + 免费下载主按钮 -->
      <div class="flex items-center gap-3">
        <!-- 主题切换 -->
        <button
          type="button"
          class="flex h-9 w-9 items-center justify-center rounded-xl border border-slate-200 bg-white text-slate-600 shadow-sm transition hover:bg-slate-100 dark:border-white/10 dark:bg-dark-800 dark:text-dark-300 dark:hover:bg-dark-700"
          :title="isDark ? '切换至亮色模式' : '切换至暗色模式'"
          @click="toggleTheme"
        >
          <Icon :name="isDark ? 'sun' : 'moon'" size="xs" />
        </button>

        <router-link
          v-if="isAuthenticated"
          :to="backTarget"
          class="hidden sm:inline-flex shrink-0 items-center justify-center gap-1.5 rounded-xl border border-slate-200 bg-white px-3.5 py-2 text-xs font-bold text-slate-700 shadow-sm transition hover:bg-slate-50 dark:border-white/10 dark:bg-dark-800 dark:text-dark-200 dark:hover:bg-dark-700"
        >
          <Icon name="arrowLeft" size="xs" />
          <span>{{ t('agentDownload.nav.backToDashboard') }}</span>
        </router-link>
        <router-link
          v-else
          :to="{ path: '/login', query: { redirect: '/download' } }"
          class="hidden sm:inline-flex shrink-0 items-center justify-center gap-1.5 rounded-xl border border-slate-200 bg-white px-3.5 py-2 text-xs font-bold text-slate-700 shadow-sm transition hover:bg-slate-50 dark:border-white/10 dark:bg-dark-800 dark:text-dark-200 dark:hover:bg-dark-700"
        >
          <span>{{ t('agentDownload.nav.login') }}</span>
        </router-link>

        <a
          href="#download-section"
          class="inline-flex shrink-0 items-center justify-center gap-1.5 rounded-xl bg-gradient-to-r from-orange-600 via-amber-600 to-rose-600 px-4 py-2 text-xs sm:text-sm font-bold text-white shadow-md shadow-orange-500/20 transition hover:opacity-95 active:scale-[0.98]"
        >
          <span>免费下载</span>
        </a>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import Icon from '@/components/icons/Icon.vue'
import { useAppStore } from '@/stores/app'
import { useAuthStore } from '@/stores/auth'

const { t } = useI18n()
const appStore = useAppStore()
const authStore = useAuthStore()

const settings = computed(() => appStore.cachedPublicSettings)
const siteName = computed(() => (settings.value?.site_name ? `魔丸 Agent (${settings.value.site_name})` : '魔丸 Agent'))
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
