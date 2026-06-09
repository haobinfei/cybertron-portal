<script setup lang="ts">
import { useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'

const route = useRoute()
const userStore = useUserStore()
</script>

<template>
  <div class="sidebar">
    <div class="logo">
      <span>Cybertron</span>
    </div>
    <el-menu
      :default-active="route.path"
      router
      background-color="#304156"
      text-color="#bfcbd9"
      active-text-color="#409eff"
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
  background-color: #304156;
}
.logo {
  height: 60px;
  line-height: 60px;
  text-align: center;
  color: #fff;
  font-size: 18px;
  font-weight: bold;
}
</style>
