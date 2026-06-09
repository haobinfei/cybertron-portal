<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import type { UserInfo } from '@/api/user'

const props = defineProps<{
  visible: boolean
  mode: 'create' | 'edit'
  user: UserInfo | null
}>()

const emit = defineEmits<{
  'update:visible': [value: boolean]
  confirm: [params: any]
}>()

const formRef = ref()
const loading = ref(false)

const form = reactive({
  username: '',
  password: '',
  nickname: '',
  email: '',
  role: 'user',
  status: 1,
})

watch(
  () => props.visible,
  (val) => {
    if (val && props.mode === 'edit' && props.user) {
      form.username = props.user.username
      form.password = ''
      form.nickname = props.user.nickname
      form.email = props.user.email
      form.role = props.user.role
      form.status = props.user.status
    } else if (val && props.mode === 'create') {
      form.username = ''
      form.password = ''
      form.nickname = ''
      form.email = ''
      form.role = 'user'
      form.status = 1
    }
  },
)

const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: props.mode === 'create'
    ? [{ required: true, message: '请输入密码', trigger: 'blur' }]
    : [],
}

async function handleConfirm() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  const params: any = {}
  if (props.mode === 'create') {
    params.username = form.username
    params.password = form.password
  }
  if (form.nickname) params.nickname = form.nickname
  if (form.email) params.email = form.email
  if (form.role) params.role = form.role
  params.status = form.status
  if (props.mode === 'edit' && form.password) {
    params.password = form.password
  }

  loading.value = true
  emit('confirm', params)
  loading.value = false
}

function handleClose() {
  emit('update:visible', false)
}
</script>

<template>
  <el-dialog
    :model-value="visible"
    :title="mode === 'create' ? '创建用户' : '编辑用户'"
    width="500px"
    @close="handleClose"
  >
    <el-form ref="formRef" :model="form" :rules="rules" label-width="80px">
      <el-form-item label="用户名" prop="username">
        <el-input v-model="form.username" :disabled="mode === 'edit'" />
      </el-form-item>
      <el-form-item label="密码" :prop="mode === 'create' ? 'password' : ''">
        <el-input
          v-model="form.password"
          type="password"
          show-password
          :placeholder="mode === 'edit' ? '留空不修改密码' : ''"
        />
      </el-form-item>
      <el-form-item label="昵称">
        <el-input v-model="form.nickname" />
      </el-form-item>
      <el-form-item label="邮箱">
        <el-input v-model="form.email" />
      </el-form-item>
      <el-form-item label="角色">
        <el-select v-model="form.role" style="width: 100%">
          <el-option label="管理员" value="admin" />
          <el-option label="普通用户" value="user" />
        </el-select>
      </el-form-item>
      <el-form-item label="状态">
        <el-radio-group v-model="form.status">
          <el-radio :value="1">启用</el-radio>
          <el-radio :value="0">禁用</el-radio>
        </el-radio-group>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button type="primary" :loading="loading" @click="handleConfirm">
        {{ mode === 'create' ? '创建' : '保存' }}
      </el-button>
    </template>
  </el-dialog>
</template>
