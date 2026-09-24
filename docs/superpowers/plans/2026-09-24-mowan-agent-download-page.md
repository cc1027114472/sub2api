# 专属智能体客户端（魔丸 Mowan Agent）下载页面 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 为 sub2api 平台构建专属的智能体客户端（魔丸 Mowan Agent）下载落地页与全站导航入口，支持安装包下载、系统探测、多平台矩阵与已登录用户 API Key 快速配置联动。

**Architecture:** 前端独立视图 (`DownloadClientView.vue`) 支持独立模式和后台嵌入模式 (`?embedded=1`)；顶栏 AppHeader、侧边栏 AppSidebar 与公开首页 HomeView 新增入口；常量配置模块 (`agentDownload.ts`) 统一管理下载地址与版本元数据；配置助手组件拉取用户 Token 实现一键复制配置。

**Tech Stack:** Vue 3 (Composition API, `<script setup>`), TypeScript, Tailwind CSS, Vue Router, Pinia, vue-i18n, Vitest.

---

### File Structure Map

- **New Files**:
  - `frontend/src/constants/agentDownload.ts`: 下载配置常量与元数据定义
  - `frontend/src/constants/__tests__/agentDownload.spec.ts`: 下载配置单元测试
  - `frontend/src/components/download/DownloadNavBar.vue`: 独立访问模式顶栏组件
  - `frontend/src/components/download/AgentDownloadContent.vue`: 下载页面核心内容（Hero、平台矩阵、配置助手、特性展示）
  - `frontend/src/components/download/__tests__/AgentDownloadContent.spec.ts`: 下载核心组件单元测试
  - `frontend/src/views/DownloadClientView.vue`: 下载主视图（支持内嵌与独立双模式）
  - `frontend/public/assets/agent/*`: 魔丸 Logo、图标等静态资产
  - `frontend/public/downloads/Mowan-Agent-Setup.exe`: 客户端安装包文件
- **Modified Files**:
  - `frontend/src/router/index.ts`: 注册 `/download` 与 `/client-download` 路由
  - `frontend/src/components/layout/AppHeader.vue`: 顶栏红框位置新增客户端下载入口
  - `frontend/src/components/layout/AppSidebar.vue`: 侧边栏菜单新增专属客户端入口
  - `frontend/src/views/HomeView.vue`: 官网公开主页新增下载入口
  - `frontend/src/i18n/locales/zh/common.ts` & `dashboard.ts`: 简体中文文案
  - `frontend/src/i18n/locales/en/common.ts` & `dashboard.ts`: 英文国际化文案

---

### Task 1: 静态资产与配置常量模块 (Static Assets & Constants)

**Files:**
- Copy: `D:\GOWorks\mowan-harness\mowan-logo-*.png`, `mowan.ico` -> `frontend/public/assets/agent/`
- Create: `frontend/src/constants/agentDownload.ts`
- Test: `frontend/src/constants/__tests__/agentDownload.spec.ts`

- [ ] **Step 1: 创建静态资源目录并复制 Logo 图标**

创建 `frontend/public/assets/agent` 目录并从 `D:\GOWorks\mowan-harness` 复制品牌图片素材。

```powershell
New-Item -ItemType Directory -Force -Path "frontend/public/assets/agent"
Copy-Item "D:\GOWorks\mowan-harness\mowan-logo-*.png" "frontend/public/assets/agent\" -Force
Copy-Item "D:\GOWorks\mowan-harness\mowan.ico" "frontend/public/assets/agent\" -Force
```

- [ ] **Step 2: 编写常量模块单元测试**

创建 `frontend/src/constants/__tests__/agentDownload.spec.ts`：

```typescript
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
```

- [ ] **Step 3: 运行测试验证失败**

运行：`pnpm --dir frontend test frontend/src/constants/__tests__/agentDownload.spec.ts`
预期：FAIL（模块尚未定义）

- [ ] **Step 4: 实现常量模块**

创建 `frontend/src/constants/agentDownload.ts`：

```typescript
/**
 * Configuration and metadata for Mowan Agent client downloads.
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
  releaseDate: '2026-09-24',
  windows: {
    installerUrl: '/downloads/Mowan-Agent-Setup.exe',
    portableUrl: '/downloads/Mowan-Agent-portable.zip',
    size: '171.5 MB',
    releaseDate: '2026-09-24',
  },
  mac: {
    armUrl: '',
    intelUrl: '',
    size: '160 MB',
    status: 'coming_soon',
  },
  linux: {
    cliUrl: 'https://github.com',
    dockerImage: 'mowan-agent:latest',
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
```

