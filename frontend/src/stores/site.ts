import { computed, ref } from 'vue'
import { defineStore } from 'pinia'
import { profileApi } from '@/api/profile'
import { noteApi } from '@/api/note'
import { toolApi } from '@/api/tool'
import type { NoteNode, Profile, Tool } from '@/api/types'

type CachedPayload<T> = {
  data: T
  cachedAt: number
}

type HomeSnapshot = {
  featuredNotes: NoteNode[]
  noteCount: number
  featuredTools: Tool[]
  toolCount: number
}

const PROFILE_CACHE_KEY = 'mdxy.site.profile'
const HOME_CACHE_KEY = 'mdxy.site.home'
const PROFILE_TTL = 1000 * 60 * 60 * 6
const HOME_TTL = 1000 * 60 * 10

let profileRequest: Promise<Profile | null> | null = null
let homeRequest: Promise<HomeSnapshot | null> | null = null

function readCache<T>(key: string, ttl: number): T | null {
  if (typeof window === 'undefined') return null

  try {
    const raw = window.localStorage.getItem(key)
    if (!raw) return null

    const cached = JSON.parse(raw) as CachedPayload<T>
    if (!cached?.cachedAt || Date.now() - cached.cachedAt > ttl) {
      window.localStorage.removeItem(key)
      return null
    }

    return cached.data ?? null
  } catch {
    return null
  }
}

function writeCache<T>(key: string, data: T) {
  if (typeof window === 'undefined') return

  try {
    const payload: CachedPayload<T> = {
      data,
      cachedAt: Date.now()
    }
    window.localStorage.setItem(key, JSON.stringify(payload))
  } catch {
    // Ignore storage failures and keep network data in memory only.
  }
}

function flattenNotes(nodes: NoteNode[]): NoteNode[] {
  return nodes.flatMap((node) => (node.type === 'file' ? [node] : flattenNotes(node.children || [])))
}

export const useSiteStore = defineStore('site', () => {
  const profile = ref<Profile | null>(null)
  const homeSnapshot = ref<HomeSnapshot | null>(null)
  const profileLoaded = ref(false)
  const homeLoaded = ref(false)

  const profileName = computed(() => profile.value?.name?.trim() || '')
  const featuredNotes = computed(() => homeSnapshot.value?.featuredNotes || [])
  const noteCount = computed(() => homeSnapshot.value?.noteCount || 0)
  const featuredTools = computed(() => homeSnapshot.value?.featuredTools || [])
  const toolCount = computed(() => homeSnapshot.value?.toolCount || 0)

  function hydrateFromCache() {
    const cachedProfile = readCache<Profile>(PROFILE_CACHE_KEY, PROFILE_TTL)
    if (cachedProfile) {
      profile.value = cachedProfile
    }

    const cachedHome = readCache<HomeSnapshot>(HOME_CACHE_KEY, HOME_TTL)
    if (cachedHome) {
      homeSnapshot.value = cachedHome
    }
  }

  async function ensureProfile(force = false) {
    if (!force && profile.value) {
      profileLoaded.value = true
      return profile.value
    }

    if (profileRequest && !force) {
      return profileRequest
    }

    profileRequest = profileApi
      .get()
      .then((res) => {
        const data = res.success && res.data ? res.data : null
        if (data) {
          profile.value = data
          writeCache(PROFILE_CACHE_KEY, data)
        }
        profileLoaded.value = true
        return data
      })
      .catch(() => {
        profileLoaded.value = true
        return profile.value
      })
      .finally(() => {
        profileRequest = null
      })

    return profileRequest
  }

  async function ensureHome(force = false) {
    if (!force && homeSnapshot.value) {
      homeLoaded.value = true
      return homeSnapshot.value
    }

    if (homeRequest && !force) {
      return homeRequest
    }

    homeRequest = Promise.allSettled([noteApi.getTree(), toolApi.list()])
      .then((results) => {
        const [notesRes, toolsRes] = results
        const nextSnapshot: HomeSnapshot = {
          featuredNotes: homeSnapshot.value?.featuredNotes || [],
          noteCount: homeSnapshot.value?.noteCount || 0,
          featuredTools: homeSnapshot.value?.featuredTools || [],
          toolCount: homeSnapshot.value?.toolCount || 0
        }

        if (notesRes.status === 'fulfilled' && notesRes.value.success && notesRes.value.data) {
          const flattened = flattenNotes(notesRes.value.data)
          nextSnapshot.noteCount = flattened.length
          nextSnapshot.featuredNotes = flattened.slice(0, 2)
        }

        if (toolsRes.status === 'fulfilled' && toolsRes.value.success && toolsRes.value.data) {
          const visibleTools = (toolsRes.value.data || []).filter((tool) => tool.is_visible)
          nextSnapshot.toolCount = visibleTools.length
          nextSnapshot.featuredTools = visibleTools.slice(0, 2)
        }

        homeSnapshot.value = nextSnapshot
        writeCache(HOME_CACHE_KEY, nextSnapshot)
        homeLoaded.value = true
        return nextSnapshot
      })
      .catch(() => {
        homeLoaded.value = true
        return homeSnapshot.value
      })
      .finally(() => {
        homeRequest = null
      })

    return homeRequest
  }

  return {
    profile,
    homeSnapshot,
    profileLoaded,
    homeLoaded,
    profileName,
    featuredNotes,
    noteCount,
    featuredTools,
    toolCount,
    hydrateFromCache,
    ensureProfile,
    ensureHome
  }
})
