<template>
  <div class="text-reader-page">
    <div class="page-header">
      <button class="btn btn-secondary" @click="goBack">
        ← 返回列表
      </button>
      <div class="header-info">
        <h2 class="page-title">{{ textStore.currentText?.title || '加载中...' }}</h2>
        <div class="header-meta" v-if="textStore.currentText">
          <span v-if="textStore.currentText.author">{{ textStore.currentText.author }}</span>
          <span v-if="textStore.currentText.dynasty" class="meta-divider">·</span>
          <span v-if="textStore.currentText.dynasty">{{ textStore.currentText.dynasty }}</span>
          <span class="meta-divider">·</span>
          <span>{{ textStore.currentText.content?.length || 0 }} 字</span>
        </div>
      </div>
      <div class="header-actions">
        <button
          class="btn"
          :class="showConfig ? 'btn-primary' : 'btn-secondary'"
          @click="showConfig = !showConfig"
        >
          ⚙ 版式设置
        </button>
        <button class="btn btn-secondary" @click="goToEdit">
          ✎ 编辑文本
        </button>
      </div>
    </div>

    <div class="reader-container">
      <aside v-if="showConfig" class="config-sidebar">
        <LayoutConfigPanel
          :layout="currentLayout"
          @change="handleLayoutChange"
          @save="handleLayoutSave"
        />
      </aside>

      <main class="reader-main">
        <div v-if="loading" class="loading-state">
          <div class="loading-spinner"></div>
          <p>加载文本中...</p>
        </div>

        <VerticalText
          v-else
          ref="verticalTextRef"
          :text="textStore.currentText?.content || ''"
          :layout="currentLayout"
          :annotations="textStore.annotations"
          :editable="true"
          @selection="handleSelection"
        />
      </main>

      <aside class="annotation-sidebar">
        <AnnotationToolbar
          :selection="currentSelection"
          :annotations="textStore.annotations"
          @add-annotation="handleAddAnnotation"
          @delete-annotation="handleDeleteAnnotation"
          @submit-all="handleSubmitAll"
          @clear-selection="handleClearSelection"
        />
      </aside>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useTextStore } from '@/stores/text'
import VerticalText from '@/components/VerticalText.vue'
import AnnotationToolbar from '@/components/AnnotationToolbar.vue'
import LayoutConfigPanel from '@/components/LayoutConfigPanel.vue'

const router = useRouter()
const route = useRoute()
const textStore = useTextStore()

const loading = ref(true)
const showConfig = ref(false)
const verticalTextRef = ref(null)

const currentSelection = reactive({
  startPos: -1,
  endPos: -1,
  text: ''
})

const currentLayout = ref({
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
})

async function loadAll() {
  loading.value = true
  try {
    const textId = route.params.id

    await Promise.all([
      textStore.fetchText(textId),
      textStore.fetchLayout(textId),
      textStore.fetchAnnotations(textId, 'default-user')
    ])

    if (textStore.layoutConfig) {
      Object.assign(currentLayout.value, textStore.layoutConfig)
    }
  } catch (e) {
    console.error('Failed to load:', e)
    alert('加载失败')
  } finally {
    loading.value = false
  }
}

function handleSelection(selection) {
  currentSelection.startPos = selection.startPos
  currentSelection.endPos = selection.endPos
  currentSelection.text = selection.text
}

function handleClearSelection() {
  currentSelection.startPos = -1
  currentSelection.endPos = -1
  currentSelection.text = ''
  verticalTextRef.value?.clearSelection()
}

async function handleAddAnnotation(annotation) {
  try {
    await textStore.addAnnotation(
      route.params.id,
      annotation,
      'default-user'
    )
    handleClearSelection()
  } catch (e) {
    console.error('Add annotation failed:', e)
    alert('添加标注失败')
  }
}

async function handleDeleteAnnotation(id) {
  try {
    await textStore.removeAnnotation(id)
  } catch (e) {
    console.error('Delete annotation failed:', e)
    alert('删除标注失败')
  }
}

async function handleSubmitAll() {
  if (textStore.annotations.length === 0) {
    alert('暂无标注可提交')
    return
  }
  if (!confirm(`确定提交 ${textStore.annotations.length} 条标注吗？`)) {
    return
  }
  try {
    await textStore.submitAllAnnotations(route.params.id, 'default-user')
    alert('提交成功！')
    await textStore.fetchAnnotations(route.params.id, 'default-user')
  } catch (e) {
    console.error('Submit failed:', e)
    alert('提交失败')
  }
}

function handleLayoutChange(layout) {
  Object.assign(currentLayout.value, layout)
}

async function handleLayoutSave(layout) {
  try {
    await textStore.saveLayoutConfig(route.params.id, layout)
    alert('版式配置已保存')
  } catch (e) {
    console.error('Save layout failed:', e)
    alert('保存版式配置失败')
  }
}

function goBack() {
  router.back()
}

function goToEdit() {
  router.push(`/texts/${route.params.id}/edit`)
}

onMounted(() => {
  loadAll()
})

onUnmounted(() => {
  textStore.resetCurrent()
})
</script>

<style scoped>
.text-reader-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
  min-height: calc(100vh - 80px);
}

.page-header {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 12px 20px;
  background: #fff;
  border-radius: 10px;
  box-shadow: 0 2px 8px rgba(139, 69, 19, 0.08);
}

.header-info {
  flex: 1;
}

.page-title {
  font-size: 20px;
  font-weight: 600;
  color: #5a4030;
  margin: 0 0 4px 0;
}

.header-meta {
  font-size: 13px;
  color: #8b7355;
}

.meta-divider {
  margin: 0 6px;
  color: #d4c8b8;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.reader-container {
  display: flex;
  gap: 16px;
  flex: 1;
  min-height: 0;
}

.config-sidebar {
  flex-shrink: 0;
}

.reader-main {
  flex: 1;
  min-width: 0;
  overflow: auto;
}

.annotation-sidebar {
  flex-shrink: 0;
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 100px 20px;
  gap: 16px;
  background: #fff;
  border-radius: 10px;
  box-shadow: 0 2px 8px rgba(139, 69, 19, 0.08);
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 3px solid #e8e0d4;
  border-top-color: #8b4513;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
