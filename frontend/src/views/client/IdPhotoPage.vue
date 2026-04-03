<template>
  <div class="section-shell">
    <div class="app-frame id-photo-page">
      <header class="page-hero">
        <span class="section-kicker">Field Kit / Photo</span>
        <h1 class="section-title">证件照裁剪</h1>
        <p class="section-description">
          上传照片，自动移除背景，生成符合中国规格的证件照。
        </p>
        <div class="page-hero__actions">
          <router-link to="/tools" class="exit-btn">返回工具箱</router-link>
        </div>
      </header>

      <div class="main-panel">
        <!-- 左侧控制区 -->
        <div class="control-panel">
          <!-- 上传区域 -->
          <div class="control-section">
            <label class="control-label">上传照片</label>
            <div
              class="upload-zone"
              :class="{ 'has-file': file }"
              @click="triggerUpload"
              @dragover.prevent="isDragging = true"
              @dragleave="isDragging = false"
              @drop.prevent="handleDrop"
            >
              <input
                ref="fileInput"
                type="file"
                accept="image/jpeg,image/png"
                hidden
                @change="handleFileSelect"
              />
              <template v-if="!file">
                <span class="upload-icon">📷</span>
                <span class="upload-text">点击或拖拽上传</span>
                <span class="upload-hint">JPG/PNG，最大 20MB</span>
              </template>
              <template v-else>
                <img v-if="previewUrl" :src="previewUrl" class="upload-preview" alt="预览" />
                <span class="upload-filename">{{ file.name }}</span>
              </template>
            </div>
          </div>

          <!-- 模型加载状态 -->
          <div v-if="loadingModel" class="control-section loading-section">
            <div class="progress-bar">
              <div class="progress-fill" :style="{ width: `${loadProgress}%` }"></div>
            </div>
            <span class="loading-text">{{ loadStatus }}</span>
            <button class="btn btn-secondary" @click="cancelLoad">取消</button>
          </div>

          <!-- 规格选择 -->
          <div class="control-section">
            <label class="control-label">规格</label>
            <div class="spec-buttons">
              <button
                v-for="s in specs"
                :key="s.id"
                class="spec-btn"
                :class="{ active: spec === s.id }"
                @click="spec = s.id"
              >
                {{ s.name }}
              </button>
            </div>
            <span class="spec-hint">{{ currentSpec?.desc }}</span>
          </div>

          <!-- 背景色 - 仅在移除背景后显示 -->
          <div v-if="hasRemovedBg" class="control-section">
            <label class="control-label">背景色</label>
            <div class="color-buttons">
              <button
                v-for="c in colors"
                :key="c.id"
                class="color-btn"
                :class="{ active: bgColor === c.id }"
                :style="{ backgroundColor: c.hex }"
                @click="bgColor = c.id"
              >
                <span v-if="c.id === 'transparent'" class="transparent-pattern"></span>
              </button>
            </div>
          </div>

          <!-- 导出分辨率 -->
          <div class="control-section">
            <label class="control-label">导出分辨率</label>
            <div class="scale-buttons">
              <button
                v-for="opt in scaleOptions"
                :key="opt.id"
                class="scale-btn"
                :class="{ active: exportScale === opt.id }"
                @click="exportScale = opt.id"
              >
                {{ opt.name }}
              </button>
            </div>
            <span class="scale-hint">{{ scaleOptions.find(o => o.id === exportScale)?.desc }} · {{ (currentSpec?.width || 295) * exportScale }}×{{ (currentSpec?.height || 413) * exportScale }}px</span>
          </div>

          <!-- 图像调整 -->
          <div v-if="originalCanvas" class="control-section">
            <label class="control-label">图像调整</label>
            <div class="slider-row">
              <span>亮度</span>
              <input v-model.number="brightness" type="range" min="-50" max="50" />
              <span>{{ brightness }}</span>
            </div>
            <div class="slider-row">
              <span>对比度</span>
              <input v-model.number="contrast" type="range" min="-50" max="50" />
              <span>{{ contrast }}</span>
            </div>
          </div>

          <!-- 模型选择 -->
          <div v-if="file" class="control-section">
            <label class="control-label">AI 模型</label>
            <select v-model="model" class="model-select">
              <option value="isnet">ISNet (~80MB，标准)</option>
              <option value="isnet_fp16">ISNet FP16 (~40MB，快速)</option>
              <option value="isnet_quint8">ISNet Quantized (~20MB，最快)</option>
            </select>
            <span v-if="modelCached" class="model-cached">✓ 模型已缓存</span>
          </div>

          <!-- 操作按钮 -->
          <div class="action-buttons">
            <button
              class="btn btn-primary"
              :disabled="!canProcess"
              @click="process"
            >
              {{ processing ? '处理中...' : '移除背景' }}
            </button>
            <button
              class="btn"
              :disabled="!noBgCanvas"
              @click="reset"
            >
              还原
            </button>
            <button
              class="btn"
              :disabled="!resultCanvas"
              @click="download"
            >
              导出
            </button>
          </div>

          <!-- 状态消息 -->
          <p v-if="message" class="status-msg">{{ message }}</p>
          <p v-if="error" class="error-msg">{{ error }}</p>
        </div>

        <!-- 右侧预览区 -->
        <div class="preview-panel">
          <div class="preview-area" :style="previewFilterStyle">
            <!-- 背景色预览层 - 在 Cropper 下方 -->
            <div
              v-if="workingUrl && hasRemovedBg && bgColor !== 'transparent' && currentColor"
              class="bg-preview-layer"
              :style="{ backgroundColor: currentColor.hex }"
            ></div>
            <div v-else-if="workingUrl && hasRemovedBg" class="bg-preview-layer transparent-bg"></div>

            <Cropper
              v-if="workingUrl"
              ref="cropperRef"
              :src="workingUrl"
              :stencil-props="{
                aspectRatio: currentSpec ? currentSpec.width / currentSpec.height : 295 / 413,
                movable: true,
                resizable: false,
                lines: true,
                handlers: true
              }"
              :min-width="100"
              :min-height="100"
              image-restriction="none"
              default-boundaries="fill"
              class="cropper-instance"
              @change="onCropChange"
              @ready="onCropperReady"
            />
            <div v-if="!workingUrl" class="preview-placeholder">
              <span>请上传照片</span>
            </div>
          </div>
          <p v-if="resultCanvas" class="preview-info">
            {{ currentSpec?.name }} · {{ resultCanvas.width }}×{{ resultCanvas.height }}px · PNG
            <span v-if="hasRemovedBg"> · {{ currentColor?.name || '透明' }}</span>
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onUnmounted } from 'vue'
import { Cropper } from 'vue-advanced-cropper'
import 'vue-advanced-cropper/dist/style.css'
import type { ModelType } from '@/utils/idPhoto/backgroundRemoval'
import {
  ID_PHOTO_SPECS,
  BACKGROUND_COLORS,
  type PhotoSpecId,
  type BackgroundColorId
} from '@/utils/idPhoto/specConfig'

