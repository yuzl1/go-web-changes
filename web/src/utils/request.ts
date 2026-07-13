import axios from 'axios'
import type { ApiResponse } from '../types'
import { ElMessage } from 'element-plus'

const request = axios.create({
  baseURL: '/api',
  timeout: 30000,
})

// 响应拦截器
request.interceptors.response.use(
  (response) => {
    const res: ApiResponse = response.data
    if (res.code !== 0) {
      ElMessage.error(res.message || '请求失败')
      return Promise.reject(new Error(res.message))
    }
    return response
  },
  (error) => {
    ElMessage.error('网络错误: ' + error.message)
    return Promise.reject(error)
  }
)

export default request
