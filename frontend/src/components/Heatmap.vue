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
const tooltipPosition = ref({ x: 0, y: 0, showAbove: true })

// 计算总打卡次数
const totalCheckins = computed(() => {
  return Object.values(props.data || {}).reduce((sum, count) => sum + count, 0)
})

// 计算活跃天数
const activeDays = computed(() => {
  return Object.values(props.data || {}).filter(count => count > 0).length
})

const weeks = computed(() => {
  const result: DayData[][] = []
  const today = new Date()
  const daysCount = props.days || 365

  // 从今天往前推 daysCount 天
  const startDate = new Date(today)
  startDate.setDate(startDate.getDate() - daysCount + 1)

  // 调整到周日开始
  const startDayOfWeek = startDate.getDay()
  startDate.setDate(startDate.getDate() - startDayOfWeek)

  let currentWeek: DayData[] = []
  const endDate = new Date(today)

  for (let d = new Date(startDate); d <= endDate; d.setDate(d.getDate() + 1)) {
    const dateStr = d.toISOString().split('T')[0] || ''
    const count = dateStr ? (props.data?.[dateStr] || 0) : 0

    currentWeek.push({
      date: dateStr,
      count,
      level: getLevel(count)
    })

    // 每周六结束一周（周日=0, 周六=6）
    if (d.getDay() === 6) {
      result.push(currentWeek)
      currentWeek = []
    }
  }

  // 添加最后一个不完整的周
  if (currentWeek.length > 0) {
    result.push(currentWeek)
  }

  return result
})

function getLevel(count: number): number {
  if (count === 0) return 0
  if (count === 1) return 1
  if (count <= 3) return 2
  if (count <= 5) return 3
  return 4
}

const levelColors = [
  'heatmap-level-0',
  'heatmap-level-1',
  'heatmap-level-2',
  'heatmap-level-3',
  'heatmap-level-4'
]

const levelLabels = ['无', '少量', '一般', '较多', '活跃']

function formatDate(dateStr: string): string {
  const date = new Date(dateStr + 'T00:00:00')
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
  const x = event.clientX - rect.left
  const y = event.clientY - rect.top

  // 检查是否靠近顶部，如果是则显示在下方
  const showAbove = y > 60

  tooltipPosition.value = {
    x: Math.max(80, Math.min(x, rect.width - 80)),
    y: showAbove ? y - 10 : y + 20,
    showAbove
  }
}

const monthLabels = computed(() => {
  const labels: { month: string; weekIndex: number }[] = []
  const today = new Date()
  const daysCount = props.days || 365

  const startDate = new Date(today)
  startDate.setDate(startDate.getDate() - daysCount + 1)
  const startDayOfWeek = startDate.getDay()
  startDate.setDate(startDate.getDate() - startDayOfWeek)

  let lastMonth = -1
  let weekIndex = 0
  const endDate = new Date(today)

  for (let d = new Date(startDate); d <= endDate; d.setDate(d.getDate() + 7)) {
    const month = d.getMonth()
    if (month !== lastMonth) {
      labels.push({
        month: d.toLocaleDateString('zh-CN', { month: 'short' }),
        weekIndex
      })
      lastMonth = month
    }
    weekIndex++
  }

  return labels
})
</script>

<template>
  <div ref="containerRef" class="heatmap-container">
    <!-- Stats bar -->
    <div class="heatmap-stats">
      <div class="stat-item">
        <span class="stat-value">{{ totalCheckins }}</span>
        <span class="stat-label">次打卡</span>
      </div>
      <div class="stat-divider"></div>
      <div class="stat-item">
        <span class="stat-value">{{ activeDays }}</span>
        <span class="stat-label">天活跃</span>
      </div>
    </div>

    <!-- Heatmap grid wrapper -->
    <div class="heatmap-scroll">
      <div class="heatmap-grid">
        <!-- Month labels -->
        <div class="month-row">
          <div class="day-label-spacer"></div>
          <div class="months">
            <span
              v-for="label in monthLabels"
              :key="label.weekIndex"
              class="month-label"
              :style="{ left: `${label.weekIndex * 14}px` }"
            >
              {{ label.month }}
            </span>
          </div>
        </div>

        <!-- Grid with day labels -->
        <div class="grid-body">
          <!-- Day labels -->
          <div class="day-labels">
            <span></span>
            <span>一</span>
            <span></span>
            <span>三</span>
            <span></span>
            <span>五</span>
            <span></span>
          </div>

          <!-- Weeks grid -->
          <div class="weeks-grid">
            <div
              v-for="(week, weekIndex) in weeks"
              :key="weekIndex"
              class="week-column"
            >
              <div
                v-for="day in week"
                :key="day.date"
                :class="['day-cell', levelColors[day.level]]"
                @mouseenter="handleMouseEnter(day, $event)"
                @mousemove="handleMouseMove"
                @mouseleave="handleMouseLeave"
              ></div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Tooltip -->
    <Transition name="tooltip">
      <div
        v-if="hoveredDay"
        class="heatmap-tooltip"
        :class="{ 'show-below': !tooltipPosition.showAbove }"
        :style="{
          left: `${tooltipPosition.x}px`,
          top: `${tooltipPosition.y}px`
        }"
      >
        <div class="tooltip-date">{{ formatDate(hoveredDay.date) }}</div>
        <div class="tooltip-count">
          {{ hoveredDay.count > 0 ? `${hoveredDay.count} 次打卡` : '无打卡记录' }}
        </div>
      </div>
    </Transition>

    <!-- Legend -->
    <div class="heatmap-footer">
      <div class="heatmap-legend">
        <span class="legend-label">少</span>
        <div class="legend-colors">
          <div
            v-for="(color, index) in levelColors"
            :key="index"
            :class="['legend-cell', color]"
            :title="levelLabels[index]"
          ></div>
        </div>
        <span class="legend-label">多</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.heatmap-container {
  position: relative;
}

