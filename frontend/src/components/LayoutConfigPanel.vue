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
            max="100"
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
            max="48"
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
            max="4"
            step="0.1"
            @change="emitChange"
          />
        </div>
        <div class="form-group">
          <label class="form-label">字距 (px)</label>
          <input
            type="number"
            v-model.number="localConfig.char_spacing"
            class="form-input"
            min="-5"
            max="30"
            step="1"
            @change="emitChange"
          />
          <div class="form-hint">调整列内字符间的垂直间距</div>
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
        <div class="section-title">系统预设</div>
        <div class="preset-buttons">
          <button
            v-for="preset in systemPresets"
            :key="preset.name"
            class="preset-btn"
            @click="applySystemPreset(preset)"
          >
            {{ preset.name }}
          </button>
        </div>
      </div>

      <div class="config-section">
        <div class="section-title">
          <span>我的预设</span>
          <button class="text-btn" @click="showSaveDialog = true" title="保存为预设">
            + 新建
          </button>
        </div>
        <div v-if="loadingPresets" class="preset-loading">加载中...</div>
        <div v-else-if="userPresets.length === 0" class="preset-empty">
          暂无个人预设
        </div>
        <div v-else class="preset-list">
          <div
            v-for="preset in userPresets"
            :key="preset.id"
            class="preset-item"
            :class="{ 'is-default': preset.is_default }"
          >
            <div class="preset-info" @click="applyUserPreset(preset)">
              <div class="preset-name">
                {{ preset.name }}
                <span v-if="preset.is_default" class="default-badge">默认</span>
              </div>
              <div class="preset-desc" v-if="preset.description">
                {{ preset.description }}
              </div>
            </div>
            <div class="preset-actions">
              <button
                v-if="!preset.is_default"
                class="icon-btn"
                title="设为默认"
                @click.stop="handleSetDefault(preset.id)"
              >
                ☆
              </button>
              <button
                class="icon-btn delete-btn"
                title="删除"
                @click.stop="deleteUserPreset(preset.id)"
              >
                ×
              </button>
            </div>
          </div>
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

    <div v-if="showSaveDialog" class="save-dialog-overlay" @click.self="showSaveDialog = false">
      <div class="save-dialog card">
        <h4 class="dialog-title">保存为预设</h4>
        <div class="form-group">
          <label class="form-label">预设名称</label>
          <input
            type="text"
            v-model="saveDialog.name"
            class="form-input"
            placeholder="给预设起个名字"
            @keyup.enter="confirmSavePreset"
          />
        </div>
        <div class="form-group">
          <label class="form-label">描述 (可选)</label>
          <textarea
            v-model="saveDialog.description"
            class="form-input"
            rows="2"
            placeholder="简单描述一下这个预设"
          ></textarea>
        </div>
        <div class="form-group">
          <label class="form-checkbox">
            <input type="checkbox" v-model="saveDialog.isDefault" />
            <span>设为默认预设</span>
          </label>
        </div>
        <div class="dialog-actions">
          <button class="btn btn-secondary" @click="showSaveDialog = false">
            取消
          </button>
          <button class="btn btn-primary" @click="confirmSavePreset" :disabled="!saveDialog.name.trim()">
            保存
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, watch, onMounted } from 'vue'
import {
  listPresets,
  createPreset,
  deletePreset,
  setDefaultPreset as setDefaultPresetApi,
  presetToLayout
} from '../api/layoutPreset'

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
  char_spacing: 0,
  text_color: '#1a1a1a',
  background_color: '#f5f0e8',
  show_border: true,
  custom_css: ''
}

const localConfig = reactive({ ...defaultConfig, ...props.layout })
const userPresets = ref([])
const loadingPresets = ref(false)
const showSaveDialog = ref(false)
const saveDialog = reactive({
  name: '',
  description: '',
  isDefault: false
})

const systemPresets = [
  {
    name: '经典古籍',
    config: {
      chars_per_column: 20,
      column_gap: 24,
      font_size: 18,
      font_family: 'SimSun, serif',
      line_height: 1.8,
      char_spacing: 0,
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
      char_spacing: 2,
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
      char_spacing: 4,
      text_color: '#000000',
      background_color: '#fff9e6',
      show_border: true
    }
  }
]

function emitChange() {
  emit('change', { ...localConfig })
}

function applySystemPreset(preset) {
  Object.assign(localConfig, defaultConfig, preset.config)
  emitChange()
}

function applyUserPreset(preset) {
  const layout = presetToLayout(preset)
  if (layout) {
    Object.assign(localConfig, defaultConfig, layout)
    emitChange()
  }
}

