<template>
  <div class="annotation-toolbar card">
    <div class="toolbar-header">
      <h3 class="toolbar-title">标注工具</h3>
      <span class="annotation-count">共 {{ annotations.length }} 条标注</span>
    </div>

    <div class="tool-section">
      <div class="section-label">断句标点</div>
      <div class="tool-buttons">
        <button
          class="tool-btn"
          :class="{ active: currentTool === 'period' }"
          @click="selectTool('period')"
          :disabled="!hasSelection"
          title="句号"
        >
          <span class="tool-icon">。</span>
          <span class="tool-text">句号</span>
        </button>
        <button
          class="tool-btn"
          :class="{ active: currentTool === 'comma' }"
          @click="selectTool('comma')"
          :disabled="!hasSelection"
          title="逗号"
        >
          <span class="tool-icon">，</span>
          <span class="tool-text">逗号</span>
        </button>
        <button
          class="tool-btn"
          :class="{ active: currentTool === 'question' }"
          @click="selectTool('question')"
          :disabled="!hasSelection"
          title="问号"
        >
          <span class="tool-icon">？</span>
          <span class="tool-text">问号</span>
        </button>
        <button
          class="tool-btn"
          :class="{ active: currentTool === 'exclaim' }"
          @click="selectTool('exclaim')"
          :disabled="!hasSelection"
          title="叹号"
        >
          <span class="tool-icon">！</span>
          <span class="tool-text">叹号</span>
        </button>
      </div>
    </div>

    <div class="tool-section">
      <div class="section-label">其他标注</div>
      <div class="tool-buttons">
        <button
          class="tool-btn"
          :class="{ active: currentTool === 'highlight' }"
          @click="selectTool('highlight')"
          :disabled="!hasSelection"
          title="高亮"
        >
          <span class="tool-icon highlight-icon">✦</span>
          <span class="tool-text">高亮</span>
        </button>
        <button
          class="tool-btn"
          :class="{ active: currentTool === 'comment' }"
          @click="selectTool('comment')"
          :disabled="!hasSelection"
          title="批注"
        >
          <span class="tool-icon">📝</span>
          <span class="tool-text">批注</span>
        </button>
      </div>
    </div>

    <div v-if="currentTool === 'highlight'" class="tool-section">
      <div class="section-label">高亮颜色</div>
      <div class="color-picker">
        <button
          v-for="color in highlightColors"
          :key="color.value"
          class="color-btn"
          :class="{ active: selectedColor === color.value }"
          :style="{ backgroundColor: color.value }"
          @click="selectedColor = color.value"
          :title="color.name"
        ></button>
      </div>
    </div>

    <div v-if="currentTool === 'comment'" class="tool-section">
      <div class="section-label">批注内容</div>
      <textarea
        v-model="commentContent"
        class="form-textarea"
        placeholder="请输入批注内容..."
        rows="4"
      ></textarea>
    </div>

    <div v-if="hasSelection" class="selection-info">
      <div class="section-label">选中内容</div>
      <div class="selection-text">{{ selectionText }}</div>
      <div class="selection-pos">位置: {{ selection.startPos }} - {{ selection.endPos }}</div>
    </div>

    <div class="action-buttons">
      <button
        class="btn btn-primary"
        :disabled="!canAddAnnotation"
        @click="addAnnotation"
      >
        添加标注
      </button>
      <button class="btn btn-secondary" @click="clearSelection">
        清除选择
      </button>
    </div>

    <div class="tool-section">
      <div class="section-label">标注列表</div>
      <div class="annotation-list">
        <div
          v-for="ann in annotations"
          :key="ann.id"
          class="annotation-item"
          :class="`type-${ann.type}`"
        >
          <div class="annotation-header">
            <span class="annotation-type">{{ getTypeLabel(ann.type) }}</span>
            <span class="annotation-pos">{{ ann.start_pos }}-{{ ann.end_pos }}</span>
          </div>
          <div v-if="ann.content" class="annotation-content">{{ ann.content }}</div>
          <button
            class="annotation-delete"
            @click="deleteAnnotation(ann.id)"
            title="删除标注"
          >
            ×
          </button>
        </div>
        <div v-if="annotations.length === 0" class="empty-list">
          暂无标注
        </div>
      </div>
    </div>

    <div class="submit-section">
      <button class="btn btn-primary btn-block" @click="submitAll">
        提交所有标注
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const props = defineProps({
  selection: {
    type: Object,
    default: () => ({ startPos: -1, endPos: -1, text: '' })
  },
  annotations: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['add-annotation', 'delete-annotation', 'submit-all', 'clear-selection'])

const currentTool = ref('period')
const selectedColor = ref('#ffd700')
const commentContent = ref('')

const highlightColors = [
  { name: '黄色', value: '#ffd700' },
  { name: '绿色', value: '#90ee90' },
  { name: '蓝色', value: '#87ceeb' },
  { name: '粉色', value: '#ffb6c1' },
  { name: '橙色', value: '#ffa500' }
]

const hasSelection = computed(() => {
  return props.selection && props.selection.startPos >= 0 && props.selection.endPos > props.selection.startPos
})

const selectionText = computed(() => {
  return props.selection?.text || ''
})

const canAddAnnotation = computed(() => {
  if (!hasSelection.value || !currentTool.value) return false
  if (currentTool.value === 'comment' && !commentContent.value.trim()) return false
  return true
})

function selectTool(tool) {
  currentTool.value = tool
}

function addAnnotation() {
  if (!canAddAnnotation.value) return

  const annotation = {
    start_pos: props.selection.startPos,
    end_pos: props.selection.endPos,
    type: currentTool.value,
    content: currentTool.value === 'comment' ? commentContent.value : '',
    color: currentTool.value === 'highlight' ? selectedColor.value : ''
  }

  emit('add-annotation', annotation)
  commentContent.value = ''
}

function deleteAnnotation(id) {
  emit('delete-annotation', id)
}

function clearSelection() {
  emit('clear-selection')
}

function submitAll() {
  emit('submit-all')
}

function getTypeLabel(type) {
  const labels = {
    period: '句号',
    comma: '逗号',
    question: '问号',
    exclaim: '叹号',
    highlight: '高亮',
    comment: '批注'
  }
  return labels[type] || type
}
</script>

<style scoped>
.annotation-toolbar {
  width: 300px;
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.toolbar-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-bottom: 12px;
  border-bottom: 1px solid #e8e0d4;
}

.toolbar-title {
  font-size: 16px;
  font-weight: 600;
  color: #5a4030;
  margin: 0;
}

.annotation-count {
  font-size: 13px;
  color: #8b7355;
  background: #f5f0e8;
  padding: 4px 10px;
  border-radius: 12px;
}

.tool-section {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.section-label {
  font-size: 13px;
  font-weight: 500;
  color: #8b7355;
}

.tool-buttons {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 8px;
}

.tool-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  padding: 12px 8px;
  background: #f8f6f3;
  border: 1px solid #e8e0d4;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
  font-size: 12px;
  color: #5a4030;
}

.tool-btn:hover:not(:disabled) {
  background: #f0ebe3;
  border-color: #d4c8b8;
}

.tool-btn.active {
  background: #8b4513;
  border-color: #8b4513;
  color: #fff;
}

.tool-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.tool-icon {
  font-size: 20px;
  font-weight: bold;
}

.highlight-icon {
  color: #ffa500;
}

.tool-btn.active .highlight-icon {
  color: #fff;
}

.tool-text {
  font-size: 12px;
}

.color-picker {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.color-btn {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  border: 2px solid transparent;
  cursor: pointer;
  transition: all 0.2s;
}

.color-btn:hover {
  transform: scale(1.1);
}

.color-btn.active {
  border-color: #5a4030;
}

.selection-info {
  background: #f5f0e8;
  padding: 12px;
  border-radius: 6px;
}

.selection-text {
  font-size: 13px;
  color: #1a1a1a;
  word-break: break-all;
  margin-bottom: 6px;
  padding: 6px;
  background: #fff;
  border-radius: 4px;
  border: 1px dashed #d4c8b8;
}

.selection-pos {
  font-size: 12px;
  color: #8b7355;
}

.action-buttons {
  display: flex;
  gap: 8px;
}

.action-buttons .btn {
  flex: 1;
}

.annotation-list {
  max-height: 240px;
  overflow-y: auto;
  border: 1px solid #e8e0d4;
  border-radius: 6px;
  background: #faf8f5;
}

.annotation-item {
  position: relative;
  padding: 10px 12px;
  border-bottom: 1px solid #e8e0d4;
  background: #fff;
}

.annotation-item:last-child {
  border-bottom: none;
}

.annotation-item.type-period .annotation-type { color: #8b4513; }
.annotation-item.type-comma .annotation-type { color: #8b4513; }
.annotation-item.type-question .annotation-type { color: #2980b9; }
.annotation-item.type-exclaim .annotation-type { color: #c0392b; }
.annotation-item.type-highlight .annotation-type { color: #f39c12; }
.annotation-item.type-comment .annotation-type { color: #27ae60; }

.annotation-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 4px;
}

.annotation-type {
  font-size: 12px;
  font-weight: 600;
}

.annotation-pos {
  font-size: 11px;
  color: #8b7355;
}

.annotation-content {
  font-size: 12px;
  color: #5a4030;
  line-height: 1.5;
}

.annotation-delete {
  position: absolute;
  top: 8px;
  right: 10px;
  width: 20px;
  height: 20px;
  background: #f5e0e0;
  color: #c0392b;
  border: none;
  border-radius: 50%;
  cursor: pointer;
  font-size: 14px;
  line-height: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.2s;
}

.annotation-item:hover .annotation-delete {
  opacity: 1;
}

.annotation-delete:hover {
  background: #e8c8c8;
}

.empty-list {
  padding: 30px;
  text-align: center;
  color: #b0a090;
  font-size: 13px;
}

.submit-section {
  margin-top: auto;
  padding-top: 16px;
  border-top: 1px solid #e8e0d4;
}

.btn-block {
  width: 100%;
}
</style>
