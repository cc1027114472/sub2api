import { describe, expect, it } from 'vitest'
import {
  buildMowanImportDeeplink,
  buildMowanWebImportUrl,
  resolveMowanEndpoint
} from '../mowanImport'

describe('mowanImport utils', () => {
  it('correctly resolves endpoint for antigravity platform', () => {
    expect(resolveMowanEndpoint('https://api.example.com', 'antigravity')).toBe(
      'https://api.example.com/antigravity'
    )
    expect(resolveMowanEndpoint('https://api.example.com/', 'antigravity')).toBe(
      'https://api.example.com/antigravity'
    )
  })

  it('correctly resolves endpoint for standard v1 platform', () => {
    expect(resolveMowanEndpoint('https://api.example.com', 'openai')).toBe(
      'https://api.example.com/v1'
    )
    expect(resolveMowanEndpoint('https://api.example.com/v1', 'openai')).toBe(
      'https://api.example.com/v1'
    )
  })

  it('builds mowan:// deeplink correctly', () => {
    const link = buildMowanImportDeeplink({
      baseUrl: 'https://api.example.com',
      platform: 'openai',
      providerName: 'TestProvider',
      apiKey: 'sk-123456',
      models: ['gpt-4o', 'claude-3-5-sonnet']
    })

    expect(link.startsWith('mowan://v1/import?')).toBe(true)
    expect(link).toContain('resource=provider')
    expect(link).toContain('name=TestProvider')
    expect(link).toContain('apiKey=sk-123456')
    expect(link).toContain('models=gpt-4o%2Cclaude-3-5-sonnet')
    expect(link).toContain('endpoint=https%3A%2F%2Fapi.example.com%2Fv1')
  })

  it('builds local web import URL targeting port 3090', () => {
    const url = buildMowanWebImportUrl({
      baseUrl: 'https://api.example.com',
      platform: 'antigravity',
      providerName: 'MowanTest',
      apiKey: 'sk-abcdef'
    })

    expect(url.startsWith('http://127.0.0.1:3090/#/settings?')).toBe(true)
    expect(url).toContain('action=import-provider')
    expect(url).toContain('name=MowanTest')
    expect(url).toContain('apiKey=sk-abcdef')
  })
})
