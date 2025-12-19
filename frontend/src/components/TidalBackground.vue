<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'
import * as THREE from 'three'

const props = defineProps<{
  phase?: 'idle' | 'contract' | 'hold' | 'relax'
  intensity?: number
  isActive?: boolean
}>()

const containerRef = ref<HTMLDivElement | null>(null)

// Three.js objects
let scene: THREE.Scene
let camera: THREE.PerspectiveCamera
let renderer: THREE.WebGLRenderer
let particles: THREE.Points
let waveGeometry: THREE.PlaneGeometry
let waveMesh: THREE.Mesh
let coreMesh: THREE.Mesh // 潮汐核心
let ringMeshes: THREE.Mesh[] = [] // 能量环
let animationId: number
let clock: THREE.Clock

// TidalCore 品牌色系 - 深海蓝绿渐变
const phaseColors = {
  idle: { primary: 0x0ea5e9, secondary: 0x0284c7, accent: 0x38bdf8 },      // 宁静海蓝
  contract: { primary: 0xf43f5e, secondary: 0xbe123c, accent: 0xfb7185 },  // 张力珊瑚红
  hold: { primary: 0xf59e0b, secondary: 0xd97706, accent: 0xfbbf24 },      // 蓄力琥珀
  relax: { primary: 0x10b981, secondary: 0x059669, accent: 0x34d399 }      // 舒缓翠绿
}

// Current target colors for smooth transitions
let currentPrimaryColor = new THREE.Color(phaseColors.idle.primary)
let targetPrimaryColor = new THREE.Color(phaseColors.idle.primary)
let currentSecondaryColor = new THREE.Color(phaseColors.idle.secondary)
let targetSecondaryColor = new THREE.Color(phaseColors.idle.secondary)
let currentAccentColor = new THREE.Color(phaseColors.idle.accent)
let targetAccentColor = new THREE.Color(phaseColors.idle.accent)

