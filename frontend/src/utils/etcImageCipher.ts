export type EtcBlockSize = 8 | 16

export interface EtcCipherOptions {
  key: string
  blockSize: EtcBlockSize
  quality: number
  enableRotateFlip: boolean
  enableChannelPermutation: boolean
  enableNegativeTransform: boolean
}

export interface EtcSerializedConfig {
  version: 'etc-v1'
  blockSize: EtcBlockSize
  quality: number
  enableRotateFlip: boolean
  enableChannelPermutation: boolean
  enableNegativeTransform: boolean
}

type Rotation = 0 | 1 | 2 | 3
type ChannelPermutation = [0 | 1 | 2, 0 | 1 | 2, 0 | 1 | 2]

interface BlockInfo {
  index: number
  x: number
  y: number
  width: number
  height: number
  groupKey: string
}

interface BlockTransform {
  rotation: Rotation
  flipX: boolean
  flipY: boolean
  channelPermutationIndex: number
  negative: boolean
}

interface CipherPlan {
  blocks: BlockInfo[]
  permutation: number[]
  transforms: BlockTransform[]
}

const CHANNEL_PERMUTATIONS: ChannelPermutation[] = [
  [0, 1, 2],
  [0, 2, 1],
  [1, 0, 2],
  [1, 2, 0],
  [2, 0, 1],
  [2, 1, 0]
]

const INVERSE_CHANNEL_PERMUTATIONS: ChannelPermutation[] = CHANNEL_PERMUTATIONS.map((perm) => {
  const inverse = [0, 0, 0] as [0 | 1 | 2, 0 | 1 | 2, 0 | 1 | 2]
  perm.forEach((value, index) => {
    inverse[value] = index as 0 | 1 | 2
  })
  return inverse
})

export const ETC_DEFAULT_CONFIG: Omit<EtcCipherOptions, 'key'> = {
  blockSize: 8,
  quality: 0.97,
  enableRotateFlip: true,
  enableChannelPermutation: false,
  enableNegativeTransform: false
}

export function normalizeEtcConfig(raw: Partial<EtcCipherOptions> & { key?: string }): EtcCipherOptions {
  return {
    key: String(raw.key ?? '').trim(),
    blockSize: normalizeBlockSize(raw.blockSize),
    quality: normalizeQuality(raw.quality),
    enableRotateFlip: Boolean(raw.enableRotateFlip ?? ETC_DEFAULT_CONFIG.enableRotateFlip),
    enableChannelPermutation: Boolean(raw.enableChannelPermutation ?? ETC_DEFAULT_CONFIG.enableChannelPermutation),
    enableNegativeTransform: Boolean(raw.enableNegativeTransform ?? ETC_DEFAULT_CONFIG.enableNegativeTransform)
  }
}

export function serializeEtcConfig(config: EtcCipherOptions): string {
  const normalized = normalizeEtcConfig(config)
  const payload: EtcSerializedConfig = {
    version: 'etc-v1',
    blockSize: normalized.blockSize,
    quality: normalized.quality,
    enableRotateFlip: normalized.enableRotateFlip,
    enableChannelPermutation: normalized.enableChannelPermutation,
    enableNegativeTransform: normalized.enableNegativeTransform
  }
  return JSON.stringify(payload, null, 2)
}

export function parseEtcConfig(raw: string): Omit<EtcCipherOptions, 'key'> | null {
  try {
    const parsed = JSON.parse(raw) as Partial<EtcSerializedConfig>
    if (parsed.version !== 'etc-v1') {
      return null
    }
    return {
      blockSize: normalizeBlockSize(parsed.blockSize),
      quality: normalizeQuality(parsed.quality),
      enableRotateFlip: Boolean(parsed.enableRotateFlip),
      enableChannelPermutation: Boolean(parsed.enableChannelPermutation),
      enableNegativeTransform: Boolean(parsed.enableNegativeTransform)
    }
  } catch {
    return null
  }
}

export function encryptImageData(source: ImageData, options: EtcCipherOptions): ImageData {
  const normalized = normalizeEtcConfig(options)
  validateKey(normalized.key)

  const output = new ImageData(source.width, source.height)
  const plan = buildCipherPlan(source.width, source.height, normalized)
  const srcData = source.data
  const dstData = output.data

  for (const srcBlock of plan.blocks) {
    const destBlock = plan.blocks[plan.permutation[srcBlock.index]]
    const transform = plan.transforms[srcBlock.index]
    copyBlockForward(srcData, source.width, srcBlock, dstData, output.width, destBlock, transform)
  }

  return output
}

export function decryptImageData(source: ImageData, options: EtcCipherOptions): ImageData {
  const normalized = normalizeEtcConfig(options)
  validateKey(normalized.key)

  const output = new ImageData(source.width, source.height)
  const plan = buildCipherPlan(source.width, source.height, normalized)
  const srcData = source.data
  const dstData = output.data

  for (const originalBlock of plan.blocks) {
    const encryptedBlock = plan.blocks[plan.permutation[originalBlock.index]]
    const transform = plan.transforms[originalBlock.index]
    copyBlockInverse(srcData, source.width, encryptedBlock, dstData, output.width, originalBlock, transform)
  }

  return output
}

