<template>
  <div class="text-editor-page">
    <div class="page-header">
      <button class="btn btn-secondary" @click="goBack">
        ← 返回列表
      </button>
      <h2 class="page-title">{{ isEdit ? '编辑古籍文本' : '新增古籍文本' }}</h2>
      <div class="header-actions">
        <button class="btn btn-secondary" @click="resetForm">
          重置
        </button>
        <button class="btn btn-primary" @click="saveText" :disabled="saving">
          {{ saving ? '保存中...' : '保存' }}
        </button>
      </div>
    </div>

    <div class="editor-container">
      <div class="form-section card">
        <h3 class="section-title">基本信息</h3>

        <div class="form-group">
          <label class="form-label">标题 <span class="required">*</span></label>
          <input
            type="text"
            v-model="form.title"
            class="form-input"
            placeholder="请输入古籍标题"
          />
        </div>

        <div class="form-row">
          <div class="form-group">
            <label class="form-label">作者</label>
            <input
              type="text"
              v-model="form.author"
              class="form-input"
              placeholder="请输入作者"
            />
          </div>
          <div class="form-group">
            <label class="form-label">朝代</label>
            <input
              type="text"
              v-model="form.dynasty"
              class="form-input"
              placeholder="如：宋、唐、明"
            />
          </div>
        </div>

        <div class="form-group">
          <label class="form-label">描述</label>
          <textarea
            v-model="form.description"
            class="form-textarea"
            rows="3"
            placeholder="请输入古籍简介、出处等描述信息"
          ></textarea>
        </div>
      </div>

      <div class="form-section card">
        <h3 class="section-title">文本内容 <span class="required">*</span></h3>
        <textarea
          v-model="form.content"
          class="form-textarea content-textarea"
          placeholder="请输入或粘贴古籍文本内容..."
          rows="15"
        ></textarea>
        <div class="char-count">
          当前字数：<strong>{{ form.content?.length || 0 }}</strong> 字
        </div>
      </div>

      <div class="form-section card" v-if="!isEdit">
        <h3 class="section-title">版式预设</h3>
        <div class="preset-options">
          <label class="preset-option" v-for="preset in layoutPresets" :key="preset.name">
            <input
              type="radio"
              v-model="selectedPreset"
              :value="preset.name"
              @change="applyPreset(preset)"
            />
            <span class="preset-name">{{ preset.name }}</span>
            <span class="preset-desc">{{ preset.desc }}</span>
          </label>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useTextStore } from '@/stores/text'

const router = useRouter()
const route = useRoute()
const textStore = useTextStore()

const saving = ref(false)
const selectedPreset = ref('classic')

const isEdit = computed(() => !!route.params.id)

const form = reactive({
  title: '',
  author: '',
  dynasty: '',
  content: '',
  description: '',
  layout_config: null
})

const layoutPresets = [
  {
    name: 'classic',
    desc: '经典古籍版式，20字/列',
    config: {
      columns: 12,
      chars_per_column: 20,
      column_gap: 24,
      font_size: 18,
      font_family: 'SimSun, serif',
      line_height: 1.8,
      text_color: '#1a1a1a',
      background_color: '#f5f0e8',
      show_border: true
    }
  },
  {
    name: 'modern',
    desc: '现代阅读版式，25字/列',
    config: {
      columns: 10,
      chars_per_column: 25,
      column_gap: 30,
      font_size: 16,
      font_family: 'KaiTi, serif',
      line_height: 2.0,
      text_color: '#333333',
      background_color: '#faf8f5',
      show_border: false
    }
  },
  {
    name: 'large',
    desc: '大字版，15字/列',
    config: {
      columns: 8,
      chars_per_column: 15,
      column_gap: 28,
      font_size: 24,
      font_family: 'SimSun, serif',
      line_height: 1.8,
      text_color: '#000000',
      background_color: '#fff9e6',
      show_border: true
    }
  }
]

function applyPreset(preset) {
  form.layout_config = { ...preset.config }
}

async function loadText(id) {
  try {
    const text = await textStore.fetchText(id)
    form.title = text.title || ''
    form.author = text.author || ''
    form.dynasty = text.dynasty || ''
    form.content = text.content || ''
    form.description = text.description || ''
  } catch (e) {
    console.error('Failed to load text:', e)
    alert('加载文本失败')
  }
}

async function saveText() {
  if (!form.title.trim()) {
    alert('请输入标题')
    return
  }
  if (!form.content.trim()) {
    alert('请输入文本内容')
    return
  }

  saving.value = true
  try {
    if (isEdit.value) {
      await textStore.updateTextData(route.params.id, {
        title: form.title,
        author: form.author,
        dynasty: form.dynasty,
        content: form.content,
        description: form.description
      })
    } else {
      const data = {
        title: form.title,
        author: form.author,
        dynasty: form.dynasty,
        content: form.content,
        description: form.description,
        layout_config: form.layout_config
      }
      const text = await textStore.createNewText(data)
      if (text?.id) {
        router.push(`/texts/${text.id}`)
        return
      }
    }
    router.push('/')
  } catch (e) {
    console.error('Save failed:', e)
    alert('保存失败：' + (e.message || '未知错误'))
  } finally {
    saving.value = false
  }
}

function resetForm() {
  if (isEdit.value && textStore.currentText) {
    const t = textStore.currentText
    form.title = t.title
    form.author = t.author
    form.dynasty = t.dynasty
    form.content = t.content
    form.description = t.description
  } else {
    form.title = ''
    form.author = ''
    form.dynasty = ''
    form.content = ''
    form.description = ''
    form.layout_config = null
    selectedPreset.value = 'classic'
    applyPreset(layoutPresets[0])
  }
}

function goBack() {
  router.back()
}

onMounted(() => {
  if (isEdit.value) {
    loadText(route.params.id)
  } else {
    applyPreset(layoutPresets[0])
  }
})
</script>

<style scoped>
.text-editor-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.page-header {
  display: flex;
  align-items: center;
  gap: 16px;
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  color: #5a4030;
  margin: 0;
  flex: 1;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.editor-container {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-section {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.section-title {
  font-size: 18px;
  font-weight: 600;
  color: #5a4030;
  margin: 0;
  padding-bottom: 12px;
  border-bottom: 1px solid #e8e0d4;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

.required {
  color: #c0392b;
}

.content-textarea {
  font-family: 'SimSun', serif;
  font-size: 14px;
  line-height: 1.8;
}

.char-count {
  text-align: right;
  font-size: 13px;
  color: #8b7355;
}

.char-count strong {
  color: #8b4513;
  font-size: 16px;
}

.preset-options {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
}

.preset-option {
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding: 16px;
  background: #f8f6f3;
  border: 2px solid #e8e0d4;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
  position: relative;
}

.preset-option input[type="radio"] {
  position: absolute;
  top: 12px;
  right: 12px;
}

.preset-option:hover {
  border-color: #d4c8b8;
  background: #f0ebe3;
}

.preset-option input[type="radio"]:checked + .preset-name,
.preset-option:has(input:checked) {
  border-color: #8b4513;
  background: #fff8f0;
}

.preset-name {
  font-size: 14px;
  font-weight: 600;
  color: #5a4030;
}

.preset-desc {
  font-size: 12px;
  color: #8b7355;
}
</style>
