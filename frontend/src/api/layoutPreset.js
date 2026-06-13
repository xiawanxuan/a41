import request from './request'

export function listPresets() {
  return request({
    url: '/layout-presets',
    method: 'get'
  })
}

export function getPreset(id) {
  return request({
    url: `/layout-presets/${id}`,
    method: 'get'
  })
}

export function getDefaultPreset() {
  return request({
    url: '/layout-presets/default',
    method: 'get'
  })
}

export function createPreset(data) {
  return request({
    url: '/layout-presets',
    method: 'post',
    data
  })
}

export function updatePreset(id, data) {
  return request({
    url: `/layout-presets/${id}`,
    method: 'put',
    data
  })
}

export function setDefaultPreset(id) {
  return request({
    url: `/layout-presets/${id}/default`,
    method: 'put'
  })
}

export function deletePreset(id) {
  return request({
    url: `/layout-presets/${id}`,
    method: 'delete'
  })
}

export function presetToLayout(preset) {
  if (!preset) return null
  return {
    columns: 12,
    chars_per_column: preset.chars_per_column || 20,
    column_gap: preset.column_gap || 20,
    font_size: preset.font_size || 18,
    font_family: preset.font_family || 'SimSun, serif',
    line_height: preset.line_height || 1.8,
    char_spacing: preset.char_spacing || 0,
    text_color: preset.text_color || '#1a1a1a',
    background_color: preset.background_color || '#f5f0e8',
    show_border: preset.show_border !== false,
    custom_css: preset.custom_css || ''
  }
}

export function layoutToPreset(layout, name, description = '') {
  return {
    name,
    description,
    chars_per_column: layout.charsPerColumn || layout.chars_per_column || 20,
    column_gap: layout.columnGap || layout.column_gap || 20,
    font_size: layout.fontSize || layout.font_size || 18,
    font_family: layout.fontFamily || layout.font_family || 'SimSun, serif',
    line_height: layout.lineHeight || layout.line_height || 1.8,
    char_spacing: layout.charSpacing ?? layout.char_spacing ?? 0,
    text_color: layout.textColor || layout.text_color || '#1a1a1a',
    background_color: layout.backgroundColor || layout.background_color || '#f5f0e8',
    show_border: layout.showBorder ?? layout.show_border ?? true,
    custom_css: layout.customCss || layout.custom_css || ''
  }
}
