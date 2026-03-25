import type { Directive, DirectiveBinding } from 'vue'

type OverflowTooltipEl = HTMLElement & {
  __overflowTooltipState__?: OverflowTooltipState
}

type OverflowTooltipState = {
  text?: string
  hideTimer: number | null
  onPointerEnter: (event: PointerEvent) => void
  onPointerMove: (event: PointerEvent) => void
  onPointerLeave: () => void
  onTouchStart: (event: TouchEvent) => void
}

let tooltipEl: HTMLDivElement | null = null
let activeEl: OverflowTooltipEl | null = null

function ensureTooltipEl() {
  if (tooltipEl) return tooltipEl

  tooltipEl = document.createElement('div')
  tooltipEl.className = 'overflow-tooltip-layer'
  document.body.appendChild(tooltipEl)
  return tooltipEl
}

function getTooltipText(el: OverflowTooltipEl) {
  const customText = el.__overflowTooltipState__?.text?.trim()
  if (customText) return customText
  return el.textContent?.trim() || ''
}

function isOverflowing(el: HTMLElement) {
  const style = window.getComputedStyle(el)
  const lineClampValue =
    style.getPropertyValue('-webkit-line-clamp') || (style as CSSStyleDeclaration & { webkitLineClamp?: string }).webkitLineClamp || ''
  const lineClamp = Number.parseInt(lineClampValue, 10)

  if (Number.isFinite(lineClamp) && lineClamp > 0 && el.clientWidth > 0) {
    const clone = el.cloneNode(true) as HTMLElement
    clone.style.position = 'fixed'
    clone.style.left = '-9999px'
    clone.style.top = '0'
    clone.style.visibility = 'hidden'
    clone.style.pointerEvents = 'none'
    clone.style.width = `${el.clientWidth}px`
    clone.style.minHeight = '0'
    clone.style.maxHeight = 'none'
    clone.style.height = 'auto'
    clone.style.overflow = 'visible'
    clone.style.display = 'block'
    clone.style.webkitLineClamp = 'unset'
    clone.style.webkitBoxOrient = 'initial'

    document.body.appendChild(clone)
    const isClamped = clone.getBoundingClientRect().height - el.getBoundingClientRect().height > 1
    clone.remove()

    if (isClamped) return true
  }

  return el.scrollWidth - el.clientWidth > 2 || el.scrollHeight - el.clientHeight > 2
}

function positionTooltip(x: number, y: number) {
  const layer = ensureTooltipEl()
  const gap = 14

  layer.style.left = '0px'
  layer.style.top = '0px'
  layer.style.visibility = 'hidden'
  layer.dataset.visible = 'true'

  const { width, height } = layer.getBoundingClientRect()
  const maxLeft = window.innerWidth - width - gap
  const maxTop = window.innerHeight - height - gap

  const left = Math.max(gap, Math.min(x, maxLeft))
  const top = Math.max(gap, Math.min(y, maxTop))

  layer.style.left = `${left}px`
  layer.style.top = `${top}px`
  layer.style.visibility = 'visible'
}

function showTooltip(el: OverflowTooltipEl, text: string, x: number, y: number) {
  const layer = ensureTooltipEl()
  layer.textContent = text
  activeEl = el
  positionTooltip(x, y)
}

function hideTooltip(el?: OverflowTooltipEl) {
  if (!tooltipEl) return
  if (el && activeEl && el !== activeEl) return
  tooltipEl.dataset.visible = 'false'
  tooltipEl.style.visibility = 'hidden'
  activeEl = null
}

function clearHideTimer(el: OverflowTooltipEl) {
  const timer = el.__overflowTooltipState__?.hideTimer
  if (timer) {
    window.clearTimeout(timer)
    if (el.__overflowTooltipState__) {
      el.__overflowTooltipState__.hideTimer = null
    }
  }
}

function updateBindingText(el: OverflowTooltipEl, binding: DirectiveBinding<string | undefined>) {
  if (el.__overflowTooltipState__) {
    el.__overflowTooltipState__.text = binding.value
  }
}

export const overflowTooltipDirective: Directive<OverflowTooltipEl, string | undefined> = {
  mounted(el, binding) {
    const state: OverflowTooltipState = {
      text: binding.value,
      hideTimer: null,
      onPointerEnter(event) {
        if (event.pointerType !== 'mouse') return
        clearHideTimer(el)
        if (!isOverflowing(el)) return

        const text = getTooltipText(el)
        if (!text) return
        showTooltip(el, text, event.clientX + 16, event.clientY + 18)
      },
      onPointerMove(event) {
        if (activeEl !== el) return
        positionTooltip(event.clientX + 16, event.clientY + 18)
      },
      onPointerLeave() {
        hideTooltip(el)
      },
      onTouchStart(event) {
        clearHideTimer(el)
        if (!isOverflowing(el)) return

        const text = getTooltipText(el)
        if (!text) return

        const touch = event.touches[0]
        const rect = el.getBoundingClientRect()
        const x = touch?.clientX ?? rect.left + rect.width / 2
        const y = rect.top - 10

        showTooltip(el, text, x, y)
        state.hideTimer = window.setTimeout(() => hideTooltip(el), 2200)
      }
    }

    el.__overflowTooltipState__ = state
    el.addEventListener('pointerenter', state.onPointerEnter)
    el.addEventListener('pointermove', state.onPointerMove)
    el.addEventListener('pointerleave', state.onPointerLeave)
    el.addEventListener('touchstart', state.onTouchStart, { passive: true })
  },

  updated(el, binding) {
    updateBindingText(el, binding)
  },

  unmounted(el) {
    const state = el.__overflowTooltipState__
    if (!state) return

    clearHideTimer(el)
    el.removeEventListener('pointerenter', state.onPointerEnter)
    el.removeEventListener('pointermove', state.onPointerMove)
    el.removeEventListener('pointerleave', state.onPointerLeave)
    el.removeEventListener('touchstart', state.onTouchStart)
    hideTooltip(el)
    delete el.__overflowTooltipState__
  }
}
