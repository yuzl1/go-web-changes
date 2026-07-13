import request from '../utils/request'
import type { ApiResponse, SmtpConfig } from '../types'

export function getSmtpConfig() {
  return request.get<ApiResponse<SmtpConfig>>('/config/smtp')
}

export function saveSmtpConfig(data: SmtpConfig) {
  return request.put<ApiResponse>('/config/smtp', data)
}

export function testSmtp() {
  return request.post<ApiResponse>('/config/smtp/test')
}
