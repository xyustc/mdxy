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
        <button class="btn-game glass-card" @click="togglePause" v-if="!gameOver">{{ paused ? '继续' : '暂停' }}</button>
        <button class="btn-game glass-card" @click="resetGame">新游戏</button>
        <button class="btn-game glass-card btn-close" aria-label="关闭贪吃蛇游戏" @click="$emit('close')">✕</button>
      </div>
    </div>

    <div class="board glass-card" ref="boardEl">
      <canvas ref="canvasEl" :width="CANVAS_SIZE" :height="CANVAS_SIZE"></canvas>
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
import { ref, onMounted, onUnmounted } from 'vue'
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

const canvasEl = ref<HTMLCanvasElement>()
const boardEl = ref<HTMLElement>()
const score = ref(0)
const bestScore = ref(0)
const gameOver = ref(false)
const paused = ref(false)
const scorePopup = ref(0)
const scorePopupKey = ref(0)

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

function initGame() {
  snake = [{ x: 10, y: 10 }, { x: 9, y: 10 }, { x: 8, y: 10 }]
  direction = { x: 1, y: 0 }
  directionQueue = []
  score.value = 0
  gameOver.value = false
  paused.value = false
  scorePopup.value = 0
  speed = INITIAL_SPEED
  dying = false
  deathFlashCount = 0
  spawnFood()
}

function spawnFood() {
  do {
    food = { x: Math.floor(Math.random() * GRID_SIZE), y: Math.floor(Math.random() * GRID_SIZE) }
  } while (snake.some(s => s.x === food.x && s.y === food.y))
}

function changeDirection(dir: 'up' | 'down' | 'left' | 'right') {
  if (paused.value || gameOver.value) return
  const vectors = { up: { x: 0, y: -1 }, down: { x: 0, y: 1 }, left: { x: -1, y: 0 }, right: { x: 1, y: 0 } }
  const newDir = vectors[dir]
  // 与当前方向或队列中最后一个方向比较，防止反向
  const lastDir = directionQueue.length > 0 ? directionQueue[directionQueue.length - 1] : direction
  if (lastDir.x + newDir.x !== 0 || lastDir.y + newDir.y !== 0) {
    if (directionQueue.length < 3) {
      directionQueue.push(newDir)
    }
  }
}

function togglePause() {
  if (gameOver.value) return
  paused.value = !paused.value
  if (!paused.value) {
    lastTime = performance.now()
    accumulated = 0
  }
  draw()
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
    setTimeout(() => { scorePopup.value = 0 }, 800)
    if (score.value > bestScore.value) {
      bestScore.value = score.value
      syncScoreToServer()
    }
    spawnFood()
    speed = Math.max(50, speed - 2)
  } else {
    snake.pop()
  }
}

function startDeathAnimation() {
  dying = true
  deathFlashCount = 0
  const flashInterval = setInterval(() => {
    deathFlashCount++
    draw()
    if (deathFlashCount >= 6) {
      clearInterval(flashInterval)
      dying = false
      gameOver.value = true
      syncScoreToServer()
      draw()
    }
  }, 100)
}

function draw() {
  const canvas = canvasEl.value
  if (!canvas) return
  const ctx = canvas.getContext('2d')
  if (!ctx) return

  const bgColor = getComputedStyle(document.documentElement).getPropertyValue('--color-bg-tertiary').trim()
  ctx.fillStyle = bgColor
  ctx.fillRect(0, 0, CANVAS_SIZE, CANVAS_SIZE)

  // 网格线
  ctx.strokeStyle = 'rgba(128, 128, 128, 0.15)'
  ctx.lineWidth = 0.5
  for (let i = 1; i < GRID_SIZE; i++) {
    const pos = i * CELL_SIZE
    ctx.beginPath(); ctx.moveTo(pos, 0); ctx.lineTo(pos, CANVAS_SIZE); ctx.stroke()
    ctx.beginPath(); ctx.moveTo(0, pos); ctx.lineTo(CANVAS_SIZE, pos); ctx.stroke()
  }

  // 食物呼吸动画
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

  // 暂停遮罩
  if (paused.value) {
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

async function syncScoreToServer() {
  try {
    const res = await gameApi.updateScore({ tool_id: TOOL_ID, player_id: playerId, score: score.value })
    if (res.success && res.data) bestScore.value = res.data.best_score
  } catch (err) { console.error('Failed to sync score:', err) }
}

async function loadBestScore() {
  try {
    const res = await gameApi.getBestScore(TOOL_ID, playerId)
    if (res.success && res.data) bestScore.value = res.data.best_score
  } catch (err) { console.error('Failed to load best score:', err) }
}

function gameFrame(time: number) {
  rafId = requestAnimationFrame(gameFrame)
  if (gameOver.value || paused.value || dying) return
  if (!lastTime) { lastTime = time; return }
  accumulated += time - lastTime
  lastTime = time
  if (accumulated >= speed) {
    accumulated -= speed
    update()
  }
  draw()
}

function resetGame() {
  if (rafId) cancelAnimationFrame(rafId)
  initGame()
  lastTime = 0
  accumulated = 0
  draw()
  rafId = requestAnimationFrame(gameFrame)
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
  if (e.cancelable) {
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

onMounted(() => {
  loadBestScore()
  initGame()
  draw()
  rafId = requestAnimationFrame(gameFrame)
  window.addEventListener('keydown', onKeyDown)
  boardEl.value?.addEventListener('touchstart', onTouchStart, { passive: true })
  boardEl.value?.addEventListener('touchmove', onTouchMove, { passive: false })
  boardEl.value?.addEventListener('touchend', onTouchEnd, { passive: true })
  boardEl.value?.addEventListener('touchcancel', onTouchCancel, { passive: true })
})

onUnmounted(() => {
  if (rafId) cancelAnimationFrame(rafId)
  window.removeEventListener('keydown', onKeyDown)
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
