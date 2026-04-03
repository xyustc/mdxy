/**
 * 背景移除工具
 * 使用 @imgly/background-removal 进行背景移除
 */

import { removeBackground as removeBg } from '@imgly/background-removal'

export type ModelType = 'isnet' | 'isnet_fp16' | 'isnet_quint8'

interface ModelConfig {
  name: string
  description: string
  estimatedSize: string
}

const MODEL_CONFIGS: Record<ModelType, ModelConfig> = {
  isnet: {
    name: 'ISNet',
    description: '标准模型，平衡速度和质量',
    estimatedSize: '~80MB'
  },
  isnet_fp16: {
    name: 'ISNet FP16',
    description: '半精度模型，更快但略低质量',
    estimatedSize: '~40MB'
  },
  isnet_quint8: {
    name: 'ISNet Quantized',
    description: '量化模型，最小最快',
    estimatedSize: '~20MB'
  }
}

type ProgressCallback = (progress: number, status: string) => void

// 模型缓存状态
const modelCache: Record<ModelType, boolean> = {
  isnet: false,
  isnet_fp16: false,
  isnet_quint8: false
}

/**
 * 检查模型是否已缓存
 */
export function isModelCached(model: ModelType): boolean {
  return modelCache[model]
}

/**
 * 初始化背景移除模型（此库不需要初始化）
 */
export async function initBackgroundRemoval(
  _model: ModelType = 'isnet',
  onProgress?: ProgressCallback
): Promise<void> {
  onProgress?.(100, '模型已就绪')
}

/**
 * 移除背景
 * 返回带有透明背景的 Canvas
 */
export async function removeBackground(
  image: HTMLImageElement | HTMLCanvasElement | ImageBitmap | string,
  model: ModelType = 'isnet',
  onProgress?: ProgressCallback
): Promise<HTMLCanvasElement> {
  // 将输入转换为 Blob
  let blob: Blob

  if (typeof image === 'string') {
    // URL 字符串
    const response = await fetch(image)
    blob = await response.blob()
  } else if (image instanceof HTMLCanvasElement) {
    blob = await new Promise<Blob>((resolve) => {
      image.toBlob((b) => resolve(b!), 'image/png')
    })
  } else if (image instanceof HTMLImageElement) {
    // 先绘制到 Canvas
    const canvas = document.createElement('canvas')
    canvas.width = image.naturalWidth
    canvas.height = image.naturalHeight
    const ctx = canvas.getContext('2d')
    if (!ctx) throw new Error('无法创建 Canvas 上下文')
    ctx.drawImage(image, 0, 0)
    blob = await new Promise<Blob>((resolve) => {
      canvas.toBlob((b) => resolve(b!), 'image/png')
    })
  } else {
    // ImageBitmap
    const canvas = document.createElement('canvas')
    canvas.width = image.width
    canvas.height = image.height
    const ctx = canvas.getContext('2d')
    if (!ctx) throw new Error('无法创建 Canvas 上下文')
    ctx.drawImage(image, 0, 0)
    blob = await new Promise<Blob>((resolve) => {
      canvas.toBlob((b) => resolve(b!), 'image/png')
    })
  }

  // 调用背景移除库
  const isCached = modelCache[model]
  onProgress?.(5, isCached ? '处理中...' : '准备处理...')

  // 进度状态
  let lastProgress = 5
  let lastDownloadKey = ''
  const startTime = Date.now()

  const resultBlob = await removeBg(blob, {
    model: model,
    output: {
      format: 'image/png',
      quality: 1
    },
    progress: (key: string, current: number, total: number) => {
      if (isCached) {
        // 模型已缓存，直接显示处理进度
        if (key === 'compute:inference') {
          // 使用时间估算进度（处理过程通常没有精确进度）
          const elapsed = Date.now() - startTime
          const estimated = 15000 // 预估 15 秒
          const progress = Math.min(5 + (elapsed / estimated) * 90, 95)
          if (progress > lastProgress) {
            lastProgress = progress
            onProgress?.(progress, `处理图像中... ${Math.round(progress)}%`)
          }
        }
        return
      }

      // 模型下载中
      if (key === 'fetch:model') {
        if (key !== lastDownloadKey && lastDownloadKey) {
          // 切换到新的下载块
          lastProgress += Math.floor(100 * (1 / 6))
        }
        lastDownloadKey = key

        if (total > 0 && current > 0) {
          const percent = (current / total)
          const downloadedMB = Math.round(current / 1024 / 1024)
          const totalMB = Math.round(total / 1024 / 1024)
          // 下载进度占 70%
          const progress = 5 + percent * 70
          lastProgress = progress
          onProgress?.(progress, `下载模型 ${downloadedMB}/${totalMB}MB`)
        }
      } else if (key === 'compute:inference') {
        // 标记模型已缓存
        modelCache[model] = true

        // 处理进度：75% - 95%
        const elapsed = Date.now() - startTime
        const estimated = 15000 // 预估 15 秒
        const progress = Math.min(75 + (elapsed / estimated) * 20, 95)
        if (progress > lastProgress) {
          lastProgress = progress
          onProgress?.(progress, `处理图像中...`)
        }
      }
    }
  })

  // 标记模型已缓存
  modelCache[model] = true
  onProgress?.(95, '生成结果图像...')

  // 将 Blob 转换为 Canvas（使用 createImageBitmap 确保 alpha 通道保留）
  const bitmap = await createImageBitmap(resultBlob as Blob)

  const canvas = document.createElement('canvas')
  canvas.width = bitmap.width
  canvas.height = bitmap.height
  const ctx = canvas.getContext('2d')
  if (!ctx) throw new Error('无法创建 Canvas 上下文')

  // 确保画布透明
  ctx.clearRect(0, 0, canvas.width, canvas.height)
  ctx.drawImage(bitmap, 0, 0)
  bitmap.close()

  // 调试：检查 canvas 的透明度
  const checkImageData = ctx.getImageData(0, 0, Math.min(canvas.width, 50), Math.min(canvas.height, 50))
  let transparentCount = 0
  let opaqueCount = 0
  for (let i = 3; i < checkImageData.data.length; i += 4) {
    if (checkImageData.data[i] < 255) transparentCount++
    else opaqueCount++
  }
  console.log(`backgroundRemoval canvas check: ${transparentCount} semi-transparent, ${opaqueCount} fully opaque (out of ${transparentCount + opaqueCount} sampled)`)

  onProgress?.(100, '处理完成')

  return canvas
}

/**
 * 获取模型信息
 */
export function getModelInfo(model: ModelType): ModelConfig {
  return MODEL_CONFIGS[model]
}

/**
 * 获取所有可用模型
 */
export function getAvailableModels(): ModelType[] {
  return ['isnet', 'isnet_fp16', 'isnet_quint8']
}