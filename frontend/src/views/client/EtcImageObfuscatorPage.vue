<template>
  <div class="section-shell">
    <div class="app-frame etc-page">
      <header class="etc-hero">
        <span class="section-kicker">Field Kit / EtC</span>
        <h1 class="section-title">EtC 图片混淆</h1>
        <p class="section-description">
          基于分块置乱与块内变换的 JPEG 友好混淆方案。使用相同参数与密钥可解混淆回接近原图，但 JPEG 为有损格式，不保证像素级一致。
        </p>
      </header>

      <section class="surface-panel etc-panel">
        <div class="etc-panel__main">
          <section class="etc-controls">
            <div class="control-group">
              <label class="control-label" for="etc-file">输入图片</label>
              <input
                id="etc-file"
                ref="fileInputRef"
                class="visually-hidden-input"
                type="file"
                accept="image/*"
                @change="onFileChange"
              />
              <div class="file-picker">
                <button class="file-picker__button" type="button" @click="openFilePicker">选择图片</button>
                <span class="file-picker__name" :class="{ 'is-empty': !selectedFileName }">
                  {{ selectedFileName || '未选择文件' }}
                </span>
              </div>
            </div>

            <div class="control-group">
              <label class="control-label" for="etc-key">密钥</label>
              <input
                id="etc-key"
                v-model.trim="form.key"
                class="control-input"
                type="password"
                placeholder="至少 4 位，解混淆需一致"
              />
            </div>

            <div class="control-grid">
              <div class="control-group">
                <label class="control-label" for="etc-block-size">块大小</label>
                <select id="etc-block-size" v-model.number="form.blockSize" class="control-input">
                  <option :value="8">8 x 8</option>
                  <option :value="16">16 x 16</option>
                </select>
              </div>

              <div class="control-group">
                <label class="control-label" for="etc-quality">JPEG 质量 {{ form.quality.toFixed(2) }}</label>
                <input
                  id="etc-quality"
                  v-model.number="form.quality"
                  class="control-range"
                  type="range"
                  min="0.7"
                  max="0.98"
                  step="0.01"
                />
              </div>
            </div>

            <div class="control-checks">
              <label><input v-model="form.enableRotateFlip" type="checkbox" /> 块内旋转/翻转</label>
              <label><input v-model="form.enableChannelPermutation" type="checkbox" /> RGB 通道置换</label>
              <label><input v-model="form.enableNegativeTransform" type="checkbox" /> 负片变换</label>
              <label><input v-model="form.usePngOnDecrypt" type="checkbox" /> 解混淆结果优先导出 PNG（推荐）</label>
            </div>

            <div class="control-group">
              <label class="control-label" for="etc-config">参数配置（可导入/导出）</label>
              <textarea
                id="etc-config"
                v-model="configPayload"
                class="control-input control-textarea"
                placeholder="点击“导出参数”生成配置 JSON；或粘贴配置后点击“导入参数”"
              ></textarea>
            </div>

            <div class="etc-actions">
              <button class="action-btn action-btn--primary" :disabled="!canRun || isBusy" @click="runTransform('encrypt')">
                {{ isBusy ? '处理中...' : '混淆' }}
              </button>
              <button class="action-btn" :disabled="!canRun || isBusy" @click="runTransform('decrypt')">解混淆</button>
              <button class="action-btn" :disabled="!selectedFile || isBusy" @click="restoreOriginal">还原</button>
              <button class="action-btn" :disabled="!workingImageData || isBusy" @click="downloadOutput">下载</button>
            </div>

            <div class="etc-actions etc-actions--secondary">
              <button class="action-btn action-btn--secondary" :disabled="isBusy" @click="applyJpegSafePreset">JPEG 友好预设</button>
              <button class="action-btn action-btn--secondary" :disabled="isBusy" @click="exportConfig">导出参数</button>
              <button class="action-btn action-btn--secondary" :disabled="isBusy" @click="importConfig">导入参数</button>
            </div>

            <p v-if="message" class="status-text">{{ message }}</p>
            <p v-if="errorMessage" class="status-text status-text--error">{{ errorMessage }}</p>
          </section>

          <section class="etc-preview">
            <article class="preview-card">
              <h3>原图</h3>
              <img v-if="originalUrl" :src="originalUrl" alt="原图预览" />
              <p v-else>请选择图片文件。</p>
            </article>

            <article class="preview-card">
              <h3>当前图像</h3>
              <img v-if="workingUrl" :src="workingUrl" alt="当前图像预览" />
              <p v-else>混淆或解混淆后会在这里展示结果。</p>
            </article>
          </section>
        </div>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, reactive, ref } from 'vue'
