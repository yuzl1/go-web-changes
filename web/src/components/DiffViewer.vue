<script setup lang="ts">
import { computed } from 'vue'
import { diffLines } from 'diff'

const props = defineProps<{
  oldText: string
  newText: string
}>()

const hasOld = computed(() => props.oldText && props.oldText.trim() !== '')
const hasNew = computed(() => props.newText && props.newText.trim() !== '')

const diffParts = computed(() => {
  if (!hasOld.value || !hasNew.value) return []
  return diffLines(props.oldText || '', props.newText || '')
})
</script>

<template>
  <div>
    <div v-if="!hasOld" style="text-align: center; color: #999; padding: 40px">
      （首次扫描，无历史内容）
    </div>
    <div v-else-if="!hasNew" style="text-align: center; color: #999; padding: 40px">
      （无扫描内容）
    </div>
    <div v-else style="display: flex; gap: 12px">
      <div style="flex: 1; border: 1px solid #dcdfe6; border-radius: 4px; overflow: hidden">
        <div style="background: #f5f7fa; padding: 8px 12px; font-size: 13px; font-weight: bold; border-bottom: 1px solid #dcdfe6">
          上次扫描内容
        </div>
        <pre style="padding: 12px; margin: 0; font-size: 12px; line-height: 1.6; overflow: auto; max-height: 500px; white-space: pre-wrap; word-break: break-all">{{ oldText }}</pre>
      </div>
      <div style="flex: 1; border: 1px solid #dcdfe6; border-radius: 4px; overflow: hidden">
        <div style="background: #f5f7fa; padding: 8px 12px; font-size: 13px; font-weight: bold; border-bottom: 1px solid #dcdfe6">
          本次扫描内容
        </div>
        <div v-if="diffParts.length === 0" style="padding: 12px; font-size: 12px; color: #999">
          内容完全相同
        </div>
        <pre v-else style="padding: 12px; margin: 0; font-size: 12px; line-height: 1.6; overflow: auto; max-height: 500px; white-space: pre-wrap; word-break: break-all"><template v-for="(part, i) in diffParts" :key="i"><span v-if="part.added" style="background-color: #d4edda; color: #155724">+ {{ part.value }}</span><span v-else-if="part.removed" style="background-color: #f8d7da; color: #721c24">- {{ part.value }}</span><span v-else>{{ part.value }}</span></template></pre>
      </div>
    </div>
  </div>
</template>