// 规格配置 (300 DPI 标准) - from specConfig
const specs = Object.values(ID_PHOTO_SPECS).map(s => ({
  id: s.id,
  name: s.name,
  width: s.widthPx,
  height: s.heightPx,
  desc: s.description
}))

const colors = Object.values(BACKGROUND_COLORS).map(c => ({
  id: c.id,
  name: c.name,
  hex: c.hex
}))

// 导出分辨率选项
const scaleOptions = [
  { id: 1, name: '标准 (1x)', desc: '网络分享' },
  { id: 3, name: '高清 (3x)', desc: '打印' },
  { id: 5, name: '超清 (5x)', desc: '高质量打印' }
]

// Refs
const fileInput = ref<HTMLInputElement | null>(null)
const cropperRef = ref<any>(null)

// 状态
const file = ref<File | null>(null)
const previewUrl = ref('')
const originalCanvas = ref<HTMLCanvasElement | null>(null)
const noBgCanvas = ref<HTMLCanvasElement | null>(null)
const resultCanvas = ref<HTMLCanvasElement | null>(null)
const hasRemovedBg = ref(false) // 追踪是否已移除背景
const workingUrl = ref('') // 用于Cropper的URL

const isDragging = ref(false)
const loadingModel = ref(false)
const loadProgress = ref(0)
const loadStatus = ref('')
const processing = ref(false)
const message = ref('')
const error = ref('')

// 选择
const spec = ref<PhotoSpecId>('one-inch')
const bgColor = ref<BackgroundColorId>('transparent')
const model = ref<ModelType>('isnet')
const brightness = ref(0)
const contrast = ref(0)
const exportScale = ref(5) // 默认5x超清

// 裁剪控制
const cropZoom = ref(1)
const cropX = ref(0)
const cropY = ref(0)

