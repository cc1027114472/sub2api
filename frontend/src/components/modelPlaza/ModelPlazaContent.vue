<template>
  <div class="space-y-6">
    <!-- 页头(独立形态下展示标题;后台形态 AppHeader 已有页面标题) -->
    <div v-if="!embedded" class="relative">
      <div class="inline-flex items-center gap-2 rounded-full border border-primary-500/20 bg-primary-500/10 px-3 py-1 text-xs font-semibold text-primary-600 dark:text-primary-300">
        <span class="flex h-1.5 w-1.5 rounded-full bg-primary-500 animate-pulse"></span>
        MODEL CATALOG & REAL-TIME PRICING
      </div>
      <h1 class="mt-3 text-3xl font-extrabold tracking-tight text-gray-900 sm:text-4xl dark:text-white">
        {{ t('modelPlaza.title') }}
      </h1>
      <p class="mt-2 max-w-3xl text-sm leading-relaxed text-gray-600 dark:text-dark-300 sm:text-base">
        {{ t('modelPlaza.description') }}
      </p>
    </div>

    <!-- 全局价格说明(管理员配置,Markdown) - Bento 细微光边框 -->
    <div
      v-if="descriptionHtml"
      class="plaza-description rounded-2xl border border-gray-200/80 bg-white/70 p-5 text-sm shadow-sm backdrop-blur-sm dark:border-white/10 dark:bg-dark-900/60 dark:text-dark-200"
      v-html="descriptionHtml"
    ></div>

    <!-- 未登录提示 (胶囊状) -->
    <div
      v-if="!isAuthenticated"
      class="flex items-center gap-2 rounded-xl border border-amber-500/20 bg-amber-50/50 px-4 py-2.5 text-xs text-amber-700 dark:border-amber-500/10 dark:bg-amber-950/20 dark:text-amber-300"
    >
      <Icon name="infoCircle" size="xs" class="h-4 w-4 shrink-0 text-amber-500" />
      <span>{{ t('modelPlaza.anonymousHint') }}</span>
    </div>

    <!-- 加载中 -->
    <div v-if="loading" class="flex min-h-[300px] flex-col items-center justify-center gap-3">
      <div class="h-9 w-9 animate-spin rounded-full border-2 border-primary-500/20 border-t-primary-600 dark:border-primary-400/20 dark:border-t-primary-400"></div>
      <span class="text-xs text-gray-400 dark:text-dark-400">加载全模型价格矩阵...</span>
    </div>

    <!-- 错误状态 -->
    <div
      v-else-if="error"
      class="rounded-2xl border border-red-200/80 bg-red-50/60 p-8 text-center text-sm text-red-600 dark:border-red-500/20 dark:bg-red-950/20 dark:text-red-400"
    >
      <Icon name="exclamationTriangle" size="md" class="mx-auto mb-2 text-red-500" />
      {{ t('modelPlaza.loadFailed') }}
    </div>

    <!-- 筛选与分组列表 -->
    <template v-else>
      <!-- 筛选区: 胶囊卡片容器包裹 -->
      <div class="rounded-2xl border border-gray-200/80 bg-white/70 p-5 shadow-sm backdrop-blur-sm dark:border-white/10 dark:bg-dark-900/60">
        <PlazaFilterBar
          :platforms="platforms"
          :groups="groupOptions"
          :rates="rates"
          :platform="selectedPlatform"
          :group-id="selectedGroupId"
          :rate="selectedRate"
          :search="searchQuery"
          @update:platform="selectedPlatform = $event"
          @update:group-id="selectedGroupId = $event"
          @update:rate="selectedRate = $event"
          @update:search="searchQuery = $event"
        />
      </div>

      <!-- 分组分节的模型清单 -->
      <div v-if="filteredGroups.length > 0" class="space-y-6">
        <PlazaGroupSection v-for="g in filteredGroups" :key="g.id" :group="g" />
      </div>

      <!-- 搜索空状态 -->
      <div
        v-else
        class="rounded-2xl border border-dashed border-gray-300/80 p-12 text-center text-sm text-gray-500 dark:border-white/10 dark:text-dark-400"
      >
        <Icon name="search" size="md" class="mx-auto mb-2 text-gray-400 opacity-60" />
        {{ searchActive ? t('modelPlaza.noSearchResult') : t('modelPlaza.empty') }}
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { marked } from 'marked'
import DOMPurify from 'dompurify'
import Icon from '@/components/icons/Icon.vue'
import PlazaFilterBar from './PlazaFilterBar.vue'
import PlazaGroupSection from './PlazaGroupSection.vue'
import type { ModelPlazaGroup, ModelPlazaResponse } from '@/api/modelPlaza'
import { useAuthStore } from '@/stores/auth'

