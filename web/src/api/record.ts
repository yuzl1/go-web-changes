import request from '../utils/request'
import type { ApiResponse, PageData, MonitorRecord, RecordDetail } from '../types'

export function getRecords(taskId: number, page = 1, pageSize = 10) {
  return request.get<ApiResponse<PageData<MonitorRecord>>>(`/tasks/${taskId}/records`, {
    params: { page, page_size: pageSize },
  })
}

export function getRecordDetail(id: number) {
  return request.get<ApiResponse<RecordDetail>>(`/records/${id}`)
}
