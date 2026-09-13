import { beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'
import { defineComponent } from 'vue'
import ChannelQuickSyncModal from './ChannelQuickSyncModal.vue'

const {
  quickSyncProbeMock,
  quickSyncCommitMock,
  getAllGroupsMock,
  showSuccessMock,
  showErrorMock,
  showWarningMock
} = vi.hoisted(() => ({
  quickSyncProbeMock: vi.fn(),
  quickSyncCommitMock: vi.fn(),
  getAllGroupsMock: vi.fn(),
  showSuccessMock: vi.fn(),
  showErrorMock: vi.fn(),
  showWarningMock: vi.fn()
}))

vi.mock('@/api/admin', () => ({
  adminAPI: {
    channels: {
      quickSyncProbe: quickSyncProbeMock,
      quickSyncCommit: quickSyncCommitMock
    },
    groups: {
      getAll: getAllGroupsMock
    }
  }
}))

vi.mock('@/stores/app', () => ({
  useAppStore: () => ({
    showSuccess: showSuccessMock,
    showError: showErrorMock,
    showWarning: showWarningMock
  })
}))

vi.mock('vue-i18n', async (importOriginal) => {
  const actual = await importOriginal<typeof import('vue-i18n')>()
  return {
    ...actual,
    useI18n: () => ({
      t: (key: string, defaultVal?: string) => defaultVal || key
    })
  }
})

const BaseDialogStub = defineComponent({
  name: 'BaseDialog',
  props: { show: { type: Boolean, default: false }, title: String, width: String },
  template: '<div v-if="show" class="base-dialog-stub"><slot /><slot name="footer" /></div>'
})

const IconStub = defineComponent({
  name: 'Icon',
  props: { name: String, size: String },
  template: '<span class="icon-stub">{{ name }}</span>'
})

const mockGroups = [
  { id: 1, name: 'Default Group', platform: 'antigravity', status: 'active', rate_multiplier: 1.0 },
  { id: 2, name: 'VIP Group', platform: 'antigravity', status: 'active', rate_multiplier: 1.5 }
]

const mockProbeModels = [
  {
    id: 'gemini-2.5-pro',
    display_name: 'Gemini 2.5 Pro',
    base_price_in: 0.000001,
    base_price_out: 0.000002,
    billing_mode: 'token' as const
  },
  {
    id: 'claude-3-7-sonnet',
    display_name: 'Claude 3.7 Sonnet',
    base_price_in: 0.000003,
    base_price_out: 0.000015,
    billing_mode: 'token' as const
  }
]

describe('ChannelQuickSyncModal.vue', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    getAllGroupsMock.mockResolvedValue(mockGroups)
    quickSyncProbeMock.mockResolvedValue({
      models: mockProbeModels,
      total: 2,
      warnings: []
    })
    quickSyncCommitMock.mockResolvedValue({
      channel_id: 101,
      account_id: 202,
      group_ids: [1],
      model_count: 2
    })
  })

  const createWrapper = (props = { visible: true }) => {
    return mount(ChannelQuickSyncModal, {
      props,
      global: {
        stubs: {
          BaseDialog: BaseDialogStub,
          Icon: IconStub
        }
      }
    })
  }

  it('renders initial step 1 with form fields and disabled probe button when empty', async () => {
    const wrapper = createWrapper()
    await flushPromises()

    expect(wrapper.find('[data-test="channel-name-input"]').exists()).toBe(true)
    expect(wrapper.find('[data-test="base-url-input"]').exists()).toBe(true)
    expect(wrapper.find('[data-test="api-key-input"]').exists()).toBe(true)
    expect(wrapper.find('[data-test="platform-select"]').exists()).toBe(true)

    const probeBtn = wrapper.find('[data-test="probe-btn"]')
    expect(probeBtn.exists()).toBe(true)
    expect((probeBtn.element as HTMLButtonElement).disabled).toBe(true)
  })

  it('calls quickSyncProbe and transitions to step 2 with populated models', async () => {
    const wrapper = createWrapper()
    await flushPromises()

    await wrapper.find('[data-test="channel-name-input"]').setValue('Antigravity-Node-1')
    await wrapper.find('[data-test="base-url-input"]').setValue('http://127.0.0.1:8080')
    await wrapper.find('[data-test="api-key-input"]').setValue('sk-test-secret-key')
    await wrapper.find('[data-test="platform-select"]').setValue('antigravity')

    const probeBtn = wrapper.find('[data-test="probe-btn"]')
    expect((probeBtn.element as HTMLButtonElement).disabled).toBe(false)

    await probeBtn.trigger('click')
    await flushPromises()

    expect(quickSyncProbeMock).toHaveBeenCalledWith({
      base_url: 'http://127.0.0.1:8080',
      api_key: 'sk-test-secret-key',
      platform: 'antigravity'
    })

    // Should now be on Step 2
    expect(wrapper.find('[data-test="default-group-select"]').exists()).toBe(true)
    expect(wrapper.find('[data-test="model-row-gemini-2.5-pro"]').exists()).toBe(true)
    expect(wrapper.find('[data-test="model-row-claude-3-7-sonnet"]').exists()).toBe(true)
  })

  it('updates model prices when changing billing strategy multiplier or per-request mode', async () => {
    const wrapper = createWrapper()
    await flushPromises()

    // Step 1 -> Step 2
    await wrapper.find('[data-test="channel-name-input"]').setValue('Antigravity-Node-1')
    await wrapper.find('[data-test="base-url-input"]').setValue('http://127.0.0.1:8080')
    await wrapper.find('[data-test="api-key-input"]').setValue('sk-test-secret-key')
    await wrapper.find('[data-test="probe-btn"]').trigger('click')
    await flushPromises()

    // 1. Ratio mode change to 2.0
    const ratioInput = wrapper.find('[data-test="strategy-ratio-input"]')
    await ratioInput.setValue(2.0)
    await ratioInput.trigger('input')
    await flushPromises()

    const priceInInput = wrapper.find('[data-test="model-price-in-gemini-2.5-pro"]')
    expect((priceInInput.element as HTMLInputElement).value).toBe('0.000002') // 0.000001 * 2

    // 2. Switch to Per-Request mode
    await wrapper.find('[data-test="strategy-per-request-btn"]').trigger('click')
    await flushPromises()

    const perReqInput = wrapper.find('[data-test="strategy-per-request-input"]')
    await perReqInput.setValue(0.05)
    await perReqInput.trigger('input')
    await flushPromises()

    const modelModeSelect = wrapper.find('[data-test="model-mode-gemini-2.5-pro"]')
    expect((modelModeSelect.element as HTMLSelectElement).value).toBe('per_request')

    const modelPerReqPrice = wrapper.find('[data-test="model-price-per-req-gemini-2.5-pro"]')
    expect((modelPerReqPrice.element as HTMLInputElement).value).toBe('0.05')
  })

  it('allows moving individual models and batch moving models to different target groups', async () => {
    const wrapper = createWrapper()
    await flushPromises()

    // Step 1 -> Step 2
    await wrapper.find('[data-test="channel-name-input"]').setValue('Antigravity-Node-1')
    await wrapper.find('[data-test="base-url-input"]').setValue('http://127.0.0.1:8080')
    await wrapper.find('[data-test="api-key-input"]').setValue('sk-test-secret-key')
    await wrapper.find('[data-test="probe-btn"]').trigger('click')
    await flushPromises()

    // 1. Change individual model target group
    const geminiGroupSelect = wrapper.find('[data-test="model-group-gemini-2.5-pro"]')
    await geminiGroupSelect.setValue(2) // VIP Group
    expect((geminiGroupSelect.element as HTMLSelectElement).value).toBe('2')

    // 2. Batch move: select all and move to Group 1
    const selectAll = wrapper.find('[data-test="select-all-checkbox"]')
    await selectAll.setValue(true)
    await selectAll.trigger('change')
    await flushPromises()

    const batchGroupSelect = wrapper.find('[data-test="batch-group-select"]')
    await batchGroupSelect.setValue(1)
    await wrapper.find('[data-test="apply-batch-group-btn"]').trigger('click')
    await flushPromises()

    expect((wrapper.find('[data-test="model-group-gemini-2.5-pro"]').element as HTMLSelectElement).value).toBe('1')
    expect((wrapper.find('[data-test="model-group-claude-3-7-sonnet"]').element as HTMLSelectElement).value).toBe('1')
  })

  it('supports creating a new independent group', async () => {
    const wrapper = createWrapper()
    await flushPromises()

    // Step 1 -> Step 2
    await wrapper.find('[data-test="channel-name-input"]').setValue('Antigravity-Node-1')
    await wrapper.find('[data-test="base-url-input"]').setValue('http://127.0.0.1:8080')
    await wrapper.find('[data-test="api-key-input"]').setValue('sk-test-secret-key')
    await wrapper.find('[data-test="probe-btn"]').trigger('click')
    await flushPromises()

    // Check create new group
    const createNewGroupCheckbox = wrapper.find('[data-test="create-new-group-checkbox"]')
    await createNewGroupCheckbox.setValue(true)
    await createNewGroupCheckbox.trigger('change')
    await flushPromises()

    expect(wrapper.find('[data-test="new-group-name-input"]').exists()).toBe(true)
    await wrapper.find('[data-test="new-group-name-input"]').setValue('Exclusive-Node-Group')
    await wrapper.find('[data-test="new-group-multiplier-input"]').setValue(1.8)
    await flushPromises()

    // Default group select should now have option 0
    expect((wrapper.find('[data-test="default-group-select"]').element as HTMLSelectElement).value).toBe('0')
  })

  it('completes the full wizard flow and commits quick sync payload on step 3', async () => {
    const wrapper = createWrapper()
    await flushPromises()

    // Step 1: Probe
    await wrapper.find('[data-test="channel-name-input"]').setValue('Antigravity-Node-1')
    await wrapper.find('[data-test="base-url-input"]').setValue('http://127.0.0.1:8080')
    await wrapper.find('[data-test="api-key-input"]').setValue('sk-test-secret-key')
    await wrapper.find('[data-test="probe-btn"]').trigger('click')
    await flushPromises()

    // Step 2: Configure & Click Next
    await wrapper.find('[data-test="step2-next-btn"]').trigger('click')
    await flushPromises()

    // Step 3: Summary view & commit
    expect(wrapper.find('[data-test="commit-btn"]').exists()).toBe(true)
    expect(wrapper.text()).toContain('Antigravity-Node-1')
    expect(wrapper.text()).toContain('http://127.0.0.1:8080')

    await wrapper.find('[data-test="commit-btn"]').trigger('click')
    await flushPromises()

    expect(quickSyncCommitMock).toHaveBeenCalledTimes(1)
    expect(quickSyncCommitMock).toHaveBeenCalledWith(expect.objectContaining({
      name: 'Antigravity-Node-1',
      base_url: 'http://127.0.0.1:8080',
      api_key: 'sk-test-secret-key',
      platform: 'antigravity',
      models: expect.arrayContaining([
        expect.objectContaining({
          model: 'gemini-2.5-pro'
        }),
        expect.objectContaining({
          model: 'claude-3-7-sonnet'
        })
      ])
    }))

    expect(wrapper.emitted('success')).toBeTruthy()
    expect(wrapper.emitted('update:visible')?.[0]).toEqual([false])
  })
})
