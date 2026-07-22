<script setup lang="ts">
import { Monitor, User, Connection, Cpu } from '@element-plus/icons-vue'

const stats = [
  { label: '在线服务', value: '12', icon: Monitor, color: '#00f0ff' },
  { label: '注册用户', value: '86', icon: User, color: '#7c3aed' },
  { label: 'API 调用', value: '2,431', icon: Connection, color: '#00f0ff' },
  { label: '系统负载', value: '32%', icon: Cpu, color: '#7c3aed' },
]
</script>

<template>
  <div class="dashboard">
    <h2 class="page-title">仪表盘</h2>
    <p class="page-subtitle">欢迎使用 Cybertron Portal 运维管理平台</p>

    <div class="stat-grid">
      <div v-for="item in stats" :key="item.label" class="stat-card tech-card">
        <el-icon :size="28" :color="item.color" class="stat-icon">
          <component :is="item.icon" />
        </el-icon>
        <div class="stat-info">
          <div class="stat-value">{{ item.value }}</div>
          <div class="stat-label">{{ item.label }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.page-title {
  margin: 0 0 8px;
  color: var(--tech-text);
  font-weight: 600;
}
.page-subtitle {
  margin: 0 0 24px;
  color: var(--tech-text-dim);
  font-size: 14px;
}
.stat-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
  gap: 16px;
}
.stat-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px;
  transition: border-color 0.2s, box-shadow 0.2s;
}
.stat-card:hover {
  border-color: var(--tech-border-strong);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.4), var(--tech-glow);
}
.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: var(--tech-text);
}
.stat-label {
  font-size: 13px;
  color: var(--tech-text-dim);
}
</style>
