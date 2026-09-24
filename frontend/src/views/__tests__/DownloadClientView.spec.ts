import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import DownloadClientView from '../DownloadClientView.vue'
import { createPinia, setActivePinia } from 'pinia'
import { useAuthStore } from '@/stores/auth'

const mockRoute = {
  query: { embedded: '' }
}

vi.mock('vue-router', () => ({
  useRoute: () => mockRoute,
  RouterLink: {
    template: '<a><slot /></a>',
  },
}))

vi.mock('@/stores/app', () => ({
  useAppStore: () => ({
    cachedPublicSettings: { site_name: '英国api.cc' },
    fetchPublicSettings: vi.fn(),
  })
}))

vi.mock('vue-i18n', async () => {
  const actual = await vi.importActual<typeof import('vue-i18n')>('vue-i18n')
  return {
    ...actual,
    useI18n: () => ({
      t: (key: string) => key,
    }),
  }
})

describe('DownloadClientView', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    mockRoute.query.embedded = ''
  })

  it('renders standalone mode by default when not embedded', () => {
    const wrapper = mount(DownloadClientView, {
      global: {
        plugins: [createPinia()],
        stubs: {
          DownloadNavBar: { template: '<div class="download-nav-bar-stub"></div>' },
          AgentDownloadContent: { template: '<div class="agent-download-content-stub"></div>' },
          AppLayout: { template: '<div class="app-layout-stub"><slot /></div>' },
          Icon: true,
        },
      },
    })

    expect(wrapper.find('.download-nav-bar-stub').exists()).toBe(true)
    expect(wrapper.find('.agent-download-content-stub').exists()).toBe(true)
    expect(wrapper.find('.app-layout-stub').exists()).toBe(false)
  })

  it('renders embedded layout when embedded=1 and authenticated', () => {
    mockRoute.query.embedded = '1'
    const pinia = createPinia()
    const authStore = useAuthStore(pinia)
    authStore.token = 'mock-token'
    authStore.user = { id: 1, email: 'test@example.com', role: 'admin' } as any

    const wrapper = mount(DownloadClientView, {
      global: {
        plugins: [pinia],
        stubs: {
          DownloadNavBar: { template: '<div class="download-nav-bar-stub"></div>' },
          AgentDownloadContent: { template: '<div class="agent-download-content-stub"></div>' },
          AppLayout: { template: '<div class="app-layout-stub"><slot /></div>' },
          Icon: true,
        },
      },
    })

    expect(wrapper.find('.app-layout-stub').exists()).toBe(true)
    expect(wrapper.find('.agent-download-content-stub').exists()).toBe(true)
    expect(wrapper.find('.download-nav-bar-stub').exists()).toBe(false)
  })
})