- [ ] **Step 5: 重新运行测试验证通过**

运行：`pnpm --dir frontend test frontend/src/constants/__tests__/agentDownload.spec.ts`
预期：PASS

- [ ] **Step 6: 提交代码**

```bash
git add frontend/public/assets/agent/ frontend/src/constants/agentDownload.ts frontend/src/constants/__tests__/agentDownload.spec.ts
git commit -m "feat(download): add agent download constants and brand assets"
```

---

### Task 2: 国际化语言包扩展 (i18n Localization)

**Files:**
- Modify: `frontend/src/i18n/locales/zh/common.ts`
- Modify: `frontend/src/i18n/locales/zh/dashboard.ts`
- Modify: `frontend/src/i18n/locales/en/common.ts`
- Modify: `frontend/src/i18n/locales/en/dashboard.ts`

- [ ] **Step 1: 在 zh 语言包添加文案**

在 `frontend/src/i18n/locales/zh/common.ts` 的 `nav` 节点中增加：
```typescript
agentDownload: '魔丸 Agent',
agentDownloadLong: '专属客户端下载',
```

在 `frontend/src/i18n/locales/zh/dashboard.ts` 中增加 `agentDownload` 命名空间：
```typescript
agentDownload: {
  title: '魔丸 (Mowan Agent) 专属客户端',
  subtitle: '深度适配本站 API 与模型体系，内置本地安全沙箱、多智能体协同循环与全自动工程工具链。',
  badge: '官方专属客户端',
  version: '最新版本',
  releaseDate: '发布日期',
  fileSize: '安装包大小',
  envIncluded: '自带独立运行环境，无需配置 Node.js 或系统依赖，解压或安装即可极速运行。',
  downloadPrimary: '立即下载 Windows 安装包 (.exe)',
  downloadPortable: '下载免安装便携版 (.zip)',
  detectedPlatform: '已智能匹配当前系统：',
  manualSelect: '或手动选择其他系统版本',
  windowsTitle: 'Windows 客户端',
  windowsDesc: '适用于 Windows 10 / 11 64位操作系统，支持一键安装与桌面快捷方式。',
  macTitle: 'macOS 客户端',
  macDesc: '支持 Apple Silicon (M系列) 与 Intel 架构芯片。',
  macComingSoon: '即将推出 / 适配中',
  linuxTitle: 'Linux / 极客模式',
  linuxDesc: '支持 Docker 部署与命令行交互，适合无头服务器与开发者环境。',
  quickStartTitle: '一键接入配置助手',
  quickStartDesc: '客户端安装完成后，只需填入本站 API 地址与您的专属密钥，即可畅享智能体极速调用。',
  step1: '1. 安装并启动客户端',
  step2: '2. 复制下方本站接口地址与 API Key',
  step3: '3. 粘贴至客户端设置，立即开启多智能体任务',
  apiBaseUrl: '接口 Base URL',
  apiKeySelect: '选择您的 API Key',
  noKeyHint: '暂无可用的 API Key，点击前往创建',
  createKeyBtn: '前往密钥管理',
  loginPrompt: '登录当前站点后，可直接在此处一键获取您的专属对接配置，秒级完成连通。',
  loginNow: '立即登录',
  copySuccess: '配置已复制到剪贴板！',
  features: {
    title: '为什么选择魔丸专属智能体？',
    loopTitle: '自主多智能体循环',
    loopDesc: '原生支持 Ralph 深度任务循环与子智能体集群协同，告别单轮对话，自主推理并完成复杂开发目标。',
    sandboxTitle: '本地安全沙箱',
    sandboxDesc: '内置安全的 PowerShell 终端执行、文件读写、网页爬取与代码运行能力，在受控环境中高效交付。',
    modelsTitle: '本站模型极速直连',
    modelsDesc: '无缝对接本站 Antigravity、Gemini、DeepSeek 等高并发优质渠道与倍率，低延迟极速响应。',
    memoryTitle: '跨会话状态持久化',
    memoryDesc: '拥有任务进度持久化记忆，支持断点恢复、长上下文工程规划与多轮会话自愈。',
  },
  nav: {
    backToDashboard: '返回控制台',
    login: '登录系统',
  }
}
```

- [ ] **Step 2: 在 en 语言包添加对应英文**