// 计算属性：获取当前工作的源图像
const workingCanvas = computed(() => noBgCanvas.value || originalCanvas.value)

const currentSpec = computed(() => specs.find(s => s.id === spec.value))
const currentColor = computed(() => colors.find(c => c.id === bgColor.value))
const canProcess = computed(() => file.value && !processing.value && !loadingModel.value)
const modelCached = ref(false)

// 预览滤镜样式 - 实时预览亮度/对比度
const previewFilterStyle = computed(() => {
  if (!workingCanvas.value) return {}
  const br = 1 + brightness.value / 100
  const ct = 1 + contrast.value / 100
  return {
    filter: `brightness(${br}) contrast(${ct})`
  }
})

// 状态
let loadCancelled = false

// 方法
function triggerUpload() {
  fileInput.value?.click()
}

function clearMsg() {
  message.value = ''
  error.value = ''
}

async function handleFileSelect(e: Event) {
  const input = e.target as HTMLInputElement
  const f = input.files?.[0]
  if (f) await loadFile(f)
}

async function handleDrop(e: DragEvent) {
  isDragging.value = false
  const f = e.dataTransfer?.files?.[0]
  if (f && f.type.match(/^image\/(jpeg|png)$/)) {
    if (f.size > 20 * 1024 * 1024) {
      error.value = '文件过大，最大 20MB'
      return
    }
    await loadFile(f)
  }
}

async function loadFile(f: File) {
  clearMsg()
  file.value = f

  // 创建预览 URL
  if (previewUrl.value) URL.revokeObjectURL(previewUrl.value)
  previewUrl.value = URL.createObjectURL(f)

  // 加载到 Canvas
  const img = new Image()
  img.onload = async () => {
    const canvas = document.createElement('canvas')
    canvas.width = img.naturalWidth
    canvas.height = img.naturalHeight
    const ctx = canvas.getContext('2d')
    ctx?.drawImage(img, 0, 0)
    originalCanvas.value = canvas

    // 创建cropper用的URL
    if (workingUrl.value) URL.revokeObjectURL(workingUrl.value)
    workingUrl.value = URL.createObjectURL(await new Promise<Blob>((resolve) => {
      canvas.toBlob((b) => resolve(b!), 'image/png')
    }))

    // 上传后立即生成预览结果
    updateResult()
  }
  img.src = previewUrl.value

  // 重置
  noBgCanvas.value = null
  hasRemovedBg.value = false
  resultCanvas.value = null
  cropZoom.value = 1
  cropX.value = 0
  cropY.value = 0

  message.value = '照片已上传，可直接裁剪或移除背景'
}

function cancelLoad() {
  loadCancelled = true
  loadingModel.value = false
  loadProgress.value = 0
  loadStatus.value = ''
  message.value = '已取消，点击"移除背景"重新处理'
}

async function process() {
  if (!originalCanvas.value) return

  clearMsg()
  processing.value = true
  loadingModel.value = true
  loadProgress.value = 0
  loadStatus.value = '正在初始化...'
  loadCancelled = false
  message.value = '正在处理...'

  try {
    // 使用新的背景移除库
    const { removeBackground: removeBg, isModelCached } = await import('@/utils/idPhoto/backgroundRemoval')

    // 检查模型是否已缓存
    modelCached.value = isModelCached(model.value as ModelType)

    const canvas = await removeBg(
      originalCanvas.value,
      model.value as ModelType,
      (progress, status) => {
        loadProgress.value = progress
        loadStatus.value = status
      }
    )

    if (loadCancelled) {
      throw new Error('已取消')
    }

    noBgCanvas.value = canvas
    hasRemovedBg.value = true
    modelCached.value = true

    // 调试：检查背景移除结果是否有透明像素
    const debugCtx = canvas.getContext('2d')
    if (debugCtx) {
      const debugData = debugCtx.getImageData(0, 0, Math.min(canvas.width, 100), Math.min(canvas.height, 100))
      const transparentPixels = debugData.data.filter((v, i) => i % 4 === 3 && v < 255).length
      console.log('Background removal result:', transparentPixels > 0 ? `has ${transparentPixels} transparent pixels` : 'fully opaque (no transparency)')
    }

    // 更新cropper用的URL（使用去背后的图像）
    if (workingUrl.value) URL.revokeObjectURL(workingUrl.value)
    workingUrl.value = URL.createObjectURL(await new Promise<Blob>((resolve) => {
      canvas.toBlob((b) => resolve(b!), 'image/png')
    }))

    // 重新生成结果（使用去背后的图像）
    updateResult()

    message.value = '背景已移除'

  } catch (e: any) {
    if (e.message !== '已取消') {
      error.value = `处理失败: ${e.message}`
    }
  } finally {
    processing.value = false
    loadingModel.value = false
  }
}

