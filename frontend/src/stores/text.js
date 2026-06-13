import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { getText, listTexts, createText, updateText, deleteText } from '@/api/text'
import { getLayoutByText, saveLayout } from '@/api/layout'
import {
  getAnnotationsByText,
  createAnnotation,
  updateAnnotation,
  deleteAnnotation,
  batchSubmitAnnotations
} from '@/api/annotation'

export const useTextStore = defineStore('text', () => {
  const texts = ref([])
  const total = ref(0)
  const currentText = ref(null)
  const layoutConfig = ref(null)
  const annotations = ref([])
  const loading = ref(false)

  const defaultLayout = {
    columns: 12,
    chars_per_column: 20,
    column_gap: 20,
    font_size: 18,
    font_family: 'SimSun, serif',
    line_height: 1.8,
    text_color: '#1a1a1a',
    background_color: '#f5f0e8',
    show_border: true,
    custom_css: ''
  }

  async function fetchTexts(params = {}) {
    loading.value = true
    try {
      const res = await listTexts(params)
      texts.value = res.items || []
      total.value = res.total || 0
      return res
    } finally {
      loading.value = false
    }
  }

  async function fetchText(id) {
    loading.value = true
    try {
      const text = await getText(id)
      currentText.value = text
      return text
    } finally {
      loading.value = false
    }
  }

  async function createNewText(data) {
    loading.value = true
    try {
      const text = await createText(data)
      return text
    } finally {
      loading.value = false
    }
  }

  async function updateTextData(id, data) {
    loading.value = true
    try {
      const text = await updateText(id, data)
      if (currentText.value?.id === id) {
        currentText.value = text
      }
      return text
    } finally {
      loading.value = false
    }
  }

  async function deleteTextData(id) {
    loading.value = true
    try {
      await deleteText(id)
      texts.value = texts.value.filter(t => t.id !== id)
    } finally {
      loading.value = false
    }
  }

  async function fetchLayout(textId) {
    try {
      const layout = await getLayoutByText(textId)
      layoutConfig.value = layout || { ...defaultLayout }
      return layoutConfig.value
    } catch (e) {
      layoutConfig.value = { ...defaultLayout }
      return layoutConfig.value
    }
  }

  async function saveLayoutConfig(textId, config) {
    loading.value = true
    try {
      const layout = await saveLayout(textId, config)
      layoutConfig.value = layout
      return layout
    } finally {
      loading.value = false
    }
  }

  async function fetchAnnotations(textId, userId = '') {
    try {
      const anns = await getAnnotationsByText(textId, userId)
      annotations.value = anns || []
      return annotations.value
    } catch (e) {
      annotations.value = []
      return []
    }
  }

  async function addAnnotation(textId, annotation, userId = '') {
    try {
      const data = {
        text_id: textId,
        user_id: userId,
        ...annotation
      }
      const ann = await createAnnotation(data)
      annotations.value.push(ann)
      return ann
    } catch (e) {
      throw e
    }
  }

  async function removeAnnotation(id) {
    try {
      await deleteAnnotation(id)
      annotations.value = annotations.value.filter(a => a.id !== id)
    } catch (e) {
      throw e
    }
  }

  async function submitAllAnnotations(textId, userId = '') {
    loading.value = true
    try {
      const data = {
        text_id: textId,
        user_id: userId,
        annotations: annotations.value.map(a => ({
          start_pos: a.start_pos,
          end_pos: a.end_pos,
          type: a.type,
          content: a.content || '',
          color: a.color || ''
        }))
      }
      const result = await batchSubmitAnnotations(data)
      return result
    } finally {
      loading.value = false
    }
  }

  function resetCurrent() {
    currentText.value = null
    layoutConfig.value = null
    annotations.value = []
  }

  return {
    texts,
    total,
    currentText,
    layoutConfig,
    annotations,
    loading,
    defaultLayout,
    fetchTexts,
    fetchText,
    createNewText,
    updateTextData,
    deleteTextData,
    fetchLayout,
    saveLayoutConfig,
    fetchAnnotations,
    addAnnotation,
    removeAnnotation,
    submitAllAnnotations,
    resetCurrent
  }
})