/* Stats bar */
.heatmap-stats {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 20px;
  padding: 12px 16px;
  background: rgba(255, 255, 255, 0.03);
  border-radius: 8px;
  border: 1px solid rgba(255, 255, 255, 0.05);
}

.stat-item {
  display: flex;
  align-items: baseline;
  gap: 4px;
}

.stat-value {
  font-size: 20px;
  font-weight: 700;
  color: #10b981;
}

.stat-label {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.5);
}

.stat-divider {
  width: 1px;
  height: 20px;
  background: rgba(255, 255, 255, 0.1);
}

/* Scroll wrapper */
.heatmap-scroll {
  overflow-x: auto;
  overflow-y: hidden;
  padding-bottom: 8px;
  margin: 0 -4px;
  padding: 0 4px;
}

.heatmap-scroll::-webkit-scrollbar {
  height: 6px;
}

.heatmap-scroll::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.03);
  border-radius: 3px;
}

.heatmap-scroll::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 3px;
}

.heatmap-scroll::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.15);
}

/* Grid layout */
.heatmap-grid {
  min-width: max-content;
}

.month-row {
  display: flex;
  margin-bottom: 4px;
}

.day-label-spacer {
  width: 28px;
  flex-shrink: 0;
}

.months {
  position: relative;
  height: 16px;
  flex: 1;
}

.month-label {
  position: absolute;
  font-size: 11px;
  color: rgba(255, 255, 255, 0.4);
  white-space: nowrap;
}

.grid-body {
  display: flex;
}

.day-labels {
  display: flex;
  flex-direction: column;
  gap: 2px;
  margin-right: 6px;
  flex-shrink: 0;
}

.day-labels span {
  height: 12px;
  line-height: 12px;
  font-size: 10px;
  color: rgba(255, 255, 255, 0.4);
  text-align: right;
  width: 20px;
}

.weeks-grid {
  display: flex;
  gap: 2px;
}

.week-column {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.day-cell {
  width: 12px;
  height: 12px;
  border-radius: 2px;
  cursor: pointer;
  transition: all 0.15s ease;
}

.day-cell:hover {
  transform: scale(1.3);
  z-index: 10;
}

/* Level colors */
.heatmap-level-0 {
  background: rgba(255, 255, 255, 0.04);
}
.heatmap-level-0:hover {
  background: rgba(255, 255, 255, 0.1);
}

.heatmap-level-1 {
  background: rgba(16, 185, 129, 0.25);
}
.heatmap-level-1:hover {
  background: rgba(16, 185, 129, 0.35);
}

.heatmap-level-2 {
  background: rgba(16, 185, 129, 0.45);
}
.heatmap-level-2:hover {
  background: rgba(16, 185, 129, 0.55);
}

.heatmap-level-3 {
  background: rgba(16, 185, 129, 0.65);
}
.heatmap-level-3:hover {
  background: rgba(16, 185, 129, 0.75);
}

.heatmap-level-4 {
  background: #10b981;
}
.heatmap-level-4:hover {
  background: #34d399;
}

/* Tooltip */
.heatmap-tooltip {
  position: absolute;
  z-index: 100;
  padding: 8px 12px;
  background: rgba(15, 23, 42, 0.95);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  backdrop-filter: blur(8px);
  pointer-events: none;
  transform: translate(-50%, -100%);
  white-space: nowrap;
}

.heatmap-tooltip.show-below {
  transform: translate(-50%, 0);
}

.tooltip-date {
  font-size: 13px;
  font-weight: 500;
  color: #fff;
}

.tooltip-count {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.6);
  margin-top: 2px;
}

/* Tooltip animation */
.tooltip-enter-active,
.tooltip-leave-active {
  transition: all 0.15s ease;
}

.tooltip-enter-from,
.tooltip-leave-to {
  opacity: 0;
  transform: translate(-50%, -100%) scale(0.95);
}

.heatmap-tooltip.show-below.tooltip-enter-from,
.heatmap-tooltip.show-below.tooltip-leave-to {
  transform: translate(-50%, 0) scale(0.95);
}

/* Footer / Legend */
.heatmap-footer {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
  padding-top: 12px;
  border-top: 1px solid rgba(255, 255, 255, 0.05);
}

.heatmap-legend {
  display: flex;
  align-items: center;
  gap: 6px;
}

.legend-label {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.4);
}

.legend-colors {
  display: flex;
  gap: 2px;
}

.legend-cell {
  width: 12px;
  height: 12px;
  border-radius: 2px;
}
</style>