function updateResult() {
  // 选择源图像：优先使用去背后的，否则用原图
  const source = noBgCanvas.value || originalCanvas.value
  if (!source || !currentSpec.value) return

  const s = currentSpec.value
  const scale = exportScale.value

  // 创建结果 Canvas (使用scale倍分辨率)
  const canvas = document.createElement('canvas')
  canvas.width = s.width * scale
  canvas.height = s.height * scale
  const ctx = canvas.getContext('2d', { alpha: true })!

  // 关键：确保画布完全透明
  ctx.clearRect(0, 0, canvas.width, canvas.height)

  // 仅在已移除背景且选择了非透明背景色时才填充背景
  // 透明背景 = 不填充任何背景，保持 canvas 透明
  const shouldFillBackground = hasRemovedBg.value
    && bgColor.value !== 'transparent'
    && currentColor.value

  if (shouldFillBackground) {
    ctx.fillStyle = currentColor.value!.hex
    ctx.fillRect(0, 0, canvas.width, canvas.height)
  }

  // 计算目标比例
  const targetRatio = s.width / s.height
  const sourceRatio = source.width / source.height

  // 基础裁剪区域（保持目标比例）
  let baseCropW: number, baseCropH: number
  if (sourceRatio > targetRatio) {
    baseCropH = source.height
    baseCropW = baseCropH * targetRatio
  } else {
    baseCropW = source.width
    baseCropH = baseCropW / targetRatio
  }

  // 应用缩放（zoom 放大意味着裁剪区域变小）
  const scaledCropW = baseCropW / cropZoom.value
  const scaledCropH = baseCropH / cropZoom.value

  // 计算居中位置
  let cropCenterX = source.width / 2
  let cropCenterY = source.height / 2

  // 应用位置偏移（百分比转换为像素）
  cropCenterX += (cropX.value / 100) * (source.width / 2)
  cropCenterY += (cropY.value / 100) * (source.height / 2)

  // 计算最终裁剪区域
  let cropX1 = cropCenterX - scaledCropW / 2
  let cropY1 = cropCenterY - scaledCropH / 2

  // 边界检查
  cropX1 = Math.max(0, Math.min(cropX1, source.width - scaledCropW))
  cropY1 = Math.max(0, Math.min(cropY1, source.height - scaledCropH))

  // 透明背景时禁用滤镜，避免 ctx.filter 破坏 alpha 通道（透明像素被渲染为白色）
  const br = 1 + brightness.value / 100
  const ct = 1 + contrast.value / 100

  // 调试：检查源 canvas 的透明度
  if (hasRemovedBg.value && bgColor.value === 'transparent') {
    const srcCtx = (source as HTMLCanvasElement).getContext('2d')
    if (srcCtx) {
      const srcData = srcCtx.getImageData(0, 0, Math.min(source.width, 50), Math.min(source.height, 50))
      let srcTransparent = 0
      let srcOpaque = 0
      for (let i = 3; i < srcData.data.length; i += 4) {
        if (srcData.data[i] < 255) srcTransparent++
        else srcOpaque++
      }
      console.log(`[updateResult] Source canvas: ${srcTransparent} transparent, ${srcOpaque} opaque`)
    }
  }

  if (!(hasRemovedBg.value && bgColor.value === 'transparent')) {
    ctx.filter = `brightness(${br}) contrast(${ct})`
  }

  // 使用高质量缩放
  ctx.imageSmoothingEnabled = true
  ctx.imageSmoothingQuality = 'high'

  // 缩放context以绘制高分辨率图像
  ctx.scale(scale, scale)

  // 从裁剪区域绘制到目标尺寸
  ctx.drawImage(
    source,
    cropX1, cropY1, scaledCropW, scaledCropH,  // 源裁剪区域
    0, 0, s.width, s.height      // 目标尺寸（scale前的逻辑尺寸）
  )

  // 调试：检查结果 canvas 的透明度
  if (hasRemovedBg.value && bgColor.value === 'transparent') {
    const resultData = ctx.getImageData(0, 0, Math.min(canvas.width, 50), Math.min(canvas.height, 50))
    let resTransparent = 0
    let resOpaque = 0
    for (let i = 3; i < resultData.data.length; i += 4) {
      if (resultData.data[i] < 255) resTransparent++
      else resOpaque++
    }
    console.log(`[updateResult] Result canvas: ${resTransparent} transparent, ${resOpaque} opaque`)
  }

  resultCanvas.value = canvas
}

