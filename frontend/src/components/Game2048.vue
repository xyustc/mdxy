<template>
  <div class="game-2048">
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
        <button class="btn-game glass-card" @click="resetGame">新游戏</button>
        <button class="btn-game glass-card btn-close" aria-label="关闭 2048 游戏" @click="$emit('close')">✕</button>
      </div>
    </div>

    <div class="board glass-card" ref="boardEl">
      <div class="board-grid">
        <div v-for="r in 4" :key="'bg-'+r" class="board-row">
          <div v-for="c in 4" :key="'bg-'+r+'-'+c" class="cell-bg"></div>
        </div>
      </div>
      <TransitionGroup name="tile" tag="div" class="tiles-layer">
        <div
          v-for="tile in flatTiles"
          :key="tile.id"
          class="tile"
          :class="['tile-' + tile.value, { 'tile-merged': tile.merged, 'tile-new': tile.isNew }]"
          :style="{ '--r': tile.row, '--c': tile.col }"
        >
          {{ tile.value }}
        </div>
      </TransitionGroup>
    </div>

    <div v-if="gameOver || showWinOverlay" class="game-overlay glass-card">
      <p class="overlay-text">{{ showWinOverlay ? '🎉 你赢了！' : '游戏结束' }}</p>
      <div class="overlay-actions">
        <button v-if="showWinOverlay" class="btn-game glass-card" @click="continueGame">继续挑战</button>
        <button class="btn-game glass-card" @click="resetGame">再来一局</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { gameApi } from '@/api/game'
import { getOrCreatePlayerId } from '@/utils/game'

defineEmits<{ close: [] }>()

interface TileData {
  id: number
  value: number
  row: number
  col: number
  merged: boolean
  isNew: boolean
}

let tileId = 0
const TOOL_ID = 1 // 2048 游戏的 tool_id
const playerId = getOrCreatePlayerId()
const MOVE_INTERVAL = 150 // 操作防抖间隔 ms
const SWIPE_TRIGGER_DISTANCE = 30

const board = ref<(TileData | null)[][]>(createEmptyBoard())
const score = ref(0)
const bestScore = ref(0)
const gameOver = ref(false)
const won = ref(false)
const keepPlaying = ref(false)
const boardEl = ref<HTMLElement>()
const scorePopup = ref(0)
const scorePopupKey = ref(0)
let lastMoveTime = 0

const showWinOverlay = computed(() => won.value && !keepPlaying.value)

const flatTiles = computed(() =>
  board.value.flat().filter((t): t is TileData => t !== null)
)

function createEmptyBoard(): (TileData | null)[][] {
  return Array.from({ length: 4 }, () => Array(4).fill(null))
}

function addRandomTile() {
  const empty: [number, number][] = []
  for (let r = 0; r < 4; r++)
    for (let c = 0; c < 4; c++)
      if (!board.value[r][c]) empty.push([r, c])
  if (!empty.length) return
  const [r, c] = empty[Math.floor(Math.random() * empty.length)]
  const value = Math.random() < 0.9 ? 2 : 4
  board.value[r][c] = { id: ++tileId, value, row: r, col: c, merged: false, isNew: true }
}

function clearFlags() {
  for (let r = 0; r < 4; r++)
    for (let c = 0; c < 4; c++) {
      const t = board.value[r][c]
      if (t) { t.merged = false; t.isNew = false }
    }
}

