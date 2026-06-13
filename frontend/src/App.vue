<template>
  <div id="app-root">
    <header class="app-header">
      <div class="header-content">
        <h1 class="app-title" @click="$router.push('/')">古籍文本版式还原与断句标注系统</h1>
        <nav class="app-nav">
          <router-link to="/" class="nav-link">文本列表</router-link>
          <router-link to="/texts/new" class="nav-link">新增文本</router-link>
        </nav>
      </div>
    </header>
    <main class="app-main">
      <router-view />
    </main>
  </div>
</template>

<script setup>
</script>

<style scoped>
.app-header {
  background: linear-gradient(135deg, #8b4513 0%, #a0522d 100%);
  color: #f5f0e8;
  padding: 0;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.header-content {
  max-width: 1600px;
  margin: 0 auto;
  padding: 16px 32px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.app-title {
  margin: 0;
  font-size: 22px;
  font-weight: 600;
  cursor: pointer;
  letter-spacing: 2px;
}

.app-nav {
  display: flex;
  gap: 24px;
}

.nav-link {
  color: #f5f0e8;
  text-decoration: none;
  padding: 8px 16px;
  border-radius: 4px;
  transition: all 0.2s;
  font-size: 15px;
}

.nav-link:hover {
  background: rgba(255, 255, 255, 0.15);
}

.nav-link.router-link-active {
  background: rgba(255, 255, 255, 0.25);
  font-weight: 500;
}

.app-main {
  max-width: 1600px;
  margin: 0 auto;
  padding: 24px 32px;
  min-height: calc(100vh - 80px);
}
</style>