// 监听调整变化
watch([brightness, contrast, bgColor, cropZoom, cropX, cropY, exportScale], () => {
  if (workingCanvas.value) updateResult()
})

// 监听背景移除状态变化，更新背景色
watch(bgColor, () => {
  if (hasRemovedBg.value && workingCanvas.value) updateResult()
})

// 处理cropper初始化完成
function onCropperReady() {
  // 初始化后设置合适的显示
  if (cropperRef.value) {
    // 确保cropper可以正常操作
  }
}

// 处理cropper变化
function onCropChange({ coordinates }: any) {
  const source = workingCanvas.value
  if (!source || !coordinates) return

  // 计算相对于源图像的crop位置
  const sourceWidth = source.width
  const sourceHeight = source.height

  // cropper返回的是相对坐标，需要转换为我们的控制参数
  // center offset as percentage
  const centerX = coordinates.left + coordinates.width / 2
  const centerY = coordinates.top + coordinates.height / 2
  const centerOffsetX = (centerX - sourceWidth / 2) / (sourceWidth / 2) * 100
  const centerOffsetY = (centerY - sourceHeight / 2) / (sourceHeight / 2) * 100

  cropX.value = Math.round(centerOffsetX)
  cropY.value = Math.round(centerOffsetY)

  // zoom based on how much of the image is visible
  const targetRatio = currentSpec.value ? currentSpec.value.width / currentSpec.value.height : 295 / 413
  let idealCropW: number, idealCropH: number
  if (sourceWidth / sourceHeight > targetRatio) {
    idealCropH = sourceHeight
    idealCropW = idealCropH * targetRatio
  } else {
    idealCropW = sourceWidth
    idealCropH = idealCropW / targetRatio
  }
  const zoom = idealCropW / coordinates.width
  cropZoom.value = Math.round(zoom * 10) / 10

  updateResult()
}

async function reset() {
  noBgCanvas.value = null
  hasRemovedBg.value = false
  brightness.value = 0
  contrast.value = 0
  bgColor.value = 'transparent'
  cropZoom.value = 1
  cropX.value = 0
  cropY.value = 0

  // 恢复cropper URL为原图
  if (originalCanvas.value) {
    if (workingUrl.value) URL.revokeObjectURL(workingUrl.value)
    workingUrl.value = URL.createObjectURL(await new Promise<Blob>((resolve) => {
      originalCanvas.value!.toBlob((b) => resolve(b!), 'image/png')
    }))
  }

  // Reset cropper position if available
  if (cropperRef.value) {
    cropperRef.value.reset()
  }
  // 重新生成结果（使用原图）
  updateResult()
  message.value = '已还原到原图'
}

async function download() {
  if (!resultCanvas.value || !file.value) return

  // 调试：检查 canvas 透明像素比例
  const ctx = resultCanvas.value.getContext('2d')
  if (ctx) {
    const w = resultCanvas.value.width
    const h = resultCanvas.value.height
    const imageData = ctx.getImageData(0, 0, w, h)
    const transparentCount = imageData.data.filter((v, i) => i % 4 === 3 && v < 255).length
    const totalPixels = imageData.data.length / 4
    console.log(`Transparency: ${transparentCount}/${totalPixels} pixels (${((transparentCount / totalPixels) * 100).toFixed(1)}%)`)
    console.log('Current bgColor:', bgColor.value)
    console.log('hasRemovedBg:', hasRemovedBg.value)
  }

  const blob = await new Promise<Blob>((resolve) => {
    resultCanvas.value!.toBlob((b) => resolve(b!), 'image/png')
  })

  const name = file.value.name.replace(/\.[^.]+$/, '')
  const specName = currentSpec.value?.name || ''
  const scaleName = exportScale.value === 5 ? '5x' : exportScale.value === 3 ? '3x' : '1x'

  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `${name}-${specName}-${scaleName}.png`
  a.click()
  URL.revokeObjectURL(url)

  message.value = '已导出 PNG'
}

