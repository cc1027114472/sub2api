<template>
  <!-- Custom Home Content: Full Page Mode -->
  <div v-if="hasHomeContent" class="min-h-screen">
    <!-- iframe mode -->
    <iframe
      v-if="isHomeContentUrl"
      :src="homeContent.trim()"
      class="h-screen w-full border-0"
      allowfullscreen
    ></iframe>
    <!-- HTML mode - SECURITY: homeContent is admin-only setting -->
    <div v-else v-html="homeContent"></div>
  </div>

  <!-- Compact Home Page -->
  <div
    v-else-if="compactHomeEnabled"
    data-testid="compact-home"
    class="flex min-h-screen flex-col bg-gray-50 text-gray-900 dark:bg-dark-950 dark:text-white"
  >
    <header class="border-b border-gray-200/80 px-4 py-4 sm:px-6 dark:border-dark-800">
      <nav class="mx-auto flex max-w-5xl flex-wrap items-center justify-between gap-3 sm:gap-4">
        <div class="flex min-w-0 flex-1 items-center gap-3">
          <img
            :src="siteLogo || '/logo.svg'"
            alt="Logo"
            class="h-9 w-9 shrink-0 rounded-lg object-contain shadow-sm"
          />
          <span class="min-w-0 truncate text-base font-semibold tracking-tight">{{ siteName }}</span>
        </div>
        <div class="flex max-w-full shrink-0 flex-wrap items-center justify-end gap-2">
          <LocaleSwitcher />
          <a
            v-if="docUrl"
            :href="docUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg text-gray-500 hover:bg-gray-100 dark:text-dark-400 dark:hover:bg-dark-800"
            :title="t('home.viewDocs')"
          >
            <Icon name="book" size="md" />
          </a>
          <router-link
            v-if="showModelPlazaEntry"
            to="/model-plaza"
            class="flex h-10 shrink-0 items-center gap-1.5 rounded-lg px-2.5 text-sm font-medium text-gray-500 hover:bg-gray-100 hover:text-gray-700 dark:text-dark-400 dark:hover:bg-dark-800 dark:hover:text-white"
            :title="t('nav.modelPlaza')"
          >
            <Icon name="grid" size="md" />
            <span class="hidden sm:inline">{{ t('nav.modelPlaza') }}</span>
          </router-link>
          <router-link
            v-if="showAgentDownloadEntry"
            to="/download"
            class="flex h-10 shrink-0 items-center gap-1.5 rounded-lg px-2.5 text-sm font-medium text-gray-500 hover:bg-gray-100 hover:text-gray-700 dark:text-dark-400 dark:hover:bg-dark-800 dark:hover:text-white"
            :title="t('nav.agentDownload')"
          >
            <Icon name="download" size="md" />
            <span class="hidden sm:inline">{{ t('nav.agentDownload') }}</span>
          </router-link>
          <button
            class="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg text-gray-500 hover:bg-gray-100 dark:text-dark-400 dark:hover:bg-dark-800"
            :title="isDark ? t('home.switchToLight') : t('home.switchToDark')"
            @click="toggleTheme"
          >
            <Icon v-if="isDark" name="sun" size="md" />
            <Icon v-else name="moon" size="md" />
          </button>
          <router-link
            :to="isAuthenticated ? dashboardPath : '/login'"
            class="inline-flex min-h-10 shrink-0 items-center justify-center rounded-lg bg-gray-900 px-4 py-2 text-sm font-medium text-white hover:bg-gray-800 dark:bg-white dark:text-gray-900 dark:hover:bg-gray-200 shadow-sm"
          >
            {{ isAuthenticated ? t('home.dashboard') : t('home.login') }}
          </router-link>
        </div>
      </nav>
    </header>

    <main class="flex min-w-0 flex-1 items-center justify-center px-4 py-16 sm:px-6">
      <div class="min-w-0 max-w-2xl text-center">
        <img
          :src="siteLogo || '/logo.svg'"
          alt="Logo"
          class="mx-auto mb-6 h-20 w-20 rounded-2xl object-contain shadow-md"
        />
        <h1 class="[overflow-wrap:anywhere] text-3xl font-bold tracking-tight md:text-4xl">{{ siteName }}</h1>
        <p class="mt-4 whitespace-pre-wrap [overflow-wrap:anywhere] text-base text-gray-600 dark:text-dark-300">{{ siteSubtitle }}</p>
        <router-link
          :to="isAuthenticated ? dashboardPath : '/login'"
          class="mt-8 inline-flex min-h-10 items-center justify-center rounded-xl bg-primary-600 px-6 py-2.5 text-sm font-medium text-white shadow-md shadow-primary-500/25 hover:bg-primary-500"
        >
          {{ isAuthenticated ? t('home.goToDashboard') : t('home.login') }}
        </router-link>
      </div>
    </main>

    <footer class="min-w-0 border-t border-gray-200/80 px-4 py-5 text-center text-sm text-gray-500 sm:px-6 dark:border-dark-800 dark:text-dark-400">
      &copy; {{ currentYear }} {{ siteName }}
    </footer>
  </div>

  <!-- Linear/Vercel Style Full Hero Home Page -->
  <div
    v-else
    class="relative flex min-h-screen flex-col overflow-x-hidden bg-[#fafafa] text-gray-900 antialiased selection:bg-primary-500/20 selection:text-primary-900 dark:bg-[#06090e] dark:text-gray-100 dark:selection:text-primary-200"
  >
    <!-- Background: 精致网格、顶部极光与微光流带 -->
    <div class="pointer-events-none fixed inset-0 z-0 overflow-hidden">
      <!-- 顶部光晕高光 -->
      <div
        class="absolute -top-[300px] left-1/2 h-[600px] w-[1000px] -translate-x-1/2 rounded-full bg-gradient-to-b from-primary-500/20 via-cyan-500/10 to-transparent blur-[120px] dark:from-primary-600/25 dark:via-cyan-600/15"
      ></div>
      <!-- 极客网格层 -->
      <div
        class="absolute inset-0 bg-[linear-gradient(to_right,rgba(99,102,241,0.04)_1px,transparent_1px),linear-gradient(to_bottom,rgba(99,102,241,0.04)_1px,transparent_1px)] bg-[size:44px_44px] [mask-image:radial-gradient(ellipse_60%_50%_at_50%_0%,#000_70%,transparent_100%)] dark:bg-[linear-gradient(to_right,rgba(255,255,255,0.03)_1px,transparent_1px),linear-gradient(to_bottom,rgba(255,255,255,0.03)_1px,transparent_1px)]"
      ></div>
      <!-- 顶部 1px 细微光带 -->
      <div class="absolute left-0 right-0 top-0 h-[1px] bg-gradient-to-r from-transparent via-primary-500/50 to-transparent"></div>
    </div>

    <!-- Floating Pill Navigation (悬浮拟态导航栏) -->
    <header class="sticky top-4 z-50 mx-auto w-full max-w-6xl px-4 sm:px-6">
      <nav
        class="flex items-center justify-between rounded-2xl border border-gray-200/80 bg-white/70 px-4 py-2.5 shadow-sm backdrop-blur-xl transition-all duration-300 dark:border-white/10 dark:bg-dark-900/70 dark:shadow-2xl dark:shadow-black/50"
      >
        <!-- Logo & Site Name -->
        <router-link to="/home" class="group flex items-center gap-3">
          <div class="relative flex h-9 w-9 items-center justify-center overflow-hidden rounded-xl border border-gray-200/60 bg-white p-1 shadow-sm transition-transform duration-200 group-hover:scale-105 dark:border-white/10 dark:bg-dark-800">
            <img :src="siteLogo || '/logo.svg'" alt="Logo" class="h-full w-full object-contain" />
          </div>
          <div class="flex flex-col">
            <span class="text-sm font-bold tracking-tight text-gray-900 dark:text-white">{{ siteName }}</span>
            <span class="text-[10px] font-medium text-primary-600 dark:text-cyan-400">AI GATEWAY</span>
          </div>
        </router-link>

        <!-- Center Nav Links (Desktop) -->
        <div class="hidden items-center gap-1 md:flex">
          <router-link
            v-if="showModelPlazaEntry"
            to="/model-plaza"
            class="flex items-center gap-1.5 rounded-lg px-3 py-1.5 text-xs font-medium text-gray-600 transition-colors hover:bg-gray-100 hover:text-gray-900 dark:text-dark-300 dark:hover:bg-white/5 dark:hover:text-white"
          >
            <Icon name="grid" size="xs" />
            <span>{{ t('nav.modelPlaza') }}</span>
          </router-link>

          <router-link
            v-if="showAgentDownloadEntry"
            to="/download"
            class="flex items-center gap-1.5 rounded-lg px-3 py-1.5 text-xs font-medium text-gray-600 transition-colors hover:bg-gray-100 hover:text-gray-900 dark:text-dark-300 dark:hover:bg-white/5 dark:hover:text-white"
          >
            <Icon name="download" size="xs" />
            <span>{{ t('nav.agentDownload') }}</span>
          </router-link>

          <router-link
            to="/key-usage"
            class="flex items-center gap-1.5 rounded-lg px-3 py-1.5 text-xs font-medium text-gray-600 transition-colors hover:bg-gray-100 hover:text-gray-900 dark:text-dark-300 dark:hover:bg-white/5 dark:hover:text-white"
            title="免登录用量速查"
          >
            <svg class="h-3.5 w-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
              <path d="M7 11V7a5 5 0 0 1 10 0v4"/>
            </svg>
            <span>{{ t('keyUsage.title') || '用量速查' }}</span>
          </router-link>

          <a
            v-if="docUrl"
            :href="docUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="flex items-center gap-1.5 rounded-lg px-3 py-1.5 text-xs font-medium text-gray-600 transition-colors hover:bg-gray-100 hover:text-gray-900 dark:text-dark-300 dark:hover:bg-white/5 dark:hover:text-white"
          >
            <Icon name="book" size="xs" />
            <span>{{ t('home.viewDocs') }}</span>
          </a>
        </div>

        <!-- Right Action Controls -->
        <div class="flex items-center gap-2">
          <LocaleSwitcher />

          <!-- Theme Toggle -->
          <button
            @click="toggleTheme"
            class="flex h-8 w-8 items-center justify-center rounded-lg text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-900 dark:text-dark-400 dark:hover:bg-white/5 dark:hover:text-white"
            :title="isDark ? t('home.switchToLight') : t('home.switchToDark')"
          >
            <Icon v-if="isDark" name="sun" size="sm" />
            <Icon v-else name="moon" size="sm" />
          </button>

          <!-- Login / Dashboard Action -->
          <router-link
            v-if="isAuthenticated"
            :to="dashboardPath"
            class="inline-flex items-center gap-2 rounded-xl bg-gray-900 px-3.5 py-1.5 text-xs font-medium text-white shadow-sm transition-all hover:bg-gray-800 dark:bg-white dark:text-gray-900 dark:hover:bg-gray-100"
          >
            <span class="flex h-4 w-4 items-center justify-center rounded-full bg-primary-500 text-[9px] font-bold text-white">
              {{ userInitial }}
            </span>
            <span>{{ t('home.dashboard') }}</span>
          </router-link>

          <router-link
            v-else
            to="/login"
            class="group relative inline-flex items-center justify-center overflow-hidden rounded-xl bg-gradient-to-r from-primary-600 to-indigo-600 px-4 py-1.5 text-xs font-semibold text-white shadow-md shadow-primary-500/20 transition-all duration-200 hover:from-primary-500 hover:to-indigo-500 hover:shadow-primary-500/30 active:scale-[0.98]"
          >
            <span>{{ t('home.login') }}</span>
            <svg class="ml-1 h-3.5 w-3.5 transition-transform group-hover:translate-x-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </router-link>
        </div>
      </nav>
    </header>

    <!-- Main Hero Section -->
    <main class="relative z-10 flex-1 px-4 sm:px-6">
      <div class="mx-auto max-w-6xl pt-12 pb-16 lg:pt-20 lg:pb-24">
        
        <!-- Hero Grid: Left Typography + Right Interactive Terminal -->
        <div class="grid items-center gap-12 lg:grid-cols-12 lg:gap-8">
          
          <!-- Left Content -->
          <div class="flex flex-col items-center text-center lg:col-span-6 lg:items-start lg:text-left">
            
            <!-- Tech Badge -->
            <div class="inline-flex items-center gap-2 rounded-full border border-primary-500/20 bg-primary-50/50 px-3 py-1 text-xs font-medium text-primary-700 backdrop-blur-md dark:border-primary-500/30 dark:bg-primary-950/30 dark:text-primary-300">
              <span class="flex h-2 w-2 rounded-full bg-emerald-500 animate-pulse"></span>
              <span>高性能 AI 算力分发 & 智能路由网关</span>
            </div>

            <!-- Big Impact Title -->
            <h1 class="mt-6 text-4xl font-extrabold tracking-tight text-gray-900 sm:text-5xl lg:text-6xl dark:text-white">
              <span class="block">聚合全球顶尖算力</span>
              <span class="mt-1 block bg-gradient-to-r from-primary-600 via-indigo-600 to-cyan-500 bg-clip-text text-transparent dark:from-primary-400 dark:via-indigo-300 dark:to-cyan-300">
                {{ t('home.heroSubtitle') || '一个密钥，畅用全系 AI' }}
              </span>
            </h1>

            <!-- Subtitle -->
            <p class="mt-5 max-w-xl text-base leading-relaxed text-gray-600 dark:text-dark-300 sm:text-lg">
              {{ siteSubtitle || t('home.heroDescription') || '无需分别订阅，一站式原生接入 Claude Code、Codex CLI、OpenAI、Gemini 与 DeepSeek 等模型，毫秒级粘性调度与实时精确计费。' }}
            </p>

            <!-- Actions Row -->
            <div class="mt-8 flex w-full flex-col gap-3.5 sm:w-auto sm:flex-row">
              <router-link
                :to="isAuthenticated ? dashboardPath : '/register'"
                class="group relative inline-flex items-center justify-center gap-2 rounded-xl bg-gradient-to-r from-primary-600 to-indigo-600 px-6 py-3.5 text-sm font-semibold text-white shadow-lg shadow-primary-500/25 transition-all duration-200 hover:from-primary-500 hover:to-indigo-500 hover:shadow-xl hover:shadow-primary-500/35 active:scale-[0.98]"
              >
                <span>{{ isAuthenticated ? t('home.goToDashboard') : (t('home.getStarted') || '免费获取 API Key') }}</span>
                <svg class="h-4 w-4 transition-transform group-hover:translate-x-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3" />
                </svg>
              </router-link>

              <router-link
                v-if="showModelPlazaEntry"
                to="/model-plaza"
                class="inline-flex items-center justify-center gap-2 rounded-xl border border-gray-200 bg-white/80 px-6 py-3.5 text-sm font-medium text-gray-700 shadow-sm backdrop-blur-sm transition-all hover:bg-gray-50 hover:text-gray-900 dark:border-white/10 dark:bg-dark-800/80 dark:text-dark-200 dark:hover:bg-dark-700 dark:hover:text-white"
              >
                <Icon name="grid" size="sm" />
                <span>{{ t('nav.modelPlaza') }} & 价格对照</span>
              </router-link>

              <router-link
                v-else
                to="/download"
                class="inline-flex items-center justify-center gap-2 rounded-xl border border-gray-200 bg-white/80 px-6 py-3.5 text-sm font-medium text-gray-700 shadow-sm backdrop-blur-sm transition-all hover:bg-gray-50 hover:text-gray-900 dark:border-white/10 dark:bg-dark-800/80 dark:text-dark-200 dark:hover:bg-dark-700 dark:hover:text-white"
              >
                <Icon name="download" size="sm" />
                <span>{{ t('nav.agentDownload') }}</span>
              </router-link>
            </div>

            <!-- Trust Highlights (Pills) -->
            <div class="mt-8 flex flex-wrap items-center justify-center gap-4 text-xs font-medium text-gray-500 dark:text-dark-400 lg:justify-start">
              <div class="flex items-center gap-1.5">
                <Icon name="checkCircle" size="xs" class="text-emerald-500" />
                <span>{{ t('home.tags.subscriptionToApi') || '订阅转 API' }}</span>
              </div>
              <span class="text-gray-300 dark:text-dark-600">·</span>
              <div class="flex items-center gap-1.5">
                <Icon name="bolt" size="xs" class="text-primary-500" />
                <span>{{ t('home.tags.stickySession') || '毫秒级粘性调度' }}</span>
              </div>
              <span class="text-gray-300 dark:text-dark-600">·</span>
              <div class="flex items-center gap-1.5">
                <Icon name="chart" size="xs" class="text-cyan-500" />
                <span>{{ t('home.tags.realtimeBilling') || '按量精准 Token 扣费' }}</span>
              </div>
            </div>

          </div>

          <!-- Right Interactive Multi-Protocol Terminal (.terminal-container) -->
          <div class="flex justify-center lg:col-span-6 lg:justify-end">
            <!-- 容器保留 .terminal-container 类名，保障单元测试和样式锚定 -->
            <div class="terminal-container w-full max-w-lg">
              
              <!-- 窗口实体 -->
              <div class="terminal-window w-full rounded-2xl border border-gray-800/80 bg-[#0d121f] text-left shadow-2xl shadow-primary-500/10 dark:border-white/10 dark:bg-[#090d16]">
                
                <!-- Terminal Header with Tabs -->
                <div class="flex items-center justify-between border-b border-white/5 bg-[#121829] px-4 py-3">
                  <!-- macOS Buttons -->
                  <div class="terminal-buttons flex items-center gap-2">
                    <span class="btn-close h-3 w-3 rounded-full bg-rose-500/80"></span>
                    <span class="btn-minimize h-3 w-3 rounded-full bg-amber-500/80"></span>
                    <span class="btn-maximize h-3 w-3 rounded-full bg-emerald-500/80"></span>
                  </div>

                  <!-- Protocol Switcher Tabs -->
                  <div class="flex items-center gap-1 rounded-lg bg-black/30 p-1">
                    <button
                      v-for="tab in terminalTabs"
                      :key="tab.id"
                      @click="activeTerminalTab = tab.id"
                      class="rounded-md px-2.5 py-1 text-[11px] font-mono transition-all"
                      :class="activeTerminalTab === tab.id ? 'bg-primary-600 text-white font-semibold shadow-sm' : 'text-gray-400 hover:text-white'"
                    >
                      {{ tab.label }}
                    </button>
                  </div>

                  <!-- Copy Button -->
                  <button
                    @click="copyActiveCode"
                    class="flex h-7 w-7 items-center justify-center rounded-lg border border-white/5 text-gray-400 transition-colors hover:bg-white/5 hover:text-white"
                    :title="copied ? '已复制' : '复制命令'"
                  >
                    <Icon v-if="copied" name="check" size="xs" class="text-emerald-400" />
                    <Icon v-else name="clipboard" size="xs" />
                  </button>
                </div>

                <!-- Terminal Body (Code Display) -->
                <div class="terminal-body p-4 sm:p-5 font-mono text-xs sm:text-[13px] leading-relaxed">
                  
                  <!-- Tab 1: Claude Code -->
                  <div v-if="activeTerminalTab === 'claude'" class="space-y-3">
                    <div class="code-line line-1 flex items-start gap-2 text-gray-400">
                      <span class="text-emerald-400 font-bold">$</span>
                      <span class="text-gray-300">export ANTHROPIC_BASE_URL="https://api.your-domain.com"</span>
                    </div>
                    <div class="code-line line-2 flex items-start gap-2 text-gray-400">
                      <span class="text-emerald-400 font-bold">$</span>
                      <span class="text-gray-300">export ANTHROPIC_API_KEY="sk-sub2api-xxxxxxxx"</span>
                    </div>
                    <div class="code-line line-3 flex items-start gap-2 text-primary-300">
                      <span class="text-emerald-400 font-bold">$</span>
                      <span class="text-cyan-300">claude</span>
                      <span class="text-gray-400">--model claude-3-7-sonnet-20250219</span>
                    </div>
                    <!-- Output Simulation -->
                    <div class="mt-3 rounded-lg border border-emerald-500/20 bg-emerald-950/20 p-2.5 text-[11px] text-emerald-400">
                      <div class="flex items-center gap-1.5 font-semibold">
                        <span class="h-1.5 w-1.5 rounded-full bg-emerald-400 animate-ping"></span>
                        <span>Gateway Connected: Claude Code Hook Ready (14ms latency)</span>
                      </div>
                    </div>
                  </div>

                  <!-- Tab 2: OpenAI / cURL -->
                  <div v-else-if="activeTerminalTab === 'curl'" class="space-y-2">
                    <div class="code-line line-1 text-gray-300">
                      <span class="text-emerald-400 font-bold">$ </span>
                      <span class="text-cyan-400">curl</span> https://api.your-domain.com/v1/chat/completions \
                    </div>
                    <div class="code-line line-2 pl-4 text-gray-400">
                      -H <span class="text-amber-300">"Authorization: Bearer sk-sub2api-xxxx"</span> \
                    </div>
                    <div class="code-line line-3 pl-4 text-gray-400">
                      -H <span class="text-amber-300">"Content-Type: application/json"</span> \
                    </div>
                    <div class="code-line line-4 pl-4 text-gray-400">
                      -d <span class="text-primary-300">'{"model": "gpt-4o", "messages": [{"role": "user", "content": "Hello"}]}'</span>
                    </div>
                    <div class="mt-3 rounded-lg border border-cyan-500/20 bg-cyan-950/20 p-2.5 text-[11px] text-cyan-300">
                      <span>HTTP/2 200 OK · tokens: { "prompt": 9, "completion": 18, "cost": "$0.00012" }</span>
                    </div>
                  </div>

                  <!-- Tab 3: Codex CLI -->
                  <div v-else-if="activeTerminalTab === 'codex'" class="space-y-3">
                    <div class="code-comment text-gray-500 italic"># ~/.codex/config.toml</div>
                    <div class="code-line line-1 text-gray-300">
                      <span class="text-primary-400">model_provider</span> = <span class="text-amber-300">"sub2api"</span>
                    </div>
                    <div class="code-line line-2 text-gray-300">
                      <span class="text-primary-400">base_url</span> = <span class="text-amber-300">"https://api.your-domain.com/v1"</span>
                    </div>
                    <div class="code-line line-3 text-gray-300">
                      <span class="text-primary-400">service_tier</span> = <span class="text-amber-300">"fast"</span>
                    </div>
                    <div class="mt-3 rounded-lg border border-indigo-500/20 bg-indigo-950/20 p-2.5 text-[11px] text-indigo-300">
                      <span>✓ Codex Fast Mode Active · Multi-Session Sticky Activated</span>
                    </div>
                  </div>

                  <!-- Tab 4: Python SDK -->
                  <div v-else class="space-y-2">
                    <div class="text-cyan-400"><span class="text-pink-400">from</span> openai <span class="text-pink-400">import</span> OpenAI</div>
                    <div class="text-gray-300">client = OpenAI(</div>
                    <div class="pl-4 text-gray-400">base_url=<span class="text-amber-300">"https://api.your-domain.com/v1"</span>,</div>
                    <div class="pl-4 text-gray-400">api_key=<span class="text-amber-300">"sk-sub2api-xxxxxxxx"</span></div>
                    <div class="text-gray-300">)</div>
                    <div class="text-gray-300">resp = client.chat.completions.create(model=<span class="text-amber-300">"claude-3-7-sonnet"</span>)</div>
                  </div>

                  <!-- Blinking Cursor -->
                  <div class="mt-3 flex items-center gap-1 text-gray-500">
                    <span class="cursor inline-block h-3.5 w-2 bg-emerald-400"></span>
                  </div>

                </div>

              </div>

            </div>
          </div>

        </div>

        <!-- Model Badges Ribbon (主流大模型矩阵) -->
        <div class="mt-20 border-y border-gray-200/60 py-6 dark:border-white/5">
          <p class="mb-4 text-center text-xs font-semibold uppercase tracking-wider text-gray-400 dark:text-dark-500">
            全量支持主流模型矩阵与专属定制通道
          </p>
          <div class="flex flex-wrap items-center justify-center gap-3 sm:gap-6">
            <span
              v-for="model in supportedModels"
              :key="model.name"
              class="inline-flex items-center gap-2 rounded-xl border border-gray-200/80 bg-white/60 px-3.5 py-1.5 text-xs font-medium text-gray-700 shadow-sm backdrop-blur-sm transition-all hover:border-primary-400 dark:border-white/5 dark:bg-dark-800/60 dark:text-dark-200 dark:hover:border-primary-500/50"
            >
              <span class="h-2 w-2 rounded-full" :class="model.dotColor"></span>
              <span class="font-mono">{{ model.name }}</span>
              <span class="text-[10px] text-gray-400 dark:text-dark-400">{{ model.tag }}</span>
            </span>
          </div>
        </div>

        <!-- Bento Grid: 核心系统优势 (6大核心能力) -->
        <div class="mt-24">
          <div class="text-center">
            <h2 class="text-xs font-bold uppercase tracking-widest text-primary-600 dark:text-cyan-400">
              ARCHITECTURE & CAPABILITIES
            </h2>
            <p class="mt-2 text-3xl font-extrabold tracking-tight text-gray-900 sm:text-4xl dark:text-white">
              企业级稳定性与灵活计费架构
            </p>
            <p class="mx-auto mt-4 max-w-2xl text-sm text-gray-600 dark:text-dark-300 sm:text-base">
              从底层协议清洗、动态上游账号轮询，到终端 IDE 粘性会话与自动化结算，全栈为生产环境护航。
            </p>
          </div>

          <!-- Bento Grid Layout -->
          <div class="mt-12 grid gap-6 md:grid-cols-2 lg:grid-cols-3">
            
            <!-- Bento Card 1: 智能调度与粘性会话 -->
            <div class="group relative overflow-hidden rounded-2xl border border-gray-200/80 bg-white/70 p-6 shadow-sm backdrop-blur-sm transition-all duration-300 hover:-translate-y-1 hover:border-primary-500/40 hover:shadow-xl hover:shadow-primary-500/5 dark:border-white/10 dark:bg-dark-900/60 lg:col-span-2">
              <div class="flex items-center gap-3">
                <div class="flex h-11 w-11 items-center justify-center rounded-xl bg-primary-500/10 text-primary-600 dark:bg-primary-500/20 dark:text-primary-300">
                  <Icon name="bolt" size="md" />
                </div>
                <div>
                  <h3 class="text-base font-bold text-gray-900 dark:text-white">{{ te('home.features.smartScheduling') ? t('home.features.smartScheduling') : '智能会话粘性调度' }}</h3>
                  <p class="text-xs text-gray-500 dark:text-dark-400">SESSION AFFINITY & FAILOVER</p>
                </div>
              </div>
              <p class="mt-4 text-sm leading-relaxed text-gray-600 dark:text-dark-300">
                {{ te('home.features.smartSchedulingDesc') ? t('home.features.smartSchedulingDesc') : '针对 Claude Code 与 Codex 的多轮长上下文会话，自动建立会话与物理上游账号的粘性绑定，避免请求在账号间频繁漂移导致上下文丢失。单账号故障时 0 毫秒感知秒级迁移。' }}
              </p>
              <!-- Mini Visual Simulation -->
              <div class="mt-6 flex flex-wrap items-center gap-2 rounded-xl bg-gray-50 p-3 font-mono text-xs dark:bg-dark-950/60">
                <span class="rounded bg-primary-100 px-2 py-0.5 text-primary-700 dark:bg-primary-950 dark:text-primary-300">Request: session_92fa</span>
                <span class="text-gray-400">→</span>
                <span class="rounded bg-emerald-100 px-2 py-0.5 text-emerald-700 dark:bg-emerald-950 dark:text-emerald-300">Upstream #03 [Active]</span>
                <span class="ml-auto text-[11px] text-gray-400">Sticky Cache: Hit (0.01ms)</span>
              </div>
            </div>

            <!-- Bento Card 2: 精准计量与 Token 实时结算 -->
            <div class="group relative overflow-hidden rounded-2xl border border-gray-200/80 bg-white/70 p-6 shadow-sm backdrop-blur-sm transition-all duration-300 hover:-translate-y-1 hover:border-cyan-500/40 hover:shadow-xl hover:shadow-cyan-500/5 dark:border-white/10 dark:bg-dark-900/60">
              <div class="flex items-center gap-3">
                <div class="flex h-11 w-11 items-center justify-center rounded-xl bg-cyan-500/10 text-cyan-600 dark:bg-cyan-500/20 dark:text-cyan-300">
                  <Icon name="chart" size="md" />
                </div>
                <div>
                  <h3 class="text-base font-bold text-gray-900 dark:text-white">{{ te('home.features.billingQuota') ? t('home.features.billingQuota') : 'Token 级精准扣费' }}</h3>
                  <p class="text-xs text-gray-500 dark:text-dark-400">REAL-TIME METERING</p>
                </div>
              </div>
              <p class="mt-4 text-sm leading-relaxed text-gray-600 dark:text-dark-300">
                {{ te('home.features.billingQuotaDesc') ? t('home.features.billingQuotaDesc') : '支持按 Token 费率精准计费，Prompt 缓存命中折扣、输入输出分级单价、图片/视频生成差异计费一应俱全。' }}
              </p>
              <div class="mt-5 flex items-center justify-between border-t border-gray-100 pt-3 text-xs text-gray-500 dark:border-white/5 dark:text-dark-400">
                <span>Prompt Cache 节省</span>
                <span class="font-bold text-emerald-500 font-mono">最高可达 90%</span>
              </div>
            </div>

            <!-- Bento Card 3: 账号池与防封风控 -->
            <div class="group relative overflow-hidden rounded-2xl border border-gray-200/80 bg-white/70 p-6 shadow-sm backdrop-blur-sm transition-all duration-300 hover:-translate-y-1 hover:border-indigo-500/40 hover:shadow-xl hover:shadow-indigo-500/5 dark:border-white/10 dark:bg-dark-900/60">
              <div class="flex items-center gap-3">
                <div class="flex h-11 w-11 items-center justify-center rounded-xl bg-indigo-500/10 text-indigo-600 dark:bg-indigo-500/20 dark:text-indigo-300">
                  <Icon name="server" size="md" />
                </div>
                <div>
                  <h3 class="text-base font-bold text-gray-900 dark:text-white">{{ te('home.features.multiAccount') ? t('home.features.multiAccount') : '多账号池弹性负载' }}</h3>
                  <p class="text-xs text-gray-500 dark:text-dark-400">LOAD BALANCING</p>
                </div>
              </div>
              <p class="mt-4 text-sm leading-relaxed text-gray-600 dark:text-dark-300">
                {{ te('home.features.multiAccountDesc') ? t('home.features.multiAccountDesc') : '支持 OAuth 账号、API Key 混合池管理，定时探活与自动限流保护，避免突发高频请求打爆单账号。' }}
              </p>
            </div>

            <!-- Bento Card 4: 协议穿透与 IDE 深度优化 -->
            <div class="group relative overflow-hidden rounded-2xl border border-gray-200/80 bg-white/70 p-6 shadow-sm backdrop-blur-sm transition-all duration-300 hover:-translate-y-1 hover:border-primary-500/40 hover:shadow-xl hover:shadow-primary-500/5 dark:border-white/10 dark:bg-dark-900/60">
              <div class="flex items-center gap-3">
                <div class="flex h-11 w-11 items-center justify-center rounded-xl bg-purple-500/10 text-purple-600 dark:bg-purple-500/20 dark:text-purple-300">
                  <Icon name="cog" size="md" />
                </div>
                <div>
                  <h3 class="text-base font-bold text-gray-900 dark:text-white">{{ te('home.features.concurrencyControl') ? t('home.features.concurrencyControl') : 'IDE 深度适配支持' }}</h3>
                  <p class="text-xs text-gray-500 dark:text-dark-400">CODEX & CLAUDE CODE</p>
                </div>
              </div>
              <p class="mt-4 text-sm leading-relaxed text-gray-600 dark:text-dark-300">
                完美兼容 Codex Fast/Flex 调度策略，下划线 Header 穿透保障、长连接 SSE 流式实时推送无卡顿。
              </p>
            </div>

            <!-- Bento Card 5: 多种支付与自助卡密兑换 -->
            <div class="group relative overflow-hidden rounded-2xl border border-gray-200/80 bg-white/70 p-6 shadow-sm backdrop-blur-sm transition-all duration-300 hover:-translate-y-1 hover:border-emerald-500/40 hover:shadow-xl hover:shadow-emerald-500/5 dark:border-white/10 dark:bg-dark-900/60">
              <div class="flex items-center gap-3">
                <div class="flex h-11 w-11 items-center justify-center rounded-xl bg-emerald-500/10 text-emerald-600 dark:bg-emerald-500/20 dark:text-emerald-300">
                  <Icon name="creditCard" size="md" />
                </div>
                <div>
                  <h3 class="text-base font-bold text-gray-900 dark:text-white">{{ te('home.features.securityStability') ? t('home.features.securityStability') : '内置多支付与卡密' }}</h3>
                  <p class="text-xs text-gray-500 dark:text-dark-400">PAYMENT & REDEEM</p>
                </div>
              </div>
              <p class="mt-4 text-sm leading-relaxed text-gray-600 dark:text-dark-300">
                集成微信官方、支付宝官方、Stripe、易支付及批量兑换卡密，支持按量充值或周期订阅套餐双轨制。
              </p>
            </div>

          </div>
        </div>

        <!-- Bottom CTA Box (底部高质感转化引导卡) -->
        <div class="relative mt-24 overflow-hidden rounded-3xl border border-primary-500/30 bg-gradient-to-br from-primary-900/40 via-dark-900 to-[#06090e] p-8 text-center shadow-2xl dark:border-white/10 dark:bg-dark-900/80 sm:p-14">
          <div class="pointer-events-none absolute -right-20 -top-20 h-64 w-64 rounded-full bg-primary-500/20 blur-3xl"></div>
          <div class="pointer-events-none absolute -bottom-20 -left-20 h-64 w-64 rounded-full bg-cyan-500/20 blur-3xl"></div>

          <div class="relative z-10 mx-auto max-w-2xl">
            <h3 class="text-2xl font-bold tracking-tight text-white sm:text-3xl">
              立即接入属于你的高可用 API 算力集群
            </h3>
            <p class="mt-4 text-sm leading-relaxed text-gray-300 sm:text-base">
              无需漫长配置，1 分钟生成专属 API Key，在任意 IDE 与终端工具畅享顶级模型。
            </p>
            <div class="mt-8 flex flex-col items-center justify-center gap-3 sm:flex-row">
              <router-link
                :to="isAuthenticated ? dashboardPath : '/register'"
                class="inline-flex w-full items-center justify-center gap-2 rounded-xl bg-white px-7 py-3 text-sm font-bold text-gray-950 shadow-lg transition-all hover:bg-gray-100 sm:w-auto"
              >
                <span>{{ isAuthenticated ? t('home.dashboard') : (te('home.getStarted') ? t('home.getStarted') : '免费注册使用') }}</span>
                <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3" />
                </svg>
              </router-link>

              <router-link
                to="/download"
                class="inline-flex w-full items-center justify-center gap-2 rounded-xl border border-white/20 bg-white/5 px-6 py-3 text-sm font-semibold text-white backdrop-blur-sm transition-all hover:bg-white/10 sm:w-auto"
              >
                <Icon name="download" size="sm" />
                <span>{{ t('nav.agentDownload') }}</span>
              </router-link>
            </div>
          </div>
        </div>

      </div>
    </main>

    <!-- Footer -->
    <footer class="relative z-10 border-t border-gray-200/80 bg-white/40 py-8 text-xs text-gray-500 backdrop-blur-md dark:border-white/5 dark:bg-dark-950/40 dark:text-dark-400">
      <div class="mx-auto flex max-w-6xl flex-col items-center justify-between gap-4 px-4 sm:flex-row sm:px-6">
        <div class="flex items-center gap-2">
          <span>&copy; {{ currentYear }} {{ siteName }}. All rights reserved.</span>
          <span v-if="icpNumber" class="text-gray-400">|</span>
          <a
            v-if="icpNumber"
            :href="icpLink || 'https://beian.miit.gov.cn/'"
            target="_blank"
            rel="noopener noreferrer"
            class="hover:text-primary-500 transition-colors"
          >
            {{ icpNumber }}
          </a>
        </div>

        <div class="flex items-center gap-4">
          <router-link to="/key-usage" class="hover:text-gray-900 dark:hover:text-white transition-colors">
            {{ t('keyUsage.title') || '用量速查' }}
          </router-link>
          <router-link v-if="showModelPlazaEntry" to="/model-plaza" class="hover:text-gray-900 dark:hover:text-white transition-colors">
            {{ t('nav.modelPlaza') }}
          </router-link>
          <router-link to="/download" class="hover:text-gray-900 dark:hover:text-white transition-colors">
            {{ t('nav.agentDownload') }}
          </router-link>
          <a v-if="docUrl" :href="docUrl" target="_blank" rel="noopener noreferrer" class="hover:text-gray-900 dark:hover:text-white transition-colors">
            {{ t('home.viewDocs') }}
          </a>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore, useAuthStore } from '@/stores'
