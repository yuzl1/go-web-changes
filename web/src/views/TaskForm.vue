<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { getTaskDetail, createTask, updateTask } from '../api/task'
import { testRules } from '../api/rule'
import type { ScanRule, StepResult } from '../types'
import { FreqCodeMap, RuleModeMap } from '../types'
import { ElMessage } from 'element-plus'
import RuleFormDialog from '../components/RuleFormDialog.vue'

const router = useRouter(); const route = useRoute()
const isEdit = computed(() => !!route.params.id)
const taskId = computed(() => Number(route.params.id))

const form = ref({ name:'', target_url:'', freq_code:4, status:1, email_notify:1, remark:'' })
const rules = ref<ScanRule[]>([])
const loading = ref(false)

const ruleDialogVisible = ref(false); const editingRuleIndex = ref(-1)
const editingRule = ref<ScanRule>({ step_order:1, rule_content:'', rule_mode:1 })
const testing = ref(false); const testResults = ref<StepResult[]>([]); const testFinalResult = ref(''); const testError = ref('')

const freqOptions = Object.entries(FreqCodeMap).map(([v,l]) => ({value:Number(v),label:l}))

const loadTask = async () => {
  if (!isEdit.value) return; loading.value = true
  try { const t = (await getTaskDetail(taskId.value)).data.data; form.value = { name:t.name, target_url:t.target_url, freq_code:t.freq_code, status:t.status, email_notify:t.email_notify??1, remark:t.remark||'' }; if (t.rules) rules.value = t.rules }
  finally { loading.value = false }
}

const openRuleDialog = (index: number) => {
  editingRuleIndex.value = index
  editingRule.value = index >= 0 ? { ...rules.value[index] } : { step_order:rules.value.length+1, rule_content:'', rule_mode:1 }
  testResults.value = []; testFinalResult.value = ''; testError.value = ''
  ruleDialogVisible.value = true
}

const handleRuleSave = (rule: ScanRule) => {
  if (editingRuleIndex.value >= 0) rules.value[editingRuleIndex.value] = rule
  else { rule.step_order = rules.value.length+1; rules.value.push(rule) }
  ruleDialogVisible.value = false
  rules.value.forEach((r,i) => r.step_order=i+1)
}
const handleRuleDelete = (i:number) => { rules.value.splice(i,1); rules.value.forEach((r,j)=>r.step_order=j+1) }
const handleRuleMove = (i:number, d:number) => { const j=i+d; if(j<0||j>=rules.value.length)return; [rules.value[i],rules.value[j]]=[rules.value[j],rules.value[i]]; rules.value.forEach((r,k)=>r.step_order=k+1) }

const handleTestRule = async (rule: ScanRule) => {
  if (!form.value.target_url) { ElMessage.warning('请先填写目标URL'); return }
  let list: ScanRule[] = editingRuleIndex.value >= 0 ? [...rules.value] : [...rules.value, rule]
  if (editingRuleIndex.value >= 0) list[editingRuleIndex.value] = rule
  testing.value = true; testError.value = ''
  try { const r = await testRules(form.value.target_url, list); testResults.value = r.data.data.steps; testFinalResult.value = r.data.data.final_result }
  catch(e:any) { testError.value = e.message||'请求失败' }
  finally { testing.value = false }
}

const submit = async () => {
  if (!form.value.name) { ElMessage.warning('请输入监听名称'); return }
  if (!/^https?:\/\//.test(form.value.target_url)) { ElMessage.warning('URL格式错误'); return }
  if (rules.value.length===0) { ElMessage.warning('请添加脚本'); return }
  loading.value = true
  try {
    const data = { ...form.value, rules: rules.value }
    if (isEdit.value) { await updateTask(taskId.value, data); ElMessage.success('更新成功') }
    else { await createTask(data); ElMessage.success('创建成功') }
    router.push('/tasks')
  } finally { loading.value = false }
}

onMounted(loadTask)
</script>

<template>
  <div>
    <el-card v-loading="loading">
      <template #header><span>{{ isEdit?'编辑':'新增' }}监听任务</span></template>
      <el-form :model="form" label-width="100px">
        <el-row :gutter="20">
          <el-col :span="12"><el-form-item label="监听名称" required><el-input v-model="form.name" maxlength="100"/></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="执行频率" required><el-select v-model="form.freq_code" style="width:100%"><el-option v-for="o in freqOptions" :key="o.value" :label="o.label" :value="o.value"/></el-select></el-form-item></el-col>
        </el-row>
        <el-form-item label="目标URL" required><el-input v-model="form.target_url" maxlength="2048" placeholder="https://example.com"/></el-form-item>
        <el-row :gutter="20">
          <el-col :span="12"><el-form-item label="状态"><el-switch v-model="form.status" :active-value="1" :inactive-value="0" active-text="启用" inactive-text="禁用"/></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="邮件通知"><el-switch v-model="form.email_notify" :active-value="1" :inactive-value="0" active-text="发送" inactive-text="不发送"/></el-form-item></el-col>
        </el-row>
        <el-form-item label="备注"><el-input v-model="form.remark" type="textarea" maxlength="500"/></el-form-item>
      </el-form>

      <div style="margin-top:24px">
        <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:12px">
          <h4>jQuery 脚本步骤</h4>
          <el-button type="primary" size="small" @click="openRuleDialog(-1)">新增步骤</el-button>
        </div>
        <el-table :data="rules" border size="small">
          <el-table-column label="步骤" width="60"><template #default="{$index}">{{ $index+1 }}</template></el-table-column>
          <el-table-column label="执行模式" width="100"><template #default="{row}">{{ RuleModeMap[row.rule_mode] }}</template></el-table-column>
          <el-table-column prop="rule_content" label="jQuery脚本" min-width="300" show-overflow-tooltip>
            <template #default="{row}"><span style="font-family:monospace;font-size:12px">{{ row.rule_content }}</span></template>
          </el-table-column>
          <el-table-column label="操作" width="160">
            <template #default="{$index}">
              <el-button size="small" @click="handleRuleMove($index,-1)" :disabled="$index===0">↑</el-button>
              <el-button size="small" @click="handleRuleMove($index,1)" :disabled="$index===rules.length-1">↓</el-button>
              <el-button size="small" @click="openRuleDialog($index)">编辑</el-button>
              <el-button size="small" type="danger" @click="handleRuleDelete($index)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
        <el-empty v-if="rules.length===0" description="暂无脚本步骤" :image-size="60"/>
      </div>

      <div style="margin-top:24px;text-align:center">
        <el-button type="primary" size="large" @click="submit" :loading="loading">{{ isEdit?'保存':'创建' }}</el-button>
        <el-button size="large" @click="router.push('/tasks')">取消</el-button>
      </div>
    </el-card>

    <RuleFormDialog v-model:visible="ruleDialogVisible" :rule="editingRule" :testing="testing"
      :test-results="testResults" :test-final-result="testFinalResult" :test-error="testError"
      @save="handleRuleSave" @test="handleTestRule"/>
  </div>
</template>
