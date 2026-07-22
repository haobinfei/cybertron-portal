<template>
  <div class="not-found tech-grid-bg">
    <h1 class="tech-glow-text">404</h1>
    <p>页面不存在</p>
    <el-button type="primary" @click="$router.push('/')">返回首页</el-button>
  </div>
</template>

<style scoped>
.not-found {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100vh;
  gap: 12px;
}
.not-found h1 {
  position: relative;
  z-index: 1;
  font-size: 96px;
  margin: 0;
  font-weight: 700;
  letter-spacing: 8px;
}
.not-found p {
  position: relative;
  z-index: 1;
  color: var(--tech-text-dim);
  margin: 0 0 12px;
}
.not-found .el-button {
  position: relative;
  z-index: 1;
}
</style>