import {
  ETC_DEFAULT_CONFIG,
  decryptImageData,
  encryptImageData,
  normalizeEtcConfig,
  parseEtcConfig,
  serializeEtcConfig,
  type EtcBlockSize
} from '@/utils/etcImageCipher'

type TransformMode = 'encrypt' | 'decrypt'

const selectedFile = ref<File | null>(null)
const fileInputRef = ref<HTMLInputElement | null>(null)
const originalImageData = ref<ImageData | null>(null)
const workingImageData = ref<ImageData | null>(null)

const originalUrl = ref('')
const workingUrl = ref('')

const configPayload = ref('')
const message = ref('')
const errorMessage = ref('')
const isBusy = ref(false)
const lastMode = ref<TransformMode>('encrypt')

const form = reactive({
  key: '',
  blockSize: ETC_DEFAULT_CONFIG.blockSize as EtcBlockSize,
  quality: ETC_DEFAULT_CONFIG.quality,
  enableRotateFlip: ETC_DEFAULT_CONFIG.enableRotateFlip,
  enableChannelPermutation: ETC_DEFAULT_CONFIG.enableChannelPermutation,
  enableNegativeTransform: ETC_DEFAULT_CONFIG.enableNegativeTransform,
  usePngOnDecrypt: true
})

const canRun = computed(() => Boolean(workingImageData.value && form.key.trim().length >= 4))
const selectedFileName = computed(() => selectedFile.value?.name || '')

function clearMessages() {
  message.value = ''
  errorMessage.value = ''
}

function openFilePicker() {
  fileInputRef.value?.click()
}

function revokeUrl(url: string) {
  if (url) {
    URL.revokeObjectURL(url)
  }
}

function setOriginalUrl(file: File) {
  revokeUrl(originalUrl.value)
  originalUrl.value = URL.createObjectURL(file)
}

async function setWorkingPreviewFromImageData(data: ImageData) {
  const canvas = imageDataToCanvas(data)
  const blob = await canvasToBlob(canvas, 'image/png', 1)
  revokeUrl(workingUrl.value)
  workingUrl.value = URL.createObjectURL(blob)
}

function cloneImageData(data: ImageData): ImageData {
  return new ImageData(new Uint8ClampedArray(data.data), data.width, data.height)
}

function imageDataToCanvas(data: ImageData): HTMLCanvasElement {
  const canvas = document.createElement('canvas')
  canvas.width = data.width
  canvas.height = data.height
  const ctx = canvas.getContext('2d')
  if (!ctx) {
    throw new Error('无法初始化图像处理上下文')
  }
  ctx.putImageData(data, 0, 0)
  return canvas
}

async function loadImageDataFromBlob(blob: Blob): Promise<ImageData> {
  const bitmap = await createImageBitmap(blob)
  const canvas = document.createElement('canvas')
  canvas.width = bitmap.width
  canvas.height = bitmap.height
  const ctx = canvas.getContext('2d')
  if (!ctx) {
    bitmap.close()
    throw new Error('无法初始化图像处理上下文')
  }
  ctx.drawImage(bitmap, 0, 0)
  bitmap.close()
  return ctx.getImageData(0, 0, canvas.width, canvas.height)
}

async function onFileChange(event: Event) {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return

  isBusy.value = true
  clearMessages()
  try {
    selectedFile.value = file
    setOriginalUrl(file)
    const loaded = await loadImageDataFromBlob(file)
    originalImageData.value = cloneImageData(loaded)
    workingImageData.value = cloneImageData(loaded)
    await setWorkingPreviewFromImageData(workingImageData.value)
    message.value = '图片已载入。工具内部现为无损链路，连续 3-5 次混淆/解混淆不会再叠加编码损失。'
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : '载入图片失败'
  } finally {
    isBusy.value = false
  }
}

