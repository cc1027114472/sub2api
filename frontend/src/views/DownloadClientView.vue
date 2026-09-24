<template>
  <!-- 后台内嵌形态: ?embedded=1 且已登录，使用管理后台主布局 -->
  <AppLayout v-if="isEmbedded">
    <div class="px-4 py-6 sm:px-6 lg:px-8">
      <AgentDownloadContent embedded />
    </div>
  </AppLayout>

  <!-- 独立全屏形态: 自带导航条（Logo/站名 + 登录/返回后台） -->
  <div v-else class="min-h-screen bg-gray-50 dark:bg-dark-950">
    <DownloadNavBar />
    <main class="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
      <AgentDownloadContent />
    </main>
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
