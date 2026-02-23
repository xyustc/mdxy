const STORAGE_KEY = 'search_history'
const MAX_ITEMS = 10

export function useSearchHistory() {
  function getHistory(): string[] {
    try {
      const raw = localStorage.getItem(STORAGE_KEY)
      return raw ? JSON.parse(raw) : []
    } catch {
      return []
    }
  }

  function addHistory(query: string) {
    if (!query.trim()) return
    const history = getHistory().filter(h => h !== query)
    history.unshift(query)
    if (history.length > MAX_ITEMS) history.pop()
    localStorage.setItem(STORAGE_KEY, JSON.stringify(history))
  }

  function clearHistory() {
    localStorage.removeItem(STORAGE_KEY)
  }

  return { getHistory, addHistory, clearHistory }
}
