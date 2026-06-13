<template>
  <div class="layout-config-panel card">
    <div class="panel-header">
      <h3 class="panel-title">版式配置</h3>
    </div>

    <div class="config-sections">
      <div class="config-section">
        <div class="section-title">列设置</div>
        <div class="form-group">
          <label class="form-label">每列字数</label>
          <input
            type="number"
            v-model.number="localConfig.chars_per_column"
            class="form-input"
            min="5"
            max="50"
            @change="emitChange"
          />
        </div>
        <div class="form-group">
          <label class="form-label">列间距 (px)</label>
          <input
            type="number"
            v-model.number="localConfig.column_gap"
            class="form-input"
            min="5"
            max="60"
            @change="emitChange"
          />
        </div>
      </div>

      <div class="config-section">
        <div class="section-title">字体设置</div>
        <div class="form-group">
          <label class="form-label">字体大小 (px)</label>
          <input
            type="number"
            v-model.number="localConfig.font_size"
            class="form-input"
            min="12"
            max="36"
            @change="emitChange"
          />
        </div>
        <div class="form-group">
          <label class="form-label">字体</label>
          <select v-model="localConfig.font_family" class="form-select" @change="emitChange">
            <option value="SimSun, serif">宋体</option>
            <option value="KaiTi, serif">楷体</option>
            <option value="FangSong, serif">仿宋</option>
            <option value="STSong, serif">华文宋体</option>
            <option value="STKaiti, serif">华文楷体</option>
            <option value="serif">衬线字体</option>
          </select>
        </div>
        <div class="form-group">
          <label class="form-label">行高</label>
          <input
            type="number"
            v-model.number="localConfig.line_height"
            class="form-input"
            min="1"
            max="3"
            step="0.1"
            @change="emitChange"
          />
        </div>
      </div>

      <div class="config-section">
        <div class="section-title">颜色设置</div>
        <div class="form-group">
          <label class="form-label">文字颜色</label>
          <div class="color-input-wrapper">
            <input
              type="color"
              v-model="localConfig.text_color"
              class="color-picker-input"
              @input="emitChange"
            />
            <input
              type="text"
              v-model="localConfig.text_color"
              class="form-input color-text-input"
              @change="emitChange"
            />
          </div>
        </div>
        <div class="form-group">
          <label class="form-label">背景颜色</label>
          <div class="color-input-wrapper">
            <input
              type="color"
              v-model="localConfig.background_color"
              class="color-picker-input"
              @input="emitChange"
            />
            <input
              type="text"
              v-model="localConfig.background_color"
              class="form-input color-text-input"
              @change="emitChange"
            />
          </div>
        </div>
      </div>

      <div class="config-section">
        <div class="section-title">其他设置</div>
        <div class="form-group">
          <label class="form-checkbox">
            <input
              type="checkbox"
              v-model="localConfig.show_border"
              @change="emitChange"
            />
            <span>显示边框</span>
          </label>
        </div>
      </div>

      <div class="config-section">
        <div class="section-title">预设模板</div>
        <div class="preset-buttons">
          <button
            v-for="preset in presets"
            :key="preset.name"
            class="preset-btn"
            @click="applyPreset(preset)"
          >
            {{ preset.name }}
          </button>
        </div>
      </div>
    </div>

    <div class="panel-actions">
      <button class="btn btn-secondary" @click="resetConfig">
        重置
      </button>
      <button class="btn btn-primary" @click="saveConfig">
        保存配置
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, watch, computed } from 'vue'

const props = defineProps({
  layout: {
    type: Object,
    default: () => ({})
  }
})

const emit = defineEmits(['change', 'save'])

const defaultConfig = {
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

const localConfig = reactive({ ...defaultConfig, ...props.layout })

const presets = [
  {
    name: '经典古籍',
    config: {
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
    name: '现代阅读',
    config: {
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
    name: '大字版',
    config: {
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

function emitChange() {
  emit('change', { ...localConfig })
}

function applyPreset(preset) {
  Object.assign(localConfig, preset.config)
  emitChange()
}

function resetConfig() {
  Object.assign(localConfig, defaultConfig)
  emitChange()
}

function saveConfig() {
  emit('save', { ...localConfig })
}

watch(() => props.layout, (newVal) => {
  if (newVal) {
    Object.assign(localConfig, defaultConfig, newVal)
  }
}, { deep: true })
</script>

<style scoped>
.layout-config-panel {
  width: 280px;
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.panel-header {
  padding-bottom: 12px;
  border-bottom: 1px solid #e8e0d4;
}

.panel-title {
  font-size: 16px;
  font-weight: 600;
  color: #5a4030;
  margin: 0;
}

.config-sections {
  display: flex;
  flex-direction: column;
  gap: 20px;
  max-height: calc(100vh - 300px);
  overflow-y: auto;
}

.config-section {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.section-title {
  font-size: 13px;
  font-weight: 600;
  color: #8b4513;
  padding-bottom: 6px;
  border-bottom: 1px dashed #e8e0d4;
}

.form-group {
  margin-bottom: 0;
}

.color-input-wrapper {
  display: flex;
  gap: 8px;
  align-items: center;
}

.color-picker-input {
  width: 40px;
  height: 38px;
  border: 1px solid #d4c8b8;
  border-radius: 6px;
  cursor: pointer;
  padding: 2px;
}

.color-text-input {
  flex: 1;
}

.form-checkbox {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  font-size: 14px;
  color: #5a4030;
}

.form-checkbox input {
  width: 16px;
  height: 16px;
  cursor: pointer;
}

.preset-buttons {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 8px;
}

.preset-btn {
  padding: 10px 12px;
  background: #f8f6f3;
  border: 1px solid #e8e0d4;
  border-radius: 6px;
  cursor: pointer;
  font-size: 13px;
  color: #5a4030;
  transition: all 0.2s;
}

.preset-btn:hover {
  background: #f0ebe3;
  border-color: #d4c8b8;
}

.panel-actions {
  display: flex;
  gap: 10px;
  padding-top: 16px;
  border-top: 1px solid #e8e0d4;
}

.panel-actions .btn {
  flex: 1;
}
</style>
