<template>
  <div
    class="vertical-text-container"
    :style="containerStyle"
    ref="containerRef"
  >
    <div
      class="vertical-text-content"
      :class="{ 'show-border': layout.showBorder }"
      :style="contentStyle"
      @mouseup="handleMouseUp"
      @mousedown="handleMouseDown"
      @mousemove="handleMouseMove"
    >
      <template v-for="(col, colIndex) in columns" :key="colIndex">
        <div class="vertical-column" :style="columnStyle">
          <span
            v-for="(char, charIndex) in col.chars"
            :key="charIndex"
            class="char-item"
            :class="{
              'has-period': char.hasPeriod,
              'has-comma': char.hasComma,
              'has-question': char.hasQuestion,
              'has-exclaim': char.hasExclaim,
              'highlighted': char.highlighted,
              'selected': isCharSelected(colIndex, charIndex)
            }"
            :style="getCharStyle(char)"
            :data-col="colIndex"
            :data-char="charIndex"
            :data-pos="char.pos"
          >
            {{ char.char }}
            <span v-if="char.hasPeriod" class="punctuation period">。</span>
            <span v-else-if="char.hasComma" class="punctuation comma">，</span>
            <span v-else-if="char.hasQuestion" class="punctuation question">？</span>
            <span v-else-if="char.hasExclaim" class="punctuation exclaim">！</span>
          </span>
        </div>
      </template>
    </div>
    <div
      v-if="selectionActive"
      class="selection-overlay"
      :style="selectionOverlayStyle"
    ></div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, nextTick } from 'vue'

