<template>
  <div class="space-y-16 sm:space-y-24">
    <!-- ==================== 1. Hero 区域 ==================== -->
    <section class="relative pt-6 pb-10 text-center sm:pt-12 sm:pb-16">
      <div class="inline-flex items-center gap-2 px-4 py-1.5 rounded-full border border-orange-200 bg-orange-50 text-xs sm:text-sm font-bold text-orange-700 shadow-sm dark:border-orange-500/30 dark:bg-orange-950/40 dark:text-orange-300 mb-6">
        <span>⚡ {{ t('agentDownload.tour.step4.quote', '人间清醒指南：是你在掌控工具，不是给 AI 当全职保姆') }}</span>
      </div>

      <h1 class="text-3xl sm:text-5xl lg:text-6xl font-black tracking-tight text-slate-950 dark:text-white leading-[1.18] mb-6">
        <span class="block">{{ t('agentDownload.title', '魔丸 Agent 专属客户端') }}</span>
        <span class="bg-gradient-to-r from-orange-600 via-pink-600 to-purple-600 bg-clip-text text-transparent">
          你在床上躺平，它在后台杀穿。
        </span>
      </h1>

      <p class="text-base sm:text-lg text-slate-600 dark:text-dark-300 max-w-3xl mx-auto font-normal leading-relaxed mb-8">
        {{ t('agentDownload.subtitle', '有些工具你得像供祖宗一样守在电脑前伺候，而魔丸只会让你体会什么叫真正的“老板级掌控”。') }}<br class="hidden sm:inline" />
        <strong class="text-slate-900 dark:text-white font-bold">局域网扫码秒连</strong> · 
        <strong class="text-slate-900 dark:text-white font-bold">免费内网穿透</strong> · 
        <strong class="text-slate-900 dark:text-white font-bold">实测连续狂飙 400+ 分钟数千步绝不降智</strong>。
      </p>

      <!-- 行动按钮 -->
      <div class="flex flex-col sm:flex-row items-center justify-center gap-4 mb-8">
        <a
          href="#download-section"
          class="w-full sm:w-auto px-8 py-4 rounded-2xl text-base font-bold text-white bg-gradient-to-r from-orange-600 via-amber-600 to-rose-600 hover:from-orange-500 hover:to-amber-500 shadow-xl shadow-orange-500/25 transition active:scale-[0.98]"
        >
          🚀 {{ t('agentDownload.downloadPrimary', '立即下载 魔丸 Agent') }}
        </a>
        <a
          href="#models"
          class="w-full sm:w-auto px-8 py-4 rounded-2xl text-base font-bold text-slate-700 bg-white hover:bg-slate-50 border border-slate-300 shadow-sm transition dark:border-white/10 dark:bg-dark-800 dark:text-dark-200 dark:hover:bg-dark-700"
        >
          ✨ 查看全明星模型矩阵 & 专精 Gemini
        </a>
      </div>

      <!-- 扎心真相条 -->
      <div class="max-w-2xl mx-auto p-4 rounded-2xl bg-amber-50/80 border border-amber-200/80 text-xs sm:text-sm text-amber-900 font-medium dark:border-amber-500/30 dark:bg-amber-950/30 dark:text-amber-300">
        💡 <strong>扎心真相：</strong>“{{ t('agentDownload.tour.step5.quote', '天天端坐在电脑前盯进度条，不是因为你敬业，而是你用的玩具 Agent 根本不配让你安心去睡大觉。') }}”
      </div>
    </section>

    <!-- ==================== 2. 全明星大模型矩阵 & 专精 Gemini ==================== -->
    <section id="models" class="pt-8 border-t border-slate-200/60 dark:border-white/10">
      <div class="text-center max-w-3xl mx-auto mb-12">
        <span class="text-xs uppercase font-extrabold tracking-widest text-indigo-600 bg-indigo-100 px-3 py-1 rounded-full dark:bg-indigo-950/60 dark:text-indigo-400">
          AI Matrix · 顶尖智商不设限
        </span>
        <h2 class="text-2xl sm:text-4xl font-black text-slate-950 dark:text-white mt-4 mb-3">
          {{ t('agentDownload.features.title', '全面适配全球顶尖大模型') }}<br />
          <span class="bg-gradient-to-r from-orange-600 via-pink-600 to-purple-600 bg-clip-text text-transparent">
            {{ t('agentDownload.features.modelsTitle', '重点深度调优 Google Gemini 旗舰') }}
          </span>
        </h2>
        <p class="text-slate-600 dark:text-dark-300 text-sm sm:text-base">
          {{ t('agentDownload.features.modelsDesc', '拒绝任何模型绑架。ChatGPT、Claude、Grok、DeepSeek、GLM、Kimi 原生即插即用；更对 Gemini 顶尖大模型实行底层级吞吐与上下文深度榨干，打造长任务算力天花板。') }}
        </p>
      </div>

      <!-- 焦点巨幕卡片：Google Gemini 顶尖旗舰模型专属调优专区 -->
      <div class="mb-12 bg-gradient-to-br from-indigo-950 via-slate-900 to-purple-950 text-white rounded-3xl p-6 sm:p-10 shadow-2xl border border-indigo-500/20 relative overflow-hidden">
        <!-- 装饰背景光效 -->
        <div class="absolute -right-20 -top-20 w-80 h-80 bg-blue-500/20 rounded-full blur-3xl pointer-events-none"></div>
        <div class="absolute -left-20 -bottom-20 w-80 h-80 bg-purple-500/20 rounded-full blur-3xl pointer-events-none"></div>

        <div class="relative z-10 grid grid-cols-1 lg:grid-cols-12 gap-8 items-center">
          <div class="lg:col-span-7 space-y-5">
            <div class="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-blue-500/20 border border-blue-400/30 text-xs font-bold text-blue-300">
              <span>💎 独家特调 · 重点核心适配</span>
            </div>

            <h3 class="text-2xl sm:text-3xl font-black tracking-tight text-white leading-tight">
              Google Gemini 顶尖大模型深度咬合：<br />
              <span class="text-transparent bg-clip-text bg-gradient-to-r from-blue-400 via-teal-300 to-indigo-300">
                百万级上下文吞吐 + 极限低延迟并发
              </span>
            </h3>

            <p class="text-slate-300 text-sm leading-relaxed">
              市面上许多 Agent 遇到百万级上下文就卡死崩溃，或者遇到长输出就格式错乱。魔丸对 <strong>Gemini 3.7 / 3.8 / 2.5 Pro & Flash</strong> 进行了指令集与上下文缓存（Context Caching）协议级深度调优，直接把长任务代码工程的性价比拉满！
            </p>

            <div class="grid grid-cols-1 sm:grid-cols-3 gap-3 pt-2">
              <div class="p-3.5 rounded-2xl bg-white/5 border border-white/10 backdrop-blur-sm">
                <div class="text-blue-400 font-bold text-sm mb-1">⚡ 专属缓存机制</div>
                <div class="text-xs text-slate-300">上下文重复命中率高达 90%+，跑通宵也不心疼 Token</div>
              </div>
              <div class="p-3.5 rounded-2xl bg-white/5 border border-white/10 backdrop-blur-sm">
                <div class="text-teal-300 font-bold text-sm mb-1">🧠 深度思考链隐藏</div>
                <div class="text-xs text-slate-300">完美适配 Thinking 推理过程解析，严谨纠错不外溢</div>
              </div>
              <div class="p-3.5 rounded-2xl bg-white/5 border border-white/10 backdrop-blur-sm">
                <div class="text-purple-300 font-bold text-sm mb-1">🚀 极速吞吐并发</div>
                <div class="text-xs text-slate-300">支持 100+ tok/s 极速代码喷涌，数千步调用不卡壳</div>
              </div>
            </div>
          </div>

          <div class="lg:col-span-5">
            <div class="p-6 rounded-2xl bg-white/10 border border-white/15 backdrop-blur-md space-y-4">
              <div class="flex items-center gap-3">
                <span class="text-2xl">⚡</span>
                <div>
                  <div class="font-bold text-base text-white">Google Gemini 调优矩阵</div>
                  <div class="text-xs text-blue-200 font-mono">GEMINI PROTOCOL TUNING</div>
                </div>
              </div>

              <div class="space-y-2.5 text-xs text-slate-200">
                <div class="flex items-center justify-between p-2 rounded-lg bg-black/20">
                  <span class="font-semibold">Gemini 3.7 Pro / Flash</span>
                  <span class="text-emerald-400 font-bold">深度 Thinking + 代码全量生成</span>
                </div>
                <div class="flex items-center justify-between p-2 rounded-lg bg-black/20">
                  <span class="font-semibold">Gemini 3.8 Flash High</span>
                  <span class="text-cyan-300 font-bold">超高速长任务主力引擎</span>
                </div>
                <div class="flex items-center justify-between p-2 rounded-lg bg-black/20">
                  <span class="font-semibold">Gemini 2.5 Pro / Flash</span>
                  <span class="text-amber-300 font-bold">百万上下文极限压榨</span>
                </div>
                <div class="flex items-center justify-between p-2 rounded-lg bg-black/20">
                  <span class="font-semibold">上下文缓存优化 (Caching)</span>
                  <span class="text-emerald-400 font-bold">支持 90%+ 缓存复用率</span>
                </div>
              </div>

              <div class="pt-2 text-[11px] text-slate-400 border-t border-white/10 flex items-center justify-between">
                <span>协议状态: 深度整合完成</span>
                <span class="text-emerald-400 font-bold">● 实时稳定</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 6 大顶尖模型矩阵网格 -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">
        <!-- 1. Claude -->
        <div class="p-5 rounded-2xl bg-white dark:bg-dark-900 border border-slate-200 dark:border-white/10 shadow-sm hover:border-slate-300 dark:hover:border-white/20 transition">
          <div class="flex items-center justify-between mb-3">
            <div class="flex items-center gap-3">
              <span class="w-9 h-9 rounded-xl bg-amber-50 text-amber-700 font-bold flex items-center justify-center text-sm border border-amber-200 dark:bg-amber-950/40 dark:text-amber-300 dark:border-amber-500/30">
                CL
              </span>
              <div>
                <h4 class="font-bold text-slate-900 dark:text-white text-sm">Anthropic Claude Code</h4>
                <span class="text-[11px] text-slate-500 dark:text-dark-400 font-mono">Claude 3.5 Sonnet / Opus 3.7</span>
              </div>
            </div>
            <span class="text-[10px] font-bold px-2 py-0.5 rounded bg-emerald-100 text-emerald-800 dark:bg-emerald-950/50 dark:text-emerald-300">完美适配</span>
          </div>
          <p class="text-xs text-slate-600 dark:text-dark-300 leading-relaxed">
            代码架构与逻辑分析天花板。无论是巨型代码库重构还是复杂算法编写，魔丸均能精准调动 Claude 的深度推理能力。
          </p>
        </div>

        <!-- 2. OpenAI GPT / Codex -->
        <div class="p-5 rounded-2xl bg-white dark:bg-dark-900 border border-slate-200 dark:border-white/10 shadow-sm hover:border-slate-300 dark:hover:border-white/20 transition">
          <div class="flex items-center justify-between mb-3">
            <div class="flex items-center gap-3">
              <span class="w-9 h-9 rounded-xl bg-emerald-50 text-emerald-700 font-bold flex items-center justify-center text-sm border border-emerald-200 dark:bg-emerald-950/40 dark:text-emerald-300 dark:border-emerald-500/30">
                GPT
              </span>
              <div>
                <h4 class="font-bold text-slate-900 dark:text-white text-sm">OpenAI GPT / Codex</h4>
                <span class="text-[11px] text-slate-500 dark:text-dark-400 font-mono">GPT-4.5 / o1 / o3-mini</span>
              </div>
            </div>
            <span class="text-[10px] font-bold px-2 py-0.5 rounded bg-emerald-100 text-emerald-800 dark:bg-emerald-950/50 dark:text-emerald-300">完美适配</span>
          </div>
          <p class="text-xs text-slate-600 dark:text-dark-300 leading-relaxed">
            极高泛化能力与通用工程解决专家。原生支持 Function Calling 与 Structured Output，指令执行精准无偏差。
          </p>
        </div>

        <!-- 3. DeepSeek -->
        <div class="p-5 rounded-2xl bg-white dark:bg-dark-900 border border-slate-200 dark:border-white/10 shadow-sm hover:border-slate-300 dark:hover:border-white/20 transition">
          <div class="flex items-center justify-between mb-3">
            <div class="flex items-center gap-3">
              <span class="w-9 h-9 rounded-xl bg-blue-50 text-blue-700 font-bold flex items-center justify-center text-sm border border-blue-200 dark:bg-blue-950/40 dark:text-blue-300 dark:border-blue-500/30">
                DS
              </span>
              <div>
                <h4 class="font-bold text-slate-900 dark:text-white text-sm">深度求索 / DeepSeek</h4>
                <span class="text-[11px] text-slate-500 dark:text-dark-400 font-mono">DeepSeek V3 / R1 满血版</span>
              </div>
            </div>
            <span class="text-[10px] font-bold px-2 py-0.5 rounded bg-blue-100 text-blue-800 dark:bg-blue-950/50 dark:text-blue-300">深度优化</span>
          </div>
          <p class="text-xs text-slate-600 dark:text-dark-300 leading-relaxed">
            极致性价比的国产大模型之光。针对 R1 的长链思考过程进行专门过滤与解析，不漏关键信息，代码推理极速交付。
          </p>
        </div>

        <!-- 4. 智谱 GLM -->
        <div class="p-5 rounded-2xl bg-white dark:bg-dark-900 border border-slate-200 dark:border-white/10 shadow-sm hover:border-slate-300 dark:hover:border-white/20 transition">
          <div class="flex items-center justify-between mb-3">
            <div class="flex items-center gap-3">
              <span class="w-9 h-9 rounded-xl bg-indigo-50 text-indigo-700 font-bold flex items-center justify-center text-sm border border-indigo-200 dark:bg-indigo-950/40 dark:text-indigo-300 dark:border-indigo-500/30">
                GLM
              </span>
              <div>
                <h4 class="font-bold text-slate-900 dark:text-white text-sm">智谱清言 / GLM</h4>
                <span class="text-[11px] text-slate-500 dark:text-dark-400 font-mono">GLM-4-Plus / GLM-4-Air</span>
              </div>
            </div>
            <span class="text-[10px] font-bold px-2 py-0.5 rounded bg-emerald-100 text-emerald-800 dark:bg-emerald-950/50 dark:text-emerald-300">完美适配</span>
          </div>
          <p class="text-xs text-slate-600 dark:text-dark-300 leading-relaxed">
            全功能国产多模态与超长上下文代表。中文语境语义理解与行业工程代码微调极其扎实，长任务执行稳定不宕机。
          </p>
        </div>

        <!-- 5. Kimi -->
        <div class="p-5 rounded-2xl bg-white dark:bg-dark-900 border border-slate-200 dark:border-white/10 shadow-sm hover:border-slate-300 dark:hover:border-white/20 transition">
          <div class="flex items-center justify-between mb-3">
            <div class="flex items-center gap-3">
              <span class="w-9 h-9 rounded-xl bg-violet-50 text-violet-700 font-bold flex items-center justify-center text-sm border border-violet-200 dark:bg-violet-950/40 dark:text-violet-300 dark:border-violet-500/30">
                KM
              </span>
              <div>
                <h4 class="font-bold text-slate-900 dark:text-white text-sm">月之暗面 / Kimi</h4>
                <span class="text-[11px] text-slate-500 dark:text-dark-400 font-mono">Kimi 1.5 · k1.5 长上下文</span>
              </div>
            </div>
            <span class="text-[10px] font-bold px-2 py-0.5 rounded bg-violet-100 text-violet-800 dark:bg-violet-950/50 dark:text-violet-300">完美适配</span>
          </div>
          <p class="text-xs text-slate-600 dark:text-dark-300 leading-relaxed">
            长文本检索与超长文档阅读理解专家。巨型项目几十个依赖库与技术文档，一网打尽无缝对照，滴水不漏。
          </p>
        </div>

        <!-- 6. 更多生态 -->
        <div class="p-5 rounded-2xl bg-white dark:bg-dark-900 border border-slate-200 dark:border-white/10 shadow-sm hover:border-slate-300 dark:hover:border-white/20 transition">
          <div class="flex items-center justify-between mb-3">
            <div class="flex items-center gap-3">
              <span class="w-9 h-9 rounded-xl bg-pink-50 text-pink-700 font-bold flex items-center justify-center text-sm border border-pink-200 dark:bg-pink-950/40 dark:text-pink-300 dark:border-pink-500/30">
                +ALL
              </span>
              <div>
                <h4 class="font-bold text-slate-900 dark:text-white text-sm">更多顶尖模型生态</h4>
                <span class="text-[11px] text-slate-500 dark:text-dark-400 font-mono">Qwen / MiniMax / Grok 3</span>
              </div>
            </div>
            <span class="text-[10px] font-bold px-2 py-0.5 rounded bg-slate-100 text-slate-700 dark:bg-dark-800 dark:text-dark-300">持续扩展</span>
          </div>
          <p class="text-xs text-slate-600 dark:text-dark-300 leading-relaxed">
            通义千问 Qwen 2.5、MiniMax abab6.5、百川智能、xAI Grok 3 等均支持一键接入，随时按需自由切换调用。
          </p>
        </div>
      </div>
    </section>

    <!-- ==================== 3. 手机躺平操控实录 ==================== -->
    <section id="lay-flat" class="pt-8 border-t border-slate-200/60 dark:border-white/10">
      <div class="text-center max-w-3xl mx-auto mb-12">
        <span class="text-xs uppercase font-extrabold tracking-widest text-orange-600 bg-orange-100 px-3 py-1 rounded-full dark:bg-orange-950/60 dark:text-orange-400">
          Real World Mobile Control
        </span>
        <h2 class="text-2xl sm:text-4xl font-black text-slate-950 dark:text-white mt-4 mb-3">
          电脑放桌上跑，你躺床上遥控
        </h2>
        <p class="text-slate-600 dark:text-dark-300 text-sm sm:text-base">
          谁说写代码一定要正襟危坐？魔丸原生内置移动端自适应引擎，手机一掏就是顶级指挥官。
        </p>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
        <!-- 卡片 1 -->
        <div class="bg-white dark:bg-dark-900 rounded-3xl border border-slate-200 dark:border-white/10 p-6 shadow-sm flex flex-col justify-between">
          <div>
            <div class="text-orange-600 dark:text-orange-400 font-black text-xl mb-1">01 / 扫码秒连</div>
            <h3 class="font-bold text-slate-900 dark:text-white text-base mb-2">局域网原生直连，零手动配置</h3>
            <p class="text-xs text-slate-600 dark:text-dark-300 leading-relaxed mb-4">
              电脑端一键启动，自动探测局域网物理网卡并生成专属二维码。手机相机一扫即入，免密认证，连 IP 地址都不用记。
            </p>
          </div>
          <div
            class="group relative cursor-zoom-in overflow-hidden rounded-xl border border-slate-200 dark:border-dark-700"
            @click="openLightbox('/assets/promo/scan-qrcode.png', '局域网扫码秒连')"
          >
            <img src="/assets/promo/scan-qrcode.png" alt="局域网扫码秒连" class="w-full h-48 object-cover object-top transition duration-200 group-hover:scale-105" />
            <div class="absolute bottom-2 right-2 rounded bg-black/70 px-2 py-0.5 text-[10px] text-white opacity-0 group-hover:opacity-100 transition">点击放大</div>
          </div>
        </div>

        <!-- 卡片 2 -->
        <div class="bg-white dark:bg-dark-900 rounded-3xl border border-slate-200 dark:border-white/10 p-6 shadow-sm flex flex-col justify-between">
          <div>
            <div class="text-pink-600 dark:text-pink-400 font-black text-xl mb-1">02 / 免费穿透</div>
            <h3 class="font-bold text-slate-900 dark:text-white text-base mb-2">免费公网穿透，出门在外随时管</h3>
            <p class="text-xs text-slate-600 dark:text-dark-300 leading-relaxed mb-4">
              内置原生 Cloudflare 穿透协议，无需公网 IP 和复杂反代。一键生成加密专属 URL，哪怕在咖啡厅或地铁上也能实时查看长任务进度。
            </p>
          </div>
          <div
            class="group relative cursor-zoom-in overflow-hidden rounded-xl border border-slate-200 dark:border-dark-700"
            @click="openLightbox('/assets/promo/lan-settings.png', '免费公网穿透配置')"
          >
            <img src="/assets/promo/lan-settings.png" alt="免费公网穿透" class="w-full h-48 object-cover object-top transition duration-200 group-hover:scale-105" />
            <div class="absolute bottom-2 right-2 rounded bg-black/70 px-2 py-0.5 text-[10px] text-white opacity-0 group-hover:opacity-100 transition">点击放大</div>
          </div>
        </div>

        <!-- 卡片 3 -->
        <div class="bg-white dark:bg-dark-900 rounded-3xl border border-slate-200 dark:border-white/10 p-6 shadow-sm flex flex-col justify-between">
          <div>
            <div class="text-purple-600 dark:text-purple-400 font-black text-xl mb-1">03 / 躺平执行</div>
            <h3 class="font-bold text-slate-900 dark:text-white text-base mb-2">手机下达指令，醒来直接收网</h3>
            <p class="text-xs text-slate-600 dark:text-dark-300 leading-relaxed mb-4">
              在手机上单手打一句话：“把整个项目做响应式重构并跑通自动化测试”，手机丢一边去睡觉，魔丸自动拆解、写代码、过门禁全闭环。
            </p>
          </div>
          <div
            class="group relative cursor-zoom-in overflow-hidden rounded-xl border border-slate-200 dark:border-dark-700"
            @click="openLightbox('/assets/promo/mobile-chat.png', '手机全自主操控实录')"
          >
            <img src="/assets/promo/mobile-chat.png" alt="手机下令执行" class="w-full h-48 object-cover object-top transition duration-200 group-hover:scale-105" />
            <div class="absolute bottom-2 right-2 rounded bg-black/70 px-2 py-0.5 text-[10px] text-white opacity-0 group-hover:opacity-100 transition">点击放大</div>
          </div>
        </div>
      </div>
    </section>

    <!-- ==================== 4. 硬核实测·2094步硬刚到底 ==================== -->
    <section id="hardcore-proof" class="pt-8 border-t border-slate-200/60 dark:border-white/10">
      <div class="text-center max-w-3xl mx-auto mb-12">
        <span class="text-xs uppercase font-extrabold tracking-widest text-rose-600 bg-rose-100 px-3 py-1 rounded-full dark:bg-rose-950/60 dark:text-rose-400">
          Hardcore Real Test
        </span>
        <h2 class="text-2xl sm:text-4xl font-black text-slate-950 dark:text-white mt-4 mb-3">
          2094 步硬刚到底 · 400+ 分钟连续狂飙绝不降智
        </h2>
        <p class="text-slate-600 dark:text-dark-300 text-sm sm:text-base">
          用事实说话。市面上绝大多数 Agent 聊到 20 轮就上下文爆炸、开始胡言乱语。看看魔丸在极端长链路任务下的真实战报：
        </p>
      </div>

      <div class="bg-gradient-to-br from-slate-900 to-slate-950 text-white rounded-3xl p-6 sm:p-10 shadow-2xl border border-slate-800">
        <div class="grid grid-cols-1 lg:grid-cols-12 gap-8 items-center">
          <div class="lg:col-span-6 space-y-5">
            <div class="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-emerald-500/20 border border-emerald-500/30 text-xs font-bold text-emerald-400">
              <span>● 真实工程实录（非 Demo）</span>
            </div>

            <h3 class="text-2xl sm:text-3xl font-black text-white leading-tight">
              单任务连续狂飙 400+ 分钟<br />
              <span class="text-transparent bg-clip-text bg-gradient-to-r from-emerald-400 via-teal-300 to-cyan-300">
                49 轮长循环 · 2094 步工具自主调度
              </span>
            </h3>

            <p class="text-slate-300 text-sm leading-relaxed">
              从底层架构重构、全量测试用例编写，到修复数十个隐蔽并发 Bug 并通过 CI/CD 门禁检测，魔丸全程死死咬住最终目标，不降智、不失忆、不幻觉！
            </p>

            <div class="grid grid-cols-3 gap-3 pt-2">
              <div class="p-3 rounded-xl bg-white/5 border border-white/10">
                <div class="text-emerald-400 text-xl sm:text-2xl font-black">2094</div>
                <div class="text-[11px] text-slate-400 font-medium">自主工具调度步数</div>
              </div>
              <div class="p-3 rounded-xl bg-white/5 border border-white/10">
                <div class="text-cyan-400 text-xl sm:text-2xl font-black">400+</div>
                <div class="text-[11px] text-slate-400 font-medium">连续执行分钟数</div>
              </div>
              <div class="p-3 rounded-xl bg-white/5 border border-white/10">
                <div class="text-amber-400 text-xl sm:text-2xl font-black">100%</div>
                <div class="text-[11px] text-slate-400 font-medium">测试用例全绿通过</div>
              </div>
            </div>
          </div>

          <div class="lg:col-span-6">
            <div
              class="group relative cursor-zoom-in overflow-hidden rounded-2xl border border-slate-700 bg-slate-950 p-2 shadow-2xl transition hover:border-slate-500"
              @click="openLightbox('/assets/promo/pc-hardcore-2094steps.png', '49轮·2094步硬核实测大图')"
            >
              <img
                src="/assets/promo/pc-hardcore-2094steps.png"
                alt="2094步硬核长任务实测战报"
                class="w-full object-cover rounded-xl transition duration-300 group-hover:scale-[1.02]"
              />
              <div class="absolute bottom-4 right-4 rounded-lg bg-black/80 px-3 py-1 text-xs text-white opacity-0 group-hover:opacity-100 transition">
                🔍 点击查看完整高清战报
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- ==================== 5. 魔丸功能全景大图 ==================== -->
    <section id="tour" class="pt-8 border-t border-slate-200/60 dark:border-white/10">
      <div class="text-center max-w-3xl mx-auto mb-12">
        <span class="text-xs uppercase font-extrabold tracking-widest text-indigo-600 bg-indigo-100 px-3 py-1 rounded-full dark:bg-indigo-950/60 dark:text-indigo-400">
          {{ t('agentDownload.tour.badge', 'Visual Tour') }}
        </span>
        <h2 class="text-2xl sm:text-4xl font-black text-slate-950 dark:text-white mt-4 mb-3">
          {{ t('agentDownload.tour.title', '魔丸功能全景大图') }}
        </h2>
        <p class="text-slate-600 dark:text-dark-300 text-sm sm:text-base">
          {{ t('agentDownload.tour.subtitle', '实机全景直出，真实直观呈现多 Agent 协作工作台、原生局域网与穿透配置及后台全自主执行。') }}
        </p>
      </div>

      <!-- 3 大全景卡片 -->
      <div class="space-y-12">
        <!-- 全景 1: 多 Agent 协作工作台 -->
        <div class="grid grid-cols-1 lg:grid-cols-12 gap-8 items-center bg-white dark:bg-dark-900 p-6 sm:p-8 rounded-3xl border border-slate-200 dark:border-white/10 shadow-sm">
          <div class="lg:col-span-5 space-y-4">
            <span class="text-xs font-bold text-orange-600 dark:text-orange-400 uppercase tracking-wider font-mono">
              View 01 · {{ t('agentDownload.features.loopTitle', '协作工作台') }}
            </span>
            <h3 class="text-xl sm:text-2xl font-black text-slate-900 dark:text-white">
              {{ t('agentDownload.tour.step1.title', '多智能体协同循环，告别单轮玩具对话') }}
            </h3>
            <p class="text-xs sm:text-sm text-slate-600 dark:text-dark-300 leading-relaxed">
              {{ t('agentDownload.features.loopDesc', '内置 Ralph 任务循环与 Subagent 集群协同。主智能体拆解目标，子智能体并发执行文件读写、网页抓取与代码验证，自动汇总决策。') }}
            </p>
          </div>
          <div class="lg:col-span-7">
            <div
              class="group relative cursor-zoom-in overflow-hidden rounded-2xl border border-slate-200 dark:border-dark-700 shadow-lg"
              @click="openLightbox('/assets/agent/guide/guide-workbench.png', '多 Agent 协作工作台')"
            >
              <img src="/assets/agent/guide/guide-workbench.png" alt="多 Agent 协作工作台" class="w-full object-cover transition duration-300 group-hover:scale-[1.02]" />
              <div class="absolute bottom-3 right-3 rounded bg-black/70 px-2.5 py-1 text-xs text-white opacity-0 group-hover:opacity-100 transition">点击放大</div>
            </div>
          </div>
        </div>

        <!-- 全景 2: 原生局域网与穿透配置 -->
        <div class="grid grid-cols-1 lg:grid-cols-12 gap-8 items-center bg-white dark:bg-dark-900 p-6 sm:p-8 rounded-3xl border border-slate-200 dark:border-white/10 shadow-sm">
          <div class="lg:col-span-7 order-2 lg:order-1">
            <div
              class="group relative cursor-zoom-in overflow-hidden rounded-2xl border border-slate-200 dark:border-dark-700 shadow-lg"
              @click="openLightbox('/assets/agent/guide/guide-settings.png', '原生局域网与穿透配置')"
            >
              <img src="/assets/agent/guide/guide-settings.png" alt="原生局域网与穿透配置" class="w-full object-cover transition duration-300 group-hover:scale-[1.02]" />
              <div class="absolute bottom-3 right-3 rounded bg-black/70 px-2.5 py-1 text-xs text-white opacity-0 group-hover:opacity-100 transition">点击放大</div>
            </div>
          </div>
          <div class="lg:col-span-5 space-y-4 order-1 lg:order-2">
            <span class="text-xs font-bold text-pink-600 dark:text-pink-400 uppercase tracking-wider font-mono">
              View 02 · 极速配置
            </span>
            <h3 class="text-xl sm:text-2xl font-black text-slate-900 dark:text-white">
              {{ t('agentDownload.tour.step2.title', 'Sub2API 一键秒级导入，绿色连通即刻起航') }}
            </h3>
            <p class="text-xs sm:text-sm text-slate-600 dark:text-dark-300 leading-relaxed">
              {{ t('agentDownload.tour.step2.desc', '支持一键唤起客户端写入 Base URL 与 API Key，系统自动探测并校验可用性。连通测试全绿，无需繁琐参数配置。') }}
            </p>
          </div>
        </div>

        <!-- 全景 3: 后台全自主沙盒执行 -->
        <div class="grid grid-cols-1 lg:grid-cols-12 gap-8 items-center bg-white dark:bg-dark-900 p-6 sm:p-8 rounded-3xl border border-slate-200 dark:border-white/10 shadow-sm">
          <div class="lg:col-span-5 space-y-4">
            <span class="text-xs font-bold text-purple-600 dark:text-purple-400 uppercase tracking-wider font-mono">
              View 03 · {{ t('agentDownload.features.sandboxTitle', '沙盒安全') }}
            </span>
            <h3 class="text-xl sm:text-2xl font-black text-slate-900 dark:text-white">
              {{ t('agentDownload.tour.step3.title', '本地安全沙箱隔离，命令与文件受控执行') }}
            </h3>
            <p class="text-xs sm:text-sm text-slate-600 dark:text-dark-300 leading-relaxed">
              {{ t('agentDownload.features.sandboxDesc', '内置安全的 PowerShell 终端执行、文件读写、网页抓取与代码运行能力，在受控环境中高效交付。') }}
            </p>
          </div>
          <div class="lg:col-span-7">
            <div
              class="group relative cursor-zoom-in overflow-hidden rounded-2xl border border-slate-200 dark:border-dark-700 shadow-lg"
              @click="openLightbox('/assets/agent/guide/guide-execution.png', '后台全自主沙盒执行')"
            >
              <img src="/assets/agent/guide/guide-execution.png" alt="后台全自主沙盒执行" class="w-full object-cover transition duration-300 group-hover:scale-[1.02]" />
              <div class="absolute bottom-3 right-3 rounded bg-black/70 px-2.5 py-1 text-xs text-white opacity-0 group-hover:opacity-100 transition">点击放大</div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- ==================== 6. 极速上手 3 步走 & 快速配置助手 & FAQ ==================== -->
    <section id="guide" class="pt-8 border-t border-slate-200/60 dark:border-white/10">
      <div class="text-center max-w-3xl mx-auto mb-12">
        <span class="text-xs uppercase font-extrabold tracking-widest text-emerald-600 bg-emerald-100 px-3 py-1 rounded-full dark:bg-emerald-950/60 dark:text-emerald-400">
          Quick Start Guide
        </span>
        <h2 class="text-2xl sm:text-4xl font-black text-slate-950 dark:text-white mt-4 mb-3">
          躺平上手只需 3 步，比点外卖还快
        </h2>
        <p class="text-slate-600 dark:text-dark-300 text-sm sm:text-base">
          不用学任何花哨的 Prompt，开箱即用，秒级连通。
        </p>
      </div>

      <!-- 3 步卡片 -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-12">
        <!-- 步骤 1 -->
        <div class="bg-white dark:bg-dark-900 rounded-3xl border border-slate-200 dark:border-white/10 p-6 shadow-sm flex flex-col justify-between">
          <div>
            <div class="flex items-center gap-3 mb-3">
              <span class="w-7 h-7 rounded-full bg-orange-600 text-white font-black text-xs flex items-center justify-center">1</span>
              <span class="text-xs font-bold text-orange-600 dark:text-orange-400 font-mono">{{ t('agentDownload.step1', 'Step 01 · 下载解压') }}</span>
            </div>
            <h3 class="text-base font-bold text-slate-900 dark:text-white mb-2">下载客户端，一键运行</h3>
            <p class="text-xs text-slate-600 dark:text-dark-300 leading-relaxed mb-4">
              下载 Windows 便携版或标准安装包。自带独立运行环境，双击即可开启服务。
            </p>
          </div>
          <div class="rounded-xl overflow-hidden border border-slate-200 dark:border-dark-700 bg-slate-50 dark:bg-dark-800 p-2">
            <img src="/assets/promo/guide-settings.png" alt="步骤1：导入配置" class="w-full h-36 object-cover object-top rounded-lg" />
          </div>
        </div>

        <!-- 步骤 2 -->
        <div class="bg-white dark:bg-dark-900 rounded-3xl border border-slate-200 dark:border-white/10 p-6 shadow-sm flex flex-col justify-between">
          <div>
            <div class="flex items-center gap-3 mb-3">
              <span class="w-7 h-7 rounded-full bg-pink-600 text-white font-black text-xs flex items-center justify-center">2</span>
              <span class="text-xs font-bold text-pink-600 dark:text-pink-400 font-mono">{{ t('agentDownload.step2', 'Step 02 · 扫码直连') }}</span>
            </div>
            <h3 class="text-base font-bold text-slate-900 dark:text-white mb-2">掏出手机扫一扫，免密秒登</h3>
            <p class="text-xs text-slate-600 dark:text-dark-300 leading-relaxed mb-4">
              点击【扫码免密直连】弹出二维码。手机相机或浏览器扫码，免密直接进入工作区。
            </p>
          </div>
          <div class="rounded-xl overflow-hidden border border-slate-200 dark:border-dark-700 bg-slate-50 dark:bg-dark-800 p-2">
            <img src="/assets/promo/scan-qrcode.png" alt="步骤2：扫码直连" class="w-full h-36 object-cover object-top rounded-lg" />
          </div>
        </div>

        <!-- 步骤 3 -->
        <div class="bg-white dark:bg-dark-900 rounded-3xl border border-slate-200 dark:border-white/10 p-6 shadow-sm flex flex-col justify-between">
          <div>
            <div class="flex items-center gap-3 mb-3">
              <span class="w-7 h-7 rounded-full bg-purple-600 text-white font-black text-xs flex items-center justify-center">3</span>
              <span class="text-xs font-bold text-purple-600 dark:text-purple-400 font-mono">{{ t('agentDownload.step3', 'Step 03 · 交代任务去睡大觉') }}</span>
            </div>
            <h3 class="text-base font-bold text-slate-900 dark:text-white mb-2">单手交代任务，静待收网</h3>
            <p class="text-xs text-slate-600 dark:text-dark-300 leading-relaxed mb-4">
              在手机上单手打一句话。手机息屏丢一边，魔丸自主干完全流程，推完代码直接看战报。
            </p>
          </div>
          <div class="rounded-xl overflow-hidden border border-slate-200 dark:border-dark-700 bg-slate-50 dark:bg-dark-800 p-2">
            <img src="/assets/promo/mobile-long-task.png" alt="步骤3：自动提交完成" class="w-full h-36 object-cover object-top rounded-lg" />
          </div>
        </div>
      </div>

      <!-- 快速接入配置助手 -->
      <div class="rounded-3xl border border-slate-200 dark:border-white/10 bg-white dark:bg-dark-900 p-6 sm:p-8 shadow-sm mb-12">
        <div class="flex flex-col justify-between gap-4 border-b border-slate-100 dark:border-white/10 pb-6 sm:flex-row sm:items-center">
          <div>
            <div class="flex items-center gap-2">
              <span class="flex h-8 w-8 items-center justify-center rounded-lg bg-orange-100 text-orange-600 dark:bg-orange-950/40 dark:text-orange-400">
                <Icon name="cog" size="sm" />
              </span>
              <h3 class="text-xl font-bold text-slate-900 dark:text-white">
                {{ t('agentDownload.quickStartTitle', '一键接入配置助手') }}
              </h3>
            </div>
            <p class="mt-1 text-xs sm:text-sm text-slate-500 dark:text-dark-400">
              {{ t('agentDownload.quickStartDesc', '客户端安装完成后，填入本站 API 地址与您的专属密钥，即可畅享智能体极速调用。') }}
            </p>
          </div>

          <div class="flex flex-wrap items-center gap-2 text-xs text-slate-600 dark:text-dark-400">
            <span class="rounded-md bg-slate-100 px-2 py-1 dark:bg-dark-800">1. 下载并运行</span>
            <span>→</span>
            <span class="rounded-md bg-orange-50 px-2 py-1 font-medium text-orange-700 dark:bg-orange-950/40 dark:text-orange-300">2. 填入 API 密钥</span>
            <span>→</span>
            <span class="rounded-md bg-slate-100 px-2 py-1 dark:bg-dark-800">3. 手机扫码躺平</span>
          </div>
        </div>

        <div class="mt-6 space-y-5">
          <!-- 基础接口地址 -->
          <div>
            <label class="block text-xs font-semibold uppercase tracking-wider text-slate-500 dark:text-dark-400">
              {{ t('agentDownload.apiBaseUrl', '接口 Base URL') }}
            </label>
            <div class="mt-2 flex items-center gap-2">
              <input
                type="text"
                readonly
                :value="currentApiBaseUrl"
                class="w-full rounded-xl border border-slate-200 bg-slate-50 px-4 py-2.5 font-mono text-xs sm:text-sm text-slate-800 dark:border-dark-700 dark:bg-dark-800 dark:text-dark-200"
              />
              <button
                type="button"
                class="shrink-0 rounded-xl bg-orange-600 px-4 py-2.5 text-xs font-bold text-white transition hover:bg-orange-500 shadow-sm"
                @click="copyText(currentApiBaseUrl)"
              >
                {{ t('agentDownload.copyUrlBtn', '复制地址') }}
              </button>
            </div>
          </div>

          <!-- API Key 选择与获取 -->
          <div v-if="isAuthenticated">
            <label class="block text-xs font-semibold uppercase tracking-wider text-slate-500 dark:text-dark-400">
              {{ t('agentDownload.apiKeySelect', '选择您的专属 API Key') }}
            </label>
            <div class="mt-2 flex flex-col sm:flex-row items-stretch sm:items-center gap-2">
              <select
                v-if="userApiKeys.length > 0"
                v-model="selectedApiKey"
                class="w-full flex-1 rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-xs sm:text-sm text-slate-800 dark:border-dark-700 dark:bg-dark-800 dark:text-dark-200"
              >
                <option v-for="key in userApiKeys" :key="key.id" :value="key.key">
                  {{ key.name || '默认密钥' }} ({{ maskKey(key.key) }})
                </option>
              </select>
              <div v-else class="flex-1 text-xs text-amber-600 dark:text-amber-400 py-2">
                {{ t('agentDownload.noKeyHint', '暂无可用的 API Key，点击右侧前往创建') }}
              </div>
              <router-link
                to="/keys"
                class="inline-flex shrink-0 items-center justify-center rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-xs font-semibold text-slate-700 hover:bg-slate-50 dark:border-dark-700 dark:bg-dark-800 dark:text-dark-200 dark:hover:bg-dark-700"
              >
                {{ t('agentDownload.createKeyBtn', '前往密钥管理') }}
              </router-link>
            </div>
          </div>

          <!-- 未登录引导 -->
          <div v-else class="rounded-xl border border-slate-200 bg-slate-50 p-4 text-xs sm:text-sm text-slate-600 dark:border-white/10 dark:bg-dark-800 dark:text-dark-300 flex flex-col sm:flex-row sm:items-center justify-between gap-3">
            <span>{{ t('agentDownload.loginPrompt', '登录当前站点后，可在此处直接一键获取您的专属 API Key 与一键连通配置。') }}</span>
            <router-link
              to="/login?redirect=/download"
              class="inline-flex shrink-0 items-center justify-center rounded-lg bg-orange-600 px-4 py-1.5 text-xs font-bold text-white shadow hover:bg-orange-500"
            >
              {{ t('agentDownload.loginNow', '立即登录') }}
            </router-link>
          </div>
        </div>
      </div>

      <!-- FAQ 问答卡片 -->
      <div class="bg-white dark:bg-dark-900 p-6 sm:p-8 rounded-3xl border border-slate-200 dark:border-white/10 shadow-sm">
        <div class="flex items-center gap-2 mb-6">
          <span class="text-xl">💬</span>
          <h3 class="text-lg sm:text-xl font-bold text-slate-900 dark:text-white">躺平操作 · 常见清醒问答 (FAQ)</h3>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-3 gap-5 text-left">
          <div class="p-4 rounded-2xl bg-slate-50 dark:bg-dark-800/80 border border-slate-200/80 dark:border-white/5">
            <h4 class="font-bold text-slate-900 dark:text-white text-xs sm:text-sm mb-2">Q: 手机扫码需要额外下载 App 吗？</h4>
            <p class="text-xs text-slate-600 dark:text-dark-300 leading-relaxed">
              <strong>完全不需要。</strong>手机连接同一局域网 Wi-Fi，用自带系统相机或微信/浏览器扫一扫，免密直接在网页端使用，体验与原生 App 一样丝滑。
            </p>
          </div>

          <div class="p-4 rounded-2xl bg-slate-50 dark:bg-dark-800/80 border border-slate-200/80 dark:border-white/5">
            <h4 class="font-bold text-slate-900 dark:text-white text-xs sm:text-sm mb-2">Q: 免费公网穿透需要自己买域名或服务器吗？</h4>
            <p class="text-xs text-slate-600 dark:text-dark-300 leading-relaxed">
              <strong>完全不需要。</strong>魔丸内置官方与 Cloudflare 隧道穿透，点一下开关自动生成专属加密 HTTPS 外网地址，零配置、零费用。
            </p>
          </div>

          <div class="p-4 rounded-2xl bg-slate-50 dark:bg-dark-800/80 border border-slate-200/80 dark:border-white/5">
            <h4 class="font-bold text-slate-900 dark:text-white text-xs sm:text-sm mb-2">Q: 手机息屏或关机了，电脑长任务会停掉吗？</h4>
            <p class="text-xs text-slate-600 dark:text-dark-300 leading-relaxed">
              <strong>绝不会停。</strong>核心计算和 Agent 调度均在电脑本地后台运行，手机只是遥控器。下完指令直接睡大觉，第二天打开手机直接看战报。
            </p>
          </div>
        </div>
      </div>
    </section>

    <!-- ==================== 7. 专治内耗：4 条扎心清醒真相 ==================== -->
    <section id="truths" class="pt-8 border-t border-slate-200/60 dark:border-white/10">
      <div class="text-center max-w-3xl mx-auto mb-12">
        <span class="text-xs uppercase font-extrabold tracking-widest text-pink-600 bg-pink-100 px-3 py-1 rounded-full dark:bg-pink-950/60 dark:text-pink-400">
          Wake Up Truths
        </span>
        <h2 class="text-2xl sm:text-4xl font-black text-slate-950 dark:text-white mt-4 mb-3">
          专治内耗：4 条扎心清醒真相
        </h2>
        <p class="text-slate-600 dark:text-dark-300 text-sm sm:text-base">
          如果你还在每天对着终端按回车，请认真读完这四句话。
        </p>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <!-- 真相 1 -->
        <div class="p-6 rounded-3xl bg-white dark:bg-dark-900 border border-slate-200 dark:border-white/10 shadow-sm">
          <div class="text-xs font-bold text-orange-600 dark:text-orange-400 uppercase tracking-widest mb-1 font-mono">
            Truth 01 / 关于工具选择
          </div>
          <h3 class="text-lg font-black text-slate-900 dark:text-white mb-2">
            “你是在追求技术，还是在给自己找罪受？”
          </h3>
          <p class="text-xs sm:text-sm text-slate-600 dark:text-dark-300 leading-relaxed">
            有些工具光配环境就得折腾半天，跑个两轮就卡死，你还得像伺候大爷一样盯着它。记住：<strong>好的工具是来当牛马的，不是让你给它当全职保姆的。</strong>
          </p>
        </div>

        <!-- 真相 2 -->
        <div class="p-6 rounded-3xl bg-white dark:bg-dark-900 border border-slate-200 dark:border-white/10 shadow-sm">
          <div class="text-xs font-bold text-pink-600 dark:text-pink-400 uppercase tracking-widest mb-1 font-mono">
            Truth 02 / 关于假装敬业
          </div>
          <h3 class="text-lg font-black text-slate-900 dark:text-white mb-2">
            “不敢离开工位，是因为它根本不配让你信任”
          </h3>
          <p class="text-xs sm:text-sm text-slate-600 dark:text-dark-300 leading-relaxed">
            你端坐在电脑前不是因为你有多自律，而是因为你知道只要你一转身，它就会把代码改崩。能让你安心上床睡觉的 Agent，才叫真正的生产力。
          </p>
        </div>

        <!-- 真相 3 -->
        <div class="p-6 rounded-3xl bg-white dark:bg-dark-900 border border-slate-200 dark:border-white/10 shadow-sm">
          <div class="text-xs font-bold text-purple-600 dark:text-purple-400 uppercase tracking-widest mb-1 font-mono">
            Truth 03 / 关于 AI 降智
          </div>
          <h3 class="text-lg font-black text-slate-900 dark:text-white mb-2">
            “别替它的无能找借口了，它就是接不住硬仗”
          </h3>
          <p class="text-xs sm:text-sm text-slate-600 dark:text-dark-300 leading-relaxed">
            遇到 20 轮就上下文崩溃失忆的叫玩具；能在 400+ 分钟、2094 步工具调用里死咬最终目标直到全绿通过的，才配叫工程智能体。
          </p>
        </div>

        <!-- 真相 4 -->
        <div class="p-6 rounded-3xl bg-white dark:bg-dark-900 border border-slate-200 dark:border-white/10 shadow-sm">
          <div class="text-xs font-bold text-emerald-600 dark:text-emerald-400 uppercase tracking-widest mb-1 font-mono">
            Truth 04 / 关于终极掌控感
          </div>
          <h3 class="text-lg font-black text-slate-900 dark:text-white mb-2">
            “最高级的掌控，是躺着看它把活干完”
          </h3>
          <p class="text-xs sm:text-sm text-slate-600 dark:text-dark-300 leading-relaxed">
            人在被窝里，手机一震：‘49 轮长任务已自主结算，所有测试用例已通过，分支已自动提交推送。’——这，才叫真正的掌控。
          </p>
        </div>
      </div>
    </section>

    <!-- ==================== 8. 客户端极速下载专区 ==================== -->
    <section id="download-section" class="pt-8 border-t border-slate-200/60 dark:border-white/10">
      <div class="bg-gradient-to-br from-orange-600 via-amber-600 to-rose-600 p-8 sm:p-12 rounded-3xl text-white text-center shadow-2xl shadow-orange-600/30 mb-12">
        <h2 class="text-2xl sm:text-4xl lg:text-5xl font-black mb-4">
          你还要在工位上假装努力多久？
        </h2>
        <p class="text-white/90 text-sm sm:text-lg mb-8 max-w-2xl mx-auto">
          现在就下载魔丸客户端，体验什么叫“躺平收割”。<br />
          局域网扫码即用 · 免费穿透不限速 · 实测连续几小时狂飙不降智
        </p>

        <div class="flex flex-col sm:flex-row items-center justify-center gap-4">
          <a
            :href="detectedPlatform.primaryUrl"
            download
            class="w-full sm:w-auto px-8 py-4 rounded-2xl text-base font-black text-orange-600 bg-white hover:bg-slate-50 shadow-xl transition active:scale-[0.98]"
          >
            🚀 立即下载 {{ detectedPlatform.platformLabel }} 客户端 ({{ downloadConfig.version }})
          </a>
          <a
            v-if="downloadConfig.windows.portableUrl"
            :href="downloadConfig.windows.portableUrl"
            download
            class="w-full sm:w-auto px-6 py-4 rounded-2xl text-sm font-bold text-white border border-white/30 bg-white/10 hover:bg-white/20 backdrop-blur-sm transition active:scale-[0.98]"
          >
            📦 {{ t('agentDownload.downloadPortable', '免安装便携版 (.zip)') }}
          </a>
        </div>

        <div class="mt-6 flex items-center justify-center gap-2 text-xs text-white/80">
          <span>{{ downloadConfig.envIncludedNote }}</span>
        </div>
      </div>

      <!-- 全平台下载版本矩阵 -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-5">
        <!-- Windows -->
        <div class="p-5 rounded-2xl bg-white dark:bg-dark-900 border border-slate-200 dark:border-white/10 shadow-sm flex flex-col justify-between">
          <div>
            <div class="flex items-center justify-between mb-3">
              <span class="font-black text-slate-900 dark:text-white text-base">{{ t('agentDownload.windowsTitle', 'Windows 客户端') }}</span>
              <span class="text-[10px] font-bold px-2 py-0.5 rounded bg-emerald-100 text-emerald-800 dark:bg-emerald-950/50 dark:text-emerald-300">首选推荐</span>
            </div>
            <p class="text-xs text-slate-500 dark:text-dark-400 mb-4">
              {{ t('agentDownload.windowsDesc', '支持 Windows 10 / 11 64位。支持标准安装包与免安装便携版。') }}
            </p>
          </div>
          <div class="space-y-2">
            <a
              :href="downloadConfig.windows.installerUrl"
              download
              class="block w-full py-2 px-3 text-center rounded-xl bg-orange-600 hover:bg-orange-500 text-white font-bold text-xs shadow-sm transition"
            >
              {{ t('agentDownload.windowsInstaller', '标准安装版 (.exe)') }}
            </a>
            <a
              v-if="downloadConfig.windows.portableUrl"
              :href="downloadConfig.windows.portableUrl"
              download
              class="block w-full py-2 px-3 text-center rounded-xl border border-slate-200 dark:border-dark-700 hover:bg-slate-50 dark:hover:bg-dark-800 text-slate-700 dark:text-dark-200 font-semibold text-xs transition"
            >
              {{ t('agentDownload.windowsPortable', '免安装便携版 (.zip)') }}
            </a>
          </div>
        </div>

        <!-- macOS -->
        <div class="p-5 rounded-2xl bg-white dark:bg-dark-900 border border-slate-200 dark:border-white/10 shadow-sm flex flex-col justify-between">
          <div>
            <div class="flex items-center justify-between mb-3">
              <span class="font-black text-slate-900 dark:text-white text-base">{{ t('agentDownload.macTitle', 'macOS 客户端') }}</span>
              <span class="text-[10px] font-bold px-2 py-0.5 rounded bg-blue-100 text-blue-800 dark:bg-blue-950/50 dark:text-blue-300">Apple Silicon</span>
            </div>
            <p class="text-xs text-slate-500 dark:text-dark-400 mb-4">
              {{ t('agentDownload.macDesc', '适配 Apple Silicon (M1/M2/M3/M4) 与 Intel 架构芯片。') }}
            </p>
          </div>
          <div class="space-y-2">
            <a
              :href="downloadConfig.mac.arm64Url || '#download-section'"
              :download="!!downloadConfig.mac.arm64Url"
              class="block w-full py-2 px-3 text-center rounded-xl bg-slate-800 hover:bg-slate-700 text-white font-bold text-xs shadow-sm transition"
            >
              Apple Silicon (.dmg)
            </a>
            <a
              :href="downloadConfig.mac.x64Url || '#download-section'"
              :download="!!downloadConfig.mac.x64Url"
              class="block w-full py-2 px-3 text-center rounded-xl border border-slate-200 dark:border-dark-700 hover:bg-slate-50 dark:hover:bg-dark-800 text-slate-700 dark:text-dark-200 font-semibold text-xs transition"
            >
              Intel 芯片 (.dmg)
            </a>
          </div>
        </div>

        <!-- Linux -->
        <div class="p-5 rounded-2xl bg-white dark:bg-dark-900 border border-slate-200 dark:border-white/10 shadow-sm flex flex-col justify-between">
          <div>
            <div class="flex items-center justify-between mb-3">
              <span class="font-black text-slate-900 dark:text-white text-base">{{ t('agentDownload.linuxTitle', 'Linux 极客版') }}</span>
              <span class="text-[10px] font-bold px-2 py-0.5 rounded bg-purple-100 text-purple-800 dark:bg-purple-950/50 dark:text-purple-300">CLI / AppImage</span>
            </div>
            <p class="text-xs text-slate-500 dark:text-dark-400 mb-4">
              {{ t('agentDownload.linuxDesc', '适用于 Ubuntu, Debian, CentOS 等发行版或 Docker 容器化运行。') }}
            </p>
          </div>
          <div class="space-y-2">
            <a
              :href="downloadConfig.linux.appImageUrl || '#download-section'"
              :download="!!downloadConfig.linux.appImageUrl"
              class="block w-full py-2 px-3 text-center rounded-xl bg-slate-800 hover:bg-slate-700 text-white font-bold text-xs shadow-sm transition"
            >
              AppImage 运行包
            </a>
            <a
              :href="downloadConfig.linux.tarUrl || '#download-section'"
              :download="!!downloadConfig.linux.tarUrl"
              class="block w-full py-2 px-3 text-center rounded-xl border border-slate-200 dark:border-dark-700 hover:bg-slate-50 dark:hover:bg-dark-800 text-slate-700 dark:text-dark-200 font-semibold text-xs transition"
            >
              tar.gz 归档包
            </a>
          </div>
        </div>

        <!-- 手机端扫码 -->
        <div class="p-5 rounded-2xl bg-white dark:bg-dark-900 border border-slate-200 dark:border-white/10 shadow-sm flex flex-col justify-between">
          <div>
            <div class="flex items-center justify-between mb-3">
              <span class="font-black text-slate-900 dark:text-white text-base">手机端躺平使用</span>
              <span class="text-[10px] font-bold px-2 py-0.5 rounded bg-amber-100 text-amber-800 dark:bg-amber-950/50 dark:text-amber-300">无需安装</span>
            </div>
            <p class="text-xs text-slate-500 dark:text-dark-400 mb-4">
              在电脑客户端点击【扫码直连】，手机微信或浏览器扫码免密直接进入！
            </p>
          </div>
          <div>
            <a
              href="#lay-flat"
              class="block w-full py-2 px-3 text-center rounded-xl bg-gradient-to-r from-orange-600 to-amber-600 hover:opacity-95 text-white font-bold text-xs shadow-sm transition"
            >
              查看扫码使用图解
            </a>
          </div>
        </div>
      </div>
    </section>

    <!-- ==================== 高清图片灯箱 Lightbox Modal ==================== -->
    <Teleport to="body">
      <div
        v-if="lightboxImage"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/85 p-4 backdrop-blur-md transition-opacity duration-200"
        @click="closeLightbox"
      >
        <div class="relative max-h-[92vh] max-w-5xl overflow-hidden rounded-2xl bg-dark-900 shadow-2xl" @click.stop>
          <div class="flex items-center justify-between border-b border-white/10 px-4 py-3 bg-dark-950/80">
            <span class="text-sm font-semibold text-white">{{ lightboxTitle || '高清大图' }}</span>
            <button
              type="button"
              class="rounded-lg p-1.5 text-gray-400 hover:bg-white/10 hover:text-white transition"
              @click="closeLightbox"
            >
              <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
          <div class="p-2 overflow-auto max-h-[82vh]">
            <img :src="lightboxImage" :alt="lightboxTitle" class="w-full h-auto object-contain rounded-lg" />
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import Icon from '@/components/icons/Icon.vue'
import { useAuthStore } from '@/stores/auth'
import { keysAPI } from '@/api/keys'
import type { ApiKey } from '@/types'