import Icon from '@/components/icons/Icon.vue'
import LocaleSwitcher from '@/components/common/LocaleSwitcher.vue'

const { t, te } = useI18n()
const appStore = useAppStore()
const authStore = useAuthStore()

// Dark mode state
const isDark = ref(false)

// Custom home content
const homeContent = computed(() => appStore.cachedPublicSettings?.home_content || '')
const hasHomeContent = computed(() => homeContent.value.trim().length > 0)
const isHomeContentUrl = computed(() => {
  const content = homeContent.value.trim()
  return content.startsWith('http://') || content.startsWith('https://')
})

// Compact home mode
const compactHomeEnabled = computed(() => Boolean(appStore.cachedPublicSettings?.compact_home_enabled))

// Model plaza configuration
const modelPlazaEnabled = computed(() => {
  const val = appStore.cachedPublicSettings?.model_plaza_enabled
  return val === undefined || Boolean(val)
})

const modelPlazaRequiresAuth = computed(() => {
  return Boolean(appStore.cachedPublicSettings?.model_plaza_require_auth)
})

const showModelPlazaEntry = computed(
  () =>
    modelPlazaEnabled.value &&
    (authStore.isAuthenticated || !modelPlazaRequiresAuth.value),
)

const showAgentDownloadEntry = computed(() => {
  const val = (appStore.cachedPublicSettings as Record<string, unknown> | null)?.agent_download_enabled
  return val === undefined || Boolean(val)
})