function validateKey(key: string) {
  if (!key || key.length < 4) {
    throw new Error('密钥长度至少需要 4 位')
  }
}

function normalizeBlockSize(blockSize?: number): EtcBlockSize {
  return blockSize === 8 ? 8 : 16
}

function normalizeQuality(quality?: number): number {
  const fallback = ETC_DEFAULT_CONFIG.quality
  if (typeof quality !== 'number' || Number.isNaN(quality)) {
    return fallback
  }
  return Math.min(0.98, Math.max(0.7, Number(quality.toFixed(2))))
}

function buildCipherPlan(width: number, height: number, options: EtcCipherOptions): CipherPlan {
  const blocks = buildBlocks(width, height, options.blockSize)
  const permutation = buildPermutation(blocks, options)
  const transforms = buildTransforms(blocks, options)
  return { blocks, permutation, transforms }
}

function buildBlocks(width: number, height: number, blockSize: EtcBlockSize): BlockInfo[] {
  const blocks: BlockInfo[] = []
  let index = 0
  for (let y = 0; y < height; y += blockSize) {
    const blockHeight = Math.min(blockSize, height - y)
    for (let x = 0; x < width; x += blockSize) {
      const blockWidth = Math.min(blockSize, width - x)
      blocks.push({
        index,
        x,
        y,
        width: blockWidth,
        height: blockHeight,
        groupKey: `${blockWidth}x${blockHeight}`
      })
      index += 1
    }
  }
  return blocks
}

function buildPermutation(blocks: BlockInfo[], options: EtcCipherOptions): number[] {
  const grouped = new Map<string, number[]>()
  blocks.forEach((block) => {
    const list = grouped.get(block.groupKey) ?? []
    list.push(block.index)
    grouped.set(block.groupKey, list)
  })

  const permutation = blocks.map((block) => block.index)
  const random = createDeterministicRandom(
    `${options.key}|perm|${options.blockSize}|${blocks.length}|${Number(options.enableRotateFlip)}${Number(options.enableChannelPermutation)}${Number(options.enableNegativeTransform)}`
  )

  grouped.forEach((indices) => {
    const shuffled = indices.slice()
    fisherYates(shuffled, random)
    indices.forEach((srcIndex, i) => {
      permutation[srcIndex] = shuffled[i]
    })
  })

  return permutation
}

function buildTransforms(blocks: BlockInfo[], options: EtcCipherOptions): BlockTransform[] {
  const random = createDeterministicRandom(
    `${options.key}|transform|${options.blockSize}|${blocks.length}|${Number(options.enableRotateFlip)}${Number(options.enableChannelPermutation)}${Number(options.enableNegativeTransform)}`
  )

  return blocks.map((block) => {
    const squareBlock = block.width === block.height
    const rotation = pickRotation(squareBlock, options.enableRotateFlip, random)
    const channelPermutationIndex = options.enableChannelPermutation
      ? Math.floor(random() * CHANNEL_PERMUTATIONS.length)
      : 0

    return {
      rotation,
      flipX: options.enableRotateFlip ? random() < 0.5 : false,
      flipY: options.enableRotateFlip ? random() < 0.5 : false,
      channelPermutationIndex,
      negative: options.enableNegativeTransform ? random() < 0.5 : false
    }
  })
}

function pickRotation(squareBlock: boolean, enableRotateFlip: boolean, random: () => number): Rotation {
  if (!enableRotateFlip) {
    return 0
  }
  if (!squareBlock) {
    return random() < 0.5 ? 0 : 2
  }
  return Math.floor(random() * 4) as Rotation
}

function copyBlockForward(
  sourceData: Uint8ClampedArray,
  sourceWidth: number,
  sourceBlock: BlockInfo,
  outputData: Uint8ClampedArray,
  outputWidth: number,
  destBlock: BlockInfo,
  transform: BlockTransform
) {
  for (let localY = 0; localY < sourceBlock.height; localY += 1) {
    for (let localX = 0; localX < sourceBlock.width; localX += 1) {
      const sourceOffset = pixelOffset(sourceBlock.x + localX, sourceBlock.y + localY, sourceWidth)
      const mapped = mapForwardCoordinate(localX, localY, sourceBlock.width, sourceBlock.height, transform)
      const destOffset = pixelOffset(destBlock.x + mapped.x, destBlock.y + mapped.y, outputWidth)
      writeForwardColor(sourceData, sourceOffset, outputData, destOffset, transform)
    }
  }
}

function copyBlockInverse(
  sourceData: Uint8ClampedArray,
  sourceWidth: number,
  encryptedBlock: BlockInfo,
  outputData: Uint8ClampedArray,
  outputWidth: number,
  originalBlock: BlockInfo,
  transform: BlockTransform
) {
  for (let localY = 0; localY < encryptedBlock.height; localY += 1) {
    for (let localX = 0; localX < encryptedBlock.width; localX += 1) {
      const sourceOffset = pixelOffset(encryptedBlock.x + localX, encryptedBlock.y + localY, sourceWidth)
      const mapped = mapInverseCoordinate(localX, localY, originalBlock.width, originalBlock.height, transform)
      const destOffset = pixelOffset(originalBlock.x + mapped.x, originalBlock.y + mapped.y, outputWidth)
      writeInverseColor(sourceData, sourceOffset, outputData, destOffset, transform)
    }
  }
}

