import request from '../utils/request'
import type { ApiResponse, ScanRule, RuleTestResult } from '../types'
export const getRules = (tid: number) => request.get<ApiResponse<ScanRule[]>>(`/tasks/${tid}/rules`)
export const createRule = (tid: number, d: Partial<ScanRule>) => request.post<ApiResponse>(`/tasks/${tid}/rules`, d)
export const updateRule = (id: number, d: Partial<ScanRule>) => request.put<ApiResponse>(`/rules/${id}`, d)
export const deleteRule = (id: number) => request.delete<ApiResponse>(`/rules/${id}`)
export const testRules = (url: string, rules: Partial<ScanRule>[]) => request.post<ApiResponse<RuleTestResult>>('/rules/test', { target_url: url, rules })
export const cachePage = (url: string) => request.post<ApiResponse<{html:string;html_length:number}>>('/cache', { target_url: url })