// Site settings
const siteName = computed(() => appStore.siteName)
const siteSubtitle = computed(() => appStore.cachedPublicSettings?.site_subtitle || '')
const siteLogo = computed(() => appStore.siteLogo)
const docUrl = computed(() => appStore.docUrl)
const icpNumber = computed(() => String((appStore.cachedPublicSettings as Record<string, unknown> | null)?.icp_number || ''))
const icpLink = computed(() => String((appStore.cachedPublicSettings as Record<string, unknown> | null)?.icp_link || ''))

// User and route
const isAuthenticated = computed(() => authStore.isAuthenticated)
const isAdmin = computed(() => authStore.isAdmin)
const dashboardPath = computed(() => (isAdmin.value ? '/admin/dashboard' : '/dashboard'))
const userInitial = computed(() => {
  const user = authStore.user
  if (!user || !user.email) return ''
  return user.email.charAt(0).toUpperCase()
})

const currentYear = computed(() => new Date().getFullYear())

// Terminal Tabs Configuration
type TerminalTabId = 'claude' | 'curl' | 'codex' | 'python'
const activeTerminalTab = ref<TerminalTabId>('claude')
const terminalTabs = [
  { id: 'claude' as TerminalTabId, label: 'Claude Code' },
  { id: 'curl' as TerminalTabId, label: 'cURL' },
  { id: 'codex' as TerminalTabId, label: 'Codex' },
  { id: 'python' as TerminalTabId, label: 'Python' }
]

