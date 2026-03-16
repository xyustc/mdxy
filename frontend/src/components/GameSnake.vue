<template>
  <div class="game-snake">
    <div class="game-header">
      <div class="scores">
        <div class="score-box glass-card">
          <span class="score-label">分数</span>
          <span class="score-value">
            {{ score }}
            <span v-if="scorePopup" class="score-popup" :key="scorePopupKey">+{{ scorePopup }}</span>
          </span>
        </div>
        <div class="score-box glass-card">
          <span class="score-label">最高分</span>
          <span class="score-value">{{ bestScore }}</span>
        </div>
      </div>
      <div class="game-actions">
        <button
          v-if="gameStatus === 'running' || gameStatus === 'paused'"
          class="btn-game glass-card"
          @click="togglePause"
        >
          {{ paused ? '继续' : '暂停' }}
        </button>
        <button class="btn-game glass-card" @click="resetGame">新游戏</button>
        <button class="btn-game glass-card btn-close" aria-label="关闭贪吃蛇游戏" @click="$emit('close')">✕</button>
      </div>
    </div>

    <div class="board glass-card" ref="boardEl">
      <canvas ref="canvasEl" :width="CANVAS_SIZE" :height="CANVAS_SIZE"></canvas>
    </div>

    <div v-if="gameStatus === 'ready'" class="game-overlay glass-card game-overlay--ready">
      <p class="overlay-text">准备就绪</p>
      <p class="overlay-score">点击开始或方向键开局</p>
      <button class="btn-game glass-card" @click="startGame">开始游戏</button>
    </div>

    <div v-if="gameOver" class="game-overlay glass-card">
      <p class="overlay-text">游戏结束</p>
      <p class="overlay-score">得分: {{ score }}</p>
      <button class="btn-game glass-card" @click="resetGame">再来一局</button>
    </div>

    <div class="mobile-controls">
      <div class="control-row">
        <button class="control-btn glass-card" aria-label="向上移动" @click="changeDirection('up')">↑</button>
      </div>
      <div class="control-row">
        <button class="control-btn glass-card" aria-label="向左移动" @click="changeDirection('left')">←</button>
        <button class="control-btn glass-card" aria-label="向下移动" @click="changeDirection('down')">↓</button>
        <button class="control-btn glass-card" aria-label="向右移动" @click="changeDirection('right')">→</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted, onUnmounted } from 'vue'
import { gameApi } from '@/api/game'
import { getOrCreatePlayerId } from '@/utils/game'

defineEmits<{ close: [] }>()

const TOOL_ID = 2
const playerId = getOrCreatePlayerId()
const CANVAS_SIZE = 400
const GRID_SIZE = 20
const CELL_SIZE = CANVAS_SIZE / GRID_SIZE
const INITIAL_SPEED = 150
const SWIPE_TRIGGER_DISTANCE = 30

interface Point { x: number; y: number }
type Direction = 'up' | 'down' | 'left' | 'right'
type GameStatus = 'ready' | 'running' | 'paused' | 'dying' | 'gameover'

const canvasEl = ref<HTMLCanvasElement>()
const boardEl = ref<HTMLElement>()
const score = ref(0)
const bestScore = ref(0)
const gameStatus = ref<GameStatus>('ready')
const scorePopup = ref(0)
const scorePopupKey = ref(0)
const gameOver = computed(() => gameStatus.value === 'gameover')
const paused = computed(() => gameStatus.value === 'paused')

let snake: Point[] = []
let direction: Point = { x: 1, y: 0 }
let directionQueue: Point[] = []
let food: Point = { x: 0, y: 0 }
let speed = INITIAL_SPEED
let rafId: number | null = null
let lastTime = 0
let accumulated = 0
let dying = false
let deathFlashCount = 0
let deathFlashTimer: ReturnType<typeof setInterval> | null = null
let scorePopupTimer: ReturnType<typeof setTimeout> | null = null
let scoreSyncTimer: ReturnType<typeof setTimeout> | null = null
let pendingSyncScore: number | null = null
let syncInFlight = false
let staticBoardLayer: HTMLCanvasElement | null = null
let cachedBgColor = ''
let cachedGridColor = 'rgba(128, 128, 128, 0.15)'
let themeObserver: MutationObserver | null = null