// ===== 潮汐波浪着色器 - 增强版 =====
const waveVertexShader = `
  uniform float uTime;
  uniform float uIntensity;
  uniform float uTidalPhase; // 潮汐相位 0-1
  varying vec2 vUv;
  varying float vElevation;
  varying float vDepth;

  // Simplex noise function
  vec3 mod289(vec3 x) { return x - floor(x * (1.0 / 289.0)) * 289.0; }
  vec4 mod289(vec4 x) { return x - floor(x * (1.0 / 289.0)) * 289.0; }
  vec4 permute(vec4 x) { return mod289(((x*34.0)+1.0)*x); }
  vec4 taylorInvSqrt(vec4 r) { return 1.79284291400159 - 0.85373472095314 * r; }

  float snoise(vec3 v) {
    const vec2 C = vec2(1.0/6.0, 1.0/3.0);
    const vec4 D = vec4(0.0, 0.5, 1.0, 2.0);

    vec3 i  = floor(v + dot(v, C.yyy));
    vec3 x0 = v - i + dot(i, C.xxx);

    vec3 g = step(x0.yzx, x0.xyz);
    vec3 l = 1.0 - g;
    vec3 i1 = min(g.xyz, l.zxy);
    vec3 i2 = max(g.xyz, l.zxy);

    vec3 x1 = x0 - i1 + C.xxx;
    vec3 x2 = x0 - i2 + C.yyy;
    vec3 x3 = x0 - D.yyy;

    i = mod289(i);
    vec4 p = permute(permute(permute(
              i.z + vec4(0.0, i1.z, i2.z, 1.0))
            + i.y + vec4(0.0, i1.y, i2.y, 1.0))
            + i.x + vec4(0.0, i1.x, i2.x, 1.0));

    float n_ = 0.142857142857;
    vec3 ns = n_ * D.wyz - D.xzx;

    vec4 j = p - 49.0 * floor(p * ns.z * ns.z);

    vec4 x_ = floor(j * ns.z);
    vec4 y_ = floor(j - 7.0 * x_);

    vec4 x = x_ *ns.x + ns.yyyy;
    vec4 y = y_ *ns.x + ns.yyyy;
    vec4 h = 1.0 - abs(x) - abs(y);

    vec4 b0 = vec4(x.xy, y.xy);
    vec4 b1 = vec4(x.zw, y.zw);

    vec4 s0 = floor(b0)*2.0 + 1.0;
    vec4 s1 = floor(b1)*2.0 + 1.0;
    vec4 sh = -step(h, vec4(0.0));

    vec4 a0 = b0.xzyw + s0.xzyw*sh.xxyy;
    vec4 a1 = b1.xzyw + s1.xzyw*sh.zzww;

    vec3 p0 = vec3(a0.xy, h.x);
    vec3 p1 = vec3(a0.zw, h.y);
    vec3 p2 = vec3(a1.xy, h.z);
    vec3 p3 = vec3(a1.zw, h.w);

    vec4 norm = taylorInvSqrt(vec4(dot(p0,p0), dot(p1,p1), dot(p2,p2), dot(p3,p3)));
    p0 *= norm.x;
    p1 *= norm.y;
    p2 *= norm.z;
    p3 *= norm.w;

    vec4 m = max(0.6 - vec4(dot(x0,x0), dot(x1,x1), dot(x2,x2), dot(x3,x3)), 0.0);
    m = m * m;
    return 42.0 * dot(m*m, vec4(dot(p0,x0), dot(p1,x1), dot(p2,x2), dot(p3,x3)));
  }

  void main() {
    vUv = uv;

    vec3 pos = position;

    // 计算到中心的距离（用于同心圆波纹）
    float distFromCenter = length(pos.xy);

    // 主波浪层 - 自然海洋波动
    float wave1 = snoise(vec3(pos.x * 0.4, pos.y * 0.4, uTime * 0.25)) * 0.6;
    float wave2 = snoise(vec3(pos.x * 0.8, pos.y * 0.8, uTime * 0.4)) * 0.3;
    float wave3 = snoise(vec3(pos.x * 1.6, pos.y * 1.6, uTime * 0.6)) * 0.15;

    // 潮汐呼吸效果 - 从中心向外扩散的波纹
    float tidalBreath = sin(uTime * 0.15) * 0.4;
    float ripple = sin(distFromCenter * 2.0 - uTime * 0.8) * 0.2 * smoothstep(4.0, 0.0, distFromCenter);

    // 中心涡旋效果
    float vortex = sin(distFromCenter * 1.5 - uTime * 0.5 + uTidalPhase * 6.28) * 0.15;
    vortex *= smoothstep(3.0, 0.5, distFromCenter);

    float elevation = (wave1 + wave2 + wave3 + tidalBreath + ripple + vortex) * uIntensity;
    pos.z = elevation;

    vElevation = elevation;
    vDepth = distFromCenter;

    gl_Position = projectionMatrix * modelViewMatrix * vec4(pos, 1.0);
  }
`

const waveFragmentShader = `
  uniform vec3 uPrimaryColor;
  uniform vec3 uSecondaryColor;
  uniform vec3 uAccentColor;
  uniform float uTime;
  varying vec2 vUv;
  varying float vElevation;
  varying float vDepth;

  void main() {
    // 基于高度的颜色混合 - 增强亮度
    float mixStrength = (vElevation + 0.5) * 0.7;
    vec3 color = mix(uSecondaryColor, uPrimaryColor, mixStrength);

    // 整体提亮
    color *= 1.2;

    // 中心高亮效果 - 增强
    float centerGlow = smoothstep(3.5, 0.0, vDepth) * 0.4;
    color = mix(color, uAccentColor, centerGlow);

    // 波光粼粼效果 - 增强
    float shimmer = sin(vUv.x * 30.0 + uTime * 1.5) * sin(vUv.y * 30.0 + uTime * 1.2) * 0.12;
    shimmer += sin(vDepth * 5.0 - uTime * 2.0) * 0.08;
    color += shimmer * uAccentColor;

    // 能量脉冲线 - 增强
    float pulse = sin(vDepth * 3.0 - uTime * 1.5) * 0.5 + 0.5;
    pulse = smoothstep(0.7, 1.0, pulse) * smoothstep(4.0, 1.0, vDepth);
    color += pulse * uAccentColor * 0.4;

    // 边缘渐变淡出 - 稍微扩大可见区域
    float edgeFade = smoothstep(0.0, 0.2, vUv.x) * smoothstep(1.0, 0.8, vUv.x) *
                     smoothstep(0.0, 0.2, vUv.y) * smoothstep(1.0, 0.8, vUv.y);

    // 增加透明度让波浪更明显
    float alpha = 0.35 * edgeFade;

    gl_FragColor = vec4(color, alpha);
  }
`

