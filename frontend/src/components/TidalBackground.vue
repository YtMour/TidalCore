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
let animationId: number
let clock: THREE.Clock

// Phase colors
const phaseColors = {
  idle: { primary: 0x64748b, secondary: 0x475569 },
  contract: { primary: 0xf43f5e, secondary: 0xec4899 },
  hold: { primary: 0xf59e0b, secondary: 0xf97316 },
  relax: { primary: 0x10b981, secondary: 0x14b8a6 }
}

// Current target colors for smooth transitions
let currentPrimaryColor = new THREE.Color(phaseColors.idle.primary)
let targetPrimaryColor = new THREE.Color(phaseColors.idle.primary)
let currentSecondaryColor = new THREE.Color(phaseColors.idle.secondary)
let targetSecondaryColor = new THREE.Color(phaseColors.idle.secondary)

// Shader for the wave plane
const waveVertexShader = `
  uniform float uTime;
  uniform float uIntensity;
  varying vec2 vUv;
  varying float vElevation;

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

    // Multiple wave layers for realistic ocean effect
    float wave1 = snoise(vec3(pos.x * 0.5, pos.y * 0.5, uTime * 0.3)) * 0.5;
    float wave2 = snoise(vec3(pos.x * 1.0, pos.y * 1.0, uTime * 0.5)) * 0.25;
    float wave3 = snoise(vec3(pos.x * 2.0, pos.y * 2.0, uTime * 0.7)) * 0.125;

    // Tidal breathing effect
    float tidal = sin(uTime * 0.2) * 0.3;

    float elevation = (wave1 + wave2 + wave3 + tidal) * uIntensity;
    pos.z = elevation;

    vElevation = elevation;

    gl_Position = projectionMatrix * modelViewMatrix * vec4(pos, 1.0);
  }
`

const waveFragmentShader = `
  uniform vec3 uPrimaryColor;
  uniform vec3 uSecondaryColor;
  uniform float uTime;
  varying vec2 vUv;
  varying float vElevation;

  void main() {
    // Color based on elevation
    float mixStrength = (vElevation + 0.5) * 0.5;
    vec3 color = mix(uSecondaryColor, uPrimaryColor, mixStrength);

    // Add shimmer effect
    float shimmer = sin(vUv.x * 20.0 + uTime) * sin(vUv.y * 20.0 + uTime * 0.7) * 0.1;
    color += shimmer;

    // Fade edges
    float edgeFade = smoothstep(0.0, 0.3, vUv.x) * smoothstep(1.0, 0.7, vUv.x) *
                     smoothstep(0.0, 0.3, vUv.y) * smoothstep(1.0, 0.7, vUv.y);

    float alpha = 0.15 * edgeFade;

    gl_FragColor = vec4(color, alpha);
  }
`

// Particle shader
const particleVertexShader = `
  uniform float uTime;
  uniform float uIntensity;
  attribute float aScale;
  attribute float aSpeed;
  attribute float aOffset;
  varying float vAlpha;

  void main() {
    vec3 pos = position;

    // Floating motion
    float floatY = sin(uTime * aSpeed + aOffset) * 0.5;
    float floatX = cos(uTime * aSpeed * 0.7 + aOffset) * 0.3;

    pos.y += floatY * uIntensity;
    pos.x += floatX * uIntensity;

    // Rising motion (like bubbles)
    pos.y += mod(uTime * aSpeed * 0.2 + aOffset, 4.0) - 2.0;

    vec4 mvPosition = modelViewMatrix * vec4(pos, 1.0);

    // Size attenuation
    gl_PointSize = aScale * 50.0 * (1.0 / -mvPosition.z) * uIntensity;

    // Alpha based on position
    vAlpha = smoothstep(-2.0, 0.0, pos.y) * smoothstep(2.0, 0.0, pos.y);

    gl_Position = projectionMatrix * mvPosition;
  }
`

const particleFragmentShader = `
  uniform vec3 uColor;
  varying float vAlpha;

  void main() {
    // Circular particle
    float dist = length(gl_PointCoord - vec2(0.5));
    if (dist > 0.5) discard;

    // Soft edges
    float alpha = smoothstep(0.5, 0.2, dist) * vAlpha * 0.6;

    gl_FragColor = vec4(uColor, alpha);
  }
`

function initThree() {
  if (!containerRef.value) return

  const width = containerRef.value.clientWidth
  const height = containerRef.value.clientHeight

  // Scene
  scene = new THREE.Scene()

  // Camera
  camera = new THREE.PerspectiveCamera(75, width / height, 0.1, 100)
  camera.position.z = 3
  camera.position.y = 1

  // Renderer
  renderer = new THREE.WebGLRenderer({
    antialias: true,
    alpha: true
  })
  renderer.setSize(width, height)
  renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2))
  renderer.setClearColor(0x000000, 0)
  containerRef.value.appendChild(renderer.domElement)

  // Clock
  clock = new THREE.Clock()

  // Create wave plane
  createWavePlane()

  // Create particles
  createParticles()

  // Start animation
  animate()

  // Handle resize
  window.addEventListener('resize', handleResize)
}

