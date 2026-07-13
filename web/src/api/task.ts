import request from '../utils/request'
import type { ApiResponse, PageData, MonitorTask } from '../types'

export function getTaskList(params: { keyword?: string; status?: number; page?: number; page_size?: number }) {
  return request.get<ApiResponse<PageData<MonitorTask>>>('/tasks', { params })
}

export function getTaskDetail(id: number) {
  return request.get<ApiResponse<MonitorTask>>(`/tasks/${id}`)
}

export function createTask(data: {
  name: string
  target_url: string
  freq_code: number
  status: number
  remark: string
  rules: any[]
}) {
  return request.post<ApiResponse>('/tasks', data)
}

export function updateTask(id: number, data: {
  name: string
  target_url: string
  freq_code: number
  status: number
  remark: string
  rules: any[]
}) {
  return request.put<ApiResponse>(`/tasks/${id}`, data)
}

export function deleteTask(id: number) {
  return request.delete<ApiResponse>(`/tasks/${id}`)
}

export function executeTask(id: number) {
  return request.post<ApiResponse>(`/tasks/${id}/execute`)
}
