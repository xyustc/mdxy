/**
 * 证件照规格配置
 * 中国常用证件照尺寸标准（300 DPI）
 */

export type PhotoSpecId = 'one-inch' | 'two-inch' | 'small-two-inch'

export interface PhotoSpec {
  id: PhotoSpecId
  name: string
  nameEn: string
  widthMm: number
  heightMm: number
  widthPx: number
  heightPx: number
  dpi: number
  headHeightRatio: number
  description: string
  usage: string[]
}

export const ID_PHOTO_SPECS: Record<PhotoSpecId, PhotoSpec> = {
  'one-inch': {
    id: 'one-inch',
    name: '一寸',
    nameEn: '1-inch',
    widthMm: 25,
    heightMm: 35,
    widthPx: 295, // 25mm * 300dpi / 25.4
    heightPx: 413, // 35mm * 300dpi / 25.4
    dpi: 300,
    headHeightRatio: 0.67,
    description: '常用于工作证、学生证、保险等',
    usage: ['工作证', '学生证', '保险', '简历', '社保卡']
  },
  'two-inch': {
    id: 'two-inch',
    name: '二寸',
    nameEn: '2-inch',
    widthMm: 35,
    heightMm: 49,
    widthPx: 413,
    heightPx: 579,
    dpi: 300,
    headHeightRatio: 0.67,
    description: '常用于护照、签证、毕业证等',
    usage: ['护照', '签证', '毕业证', '学位证', '教师资格证']
  },
  'small-two-inch': {
    id: 'small-two-inch',
    name: '小二寸',
    nameEn: 'Small 2-inch',
    widthMm: 33,
    heightMm: 48,
    widthPx: 390,
    heightPx: 567,
    dpi: 300,
    headHeightRatio: 0.67,
    description: '常用于护照、港澳通行证等',
    usage: ['护照', '港澳通行证', '台湾通行证']
  }
} as const

export type BackgroundColorId = 'transparent' | 'white' | 'blue' | 'red'

export interface BackgroundColor {
  id: BackgroundColorId
  name: string
  hex: string
  usage: string[]
}

export const BACKGROUND_COLORS: Record<BackgroundColorId, BackgroundColor> = {
  transparent: {
    id: 'transparent',
    name: '透明',
    hex: 'transparent',
    usage: ['自定义背景']
  },
  white: {
    id: 'white',
    name: '白底',
    hex: '#FFFFFF',
    usage: ['护照', '签证', '身份证', '驾驶证']
  },
  blue: {
    id: 'blue',
    name: '蓝底',
    hex: '#438EDB',
    usage: ['毕业证', '工作证', '简历', '职称证']
  },
  red: {
    id: 'red',
    name: '红底',
    hex: '#D12E2E',
    usage: ['保险', '医保', '结婚照', '居住证']
  }
} as const

export const DEFAULT_SPEC_ID: PhotoSpecId = 'one-inch'
export const DEFAULT_BACKGROUND_ID: BackgroundColorId = 'transparent'

/**
 * 获取规格列表（用于 UI 遍历）
 */
export function getSpecList(): PhotoSpec[] {
  return Object.values(ID_PHOTO_SPECS)
}

/**
 * 获取背景色列表（用于 UI 遍历）
 */
export function getBackgroundList(): BackgroundColor[] {
  return Object.values(BACKGROUND_COLORS)
}

/**
 * 根据像素计算裁剪区域
 */
export function calculateCropArea(
  imageWidth: number,
  imageHeight: number,
  spec: PhotoSpec
): { x: number; y: number; width: number; height: number } {
  const targetRatio = spec.widthPx / spec.heightPx
  const imageRatio = imageWidth / imageHeight

  if (imageRatio > targetRatio) {
    // 图片更宽，以高度为准
    const height = imageHeight
    const width = height * targetRatio
    const x = (imageWidth - width) / 2
    return { x, y: 0, width, height }
  } else {
    // 图片更高，以宽度为准
    const width = imageWidth
    const height = width / targetRatio
    const y = (imageHeight - height) / 2
    return { x: 0, y, width, height }
  }
}