function move(dir: 'up' | 'down' | 'left' | 'right') {
  if (gameOver.value || showWinOverlay.value) return
  const now = Date.now()
  if (now - lastMoveTime < MOVE_INTERVAL) return
  lastMoveTime = now

  clearFlags()
  let moved = false
  let moveScore = 0
  const b = board.value

  const traverse = (cb: (r: number, c: number) => void) => {
    const rows = dir === 'down' ? [3, 2, 1, 0] : [0, 1, 2, 3]
    const cols = dir === 'right' ? [3, 2, 1, 0] : [0, 1, 2, 3]
    for (const r of rows) for (const c of cols) cb(r, c)
  }

  const vector = { up: [-1, 0], down: [1, 0], left: [0, -1], right: [0, 1] }[dir]

  traverse((r, c) => {
    const tile = b[r][c]
    if (!tile) return
    let nr = r, nc = c
    while (true) {
      const tr = nr + vector[0], tc = nc + vector[1]
      if (tr < 0 || tr > 3 || tc < 0 || tc > 3) break
      if (!b[tr][tc]) { nr = tr; nc = tc; continue }
      if (b[tr][tc]!.value === tile.value && !b[tr][tc]!.merged) { nr = tr; nc = tc }
      break
    }
    if (nr !== r || nc !== c) {
      moved = true
      b[r][c] = null
      const target = b[nr][nc]
      if (target && target.value === tile.value) {
        const newVal = tile.value * 2
        b[nr][nc] = { id: ++tileId, value: newVal, row: nr, col: nc, merged: true, isNew: false }
        score.value += newVal
        moveScore += newVal
        if (newVal === 2048 && !keepPlaying.value) won.value = true
      } else {
        tile.row = nr; tile.col = nc
        b[nr][nc] = tile
      }
    }
  })

  if (moved) {
    addRandomTile()
    if (moveScore > 0) {
      scorePopup.value = moveScore
      scorePopupKey.value++
      setTimeout(() => { scorePopup.value = 0 }, 800)
    }
    if (score.value > bestScore.value) {
      bestScore.value = score.value
      syncScoreToServer()
    }
    checkGameOver()
  }
}

async function syncScoreToServer() {
  try {
    const res = await gameApi.updateScore({
      tool_id: TOOL_ID,
      player_id: playerId,
      score: score.value
    })
    if (res.success && res.data) {
      bestScore.value = res.data.best_score
    }
  } catch (err) {
    console.error('Failed to sync score:', err)
  }
}

async function loadBestScore() {
  try {
    const res = await gameApi.getBestScore(TOOL_ID, playerId)
    if (res.success && res.data) {
      bestScore.value = res.data.best_score
    }
  } catch (err) {
    console.error('Failed to load best score:', err)
  }
}

function checkGameOver() {
  for (let r = 0; r < 4; r++)
    for (let c = 0; c < 4; c++) {
      if (!board.value[r][c]) return
      const v = board.value[r][c]!.value
      if (r < 3 && board.value[r + 1][c]?.value === v) return
      if (c < 3 && board.value[r][c + 1]?.value === v) return
    }
  gameOver.value = true
}

function resetGame() {
  board.value = createEmptyBoard()
  score.value = 0
  gameOver.value = false
  won.value = false
  keepPlaying.value = false
  scorePopup.value = 0
  addRandomTile()
  addRandomTile()
}

function continueGame() {
  keepPlaying.value = true
}

function onKeyDown(e: KeyboardEvent) {
  const map: Record<string, 'up' | 'down' | 'left' | 'right'> = {
    ArrowUp: 'up', ArrowDown: 'down', ArrowLeft: 'left', ArrowRight: 'right',
    w: 'up', s: 'down', a: 'left', d: 'right'
  }
  const dir = map[e.key]
  if (dir) { e.preventDefault(); move(dir) }
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
  const absDx = Math.abs(dx), absDy = Math.abs(dy)
  if (Math.max(absDx, absDy) < SWIPE_TRIGGER_DISTANCE) return
  if (absDx > absDy) move(dx > 0 ? 'right' : 'left')
  else move(dy > 0 ? 'down' : 'up')
}

function onTouchCancel() {
  touchStartX = 0
  touchStartY = 0
}

