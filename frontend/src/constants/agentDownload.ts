/**
 * Configuration and metadata for mowan-harness client downloads.
 */

export interface AgentPlatformWindows {
  installerUrl: string
  portableUrl?: string
  size: string
  releaseDate: string
  sha256?: string
}

export interface AgentPlatformMac {
  armUrl?: string
  intelUrl?: string
  size?: string
  status: 'available' | 'coming_soon'
}

export interface AgentPlatformLinux {
  cliUrl?: string
  dockerImage?: string
  status: 'available' | 'coming_soon'
}

export interface AgentDownloadConfig {
  version: string
  releaseDate: string
  windows: AgentPlatformWindows
  mac: AgentPlatformMac
  linux: AgentPlatformLinux
}

export const defaultAgentDownloadConfig: AgentDownloadConfig = {
  version: 'v1.0.0',
  releaseDate: '2026-09-25',
  windows: {
    installerUrl: '/downloads/Mowan-Harness-Setup.exe',
    portableUrl: '/downloads/Mowan-Harness-Setup.exe',
    size: '174.3 MB',
    releaseDate: '2026-09-25',
  },
  mac: {
    armUrl: '',
    intelUrl: '',
    size: '160 MB',
    status: 'coming_soon',
  },
  linux: {
    cliUrl: 'https://github.com',
    dockerImage: 'mowan-harness:latest',
    status: 'coming_soon',
  },
}

export interface ResolvedPlatform {
  isWindows: boolean
  isMac: boolean
  isLinux: boolean
  primaryUrl: string
  platformLabel: string
}

export function resolvePlatformDownload(userAgent: string = ''): ResolvedPlatform {
  const ua = userAgent || (typeof navigator !== 'undefined' ? navigator.userAgent : '')
  const isWindows = /windows|win32|win64/i.test(ua)
  const isMac = /macintosh|mac os x/i.test(ua)
  const isLinux = /linux/i.test(ua) && !/android/i.test(ua)

  let primaryUrl = defaultAgentDownloadConfig.windows.installerUrl
  let platformLabel = 'Windows 64位'

  if (isMac) {
    primaryUrl = defaultAgentDownloadConfig.mac.armUrl || defaultAgentDownloadConfig.windows.installerUrl
    platformLabel = 'macOS'
  } else if (isLinux) {
    primaryUrl = defaultAgentDownloadConfig.linux.cliUrl || defaultAgentDownloadConfig.windows.installerUrl
    platformLabel = 'Linux / Docker'
  }

  return {
    isWindows,
    isMac,
    isLinux,
    primaryUrl,
    platformLabel,
  }
}
