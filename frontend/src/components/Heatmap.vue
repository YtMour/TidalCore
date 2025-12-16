<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  data: Record<string, number>
  days?: number
}>()

interface DayData {
  date: string
  count: number
  level: number
}

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

function getLevel(count: number): number {
  if (count === 0) return 0
  if (count === 1) return 1
  if (count <= 3) return 2
  if (count <= 5) return 3
  return 4
}

const levelColors = [
  'bg-white/5',
  'bg-emerald-900/50',
  'bg-emerald-700/60',
  'bg-emerald-500/70',
  'bg-emerald-400'
]

function formatDate(dateStr: string): string {
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN', {
    month: 'short',
    day: 'numeric'
  })
}
</script>

<template>
  <div class="overflow-x-auto">
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
          class="w-3 h-3 rounded-sm cursor-pointer transition-all hover:ring-2 hover:ring-white/30"
          :title="`${formatDate(day.date)}: ${day.count} 次打卡`"
        ></div>
      </div>
    </div>

    <div class="flex items-center justify-end gap-2 mt-4 text-xs text-white/50">
      <span>少</span>
      <div
        v-for="(color, index) in levelColors"
        :key="index"
        :class="color"
        class="w-3 h-3 rounded-sm"
      ></div>
      <span>多</span>
    </div>
  </div>
</template>