onUnmounted(() => {
  if (previewUrl.value) URL.revokeObjectURL(previewUrl.value)
  if (workingUrl.value) URL.revokeObjectURL(workingUrl.value)
})
</script>

<style scoped>
.id-photo-page {
  height: calc(100vh - 60px);
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  overflow: hidden;
}

.page-hero {
  display: flex;
  flex-direction: column;
  gap: 0.2rem;
  flex: 0 0 auto;
}

.page-hero__actions {
  padding-top: 0.2rem;
}

.exit-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-height: 2.5rem;
  padding: 0.48rem 1rem;
  border-radius: var(--radius-pill);
  border: 1px solid var(--border-primary);
  color: var(--text-primary);
  background: var(--bg-panel);
  transition: transform var(--duration-fast) var(--ease-standard), border-color var(--duration-fast) var(--ease-standard);
}

.exit-btn:hover {
  transform: translateY(-1px);
  border-color: var(--border-strong);
}

.main-panel {
  display: grid;
  grid-template-columns: 260px 1fr;
  gap: 0.75rem;
  padding: 0.75rem;
  background: var(--bg-panel);
  border-radius: var(--radius-lg);
  flex: 1;
  min-height: 0;
  overflow: hidden;
}

.control-panel {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  overflow-y: auto;
  max-height: 100%;
  padding-right: 0.25rem;
}

.control-section {
  display: flex;
  flex-direction: column;
  gap: 0.2rem;
}

.control-label {
  font-size: 0.65rem;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.upload-zone {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 60px;
  border: 2px dashed var(--border-primary);
  border-radius: var(--radius-md);
  cursor: pointer;
  padding: 0.4rem;
  text-align: center;
}

.upload-zone:hover {
  border-color: var(--text-muted);
}

.upload-zone.has-file {
  border-style: solid;
  min-height: auto;
  padding: 0.3rem;
}

.upload-icon {
  font-size: 1.5rem;
}

.upload-text {
  color: var(--text-secondary);
  font-size: 0.8rem;
}

.upload-hint {
  font-size: 0.7rem;
  color: var(--text-muted);
}

.upload-preview {
  max-width: 100%;
  max-height: 40px;
  border-radius: var(--radius-sm);
}

.upload-filename {
  font-size: 0.75rem;
  color: var(--text-secondary);
  word-break: break-all;
}

.loading-section {
  padding: 0.5rem;
  background: var(--accent-soft);
  border-radius: var(--radius-md);
}

.progress-bar {
  height: 4px;
  background: var(--bg-panel-strong);
  border-radius: 2px;
  overflow: hidden;
  margin-bottom: 0.3rem;
}

.progress-fill {
  height: 100%;
  background: var(--accent-primary);
  transition: width 0.3s;
}

.loading-text {
  font-size: 0.75rem;
  color: var(--text-secondary);
  display: block;
  margin-bottom: 0.3rem;
}

.spec-buttons {
  display: flex;
  gap: 0.3rem;
}

.spec-btn {
  flex: 1;
  padding: 0.3rem 0.4rem;
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-pill);
  background: var(--bg-panel-strong);
  color: var(--text-secondary);
  cursor: pointer;
  font-size: 0.8rem;
}

.spec-btn.active {
  background: var(--accent-soft);
  border-color: var(--border-strong);
  color: var(--text-primary);
}

.spec-hint {
  font-size: 0.65rem;
  color: var(--text-muted);
}

.scale-buttons {
  display: flex;
  gap: 0.3rem;
}

.scale-btn {
  flex: 1;
  padding: 0.3rem 0.4rem;
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-pill);
  background: var(--bg-panel-strong);
  color: var(--text-secondary);
  cursor: pointer;
  font-size: 0.8rem;
}

.scale-btn.active {
  background: var(--accent-soft);
  border-color: var(--border-strong);
  color: var(--text-primary);
}

.scale-hint {
  font-size: 0.65rem;
  color: var(--text-muted);
}

.color-buttons {
  display: flex;
  gap: 0.4rem;
}

.color-btn {
  width: 1.75rem;
  height: 1.75rem;
  border-radius: 50%;
  border: 2px solid var(--border-primary);
  cursor: pointer;
}

.color-btn.active {
  border-width: 3px;
  border-color: var(--text-primary);
}

.slider-row {
  display: flex;
  align-items: center;
  gap: 0.3rem;
  font-size: 0.75rem;
}

