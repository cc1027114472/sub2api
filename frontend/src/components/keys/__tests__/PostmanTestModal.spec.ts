import { afterEach, describe, expect, it, vi } from 'vitest'
import { mount } from '@vue/test-utils'

const { copyToClipboardMock } = vi.hoisted(() => ({
  copyToClipboardMock: vi.fn().mockResolvedValue(true)
}))

vi.mock('vue-i18n', () => ({
  useI18n: () => ({
    t: (key: string) => key
  })
}))

vi.mock('@/composables/useClipboard', () => ({
  useClipboard: () => ({
    copyToClipboard: copyToClipboardMock
  })
}))

import PostmanTestModal from '../PostmanTestModal.vue'

describe('PostmanTestModal', () => {
  afterEach(() => {
    vi.clearAllMocks()
  })

  it('renders 4 protocol tabs correctly', () => {
    const wrapper = mount(PostmanTestModal, {
      props: {
        show: true,
        apiKey: 'sk-test-1234567890abcdef',
        baseUrl: 'https://ukapi.cc',
        platform: 'antigravity',
        allowedModels: ['gemini-3-flash', 'claude-sonnet-4-6']
      },
      global: {
        stubs: {
          BaseDialog: {
            template: '<div><slot /></div>'
          },
          Icon: {
            template: '<span />'
          }
        }
      }
    })

    const buttons = wrapper.findAll('nav[aria-label="Protocol Tabs"] button')
    expect(buttons.length).toBe(4)
    expect(buttons[0].text()).toContain('OpenAI')
    expect(buttons[1].text()).toContain('Claude')
    expect(buttons[2].text()).toContain('Gemini')
    expect(buttons[3].text()).toContain('Responses')
  })

  it('generates correct cURL for OpenAI Chat Completions by default', () => {
    const wrapper = mount(PostmanTestModal, {
      props: {
        show: true,
        apiKey: 'sk-test-key-123456',
        baseUrl: 'https://ukapi.cc',
        platform: 'antigravity',
        allowedModels: ['gemini-3-flash']
      },
      global: {
        stubs: {
          BaseDialog: { template: '<div><slot /></div>' },
          Icon: { template: '<span />' }
        }
      }
    })

    const code = wrapper.find('pre code').text()
    expect(code).toContain('curl -X POST "https://ukapi.cc/v1/chat/completions"')
    expect(code).toContain('-H "Authorization: Bearer sk-test-key-123456"')
    expect(code).toContain('"model": "gemini-3-flash"')
  })

  it('switches to Claude Messages and updates cURL command', async () => {
    const wrapper = mount(PostmanTestModal, {
      props: {
        show: true,
        apiKey: 'sk-test-key-123456',
        baseUrl: 'https://ukapi.cc',
        platform: 'antigravity',
        allowedModels: ['claude-sonnet-4-6']
      },
      global: {
        stubs: {
          BaseDialog: { template: '<div><slot /></div>' },
          Icon: { template: '<span />' }
        }
      }
    })

    const buttons = wrapper.findAll('nav[aria-label="Protocol Tabs"] button')
    await buttons[1].trigger('click') // Click Claude tab

    const code = wrapper.find('pre code').text()
    expect(code).toContain('curl -X POST "https://ukapi.cc/v1/messages"')
    expect(code).toContain('-H "x-api-key: sk-test-key-123456"')
    expect(code).toContain('-H "anthropic-version: 2023-06-01"')
    expect(code).toContain('"model": "claude-sonnet-4-6"')
  })

  it('switches to Gemini Native and updates cURL command with query key', async () => {
    const wrapper = mount(PostmanTestModal, {
      props: {
        show: true,
        apiKey: 'sk-test-key-123456',
        baseUrl: 'https://ukapi.cc',
        platform: 'antigravity',
        allowedModels: ['gemini-2.5-flash']
      },
      global: {
        stubs: {
          BaseDialog: { template: '<div><slot /></div>' },
          Icon: { template: '<span />' }
        }
      }
    })

    const buttons = wrapper.findAll('nav[aria-label="Protocol Tabs"] button')
    await buttons[2].trigger('click') // Click Gemini tab

    const code = wrapper.find('pre code').text()
    expect(code).toContain('curl -X POST "https://ukapi.cc/v1beta/models/gemini-2.5-flash:generateContent?key=sk-test-key-123456"')
    expect(code).toContain('"contents"')
  })

  it('switches to OpenAI Responses and updates cURL command', async () => {
    const wrapper = mount(PostmanTestModal, {
      props: {
        show: true,
        apiKey: 'sk-test-key-123456',
        baseUrl: 'https://ukapi.cc',
        platform: 'antigravity',
        allowedModels: ['gemini-3-flash']
      },
      global: {
        stubs: {
          BaseDialog: { template: '<div><slot /></div>' },
          Icon: { template: '<span />' }
        }
      }
    })

    const buttons = wrapper.findAll('nav[aria-label="Protocol Tabs"] button')
    await buttons[3].trigger('click') // Click Responses tab

    const code = wrapper.find('pre code').text()
    expect(code).toContain('curl -X POST "https://ukapi.cc/v1/responses"')
    expect(code).toContain('-H "Authorization: Bearer sk-test-key-123456"')
    expect(code).toContain('"input":')
  })
})