async function runTransform(mode: TransformMode) {
  if (!canRun.value || !workingImageData.value) return

  clearMessages()
  isBusy.value = true
  lastMode.value = mode

  try {
    const options = normalizeEtcConfig({
      key: form.key,
      blockSize: form.blockSize,
      quality: form.quality,
      enableRotateFlip: form.enableRotateFlip,
      enableChannelPermutation: form.enableChannelPermutation,
      enableNegativeTransform: form.enableNegativeTransform
    })

    const source = cloneImageData(workingImageData.value)
    const transformed = mode === 'encrypt' ? encryptImageData(source, options) : decryptImageData(source, options)
    workingImageData.value = transformed
    await setWorkingPreviewFromImageData(transformed)

    if (mode === 'encrypt') {
      message.value = '混淆完成。当前操作链路未做 JPEG 重编码，可连续多次处理。'
    } else {
      message.value = form.usePngOnDecrypt
        ? '解混淆完成。建议下载 PNG 以保持当前质量。'
        : '解混淆完成。若下载 JPEG，多次导出会累积有损误差。'
    }
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : '处理失败，请检查参数后重试。'
  } finally {
    isBusy.value = false
  }
}

async function restoreOriginal() {
  if (!originalImageData.value) return
  clearMessages()
  isBusy.value = true
  try {
    workingImageData.value = cloneImageData(originalImageData.value)
    await setWorkingPreviewFromImageData(workingImageData.value)
    message.value = '已恢复到原始上传图像（内存无损副本）。'
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : '还原失败'
  } finally {
    isBusy.value = false
  }
}

async function downloadOutput() {
  if (!workingImageData.value || !selectedFile.value) return

  try {
    const preferPng = lastMode.value === 'decrypt' && form.usePngOnDecrypt
    const type: 'image/png' | 'image/jpeg' = preferPng ? 'image/png' : 'image/jpeg'
    const extension = preferPng ? '.png' : '.jpg'
    const canvas = imageDataToCanvas(workingImageData.value)
    const blob = await canvasToBlob(canvas, type, form.quality)
    const baseName = selectedFile.value.name.replace(/\.[^.]+$/, '')
    const suffix = lastMode.value === 'encrypt' ? '-etc-encrypted' : '-etc-decrypted'
    const href = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = href
    link.download = `${baseName}${suffix}${extension}`
    link.click()
    URL.revokeObjectURL(href)
  } catch (error) {
    clearMessages()
    errorMessage.value = error instanceof Error ? error.message : '下载失败'
  }
}

function exportConfig() {
  const configText = serializeEtcConfig(
    normalizeEtcConfig({
      key: form.key,
      blockSize: form.blockSize,
      quality: form.quality,
      enableRotateFlip: form.enableRotateFlip,
      enableChannelPermutation: form.enableChannelPermutation,
      enableNegativeTransform: form.enableNegativeTransform
    })
  )
  configPayload.value = configText
  clearMessages()
  message.value = '参数已导出到文本框，可复制保存。'
}

function applyJpegSafePreset() {
  form.blockSize = 8
  form.quality = 0.97
  form.enableRotateFlip = true
  form.enableChannelPermutation = false
  form.enableNegativeTransform = false
  clearMessages()
  message.value = '已应用 JPEG 友好预设：8x8 / 0.97 / 关闭通道置换与负片。'
}

function importConfig() {
  const parsed = parseEtcConfig(configPayload.value)
  if (!parsed) {
    clearMessages()
    errorMessage.value = '参数格式无效，请确认包含 version=etc-v1。'
    return
  }

  form.blockSize = parsed.blockSize
  form.quality = parsed.quality
  form.enableRotateFlip = parsed.enableRotateFlip
  form.enableChannelPermutation = parsed.enableChannelPermutation
  form.enableNegativeTransform = parsed.enableNegativeTransform
  clearMessages()
  message.value = '参数导入成功。'
}

async function canvasToBlob(canvas: HTMLCanvasElement, type: 'image/jpeg' | 'image/png', quality: number): Promise<Blob> {
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

onBeforeUnmount(() => {
  revokeUrl(originalUrl.value)
  revokeUrl(workingUrl.value)
})
</script>

<style scoped>
.etc-page {
  display: grid;
  gap: var(--space-xl);
}

.etc-hero {
  display: grid;
  gap: var(--space-sm);
}

.etc-hero .section-title {
  max-width: 12ch;
}

.etc-panel {
  padding: var(--space-xl);
}

.etc-panel__main {
  display: grid;
  grid-template-columns: minmax(280px, 420px) minmax(0, 1fr);
  gap: var(--space-xl);
}

.etc-controls {
  display: grid;
  gap: var(--space-md);
}

.control-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: var(--space-sm);
}