// ===== 潮汐核心着色器 =====
const coreVertexShader = `
  uniform float uTime;
  uniform float uIntensity;
  varying vec3 vNormal;
  varying vec3 vPosition;
  varying float vPulse;

  void main() {
    vNormal = normalize(normalMatrix * normal);
    vPosition = position;

    // 呼吸脉动
    float breathe = sin(uTime * 0.3) * 0.1 + 1.0;
    float pulse = sin(uTime * 2.0) * 0.02;
    vPulse = pulse;

    vec3 pos = position * (breathe + pulse * uIntensity);

    gl_Position = projectionMatrix * modelViewMatrix * vec4(pos, 1.0);
  }
`

const coreFragmentShader = `
  uniform vec3 uPrimaryColor;
  uniform vec3 uAccentColor;
  uniform float uTime;
  varying vec3 vNormal;
  varying vec3 vPosition;
  varying float vPulse;

  void main() {
    // 菲涅尔边缘发光 - 增强
    vec3 viewDir = normalize(cameraPosition - vPosition);
    float fresnel = pow(1.0 - abs(dot(viewDir, vNormal)), 2.5);

    // 核心渐变 - 提亮
    vec3 color = mix(uPrimaryColor * 0.5, uAccentColor, fresnel);

    // 能量脉冲 - 增强
    float energyPulse = sin(uTime * 3.0) * 0.5 + 0.5;
    color += uAccentColor * energyPulse * 0.35;

    // 内部发光 - 增强
    float innerGlow = smoothstep(1.0, 0.0, length(vPosition)) * 0.65;
    color += uAccentColor * innerGlow;

    float alpha = 0.7 + fresnel * 0.3;

    gl_FragColor = vec4(color, alpha);
  }
`

// ===== 粒子着色器 - 增强版 =====
const particleVertexShader = `
  uniform float uTime;
  uniform float uIntensity;
  attribute float aScale;
  attribute float aSpeed;
  attribute float aOffset;
  attribute float aOrbit; // 轨道半径
  varying float vAlpha;
  varying float vGlow;

  void main() {
    vec3 pos = position;

    // 围绕中心旋转的轨道运动
    float angle = uTime * aSpeed * 0.3 + aOffset;
    float orbitRadius = aOrbit * (1.0 + sin(uTime * 0.5 + aOffset) * 0.2);

    pos.x += cos(angle) * orbitRadius * 0.5;
    pos.z += sin(angle) * orbitRadius * 0.3;

    // 上升漂浮
    float floatY = sin(uTime * aSpeed + aOffset) * 0.3;
    pos.y += floatY * uIntensity;
    pos.y += mod(uTime * aSpeed * 0.15 + aOffset, 5.0) - 2.5;

    vec4 mvPosition = modelViewMatrix * vec4(pos, 1.0);

    // 大小衰减
    gl_PointSize = aScale * 60.0 * (1.0 / -mvPosition.z) * uIntensity;

    // 基于位置的透明度
    float distFromCenter = length(pos.xz);
    vAlpha = smoothstep(-2.5, 0.0, pos.y) * smoothstep(2.5, 0.0, pos.y);
    vAlpha *= smoothstep(4.0, 1.0, distFromCenter);

    // 发光强度
    vGlow = sin(uTime * 2.0 + aOffset * 3.0) * 0.5 + 0.5;

    gl_Position = projectionMatrix * mvPosition;
  }
`

