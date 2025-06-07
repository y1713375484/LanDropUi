<template>
  <div id="app">
    <header class="app-header">
      <h1 class="app-title">局域网文件传输-接收端</h1>
    </header>

    <main class="app-content">
      <!-- 路由匹配的组件将渲染在这里 -->
      <router-view></router-view>
    </main>

    <footer class="app-footer">
      <p></p>
    </footer>
  </div>
</template>

<script>
export default {
  methods: {
    goToPage(name) {
      if (this.$route.name !== name) {
        this.$router.push({ name })
      }
    }
  }
}
</script>

<style>
/* 全局样式 */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
  background-color: #f5f7fa;
}

#app {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

/* 头部样式 */
.app-header {
  background-color: #3498db;
  color: white;
  padding: 1rem 2rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.app-title {
  font-size: 1.5rem;
  font-weight: bold;
}

.app-nav {
  display: flex;
  gap: 1rem;
}

.nav-button {
  background-color: transparent;
  border: 1px solid white;
  color: white;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s, transform 0.2s;
}

.nav-button:hover {
  background-color: rgba(255, 255, 255, 0.2);
  transform: translateY(-2px);
}

/* 内容区域样式 */
.app-content {
  flex: 1;
  padding: 2rem;
  max-width: 1200px;
  width: 100%;
  margin: 0 auto;
}

/* 页脚样式 */
.app-footer {
  background-color: #2c3e50;
  color: #ecf0f1;
  text-align: center;
  padding: 1rem;
  margin-top: auto;
}
</style>