同步在 `frontend/src/i18n/locales/en/common.ts` 与 `en/dashboard.ts` 中补充英文对应定义。

- [ ] **Step 3: 运行 typecheck 检查类型完整性**

运行：`pnpm --dir frontend typecheck`
预期：类型检查通过无报错。

- [ ] **Step 4: 提交代码**

```bash
git add frontend/src/i18n/locales/
git commit -m "feat(i18n): add agent download translations in zh and en"
```

---

### Task 3: 下载页面核心组件开发 (Download Content & Components)

**Files:**
- Create: `frontend/src/components/download/DownloadNavBar.vue`
- Create: `frontend/src/components/download/AgentDownloadContent.vue`
- Test: `frontend/src/components/download/__tests__/AgentDownloadContent.spec.ts`

- [ ] **Step 1: 编写核心组件测试用例**

创建 `frontend/src/components/download/__tests__/AgentDownloadContent.spec.ts`：

```typescript
import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import AgentDownloadContent from '../AgentDownloadContent.vue'
import { createTestingPinia } from '@pinia/testing'

describe('AgentDownloadContent', () => {
  it('renders brand title and download buttons', () => {
    const wrapper = mount(AgentDownloadContent, {
      global: {
        plugins: [createTestingPinia({ createSpy: vi.fn })],
        mocks: {
          t: (key: string) => key,
        },
        stubs: {
          Icon: true,
          RouterLink: true,
        },
      },
    })

    expect(wrapper.text()).toContain('agentDownload.title')
    expect(wrapper.text()).toContain('agentDownload.downloadPrimary')
  })
})
```

- [ ] **Step 2: 运行测试验证失败**

运行：`pnpm --dir frontend test frontend/src/components/download/__tests__/AgentDownloadContent.spec.ts`
预期：FAIL（组件不存在）

- [ ] **Step 3: 实现独立顶栏组件 `DownloadNavBar.vue`**

创建 `frontend/src/components/download/DownloadNavBar.vue`，实现与 `PlazaNavBar.vue` 风格一致的高质感顶栏（展示站点 Logo、名称、主题切换、登录/返回控制台按钮）。

- [ ] **Step 4: 实现下载页面主体组件 `AgentDownloadContent.vue`**

创建 `frontend/src/components/download/AgentDownloadContent.vue`：
- Hero 模块：带有渐变背景、魔丸 Logo 图片、智能操作系统检测、主下载按钮与版本信息。
- 快速接入配置助手：
  - 读取当前 `window.location.origin + '/v1'` 作为 Base URL；
  - 若已登录，调用 `keysAPI.list(1, 20, { status: 'active' })` 异步获取用户的有效密钥；
  - 提供单选/下拉选择并一键生成完整 YAML / ENV 配置；
  - 未登录状态友好提示登录。
- 多平台卡片列表：Windows、macOS、Linux 卡片。
- 四大核心特性展示卡片。

- [ ] **Step 5: 运行测试验证通过**

运行：`pnpm --dir frontend test frontend/src/components/download/__tests__/AgentDownloadContent.spec.ts`
预期：PASS

- [ ] **Step 6: 提交代码**

```bash
git add frontend/src/components/download/
git commit -m "feat(download): implement AgentDownloadContent and DownloadNavBar components"
```

---

### Task 4: 视图封装与路由挂载 (View & Router Integration)

**Files:**
- Create: `frontend/src/views/DownloadClientView.vue`
- Modify: `frontend/src/router/index.ts`
- Test: `frontend/src/views/__tests__/DownloadClientView.spec.ts`

- [ ] **Step 1: 编写视图包装组件测试**

创建 `frontend/src/views/__tests__/DownloadClientView.spec.ts`，验证在 `isEmbedded` 为真时挂载 `AppLayout`，为假时挂载 `DownloadNavBar`。

- [ ] **Step 2: 实现主视图组件 `DownloadClientView.vue`**

创建 `frontend/src/views/DownloadClientView.vue`：
```vue
<template>
  <AppLayout v-if="isEmbedded">
    <AgentDownloadContent embedded />
  </AppLayout>
  <div v-else class="min-h-screen bg-gray-50 dark:bg-dark-950">
    <DownloadNavBar />
    <main class="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
      <AgentDownloadContent />
    </main>
  </div>
</template>
```

- [ ] **Step 3: 在 `frontend/src/router/index.ts` 注册路由**

