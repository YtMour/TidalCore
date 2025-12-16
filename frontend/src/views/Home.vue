<script setup lang="ts">
import { ref, onMounted } from 'vue'
import MainLayout from '@/layouts/MainLayout.vue'
import Heatmap from '@/components/Heatmap.vue'
import { Card, Skeleton } from '@/components/ui'
import { getGlobalHeatmap } from '@/api/checkin'
import { Timer, Flame, Trophy, Shield, Zap, Heart, Sparkles, ArrowRight, Star, Users } from 'lucide-vue-next'

const heatmapData = ref<Record<string, number>>()
const loading = ref(true)
const mounted = ref(false)

onMounted(async () => {
  setTimeout(() => {
    mounted.value = true
  }, 100)

  try {
    heatmapData.value = await getGlobalHeatmap(365)
  } catch {
    // 静默处理错误
  } finally {
    loading.value = false
  }
})

const features = [
  {
    icon: Timer,
    title: '科学计时',
    description: '自定义收缩-保持-放松循环，专业训练节奏引导',
    color: 'text-violet-400',
    bg: 'bg-gradient-to-br from-violet-500/20 to-purple-500/10',
    glow: 'group-hover:shadow-violet-500/20'
  },
  {
    icon: Flame,
    title: '连续打卡',
    description: '记录你的坚持轨迹，见证每一天的进步与成长',
    color: 'text-orange-400',
    bg: 'bg-gradient-to-br from-orange-500/20 to-amber-500/10',
    glow: 'group-hover:shadow-orange-500/20'
  },
  {
    icon: Trophy,
    title: '毅力排行',
    description: '匿名排行榜，与全站用户一起坚持，互相激励',
    color: 'text-amber-400',
    bg: 'bg-gradient-to-br from-amber-500/20 to-yellow-500/10',
    glow: 'group-hover:shadow-amber-500/20'
  },
  {
    icon: Shield,
    title: '隐私保护',
    description: '仅需用户名密码，不收集任何敏感个人信息',
    color: 'text-emerald-400',
    bg: 'bg-gradient-to-br from-emerald-500/20 to-teal-500/10',
    glow: 'group-hover:shadow-emerald-500/20'
  },
  {
    icon: Zap,
    title: '轻量高效',
    description: '每天只需几分钟，随时随地开始训练',
    color: 'text-cyan-400',
    bg: 'bg-gradient-to-br from-cyan-500/20 to-blue-500/10',
    glow: 'group-hover:shadow-cyan-500/20'
  },
  {
    icon: Heart,
    title: '健康习惯',
    description: '科学的训练方法，帮助建立长期健康习惯',
    color: 'text-pink-400',
    bg: 'bg-gradient-to-br from-pink-500/20 to-rose-500/10',
    glow: 'group-hover:shadow-pink-500/20'
  }
]

const stats = [
  { value: '10K+', label: '活跃用户', icon: Users },
  { value: '100K+', label: '累计打卡', icon: Flame },
  { value: '365', label: '最长连续', icon: Trophy },
  { value: '4.9', label: '用户评分', icon: Star }
]
</script>