const particleFragmentShader = `
  uniform vec3 uColor;
  uniform vec3 uAccentColor;
  varying float vAlpha;
  varying float vGlow;

  void main() {
    // 圆形粒子
    float dist = length(gl_PointCoord - vec2(0.5));
    if (dist > 0.5) discard;

    // 柔和边缘 - 增强透明度
    float alpha = smoothstep(0.5, 0.1, dist) * vAlpha * 0.9;

    // 颜色混合 - 更亮
    vec3 color = mix(uColor, uAccentColor, vGlow * 0.6);

    // 中心发光 - 增强
    float centerGlow = smoothstep(0.3, 0.0, dist) * vGlow;
    color += uAccentColor * centerGlow * 0.7;

    gl_FragColor = vec4(color, alpha);
  }
`

function initThree() {
  if (!containerRef.value) return

  const width = containerRef.value.clientWidth
  const height = containerRef.value.clientHeight

  // Scene
  scene = new THREE.Scene()

  // Camera - 调整视角以更好展示潮汐核心和大海浪
  camera = new THREE.PerspectiveCamera(70, width / height, 0.1, 100)
  camera.position.z = 5
  camera.position.y = 2.5

  // Renderer
  renderer = new THREE.WebGLRenderer({
    antialias: true,
    alpha: true,
    powerPreference: 'high-performance'
  })
  renderer.setSize(width, height)
  renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2))
  renderer.setClearColor(0x000000, 0)
  containerRef.value.appendChild(renderer.domElement)

  // Clock
  clock = new THREE.Clock()

  // Create scene elements
  createWavePlane()
  createTidalCore()
  createEnergyRings()
  createParticles()

  // Start animation
  animate()

  // Handle resize
  window.addEventListener('resize', handleResize)
}

function createWavePlane() {
  // 更大的波浪平面 - 增大尺寸覆盖更多区域
  waveGeometry = new THREE.PlaneGeometry(20, 20, 200, 200)

  const waveMaterial = new THREE.ShaderMaterial({
    vertexShader: waveVertexShader,
    fragmentShader: waveFragmentShader,
    transparent: true,
    side: THREE.DoubleSide,
    uniforms: {
      uTime: { value: 0 },
      uTidalPhase: { value: 0 },
      uIntensity: { value: props.intensity ?? 1 },
      uPrimaryColor: { value: currentPrimaryColor },
      uSecondaryColor: { value: currentSecondaryColor },
      uAccentColor: { value: currentAccentColor }
    }
  })

  waveMesh = new THREE.Mesh(waveGeometry, waveMaterial)
  waveMesh.rotation.x = -Math.PI / 2.2
  waveMesh.position.y = -1.8
  scene.add(waveMesh)
}

// 创建潮汐核心 - 中央发光球体
function createTidalCore() {
  const coreGeometry = new THREE.SphereGeometry(0.3, 64, 64)

  const coreMaterial = new THREE.ShaderMaterial({
    vertexShader: coreVertexShader,
    fragmentShader: coreFragmentShader,
    transparent: true,
    side: THREE.DoubleSide,
    uniforms: {
      uTime: { value: 0 },
      uIntensity: { value: props.intensity ?? 1 },
      uPrimaryColor: { value: currentPrimaryColor },
      uAccentColor: { value: currentAccentColor }
    }
  })

  coreMesh = new THREE.Mesh(coreGeometry, coreMaterial)
  coreMesh.position.set(0, 0.3, 0)
  scene.add(coreMesh)
}

// 创建能量环 - 围绕核心的同心圆环 - 增大尺寸
function createEnergyRings() {
  const ringCount = 4
  const ringColors = [0x38bdf8, 0x22d3ee, 0x0ea5e9, 0x0284c7]

  for (let i = 0; i < ringCount; i++) {
    const radius = 0.6 + i * 0.35
    const ringGeometry = new THREE.TorusGeometry(radius, 0.015, 16, 100)

    const ringMaterial = new THREE.MeshBasicMaterial({
      color: ringColors[i],
      transparent: true,
      opacity: 0.5 - i * 0.08
    })

    const ring = new THREE.Mesh(ringGeometry, ringMaterial)
    ring.position.set(0, 0.3, 0)
    ring.rotation.x = Math.PI / 2

    ringMeshes.push(ring)
    scene.add(ring)
  }
}

