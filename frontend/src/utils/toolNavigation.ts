import type { Router } from 'vue-router'

export function isResolvableToolRoute(router: Router, url?: string): boolean {
  if (!url || !url.startsWith('/')) {
    return false
  }
  return router.resolve(url).matched.length > 0
}