const props = defineProps({
  text: {
    type: String,
    default: ''
  },
  layout: {
    type: Object,
    default: () => ({
      columns: 12,
      charsPerColumn: 20,
      columnGap: 20,
      fontSize: 18,
      fontFamily: 'SimSun, serif',
      lineHeight: 1.8,
      textColor: '#1a1a1a',
      backgroundColor: '#f5f0e8',
      showBorder: true,
      customCSS: ''
    })
  },
  annotations: {
    type: Array,
    default: () => []
  },
  editable: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['selection', 'annotation-click'])

const containerRef = ref(null)
const selectionStart = ref(null)
const selectionEnd = ref(null)
const selectionActive = ref(false)
const isSelecting = ref(false)

const containerStyle = computed(() => ({
  backgroundColor: props.layout.backgroundColor,
  padding: '20px',
  minHeight: '500px',
  overflowX: 'auto',
  position: 'relative'
}))

const contentStyle = computed(() => ({
  fontFamily: props.layout.fontFamily,
  fontSize: props.layout.fontSize + 'px',
  color: props.layout.textColor,
  lineHeight: props.layout.lineHeight,
  writingMode: 'vertical-rl',
  textOrientation: 'upright',
  direction: 'rtl',
  display: 'flex',
  gap: props.layout.columnGap + 'px',
  padding: '20px',
  border: props.layout.showBorder ? '1px solid #d4c8b8' : 'none',
  borderRadius: '4px'
}))

const columnStyle = computed(() => ({
  height: props.layout.charsPerColumn * props.layout.fontSize * props.layout.lineHeight + 'px',
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center'
}))

const charStyle = computed(() => ({
  width: props.layout.fontSize * 1.2 + 'px',
  height: props.layout.fontSize * props.layout.lineHeight + 'px',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  position: 'relative',
  cursor: props.editable ? 'text' : 'default',
  userSelect: 'none'
}))

const columns = computed(() => {
  const chars = props.text.split('')
  const cols = []
  const charsPerCol = props.layout.charsPerColumn
  const numCols = Math.ceil(chars.length / charsPerCol)

  for (let i = 0; i < numCols; i++) {
    const colChars = []
    for (let j = 0; j < charsPerCol; j++) {
      const pos = i * charsPerCol + j
      if (pos < chars.length) {
        const char = chars[pos]
        colChars.push({
          char,
          pos,
          hasPeriod: false,
          hasComma: false,
          hasQuestion: false,
          hasExclaim: false,
          highlighted: false
        })
      }
    }
    cols.push({ chars: colChars })
  }

  props.annotations.forEach(ann => {
    if (ann.type === 'period' || ann.type === 'comma' || ann.type === 'question' || ann.type === 'exclaim') {
      const endPos = ann.end_pos - 1
      const colIndex = Math.floor(endPos / charsPerCol)
      const charIndex = endPos % charsPerCol
      if (cols[colIndex] && cols[colIndex].chars[charIndex]) {
        if (ann.type === 'period') cols[colIndex].chars[charIndex].hasPeriod = true
        if (ann.type === 'comma') cols[colIndex].chars[charIndex].hasComma = true
        if (ann.type === 'question') cols[colIndex].chars[charIndex].hasQuestion = true
        if (ann.type === 'exclaim') cols[colIndex].chars[charIndex].hasExclaim = true
      }
    } else if (ann.type === 'highlight') {
      for (let pos = ann.start_pos; pos < ann.end_pos; pos++) {
        const colIndex = Math.floor(pos / charsPerCol)
        const charIndex = pos % charsPerCol
        if (cols[colIndex] && cols[colIndex].chars[charIndex]) {
          cols[colIndex].chars[charIndex].highlighted = true
          cols[colIndex].chars[charIndex].highlightColor = ann.color
        }
      }
    }
  })

  return cols
})

const selectionOverlayStyle = computed(() => {
  if (!selectionStart.value || !selectionEnd.value) return {}
  return {
    position: 'absolute',
    background: 'rgba(139, 69, 19, 0.2)',
    pointerEvents: 'none',
    borderRadius: '2px'
  }
})

function isCharSelected(colIndex, charIndex) {
  if (!selectionStart.value || !selectionEnd.value) return false
  const pos = colIndex * props.layout.charsPerColumn + charIndex
  const startPos = Math.min(selectionStart.value, selectionEnd.value)
  const endPos = Math.max(selectionStart.value, selectionEnd.value)
  return pos >= startPos && pos < endPos
}

function getCharStyle(char) {
  const base = {
    width: props.layout.fontSize * 1.2 + 'px',
    height: props.layout.fontSize * props.layout.lineHeight + 'px',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    position: 'relative',
    cursor: props.editable ? 'text' : 'default',
    userSelect: 'none'
  }
  if (char.highlighted && char.highlightColor) {
    base.backgroundColor = char.highlightColor + '4D'
  }
  return base
}

function handleMouseDown(e) {
  if (!props.editable) return
  const target = e.target.closest('.char-item')
  if (target) {
    const pos = parseInt(target.dataset.pos)
    selectionStart.value = pos
    selectionEnd.value = pos + 1
    isSelecting.value = true
    selectionActive.value = true
  }
}

function handleMouseMove(e) {
  if (!isSelecting.value || !props.editable) return
  const target = document.elementFromPoint(e.clientX, e.clientY)
  const charItem = target?.closest('.char-item')
  if (charItem) {
    const pos = parseInt(charItem.dataset.pos)
    selectionEnd.value = pos + 1
  }
}

function handleMouseUp() {
  if (!isSelecting.value) return
  isSelecting.value = false

  if (selectionStart.value !== null && selectionEnd.value !== null) {
    const start = Math.min(selectionStart.value, selectionEnd.value)
    const end = Math.max(selectionStart.value, selectionEnd.value)
    if (end > start) {
      emit('selection', {
        startPos: start,
        endPos: end,
        text: props.text.substring(start, end)
      })
    }
  }

  setTimeout(() => {
    selectionActive.value = false
    selectionStart.value = null
    selectionEnd.value = null
  }, 200)
}

function clearSelection() {
  selectionActive.value = false
  selectionStart.value = null
  selectionEnd.value = null
}

defineExpose({ clearSelection })

onMounted(() => {
  nextTick(() => {})
})

watch(() => props.layout, () => {
  nextTick(() => {})
}, { deep: true })
</script>

<style scoped>
.vertical-text-container {
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(139, 69, 19, 0.1);
}

.vertical-text-content {
  background: linear-gradient(to right, transparent 0%, rgba(139, 69, 19, 0.02) 50%, transparent 100%);
}

.show-border {
  box-shadow: inset 0 0 20px rgba(139, 69, 19, 0.05);
}

.vertical-column {
  position: relative;
}

.char-item {
  transition: background-color 0.15s;
}

.char-item:hover {
  background-color: rgba(139, 69, 19, 0.08);
}

.char-item.selected {
  background-color: rgba(139, 69, 19, 0.2);
}

.char-item.highlighted {
  background-color: rgba(255, 215, 0, 0.3);
}

.punctuation {
  position: absolute;
  font-size: 0.7em;
  color: #8b4513;
  font-weight: bold;
}

.punctuation.period,
.punctuation.comma {
  right: -2px;
  bottom: 0.2em;
}

.punctuation.question,
.punctuation.exclaim {
  right: -2px;
  bottom: 0.2em;
}
</style>