function createParticles() {
  const particleCount = 450 // 增加粒子数量
  const positions = new Float32Array(particleCount * 3)
  const scales = new Float32Array(particleCount)
  const speeds = new Float32Array(particleCount)
  const offsets = new Float32Array(particleCount)
  const orbits = new Float32Array(particleCount)

  for (let i = 0; i < particleCount; i++) {
    // 以核心为中心分布 - 扩大范围
    const theta = Math.random() * Math.PI * 2
    const phi = Math.random() * Math.PI
    const radius = 0.8 + Math.random() * 5

    positions[i * 3] = Math.sin(phi) * Math.cos(theta) * radius
    positions[i * 3 + 1] = (Math.random() - 0.5) * 6
    positions[i * 3 + 2] = Math.sin(phi) * Math.sin(theta) * radius

    scales[i] = Math.random() * 0.8 + 0.5
    speeds[i] = Math.random() * 0.4 + 0.3
    offsets[i] = Math.random() * Math.PI * 2
    orbits[i] = Math.random() * 3 + 0.8
  }

  const particleGeometry = new THREE.BufferGeometry()
  particleGeometry.setAttribute('position', new THREE.BufferAttribute(positions, 3))
  particleGeometry.setAttribute('aScale', new THREE.BufferAttribute(scales, 1))
  particleGeometry.setAttribute('aSpeed', new THREE.BufferAttribute(speeds, 1))
  particleGeometry.setAttribute('aOffset', new THREE.BufferAttribute(offsets, 1))
  particleGeometry.setAttribute('aOrbit', new THREE.BufferAttribute(orbits, 1))

  const particleMaterial = new THREE.ShaderMaterial({
    vertexShader: particleVertexShader,
    fragmentShader: particleFragmentShader,
    transparent: true,
    depthWrite: false,
    blending: THREE.AdditiveBlending,
    uniforms: {
      uTime: { value: 0 },
      uIntensity: { value: props.intensity ?? 1 },
      uColor: { value: currentPrimaryColor },
      uAccentColor: { value: currentAccentColor }
    }
  })

  particles = new THREE.Points(particleGeometry, particleMaterial)
  scene.add(particles)
}