.slider-row input {
  flex: 1;
  height: 4px;
}

.model-select {
  width: 100%;
  padding: 0.3rem;
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-md);
  background: var(--bg-panel-strong);
  color: var(--text-primary);
  font-size: 0.8rem;
}

.model-cached {
  font-size: 0.65rem;
  color: var(--color-success, #22c55e);
}

.action-buttons {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 0.3rem;
}

.btn {
  padding: 0.4rem;
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-pill);
  background: var(--bg-panel-strong);
  color: var(--text-primary);
  cursor: pointer;
  font-size: 0.85rem;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-primary {
  background: var(--accent-soft);
  border-color: var(--border-strong);
}

.status-msg {
  font-size: 0.75rem;
  color: var(--text-secondary);
}

.error-msg {
  font-size: 0.75rem;
  color: var(--color-error);
}

/* 预览区 */
.preview-panel {
  display: flex;
  flex-direction: column;
  min-width: 0;
  flex: 1;
  min-height: 0;
}

.preview-area {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  flex: 1;
  min-height: 200px;
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-md);
  background: var(--bg-panel-strong);
  overflow: hidden;
}

.preview-placeholder {
  color: var(--text-muted);
  padding: 2rem;
  position: absolute;
  z-index: 1;
}

/* 背景预览层 - 放在 Cropper 下方 */
.bg-preview-layer {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 0;
  pointer-events: none;
}

/* 透明背景的棋盘格图案 */
.transparent-bg {
  background-image:
    linear-gradient(45deg, #ccc 25%, transparent 25%),
    linear-gradient(-45deg, #ccc 25%, transparent 25%),
    linear-gradient(45deg, transparent 75%, #ccc 75%),
    linear-gradient(-45deg, transparent 75%, #ccc 75%);
  background-size: 16px 16px;
  background-position: 0 0, 0 8px, 8px -8px, -8px 0;
  background-color: #fff;
}

.cropper-instance {
  width: 100%;
  height: 100%;
  position: relative;
  z-index: 1;
}

.cropper-instance :deep(.vue-advanced-cropper) {
  width: 100%;
  height: 100%;
}

.cropper-instance :deep(.vue-advanced-cropper__background) {
  background: transparent !important;
}

/* 移除玻璃遮罩层 */
.cropper-instance :deep(.vue-advanced-cropper__foreground) {
  background: transparent !important;
}

/* 裁剪框样式 - 清晰边框 + 中点线 */
.cropper-instance :deep(.vue-rectangle-stencil) {
  border: 1.5px solid rgba(255, 255, 255, 0.8) !important;
  box-shadow: 0 0 0 9999px rgba(0, 0, 0, 0.5);
}

/* 网格线 */
.cropper-instance :deep(.vue-simple-line) {
  border-color: rgba(255, 255, 255, 0.3) !important;
}

/* 四角拖拽点 - 所有可能的 class 名 */
.cropper-instance :deep(.vue-handler-wrapper) {
  width: 20px !important;
  height: 20px !important;
}

.cropper-instance :deep(.vue-simple-handler) {
  background: #fff !important;
  border: none !important;
  width: 10px !important;
  height: 10px !important;
  border-radius: 50% !important;
  box-shadow: 0 0 0 2px rgba(0,0,0,0.3);
}

/* 四角 L 形装饰 */
.cropper-instance :deep(.vue-rectangle-stencil)::before,
.cropper-instance :deep(.vue-rectangle-stencil)::after {
  display: none;
}

/* 透明背景按钮图案 */
.color-btn {
  position: relative;
  overflow: hidden;
}

.transparent-pattern {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-image:
    linear-gradient(45deg, #ccc 25%, transparent 25%),
    linear-gradient(-45deg, #ccc 25%, transparent 25%),
    linear-gradient(45deg, transparent 75%, #ccc 75%),
    linear-gradient(-45deg, transparent 75%, #ccc 75%);
  background-size: 8px 8px;
  background-position: 0 0, 0 4px, 4px -4px, -4px 0;
}

.preview-info {
  font-size: 0.7rem;
  color: var(--text-muted);
  text-align: center;
  padding: 0.2rem 0;
}

/* 响应式 */
@media (max-width: 900px) {
  .id-photo-page {
    height: auto;
    overflow: visible;
  }

  .main-panel {
    grid-template-columns: 1fr;
    flex: none;
  }

  .preview-area {
    min-height: 220px;
    flex: none;
  }
}
</style>