const SCORE_SYNC_DELAY = 800

function initGame() {
  snake = [{ x: 10, y: 10 }, { x: 9, y: 10 }, { x: 8, y: 10 }]
  direction = { x: 1, y: 0 }
  directionQueue = []
  score.value = 0
  gameStatus.value = 'ready'
  scorePopup.value = 0
  speed = INITIAL_SPEED
  dying = false
  deathFlashCount = 0
  if (deathFlashTimer) {
    clearInterval(deathFlashTimer)
    deathFlashTimer = null
  }
  spawnFood()
}

function spawnFood() {
  do {
    food = { x: Math.floor(Math.random() * GRID_SIZE), y: Math.floor(Math.random() * GRID_SIZE) }
  } while (snake.some(s => s.x === food.x && s.y === food.y))
}

function changeDirection(dir: Direction) {
  if (gameStatus.value === 'ready') {
    startGame()
  }
  if (gameStatus.value !== 'running') return
  const vectors: Record<Direction, Point> = { up: { x: 0, y: -1 }, down: { x: 0, y: 1 }, left: { x: -1, y: 0 }, right: { x: 1, y: 0 } }
  const newDir = vectors[dir]
  const lastDir = directionQueue.length > 0 ? directionQueue[directionQueue.length - 1] : direction
  if (lastDir.x + newDir.x !== 0 || lastDir.y + newDir.y !== 0) {
    if (directionQueue.length < 3) {
      directionQueue.push(newDir)
    }
  }
}

function togglePause() {
  if (gameStatus.value === 'ready') {
    startGame()
    return
  }
  if (gameStatus.value === 'gameover' || gameStatus.value === 'dying') return
  if (gameStatus.value === 'running') {
    gameStatus.value = 'paused'
    stopLoop()
  } else if (gameStatus.value === 'paused') {
    gameStatus.value = 'running'
    startLoop()
  }
  draw()
}

function startGame() {
  if (gameStatus.value === 'running' || gameStatus.value === 'dying') return
  gameStatus.value = 'running'
  startLoop()
}

function update() {
  if (directionQueue.length > 0) {
    direction = directionQueue.shift()!
  }
  const head = { x: snake[0].x + direction.x, y: snake[0].y + direction.y }

  if (head.x < 0 || head.x >= GRID_SIZE || head.y < 0 || head.y >= GRID_SIZE ||
      snake.some(s => s.x === head.x && s.y === head.y)) {
    startDeathAnimation()
    return
  }

  snake.unshift(head)

  if (head.x === food.x && head.y === food.y) {
    score.value += 10
    scorePopup.value = 10
    scorePopupKey.value++
    if (scorePopupTimer) clearTimeout(scorePopupTimer)
    scorePopupTimer = setTimeout(() => { scorePopup.value = 0 }, 500)
    if (score.value > bestScore.value) {
      bestScore.value = score.value
      scheduleScoreSync()
    }
    spawnFood()
    speed = Math.max(50, speed - 2)
  } else {
    snake.pop()
  }
}

function startDeathAnimation() {
  dying = true
  gameStatus.value = 'dying'
  stopLoop()
  deathFlashCount = 0
  if (deathFlashTimer) clearInterval(deathFlashTimer)
  deathFlashTimer = setInterval(() => {
    deathFlashCount++
    draw()
    if (deathFlashCount >= 6) {
      if (deathFlashTimer) {
        clearInterval(deathFlashTimer)
        deathFlashTimer = null
      }
      dying = false
      gameStatus.value = 'gameover'
      scheduleScoreSync(true)
      draw()
    }
  }, 100)
}

function refreshBoardStyleCache() {
  const styles = getComputedStyle(document.documentElement)
  cachedBgColor = styles.getPropertyValue('--color-bg-tertiary').trim() || '#1f2226'
  cachedGridColor = 'rgba(128, 128, 128, 0.15)'
  buildStaticBoardLayer()
}

