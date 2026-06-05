<script setup lang="ts">
</script>

<template>
  <div class="login-container">
    <el-card class="login-card">
      <h2>登录</h2>
      <el-form>
        <el-form-item>
          <el-input placeholder="用户名" />
        </el-form-item>
        <el-form-item>
          <el-input type="password" placeholder="密码" show-password />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" style="width: 100%">登录</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background: #f0f2f5;
}
.login-card {
  width: 400px;
}
.login-card h2 {
  text-align: center;
  margin-bottom: 24px;
}
</style>
