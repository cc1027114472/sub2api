<template>
  <!-- 后台内嵌形态: ?embedded=1 且已登录，使用管理后台主布局 -->
  <AppLayout v-if="isEmbedded">
    <div class="px-4 py-6 sm:px-6 lg:px-8">
      <AgentDownloadContent embedded />
    </div>
  </AppLayout>

  <!-- 独立全屏落地页形态: 带有专属网格背景、顶部导航条、全景内容与页脚 -->
  <div
    v-else
    class="relative min-h-screen overflow-x-hidden bg-[#f8fafc] text-slate-900 antialiased selection:bg-orange-500 selection:text-white dark:bg-[#06090e] dark:text-gray-100"
  >
    <!-- Background: 精致点状网格与微光光晕 -->
    <div class="pointer-events-none fixed inset-0 z-0 overflow-hidden">
      <!-- 顶部中央微光光晕 -->
      <div
        class="absolute -top-[250px] left-1/2 h-[500px] w-[900px] -translate-x-1/2 rounded-full bg-gradient-to-b from-orange-500/10 via-pink-500/5 to-transparent blur-[120px] dark:from-orange-600/15 dark:via-purple-600/10"
      ></div>
      <!-- Radial Dot Grid Pattern -->
      <div
        class="absolute inset-0 bg-[radial-gradient(#cbd5e1_1.2px,transparent_1.2px)] [background-size:24px_24px] dark:bg-[radial-gradient(#334155_1.2px,transparent_1.2px)]"
      ></div>
    </div>

    <!-- 顶部导航 -->
    <div class="relative z-20">
      <DownloadNavBar />
    </div>

    <!-- 主体内容 -->
    <main class="relative z-10 mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
      <AgentDownloadContent />
    </main>

    <!-- 品牌页脚 -->
    <footer class="relative z-10 border-t border-slate-200/80 dark:border-white/10 py-10 text-center text-xs text-slate-500 dark:text-dark-400 bg-white/50 dark:bg-dark-900/50 backdrop-blur-sm mt-16">
      <div class="max-w-7xl mx-auto px-4 flex flex-col sm:flex-row items-center justify-between gap-4">
        <div class="flex items-center gap-2">
          <img src="/logo.svg?v=2.1.0" alt="Logo" class="w-5 h-5 object-contain" />
          <span class="font-bold text-slate-700 dark:text-slate-300">{{ appStore.siteName }} · 魔丸 Agent 官方专属客户端</span>
        </div>
        <div>
          &copy; {{ new Date().getFullYear() }} {{ appStore.siteName }}. 别在工位当牛马，你在床上躺平，它在后台杀穿。
        </div>
      </div>
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