onMounted(() => {
  loadBestScore()
  resetGame()
  window.addEventListener('keydown', onKeyDown)
  nextTick(() => {
    boardEl.value?.addEventListener('touchstart', onTouchStart, { passive: true })
    boardEl.value?.addEventListener('touchmove', onTouchMove, { passive: false })
    boardEl.value?.addEventListener('touchend', onTouchEnd, { passive: true })
    boardEl.value?.addEventListener('touchcancel', onTouchCancel, { passive: true })
  })
})

onUnmounted(() => {
  window.removeEventListener('keydown', onKeyDown)
  boardEl.value?.removeEventListener('touchstart', onTouchStart)
  boardEl.value?.removeEventListener('touchmove', onTouchMove)
  boardEl.value?.removeEventListener('touchend', onTouchEnd)
  boardEl.value?.removeEventListener('touchcancel', onTouchCancel)
})
</script>

<style scoped>
.game-2048 {
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
  padding: 8px;
  aspect-ratio: 1;
  touch-action: none;
  overscroll-behavior: contain;
}

.board-grid { display: flex; flex-direction: column; gap: 8px; }

.board-row { display: flex; gap: 8px; }

.cell-bg {
  width: calc((100% - 24px) / 4);
  aspect-ratio: 1;
  border-radius: 8px;
  background: var(--color-bg-tertiary);
  flex-shrink: 0;
}

.tiles-layer {
  position: absolute;
  inset: 8px;
}

.tile {
  position: absolute;
  width: calc((100% - 24px) / 4);
  aspect-ratio: 1;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 28px;
  color: var(--color-text-primary);
  top: calc(var(--r) * (25% + 2px));
  left: calc(var(--c) * (25% + 2px));
  transition: top 0.15s ease, left 0.15s ease;
}

.tile-2 { background: #eee4da; color: #776e65; }
.tile-4 { background: #ede0c8; color: #776e65; }
.tile-8 { background: #f2b179; color: #f9f6f2; }
.tile-16 { background: #f59563; color: #f9f6f2; }
.tile-32 { background: #f67c5f; color: #f9f6f2; }
.tile-64 { background: #f65e3b; color: #f9f6f2; }
.tile-128 { background: #edcf72; color: #f9f6f2; font-size: 24px; }
.tile-256 { background: #edcc61; color: #f9f6f2; font-size: 24px; }
.tile-512 { background: #edc850; color: #f9f6f2; font-size: 24px; }
.tile-1024 { background: #edc53f; color: #f9f6f2; font-size: 20px; }
.tile-2048 { background: #edc22e; color: #f9f6f2; font-size: 20px; }

/* 超过 2048 的高数值方块兜底 */
.tile-4096 { background: #3c3a32; color: #f9f6f2; font-size: 18px; }
.tile-8192 { background: #3c3a32; color: #f9f6f2; font-size: 18px; }

.tile-new {
  animation: tile-appear 0.2s ease;
}

.tile-merged {
  animation: tile-pop 0.25s cubic-bezier(0.18, 0.89, 0.32, 1.28);
}

@keyframes tile-appear {
  0% { transform: scale(0); opacity: 0; }
  60% { transform: scale(1.05); opacity: 1; }
  100% { transform: scale(1); opacity: 1; }
}

@keyframes tile-pop {
  0% { transform: scale(0.8); }
  40% { transform: scale(1.2); }
  100% { transform: scale(1); }
}

.game-overlay {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 16px;
  background: rgba(0, 0, 0, 0.45);
  backdrop-filter: blur(4px);
  border-radius: var(--radius-lg);
  z-index: 10;
  animation: overlay-fade 0.3s ease;
}

.overlay-actions {
  display: flex;
  gap: 12px;
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

@media (max-width: 480px) {
  .game-2048 { max-width: 100%; }
  .tile { font-size: 20px; }
  .tile-128, .tile-256, .tile-512 { font-size: 18px; }
  .tile-1024, .tile-2048 { font-size: 15px; }
  .score-box { min-width: 64px; padding: 6px 10px; }
  .score-value { font-size: 16px; }
}
</style>
