import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export type TrainingPhase = 'idle' | 'contract' | 'hold' | 'relax'
export type DifficultyLevel = 'beginner' | 'easy' | 'medium' | 'hard' | 'random' | 'custom'

export interface TrainingSettings {
  contractTime: number
  holdTime: number
  relaxTime: number
  cycles: number
}

export interface DifficultyPreset {
  id: DifficultyLevel
  name: string
  description: string
  icon: string
  settings: TrainingSettings
}

// 基于凯格尔训练健康建议的难度预设
export const DIFFICULTY_PRESETS: DifficultyPreset[] = [
  {
    id: 'beginner',
    name: '入门',
    description: '适合初学者，轻松上手',
    icon: '🌊',
    settings: { contractTime: 3, holdTime: 2, relaxTime: 4, cycles: 8 }
  },
  {
    id: 'easy',
    name: '简单',
    description: '基础训练，循序渐进',
    icon: '🌴',
    settings: { contractTime: 4, holdTime: 3, relaxTime: 4, cycles: 10 }
  },
  {
    id: 'medium',
    name: '中等',
    description: '标准训练，稳步提升',
    icon: '⚡',
    settings: { contractTime: 5, holdTime: 5, relaxTime: 5, cycles: 12 }
  },
  {
    id: 'hard',
    name: '困难',
    description: '进阶挑战，强化训练',
    icon: '🔥',
    settings: { contractTime: 8, holdTime: 8, relaxTime: 6, cycles: 15 }
  }
]

const DEFAULT_SETTINGS: TrainingSettings = {
  contractTime: 3,
  holdTime: 3,
  relaxTime: 3,
  cycles: 10
}

const STORAGE_KEY = 'tidalcore_training_settings'
const DIFFICULTY_KEY = 'tidalcore_difficulty'

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

function loadDifficulty(): DifficultyLevel {
  try {
    const saved = localStorage.getItem(DIFFICULTY_KEY)
    if (saved && ['beginner', 'easy', 'medium', 'hard', 'random', 'custom'].includes(saved)) {
      return saved as DifficultyLevel
    }
  } catch {
    // ignore
  }
  return 'easy' // 默认简单难度
}

function saveDifficulty(difficulty: DifficultyLevel) {
  try {
    localStorage.setItem(DIFFICULTY_KEY, difficulty)
  } catch {
    // ignore
  }
}

function generateRandomSettings(): TrainingSettings {
  // 随机生成合理范围内的训练参数
  return {
    contractTime: Math.floor(Math.random() * 6) + 3, // 3-8秒
    holdTime: Math.floor(Math.random() * 6) + 2,     // 2-7秒
    relaxTime: Math.floor(Math.random() * 4) + 3,    // 3-6秒
    cycles: Math.floor(Math.random() * 8) + 8        // 8-15次
  }
}

export const useTrainingStore = defineStore('training', () => {
  const settings = ref<TrainingSettings>(loadSettings())
  const difficulty = ref<DifficultyLevel>(loadDifficulty())
  const phase = ref<TrainingPhase>('idle')
  const countdown = ref(0)
  const currentCycle = ref(0)
  const totalDuration = ref(0)
  const isRunning = ref(false)
  const isCompleted = ref(false)

  let timer: ReturnType<typeof setInterval> | null = null
  let durationTimer: ReturnType<typeof setInterval> | null = null

  // 当前难度的预设信息
  const currentPreset = computed(() => {
    if (difficulty.value === 'custom' || difficulty.value === 'random') {
      return null
    }
    return DIFFICULTY_PRESETS.find(p => p.id === difficulty.value) || null
  })

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
      const nextIdx = currentIndex + 1
      phase.value = phases[nextIdx] as TrainingPhase
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

    // 用户手动修改设置时，切换到自定义模式
    difficulty.value = 'custom'
    saveDifficulty('custom')
  }

  function setDifficulty(level: DifficultyLevel) {
    if (isRunning.value) return

    difficulty.value = level
    saveDifficulty(level)

    if (level === 'random') {
      // 随机模式：生成随机参数
      const randomSettings = generateRandomSettings()
      settings.value = randomSettings
      saveSettings(randomSettings)
    } else if (level !== 'custom') {
      // 预设难度：应用预设参数
      const preset = DIFFICULTY_PRESETS.find(p => p.id === level)
      if (preset) {
        settings.value = { ...preset.settings }
        saveSettings(preset.settings)
      }
    }
    // custom 模式保持当前设置不变
  }

  return {
    settings,
    difficulty,
    currentPreset,
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
    setDifficulty,
    clearTimers
  }
})
