<script setup lang="ts">
import { useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'

const route = useRoute()
const userStore = useUserStore()
</script>

<template>
  <div class="sidebar">
    <div class="logo">
      <span class="tech-glow-text">CYBERTRON</span>
    </div>
    <el-menu
      :default-active="route.path"
      router
      background-color="transparent"
      text-color="#94a3b8"
      active-text-color="#00f0ff"
    >
      <el-menu-item index="/dashboard">
        <el-icon><Monitor /></el-icon>
        <span>仪表盘</span>
      </el-menu-item>
      <el-menu-item v-if="userStore.isAdmin" index="/users">
        <el-icon><User /></el-icon>
        <span>用户管理</span>
      </el-menu-item>
    </el-menu>
  </div>
</template>

<script lang="ts">
import { Monitor, User } from '@element-plus/icons-vue'
export default { components: { Monitor, User } }
</script>

<style scoped>
.sidebar {
  height: 100%;
  background: linear-gradient(180deg, #0a0e27 0%, #0d1230 100%);
  border-right: 1px solid var(--tech-border);
}
.logo {
  height: 60px;
  line-height: 60px;
  text-align: center;
  font-size: 18px;
  font-weight: 700;
  letter-spacing: 3px;
  border-bottom: 1px solid var(--tech-border);
}
.el-menu {
  border-right: none;
}
.el-menu :deep(.el-menu-item) {
  margin: 4px 8px;
  border-radius: 6px;
  border: 1px solid transparent;
}
.el-menu :deep(.el-menu-item:hover) {
  background-color: var(--tech-primary-dim);
  color: var(--tech-primary);
}
.el-menu :deep(.el-menu-item.is-active) {
  background-color: var(--tech-primary-dim);
  border-color: var(--tech-border-strong);
  box-shadow: var(--tech-glow);
}
</style>