const props = defineProps<{
  response: ModelPlazaResponse | null
  loading: boolean
  error?: boolean
  /** 后台内嵌形态(AppLayout 内):隐藏页头。 */
  embedded?: boolean
}>()

const { t } = useI18n()
const authStore = useAuthStore()
const isAuthenticated = computed(() => authStore.isAuthenticated)

const selectedPlatform = ref<string>('all')
const selectedGroupId = ref<number | 'all'>('all')
const selectedRate = ref<number | 'all'>('all')
const searchQuery = ref('')

const searchActive = computed(() => searchQuery.value.trim() !== '')

const descriptionHtml = computed(() => {
  const md = props.response?.description?.trim()
  if (!md) return ''
  return DOMPurify.sanitize(marked.parse(md) as string)
})

/** 生效倍率 = 用户专属倍率 ?? 分组默认倍率。 */
function effectiveRate(g: ModelPlazaGroup): number {
  return g.user_rate_multiplier ?? g.rate_multiplier
}

const platforms = computed(() =>
  [...new Set((props.response?.groups ?? []).map((g) => g.platform).filter(Boolean))].sort()
)

const groupOptions = computed(() =>
  (props.response?.groups ?? []).map((g) => ({
    id: g.id,
    name: g.name,
    platform: g.platform,
    rate: effectiveRate(g)
  }))
)

/** 全量生效倍率;当前组合下不可用的项由 FilterBar 置灰而非隐藏。 */
const rates = computed(() =>
  [...new Set((props.response?.groups ?? []).map((g) => effectiveRate(g)))].sort(
    (a, b) => a - b
  )
)

/** 重置级联:切平台时,如果当前选中的分组不属于新平台,重置 groupId 为 all。 */
watch(selectedPlatform, (newP) => {
  if (newP === 'all' || selectedGroupId.value === 'all') return
  const g = props.response?.groups.find((x) => x.id === selectedGroupId.value)
  if (g && g.platform !== newP) {
    selectedGroupId.value = 'all'
  }
})

/** 过滤出要展示的分组及分组内部过滤后的模型。 */
const filteredGroups = computed<ModelPlazaGroup[]>(() => {
  if (!props.response) return []
  const q = searchQuery.value.trim().toLowerCase()

  return props.response.groups
    .filter((g) => {
      if (selectedPlatform.value !== 'all' && g.platform !== selectedPlatform.value) {
        return false
      }
      if (selectedGroupId.value !== 'all' && g.id !== selectedGroupId.value) {
        return false
      }
      if (selectedRate.value !== 'all' && effectiveRate(g) !== selectedRate.value) {
        return false
      }
      return true
    })
    .map((g) => {
      if (!q) return g
      const models = g.models.filter((m) => m.name.toLowerCase().includes(q))
      return { ...g, models }
    })
    .filter((g) => !q || g.models.length > 0)
    .sort((a, b) => effectiveRate(a) - effectiveRate(b))
})
</script>

<style scoped>
:deep(.plaza-description) h1,
:deep(.plaza-description) h2,
:deep(.plaza-description) h3 {
  font-weight: 600;
  margin-top: 0.5rem;
  margin-bottom: 0.25rem;
}
:deep(.plaza-description) p {
  margin-bottom: 0.5rem;
}
:deep(.plaza-description) p:last-child {
  margin-bottom: 0;
}
:deep(.plaza-description) ul {
  list-style: disc;
  padding-left: 1.25rem;
  margin-bottom: 0.5rem;
}
:deep(.plaza-description) a {
  color: #6366f1;
  text-decoration: underline;
}
</style>
