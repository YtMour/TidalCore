import { createAvatar } from '@dicebear/core'
import { thumbs } from '@dicebear/collection'

// 海洋主题配色
const oceanColors = [
  '0ea5e9', // sky-500
  '0284c7', // sky-600
  '0369a1', // sky-700
  '22d3ee', // cyan-400
  '06b6d4', // cyan-500
  '0891b2', // cyan-600
  '38bdf8', // sky-400
  '14b8a6', // teal-500
  '2dd4bf', // teal-400
]

/**
 * 生成用户头像 Data URI
 * @param seed 种子值（用户 ID 或用户名）
 * @param size 头像尺寸
 * @returns SVG Data URI
 */
export function generateAvatar(seed: string | number, size = 64): string {
  const avatar = createAvatar(thumbs, {
    seed: String(seed),
    size,
    backgroundColor: oceanColors,
    backgroundType: ['gradientLinear'],
  })
  return avatar.toDataUri()
}

/**
 * 生成用户头像 SVG 字符串
 * @param seed 种子值（用户 ID 或用户名）
 * @param size 头像尺寸
 * @returns SVG 字符串
 */
export function generateAvatarSvg(seed: string | number, size = 64): string {
  const avatar = createAvatar(thumbs, {
    seed: String(seed),
    size,
    backgroundColor: oceanColors,
    backgroundType: ['gradientLinear'],
  })
  return avatar.toString()
}
