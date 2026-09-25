<template>
  <!-- 后台内嵌形态: ?embedded=1 且已登录，使用管理后台主布局 -->
  <AppLayout v-if="isEmbedded">
    <div class="px-4 py-6 sm:px-6 lg:px-8">
      <AgentDownloadContent embedded />
    </div>
  </AppLayout>

  <!-- 独立全屏形态: 自带导航条（Logo/站名 + 登录/返回后台）与 Linear 极客精细网格质感 -->
  <div
    v-else
    class="relative min-h-screen overflow-x-hidden bg-[#fafafa] text-gray-900 antialiased selection:bg-primary-500/20 selection:text-primary-900 dark:bg-[#06090e] dark:text-gray-100 dark:selection:text-primary-200"
  >
    <!-- Background: 精致网格与微光光晕 -->
    <div class="pointer-events-none fixed inset-0 z-0 overflow-hidden">
      <div
        class="absolute -top-[300px] left-1/2 h-[500px] w-[900px] -translate-x-1/2 rounded-full bg-gradient-to-b from-primary-500/15 via-cyan-500/10 to-transparent blur-[120px] dark:from-primary-600/20 dark:via-cyan-600/10"
      ></div>
      <div
        class="absolute inset-0 bg-[linear-gradient(to_right,rgba(99,102,241,0.03)_1px,transparent_1px),linear-gradient(to_bottom,rgba(99,102,241,0.03)_1px,transparent_1px)] bg-[size:44px_44px] [mask-image:radial-gradient(ellipse_70%_50%_at_50%_0%,#000_70%,transparent_100%)] dark:bg-[linear-gradient(to_right,rgba(255,255,255,0.025)_1px,transparent_1px),linear-gradient(to_bottom,rgba(255,255,255,0.025)_1px,transparent_1px)]"
      ></div>
      <div class="absolute left-0 right-0 top-0 h-[1px] bg-gradient-to-r from-transparent via-primary-500/40 to-transparent"></div>
    </div>

    <!-- 顶部导航 -->
    <div class="relative z-20">
      <DownloadNavBar />
    </div>

    <!-- 主体内容 -->
    <main class="relative z-10 mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
      <AgentDownloadContent />
    </main>

    <!-- 极简页脚 -->
    <footer class="relative z-10 border-t border-gray-200/60 py-6 text-center text-xs text-gray-500 dark:border-white/5 dark:text-dark-500">
      &copy; {{ new Date().getFullYear() }} {{ appStore.siteName }}. All rights reserved.
    </footer>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import AppLayout from '@/components/layout/AppLayout.vue'
import DownloadNavBar from '@/components/download/DownloadNavBar.vue'
import AgentDownloadContent from '@/components/download/AgentDownloadContent.vue'
import { useAppStore } from '@/stores/app'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const appStore = useAppStore()
const authStore = useAuthStore()

// embedded=1 但未登录（如复制的链接）自动降级为独立形态
const isEmbedded = computed(() => route.query.embedded === '1' && authStore.isAuthenticated)

onMounted(() => {
  void appStore.fetchPublicSettings()
})
</script>
