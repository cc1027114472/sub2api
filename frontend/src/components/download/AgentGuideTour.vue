<template>
  <div class="space-y-12">
    <!-- 模块标题与导语 -->
    <div class="text-center">
      <div class="inline-flex items-center gap-1.5 rounded-full bg-primary-100/80 px-3.5 py-1 text-xs font-semibold text-primary-700 dark:bg-primary-900/40 dark:text-primary-300">
        <Icon name="sparkles" size="xs" />
        <span>{{ t('agentDownload.tour.badge', '图文全景体验') }}</span>
      </div>
      <h2 class="mt-3 text-2xl font-extrabold tracking-tight text-gray-900 dark:text-white sm:text-3xl lg:text-4xl">
        {{ t('agentDownload.tour.title') }}
      </h2>
      <p class="mx-auto mt-3 max-w-2xl text-sm text-gray-600 dark:text-dark-300 sm:text-base">
        {{ t('agentDownload.tour.subtitle') }}
      </p>
    </div>

    <!-- 图文卡片流 (方案 A: 交错排列) -->
    <div class="space-y-16 lg:space-y-24">
      <!-- Step 1: 极简沉浸的工作台 (左文右图) -->
      <div class="grid items-center gap-8 lg:grid-cols-12 lg:gap-12">
        <div class="space-y-5 lg:col-span-5">
          <div class="inline-flex items-center gap-2 rounded-xl bg-purple-50 px-3 py-1.5 text-xs font-bold text-purple-700 ring-1 ring-inset ring-purple-600/20 dark:bg-purple-950/40 dark:text-purple-300 dark:ring-purple-400/20">
            <span class="flex h-5 w-5 items-center justify-center rounded-full bg-purple-600 text-[10px] font-black text-white">01</span>
            <span>{{ t('agentDownload.tour.step1.badge') }}</span>
          </div>

          <h3 class="text-2xl font-bold tracking-tight text-gray-900 dark:text-white sm:text-3xl">
            {{ t('agentDownload.tour.step1.title') }}
          </h3>

          <div class="rounded-xl border-l-4 border-purple-500 bg-purple-50/50 p-3 text-sm italic font-medium text-purple-900 dark:bg-purple-950/20 dark:text-purple-200">
            {{ t('agentDownload.tour.step1.quote') }}
          </div>

          <p class="text-sm leading-relaxed text-gray-600 dark:text-dark-300 sm:text-base">
            {{ t('agentDownload.tour.step1.desc') }}
          </p>

          <ul class="space-y-2.5 pt-2 text-sm text-gray-700 dark:text-dark-200">
            <li class="flex items-center gap-2.5">
              <span class="flex h-5 w-5 shrink-0 items-center justify-center rounded-full bg-emerald-100 text-emerald-600 dark:bg-emerald-950/50 dark:text-emerald-400">
                <Icon name="check" size="xs" />
              </span>
              <span>{{ t('agentDownload.tour.step1.feature1') }}</span>
            </li>
            <li class="flex items-center gap-2.5">
              <span class="flex h-5 w-5 shrink-0 items-center justify-center rounded-full bg-emerald-100 text-emerald-600 dark:bg-emerald-950/50 dark:text-emerald-400">
                <Icon name="check" size="xs" />
              </span>
              <span>{{ t('agentDownload.tour.step1.feature2') }}</span>
            </li>
            <li class="flex items-center gap-2.5">
              <span class="flex h-5 w-5 shrink-0 items-center justify-center rounded-full bg-emerald-100 text-emerald-600 dark:bg-emerald-950/50 dark:text-emerald-400">
                <Icon name="check" size="xs" />
              </span>
              <span>{{ t('agentDownload.tour.step1.feature3') }}</span>
            </li>
          </ul>
        </div>

        <div class="lg:col-span-7">
          <div
            class="group relative cursor-zoom-in overflow-hidden rounded-2xl border border-gray-200/80 bg-gray-900 p-2 shadow-xl shadow-purple-500/5 transition-all duration-300 hover:shadow-2xl hover:shadow-purple-500/10 dark:border-dark-700"
            @click="openLightbox('/assets/agent/guide/guide-workbench.png', t('agentDownload.tour.step1.title'))"
          >
            <!-- 仿视窗顶栏 -->
            <div class="mb-2 flex items-center gap-1.5 px-2 pt-1">
              <div class="h-2.5 w-2.5 rounded-full bg-red-400/80"></div>
              <div class="h-2.5 w-2.5 rounded-full bg-yellow-400/80"></div>
              <div class="h-2.5 w-2.5 rounded-full bg-green-400/80"></div>
              <span class="ml-2 text-[11px] font-mono text-gray-400">mowan-harness - Workbench</span>
            </div>
            <div class="overflow-hidden rounded-xl bg-white dark:bg-dark-950">
              <img
                src="/assets/agent/guide/guide-workbench.png"
                :alt="t('agentDownload.tour.step1.title')"
                class="w-full object-cover transition-transform duration-300 group-hover:scale-[1.015]"
                loading="lazy"
              />
            </div>
            <!-- 悬浮提示 -->
            <div class="absolute bottom-4 right-4 flex items-center gap-1.5 rounded-lg bg-gray-900/80 px-2.5 py-1 text-xs font-medium text-white opacity-0 backdrop-blur-sm transition-opacity group-hover:opacity-100">
              <Icon name="search" size="xs" />
              <span>{{ t('agentDownload.tour.clickToEnlarge') }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Step 2: Sub2API 一键秒级导入 (左图右文) -->
      <div class="grid items-center gap-8 lg:grid-cols-12 lg:gap-12">
        <div class="order-2 lg:order-1 lg:col-span-7">
          <div
            class="group relative cursor-zoom-in overflow-hidden rounded-2xl border border-gray-200/80 bg-gray-900 p-2 shadow-xl shadow-primary-500/5 transition-all duration-300 hover:shadow-2xl hover:shadow-primary-500/10 dark:border-dark-700"
            @click="openLightbox('/assets/agent/guide/guide-settings.png', t('agentDownload.tour.step2.title'))"
          >
            <!-- 仿视窗顶栏 -->
            <div class="mb-2 flex items-center gap-1.5 px-2 pt-1">
              <div class="h-2.5 w-2.5 rounded-full bg-red-400/80"></div>
              <div class="h-2.5 w-2.5 rounded-full bg-yellow-400/80"></div>
              <div class="h-2.5 w-2.5 rounded-full bg-green-400/80"></div>
              <span class="ml-2 text-[11px] font-mono text-gray-400">Settings - Model Provider Auto Import</span>
            </div>
            <div class="overflow-hidden rounded-xl bg-white dark:bg-dark-950">
              <img
                src="/assets/agent/guide/guide-settings.png"
                :alt="t('agentDownload.tour.step2.title')"
                class="w-full object-cover transition-transform duration-300 group-hover:scale-[1.015]"
                loading="lazy"
              />
            </div>
            <!-- 悬浮提示 -->
            <div class="absolute bottom-4 right-4 flex items-center gap-1.5 rounded-lg bg-gray-900/80 px-2.5 py-1 text-xs font-medium text-white opacity-0 backdrop-blur-sm transition-opacity group-hover:opacity-100">
              <Icon name="search" size="xs" />
              <span>{{ t('agentDownload.tour.clickToEnlarge') }}</span>
            </div>
          </div>
        </div>

        <div class="order-1 space-y-5 lg:order-2 lg:col-span-5">
          <div class="inline-flex items-center gap-2 rounded-xl bg-blue-50 px-3 py-1.5 text-xs font-bold text-blue-700 ring-1 ring-inset ring-blue-600/20 dark:bg-blue-950/40 dark:text-blue-300 dark:ring-blue-400/20">
            <span class="flex h-5 w-5 items-center justify-center rounded-full bg-blue-600 text-[10px] font-black text-white">02</span>
            <span>{{ t('agentDownload.tour.step2.badge') }}</span>
          </div>

          <h3 class="text-2xl font-bold tracking-tight text-gray-900 dark:text-white sm:text-3xl">
            {{ t('agentDownload.tour.step2.title') }}
          </h3>

          <p class="text-sm leading-relaxed text-gray-600 dark:text-dark-300 sm:text-base">
            {{ t('agentDownload.tour.step2.desc') }}
          </p>

          <ul class="space-y-2.5 pt-2 text-sm text-gray-700 dark:text-dark-200">
            <li class="flex items-center gap-2.5">
              <span class="flex h-5 w-5 shrink-0 items-center justify-center rounded-full bg-emerald-100 text-emerald-600 dark:bg-emerald-950/50 dark:text-emerald-400">
                <Icon name="check" size="xs" />
              </span>
              <span>{{ t('agentDownload.tour.step2.feature1') }}</span>
            </li>
            <li class="flex items-center gap-2.5">
              <span class="flex h-5 w-5 shrink-0 items-center justify-center rounded-full bg-emerald-100 text-emerald-600 dark:bg-emerald-950/50 dark:text-emerald-400">
                <Icon name="check" size="xs" />
              </span>
              <span>{{ t('agentDownload.tour.step2.feature2') }}</span>
            </li>
            <li class="flex items-center gap-2.5">
              <span class="flex h-5 w-5 shrink-0 items-center justify-center rounded-full bg-emerald-100 text-emerald-600 dark:bg-emerald-950/50 dark:text-emerald-400">
                <Icon name="check" size="xs" />
              </span>
              <span>{{ t('agentDownload.tour.step2.feature3') }}</span>
            </li>
          </ul>
        </div>
      </div>

      <!-- Step 3: 强悍的代码生成与透明监控 (左文右图) -->
      <div class="grid items-center gap-8 lg:grid-cols-12 lg:gap-12">
        <div class="space-y-5 lg:col-span-5">
          <div class="inline-flex items-center gap-2 rounded-xl bg-emerald-50 px-3 py-1.5 text-xs font-bold text-emerald-700 ring-1 ring-inset ring-emerald-600/20 dark:bg-emerald-950/40 dark:text-emerald-300 dark:ring-emerald-400/20">
            <span class="flex h-5 w-5 items-center justify-center rounded-full bg-emerald-600 text-[10px] font-black text-white">03</span>
            <span>{{ t('agentDownload.tour.step3.badge') }}</span>
          </div>

          <h3 class="text-2xl font-bold tracking-tight text-gray-900 dark:text-white sm:text-3xl">
            {{ t('agentDownload.tour.step3.title') }}
          </h3>

          <p class="text-sm leading-relaxed text-gray-600 dark:text-dark-300 sm:text-base">
            {{ t('agentDownload.tour.step3.desc') }}
          </p>

          <ul class="space-y-2.5 pt-2 text-sm text-gray-700 dark:text-dark-200">
            <li class="flex items-center gap-2.5">
              <span class="flex h-5 w-5 shrink-0 items-center justify-center rounded-full bg-emerald-100 text-emerald-600 dark:bg-emerald-950/50 dark:text-emerald-400">
                <Icon name="check" size="xs" />
              </span>
              <span>{{ t('agentDownload.tour.step3.feature1') }}</span>
            </li>
            <li class="flex items-center gap-2.5">
              <span class="flex h-5 w-5 shrink-0 items-center justify-center rounded-full bg-emerald-100 text-emerald-600 dark:bg-emerald-950/50 dark:text-emerald-400">
                <Icon name="check" size="xs" />
              </span>
              <span>{{ t('agentDownload.tour.step3.feature2') }}</span>
            </li>
            <li class="flex items-center gap-2.5">
              <span class="flex h-5 w-5 shrink-0 items-center justify-center rounded-full bg-emerald-100 text-emerald-600 dark:bg-emerald-950/50 dark:text-emerald-400">
                <Icon name="check" size="xs" />
              </span>
              <span>{{ t('agentDownload.tour.step3.feature3') }}</span>
            </li>
          </ul>
        </div>

        <div class="lg:col-span-7">
          <div
            class="group relative cursor-zoom-in overflow-hidden rounded-2xl border border-gray-200/80 bg-gray-900 p-2 shadow-xl shadow-emerald-500/5 transition-all duration-300 hover:shadow-2xl hover:shadow-emerald-500/10 dark:border-dark-700"
            @click="openLightbox('/assets/agent/guide/guide-execution.png', t('agentDownload.tour.step3.title'))"
          >
            <!-- 仿视窗顶栏 -->
            <div class="mb-2 flex items-center gap-1.5 px-2 pt-1">
              <div class="h-2.5 w-2.5 rounded-full bg-red-400/80"></div>
              <div class="h-2.5 w-2.5 rounded-full bg-yellow-400/80"></div>
              <div class="h-2.5 w-2.5 rounded-full bg-green-400/80"></div>
              <span class="ml-2 text-[11px] font-mono text-gray-400">Execution & Metrics Monitoring</span>
            </div>
            <div class="overflow-hidden rounded-xl bg-white dark:bg-dark-950">
              <img
                src="/assets/agent/guide/guide-execution.png"
                :alt="t('agentDownload.tour.step3.title')"
                class="w-full object-cover transition-transform duration-300 group-hover:scale-[1.015]"
                loading="lazy"
              />
            </div>
            <!-- 悬浮提示 -->
            <div class="absolute bottom-4 right-4 flex items-center gap-1.5 rounded-lg bg-gray-900/80 px-2.5 py-1 text-xs font-medium text-white opacity-0 backdrop-blur-sm transition-opacity group-hover:opacity-100">
              <Icon name="search" size="xs" />
              <span>{{ t('agentDownload.tour.clickToEnlarge') }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 高清大图预览弹窗 (Lightbox) -->
    <Teleport to="body">
      <div
        v-if="lightboxOpen"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/85 p-4 backdrop-blur-sm transition-opacity animate-in fade-in duration-200"
        @click="closeLightbox"
      >
        <div class="relative max-h-[92vh] max-w-[94vw] overflow-hidden rounded-2xl bg-gray-950 p-2 shadow-2xl" @click.stop>
          <div class="mb-2 flex items-center justify-between border-b border-gray-800 px-3 py-1.5 text-xs text-gray-300">
            <span class="font-medium truncate max-w-[80vw]">{{ activeTitle }}</span>
            <button
              type="button"
              class="rounded-lg p-1 text-gray-400 transition-colors hover:bg-gray-800 hover:text-white"
              @click="closeLightbox"
            >
              <Icon name="x" size="sm" />
            </button>
          </div>
          <div class="max-h-[84vh] overflow-auto rounded-xl">
            <img :src="activeImg" :alt="activeTitle" class="h-auto max-w-full rounded-lg object-contain shadow-lg" />
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import Icon from '@/components/icons/Icon.vue'

const { t } = useI18n()

const lightboxOpen = ref(false)
const activeImg = ref('')
const activeTitle = ref('')

const openLightbox = (src: string, title: string) => {
  activeImg.value = src
  activeTitle.value = title
  lightboxOpen.value = true
  document.body.style.overflow = 'hidden'
}

const closeLightbox = () => {
  lightboxOpen.value = false
  document.body.style.overflow = ''
}

const handleKeyDown = (e: KeyboardEvent) => {
  if (e.key === 'Escape' && lightboxOpen.value) {
    closeLightbox()
  }
}

onMounted(() => {
  window.addEventListener('keydown', handleKeyDown)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeyDown)
  document.body.style.overflow = ''
})
</script>
