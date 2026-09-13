import { describe, expect, it, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import ModelPlazaContent from '../ModelPlazaContent.vue'
import type { ModelPlazaResponse } from '@/api/modelPlaza'

vi.mock('vue-i18n', async () => {
  const actual = await vi.importActual<typeof import('vue-i18n')>('vue-i18n')
  return {
    ...actual,
    useI18n: () => ({
      t: (key: string) => key
    })
  }
})

vi.mock('@/stores/app', () => ({
  useAppStore: () => ({ cachedPublicSettings: null })
}))

beforeEach(() => {
  setActivePinia(createPinia())
})

function mockResponse(): ModelPlazaResponse {
  return {
    is_authenticated: false,
    groups: [
      {
        id: 1,
        name: 'OpenAI Group',
        description: '',
        platform: 'openai',
        subscription_type: 'standard',
        rate_multiplier: 1,
        peak_rate_enabled: false,
        peak_start: '',
        peak_end: '',
        peak_rate_multiplier: 1,
        is_exclusive: false,
        image_rate_independent: false,
        image_rate_multiplier: 1,
        long_context_pricing_enabled: true,
        models: [
          {
            name: 'gpt-4o',
            platform: 'openai',
            pricing: {
              billing_mode: 'token',
              input_price: 5e-6,
              output_price: 1.5e-5,
              cache_write_price: null,
              cache_read_price: null,
              image_input_price: null,
              image_output_price: null,
              per_request_price: null,
              intervals: []
            },
            official_pricing: null
          },
          {
            name: 'gemini-3.8-flash-high',
            platform: 'openai',
            pricing: {
              billing_mode: 'token',
              input_price: 1e-6,
              output_price: 3e-6,
              cache_write_price: null,
              cache_read_price: null,
              image_input_price: null,
              image_output_price: null,
              per_request_price: null,
              intervals: []
            },
            official_pricing: null
          }
        ]
      }
    ]
  }
}

describe('ModelPlazaContent 多关键词分词搜索', () => {
  it('支持空格分隔多个关键词同时命中模型名称', async () => {
    const wrapper = mount(ModelPlazaContent, {
      props: {
        loading: false,
        error: null,
        response: mockResponse()
      },
      global: {
        stubs: {
          PlazaNavBar: true,
          PlazaFilterBar: {
            template: '<div class="filter-bar-stub" />',
            props: ['search', 'platform', 'groupId', 'rate']
          },
          PlazaGroupSection: {
            template: '<div class="group-section-stub">{{ group.models.map(m => m.name).join(",") }}</div>',
            props: ['group']
          }
        }
      }
    })

    // 默认展示全部模型
    expect(wrapper.text()).toContain('gpt-4o,gemini-3.8-flash-high')

    // 模拟搜索 "3.8 flash" (空格分隔关键词)
    const vm = wrapper.vm as any
    vm.searchQuery = '3.8 flash'
    await wrapper.vm.$nextTick()

    // 命中 gemini-3.8-flash-high，排除 gpt-4o
    expect(wrapper.text()).toContain('gemini-3.8-flash-high')
    expect(wrapper.text()).not.toContain('gpt-4o')
  })
})
