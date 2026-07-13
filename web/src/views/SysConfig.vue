<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getSmtpConfig, saveSmtpConfig, testSmtp } from '../api/config'
import type { SmtpConfig } from '../types'
import { ElMessage } from 'element-plus'

const form = ref<SmtpConfig>({
  smtp_host: '',
  smtp_port: 587,
  smtp_encryption: 'TLS',
  smtp_sender_email: '',
  smtp_sender_name: '',
  smtp_username: '',
  smtp_password: '',
  smtp_receiver_emails: '',
})

const loading = ref(false)
const testing = ref(false)

const encryptionOptions = [
  { value: 'PLAIN', label: 'PLAIN — 无加密，明文传输' },
  { value: 'TLS', label: 'TLS — TLS加密传输' },
]

const portHint = () => {
  if (form.value.smtp_encryption === 'PLAIN') return '常用端口: 25'
  return '常用端口: 587 (STARTTLS) / 465 (SMTPS)'
}

const loadConfig = async () => {
  loading.value = true
  try {
    const res = await getSmtpConfig()
    if (res.data.data) {
      form.value = { ...form.value, ...res.data.data }
      // 密码不显示
      form.value.smtp_password = ''
    }
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  if (!form.value.smtp_host) { ElMessage.warning('请输入SMTP服务器'); return }
  if (!form.value.smtp_port) { ElMessage.warning('请输入SMTP端口'); return }
  if (!form.value.smtp_sender_email) { ElMessage.warning('请输入发件人邮箱'); return }
  if (!form.value.smtp_username) { ElMessage.warning('请输入邮箱账号'); return }
  if (!form.value.smtp_receiver_emails) { ElMessage.warning('请输入收件人邮箱'); return }

  loading.value = true
  try {
    await saveSmtpConfig(form.value)
    ElMessage.success('SMTP 配置保存成功')
    form.value.smtp_password = '' // 清空密码字段
  } finally {
    loading.value = false
  }
}

const handleTest = async () => {
  // 先保存再测试
  await handleSave()
  testing.value = true
  try {
    await testSmtp()
    ElMessage.success('测试邮件已发送，请检查收件邮箱')
  } catch {
    // 错误已在 request 拦截器中处理
  } finally {
    testing.value = false
  }
}

onMounted(loadConfig)
</script>

<template>
  <div>
    <el-card v-loading="loading">
      <template #header>
        <span>SMTP 邮件配置</span>
      </template>

      <el-form :model="form" label-width="140px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="SMTP服务器" required>
              <el-input v-model="form.smtp_host" placeholder="smtp.example.com" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="SMTP端口" required>
              <el-input-number v-model="form.smtp_port" :min="1" :max="65535" style="width: 100%" />
              <div style="font-size: 12px; color: #999; margin-top: 4px">{{ portHint() }}</div>
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="加密方式" required>
          <el-radio-group v-model="form.smtp_encryption">
            <el-radio v-for="o in encryptionOptions" :key="o.value" :value="o.value">{{ o.label }}</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="发件人邮箱" required>
              <el-input v-model="form.smtp_sender_email" placeholder="sender@example.com" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="发件人名称">
              <el-input v-model="form.smtp_sender_name" placeholder="网页变更监听" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="邮箱账号" required>
              <el-input v-model="form.smtp_username" placeholder="SMTP认证用户名" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="邮箱密码">
              <el-input v-model="form.smtp_password" type="password" show-password placeholder="留空则不修改密码" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="收件人邮箱" required>
          <el-input v-model="form.smtp_receiver_emails" placeholder="receiver@example.com（多个用逗号分隔）" />
        </el-form-item>
      </el-form>

      <div style="text-align: center; margin-top: 20px">
        <el-button type="primary" size="large" @click="handleSave" :loading="loading">保存配置</el-button>
        <el-button type="success" size="large" @click="handleTest" :loading="testing">测试发送</el-button>
      </div>
    </el-card>
  </div>
</template>
