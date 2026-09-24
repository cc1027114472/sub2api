import type { GroupPlatform } from '@/types'

export interface MowanImportDeeplinkInput {
  baseUrl: string
  platform?: GroupPlatform | null
  providerName: string
  apiKey: string
  models?: string[]
}

/**
 * 根据平台类型处理 endpoint 地址
 */
export function resolveMowanEndpoint(baseUrl: string, platform?: GroupPlatform | null): string {
  const normalized = baseUrl.replace(/\/+$/, '')
  if (platform === 'antigravity') {
    return `${normalized}/antigravity`
  }
  return normalized.endsWith('/v1') ? normalized : `${normalized}/v1`
}

/**
 * 构造 mowan:// 协议深度链接，用于唤起桌面客户端
 */
export function buildMowanImportDeeplink(input: MowanImportDeeplinkInput): string {
  const endpoint = resolveMowanEndpoint(input.baseUrl, input.platform)
  const params = new URLSearchParams([
    ['resource', 'provider'],
    ['name', input.providerName],
    ['endpoint', endpoint],
    ['apiKey', input.apiKey],
    ['platform', input.platform || 'openai']
  ])

  if (input.models && input.models.length > 0) {
    params.set('models', input.models.join(','))
  }

  return `mowan://v1/import?${params.toString()}`
}

/**
 * 构造本地 Web GUI 的直连导入链接 (默认 3090 端口)
 */
export function buildMowanWebImportUrl(input: MowanImportDeeplinkInput, localPort = 3090): string {
  const endpoint = resolveMowanEndpoint(input.baseUrl, input.platform)
  const params = new URLSearchParams([
    ['action', 'import-provider'],
    ['name', input.providerName],
    ['endpoint', endpoint],
    ['apiKey', input.apiKey],
    ['platform', input.platform || 'openai']
  ])

  if (input.models && input.models.length > 0) {
    params.set('models', input.models.join(','))
  }

  return `http://127.0.0.1:${localPort}/#/settings?${params.toString()}`
}
