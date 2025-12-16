import { defineStore } from 'pinia'
import { ref, computed, onUnmounted } from 'vue'

export type TrainingPhase = 'idle' | 'contract' | 'hold' | 'relax'

export interface TrainingSettings {
  contractTime: number
  holdTime: number
  relaxTime: number
  cycles: number
}

const DEFAULT_SETTINGS: TrainingSettings = {
  contractTime: 3,
  holdTime: 3,
  relaxTime: 3,
  cycles: 10
}

const STORAGE_KEY = 'tidalcore_training_settings'

function loadSettings(): TrainingSettings {
  try {
    const saved = localStorage.getItem(STORAGE_KEY)
    if (saved) {
      const parsed = JSON.parse(saved)
      return { ...DEFAULT_SETTINGS, ...parsed }
    }
  } catch {
    // ignore
  }
  return { ...DEFAULT_SETTINGS }
}

function saveSettings(settings: TrainingSettings) {
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(settings))
  } catch {
    // ignore
  }
}

export const useTrainingStore = defineStore('training', () => {
  const settings = ref<TrainingSettings>(loadSettings())
  const phase = ref<TrainingPhase>('idle')
  const countdown = ref(0)
  const currentCycle = ref(0)
  const totalDuration = ref(0)
  const isRunning = ref(false)
  const isCompleted = ref(false)

  let timer: ReturnType<typeof setInterval> | null = null
  let durationTimer: ReturnType<typeof setInterval> | null = null

  const progress = computed(() => {
    if (settings.value.cycles === 0) return 0
    return (currentCycle.value / settings.value.cycles) * 100
  })

  const phaseProgress = computed(() => {
    if (phase.value === 'idle') return 0
    const phaseTime = getPhaseTime(phase.value)
    if (phaseTime === 0) return 0
    return ((phaseTime - countdown.value) / phaseTime) * 100
  })

  const phaseLabel = computed(() => {
    const labels: Record<TrainingPhase, string> = {
      idle: '准备开始',
      contract: '收缩',
      hold: '保持',
      relax: '放松'
    }
    return labels[phase.value]
  })

  function getPhaseTime(p: TrainingPhase): number {
    switch (p) {
      case 'contract': return settings.value.contractTime
      case 'hold': return settings.value.holdTime
      case 'relax': return settings.value.relaxTime
      default: return 0
    }
  }

  function nextPhase() {
    const phases: TrainingPhase[] = ['contract', 'hold', 'relax']
    const currentIndex = phases.indexOf(phase.value)

    if (currentIndex === -1) {
      phase.value = 'contract'
    } else if (currentIndex === 2) {
      currentCycle.value++
      if (currentCycle.value >= settings.value.cycles) {
        complete()
        return
      }
      phase.value = 'contract'
    } else {
      phase.value = phases[currentIndex + 1]
    }

    countdown.value = getPhaseTime(phase.value)
  }

  function start() {
    if (isRunning.value) return

    isRunning.value = true
    isCompleted.value = false
    currentCycle.value = 0
    totalDuration.value = 0
    phase.value = 'contract'
    countdown.value = settings.value.contractTime

    timer = setInterval(() => {
      if (countdown.value > 0) {
        countdown.value--
        if (countdown.value === 0) {
          nextPhase()
        }
      }
    }, 1000)

    durationTimer = setInterval(() => {
      totalDuration.value++
    }, 1000)
  }

  function complete() {
    isRunning.value = false
    isCompleted.value = true
    phase.value = 'idle'
    clearTimers()
  }

  function stop() {
    isRunning.value = false
    phase.value = 'idle'
    clearTimers()
  }

  function clearTimers() {
    if (timer) {
      clearInterval(timer)
      timer = null
    }
    if (durationTimer) {
      clearInterval(durationTimer)
      durationTimer = null
    }
  }

  function reset() {
    stop()
    countdown.value = 0
    currentCycle.value = 0
    totalDuration.value = 0
    isCompleted.value = false
  }

  function updateSettings(newSettings: Partial<TrainingSettings>) {
    if (isRunning.value) return

    const updated = { ...settings.value, ...newSettings }
    updated.contractTime = Math.max(1, Math.min(30, updated.contractTime))
    updated.holdTime = Math.max(1, Math.min(30, updated.holdTime))
    updated.relaxTime = Math.max(1, Math.min(30, updated.relaxTime))
    updated.cycles = Math.max(1, Math.min(50, updated.cycles))

    settings.value = updated
    saveSettings(updated)
  }

  return {
    settings,
    phase,
    countdown,
    currentCycle,
    totalDuration,
    isRunning,
    isCompleted,
    progress,
    phaseProgress,
    phaseLabel,
    getPhaseTime,
    start,
    stop,
    reset,
    updateSettings,
    clearTimers
  }
})