在 Public Routes 区域添加：
```typescript
{
  path: '/download',
  alias: '/client-download',
  name: 'AgentDownload',
  component: () => import('@/views/DownloadClientView.vue'),
  meta: {
    requiresAuth: false,
    title: 'Download Agent',
    titleKey: 'agentDownload.title'
  }
}
```

- [ ] **Step 4: 运行测试验证通过**

运行：`pnpm --dir frontend test frontend/src/views/__tests__/DownloadClientView.spec.ts`
预期：PASS

- [ ] **Step 5: 提交代码**

```bash
git add frontend/src/views/DownloadClientView.vue frontend/src/router/index.ts frontend/src/views/__tests__/DownloadClientView.spec.ts
git commit -m "feat(router): add /download and /client-download routes with DownloadClientView"
```

---

### Task 5: 导航入口集成 (AppHeader, AppSidebar, HomeView)

**Files:**
- Modify: `frontend/src/components/layout/AppHeader.vue:40-52` (红框位置：铃铛与模型广场之间)
- Modify: `frontend/src/components/layout/AppSidebar.vue:720-745` (左侧菜单项)
- Modify: `frontend/src/views/HomeView.vue:42-53` (公开主页顶栏)

- [ ] **Step 1: 在 AppHeader.vue 红框处添加客户端下载按钮**

在 `frontend/src/components/layout/AppHeader.vue` 的 `AnnouncementBell` 与 `router-link to="/model-plaza"` 之间添加：
```vue
<!-- Agent Client Download Entry -->
<router-link
  v-if="user"
  :to="{ path: '/download', query: { embedded: '1' } }"
  class="hidden items-center gap-1.5 rounded-lg px-2.5 py-1.5 text-sm font-medium text-gray-600 transition-colors hover:bg-gray-100 hover:text-gray-900 dark:text-dark-400 dark:hover:bg-dark-800 dark:hover:text-white sm:flex"
  :title="t('nav.agentDownload')"
>
  <Icon name="download" size="sm" />
  <span class="hidden sm:inline">{{ t('nav.agentDownload') }}</span>
</router-link>
```

- [ ] **Step 2: 在 AppSidebar.vue 菜单添加专属客户端项**

在 `buildSelfNavItems` 中加入：
```typescript
{
  path: '/download',
  label: t('nav.agentDownload'),
  icon: DownloadIcon, // 或者复用现有图标
}
```

- [ ] **Step 3: 在 HomeView.vue 顶部导航添加下载链接**

在模型广场链接旁边加入：
```vue
<router-link
  to="/download"
  class="flex h-10 shrink-0 items-center gap-1.5 rounded-lg px-2.5 text-sm font-medium text-gray-500 hover:bg-gray-100 hover:text-gray-700 dark:text-dark-400 dark:hover:bg-dark-800 dark:hover:text-white"
  :title="t('nav.agentDownload')"
>
  <Icon name="download" size="md" />
  <span class="hidden sm:inline">{{ t('nav.agentDownload') }}</span>
</router-link>
```

- [ ] **Step 4: 运行前端自动化测试**

运行：`pnpm --dir frontend test`
预期：所有导航与布局测试全部 PASS。

- [ ] **Step 5: 提交代码**

```bash
git add frontend/src/components/layout/ frontend/src/views/HomeView.vue
git commit -m "feat(nav): add agent download entry points to header, sidebar and home page"
```

---

### Task 6: 安装包放置与端到端系统验证 (Deployment & Verification)

**Files:**
- Create: `frontend/public/downloads/Mowan-Agent-Setup.exe` (复制自 `D:\GOWorks\mowan-harness\dist\Mowan-Agent-Setup.exe`)

- [ ] **Step 1: 复制安装包到静态资源目录**

```powershell
New-Item -ItemType Directory -Force -Path "frontend/public/downloads"
Copy-Item "D:\GOWorks\mowan-harness\dist\Mowan-Agent-Setup.exe" "frontend/public/downloads\Mowan-Agent-Setup.exe" -Force
```

- [ ] **Step 2: 执行前端代码类型与完整测试检查**

运行：
```powershell
pnpm --dir frontend typecheck
pnpm --dir frontend test
```
预期：0 errors, 所有测试通过。

- [ ] **Step 3: 验证构建产物**

运行：
```powershell
pnpm --dir frontend build
```
预期：Vite 构建成功，无打包警告与错误。

- [ ] **Step 4: 提交并完成功能集成**

```bash
git add .
git commit -m "feat: complete mowan agent dedicated download page integration"
```