defineProps<{
  embedded?: boolean
}>()

const { t } = useI18n()
const authStore = useAuthStore()

const isAuthenticated = computed(() => authStore.isAuthenticated)
const currentApiBaseUrl = computed(() => (typeof window !== 'undefined' ? `${window.location.origin}/v1` : 'https://api.yourdomain.com/v1'))

const userApiKeys = ref<ApiKey[]>([])
const selectedApiKey = ref('')

// 灯箱状态
const lightboxImage = ref<string | null>(null)
const lightboxTitle = ref<string>('')

const openLightbox = (image: string, title?: string) => {
  lightboxImage.value = image
  lightboxTitle.value = title || ''
}

const closeLightbox = () => {
  lightboxImage.value = null
  lightboxTitle.value = ''
}

const handleKeyDown = (e: KeyboardEvent) => {
  if (e.key === 'Escape' && lightboxImage.value) {
    closeLightbox()
  }
}

// 下载包版本与链接配置
const downloadConfig = ref({
  version: 'v2.1.0',
  windows: {
    installerUrl: '/downloads/mowan-harness-setup.exe',
    portableUrl: '/downloads/mowan-harness-portable.zip',
    size: '142 MB',
  },
  mac: {
    arm64Url: '/downloads/mowan-harness-arm64.dmg',
    x64Url: '/downloads/mowan-harness-x64.dmg',
  },
  linux: {
    appImageUrl: '/downloads/mowan-harness.AppImage',
    tarUrl: '/downloads/mowan-harness-linux.tar.gz',
  },
  envIncludedNote: '自带独立运行环境，无需安装 Node.js 或繁琐配置，下载后一键极速运行。',
})