const copied = ref(false)
function copyActiveCode() {
  let text = ''
  if (activeTerminalTab.value === 'claude') {
    text = `export ANTHROPIC_BASE_URL="https://api.your-domain.com"\nexport ANTHROPIC_API_KEY="sk-sub2api-xxxxxxxx"\nclaude --model claude-3-7-sonnet-20250219`
  } else if (activeTerminalTab.value === 'curl') {
    text = `curl https://api.your-domain.com/v1/chat/completions \\\n  -H "Authorization: Bearer sk-sub2api-xxxx" \\\n  -H "Content-Type: application/json" \\\n  -d '{"model": "gpt-4o", "messages": [{"role": "user", "content": "Hello"}]}'`
  } else if (activeTerminalTab.value === 'codex') {
    text = `model_provider = "sub2api"\nbase_url = "https://api.your-domain.com/v1"\nservice_tier = "fast"`
  } else {
    text = `from openai import OpenAI\nclient = OpenAI(base_url="https://api.your-domain.com/v1", api_key="sk-sub2api-xxxxxxxx")\nresp = client.chat.completions.create(model="claude-3-7-sonnet")`
  }
  navigator.clipboard?.writeText(text)
  copied.value = true
  setTimeout(() => {
    copied.value = false
  }, 2000)
}