function writeForwardColor(
  source: Uint8ClampedArray,
  sourceOffset: number,
  output: Uint8ClampedArray,
  outputOffset: number,
  transform: BlockTransform
) {
  const permutation = CHANNEL_PERMUTATIONS[transform.channelPermutationIndex]
  const channels: [number, number, number] = [
    source[sourceOffset + permutation[0]],
    source[sourceOffset + permutation[1]],
    source[sourceOffset + permutation[2]]
  ]

  output[outputOffset] = transform.negative ? 255 - channels[0] : channels[0]
  output[outputOffset + 1] = transform.negative ? 255 - channels[1] : channels[1]
  output[outputOffset + 2] = transform.negative ? 255 - channels[2] : channels[2]
  output[outputOffset + 3] = source[sourceOffset + 3]
}

function writeInverseColor(
  source: Uint8ClampedArray,
  sourceOffset: number,
  output: Uint8ClampedArray,
  outputOffset: number,
  transform: BlockTransform
) {
  const baseR = transform.negative ? 255 - source[sourceOffset] : source[sourceOffset]
  const baseG = transform.negative ? 255 - source[sourceOffset + 1] : source[sourceOffset + 1]
  const baseB = transform.negative ? 255 - source[sourceOffset + 2] : source[sourceOffset + 2]
  const inversePermutation = INVERSE_CHANNEL_PERMUTATIONS[transform.channelPermutationIndex]
  const channels: [number, number, number] = [baseR, baseG, baseB]

  output[outputOffset] = channels[inversePermutation[0]]
  output[outputOffset + 1] = channels[inversePermutation[1]]
  output[outputOffset + 2] = channels[inversePermutation[2]]
  output[outputOffset + 3] = source[sourceOffset + 3]
}

function mapForwardCoordinate(x: number, y: number, width: number, height: number, transform: BlockTransform) {
  let mappedX = x
  let mappedY = y
  let rotateWidth = width
  let rotateHeight = height

  switch (transform.rotation) {
    case 1:
      mappedX = height - 1 - y
      mappedY = x
      rotateWidth = height
      rotateHeight = width
      break
    case 2:
      mappedX = width - 1 - x
      mappedY = height - 1 - y
      break
    case 3:
      mappedX = y
      mappedY = width - 1 - x
      rotateWidth = height
      rotateHeight = width
      break
  }

  if (transform.flipX) {
    mappedX = rotateWidth - 1 - mappedX
  }
  if (transform.flipY) {
    mappedY = rotateHeight - 1 - mappedY
  }

  return { x: mappedX, y: mappedY }
}

function mapInverseCoordinate(x: number, y: number, width: number, height: number, transform: BlockTransform) {
  const rotateWidth = transform.rotation % 2 === 0 ? width : height
  const rotateHeight = transform.rotation % 2 === 0 ? height : width

  let mappedX = transform.flipX ? rotateWidth - 1 - x : x
  let mappedY = transform.flipY ? rotateHeight - 1 - y : y

  switch (transform.rotation) {
    case 1:
      return { x: mappedY, y: height - 1 - mappedX }
    case 2:
      return { x: width - 1 - mappedX, y: height - 1 - mappedY }
    case 3:
      return { x: width - 1 - mappedY, y: mappedX }
    default:
      return { x: mappedX, y: mappedY }
  }
}

function pixelOffset(x: number, y: number, width: number): number {
  return (y * width + x) * 4
}

function fisherYates<T>(arr: T[], random: () => number) {
  for (let i = arr.length - 1; i > 0; i -= 1) {
    const j = Math.floor(random() * (i + 1))
    ;[arr[i], arr[j]] = [arr[j], arr[i]]
  }
}

function createDeterministicRandom(seedText: string): () => number {
  const seed = xmur3(seedText)
  return mulberry32(seed())
}

function xmur3(str: string): () => number {
  let h = 1779033703 ^ str.length
  for (let i = 0; i < str.length; i += 1) {
    h = Math.imul(h ^ str.charCodeAt(i), 3432918353)
    h = (h << 13) | (h >>> 19)
  }
  return () => {
    h = Math.imul(h ^ (h >>> 16), 2246822507)
    h = Math.imul(h ^ (h >>> 13), 3266489909)
    h ^= h >>> 16
    return h >>> 0
  }
}

function mulberry32(seed: number): () => number {
  return () => {
    let t = (seed += 0x6d2b79f5)
    t = Math.imul(t ^ (t >>> 15), t | 1)
    t ^= t + Math.imul(t ^ (t >>> 7), t | 61)
    return ((t ^ (t >>> 14)) >>> 0) / 4294967296
  }
}