// 智能检测当前操作系统
const detectedPlatform = computed(() => {
  if (typeof window === 'undefined') {
    return {
      name: 'windows',
      platformLabel: 'Windows 64位',
      primaryUrl: downloadConfig.value.windows.installerUrl,
    }
  }

  const ua = window.navigator.userAgent.toLowerCase()
  if (ua.includes('mac') || ua.includes('darwin')) {
    return {
      name: 'mac',
      platformLabel: 'macOS (Apple Silicon / Intel)',
      primaryUrl: downloadConfig.value.mac.arm64Url,
    }
  }
  if (ua.includes('linux')) {
    return {
      name: 'linux',
      platformLabel: 'Linux',
      primaryUrl: downloadConfig.value.linux.appImageUrl,
    }
  }
  return {
    name: 'windows',
    platformLabel: 'Windows 64位',
    primaryUrl: downloadConfig.value.windows.installerUrl,
  }
})

const maskKey = (key: string) => {
  if (!key || key.length < 10) return key
  return `${key.slice(0, 7)}...${key.slice(-4)}`
}

const copyText = async (text: string) => {
  try {
    await navigator.clipboard.writeText(text)
    alert('已复制到剪贴板！')
  } catch {
    // fallback
  }
}

// 加载用户的 API Keys
const fetchUserApiKeys = async () => {
  if (!isAuthenticated.value) return
  try {
    const res = await keysAPI.list(1, 100)
    if (res?.items && Array.isArray(res.items)) {
      userApiKeys.value = res.items
      if (userApiKeys.value.length > 0 && !selectedApiKey.value) {
        selectedApiKey.value = userApiKeys.value[0].key
      }
    }
  } catch {
    // ignore
  }
}

onMounted(() => {
  if (isAuthenticated.value) {
    void fetchUserApiKeys()
  }
  window.addEventListener('keydown', handleKeyDown)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeyDown)
})
</script>
