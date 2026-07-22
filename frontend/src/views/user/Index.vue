<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getUserList, createUser, updateUser, deleteUser, type UserInfo, type CreateUserParams, type UpdateUserParams } from '@/api/user'
import UserDialog from './UserDialog.vue'

const loading = ref(false)
const users = ref<UserInfo[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

const dialogVisible = ref(false)
const dialogMode = ref<'create' | 'edit'>('create')
const currentUser = ref<UserInfo | null>(null)

async function fetchUsers() {
  loading.value = true
  try {
    const res = await getUserList({ page: page.value, page_size: pageSize.value })
    users.value = res.list
    total.value = res.total
  } finally {
    loading.value = false
  }
}

function handleCreate() {
  dialogMode.value = 'create'
  currentUser.value = null
  dialogVisible.value = true
}

function handleEdit(user: UserInfo) {
  dialogMode.value = 'edit'
  currentUser.value = { ...user }
  dialogVisible.value = true
}

async function handleDelete(user: UserInfo) {
  try {
    await ElMessageBox.confirm(`确定要删除用户「${user.username}」吗？`, '确认删除', {
      type: 'warning',
    })
    await deleteUser(user.id)
    ElMessage.success('删除成功')
    fetchUsers()
  } catch {
    // cancelled
  }
}

async function handleDialogConfirm(params: CreateUserParams | UpdateUserParams) {
  if (dialogMode.value === 'create') {
    await createUser(params as CreateUserParams)
    ElMessage.success('创建用户成功')
  } else {
    await updateUser(currentUser.value!.id, params as UpdateUserParams)
    ElMessage.success('更新用户成功')
  }
  dialogVisible.value = false
  fetchUsers()
}

function handlePageChange(p: number) {
  page.value = p
  fetchUsers()
}

function handleSizeChange(s: number) {
  pageSize.value = s
  page.value = 1
  fetchUsers()
}

onMounted(() => {
  fetchUsers()
})
</script>

<template>
  <div class="user-manage">
    <div class="toolbar">
      <h3>用户管理</h3>
      <el-button type="primary" @click="handleCreate">
        <el-icon><Plus /></el-icon>
        创建用户
      </el-button>
    </div>

    <el-table :data="users" border stripe v-loading="loading">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="username" label="用户名" width="150" />
      <el-table-column prop="nickname" label="昵称" width="150" />
      <el-table-column prop="email" label="邮箱" min-width="200" />
      <el-table-column prop="role" label="角色" width="100">
        <template #default="{ row }">
          <el-tag :type="row.role === 'admin' ? 'danger' : 'info'" size="small">
            {{ row.role === 'admin' ? '管理员' : '普通用户' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="80">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'warning'" size="small">
            {{ row.status === 1 ? '启用' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="last_login_at" label="最后登录" width="180" />
      <el-table-column prop="created_at" label="创建时间" width="180" />
      <el-table-column label="操作" width="180" fixed="right">
        <template #default="{ row }">
          <el-button size="small" @click="handleEdit(row)">编辑</el-button>
          <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <div class="pagination">
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next"
        @current-change="handlePageChange"
        @size-change="handleSizeChange"
      />
    </div>

    <UserDialog
      v-model:visible="dialogVisible"
      :mode="dialogMode"
      :user="currentUser"
      @confirm="handleDialogConfirm"
    />
  </div>
</template>

<script lang="ts">
import { Plus } from '@element-plus/icons-vue'
export default { components: { Plus } }
</script>

<style scoped>
.user-manage {
  padding: 4px;
}
.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}
.toolbar h3 {
  margin: 0;
  color: var(--tech-text);
  font-weight: 600;
}
.pagination {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
</style>
