<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getRecords, getRecordDetail } from '../api/record'
import { getTaskDetail } from '../api/task'
import type { MonitorRecord, MonitorTask, RecordDetail } from '../types'
import DiffViewer from '../components/DiffViewer.vue'

const route = useRoute()
const taskId = Number(route.params.id)

const task = ref<MonitorTask | null>(null)
const list = ref<MonitorRecord[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const loading = ref(false)

// 详情抽屉
const drawerVisible = ref(false)
const detail = ref<RecordDetail | null>(null)
const detailLoading = ref(false)

const fetchTask = async () => {
  const res = await getTaskDetail(taskId)
  task.value = res.data.data
}

const fetchList = async () => {
  loading.value = true
  try {
    const res = await getRecords(taskId, page.value, pageSize.value)
    list.value = res.data.data.list
    total.value = res.data.data.total
  } finally {
    loading.value = false
  }
}

const handlePageChange = (p: number) => {
  page.value = p
  fetchList()
}

const openDetail = async (recordId: number) => {
  drawerVisible.value = true
  detailLoading.value = true
  try {
    const res = await getRecordDetail(recordId)
    detail.value = res.data.data
  } finally {
    detailLoading.value = false
  }
}

const changedTag = (isChanged: number) => {
  if (isChanged === 1) return { type: 'success' as const, text: '有变更' }
  if (isChanged === 0) return { type: 'info' as const, text: '无变更' }
  return { type: 'danger' as const, text: '失败' }
}

onMounted(() => {
  fetchTask()
  fetchList()
})
</script>

<template>
  <div>
    <el-card>
      <template #header>
        <el-breadcrumb separator="/">
          <el-breadcrumb-item :to="{ path: '/tasks' }">监听列表</el-breadcrumb-item>
          <el-breadcrumb-item>扫描历史</el-breadcrumb-item>
        </el-breadcrumb>
      </template>

      <!-- 任务摘要 -->
      <el-descriptions v-if="task" :column="2" border size="small" style="margin-bottom: 20px">
        <el-descriptions-item label="任务名称">{{ task.name }}</el-descriptions-item>
        <el-descriptions-item label="目标URL">{{ task.target_url }}</el-descriptions-item>
      </el-descriptions>

      <el-table :data="list" v-loading="loading" stripe>
        <el-table-column prop="scan_time" label="扫描时间" width="180" />
        <el-table-column prop="target_url" label="目标URL" min-width="200" show-overflow-tooltip />
        <el-table-column label="是否有变更" width="100">
          <template #default="{ row }">
            <el-tag :type="changedTag(row.is_changed).type" size="small">
              {{ changedTag(row.is_changed).text }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="提取内容" min-width="180" show-overflow-tooltip>
          <template #default="{ row }">{{ row.scan_preview || '-' }}</template>
        </el-table-column>
        <el-table-column label="邮件" width="80">
          <template #default="{ row }">
            <el-tag v-if="row.email_sent === 1" type="success" size="small">已发送</el-tag>
            <el-tag v-else-if="row.email_sent === -1" type="danger" size="small">失败</el-tag>
            <span v-else style="color: #999; font-size: 12px">未发送</span>
          </template>
        </el-table-column>
        <el-table-column prop="error_msg" label="错误信息" min-width="140" show-overflow-tooltip>
          <template #default="{ row }">{{ row.error_msg || '-' }}</template>
        </el-table-column>
        <el-table-column label="操作" width="100">
          <template #default="{ row }">
            <el-button size="small" type="primary" @click="openDetail(row.id)" :disabled="row.is_changed === -1">
              查看详情
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-if="total > 0"
        style="margin-top: 16px; justify-content: center"
        background
        layout="prev, pager, next"
        :total="total"
        :page-size="pageSize"
        :current-page="page"
        @current-change="handlePageChange"
      />
    </el-card>

    <!-- 详情抽屉 -->
    <el-drawer v-model="drawerVisible" title="扫描详情" size="80%" direction="rtl">
      <div v-loading="detailLoading">
        <template v-if="detail">
          <el-descriptions :column="2" border size="small" style="margin-bottom: 20px">
            <el-descriptions-item label="扫描时间">{{ detail.scan_time }}</el-descriptions-item>
            <el-descriptions-item label="目标URL">{{ detail.target_url }}</el-descriptions-item>
            <el-descriptions-item label="是否有变更">
              <el-tag :type="changedTag(detail.is_changed).type" size="small">
                {{ changedTag(detail.is_changed).text }}
              </el-tag>
            </el-descriptions-item>
          </el-descriptions>

          <!-- 本次提取内容 -->
          <div v-if="detail.scan_result" style="margin-bottom: 16px; border: 1px solid #dcdfe6; border-radius: 4px; overflow: hidden">
            <div style="background: #e1f3d8; padding: 8px 12px; font-size: 13px; font-weight: bold; border-bottom: 1px solid #dcdfe6">
              本次扫描提取内容
            </div>
            <pre style="padding: 12px; margin: 0; font-size: 12px; line-height: 1.6; overflow: auto; max-height: 300px; white-space: pre-wrap; word-break: break-all; background: #fafafa">{{ detail.scan_result }}</pre>
          </div>

          <DiffViewer :old-text="detail.prev_scan_result" :new-text="detail.scan_result" />
        </template>
      </div>
    </el-drawer>
  </div>
</template>
