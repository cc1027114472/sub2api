import { beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'
import { defineComponent } from 'vue'
import ChannelProbeModelsModal from './ChannelProbeModelsModal.vue'

const {
  quickSyncProbeMock
} = vi.hoisted(() => ({
  quickSyncProbeMock: vi.fn()
}))

vi.mock('@/api/admin', () => ({
  default: {
    channels: {
      quickSyncProbe: quickSyncProbeMock
    }
  }
}))

vi.mock('vue-i18n', async (importOriginal) => {
  const actual = await importOriginal<typeof import('vue-i18n')>()
  const messages: Record<string, string> = {
    'admin.channels.probeModal.title': '从上游探测并选择模型',
    'admin.channels.probeModal.startProbe': '连接并探测模型',
    'admin.channels.probeModal.alreadyInChannel': '渠道中已有',
    'admin.channels.probeModal.canAdd': '未添加',
    'admin.channels.probeModal.confirmImport': '导入已选'
  }
  return {
    ...actual,
    useI18n: () => ({
      t: (key: string, valuesOrDefault?: any, defaultVal?: string) => {
        if (messages[key]) return messages[key]
        if (typeof valuesOrDefault === 'string') return valuesOrDefault
        if (defaultVal) return defaultVal
        return key
      }
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
  template: '<i :data-icon="name" />'
})

describe('ChannelProbeModelsModal', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    localStorage.clear()
  })

  it('renders correctly and initializes form from localStorage if present', () => {
    localStorage.setItem('sub2api_upstream_base_url', 'http://154.36.173.146')
    localStorage.setItem('sub2api_upstream_api_key', 'test-key-123')

    const wrapper = mount(ChannelProbeModelsModal, {
      props: {
        visible: true,
        platform: 'antigravity',
        existingModels: ['claude-sonnet-4-6'],
        pricingEntriesCount: 1
      },
      global: {
        stubs: {
          BaseDialog: BaseDialogStub,
          Icon: IconStub
        }
      }
    })

    const inputs = wrapper.findAll('input')
    expect(inputs[0].element.value).toBe('http://154.36.173.146')
    expect(inputs[1].element.value).toBe('test-key-123')
  })

  it('probes upstream models and marks existing models disabled', async () => {
    quickSyncProbeMock.mockResolvedValueOnce({
      models: [
        { id: 'claude-sonnet-4-6', name: 'Claude Sonnet 4.6' },
        { id: 'gemini-3-pro', name: 'Gemini 3 Pro' },
        { id: 'gemini-3-flash', name: 'Gemini 3 Flash' }
      ]
    })

    const wrapper = mount(ChannelProbeModelsModal, {
      props: {
        visible: true,
        platform: 'antigravity',
        existingModels: ['claude-sonnet-4-6'],
        pricingEntriesCount: 1
      },
      global: {
        stubs: {
          BaseDialog: BaseDialogStub,
          Icon: IconStub
        }
      }
    })

    // Set base url and api key
    const inputs = wrapper.findAll('input')
    await inputs[0].setValue('http://154.36.173.146')
    await inputs[1].setValue('test-key-123')

    // Find and click the probe button
    const probeBtn = wrapper.findAll('button').find(b => b.text().includes('连接并探测模型'))
    expect(probeBtn).toBeDefined()
    await probeBtn!.trigger('click')
    await flushPromises()

    expect(quickSyncProbeMock).toHaveBeenCalledWith({
      base_url: 'http://154.36.173.146',
      api_key: 'test-key-123',
      platform: 'antigravity'
    })

    // Verify localStorage persistence
    expect(localStorage.getItem('sub2api_upstream_base_url')).toBe('http://154.36.173.146')
    expect(localStorage.getItem('sub2api_upstream_api_key')).toBe('test-key-123')

    // claude-sonnet-4-6 should be disabled and marked as already in channel
    expect(wrapper.text()).toContain('渠道中已有')
    expect(wrapper.text()).toContain('未添加')

    // gemini-3-pro and gemini-3-flash are unadded, should be selected by default
    expect(wrapper.text()).toContain('gemini-3-pro')
    expect(wrapper.text()).toContain('gemini-3-flash')
  })

  it('emits import-models with selected models when confirmed', async () => {
    quickSyncProbeMock.mockResolvedValueOnce({
      models: [
        { id: 'claude-sonnet-4-6', name: 'Claude Sonnet 4.6' },
        { id: 'gemini-3-pro', name: 'Gemini 3 Pro' }
      ]
    })

    const wrapper = mount(ChannelProbeModelsModal, {
      props: {
        visible: true,
        platform: 'antigravity',
        existingModels: ['claude-sonnet-4-6'],
        pricingEntriesCount: 2
      },
      global: {
        stubs: {
          BaseDialog: BaseDialogStub,
          Icon: IconStub
        }
      }
    })

    const inputs = wrapper.findAll('input')
    await inputs[0].setValue('http://154.36.173.146')
    await inputs[1].setValue('test-key-123')

    const probeBtn = wrapper.findAll('button').find(b => b.text().includes('连接并探测模型'))
    await probeBtn!.trigger('click')
    await flushPromises()

    // Confirm import
    const confirmBtn = wrapper.findAll('button').find(b => b.text().includes('导入已选'))
    expect(confirmBtn).toBeDefined()
    await confirmBtn!.trigger('click')

    const emitted = wrapper.emitted('import-models')
    expect(emitted).toBeDefined()
    expect(emitted![0][0]).toEqual({
      models: ['gemini-3-pro'],
      targetEntryIndex: 'new'
    })
  })
})
