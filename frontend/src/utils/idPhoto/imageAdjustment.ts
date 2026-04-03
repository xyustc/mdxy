/**
 * 图像调整工具
 * 使用 Canvas 2D API 进行亮度、对比度调整
 */

export interface ImageAdjustments {
  brightness: number // -100 to 100
  contrast: number // -100 to 100
}

export const DEFAULT_ADJUSTMENTS: ImageAdjustments = {
  brightness: 0,
  contrast: 0
}

/**
 * 将调整值转换为 CSS filter 字符串
 */
export function adjustmentsToFilter(adjustments: ImageAdjustments): string {
  const brightness = 1 + adjustments.brightness / 100
  const contrast = 1 + adjustments.contrast / 100

  return `brightness(${brightness}) contrast(${contrast})`
}

/**
 * 应用背景色到透明图像
 */
export function applyBackground(
  source: HTMLImageElement | HTMLCanvasElement | ImageBitmap,
  backgroundColor: string
): HTMLCanvasElement {
  const canvas = document.createElement('canvas')
  canvas.width = source.width
  canvas.height = source.height

  const ctx = canvas.getContext('2d')
  if (!ctx) {
    throw new Error('无法创建 Canvas 上下文')
  }

  // 先填充背景色
  ctx.fillStyle = backgroundColor
  ctx.fillRect(0, 0, canvas.width, canvas.height)

  // 再绘制源图像（如果有透明区域，背景色会显示出来）
  ctx.drawImage(source, 0, 0)

  return canvas
}

/**
 * 将 Canvas 转换为 Blob
 */
export function canvasToBlob(
  canvas: HTMLCanvasElement,
  type: 'image/jpeg' | 'image/png' = 'image/jpeg',
  quality: number = 0.95
): Promise<Blob> {
  return new Promise((resolve, reject) => {
    canvas.toBlob(
      (blob) => {
        if (!blob) {
          reject(new Error('导出图像失败'))
          return
        }
        resolve(blob)
      },
      type,
      type === 'image/jpeg' ? quality : undefined
    )
  })
}

/**
 * 从 Blob 加载图像
 */
export async function loadImageFromBlob(blob: Blob): Promise<HTMLCanvasElement> {
  const bitmap = await createImageBitmap(blob)
  const canvas = document.createElement('canvas')
  canvas.width = bitmap.width
  canvas.height = bitmap.height

  const ctx = canvas.getContext('2d')
  if (!ctx) {
    bitmap.close()
    throw new Error('无法创建 Canvas 上下文')
  }

  ctx.drawImage(bitmap, 0, 0)
  bitmap.close()

  return canvas
}