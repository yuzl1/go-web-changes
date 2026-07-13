<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getTaskList, deleteTask, executeTask } from '../api/task'
import type { MonitorTask } from '../types'
import { FreqCodeMap } from '../types'
import { ElMessage, ElMessageBox } from 'element-plus'

const router = useRouter()
const loading = ref(false)
const list = ref<MonitorTask[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const keyword = ref('')
const executing = ref<Set<number>>(new Set())

const fetchList = async () => {
  loading.value = true
  try {
    const res = await getTaskList({ keyword: keyword.value, page: page.value, page_size: pageSize.value })
    list.value = res.data.data.list
    total.value = res.data.data.total
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  page.value = 1
  fetchList()
}

const handlePageChange = (p: number) => {
  page.value = p
  fetchList()
}

const handleDelete = async (row: MonitorTask) => {
  try {
    await ElMessageBox.confirm('确认删除该监听任务及其所有扫描历史和规则？', '删除确认', {
      type: 'warning',
      confirmButtonText: '确认删除',
      cancelButtonText: '取消',
    })
    await deleteTask(row.id)
    ElMessage.success('删除成功')
    fetchList()
  } catch { /* cancelled */ }
}

const handleExecute = async (row: MonitorTask) => {
  executing.value.add(row.id)
  try {
    const res = await executeTask(row.id)
    const d = res.data.data
    if (d.is_changed === 1) {
      ElMessage.success('扫描完成 — 检测到变更！')
    } else if (d.is_changed === 0) {
      ElMessage.success('扫描完成 — 无变更')
    } else {
      ElMessage.warning('扫描失败: ' + (d.error || '未知错误'))
    }
    fetchList()
  } finally {
    executing.value.delete(row.id)
  }
}

onMounted(fetchList)
</script>

<template>
  <div>
    <el-card>
      <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px">
        <div style="display: flex; gap: 12px">
          <el-input v-model="keyword" placeholder="搜索监听名称" clearable style="width: 280px" @keyup.enter="handleSearch" />
          <el-button type="primary" @click="handleSearch">搜索</el-button>
        </div>
        <el-button type="primary" @click="router.push('/tasks/new')">新增任务</el-button>
      </div>

      <el-table :data="list" v-loading="loading" stripe>
        <el-table-column label="序号" type="index" width="60" />
        <el-table-column prop="name" label="监听名称" min-width="160" />
        <el-table-column prop="target_url" label="目标URL" min-width="240" show-overflow-tooltip />
        <el-table-column label="执行频率" width="110">
          <template #default="{ row }">{{ FreqCodeMap[row.freq_code] || '-' }}</template>
        </el-table-column>
        <el-table-column prop="last_scan_time" label="上次扫描时间" width="170">
          <template #default="{ row }">{{ row.last_scan_time || '从未扫描' }}</template>
        </el-table-column>
        <el-table-column label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="280" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="router.push(`/tasks/${row.id}/edit`)">编辑</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
            <el-button size="small" type="success" :loading="executing.has(row.id)" @click="handleExecute(row)">立即执行</el-button>
            <el-button size="small" @click="router.push(`/tasks/${row.id}/records`)">历史</el-button>
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
  </div>
</template>