function buildStaticBoardLayer() {
  const layer = document.createElement('canvas')
  layer.width = CANVAS_SIZE
  layer.height = CANVAS_SIZE
  const ctx = layer.getContext('2d')
  if (!ctx) return
  ctx.fillStyle = cachedBgColor
  ctx.fillRect(0, 0, CANVAS_SIZE, CANVAS_SIZE)

  ctx.strokeStyle = cachedGridColor
  ctx.lineWidth = 0.5
  for (let i = 1; i < GRID_SIZE; i++) {
    const pos = i * CELL_SIZE
    ctx.beginPath(); ctx.moveTo(pos, 0); ctx.lineTo(pos, CANVAS_SIZE); ctx.stroke()
    ctx.beginPath(); ctx.moveTo(0, pos); ctx.lineTo(CANVAS_SIZE, pos); ctx.stroke()
  }
  staticBoardLayer = layer
}

function draw() {
  const canvas = canvasEl.value
  if (!canvas) return
  const ctx = canvas.getContext('2d')
  if (!ctx) return

  if (staticBoardLayer) {
    ctx.drawImage(staticBoardLayer, 0, 0)
  } else {
    ctx.fillStyle = cachedBgColor || '#1f2226'
    ctx.fillRect(0, 0, CANVAS_SIZE, CANVAS_SIZE)
  }

  const pulse = Math.sin(performance.now() / 300) * 2
  ctx.fillStyle = '#ff6b6b'
  ctx.beginPath()
  ctx.roundRect(
    food.x * CELL_SIZE + 2 - pulse / 2, food.y * CELL_SIZE + 2 - pulse / 2,
    CELL_SIZE - 4 + pulse, CELL_SIZE - 4 + pulse, 4
  )
  ctx.fill()

  // 蛇身渐变 + 死亡闪烁
  const isFlashRed = dying && deathFlashCount % 2 === 1
  const len = snake.length
  snake.forEach((s, i) => {
    const t = len > 1 ? i / (len - 1) : 0
    if (isFlashRed) {
      ctx.fillStyle = '#ff4444'
    } else {
      const r = Math.round(0 + t * 0)
      const g = Math.round(212 - t * 80)
      const b = Math.round(170 - t * 60)
      ctx.fillStyle = `rgb(${r},${g},${b})`
    }
    ctx.globalAlpha = 1 - t * 0.3
    const radius = i === 0 ? 6 : 3
    ctx.beginPath()
    ctx.roundRect(s.x * CELL_SIZE + 1, s.y * CELL_SIZE + 1, CELL_SIZE - 2, CELL_SIZE - 2, radius)
    ctx.fill()
  })
  ctx.globalAlpha = 1

  if (gameStatus.value === 'paused') {
    ctx.fillStyle = 'rgba(0, 0, 0, 0.4)'
    ctx.fillRect(0, 0, CANVAS_SIZE, CANVAS_SIZE)
    ctx.fillStyle = '#fff'
    ctx.font = 'bold 24px sans-serif'
    ctx.textAlign = 'center'
    ctx.textBaseline = 'middle'
    ctx.fillText('已暂停', CANVAS_SIZE / 2, CANVAS_SIZE / 2)
    ctx.font = '14px sans-serif'
    ctx.fillText('按空格键继续', CANVAS_SIZE / 2, CANVAS_SIZE / 2 + 30)
  }
}

function scheduleScoreSync(force = false) {
  pendingSyncScore = Math.max(pendingSyncScore ?? 0, score.value)
  if (force) {
    if (scoreSyncTimer) {
      clearTimeout(scoreSyncTimer)
      scoreSyncTimer = null
    }
    void flushScoreSync()
    return
  }
  if (scoreSyncTimer) return
  scoreSyncTimer = setTimeout(() => {
    scoreSyncTimer = null
    void flushScoreSync()
  }, SCORE_SYNC_DELAY)
}

