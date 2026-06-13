<template>
  <div class="text-list-page">
    <div class="page-header">
      <h2 class="page-title">古籍文本列表</h2>
      <div class="header-actions">
        <div class="search-box">
          <input
            type="text"
            v-model="keyword"
            class="form-input"
            placeholder="搜索标题、作者..."
            @keyup.enter="searchTexts"
          />
          <button class="btn btn-primary btn-sm" @click="searchTexts">搜索</button>
        </div>
        <button class="btn btn-primary" @click="goToCreate">
          + 新增文本
        </button>
      </div>
    </div>

    <div v-if="textStore.loading && textStore.texts.length === 0" class="loading-state">
      <div class="loading-spinner"></div>
      <p>加载中...</p>
    </div>

    <div v-else-if="textStore.texts.length === 0" class="empty-state">
      <div class="empty-icon">📜</div>
      <h3>暂无古籍文本</h3>
      <p>点击右上角按钮添加您的第一本古籍</p>
      <button class="btn btn-primary" @click="goToCreate">
        立即添加
      </button>
    </div>

    <div v-else class="text-grid">
      <div
        v-for="text in textStore.texts"
        :key="text.id"
        class="text-card"
        @click="goToDetail(text.id)"
      >
        <div class="card-header">
          <h3 class="card-title">{{ text.title }}</h3>
          <span v-if="text.dynasty" class="card-dynasty">{{ text.dynasty }}</span>
        </div>
        <div class="card-author" v-if="text.author">
          作者：{{ text.author }}
        </div>
        <div class="card-preview">
          {{ truncateText(text.content, 80) }}
        </div>
        <div class="card-footer">
          <span class="char-count">{{ text.content?.length || 0 }} 字</span>
          <span class="card-date">{{ formatDate(text.created_at) }}</span>
        </div>
        <div class="card-actions" @click.stop>
          <button class="btn btn-secondary btn-sm" @click="goToEdit(text.id)">
            编辑
          </button>
          <button class="btn btn-danger btn-sm" @click="handleDelete(text.id)">
            删除
          </button>
        </div>
      </div>
    </div>

    <div v-if="textStore.total > pageSize" class="pagination">
      <button
        class="page-btn"
        :disabled="page <= 1"
        @click="changePage(page - 1)"
      >
        上一页
      </button>
      <span class="page-info">
        第 {{ page }} / {{ totalPages }} 页，共 {{ textStore.total }} 条
      </span>
      <button
        class="page-btn"
        :disabled="page >= totalPages"
        @click="changePage(page + 1)"
      >
        下一页
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useTextStore } from '@/stores/text'

const router = useRouter()
const textStore = useTextStore()

const keyword = ref('')
const page = ref(1)
const pageSize = ref(12)

const totalPages = computed(() => {
  return Math.ceil(textStore.total / pageSize.value) || 1
})

async function loadTexts() {
  const params = {
    page: page.value,
    page_size: pageSize.value
  }
  if (keyword.value) {
    params.keyword = keyword.value
  }
  await textStore.fetchTexts(params)
}

function searchTexts() {
  page.value = 1
  loadTexts()
}

function changePage(p) {
  page.value = p
  loadTexts()
}

function goToCreate() {
  router.push('/texts/new')
}

function goToDetail(id) {
  router.push(`/texts/${id}`)
}

function goToEdit(id) {
  router.push(`/texts/${id}/edit`)
}

async function handleDelete(id) {
  if (confirm('确定要删除这篇古籍文本吗？')) {
    await textStore.deleteTextData(id)
    if (textStore.texts.length === 0 && page.value > 1) {
      page.value--
    }
    loadTexts()
  }
}

function truncateText(text, len) {
  if (!text) return ''
  if (text.length <= len) return text
  return text.substring(0, len) + '...'
}

function formatDate(dateStr) {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN')
}

onMounted(() => {
  loadTexts()
})
</script>

<style scoped>
.text-list-page {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  color: #5a4030;
  margin: 0;
}

.header-actions {
  display: flex;
  gap: 16px;
  align-items: center;
}

.search-box {
  display: flex;
  gap: 8px;
}

.search-box .form-input {
  width: 240px;
}

.loading-state,
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 20px;
  gap: 16px;
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

.empty-icon {
  font-size: 64px;
}

.empty-state h3 {
  font-size: 18px;
  color: #5a4030;
  margin: 0;
}

.empty-state p {
  color: #8b7355;
  margin: 0;
}

.text-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
}

.text-card {
  background: #fff;
  border-radius: 10px;
  padding: 20px;
  box-shadow: 0 2px 12px rgba(139, 69, 19, 0.08);
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  flex-direction: column;
  gap: 12px;
  position: relative;
}

.text-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 20px rgba(139, 69, 19, 0.15);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 10px;
}

.card-title {
  font-size: 18px;
  font-weight: 600;
  color: #5a4030;
  margin: 0;
  line-height: 1.4;
}

.card-dynasty {
  font-size: 12px;
  color: #8b7355;
  background: #f5f0e8;
  padding: 3px 8px;
  border-radius: 4px;
  white-space: nowrap;
}

.card-author {
  font-size: 14px;
  color: #8b7355;
}

.card-preview {
  font-size: 13px;
  color: #6b5848;
  line-height: 1.6;
  flex: 1;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.card-footer {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: #a08970;
  padding-top: 10px;
  border-top: 1px solid #f0ebe3;
}

.card-actions {
  display: flex;
  gap: 8px;
  opacity: 0;
  transition: opacity 0.2s;
}

.text-card:hover .card-actions {
  opacity: 1;
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 16px;
  padding: 20px 0;
}

.page-btn {
  padding: 8px 20px;
  background: #fff;
  border: 1px solid #d4c8b8;
  border-radius: 6px;
  cursor: pointer;
  color: #5a4030;
  font-size: 14px;
  transition: all 0.2s;
}

.page-btn:hover:not(:disabled) {
  background: #f5f0e8;
  border-color: #8b4513;
}

.page-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.page-info {
  font-size: 14px;
  color: #8b7355;
}
</style>
