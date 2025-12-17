<script setup lang="ts">
import { computed, ref } from 'vue'

const props = defineProps<{
  data: Record<string, number>
  days?: number
}>()

interface DayData {
  date: string
  count: number
  level: number
}

const containerRef = ref<HTMLElement>()
const hoveredDay = ref<DayData | null>(null)
const tooltipPosition = ref({ x: 0, y: 0 })

const weeks = computed(() => {
  const result: DayData[][] = []
  const today = new Date()
  const daysCount = props.days || 365

  let currentWeek: DayData[] = []

  for (let i = daysCount - 1; i >= 0; i--) {
    const date = new Date(today)
    date.setDate(date.getDate() - i)
    const dateStr = date.toISOString().split('T')[0] as string
    const count = props.data[dateStr] || 0

    currentWeek.push({
      date: dateStr,
      count,
      level: getLevel(count)
    })

    if (date.getDay() === 6 || i === 0) {
      result.push(currentWeek)
      currentWeek = []
    }
  }

  return result
})

// Calculate total check-ins
const totalCheckins = computed(() => {
  return Object.values(props.data).reduce((sum, count) => sum + count, 0)
})

// Calculate active days
const activeDays = computed(() => {
  return Object.values(props.data).filter(count => count > 0).length
})

function getLevel(count: number): number {
  if (count === 0) return 0
  if (count === 1) return 1
  if (count <= 3) return 2
  if (count <= 5) return 3
  return 4
}

const levelColors = [
  'bg-white/[0.04] hover:bg-white/[0.08]',
  'bg-indigo-900/50 hover:bg-indigo-800/60',
  'bg-indigo-700/60 hover:bg-indigo-600/70',
  'bg-indigo-500/70 hover:bg-indigo-400/80',
  'bg-indigo-400 hover:bg-indigo-300'
]

const levelLabels = ['无', '少量', '一般', '较多', '活跃']

function formatDate(dateStr: string): string {
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    weekday: 'short'
  })
}

function handleMouseEnter(day: DayData, event: MouseEvent) {
  hoveredDay.value = day
  updateTooltipPosition(event)
}

function handleMouseMove(event: MouseEvent) {
  if (hoveredDay.value) {
    updateTooltipPosition(event)
  }
}

function handleMouseLeave() {
  hoveredDay.value = null
}

function updateTooltipPosition(event: MouseEvent) {
  const container = containerRef.value
  if (!container) return

  const rect = container.getBoundingClientRect()
  tooltipPosition.value = {
    x: event.clientX - rect.left,
    y: event.clientY - rect.top - 10
  }
}

const monthLabels = computed(() => {
  const labels: { month: string; index: number }[] = []
  const today = new Date()
  const daysCount = props.days || 365
  let lastMonth = -1

  for (let i = daysCount - 1; i >= 0; i--) {
    const date = new Date(today)
    date.setDate(date.getDate() - i)
    const month = date.getMonth()
    const weekIndex = Math.floor((daysCount - 1 - i) / 7)

    if (month !== lastMonth) {
      labels.push({
        month: date.toLocaleDateString('zh-CN', { month: 'short' }),
        index: weekIndex
      })
      lastMonth = month
    }
  }

  return labels
})
</script>

<template>
  <div ref="containerRef" class="heatmap-wrapper relative">
    <!-- Month labels -->
    <div class="month-labels flex mb-2 text-xs pl-8 relative h-4">
      <div
        v-for="label in monthLabels"
        :key="label.index"
        class="absolute"
        :style="{ left: `${label.index * 14 + 32}px` }"
      >
        {{ label.month }}
      </div>
    </div>

    <!-- Heatmap grid -->
    <div class="flex gap-1 overflow-x-auto pb-2">
      <!-- Day labels -->
      <div class="day-labels flex flex-col gap-1 text-xs pr-2 shrink-0">
        <div class="h-3"></div>
        <div class="h-3 leading-3">一</div>
        <div class="h-3"></div>
        <div class="h-3 leading-3">三</div>
        <div class="h-3"></div>
        <div class="h-3 leading-3">五</div>
        <div class="h-3"></div>
      </div>

      <!-- Weeks -->
      <div class="flex gap-1 min-w-max">
        <div
          v-for="(week, weekIndex) in weeks"
          :key="weekIndex"
          class="flex flex-col gap-1"
        >
          <div
            v-for="day in week"
            :key="day.date"
            :class="levelColors[day.level]"
            class="w-3 h-3 rounded cursor-pointer heatmap-cell relative"
            @mouseenter="handleMouseEnter(day, $event)"
            @mousemove="handleMouseMove"
            @mouseleave="handleMouseLeave"
          ></div>
        </div>
      </div>
    </div>

    <!-- Tooltip -->
    <Transition name="fade">
      <div
        v-if="hoveredDay"
        class="heatmap-tooltip absolute z-50 px-3 py-2 glass-card text-sm pointer-events-none transform -translate-x-1/2 -translate-y-full"
        :style="{
          left: `${tooltipPosition.x}px`,
          top: `${tooltipPosition.y}px`
        }"
      >
        <div class="tooltip-date font-medium">{{ formatDate(hoveredDay.date) }}</div>
        <div class="tooltip-count mt-0.5">
          {{ hoveredDay.count > 0 ? `${hoveredDay.count} 次打卡` : '无打卡记录' }}
        </div>
      </div>
    </Transition>

    <!-- Legend and Stats -->
    <div class="heatmap-footer flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 mt-4 pt-4">
      <!-- Stats -->
      <div class="heatmap-stats flex items-center gap-6 text-sm">
        <div class="stat-item">
          <span class="stat-value">{{ totalCheckins }}</span> 次打卡
        </div>
        <div class="stat-item">
          <span class="stat-value">{{ activeDays }}</span> 天活跃
        </div>
      </div>

      <!-- Legend -->
      <div class="heatmap-legend flex items-center gap-2 text-xs">
        <span>少</span>
        <div class="flex gap-1">
          <div
            v-for="(color, index) in levelColors"
            :key="index"
            :class="color"
            class="w-3 h-3 rounded transition-colors"
            :title="levelLabels[index]"
          ></div>
        </div>
        <span>多</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.heatmap-wrapper {
  color: rgba(255, 255, 255, 0.5);
}

.month-labels {
  color: rgba(255, 255, 255, 0.4);
}

.day-labels {
  color: rgba(255, 255, 255, 0.4);
}

.heatmap-footer {
  border-top: 1px solid rgba(255, 255, 255, 0.05);
}

.heatmap-stats .stat-item {
  color: rgba(255, 255, 255, 0.5);
}

.heatmap-stats .stat-value {
  color: #fff;
  font-weight: 500;
}

.heatmap-legend {
  color: rgba(255, 255, 255, 0.5);
}

.heatmap-tooltip .tooltip-date {
  color: #fff;
}

.heatmap-tooltip .tooltip-count {
  color: rgba(255, 255, 255, 0.6);
}

</style>
