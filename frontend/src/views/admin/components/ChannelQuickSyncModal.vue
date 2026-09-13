<template>
  <BaseDialog
    :show="visible"
    :title="dialogTitle"
    width="extra-wide"
    @close="handleClose"
  >
    <!-- Wizard Step Indicator -->
    <div class="mb-6">
      <div class="flex items-center justify-between">
        <div
          v-for="s in steps"
          :key="s.step"
          class="flex flex-1 items-center"
          :class="{ 'opacity-50': currentStep < s.step }"
        >
          <div class="flex items-center gap-2">
            <div
              :class="[
                'flex h-7 w-7 items-center justify-center rounded-full text-xs font-semibold transition-colors',
                currentStep === s.step
                  ? 'bg-primary-600 text-white'
                  : currentStep > s.step
                  ? 'bg-green-500 text-white'
                  : 'bg-gray-200 text-gray-700 dark:bg-dark-600 dark:text-gray-300'
              ]"
            >
              <Icon v-if="currentStep > s.step" name="check" size="xs" />
              <span v-else>{{ s.step }}</span>
            </div>
            <span
              :class="[
                'text-sm font-medium whitespace-nowrap',
                currentStep === s.step
                  ? 'text-primary-600 dark:text-primary-400 font-semibold'
                  : 'text-gray-700 dark:text-gray-300'
              ]"
            >
              {{ s.title }}
            </span>
          </div>
          <div
            v-if="s.step < steps.length"
            class="mx-4 h-0.5 flex-1 bg-gray-200 dark:bg-dark-600"
          />
        </div>
      </div>
    </div>

    <!-- Step 1: 连接与探测 (Connection & Probe) -->
    <div v-if="currentStep === 1" class="space-y-4">
      <div>
        <label class="input-label">{{ t('admin.channels.quickSync.channelName', '节点名称') }} <span class="text-red-500">*</span></label>
        <input
          v-model="form.name"
          type="text"
          class="input"
          data-test="channel-name-input"
          :placeholder="t('admin.channels.quickSync.channelNamePlaceholder', '例如：Antigravity-Node-1')"
        />
      </div>

      <div>
        <label class="input-label">{{ t('admin.channels.quickSync.baseUrl', '基础接口地址 (Base URL)') }} <span class="text-red-500">*</span></label>
        <input
          v-model="form.base_url"
          type="text"
          class="input font-mono text-sm"
          data-test="base-url-input"
          :placeholder="t('admin.channels.quickSync.baseUrlPlaceholder', 'http://127.0.0.1:8080 或上游 API 接口 URL')"
        />
      </div>

      <div>
        <label class="input-label">{{ t('admin.channels.quickSync.apiKey', 'API 密钥 (API Key)') }} <span class="text-red-500">*</span></label>
        <div class="relative">
          <input
            v-model="form.api_key"
            :type="showApiKey ? 'text' : 'password'"
            class="input pr-10 font-mono text-sm"
            data-test="api-key-input"
            :placeholder="t('admin.channels.quickSync.apiKeyPlaceholder', 'sk-... 或上游访问令牌')"
          />
          <button
            type="button"
            class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
            @click="showApiKey = !showApiKey"
          >
            <Icon :name="showApiKey ? 'eyeOff' : 'eye'" size="sm" />
          </button>
        </div>
      </div>

      <div>
        <label class="input-label">{{ t('admin.channels.quickSync.platform', '平台类型') }}</label>
        <select v-model="form.platform" class="input" data-test="platform-select">
          <option value="antigravity">Antigravity (多协议自动桥接)</option>
          <option value="gemini">Gemini</option>
          <option value="openai">OpenAI</option>
        </select>
        <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
          {{ t('admin.channels.quickSync.platformHelp', '选择 antigravity 将自动适配 Gemini、OpenAI、Claude 与 Codex 协议。') }}
        </p>
      </div>

      <div v-if="probeError" class="rounded-lg bg-red-50 p-3 text-sm text-red-600 dark:bg-red-900/20 dark:text-red-400">
        {{ probeError }}
      </div>
    </div>

    <!-- Step 2: 分组与计费配置 (Group & Pricing) -->
    <div v-else-if="currentStep === 2" class="space-y-6">
      <!-- 1. 分组配置卡片 -->
      <div class="rounded-lg border border-gray-200 bg-gray-50/50 p-4 dark:border-dark-700 dark:bg-dark-800/50">
        <h4 class="mb-3 text-sm font-semibold text-gray-900 dark:text-white flex items-center gap-2">
          <Icon name="users" size="sm" class="text-primary-600 dark:text-primary-400" />
          {{ t('admin.channels.quickSync.groupConfig', '目标分组配置') }}
        </h4>
        <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
          <div>
            <label class="input-label">{{ t('admin.channels.quickSync.defaultGroup', '默认归属分组') }}</label>
            <select
              v-model.number="form.default_group_id"
              class="input"
              data-test="default-group-select"
              @change="onDefaultGroupChanged"
            >
              <option v-if="form.createNewGroup" :value="0">
                ⭐ [新建] {{ form.newGroupName || t('admin.channels.quickSync.newGroupDefault', '新建独立分组') }}
              </option>
              <option v-for="g in existingGroups" :key="g.id" :value="g.id">
                {{ g.name }} (ID: {{ g.id }})
              </option>
            </select>
          </div>

          <div class="flex items-center">
            <label class="flex items-center gap-2 cursor-pointer pt-6">
              <input
                v-model="form.createNewGroup"
                type="checkbox"
                class="rounded border-gray-300 text-primary-600 focus:ring-primary-500 dark:border-dark-600 dark:bg-dark-700"
                data-test="create-new-group-checkbox"
                @change="onToggleCreateNewGroup"
              />
              <span class="text-sm font-medium text-gray-700 dark:text-gray-300">
                + {{ t('admin.channels.quickSync.createNewGroup', '新建独立分组') }}
              </span>
            </label>
          </div>
        </div>

        <!-- 新建分组展开表单 -->
        <div v-if="form.createNewGroup" class="mt-4 grid grid-cols-1 gap-4 border-t border-gray-200 pt-4 dark:border-dark-700 sm:grid-cols-2">
          <div>
            <label class="input-label">{{ t('admin.channels.quickSync.newGroupName', '新分组名称') }} <span class="text-red-500">*</span></label>
            <input
              v-model="form.newGroupName"
              type="text"
              class="input"
              data-test="new-group-name-input"
              :placeholder="t('admin.channels.quickSync.newGroupNamePlaceholder', '输入新分组名称')"
            />
          </div>
          <div>
            <label class="input-label">{{ t('admin.channels.quickSync.newGroupMultiplier', '新分组费率倍率') }}</label>
            <input
              v-model.number="form.newGroupRateMultiplier"
              type="number"
              step="0.1"
              min="0.01"
              class="input"
              data-test="new-group-multiplier-input"
            />
          </div>
        </div>
      </div>

      <!-- 2. 计费策略配置卡片 -->
      <div class="rounded-lg border border-gray-200 bg-gray-50/50 p-4 dark:border-dark-700 dark:bg-dark-800/50">
        <h4 class="mb-3 text-sm font-semibold text-gray-900 dark:text-white flex items-center gap-2">
          <Icon name="chart" size="sm" class="text-primary-600 dark:text-primary-400" />
          {{ t('admin.channels.quickSync.billingStrategy', '计费策略设定') }}
        </h4>
        <div class="flex flex-wrap gap-4 items-center">
          <div class="flex rounded-lg bg-gray-200 p-0.5 dark:bg-dark-700">
            <button
              type="button"
              :class="[
                'px-3 py-1.5 text-xs font-medium rounded-md transition-colors',
                billingStrategy.mode === 'ratio'
                  ? 'bg-white text-primary-600 shadow-sm dark:bg-dark-600 dark:text-primary-400'
                  : 'text-gray-600 hover:text-gray-900 dark:text-gray-400 dark:hover:text-white'
              ]"
              data-test="strategy-ratio-btn"
              @click="setBillingStrategyMode('ratio')"
            >
              {{ t('admin.channels.quickSync.ratioMode', '倍率计费模式') }}
            </button>
            <button
              type="button"
              :class="[
                'px-3 py-1.5 text-xs font-medium rounded-md transition-colors',
                billingStrategy.mode === 'per_request'
                  ? 'bg-white text-primary-600 shadow-sm dark:bg-dark-600 dark:text-primary-400'
                  : 'text-gray-600 hover:text-gray-900 dark:text-gray-400 dark:hover:text-white'
              ]"
              data-test="strategy-per-request-btn"
              @click="setBillingStrategyMode('per_request')"
            >
              {{ t('admin.channels.quickSync.perRequestMode', '按次计费模式') }}
            </button>
            <button
              type="button"
              :class="[
                'px-3 py-1.5 text-xs font-medium rounded-md transition-colors',
                billingStrategy.mode === 'fixed'
                  ? 'bg-white text-primary-600 shadow-sm dark:bg-dark-600 dark:text-primary-400'
                  : 'text-gray-600 hover:text-gray-900 dark:text-gray-400 dark:hover:text-white'
              ]"
              data-test="strategy-fixed-btn"
              @click="setBillingStrategyMode('fixed')"
            >
              {{ t('admin.channels.quickSync.customMode', '自定义调整') }}
            </button>
          </div>

          <!-- 倍率输入 -->
          <div v-if="billingStrategy.mode === 'ratio'" class="flex items-center gap-2">
            <span class="text-xs text-gray-500 dark:text-gray-400">{{ t('admin.channels.quickSync.multiplierLabel', '基准倍率:') }}</span>
            <input
              v-model.number="billingStrategy.ratio"
              type="number"
              step="0.1"
              min="0.01"
              class="input w-24 py-1 text-sm"
              data-test="strategy-ratio-input"
              @input="applyBillingStrategy"
            />
            <span class="text-xs text-gray-400">({{ t('admin.channels.quickSync.ratioHint', '自动按倍率更新各模型输入/输出价格') }})</span>
          </div>

          <!-- 按次价格输入 -->
          <div v-if="billingStrategy.mode === 'per_request'" class="flex items-center gap-2">
            <span class="text-xs text-gray-500 dark:text-gray-400">{{ t('admin.channels.quickSync.perReqPriceLabel', '单次费用 ($):') }}</span>
            <input
              v-model.number="billingStrategy.per_request_price"
              type="number"
              step="0.001"
              min="0"
              class="input w-28 py-1 text-sm"
              data-test="strategy-per-request-input"
              @input="applyBillingStrategy"
            />
            <span class="text-xs text-gray-400">({{ t('admin.channels.quickSync.perReqHint', '将所有模型设置为按次计费') }})</span>
          </div>
        </div>
      </div>

      <!-- 3. 渠道健康监控探针配置卡片 -->
      <div class="rounded-lg border border-gray-200 bg-gray-50/50 p-4 dark:border-dark-700 dark:bg-dark-800/50">
        <div class="flex items-center justify-between">
          <h4 class="text-sm font-semibold text-gray-900 dark:text-white flex items-center gap-2">
            <Icon name="bolt" size="sm" class="text-primary-600 dark:text-primary-400" />
            {{ t('admin.channels.quickSync.monitorTitle', '渠道健康监控') }}
          </h4>
          <label class="inline-flex items-center cursor-pointer">
            <input
              v-model="form.enableMonitor"
              type="checkbox"
              class="sr-only peer"
              data-test="enable-monitor-toggle"
            />
            <div class="relative w-9 h-5 bg-gray-200 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-4 after:w-4 after:transition-all dark:border-gray-600 peer-checked:bg-primary-600 dark:bg-dark-700"></div>
            <span class="ms-2 text-xs font-medium text-gray-700 dark:text-gray-300">
              {{ form.enableMonitor ? t('common.enabled', '已启用') : t('common.disabled', '未启用') }}
            </span>
          </label>
        </div>
        <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
          {{ t('admin.channels.quickSync.monitorDesc', '自动创建健康巡检探针，定时拨测节点延迟与可用率') }}
        </p>

        <div v-if="form.enableMonitor" class="mt-3 grid grid-cols-1 md:grid-cols-2 gap-4 border-t border-gray-200/60 pt-3 dark:border-dark-700/60">
          <div>
            <label class="input-label">{{ t('admin.channels.quickSync.monitorPrimaryModel', '主探测模型') }}</label>
            <select
              v-model="form.monitorModel"
              class="input w-full"
              data-test="monitor-model-select"
            >
              <option v-for="m in models" :key="m.id" :value="m.id">
                {{ m.id }}
              </option>
            </select>
          </div>
          <div>
            <label class="input-label">{{ t('admin.channels.quickSync.monitorInterval', '巡检周期 (秒)') }}</label>
            <input
              v-model.number="form.monitorInterval"
              type="number"
              min="10"
              step="5"
              class="input w-full"
              data-test="monitor-interval-input"
              placeholder="60"
            />
          </div>
        </div>
      </div>

      <!-- 4. 模型调配表格与批量操作 -->
      <div>
        <div class="mb-3 flex flex-wrap items-center justify-between gap-2">
          <div class="flex items-center gap-2">
            <h4 class="text-sm font-semibold text-gray-900 dark:text-white">
              {{ t('admin.channels.quickSync.modelListTitle', '探测模型列表') }} ({{ models.length }})
            </h4>
            <span class="text-xs text-gray-500 dark:text-gray-400">
              已选 {{ selectedModelCount }} 个
            </span>
          </div>

          <!-- 批量操作工具栏 -->
          <div v-if="selectedModelCount > 0" class="flex flex-wrap items-center gap-2 rounded-md bg-primary-50 px-3 py-1.5 dark:bg-primary-900/20">
            <span class="text-xs font-medium text-primary-700 dark:text-primary-300">{{ t('admin.channels.quickSync.batchActions', '批量操作:') }}</span>
            
            <select
              v-model.number="batchTargetGroupId"
              class="input py-0.5 text-xs w-36"
              data-test="batch-group-select"
            >
              <option :value="null" disabled>{{ t('admin.channels.quickSync.selectGroup', '选择目标分组') }}</option>
              <option v-if="form.createNewGroup" :value="0">⭐ [新建] {{ form.newGroupName || '新分组' }}</option>
              <option v-for="g in existingGroups" :key="g.id" :value="g.id">{{ g.name }}</option>
            </select>
            <button
              type="button"
              class="btn btn-secondary py-0.5 px-2 text-xs"
              data-test="apply-batch-group-btn"
              :disabled="batchTargetGroupId === null"
              @click="applyBatchGroup"
            >
              {{ t('admin.channels.quickSync.applyGroup', '移动到分组') }}
            </button>

            <div class="h-4 w-px bg-gray-300 dark:bg-dark-600 mx-1" />

            <button
              type="button"
              class="btn btn-secondary py-0.5 px-2 text-xs"
              data-test="batch-mode-token-btn"
              @click="applyBatchBillingMode('token')"
            >
              {{ t('admin.channels.quickSync.setTokenMode', '设为按Token') }}
            </button>
            <button
              type="button"
              class="btn btn-secondary py-0.5 px-2 text-xs"
              data-test="batch-mode-per-req-btn"
              @click="applyBatchBillingMode('per_request')"
            >
              {{ t('admin.channels.quickSync.setPerReqMode', '设为按次') }}
            </button>
          </div>
        </div>

        <!-- 表格 -->
        <div class="max-h-80 overflow-y-auto rounded-lg border border-gray-200 dark:border-dark-700">
          <table class="w-full text-left text-xs">
            <thead class="sticky top-0 bg-gray-100 text-gray-700 dark:bg-dark-700 dark:text-gray-300">
              <tr>
                <th class="w-10 px-3 py-2 text-center">
                  <input
                    type="checkbox"
                    :checked="allSelected"
                    class="rounded border-gray-300 text-primary-600 focus:ring-primary-500 dark:border-dark-600 dark:bg-dark-700"
                    data-test="select-all-checkbox"
                    @change="toggleSelectAll"
                  />
                </th>
                <th class="px-3 py-2">{{ t('admin.channels.quickSync.modelName', '模型 ID') }}</th>
                <th class="px-3 py-2 w-44">{{ t('admin.channels.quickSync.targetGroup', '归属分组') }}</th>
                <th class="px-3 py-2 w-28">{{ t('admin.channels.quickSync.billingMode', '计费模式') }}</th>
                <th class="px-3 py-2 w-32">{{ t('admin.channels.quickSync.inputPrice', '输入价 ($/1M)') }}</th>
                <th class="px-3 py-2 w-32">{{ t('admin.channels.quickSync.outputPrice', '输出价 ($/1M)') }}</th>
                <th class="px-3 py-2 w-32">{{ t('admin.channels.quickSync.perReqPrice', '按次单价 ($)') }}</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-200 dark:divide-dark-700">
              <tr
                v-for="item in models"
                :key="item.id"
                :class="[
                  'transition-colors hover:bg-gray-50 dark:hover:bg-dark-800',
                  item.selected ? 'bg-primary-50/40 dark:bg-primary-900/10' : ''
                ]"
                :data-test="`model-row-${item.id}`"
              >
                <td class="px-3 py-2 text-center">
                  <input
                    v-model="item.selected"
                    type="checkbox"
                    class="rounded border-gray-300 text-primary-600 focus:ring-primary-500 dark:border-dark-600 dark:bg-dark-700"
                    :data-test="`model-select-${item.id}`"
                  />
                </td>
                <td class="px-3 py-2 font-medium text-gray-900 dark:text-white">
                  <div class="flex items-center gap-1.5">
                    <span class="font-mono">{{ item.id }}</span>
                    <span
                      v-if="item.display_name && item.display_name !== item.id"
                      class="text-gray-400 text-[10px]"
                    >
                      ({{ item.display_name }})
                    </span>
                  </div>
                </td>
                <td class="px-3 py-2">
                  <select
                    v-model.number="item.target_group_id"
                    class="input py-1 text-xs"
                    :data-test="`model-group-${item.id}`"
                  >
                    <option v-if="form.createNewGroup" :value="0">⭐ [新建] {{ form.newGroupName || '新分组' }}</option>
                    <option v-for="g in existingGroups" :key="g.id" :value="g.id">{{ g.name }}</option>
                  </select>
                </td>
                <td class="px-3 py-2">
                  <select
                    v-model="item.billing_mode"
                    class="input py-1 text-xs"
                    :data-test="`model-mode-${item.id}`"
                  >
                    <option value="token">{{ t('admin.availableChannels.pricing.billingModeToken', '按 Token') }}</option>
                    <option value="per_request">{{ t('admin.availableChannels.pricing.billingModePerRequest', '按次') }}</option>
                  </select>
                </td>
                <td class="px-3 py-2">
                  <div v-if="item.billing_mode === 'token'">
                    <input
                      v-model.number="item.price_in"
                      type="number"
                      step="0.000001"
                      class="input py-1 text-xs font-mono"
                      :data-test="`model-price-in-${item.id}`"
                    />
                  </div>
                  <span v-else class="text-gray-400">-</span>
                </td>
                <td class="px-3 py-2">
                  <div v-if="item.billing_mode === 'token'">
                    <input
                      v-model.number="item.price_out"
                      type="number"
                      step="0.000001"
                      class="input py-1 text-xs font-mono"
                      :data-test="`model-price-out-${item.id}`"
                    />
                  </div>
                  <span v-else class="text-gray-400">-</span>
                </td>
                <td class="px-3 py-2">
                  <div v-if="item.billing_mode === 'per_request'">
                    <input
                      v-model.number="item.per_request_price"
                      type="number"
                      step="0.001"
                      class="input py-1 text-xs font-mono"
                      :data-test="`model-price-per-req-${item.id}`"
                    />
                  </div>
                  <span v-else class="text-gray-400">-</span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- Step 3: 确认与导入 (Summary & Commit) -->
    <div v-else-if="currentStep === 3" class="space-y-6">
      <div class="grid grid-cols-1 gap-4 md:grid-cols-2">
        <!-- 节点信息卡片 -->
        <div class="rounded-lg border border-gray-200 bg-gray-50/50 p-4 dark:border-dark-700 dark:bg-dark-800/50">
          <h4 class="mb-3 text-sm font-semibold text-gray-900 dark:text-white flex items-center gap-2">
            <Icon name="link" size="sm" class="text-primary-600 dark:text-primary-400" />
            {{ t('admin.channels.quickSync.summaryChannel', '节点与凭据配置') }}
          </h4>
          <dl class="space-y-2 text-xs">
            <div class="flex justify-between">
              <dt class="text-gray-500 dark:text-gray-400">{{ t('admin.channels.quickSync.channelName', '节点名称') }}:</dt>
              <dd class="font-medium text-gray-900 dark:text-white">{{ form.name }}</dd>
            </div>
            <div class="flex justify-between">
              <dt class="text-gray-500 dark:text-gray-400">{{ t('admin.channels.quickSync.platform', '平台类型') }}:</dt>
              <dd class="font-medium text-gray-900 dark:text-white uppercase">{{ form.platform }}</dd>
            </div>
            <div class="flex justify-between">
              <dt class="text-gray-500 dark:text-gray-400">{{ t('admin.channels.quickSync.baseUrl', 'Base URL') }}:</dt>
              <dd class="font-mono text-gray-900 dark:text-white truncate max-w-xs" :title="form.base_url">{{ form.base_url }}</dd>
            </div>
          </dl>
        </div>

        <!-- 分组与策略卡片 -->
        <div class="rounded-lg border border-gray-200 bg-gray-50/50 p-4 dark:border-dark-700 dark:bg-dark-800/50">
          <h4 class="mb-3 text-sm font-semibold text-gray-900 dark:text-white flex items-center gap-2">
            <Icon name="users" size="sm" class="text-primary-600 dark:text-primary-400" />
            {{ t('admin.channels.quickSync.summaryGroups', '目标分组与模型统计') }}
          </h4>
          <dl class="space-y-2 text-xs">
            <div class="flex justify-between">
              <dt class="text-gray-500 dark:text-gray-400">{{ t('admin.channels.quickSync.newGroupAction', '新建分组') }}:</dt>
              <dd class="font-medium text-gray-900 dark:text-white">
                <span v-if="form.createNewGroup" class="text-green-600 dark:text-green-400">
                  {{ form.newGroupName }} (倍率: {{ form.newGroupRateMultiplier }}x)
                </span>
                <span v-else class="text-gray-400">不新建</span>
              </dd>
            </div>
            <div class="flex justify-between">
              <dt class="text-gray-500 dark:text-gray-400">{{ t('admin.channels.quickSync.totalModelsCount', '模型总数') }}:</dt>
              <dd class="font-bold text-gray-900 dark:text-white">{{ models.length }}</dd>
            </div>
            <div class="flex justify-between">
              <dt class="text-gray-500 dark:text-gray-400">{{ t('admin.channels.quickSync.billingBreakdown', '计费模式分布') }}:</dt>
              <dd class="text-gray-900 dark:text-white">
                Token: {{ tokenModelCount }} 个 / 按次: {{ perReqModelCount }} 个
              </dd>
            </div>
            <div class="flex justify-between">
              <dt class="text-gray-500 dark:text-gray-400">{{ t('admin.channels.quickSync.monitorStatus', '健康监控') }}:</dt>
              <dd class="font-medium text-gray-900 dark:text-white">
                <span v-if="form.enableMonitor" class="text-green-600 dark:text-green-400">
                  {{ form.monitorModel || '默认主模型' }} ({{ form.monitorInterval || 60 }}s 巡检)
                </span>
                <span v-else class="text-gray-400">未启用</span>
              </dd>
            </div>
          </dl>
        </div>
      </div>

      <!-- 分组归属分布明细 -->
      <div>
        <h4 class="mb-2 text-xs font-semibold text-gray-700 dark:text-gray-300">
          {{ t('admin.channels.quickSync.groupBreakdown', '各分组包含模型明细') }}:
        </h4>
        <div class="space-y-2 max-h-48 overflow-y-auto">
          <div
            v-for="groupSummary in groupSummaries"
            :key="groupSummary.name"
            class="rounded border border-gray-200 bg-white p-2.5 text-xs dark:border-dark-700 dark:bg-dark-800"
          >
            <div class="flex items-center justify-between font-medium text-gray-800 dark:text-gray-200">
              <span>{{ groupSummary.name }}</span>
              <span class="rounded bg-primary-100 px-2 py-0.5 text-primary-700 dark:bg-primary-900/30 dark:text-primary-300">
                {{ groupSummary.models.length }} 个模型
              </span>
            </div>
            <div class="mt-1 flex flex-wrap gap-1">
              <span
                v-for="m in groupSummary.models"
                :key="m"
                class="rounded bg-gray-100 px-1.5 py-0.5 text-[11px] font-mono text-gray-600 dark:bg-dark-700 dark:text-gray-400"
              >
                {{ m }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Footer Actions -->
    <template #footer>
      <div class="flex items-center justify-between w-full">
        <div>
          <button
            v-if="currentStep > 1"
            type="button"
            class="btn btn-secondary"
            data-test="prev-step-btn"
            :disabled="loading"
            @click="prevStep"
          >
            <Icon name="chevronLeft" size="sm" class="mr-1" />
            {{ t('common.back', '上一步') }}
          </button>
        </div>

        <div class="flex items-center gap-2">
          <button
            type="button"
            class="btn btn-secondary"
            data-test="cancel-btn"
            :disabled="loading"
            @click="handleClose"
          >
            {{ t('common.cancel', '取消') }}
          </button>

          <!-- Step 1 Next Button: Probe -->
          <button
            v-if="currentStep === 1"
            type="button"
            class="btn btn-primary"
            data-test="probe-btn"
            :disabled="loading || !isStep1Valid"
            @click="handleProbe"
          >
            <Icon v-if="loading" name="refresh" size="sm" class="mr-1 animate-spin" />
            <Icon v-else name="play" size="sm" class="mr-1" />
            {{ loading ? t('admin.channels.quickSync.probing', '正在探测...') : t('admin.channels.quickSync.probeButton', '测试连接并探测模型') }}
          </button>

          <!-- Step 2 Next Button -->
          <button
            v-else-if="currentStep === 2"
            type="button"
            class="btn btn-primary"
            data-test="step2-next-btn"
            :disabled="!isStep2Valid"
            @click="nextStep"
          >
            {{ t('admin.channels.quickSync.nextConfirm', '下一步：确认导入') }}
            <Icon name="chevronRight" size="sm" class="ml-1" />
          </button>

          <!-- Step 3 Commit Button -->
          <button
            v-else-if="currentStep === 3"
            type="button"
            class="btn btn-primary"
            data-test="commit-btn"
            :disabled="loading"
            @click="handleCommit"
          >
            <Icon v-if="loading" name="refresh" size="sm" class="mr-1 animate-spin" />
            <Icon v-else name="check" size="sm" class="mr-1" />
            {{ loading ? t('admin.channels.quickSync.committing', '正在导入...') : t('admin.channels.quickSync.commitButton', '确认导入') }}
          </button>
        </div>
      </div>
    </template>
  </BaseDialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { adminAPI } from '@/api/admin'
import { useAppStore } from '@/stores/app'
import BaseDialog from '@/components/common/BaseDialog.vue'
import Icon from '@/components/icons/Icon.vue'
import type { AdminGroup } from '@/types'
import type {
  QuickSyncProbeModel,
  QuickSyncCommitParams,
  QuickSyncCommitResult
} from '@/api/admin/channels'

export interface EditableModelItem extends QuickSyncProbeModel {
  selected: boolean
}

const props = defineProps<{
  visible: boolean
}>()

const emit = defineEmits<{
  (e: 'update:visible', value: boolean): void
  (e: 'success', data?: QuickSyncCommitResult): void
}>()

const { t } = useI18n()
const appStore = useAppStore()

const currentStep = ref<1 | 2 | 3>(1)
const loading = ref(false)
const showApiKey = ref(false)
const probeError = ref('')

const steps = computed(() => [
  { step: 1, title: t('admin.channels.quickSync.step1', '1. 连接与探测') },
  { step: 2, title: t('admin.channels.quickSync.step2', '2. 分组与计费配置') },
  { step: 3, title: t('admin.channels.quickSync.step3', '3. 确认与导入') }
])

const dialogTitle = computed(() => {
  return t('admin.channels.quickSync.modalTitle', '一键同步节点与多协议配置向导')
})

const form = reactive({
  name: '',
  base_url: '',
  api_key: '',
  platform: 'antigravity',
  default_group_id: 0 as number,
  createNewGroup: false,
  newGroupName: '',
  newGroupRateMultiplier: 1.0,
  enableMonitor: true,
  monitorModel: '',
  monitorInterval: 60
})

const billingStrategy = reactive({
  mode: 'ratio' as 'ratio' | 'per_request' | 'fixed',
  ratio: 1.0,
  per_request_price: 0.01 as number | null
})

const existingGroups = ref<AdminGroup[]>([])
const models = ref<EditableModelItem[]>([])
const batchTargetGroupId = ref<number | null>(null)

const isStep1Valid = computed(() => {
  return form.name.trim() !== '' && form.base_url.trim() !== '' && form.api_key.trim() !== ''
})

const isStep2Valid = computed(() => {
  if (models.value.length === 0) return false
  if (form.createNewGroup && form.newGroupName.trim() === '') return false
  return true
})

const selectedModelCount = computed(() => {
  return models.value.filter(m => m.selected).length
})

const allSelected = computed(() => {
  return models.value.length > 0 && models.value.every(m => m.selected)
})

const tokenModelCount = computed(() => {
  return models.value.filter(m => m.billing_mode === 'token').length
})

const perReqModelCount = computed(() => {
  return models.value.filter(m => m.billing_mode === 'per_request').length
})

const groupSummaries = computed(() => {
  const map: Record<string, string[]> = {}
  for (const m of models.value) {
    let groupName = '未知分组'
    if (m.target_group_id === 0) {
      groupName = form.createNewGroup
        ? `⭐ [新建] ${form.newGroupName || '新分组'}`
        : '默认分组'
    } else {
      const g = existingGroups.value.find(eg => eg.id === m.target_group_id)
      groupName = g ? `${g.name} (ID: ${g.id})` : `分组 ID: ${m.target_group_id}`
    }
    if (!map[groupName]) {
      map[groupName] = []
    }
    map[groupName].push(m.id)
  }
  return Object.entries(map).map(([name, modelIds]) => ({
    name,
    models: modelIds
  }))
})

// Load groups on mount or when modal opens
const loadGroups = async () => {
  try {
    if (adminAPI?.groups?.getAll) {
      const groups = await adminAPI.groups.getAll()
      existingGroups.value = groups || []
    } else if (adminAPI?.groups?.list) {
      const resp = await adminAPI.groups.list(1, 100)
      existingGroups.value = resp.items || []
    }
    if (existingGroups.value.length > 0 && form.default_group_id === 0 && !form.createNewGroup) {
      form.default_group_id = existingGroups.value[0].id
    }
  } catch (err: any) {
    console.error('Failed to load groups for quick sync:', err)
  }
}

watch(
  () => props.visible,
  (val) => {
    if (val) {
      loadGroups()
    } else {
      resetState()
    }
  },
  { immediate: true }
)

const resetState = () => {
  currentStep.value = 1
  probeError.value = ''
  loading.value = false
  form.name = ''
  form.base_url = ''
  form.api_key = ''
  form.platform = 'antigravity'
  form.default_group_id = existingGroups.value[0]?.id ?? 0
  form.createNewGroup = false
  form.newGroupName = ''
  form.newGroupRateMultiplier = 1.0
  billingStrategy.mode = 'ratio'
  billingStrategy.ratio = 1.0
  billingStrategy.per_request_price = 0.01
  models.value = []
}

const handleClose = () => {
  emit('update:visible', false)
}

const prevStep = () => {
  if (currentStep.value > 1) {
    currentStep.value = (currentStep.value - 1) as 1 | 2 | 3
  }
}

const nextStep = () => {
  if (currentStep.value < 3) {
    currentStep.value = (currentStep.value + 1) as 1 | 2 | 3
  }
}

// 探测上游模型
const handleProbe = async () => {
  if (!isStep1Valid.value) return
  loading.value = true
  probeError.value = ''

  try {
    const result = await adminAPI.channels.quickSyncProbe({
      base_url: form.base_url.trim(),
      api_key: form.api_key.trim(),
      platform: form.platform
    })

    const initialGroupId = form.createNewGroup ? 0 : form.default_group_id

    models.value = (result.models || []).map(m => {
      const baseIn = m.base_price_in ?? null
      const baseOut = m.base_price_out ?? null
      const priceIn = m.price_in ?? baseIn
      const priceOut = m.price_out ?? baseOut
      return {
        ...m,
        target_group_id: m.target_group_id ?? initialGroupId,
        price_in: priceIn,
        price_out: priceOut,
        per_request_price: m.per_request_price ?? null,
        selected: true
      }
    })

    if (result.warnings && result.warnings.length > 0) {
      appStore?.showWarning?.(result.warnings.join('; '))
    }

    if (models.value.length > 0 && !form.monitorModel) {
      form.monitorModel = models.value[0].id
    }

    // Auto set new group name if empty
    if (!form.newGroupName) {
      form.newGroupName = form.name ? `${form.name}-group` : 'quick-sync-group'
    }

    currentStep.value = 2
  } catch (err: any) {
    probeError.value = err.message || t('admin.channels.quickSync.probeFailed', '探测失败，请检查 Base URL 与 API Key 是否正确')
    appStore?.showError?.(probeError.value)
  } finally {
    loading.value = false
  }
}

const onDefaultGroupChanged = () => {
  const newDefault = form.default_group_id
  for (const m of models.value) {
    m.target_group_id = newDefault
  }
}

const onToggleCreateNewGroup = () => {
  if (form.createNewGroup) {
    if (!form.newGroupName) {
      form.newGroupName = form.name ? `${form.name}-group` : 'quick-sync-group'
    }
    form.default_group_id = 0
    for (const m of models.value) {
      m.target_group_id = 0
    }
  } else {
    form.default_group_id = existingGroups.value[0]?.id ?? 0
    for (const m of models.value) {
      m.target_group_id = form.default_group_id
    }
  }
}

const setBillingStrategyMode = (mode: 'ratio' | 'per_request' | 'fixed') => {
  billingStrategy.mode = mode
  applyBillingStrategy()
}

const applyBillingStrategy = () => {
  if (billingStrategy.mode === 'ratio') {
    const r = billingStrategy.ratio || 1.0
    for (const m of models.value) {
      m.billing_mode = 'token'
      m.price_in = m.base_price_in != null ? Number((m.base_price_in * r).toFixed(8)) : null
      m.price_out = m.base_price_out != null ? Number((m.base_price_out * r).toFixed(8)) : null
      m.per_request_price = null
    }
  } else if (billingStrategy.mode === 'per_request') {
    const fixed = billingStrategy.per_request_price != null ? billingStrategy.per_request_price : 0.01
    for (const m of models.value) {
      m.billing_mode = 'per_request'
      m.per_request_price = fixed
      m.price_in = null
      m.price_out = null
    }
  }
}

const toggleSelectAll = (e: Event) => {
  const checked = (e.target as HTMLInputElement).checked
  for (const m of models.value) {
    m.selected = checked
  }
}

const applyBatchGroup = () => {
  if (batchTargetGroupId.value === null) return
  const targetId = batchTargetGroupId.value
  for (const m of models.value) {
    if (m.selected) {
      m.target_group_id = targetId
    }
  }
  appStore?.showSuccess?.(t('admin.channels.quickSync.batchGroupApplied', '已批量设置所选模型分组'))
}

const applyBatchBillingMode = (mode: 'token' | 'per_request') => {
  for (const m of models.value) {
    if (m.selected) {
      m.billing_mode = mode
      if (mode === 'token') {
        const r = billingStrategy.ratio || 1.0
        m.price_in = m.base_price_in != null ? Number((m.base_price_in * r).toFixed(8)) : null
        m.price_out = m.base_price_out != null ? Number((m.base_price_out * r).toFixed(8)) : null
        m.per_request_price = null
      } else {
        m.per_request_price = billingStrategy.per_request_price ?? 0.01
        m.price_in = null
        m.price_out = null
      }
    }
  }
}

// 提交导入落库
const handleCommit = async () => {
  loading.value = true
  try {
    const payload: QuickSyncCommitParams = {
      name: form.name.trim(),
      base_url: form.base_url.trim(),
      api_key: form.api_key.trim(),
      platform: form.platform,
      default_group_id: form.default_group_id || undefined,
      new_group: form.createNewGroup
        ? {
            create: true,
            name: form.newGroupName.trim() || form.name.trim(),
            rate_multiplier: form.newGroupRateMultiplier || 1.0
          }
        : undefined,
      billing_strategy: {
        mode: billingStrategy.mode,
        ratio: billingStrategy.ratio,
        per_request_price: billingStrategy.per_request_price
      },
      models: models.value.map(m => ({
        model: m.id,
        target_group_id: m.target_group_id ?? 0,
        billing_mode: m.billing_mode,
        input_price: m.billing_mode === 'token' ? m.price_in : undefined,
        output_price: m.billing_mode === 'token' ? m.price_out : undefined,
        per_request_price: m.billing_mode === 'per_request' ? m.per_request_price : undefined
      })),
      enable_monitor: form.enableMonitor,
      monitor_model: form.enableMonitor ? (form.monitorModel || models.value[0]?.id) : undefined,
      monitor_interval: form.enableMonitor ? (form.monitorInterval || 60) : undefined
    }

    const result = await adminAPI.channels.quickSyncCommit(payload)

    appStore?.showSuccess?.(t('admin.channels.quickSync.commitSuccess', '一键同步节点配置成功！'))
    emit('success', result)
    emit('update:visible', false)
  } catch (err: any) {
    appStore?.showError?.(err.message || t('admin.channels.quickSync.commitFailed', '导入失败，请重试'))
  } finally {
    loading.value = false
  }
}
</script>