async function flushScoreSync() {
  if (syncInFlight || pendingSyncScore === null) return
  syncInFlight = true
  const syncScore = pendingSyncScore
  pendingSyncScore = null
  try {
    const res = await gameApi.updateScore({ tool_id: TOOL_ID, player_id: playerId, score: syncScore })
    if (res.success && res.data) bestScore.value = res.data.best_score
  } catch (err) {
    console.error('Failed to sync score:', err)
  } finally {
    syncInFlight = false
    if (pendingSyncScore !== null) {
      void flushScoreSync()
    }
  }
}

async function loadBestScore() {
  try {
    const res = await gameApi.getBestScore(TOOL_ID, playerId)
    if (res.success && res.data) bestScore.value = res.data.best_score
  } catch (err) { console.error('Failed to load best score:', err) }
}

function gameFrame(time: number) {
  if (gameStatus.value !== 'running') {
    rafId = null
    return
  }
  rafId = requestAnimationFrame(gameFrame)
  if (!lastTime) { lastTime = time; return }
  accumulated += time - lastTime
  lastTime = time
  while (accumulated >= speed && gameStatus.value === 'running') {
    accumulated -= speed
    update()
  }
  draw()
}

function startLoop() {
  if (rafId !== null) return
  lastTime = 0
  accumulated = 0
  rafId = requestAnimationFrame(gameFrame)
}

function stopLoop() {
  if (rafId !== null) {
    cancelAnimationFrame(rafId)
    rafId = null
  }
}

function resetGame() {
  stopLoop()
  initGame()
  draw()
}

function onKeyDown(e: KeyboardEvent) {
  if (e.key === ' ') { e.preventDefault(); togglePause(); return }
  const map: Record<string, 'up' | 'down' | 'left' | 'right'> = {
    ArrowUp: 'up', ArrowDown: 'down', ArrowLeft: 'left', ArrowRight: 'right',
    w: 'up', s: 'down', a: 'left', d: 'right'
  }
  const dir = map[e.key]
  if (dir) { e.preventDefault(); changeDirection(dir) }
}

let touchStartX = 0
let touchStartY = 0

function onTouchStart(e: TouchEvent) {
  touchStartX = e.touches[0].clientX
  touchStartY = e.touches[0].clientY
}

function onTouchMove(e: TouchEvent) {
  if (!e.touches.length) return
  const dx = e.touches[0].clientX - touchStartX
  const dy = e.touches[0].clientY - touchStartY
  if (Math.max(Math.abs(dx), Math.abs(dy)) > 8 && e.cancelable) {
    e.preventDefault()
  }
}

function onTouchEnd(e: TouchEvent) {
  const dx = e.changedTouches[0].clientX - touchStartX
  const dy = e.changedTouches[0].clientY - touchStartY
  if (Math.max(Math.abs(dx), Math.abs(dy)) < SWIPE_TRIGGER_DISTANCE) return
  if (Math.abs(dx) > Math.abs(dy)) changeDirection(dx > 0 ? 'right' : 'left')
  else changeDirection(dy > 0 ? 'down' : 'up')
}

function onTouchCancel() {
  touchStartX = 0
  touchStartY = 0
}

function onVisibilityChange() {
  if (document.visibilityState === 'hidden') {
    if (gameStatus.value === 'running') {
      gameStatus.value = 'paused'
      stopLoop()
      draw()
    }
    if (pendingSyncScore !== null) {
      scheduleScoreSync(true)
    }
  }
}

onMounted(() => {
  loadBestScore()
  initGame()
  refreshBoardStyleCache()
  draw()
  window.addEventListener('keydown', onKeyDown)
  document.addEventListener('visibilitychange', onVisibilityChange)
  themeObserver = new MutationObserver(() => {
    refreshBoardStyleCache()
    draw()
  })
  themeObserver.observe(document.documentElement, { attributes: true, attributeFilter: ['data-theme'] })
  boardEl.value?.addEventListener('touchstart', onTouchStart, { passive: true })
  boardEl.value?.addEventListener('touchmove', onTouchMove, { passive: false })
  boardEl.value?.addEventListener('touchend', onTouchEnd, { passive: true })
  boardEl.value?.addEventListener('touchcancel', onTouchCancel, { passive: true })
})

