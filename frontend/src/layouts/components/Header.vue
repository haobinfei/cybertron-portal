<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useAppStore } from '@/stores/app'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const appStore = useAppStore()
const userStore = useUserStore()

async function handleLogout() {
  await userStore.logout()
  router.push('/login')
}
</script>

<template>
  <div class="header">
    <div class="header-left">
      <el-icon class="collapse-btn" @click="appStore.toggleSidebar">
        <Fold v-if="!appStore.sidebarCollapsed" />
        <Expand v-else />
      </el-icon>
    </div>
    <div class="header-right">
      <el-dropdown>
        <span class="user-info">
          <el-avatar :size="32" />
          <span class="username">{{ userStore.userInfo?.nickname || userStore.userInfo?.username || '管理员' }}</span>
        </span>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item @click="handleLogout">退出登录</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </div>
</template>

<script lang="ts">
import { Expand, Fold } from '@element-plus/icons-vue'
export default { components: { Expand, Fold } }
</script>

<style scoped>
.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 100%;
  padding: 0 16px;
  background: var(--tech-bg-card);
  backdrop-filter: blur(14px);
  -webkit-backdrop-filter: blur(14px);
  border-bottom: 1px solid var(--tech-border);
}
.collapse-btn {
  font-size: 20px;
  cursor: pointer;
  color: var(--tech-text-dim);
  transition: color 0.2s;
}
.collapse-btn:hover {
  color: var(--tech-primary);
  text-shadow: 0 0 8px rgba(0, 240, 255, 0.6);
}
.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  color: var(--tech-text);
}
</style>
