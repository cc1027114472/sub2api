import { describe, it, expect } from 'vitest'
import { defaultAgentDownloadConfig, resolvePlatformDownload } from '../agentDownload'

describe('agentDownload config', () => {
  it('provides default Windows setup installer configuration', () => {
    expect(defaultAgentDownloadConfig.version).toBeDefined()
    expect(defaultAgentDownloadConfig.windows.installerUrl).toContain('Mowan-Agent-Setup.exe')
    expect(defaultAgentDownloadConfig.windows.size).toBeDefined()
  })

  it('resolves correct platform info for windows user agent', () => {
    const info = resolvePlatformDownload('Mozilla/5.0 (Windows NT 10.0; Win64; x64)')
    expect(info.isWindows).toBe(true)
    expect(info.isMac).toBe(false)
    expect(info.primaryUrl).toBe(defaultAgentDownloadConfig.windows.installerUrl)
  })

  it('resolves correct platform info for mac user agent', () => {
    const info = resolvePlatformDownload('Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)')
    expect(info.isWindows).toBe(false)
    expect(info.isMac).toBe(true)
  })
})
