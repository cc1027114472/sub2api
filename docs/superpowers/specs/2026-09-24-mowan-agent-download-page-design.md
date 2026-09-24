# 专属智能体客户端（魔丸 Mowan Agent）下载页面设计规格

- **日期**：2026-09-24
- **状态**：已批准 (Approved)
- **目标组件**：前端专属客户端下载中心 (`/download`) 与全站导航集成
- **关联项目**：`sub2api` (Web 前端) 与 `mowan-harness` (客户端安装包及资产来源)

---

## 1. 背景与目标

为了让平台用户更顺畅地使用 Antigravity / DeepSeek 生态下的智能体能力，需要为本站用户打造专属的客户端软件（魔丸 Mowan Agent）下载页面。该客户端具备自动化代码执行、沙箱隔离、多智能体协同循环以及长会话记忆能力。

本功能旨在解决以下痛点：
1. **统一入口**：在顶栏红框显著位置（消息通知与模型广场之间）以及侧边栏提供直观的客户端入口。
2. **便捷分发**：支持直接下载 Windows 安装包（`Mowan-Agent-Setup.exe`）、便携版，并预告 macOS 支持。
3. **闭环体验**：不仅提供下载，而且提供“一键接入配置助手”，已登录用户可一键带入当前站点的 Base URL 与 API Key，无需手动繁琐配置。
4. **灵活托管**：支持本站静态直接下载与外部 CDN / 对象存储链接覆盖。

---

## 2. 架构与路由设计

### 2.1 路由定义
* **路径**：`/download`
* **别名**：`/client-download`
* **访问权限**：`requiresAuth: false`（完全支持未登录与公开分享访问）
* **双模式展示**：
  * **嵌入模式 (`?embedded=1`)**：用户在登录态下从顶栏或侧边栏点击进入，保留后台侧边栏和顶栏，融入主系统。
  * **独立模式**：未登录用户或外链直达时，作为独立高质感 Landing Page 展示，顶部提供登录/注册及返回按钮（对标现有 `/model-plaza` 机制）。

### 2.2 导航入口设计
1. **顶栏 AppHeader**：
   * **位置**：位于消息铃铛 (`AnnouncementBell`) 与模型广场 (`router-link to="/model-plaza"`) 之间。
   * **展示形式**：带下载/终端图标的导航按钮，文案为“专属客户端”或“魔丸 Agent”。
   * **交互**：点击带 `?embedded=1` 跳转至 `/download`。
2. **左侧侧边栏 AppSidebar**：
   * **位置**：常规功能菜单中，提供“专属客户端”菜单项。
3. **未登录主页 HomeView**：
   * 顶部导航栏增加“客户端下载”链接，便于新用户直接了解并下载。

---

## 3. 界面布局与模块交互

### 3.1 Hero 视觉区
* **品牌展示**：居中展示魔丸（Mowan Agent）Logo，辅以科技感微动效光晕。
* **标语**：
  * 主标题：“魔丸 (Mowan Agent) —— 专属智能体客户端”
  * 副标题：“深度适配本站 API 与模型体系，内置本地安全沙箱、多智能体协同与全自动工具链。”
* **系统自适应主按钮**：
  * 浏览器环境嗅探：判断访客操作系统是否为 Windows；
  * 若是 Windows：首选大按钮高亮“下载 Windows 安装包 (.exe)”，副按钮“下载免安装便携版 (.zip)”；
  * 显示版本号（如 `v1.0.0`）、文件大小、环境内置说明（自带轻量 Node.js 运行时，免配置环境变量）。

### 3.2 多平台下载卡片矩阵
* **Windows 卡片（推荐/主推）**：
  * 安装向导版：`Mowan-Agent-Setup.exe`，支持一键创建桌面快捷方式与全局注册。
  * 便携免安装版：解压即用。
  * 支持系统：Windows 10 / 11 64位。
* **macOS 卡片**：
  * Apple Silicon（M 系列）与 Intel 架构。
  * 提供敬请期待或源码构建指引。
* **Linux & CLI 卡片**：
  * 提供 Docker 运行方式与源码/开发者启动命令。

### 3.3 快速接入配置助手（Quick Setup Assistant）
用户安装后需配置 API 地址与密钥，本模块实现无缝联动：
* **已登录用户**：
  * **接口 Base URL**：显示当前站点地址（如 `https://域名/v1`），支持一键复制。
  * **Token 密钥联动**：自动拉取用户已创建的可用 API Token 列表供下拉选择，同时提供“新建 Token”入口。
  * **快速配置生成器**：
    * 提供“配置文件代码块 (config.yaml / json)”与“环境变量命令 (export / set)”，支持一键复制；
    * 3 步简要指引：① 安装并打开客户端 ➔ ② 粘贴当前配置 ➔ ③ 开启智能体对话与任务。
* **未登录用户**：
  * 提示“登录账户后，可在此一键生成专属接入配置，直接复制至客户端秒级生效”，并提供登录按钮。

### 3.4 核心能力与特性介绍（4大亮点）
1. **多智能体深度协同**：支持 Ralph 深度任务循环与子智能体集群调度，自动分解复杂目标。
2. **真实工具链与本地沙箱**：内置安全沙箱模式，支持 PowerShell 终端、文件读写、网页抓取与代码执行。
3. **本站模型无缝直连**：与本站高并发渠道与模型（Antigravity / Gemini / DeepSeek 等）深度联动，低延迟极速响应。
4. **会话持久化与状态机**：跨会话目标持久追踪，支持断点恢复与多分支任务。

---

## 4. 资源管理与可维护性设计

### 4.1 静态资源
* 将 `D:\GOWorks\mowan-harness` 中的魔丸 Logo（`mowan-logo-128.png`、`mowan-logo-square.png`、`mowan.ico`）复制到 `frontend/public/assets/agent/`。
* 默认安装包存放在 `frontend/public/downloads/Mowan-Agent-Setup.exe`。

### 4.2 配置解耦 (`frontend/src/constants/agentDownload.ts`)
```typescript
export interface AgentDownloadConfig {
  version: string
  releaseDate: string
  windows: {
    installerUrl: string
    portableUrl?: string
    size: string
    sha256?: string
  }
  mac: {
    armUrl?: string
    intelUrl?: string
    size?: string
    status: 'available' | 'coming_soon'
  }
  linux: {
    cliUrl?: string
    status: 'available' | 'coming_soon'
  }
}
```
通过该配置常量集中管理下载链接，后续支持通过后台配置或环境变量覆盖，无需侵入页面代码。

### 4.3 国际化多语言
在 `frontend/src/i18n/locales/zh/` 和 `en/` 中新增 `agentDownload` 相关文案，包括标题、副标题、下载按钮、系统要求、安装指引和特性说明等。

---

## 5. 测试与验证策略

1. **路由与权限验证**：
   * 未登录状态访问 `/download` 正常展示独立页面；
   * 登录状态下点击顶栏/侧边栏入口，以嵌入模式正常展示。
2. **交互验证**：
   * 点击下载按钮正常触发安装包下载；
   * 已登录状态下拉取 Token 并一键复制配置，未登录提示登录；
   * 明亮/暗黑主题切换样式正常无错位；
   * 移动端与桌面端自适应布局测试。
