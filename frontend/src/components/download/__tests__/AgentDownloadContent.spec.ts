import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import AgentDownloadContent from '../AgentDownloadContent.vue'
import { createPinia, setActivePinia } from 'pinia'

vi.mock('vue-i18n', async () => {
  const actual = await vi.importActual<typeof import('vue-i18n')>('vue-i18n')
  return {
    ...actual,
    useI18n: () => ({
      t: (key: string) => key,
    }),
  }
})

describe('AgentDownloadContent', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('renders brand title and download buttons', () => {
    const wrapper = mount(AgentDownloadContent, {
      global: {
        plugins: [createPinia()],
        stubs: {
          Icon: true,
          RouterLink: true,
        },
      },
    })

    expect(wrapper.text()).toContain('agentDownload.title')
    expect(wrapper.text()).toContain('agentDownload.downloadPrimary')
    expect(wrapper.text()).toContain('agentDownload.quickStartTitle')
    expect(wrapper.text()).toContain('agentDownload.features.title')
  })
})