// Supported Models Banner List
const supportedModels = [
  { name: 'Claude 3.7 Sonnet', tag: 'Anthropic', dotColor: 'bg-amber-400' },
  { name: 'GPT-4.5 / o1', tag: 'OpenAI', dotColor: 'bg-emerald-400' },
  { name: 'Gemini 2.0 Flash', tag: 'Google', dotColor: 'bg-blue-400' },
  { name: 'Grok 3 (Beta)', tag: 'xAI', dotColor: 'bg-purple-400' },
  { name: 'DeepSeek R1 / V3', tag: 'DeepSeek', dotColor: 'bg-cyan-400' }
]

function toggleTheme() {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark', isDark.value)
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
}

function initTheme() {
  const savedTheme = localStorage.getItem('theme')
  if (
    savedTheme === 'dark' ||
    (!savedTheme && window.matchMedia('(prefers-color-scheme: dark)').matches)
  ) {
    isDark.value = true
    document.documentElement.classList.add('dark')
  }
}

onMounted(() => {
  initTheme()
  authStore.checkAuth()
  if (!appStore.publicSettingsLoaded) {
    appStore.fetchPublicSettings()
  }
})
</script>

<style scoped>
/* Terminal Container & Styling 保障测试兼容与动画 */
.terminal-container {
  position: relative;
  display: block;
}

.terminal-window {
  box-shadow:
    0 25px 50px -12px rgba(0, 0, 0, 0.4),
    0 0 0 1px rgba(255, 255, 255, 0.08),
    inset 0 1px 0 rgba(255, 255, 255, 0.1);
  overflow: hidden;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.terminal-window:hover {
  transform: translateY(-2px);
  box-shadow:
    0 30px 60px -15px rgba(0, 0, 0, 0.5),
    0 0 0 1px rgba(99, 102, 241, 0.3),
    0 0 30px rgba(99, 102, 241, 0.15);
}

.terminal-header {
  display: flex;
  align-items: center;
}

.terminal-buttons span {
  display: inline-block;
}

.code-line {
  animation: line-appear 0.3s ease forwards;
}

@keyframes line-appear {
  from {
    opacity: 0;
    transform: translateY(3px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.cursor {
  animation: blink 1s step-end infinite;
}

@keyframes blink {
  0%, 50% { opacity: 1; }
  51%, 100% { opacity: 0; }
}
</style>