function createWavePlane() {
  waveGeometry = new THREE.PlaneGeometry(8, 8, 128, 128)

  const waveMaterial = new THREE.ShaderMaterial({
    vertexShader: waveVertexShader,
    fragmentShader: waveFragmentShader,
    transparent: true,
    side: THREE.DoubleSide,
    uniforms: {
      uTime: { value: 0 },
      uIntensity: { value: props.intensity ?? 1 },
      uPrimaryColor: { value: currentPrimaryColor },
      uSecondaryColor: { value: currentSecondaryColor }
    }
  })

  waveMesh = new THREE.Mesh(waveGeometry, waveMaterial)
  waveMesh.rotation.x = -Math.PI / 2.5
  waveMesh.position.y = -1
  scene.add(waveMesh)
}

function createParticles() {
  const particleCount = 200
  const positions = new Float32Array(particleCount * 3)
  const scales = new Float32Array(particleCount)
  const speeds = new Float32Array(particleCount)
  const offsets = new Float32Array(particleCount)

  for (let i = 0; i < particleCount; i++) {
    positions[i * 3] = (Math.random() - 0.5) * 6
    positions[i * 3 + 1] = (Math.random() - 0.5) * 4
    positions[i * 3 + 2] = (Math.random() - 0.5) * 4

    scales[i] = Math.random() * 0.5 + 0.5
    speeds[i] = Math.random() * 0.5 + 0.5
    offsets[i] = Math.random() * Math.PI * 2
  }

  const particleGeometry = new THREE.BufferGeometry()
  particleGeometry.setAttribute('position', new THREE.BufferAttribute(positions, 3))
  particleGeometry.setAttribute('aScale', new THREE.BufferAttribute(scales, 1))
  particleGeometry.setAttribute('aSpeed', new THREE.BufferAttribute(speeds, 1))
  particleGeometry.setAttribute('aOffset', new THREE.BufferAttribute(offsets, 1))

  const particleMaterial = new THREE.ShaderMaterial({
    vertexShader: particleVertexShader,
    fragmentShader: particleFragmentShader,
    transparent: true,
    depthWrite: false,
    blending: THREE.AdditiveBlending,
    uniforms: {
      uTime: { value: 0 },
      uIntensity: { value: props.intensity ?? 1 },
      uColor: { value: currentPrimaryColor }
    }
  })

  particles = new THREE.Points(particleGeometry, particleMaterial)
  scene.add(particles)
}

function animate() {
  animationId = requestAnimationFrame(animate)

  const elapsedTime = clock.getElapsedTime()

  // Smooth color transitions
  currentPrimaryColor.lerp(targetPrimaryColor, 0.02)
  currentSecondaryColor.lerp(targetSecondaryColor, 0.02)

  // Update wave uniforms
  if (waveMesh) {
    const waveMaterial = waveMesh.material as THREE.ShaderMaterial
    if (waveMaterial.uniforms.uTime) waveMaterial.uniforms.uTime.value = elapsedTime
    if (waveMaterial.uniforms.uIntensity) waveMaterial.uniforms.uIntensity.value = props.isActive ? (props.intensity ?? 1) : 0.5
    if (waveMaterial.uniforms.uPrimaryColor) waveMaterial.uniforms.uPrimaryColor.value = currentPrimaryColor
    if (waveMaterial.uniforms.uSecondaryColor) waveMaterial.uniforms.uSecondaryColor.value = currentSecondaryColor
  }

  // Update particle uniforms
  if (particles) {
    const particleMaterial = particles.material as THREE.ShaderMaterial
    if (particleMaterial.uniforms.uTime) particleMaterial.uniforms.uTime.value = elapsedTime
    if (particleMaterial.uniforms.uIntensity) particleMaterial.uniforms.uIntensity.value = props.isActive ? (props.intensity ?? 1) : 0.5
    if (particleMaterial.uniforms.uColor) particleMaterial.uniforms.uColor.value = currentPrimaryColor

    // Subtle rotation
    particles.rotation.y = elapsedTime * 0.05
  }

  // Camera subtle movement
  camera.position.x = Math.sin(elapsedTime * 0.1) * 0.2
  camera.position.y = 1 + Math.cos(elapsedTime * 0.15) * 0.1

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

// Watch for phase changes
watch(() => props.phase, (newPhase) => {
  if (newPhase) {
    const colors = phaseColors[newPhase]
    targetPrimaryColor = new THREE.Color(colors.primary)
    targetSecondaryColor = new THREE.Color(colors.secondary)
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

  // Cleanup
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