onUnmounted(() => {
  stopLoop()
  if (deathFlashTimer) clearInterval(deathFlashTimer)
  if (scorePopupTimer) clearTimeout(scorePopupTimer)
  if (scoreSyncTimer) clearTimeout(scoreSyncTimer)
  if (pendingSyncScore !== null) {
    void flushScoreSync()
  }
  themeObserver?.disconnect()
  window.removeEventListener('keydown', onKeyDown)
  document.removeEventListener('visibilitychange', onVisibilityChange)
  boardEl.value?.removeEventListener('touchstart', onTouchStart)
  boardEl.value?.removeEventListener('touchmove', onTouchMove)
  boardEl.value?.removeEventListener('touchend', onTouchEnd)
  boardEl.value?.removeEventListener('touchcancel', onTouchCancel)
})
</script>

<style scoped>
.game-snake {
  max-width: 420px;
  margin: 0 auto;
  position: relative;
  user-select: none;
}

.game-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 16px;
  gap: 12px;
}

.scores { display: flex; gap: 10px; }

.score-box {
  padding: 8px 16px;
  text-align: center;
  min-width: 80px;
}

.score-label {
  display: block;
  font-size: 11px;
  color: var(--color-text-tertiary);
  text-transform: uppercase;
  letter-spacing: 1px;
}

.score-value {
  display: block;
  font-size: 20px;
  font-weight: 700;
  color: var(--color-text-primary);
  position: relative;
}

.score-popup {
  position: absolute;
  top: -8px;
  right: -10px;
  font-size: 14px;
  font-weight: 700;
  color: var(--accent-cyan);
  animation: score-float 0.8s ease-out forwards;
  pointer-events: none;
}

@keyframes score-float {
  0% { opacity: 1; transform: translateY(0); }
  100% { opacity: 0; transform: translateY(-20px); }
}

.game-actions { display: flex; gap: 8px; }

.btn-game {
  padding: 8px 16px;
  border: none;
  cursor: pointer;
  font-size: 14px;
  font-weight: 600;
  color: var(--color-text-primary);
  transition: var(--transition);
}

.btn-game:hover { color: var(--accent-cyan); }
.btn-close { padding: 8px 12px; }

.board {
  position: relative;
  width: min(100%, 416px);
  margin: 0 auto;
  padding: 8px;
  display: flex;
  justify-content: center;
  align-items: center;
  touch-action: none;
  overscroll-behavior: contain;
}

canvas {
  display: block;
  width: 100%;
  border-radius: 8px;
  height: auto;
}

.game-overlay {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  background: rgba(0, 0, 0, 0.45);
  backdrop-filter: blur(4px);
  border-radius: var(--radius-lg);
  z-index: 10;
  animation: overlay-fade 0.3s ease;
}

@keyframes overlay-fade {
  0% { opacity: 0; }
  100% { opacity: 1; }
}

.overlay-text {
  font-size: 28px;
  font-weight: 700;
  color: #fff;
}

.overlay-score {
  font-size: 18px;
  color: #fff;
  opacity: 0.9;
}

.mobile-controls {
  display: none;
  margin-top: 16px;
  gap: 8px;
  flex-direction: column;
  align-items: center;
}

.control-row {
  display: flex;
  gap: 8px;
}

.control-btn {
  width: 56px;
  height: 56px;
  border: none;
  cursor: pointer;
  font-size: 24px;
  color: var(--color-text-primary);
  transition: var(--transition);
  display: flex;
  align-items: center;
  justify-content: center;
}

.control-btn:active {
  transform: scale(0.95);
  color: var(--accent-cyan);
}

@media (max-width: 480px) {
  .game-snake { max-width: 100%; }
  .score-box { min-width: 64px; padding: 6px 10px; }
  .score-value { font-size: 16px; }
  .mobile-controls { display: flex; }
}
</style>
