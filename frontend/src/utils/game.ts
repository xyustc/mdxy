// 游戏工具函数
export function getOrCreatePlayerId(): string {
  const STORAGE_KEY = 'game_player_id'
  let playerId = localStorage.getItem(STORAGE_KEY)

  if (!playerId) {
    // 生成 UUID v4
    playerId = 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, (c) => {
      const r = Math.random() * 16 | 0
      const v = c === 'x' ? r : (r & 0x3 | 0x8)
      return v.toString(16)
    })
    localStorage.setItem(STORAGE_KEY, playerId)
  }

  return playerId
}
