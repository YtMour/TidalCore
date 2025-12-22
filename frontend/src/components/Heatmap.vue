<script setup lang="ts">
import { computed, ref } from 'vue'

const props = withDefaults(defineProps<{
  data: Record<string, number>
  days?: number
  label?: string
}>(), {
  days: 365,
  label: '打卡'
})

interface DayData {
  date: string
  count: number
  level: number
  dayOfWeek: number
}

const containerRef = ref<HTMLElement>()
const hoveredDay = ref<DayData | null>(null)
const tooltipPosition = ref({ x: 0, y: 0, showAbove: true })

// 标准化日期键 - 支持多种格式
function normalizeDataKeys(data: Record<string, number>): Record<string, number> {
  const normalized: Record<string, number> = {}
  for (const [key, value] of Object.entries(data)) {
    // 处理 ISO 格式 "2025-12-17T00:00:00+08:00" 或简单格式 "2025-12-17"
    const dateOnly = key.split('T')[0] ?? key
    normalized[dateOnly] = (normalized[dateOnly] || 0) + value
  }
  return normalized
}

// 格式化日期为 YYYY-MM-DD（本地时区）
function formatDateKey(date: Date): string {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

// 标准化后的数据
const normalizedData = computed(() => normalizeDataKeys(props.data || {}))

// 生成热力图网格数据
const weeks = computed(() => {
  const result: DayData[][] = []
  const today = new Date()
  today.setHours(0, 0, 0, 0)

  const daysCount = props.days || 365

  // 计算开始日期
  const startDate = new Date(today)
  startDate.setDate(startDate.getDate() - daysCount + 1)

  // 调整到该周的周日（周日=0）
  const dayOfWeek = startDate.getDay()
  startDate.setDate(startDate.getDate() - dayOfWeek)

  let currentWeek: DayData[] = []
  const current = new Date(startDate)

  while (current <= today) {
    const dateKey = formatDateKey(current)
    const count = normalizedData.value[dateKey] || 0

    currentWeek.push({
      date: dateKey,
      count,
      level: getLevel(count),
      dayOfWeek: current.getDay()
    })

    // 周六结束当前周
    if (current.getDay() === 6) {
      result.push(currentWeek)
      currentWeek = []
    }

    current.setDate(current.getDate() + 1)
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

const levelLabels = ['无记录', '1次', '2-3次', '4-5次', '5次以上']

function formatDisplayDate(dateStr: string): string {
  const parts = dateStr.split('-').map(Number)
  const year = parts[0] ?? 0
  const month = parts[1] ?? 1
  const day = parts[2] ?? 1
  const date = new Date(year, month - 1, day)
  const weekdays = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']
  return `${year}年${month}月${day}日 ${weekdays[date.getDay()]}`
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

  const showAbove = y > 80

  // 限制 tooltip 在容器内，增加边距防止溢出
  const tooltipWidth = 140
  const minX = tooltipWidth / 2 + 10
  const maxX = rect.width - tooltipWidth / 2 - 10

  tooltipPosition.value = {
    x: Math.max(minX, Math.min(x, maxX)),
    y: showAbove ? y - 8 : y + 24,
    showAbove
  }
}

// 计算每个月份标签的位置（按周计算）
const monthPositions = computed(() => {
  const positions: { month: string; weekIndex: number }[] = []
  let lastMonth = -1
  const monthNames = ['1月', '2月', '3月', '4月', '5月', '6月', '7月', '8月', '9月', '10月', '11月', '12月']

  weeks.value.forEach((week, weekIndex) => {
    if (week.length > 0) {
      const firstDay = week[0]
      if (firstDay) {
        const parts = firstDay.date.split('-').map(Number)
        const month = parts[1] ?? 1
        if (month !== lastMonth) {
          positions.push({
            month: monthNames[month - 1] ?? '1月',
            weekIndex
          })
          lastMonth = month
        }
      }
    }
  })

  return positions
})

// 计算总周数用于布局
const totalWeeks = computed(() => weeks.value.length)
</script>

<template>
  <div ref="containerRef" class="heatmap-wrap">
    <!-- 热力图主体 -->
    <div class="hm-main">
      <!-- 月份标签行 -->
      <div class="month-row">
        <div class="weekday-placeholder"></div>
        <div class="month-track">
          <span
            v-for="(m, idx) in monthPositions"
            :key="idx"
            class="month-label"
            :style="{ left: `calc(${(m.weekIndex / totalWeeks) * 100}%)` }"
          >{{ m.month }}</span>
        </div>
      </div>

      <!-- 网格区域 -->
      <div class="grid-area">
        <!-- 星期标签 -->
        <div class="weekday-col">
          <span>日</span>
          <span>一</span>
          <span>二</span>
          <span>三</span>
          <span>四</span>
          <span>五</span>
          <span>六</span>
        </div>

        <!-- 热力图网格 -->
        <div class="grid-box">
          <div
            v-for="(week, wIdx) in weeks"
            :key="wIdx"
            class="week"
          >
            <div
              v-for="day in week"
              :key="day.date"
              class="cell"
              :class="`lv-${day.level}`"
              @mouseenter="handleMouseEnter(day, $event)"
              @mousemove="handleMouseMove"
              @mouseleave="handleMouseLeave"
            />
          </div>
        </div>
      </div>
    </div>

    <!-- Tooltip -->
    <Transition name="tip">
      <div
        v-if="hoveredDay"
        class="hm-tooltip"
        :class="{ 'below': !tooltipPosition.showAbove }"
        :style="{ left: `${tooltipPosition.x}px`, top: `${tooltipPosition.y}px` }"
      >
        <div class="tip-date">{{ formatDisplayDate(hoveredDay.date) }}</div>
        <div class="tip-count" :class="{ active: hoveredDay.count > 0 }">
          {{ hoveredDay.count > 0 ? `${hoveredDay.count} 次${props.label}` : `无${props.label}记录` }}
        </div>
      </div>
    </Transition>

    <!-- 图例 -->
    <div class="hm-legend">
      <span class="lg-text">少</span>
      <div class="lg-items">
        <div
          v-for="i in 5"
          :key="i"
          class="lg-box"
          :class="`lv-${i - 1}`"
          :title="levelLabels[i - 1]"
        />
      </div>
      <span class="lg-text">多</span>
    </div>
  </div>
</template>

<style scoped>
.heatmap-wrap {
  position: relative;
  overflow: visible;
}

/* 热力图主体 */
.hm-main {
  margin-bottom: 12px;
  overflow: hidden;
}

/* 月份行 */
.month-row {
  display: flex;
  margin-bottom: 6px;
}

.weekday-placeholder {
  width: 24px;
  flex-shrink: 0;
}

.month-track {
  position: relative;
  height: 18px;
  flex: 1;
}

.month-label {
  position: absolute;
  font-size: 11px;
  color: rgba(255, 255, 255, 0.6);
  white-space: nowrap;
  transform: translateX(0);
}

/* 网格区域 */
.grid-area {
  display: flex;
  width: 100%;
}

.weekday-col {
  display: flex;
  flex-direction: column;
  width: 24px;
  flex-shrink: 0;
  gap: 3px;
}

.weekday-col span {
  height: 14px;
  line-height: 14px;
  font-size: 10px;
  color: rgba(255, 255, 255, 0.55);
  text-align: right;
  padding-right: 6px;
}

/* 网格容器 - 使用flex自动填充 */
.grid-box {
  display: flex;
  flex: 1;
  gap: 3px;
  justify-content: space-between;
}

.week {
  display: flex;
  flex-direction: column;
  gap: 3px;
  flex: 1;
  max-width: 14px;
}

.cell {
  width: 100%;
  aspect-ratio: 1;
  max-width: 14px;
  max-height: 14px;
  border-radius: 3px;
  cursor: pointer;
  transition: transform 0.1s, box-shadow 0.1s;
}

.cell:hover {
  transform: scale(1.4);
  z-index: 10;
  position: relative;
}

/* 等级颜色 */
.lv-0 {
  background: rgba(255, 255, 255, 0.06);
}
.lv-0:hover {
  background: rgba(255, 255, 255, 0.12);
  box-shadow: 0 0 6px rgba(255, 255, 255, 0.1);
}

.lv-1 {
  background: rgba(16, 185, 129, 0.35);
}
.lv-1:hover {
  background: rgba(16, 185, 129, 0.5);
  box-shadow: 0 0 8px rgba(16, 185, 129, 0.4);
}

.lv-2 {
  background: rgba(16, 185, 129, 0.55);
}
.lv-2:hover {
  background: rgba(16, 185, 129, 0.7);
  box-shadow: 0 0 8px rgba(16, 185, 129, 0.5);
}

.lv-3 {
  background: rgba(16, 185, 129, 0.75);
}
.lv-3:hover {
  background: rgba(16, 185, 129, 0.9);
  box-shadow: 0 0 10px rgba(16, 185, 129, 0.6);
}

.lv-4 {
  background: #10b981;
}
.lv-4:hover {
  background: #34d399;
  box-shadow: 0 0 12px rgba(16, 185, 129, 0.7);
}

/* Tooltip */
.hm-tooltip {
  position: absolute;
  z-index: 100;
  padding: 10px 14px;
  background: rgba(10, 15, 25, 0.95);
  border: 1px solid rgba(255, 255, 255, 0.12);
  border-radius: 8px;
  backdrop-filter: blur(8px);
  pointer-events: none;
  transform: translate(-50%, -100%);
  white-space: nowrap;
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.4);
}

.hm-tooltip.below {
  transform: translate(-50%, 0);
}

.tip-date {
  font-size: 12px;
  font-weight: 600;
  color: #fff;
  margin-bottom: 3px;
}

.tip-count {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.55);
}

.tip-count.active {
  color: #10b981;
  font-weight: 500;
}

/* Tooltip 动画 */
.tip-enter-active,
.tip-leave-active {
  transition: all 0.12s ease;
}

.tip-enter-from,
.tip-leave-to {
  opacity: 0;
  transform: translate(-50%, -100%) scale(0.92);
}

.hm-tooltip.below.tip-enter-from,
.hm-tooltip.below.tip-leave-to {
  transform: translate(-50%, 0) scale(0.92);
}

/* 图例 */
.hm-legend {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 8px;
  margin-top: 8px;
}

.lg-text {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.55);
}

.lg-items {
  display: flex;
  gap: 4px;
}

.lg-box {
  width: 14px;
  height: 14px;
  border-radius: 3px;
}
</style>