function resetConfig() {
  Object.assign(localConfig, defaultConfig)
  emitChange()
}

function saveConfig() {
  emit('save', { ...localConfig })
}

async function loadUserPresets() {
  loadingPresets.value = true
  try {
    const res = await listPresets()
    if (res && res.data) {
      userPresets.value = res.data
    }
  } catch (e) {
    console.error('Failed to load presets:', e)
  } finally {
    loadingPresets.value = false
  }
}

async function confirmSavePreset() {
  if (!saveDialog.name.trim()) return
  try {
    const data = {
      name: saveDialog.name.trim(),
      description: saveDialog.description,
      chars_per_column: localConfig.chars_per_column,
      column_gap: localConfig.column_gap,
      font_size: localConfig.font_size,
      font_family: localConfig.font_family,
      line_height: localConfig.line_height,
      char_spacing: localConfig.char_spacing,
      text_color: localConfig.text_color,
      background_color: localConfig.background_color,
      show_border: localConfig.show_border,
      is_default: saveDialog.isDefault,
      custom_css: localConfig.custom_css
    }
    await createPreset(data)
    showSaveDialog.value = false
    saveDialog.name = ''
    saveDialog.description = ''
    saveDialog.isDefault = false
    await loadUserPresets()
  } catch (e) {
    console.error('Failed to create preset:', e)
    alert('保存失败，请稍后重试')
  }
}

async function handleSetDefault(id) {
  try {
    await setDefaultPresetApi(id)
    await loadUserPresets()
  } catch (e) {
    console.error('Failed to set default preset:', e)
  }
}

async function deleteUserPreset(id) {
  if (!confirm('确定要删除这个预设吗？')) return
  try {
    await deletePreset(id)
    await loadUserPresets()
  } catch (e) {
    console.error('Failed to delete preset:', e)
  }
}

watch(() => props.layout, (newVal) => {
  if (newVal) {
    Object.assign(localConfig, defaultConfig, newVal)
  }
}, { deep: true })

onMounted(() => {
  loadUserPresets()
})

defineExpose({ loadUserPresets })
</script>

<style scoped>
.layout-config-panel {
  width: 300px;
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  position: relative;
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
  max-height: calc(100vh - 320px);
  overflow-y: auto;
  padding-right: 4px;
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
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.text-btn {
  background: none;
  border: none;
  color: #8b4513;
  font-size: 12px;
  cursor: pointer;
  padding: 2px 6px;
  border-radius: 4px;
}

.text-btn:hover {
  background: #f5ebe0;
}

.form-group {
  margin-bottom: 0;
}

.form-hint {
  font-size: 11px;
  color: #999;
  margin-top: 4px;
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

.preset-loading,
.preset-empty {
  font-size: 13px;
  color: #999;
  padding: 12px 0;
  text-align: center;
}

.preset-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.preset-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 10px;
  background: #f8f6f3;
  border: 1px solid #e8e0d4;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
}

.preset-item:hover {
  background: #f0ebe3;
  border-color: #d4c8b8;
}

.preset-item.is-default {
  border-color: #c9a96e;
  background: #fdf8ef;
}

.preset-info {
  flex: 1;
  min-width: 0;
}

.preset-name {
  font-size: 13px;
  color: #5a4030;
  font-weight: 500;
  display: flex;
  align-items: center;
  gap: 6px;
}

.default-badge {
  font-size: 10px;
  background: #c9a96e;
  color: #fff;
  padding: 1px 5px;
  border-radius: 3px;
  font-weight: normal;
}

.preset-desc {
  font-size: 11px;
  color: #999;
  margin-top: 2px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.preset-actions {
  display: flex;
  gap: 4px;
  margin-left: 8px;
}

.icon-btn {
  width: 24px;
  height: 24px;
  border: none;
  background: transparent;
  cursor: pointer;
  border-radius: 4px;
  font-size: 14px;
  color: #999;
  display: flex;
  align-items: center;
  justify-content: center;
}

.icon-btn:hover {
  background: #e8e0d4;
  color: #5a4030;
}

.icon-btn.delete-btn:hover {
  background: #ffe0e0;
  color: #c0392b;
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

.save-dialog-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.save-dialog {
  width: 340px;
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.dialog-title {
  font-size: 16px;
  font-weight: 600;
  color: #5a4030;
  margin: 0;
}

.dialog-actions {
  display: flex;
  gap: 10px;
  justify-content: flex-end;
  margin-top: 8px;
}

.dialog-actions .btn {
  padding: 8px 20px;
}
</style>
