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

type Direction = 'up' | 'down' | 'left' | 'right'
type BoardState = (TileData | null)[][]

let tileId = 0
const TOOL_ID = 1 // 2048 游戏的 tool_id
const playerId = getOrCreatePlayerId()
const MOVE_INTERVAL = 100
const SCORE_SYNC_DELAY = 800
const SWIPE_TRIGGER_DISTANCE = 30

const board = ref<BoardState>(createEmptyBoard())
const score = ref(0)
const bestScore = ref(0)
const gameOver = ref(false)
const won = ref(false)
const keepPlaying = ref(false)
const boardEl = ref<HTMLElement>()
const scorePopup = ref(0)
const scorePopupKey = ref(0)
let lastMoveTime = 0
let queuedDirection: Direction | null = null
let moveQueueTimer: ReturnType<typeof setTimeout> | null = null
let scorePopupTimer: ReturnType<typeof setTimeout> | null = null
let scoreSyncTimer: ReturnType<typeof setTimeout> | null = null
let pendingSyncScore: number | null = null
let syncInFlight = false

const showWinOverlay = computed(() => won.value && !keepPlaying.value)

const flatTiles = computed(() =>
  board.value.flat().filter((t): t is TileData => t !== null)
)

function createEmptyBoard(): BoardState {
  return Array.from({ length: 4 }, () => Array(4).fill(null))
}

function cloneBoardState(source: BoardState): BoardState {
  return source.map((row) => row.map((tile) => (tile ? { ...tile, merged: false, isNew: false } : null)))
}

function addRandomTile(targetBoard: BoardState) {
  const empty: [number, number][] = []
  for (let r = 0; r < 4; r++)
    for (let c = 0; c < 4; c++)
      if (!targetBoard[r][c]) empty.push([r, c])
  if (!empty.length) return
  const [r, c] = empty[Math.floor(Math.random() * empty.length)]
  const value = Math.random() < 0.9 ? 2 : 4
  targetBoard[r][c] = { id: ++tileId, value, row: r, col: c, merged: false, isNew: true }
}

function scheduleQueuedMove() {
  if (!queuedDirection || moveQueueTimer) return
  const elapsed = Date.now() - lastMoveTime
  const delay = Math.max(0, MOVE_INTERVAL - elapsed)
  moveQueueTimer = setTimeout(() => {
    moveQueueTimer = null
    const dir = queuedDirection
    queuedDirection = null
    if (dir) move(dir)
  }, delay)
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
    const res = await gameApi.updateScore({
      tool_id: TOOL_ID,
      player_id: playerId,
      score: syncScore
    })
    if (res.success && res.data) {
      bestScore.value = Math.max(bestScore.value, res.data.best_score)
    }
  } catch (err) {
    console.error('Failed to sync score:', err)
  } finally {
    syncInFlight = false
    if (pendingSyncScore !== null) {
      void flushScoreSync()
    }
  }
}

function move(dir: Direction) {
  if (gameOver.value || showWinOverlay.value) return
  const now = Date.now()
  if (now - lastMoveTime < MOVE_INTERVAL) {
    queuedDirection = dir
    scheduleQueuedMove()
    return
  }
  lastMoveTime = now

  const b = cloneBoardState(board.value)
  let moved = false
  let moveScore = 0

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
      if (b[tr][tc]!.value === tile.value && !b[tr][tc]!.merged) {
        nr = tr
        nc = tc
      }
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
        b[nr][nc] = { ...tile, row: nr, col: nc }
      }
    }
  })

  if (moved) {
    addRandomTile(b)
    board.value = b
    if (moveScore > 0) {
      scorePopup.value = moveScore
      scorePopupKey.value++
      if (scorePopupTimer) clearTimeout(scorePopupTimer)
      scorePopupTimer = setTimeout(() => { scorePopup.value = 0 }, 500)
    }
    if (score.value > bestScore.value) {
      bestScore.value = score.value
      scheduleScoreSync()
    }
    checkGameOver()
  }
  if (queuedDirection) {
    scheduleQueuedMove()
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

function checkGameOver(state: BoardState = board.value) {
  for (let r = 0; r < 4; r++)
    for (let c = 0; c < 4; c++) {
      if (!state[r][c]) return
      const v = state[r][c]!.value
      if (r < 3 && state[r + 1][c]?.value === v) return
      if (c < 3 && state[r][c + 1]?.value === v) return
    }
  gameOver.value = true
}

function resetGame() {
  const next = createEmptyBoard()
  addRandomTile(next)
  addRandomTile(next)
  board.value = next
  score.value = 0
  gameOver.value = false
  won.value = false
  keepPlaying.value = false
  scorePopup.value = 0
  queuedDirection = null
  if (moveQueueTimer) {
    clearTimeout(moveQueueTimer)
    moveQueueTimer = null
  }
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
  const absDx = Math.abs(dx), absDy = Math.abs(dy)
  if (Math.max(absDx, absDy) < SWIPE_TRIGGER_DISTANCE) return
  if (absDx > absDy) move(dx > 0 ? 'right' : 'left')
  else move(dy > 0 ? 'down' : 'up')
}

function onTouchCancel() {
  touchStartX = 0
  touchStartY = 0
}

function onVisibilityChange() {
  if (document.visibilityState === 'hidden' && pendingSyncScore !== null) {
    scheduleScoreSync(true)
  }
}

onMounted(() => {
  loadBestScore()
  resetGame()
  window.addEventListener('keydown', onKeyDown)
  document.addEventListener('visibilitychange', onVisibilityChange)
  nextTick(() => {
    boardEl.value?.addEventListener('touchstart', onTouchStart, { passive: true })
    boardEl.value?.addEventListener('touchmove', onTouchMove, { passive: false })
    boardEl.value?.addEventListener('touchend', onTouchEnd, { passive: true })
    boardEl.value?.addEventListener('touchcancel', onTouchCancel, { passive: true })
  })
})

onUnmounted(() => {
  window.removeEventListener('keydown', onKeyDown)
  document.removeEventListener('visibilitychange', onVisibilityChange)
  boardEl.value?.removeEventListener('touchstart', onTouchStart)
  boardEl.value?.removeEventListener('touchmove', onTouchMove)
  boardEl.value?.removeEventListener('touchend', onTouchEnd)
  boardEl.value?.removeEventListener('touchcancel', onTouchCancel)
  if (moveQueueTimer) clearTimeout(moveQueueTimer)
  if (scorePopupTimer) clearTimeout(scorePopupTimer)
  if (scoreSyncTimer) clearTimeout(scoreSyncTimer)
  if (pendingSyncScore !== null) {
    void flushScoreSync()
  }
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
  top: 0;
  left: 0;
  transform: translate3d(calc(var(--c) * (100% + 8px)), calc(var(--r) * (100% + 8px)), 0);
  transition: transform 0.11s cubic-bezier(0.22, 0.61, 0.36, 1);
  will-change: transform;
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
  animation: tile-appear 0.16s ease;
}

.tile-merged {
  animation: tile-pop 0.18s ease;
}

@keyframes tile-appear {
  0% { opacity: 0.5; }
  100% { opacity: 1; }
}

@keyframes tile-pop {
  0% { filter: brightness(1); }
  50% { filter: brightness(1.16); }
  100% { filter: brightness(1); }
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