<template>
  <MainLayout>
    <div class="space-y-24">
      <!-- Hero Section -->
      <section class="relative py-20 md:py-32 text-center overflow-hidden">
        <!-- Enhanced Decorative elements -->
        <div class="absolute inset-0 overflow-hidden pointer-events-none">
          <!-- Animated gradient orbs -->
          <div class="absolute top-0 left-1/4 w-96 h-96 bg-gradient-to-br from-violet-500/30 to-purple-500/20 rounded-full blur-3xl animate-float-slow"></div>
          <div class="absolute bottom-0 right-1/4 w-80 h-80 bg-gradient-to-br from-pink-500/25 to-rose-500/15 rounded-full blur-3xl animate-float-medium" style="animation-delay: -2s;"></div>
          <div class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[600px] h-[600px] bg-gradient-to-br from-cyan-500/10 to-blue-500/5 rounded-full blur-3xl"></div>

          <!-- Decorative rings -->
          <div class="absolute top-20 right-20 w-32 h-32 border border-white/5 rounded-full animate-pulse-soft"></div>
          <div class="absolute bottom-32 left-16 w-24 h-24 border border-violet-500/10 rounded-full animate-float-fast"></div>

          <!-- Sparkle decorations -->
          <Sparkles class="absolute top-1/4 right-1/3 w-6 h-6 text-violet-400/40 animate-pulse" />
          <Star class="absolute bottom-1/3 left-1/4 w-4 h-4 text-pink-400/30 animate-float-fast" />
        </div>

        <div class="relative z-10">
          <!-- Logo with enhanced animation -->
          <div
            class="inline-block mb-8 transition-all duration-1000"
            :class="mounted ? 'opacity-100 translate-y-0 scale-100' : 'opacity-0 translate-y-8 scale-90'"
          >
            <div class="relative">
              <span class="text-8xl md:text-9xl animate-float inline-block drop-shadow-2xl">🌊</span>
              <!-- Glow effect behind emoji -->
              <div class="absolute inset-0 blur-2xl bg-gradient-to-br from-cyan-400/30 to-blue-500/20 rounded-full scale-150 -z-10"></div>
            </div>
          </div>

          <!-- Title with enhanced gradient -->
          <h1
            class="text-display mb-6 transition-all duration-1000 delay-150"
            :class="mounted ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-8'"
          >
            <span class="gradient-text text-glow-purple">TidalCore</span>
          </h1>

          <!-- Tagline badge -->
          <div
            class="inline-flex items-center gap-2 px-4 py-2 rounded-lg bg-white/[0.03] border border-white/[0.08] mb-6 transition-all duration-700 delay-200"
            :class="mounted ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-8'"
          >
            <span class="w-2 h-2 rounded-full bg-emerald-400 animate-pulse"></span>
            <span class="text-sm text-white/70">开源 · 免费 · 隐私优先</span>
          </div>

          <!-- Subtitle -->
          <p
            class="text-lg md:text-xl text-white/60 mb-12 max-w-2xl mx-auto leading-relaxed transition-all duration-1000 delay-300"
            :class="mounted ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-8'"
          >
            专业的盆底肌训练平台，帮助你建立健康的训练习惯。
            <br class="hidden sm:block" />
            <span class="text-white/80">每天几分钟，坚持就是胜利。</span>
          </p>

          <!-- CTA Buttons -->
          <div
            class="flex flex-col sm:flex-row justify-center gap-4 transition-all duration-1000 delay-400"
            :class="mounted ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-8'"
          >
            <RouterLink
              to="/train"
              class="group relative btn-gradient btn-gradient-glow px-8 py-3.5 rounded-xl text-white font-semibold text-lg overflow-hidden"
            >
              <span class="relative z-10 flex items-center justify-center gap-2">
                <Timer class="w-5 h-5" />
                开始训练
                <ArrowRight class="w-4 h-4 group-hover:translate-x-1 transition-transform" />
              </span>
            </RouterLink>
            <RouterLink
              to="/leaderboard"
              class="group px-8 py-3.5 rounded-xl bg-white/[0.04] hover:bg-white/[0.08] text-white font-semibold text-lg transition-all duration-200 border border-white/[0.08] hover:border-white/[0.15] flex items-center justify-center gap-2"
            >
              <Trophy class="w-5 h-5 text-amber-400" />
              查看排行
            </RouterLink>
          </div>

          <!-- Stats Row -->
          <div
            class="grid grid-cols-2 md:grid-cols-4 gap-4 md:gap-8 mt-16 max-w-3xl mx-auto transition-all duration-1000 delay-500"
            :class="mounted ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-8'"
          >
            <div
              v-for="(stat, index) in stats"
              :key="stat.label"
              class="text-center p-4 rounded-xl bg-white/[0.02] border border-white/[0.06] hover:border-white/[0.1] transition-all duration-200 hover:bg-white/[0.04]"
              :style="{ transitionDelay: `${500 + index * 100}ms` }"
            >
              <component :is="stat.icon" class="w-5 h-5 mx-auto mb-2 text-violet-400" />
              <div class="text-2xl md:text-3xl font-bold text-white mb-1">{{ stat.value }}</div>
              <div class="text-xs text-white/50">{{ stat.label }}</div>
            </div>
          </div>
        </div>
      </section>

      <!-- Global Heatmap Section -->
      <section
        class="transition-all duration-1000 delay-600"
        :class="mounted ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-8'"
      >
        <Card padding="lg" class="relative overflow-hidden">
          <!-- Decorative gradient -->
          <div class="absolute top-0 right-0 w-64 h-64 bg-gradient-to-bl from-violet-500/10 to-transparent rounded-full blur-3xl pointer-events-none"></div>

          <div class="relative">
            <div class="flex items-center justify-between mb-8">
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 rounded-lg bg-gradient-to-br from-emerald-500/20 to-teal-500/10 flex items-center justify-center">
                  <Flame class="w-5 h-5 text-emerald-400" />
                </div>
                <div>
                  <h2 class="text-xl font-bold text-white">全站打卡热力图</h2>
                  <p class="text-sm text-white/40">社区成员的坚持轨迹</p>
                </div>
              </div>
              <span class="badge flex items-center gap-1.5">
                <span class="w-1.5 h-1.5 rounded-full bg-emerald-400 animate-pulse"></span>
                实时更新
              </span>
            </div>

            <div v-if="loading" class="py-8 space-y-4">
              <Skeleton height="1rem" width="60%" />
              <div class="flex gap-1">
                <Skeleton v-for="i in 20" :key="i" width="3rem" height="5rem" />
              </div>
            </div>
            <Heatmap v-else :data="heatmapData || {}" :days="365" />
          </div>
        </Card>
      </section>

      <!-- Features Section -->
      <section
        class="transition-all duration-1000 delay-700"
        :class="mounted ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-8'"
      >
        <div class="text-center mb-16">
          <div class="inline-flex items-center gap-2 px-3 py-1.5 rounded-lg bg-violet-500/10 border border-violet-500/20 mb-4">
            <Sparkles class="w-4 h-4 text-violet-400" />
            <span class="text-sm text-violet-300">核心功能</span>
          </div>
          <h2 class="text-heading text-white mb-4">为什么选择 TidalCore</h2>
          <p class="text-white/50 max-w-xl mx-auto">
            简单、专注、有效的训练体验，帮助你养成健康习惯
          </p>
        </div>

        <div class="grid sm:grid-cols-2 lg:grid-cols-3 gap-6">
          <div
            v-for="(feature, index) in features"
            :key="feature.title"
            class="group glass-card glass-card-hover p-6 transition-all duration-500 hover:shadow-xl"
            :class="feature.glow"
            :style="{ transitionDelay: `${index * 100}ms` }"
          >
            <div
              class="w-14 h-14 rounded-xl flex items-center justify-center mb-5 transition-transform duration-300 group-hover:scale-110"
              :class="feature.bg"
            >
              <component :is="feature.icon" class="w-7 h-7" :class="feature.color" />
            </div>
            <h3 class="text-lg font-semibold text-white mb-2 group-hover:text-violet-200 transition-colors">{{ feature.title }}</h3>
            <p class="text-white/50 text-sm leading-relaxed">{{ feature.description }}</p>
          </div>
        </div>
      </section>

      <!-- How It Works Section -->
      <section
        class="transition-all duration-1000 delay-800"
        :class="mounted ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-8'"
      >
        <div class="text-center mb-16">
          <div class="inline-flex items-center gap-2 px-3 py-1.5 rounded-lg bg-pink-500/10 border border-pink-500/20 mb-4">
            <Zap class="w-4 h-4 text-pink-400" />
            <span class="text-sm text-pink-300">快速上手</span>
          </div>
          <h2 class="text-heading text-white mb-4">三步开始训练</h2>
          <p class="text-white/50 max-w-xl mx-auto">
            无需复杂设置，立即开始你的健康之旅
          </p>
        </div>

        <div class="grid md:grid-cols-3 gap-8">
          <div class="relative text-center">
            <div class="w-16 h-16 mx-auto mb-6 rounded-xl bg-gradient-to-br from-violet-500 to-purple-600 flex items-center justify-center text-2xl font-bold text-white shadow-lg shadow-violet-500/30">
              1
            </div>
            <h3 class="text-lg font-semibold text-white mb-2">注册账号</h3>
            <p class="text-white/50 text-sm">仅需用户名和密码，30秒完成注册</p>
            <!-- Connector line -->
            <div class="hidden md:block absolute top-8 left-[60%] w-[80%] h-px bg-gradient-to-r from-violet-500/50 to-transparent"></div>
          </div>

          <div class="relative text-center">
            <div class="w-16 h-16 mx-auto mb-6 rounded-xl bg-gradient-to-br from-pink-500 to-rose-600 flex items-center justify-center text-2xl font-bold text-white shadow-lg shadow-pink-500/30">
              2
            </div>
            <h3 class="text-lg font-semibold text-white mb-2">开始训练</h3>
            <p class="text-white/50 text-sm">跟随计时器节奏，完成每日训练</p>
            <!-- Connector line -->
            <div class="hidden md:block absolute top-8 left-[60%] w-[80%] h-px bg-gradient-to-r from-pink-500/50 to-transparent"></div>
          </div>

          <div class="text-center">
            <div class="w-16 h-16 mx-auto mb-6 rounded-xl bg-gradient-to-br from-amber-500 to-orange-600 flex items-center justify-center text-2xl font-bold text-white shadow-lg shadow-amber-500/30">
              3
            </div>
            <h3 class="text-lg font-semibold text-white mb-2">坚持打卡</h3>
            <p class="text-white/50 text-sm">记录进度，冲击排行榜</p>
          </div>
        </div>
      </section>

      <!-- Bottom CTA -->
      <section class="text-center py-16">
        <Card padding="lg" class="max-w-3xl mx-auto relative overflow-hidden">
          <!-- Decorative elements -->
          <div class="absolute top-0 left-0 w-full h-full overflow-hidden pointer-events-none">
            <div class="absolute -top-24 -left-24 w-48 h-48 bg-gradient-to-br from-violet-500/20 to-transparent rounded-full blur-3xl"></div>
            <div class="absolute -bottom-24 -right-24 w-48 h-48 bg-gradient-to-br from-pink-500/20 to-transparent rounded-full blur-3xl"></div>
          </div>

          <div class="relative">
            <div class="w-16 h-16 mx-auto mb-6 rounded-xl bg-gradient-to-br from-violet-500 to-pink-500 flex items-center justify-center shadow-lg shadow-violet-500/30">
              <Heart class="w-8 h-8 text-white" />
            </div>
            <h3 class="text-2xl md:text-3xl font-bold text-white mb-4">准备好开始了吗？</h3>
            <p class="text-white/50 mb-8 max-w-md mx-auto">
              加入数千名用户，开始你的健康训练之旅。完全免费，永久开源。
            </p>
            <div class="flex flex-col sm:flex-row justify-center gap-4">
              <RouterLink
                to="/register"
                class="group btn-gradient btn-gradient-glow px-8 py-3.5 rounded-xl text-white font-semibold text-lg"
              >
                <span class="flex items-center justify-center gap-2">
                  立即注册，免费使用
                  <ArrowRight class="w-5 h-5 group-hover:translate-x-1 transition-transform" />
                </span>
              </RouterLink>
              <RouterLink
                to="/train"
                class="px-8 py-3.5 rounded-xl bg-white/[0.04] hover:bg-white/[0.08] text-white font-semibold text-lg transition-all duration-200 border border-white/[0.08] hover:border-white/[0.15]"
              >
                先体验一下
              </RouterLink>
            </div>
          </div>
        </Card>
      </section>
    </div>
  </MainLayout>
</template>
