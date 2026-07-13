export interface MonitorTask {
  id: number; name: string; target_url: string; freq_code: number; freq_desc?: string
  status: number; email_notify: number
  last_scan_time: string | null; last_scan_content?: string; remark: string
  rules?: ScanRule[]; created_at: string; updated_at: string
}
export interface ScanRule {
  id?: number; task_id?: number; step_order: number
  rule_content: string   // jQuery脚本
  rule_mode: number      // 1=必须成功 2=失败跳过
}
export interface MonitorRecord {
  id: number; task_id: number; target_url: string; scan_time: string
  scan_result?: string; is_changed: number; error_msg: string | null
  email_sent: number; scan_preview?: string; created_at: string
}
export interface StepResult {
  step_order: number; status: 'success'|'error'|'skipped'; elapsed_ms: number
  output: string; count: number; error?: string
}
export interface RuleTestResult { steps: StepResult[]; final_result: string; message?: string }
export interface ApiResponse<T=any> { code: number; message: string; data: T }
export interface PageData<T=any> { list: T[]; total: number; page: number; page_size: number }
export const FreqCodeMap: Record<number,string> = {1:'每5分钟',2:'每15分钟',3:'每30分钟',4:'每1小时',5:'每6小时',6:'每12小时',7:'每天1次'}
export const RuleModeMap: Record<number,string> = {1:'必须成功',2:'失败跳过'}