function animate() {
  animationId = requestAnimationFrame(animate)

  const elapsedTime = clock.getElapsedTime()

  // 平滑颜色过渡
  currentPrimaryColor.lerp(targetPrimaryColor, 0.025)
  currentSecondaryColor.lerp(targetSecondaryColor, 0.025)
  currentAccentColor.lerp(targetAccentColor, 0.025)

  // 计算潮汐相位 (用于同步所有元素)
  const tidalPhase = (Math.sin(elapsedTime * 0.15) + 1) / 2

  // 更新波浪 uniforms
  if (waveMesh) {
    const waveMaterial = waveMesh.material as THREE.ShaderMaterial
    if (waveMaterial.uniforms.uTime) waveMaterial.uniforms.uTime.value = elapsedTime
    if (waveMaterial.uniforms.uTidalPhase) waveMaterial.uniforms.uTidalPhase.value = tidalPhase
    if (waveMaterial.uniforms.uIntensity) waveMaterial.uniforms.uIntensity.value = props.isActive ? (props.intensity ?? 1) : 0.6
    if (waveMaterial.uniforms.uPrimaryColor) waveMaterial.uniforms.uPrimaryColor.value = currentPrimaryColor
    if (waveMaterial.uniforms.uSecondaryColor) waveMaterial.uniforms.uSecondaryColor.value = currentSecondaryColor
    if (waveMaterial.uniforms.uAccentColor) waveMaterial.uniforms.uAccentColor.value = currentAccentColor
  }

  // 更新核心 uniforms
  if (coreMesh) {
    const coreMaterial = coreMesh.material as THREE.ShaderMaterial
    if (coreMaterial.uniforms.uTime) coreMaterial.uniforms.uTime.value = elapsedTime
    if (coreMaterial.uniforms.uIntensity) coreMaterial.uniforms.uIntensity.value = props.isActive ? (props.intensity ?? 1) : 0.6
    if (coreMaterial.uniforms.uPrimaryColor) coreMaterial.uniforms.uPrimaryColor.value = currentPrimaryColor
    if (coreMaterial.uniforms.uAccentColor) coreMaterial.uniforms.uAccentColor.value = currentAccentColor

    // 核心微妙旋转
    coreMesh.rotation.y = elapsedTime * 0.1
  }

  // 更新能量环
  ringMeshes.forEach((ring, index) => {
    // 不同速度的旋转
    ring.rotation.z = elapsedTime * (0.2 + index * 0.1)

    // 脉动效果
    const pulse = Math.sin(elapsedTime * 2 + index * Math.PI / 3) * 0.02 + 1
    ring.scale.set(pulse, pulse, 1)

    // 更新颜色
    const material = ring.material as THREE.MeshBasicMaterial
    material.color.lerp(currentAccentColor, 0.02)
    material.opacity = (0.4 - index * 0.1) * (props.isActive ? 1 : 0.5)
  })

  // 更新粒子 uniforms
  if (particles) {
    const particleMaterial = particles.material as THREE.ShaderMaterial
    if (particleMaterial.uniforms.uTime) particleMaterial.uniforms.uTime.value = elapsedTime
    if (particleMaterial.uniforms.uIntensity) particleMaterial.uniforms.uIntensity.value = props.isActive ? (props.intensity ?? 1) : 0.5
    if (particleMaterial.uniforms.uColor) particleMaterial.uniforms.uColor.value = currentPrimaryColor
    if (particleMaterial.uniforms.uAccentColor) particleMaterial.uniforms.uAccentColor.value = currentAccentColor

    // 粒子整体缓慢旋转
    particles.rotation.y = elapsedTime * 0.03
  }

  // 相机呼吸运动 - 更加流畅
  camera.position.x = Math.sin(elapsedTime * 0.08) * 0.3
  camera.position.y = 1.5 + Math.cos(elapsedTime * 0.1) * 0.15
  camera.lookAt(0, 0, 0)

  renderer.render(scene, camera)
}

function handleResize() {
  if (!containerRef.value) return

  const width = containerRef.value.clientWidth
  const height = containerRef.value.clientHeight

  camera.aspect = width / height
  camera.updateProjectionMatrix()

  renderer.setSize(width, height)
}

// 监听相位变化
watch(() => props.phase, (newPhase) => {
  if (newPhase) {
    const colors = phaseColors[newPhase]
    targetPrimaryColor = new THREE.Color(colors.primary)
    targetSecondaryColor = new THREE.Color(colors.secondary)
    targetAccentColor = new THREE.Color(colors.accent)
  }
}, { immediate: true })

onMounted(() => {
  initThree()
})

onUnmounted(() => {
  if (animationId) {
    cancelAnimationFrame(animationId)
  }

  window.removeEventListener('resize', handleResize)

  // 清理资源
  if (renderer) {
    renderer.dispose()
    if (containerRef.value && renderer.domElement) {
      containerRef.value.removeChild(renderer.domElement)
    }
  }

  if (waveGeometry) waveGeometry.dispose()
  if (waveMesh) {
    (waveMesh.material as THREE.Material).dispose()
  }

  if (coreMesh) {
    coreMesh.geometry.dispose()
    ;(coreMesh.material as THREE.Material).dispose()
  }

  ringMeshes.forEach(ring => {
    ring.geometry.dispose()
    ;(ring.material as THREE.Material).dispose()
  })
  ringMeshes = []

  if (particles) {
    particles.geometry.dispose()
    ;(particles.material as THREE.Material).dispose()
  }
})
</script>

<template>
  <div ref="containerRef" class="tidal-background" />
</template>

<style scoped>
.tidal-background {
  position: absolute;
  inset: 0;
  z-index: 0;
  pointer-events: none;
  overflow: hidden;
}
</style>
