/**
 * 防复制保护模块
 * 多层次防止用户复制、截图、打印笔记内容
 */

// 禁用的按键组合
const DISABLED_SHORTCUTS = [
  { ctrl: true, key: 'c' },      // Ctrl+C 复制
  { ctrl: true, key: 'x' },      // Ctrl+X 剪切
  { ctrl: true, key: 'a' },      // Ctrl+A 全选
  { ctrl: true, key: 's' },      // Ctrl+S 保存
  { ctrl: true, key: 'p' },      // Ctrl+P 打印
  { ctrl: true, key: 'u' },      // Ctrl+U 查看源代码
  { ctrl: true, shift: true, key: 'i' }, // Ctrl+Shift+I 开发者工具
  { ctrl: true, shift: true, key: 'j' }, // Ctrl+Shift+J 控制台
  { ctrl: true, shift: true, key: 'c' }, // Ctrl+Shift+C 元素检查
  { key: 'F12' },                // F12 开发者工具
  { key: 'PrintScreen' },        // 截图键
]

/**
 * 初始化防复制保护
 */
export function initCopyProtection() {
  // 1. 禁用右键菜单
  disableContextMenu()
  
  // 2. 禁用文本选择
  disableTextSelection()
  
  // 3. 禁用键盘快捷键
  disableKeyboardShortcuts()
  
  // 4. 禁用拖拽
  disableDragAndDrop()
  
  // 5. 禁用打印
  disablePrint()
  
  // 6. 检测开发者工具
  detectDevTools()
  
  // 7. 禁用复制事件
  disableCopyEvent()
  
  console.log('🔒 内容保护已启用')
}

/**
 * 禁用右键菜单
 */
function disableContextMenu() {
  document.addEventListener('contextmenu', (e) => {
    e.preventDefault()
    showWarning('右键菜单已禁用')
    return false
  })
}

/**
 * 禁用文本选择（通过 CSS 和 JS 双重保护）
 */
function disableTextSelection() {
  // JS 层面禁用
  document.addEventListener('selectstart', (e) => {
    // 允许输入框选择
    if (e.target.tagName === 'INPUT' || e.target.tagName === 'TEXTAREA') {
      return true
    }
    e.preventDefault()
    return false
  })
  
  // 禁用双击选择
  document.addEventListener('mousedown', (e) => {
    if (e.detail > 1) {
      e.preventDefault()
    }
  })
}

/**
 * 禁用键盘快捷键
 */
function disableKeyboardShortcuts() {
  document.addEventListener('keydown', (e) => {
    const key = e.key.toLowerCase()
    
    for (const shortcut of DISABLED_SHORTCUTS) {
      const ctrlMatch = shortcut.ctrl ? (e.ctrlKey || e.metaKey) : true
      const shiftMatch = shortcut.shift ? e.shiftKey : !e.shiftKey || shortcut.shift === undefined
      const keyMatch = shortcut.key.toLowerCase() === key || shortcut.key === e.key
      
      if (ctrlMatch && shiftMatch && keyMatch) {
        e.preventDefault()
        e.stopPropagation()
        
        if (key === 'c' || key === 'x') {
          showWarning('复制功能已禁用')
        } else if (key === 'p') {
          showWarning('打印功能已禁用')
        } else if (key === 'F12' || (shortcut.shift && key === 'i')) {
          showWarning('开发者工具已禁用')
        }
        
        return false
      }
    }
  }, true)
}

/**
 * 禁用拖拽
 */
function disableDragAndDrop() {
  document.addEventListener('dragstart', (e) => {
    e.preventDefault()
    return false
  })
  
  document.addEventListener('drop', (e) => {
    e.preventDefault()
    return false
  })
}

/**
 * 禁用打印
 */
function disablePrint() {
  // 监听打印前事件
  window.addEventListener('beforeprint', () => {
    document.body.style.visibility = 'hidden'
  })
  
  window.addEventListener('afterprint', () => {
    document.body.style.visibility = 'visible'
  })
  
  // 通过 CSS 媒体查询隐藏打印内容
  const style = document.createElement('style')
  style.textContent = `
    @media print {
      body * {
        display: none !important;
      }
      body::after {
        content: "打印功能已禁用";
        display: block !important;
        font-size: 24px;
        text-align: center;
        padding: 100px;
      }
    }
  `
  document.head.appendChild(style)
}

/**
 * 检测开发者工具
 */
function detectDevTools() {
  const threshold = 160
  
  const checkDevTools = () => {
    const widthThreshold = window.outerWidth - window.innerWidth > threshold
    const heightThreshold = window.outerHeight - window.innerHeight > threshold
    
    if (widthThreshold || heightThreshold) {
      // 开发者工具可能已打开
      // 可以选择：清空内容、跳转、显示警告等
      console.clear()
      console.log('%c⚠️ 检测到开发者工具', 'font-size: 24px; color: red;')
      console.log('%c请尊重知识产权，禁止复制内容', 'font-size: 16px; color: orange;')
    }
  }
  
  // 定期检测
  setInterval(checkDevTools, 1000)
  
  // 监听窗口大小变化
  window.addEventListener('resize', checkDevTools)
}

/**
 * 禁用复制事件
 */
function disableCopyEvent() {
  document.addEventListener('copy', (e) => {
    e.preventDefault()
    // 可以替换剪贴板内容
    e.clipboardData?.setData('text/plain', '复制功能已禁用，请尊重知识产权。')
    showWarning('复制功能已禁用')
    return false
  })
  
  document.addEventListener('cut', (e) => {
    e.preventDefault()
    showWarning('剪切功能已禁用')
    return false
  })
}

/**
 * 显示警告提示
 */
let warningTimeout = null
function showWarning(message) {
  // 移除已有的警告
  const existing = document.querySelector('.copy-warning')
  if (existing) {
    existing.remove()
  }
  
  // 创建警告元素
  const warning = document.createElement('div')
  warning.className = 'copy-warning'
  warning.innerHTML = `
    <span class="warning-icon">🔒</span>
    <span class="warning-text">${message}</span>
  `
  document.body.appendChild(warning)
  
  // 动画显示
  requestAnimationFrame(() => {
    warning.classList.add('show')
  })
  
  // 自动隐藏
  clearTimeout(warningTimeout)
  warningTimeout = setTimeout(() => {
    warning.classList.remove('show')
    setTimeout(() => warning.remove(), 300)
  }, 2000)
}

/**
 * 添加水印
 */
export function addWatermark(text = '仅供个人学习') {
  const canvas = document.createElement('canvas')
  const ctx = canvas.getContext('2d')
  
  canvas.width = 200
  canvas.height = 150
  
  ctx.rotate(-20 * Math.PI / 180)
  ctx.font = '14px Arial'
  ctx.fillStyle = 'rgba(180, 180, 180, 0.15)'
  ctx.textAlign = 'center'
  ctx.fillText(text, 100, 100)
  
  const watermarkDiv = document.createElement('div')
  watermarkDiv.className = 'watermark-layer'
  watermarkDiv.style.cssText = `
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    pointer-events: none;
    z-index: 9999;
    background-image: url(${canvas.toDataURL()});
    background-repeat: repeat;
  `
  
  document.body.appendChild(watermarkDiv)
  
  // 防止水印被删除
  const observer = new MutationObserver(() => {
    if (!document.querySelector('.watermark-layer')) {
      document.body.appendChild(watermarkDiv.cloneNode(true))
    }
  })
  
  observer.observe(document.body, { childList: true })
}
