<script setup lang="ts">
import { computed } from 'vue'
import { generateAvatar } from '@/utils/avatar'

interface Props {
  userId?: number | string
  username?: string
  size?: number
  variant?: 'default' | 'gold' | 'silver' | 'bronze'
  borderRadius?: number
}

const props = withDefaults(defineProps<Props>(), {
  size: 40,
  variant: 'default',
  borderRadius: 10,
})

// 使用 userId 作为种子，如果没有则使用 username
const seed = computed(() => props.userId || props.username || 'default')

// 生成头像 Data URI
const avatarSrc = computed(() => generateAvatar(seed.value, props.size))

// 计算样式
const avatarStyle = computed(() => ({
  width: `${props.size}px`,
  height: `${props.size}px`,
  borderRadius: `${props.borderRadius}px`,
}))

// 变体对应的边框和阴影样式
const variantClass = computed(() => {
  switch (props.variant) {
    case 'gold':
      return 'avatar-gold'
    case 'silver':
      return 'avatar-silver'
    case 'bronze':
      return 'avatar-bronze'
    default:
      return ''
  }
})
</script>

<template>
  <div class="user-avatar" :class="variantClass" :style="avatarStyle">
    <img :src="avatarSrc" :alt="username || 'avatar'" class="avatar-img" />
  </div>
</template>

<style scoped>
.user-avatar {
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  flex-shrink: 0;
}

.avatar-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

/* 金色变体 - 第一名 */
.avatar-gold {
  box-shadow: 0 12px 30px rgba(251, 191, 36, 0.4);
  border: 2px solid rgba(251, 191, 36, 0.6);
}

/* 银色变体 - 第二名 */
.avatar-silver {
  box-shadow: 0 10px 25px rgba(148, 163, 184, 0.35);
  border: 2px solid rgba(148, 163, 184, 0.5);
}

/* 铜色变体 - 第三名 */
.avatar-bronze {
  box-shadow: 0 10px 25px rgba(217, 119, 6, 0.35);
  border: 2px solid rgba(217, 119, 6, 0.5);
}
</style>
