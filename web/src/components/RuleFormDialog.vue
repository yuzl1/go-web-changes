<script setup lang="ts">
import { ref, watch } from 'vue'
import type { ScanRule, StepResult } from '../types'

const props = defineProps<{
  visible: boolean; rule: ScanRule
  testing: boolean; testResults: StepResult[]; testFinalResult: string; testError: string
}>()
const emit = defineEmits<{ 'update:visible': [v: boolean]; save: [rule: ScanRule]; test: [rule: ScanRule] }>()

const localRule = ref<ScanRule>({ step_order:1, rule_content:'', rule_mode:1 })
watch(() => props.visible, v => { if(v) localRule.value = { ...props.rule } })

const examples = [
  "$('h1').text()",
  "$('div.price').text().trim()",
  "$('a').attr('href')",
  "$('img').map(function(){return $(this).attr('src')}).get().join('\\n')",
  "var arr=[];$('.item').each(function(){arr.push($(this).text())});arr.join('|')",
]

const handleSave = () => { if(!localRule.value.rule_content.trim())return; emit('save',{...localRule.value}) }
const handleTest = () => { if(!localRule.value.rule_content.trim())return; emit('test',{...localRule.value}) }
const quickFill = (ex: string) => { localRule.value.rule_content = ex }

const stepIcon = (s: StepResult) => s.status==='skipped'?'⚠️':s.status==='error'?'✗':'✓'
const stepColor = (s: StepResult) => s.status==='skipped'?'#e6a23c':s.status==='error'?'#f56c6c':'#67c23a'
</script>

<template>
  <el-dialog :model-value="visible" @update:model-value="emit('update:visible',$event)" title="编辑jQuery脚本" width="700px" :close-on-click-modal="false">
    <el-form :model="localRule" label-width="80px">
      <el-row :gutter="16">
        <el-col :span="12">
          <el-form-item label="执行模式">
            <el-radio-group v-model="localRule.rule_mode">
              <el-radio :value="1">必须成功</el-radio>
              <el-radio :value="2">失败跳过</el-radio>
            </el-radio-group>
          </el-form-item>
        </el-col>
      </el-row>
      <el-form-item label="jQuery脚本" required>
        <el-input v-model="localRule.rule_content" type="textarea" :rows="6"
          placeholder="$('h1').text()&#10;$('div.price').text().trim()&#10;$('a').attr('href')"
          style="font-family:monospace;font-size:13px"/>
        <div style="margin-top:6px;display:flex;flex-wrap:wrap;gap:4px">
          <span style="font-size:11px;color:#999;line-height:24px">快捷填入:</span>
          <el-button v-for="ex in examples" :key="ex" size="small" text type="primary" @click="quickFill(ex)">{{ ex.substring(0,35) }}{{ ex.length>35?'...':'' }}</el-button>
        </div>
        <div style="font-size:11px;color:#909399;margin-top:4px;line-height:1.6">
          引用上一步结果: <b>window.__prev_result</b>(文本) 或 <b>window.__prev_html</b>(HTML)<br/>
          例: 上一步提取了HTML, 下一步 <b>$(window.__prev_html).find('.item').text()</b>
        </div>
      </el-form-item>
    </el-form>

    <!-- 测试结果 -->
    <div v-if="testResults.length>0||testError" style="margin-top:12px;border:1px solid #dcdfe6;border-radius:4px;padding:12px;background:#fafafa">
      <h4 style="margin:0 0 8px">测试结果</h4>
      <div v-for="(step,i) in testResults" :key="i" style="margin-bottom:8px;padding:8px;border-radius:4px;background:#fff;border:1px solid #ebeef5">
        <div style="font-size:13px">
          <span :style="{color:stepColor(step),fontWeight:'bold'}">{{ stepIcon(step) }}</span>
          步骤{{ step.step_order }} <span style="color:#999">({{ step.elapsed_ms }}ms)</span>
          <span v-if="step.count >= 0" style="font-size:12px;color:#909399;margin-left:4px">匹配 <b>{{ step.count }}</b> 个元素</span>
          <span v-else-if="step.status==='success' && !step.output" style="font-size:12px;color:#f56c6c;margin-left:4px">⚠ 未匹配到元素或返回为空</span>
        </div>
        <div v-if="step.error" style="font-size:12px;color:#f56c6c;margin-top:2px">{{ step.error }}</div>
        <div v-if="step.output" style="margin-top:4px"><div style="font-size:11px;color:#999">输出:</div><pre style="font-size:12px;background:#f5f7fa;padding:6px 8px;border-radius:3px;overflow:auto;max-height:120px;margin:2px 0 0">{{ step.output }}</pre></div>
      </div>
      <div v-if="testFinalResult" style="margin-top:8px;border-top:1px dashed #dcdfe6;padding-top:8px"><div style="font-size:12px;color:#999;margin-bottom:4px">最终结果:</div><pre style="font-size:12px;background:#f0f9eb;padding:8px;border-radius:4px;overflow:auto;max-height:200px;margin:0;border-left:3px solid #67c23a">{{ testFinalResult }}</pre></div>
      <div v-if="testError" style="color:#f56c6c;font-size:13px;margin-top:4px">{{ testError }}</div>
    </div>

    <template #footer>
      <div style="display:flex;justify-content:space-between">
        <el-button type="warning" :loading="testing" @click="handleTest">▶ 测试执行</el-button>
        <div><el-button type="primary" @click="handleSave">确定</el-button><el-button @click="emit('update:visible',false)">取消</el-button></div>
      </div>
    </template>
  </el-dialog>
</template>
