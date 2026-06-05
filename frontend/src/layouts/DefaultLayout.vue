<script setup lang="ts">
import Sidebar from './components/Sidebar.vue'
import Header from './components/Header.vue'
</script>

<template>
  <el-container class="layout-container">
    <el-aside :width="sidebarWidth">
      <Sidebar />
    </el-aside>
    <el-container>
      <el-header height="60px">
        <Header />
      </el-header>
      <el-main>
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script lang="ts">
import { computed } from 'vue'
import { useAppStore } from '@/stores/app'

export default {
  setup() {
    const appStore = useAppStore()
    const sidebarWidth = computed(() =>
      appStore.sidebarCollapsed ? '64px' : '200px',
    )
    return { sidebarWidth }
  },
}
</script>

<style scoped>
.layout-container {
  height: 100vh;
}
.el-header {
  padding: 0;
}
</style>