.control-group {
  display: grid;
  gap: 0.4rem;
}

.control-label {
  font-size: 0.82rem;
  color: var(--text-muted);
  font-family: var(--font-mono);
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

.control-input {
  width: 100%;
  min-height: 2.7rem;
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-md);
  background: var(--bg-panel-strong);
  color: var(--text-primary);
  padding: 0.55rem 0.75rem;
}

.visually-hidden-input {
  position: absolute;
  width: 1px;
  height: 1px;
  overflow: hidden;
  clip: rect(0 0 0 0);
  clip-path: inset(50%);
  white-space: nowrap;
}

.file-picker {
  display: grid;
  grid-template-columns: auto minmax(0, 1fr);
  gap: var(--space-sm);
  align-items: center;
  min-height: 2.8rem;
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-md);
  background: var(--bg-panel-strong);
  padding: 0.4rem;
}

.file-picker__button {
  min-height: 2rem;
  padding: 0.42rem 0.95rem;
  border-radius: var(--radius-pill);
  border: 1px solid var(--border-primary);
  color: var(--text-primary);
  background: var(--accent-soft);
  white-space: nowrap;
  transition: transform var(--duration-fast) var(--ease-standard), border-color var(--duration-fast) var(--ease-standard);
}

.file-picker__button:hover {
  transform: translateY(-1px);
  border-color: var(--border-strong);
}

.file-picker__name {
  min-width: 0;
  color: var(--text-secondary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  padding-right: 0.2rem;
}

.file-picker__name.is-empty {
  color: var(--text-muted);
}

.control-textarea {
  min-height: 6.4rem;
  resize: vertical;
  line-height: 1.55;
  font-family: var(--font-mono);
  font-size: 0.82rem;
}

.control-range {
  accent-color: var(--accent-primary);
}

.control-checks {
  display: grid;
  gap: 0.5rem;
  color: var(--text-secondary);
}

.control-checks label {
  display: inline-flex;
  align-items: center;
  gap: 0.45rem;
}

.etc-actions {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: var(--space-sm);
}

.etc-actions--secondary {
  grid-template-columns: repeat(3, minmax(0, 1fr));
}

.action-btn {
  min-height: 2.6rem;
  border-radius: var(--radius-pill);
  border: 1px solid var(--border-primary);
  color: var(--text-primary);
  background: rgba(255, 255, 255, 0.34);
  transition: transform var(--duration-fast) var(--ease-standard), border-color var(--duration-fast) var(--ease-standard);
}

.action-btn:hover:not(:disabled) {
  transform: translateY(-1px);
  border-color: var(--border-strong);
}

.action-btn:disabled {
  cursor: not-allowed;
  opacity: 0.52;
}

.action-btn--primary {
  background: var(--accent-soft);
  border-color: var(--border-strong);
}

.action-btn--secondary {
  background: var(--bg-panel);
}

.status-text {
  color: var(--text-secondary);
  line-height: 1.65;
}

.status-text--error {
  color: var(--color-error);
}

.etc-preview {
  display: grid;
  gap: var(--space-md);
}

.preview-card {
  border: 1px solid var(--border-primary);
  border-radius: var(--radius-lg);
  background: var(--bg-panel-strong);
  padding: var(--space-md);
  display: grid;
  gap: var(--space-sm);
}

.preview-card h3 {
  font-family: var(--font-display);
  font-size: 1.2rem;
  letter-spacing: -0.02em;
}

.preview-card p {
  color: var(--text-muted);
}

.preview-card img {
  width: 100%;
  max-height: 460px;
  object-fit: contain;
  border-radius: var(--radius-md);
  border: 1px solid var(--border-primary);
  background: var(--bg-panel);
}

@media (max-width: 1080px) {
  .etc-panel__main {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 760px) {
  .etc-panel {
    padding: var(--space-md);
  }

  .etc-actions {
    grid-template-columns: 1fr 1fr;
  }

  .control-grid {
    grid-template-columns: 1fr;
  }

  .file-picker {
    grid-template-columns: 1fr;
    gap: 0.5rem;
  }
}
</style>
