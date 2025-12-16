import axios, { type AxiosResponse, type InternalAxiosRequestConfig } from 'axios'
import { useUserStore } from '@/store/user'
import router from '@/router'

const request = axios.create({
  baseURL: '/api/v1',
  timeout: 15000
})

request.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const userStore = useUserStore()
    if (userStore.token) {
      config.headers.Authorization = `Bearer ${userStore.token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

request.interceptors.response.use(
  (response: AxiosResponse) => {
    const { code, msg, data } = response.data || {}

    if (code === 200) {
      return data
    }

    return Promise.reject(new Error(msg || '请求失败'))
  },
  (error) => {
    if (error.response) {
      const { status, data } = error.response

      if (status === 401) {
        const userStore = useUserStore()
        userStore.logout()
        router.push({
          name: 'login',
          query: { redirect: router.currentRoute.value.fullPath, expired: '1' }
        })
        return Promise.reject(new Error('登录已过期，请重新登录'))
      }

      if (data?.msg) {
        return Promise.reject(new Error(data.msg))
      }
    }

    if (error.code === 'ECONNABORTED') {
      return Promise.reject(new Error('请求超时，请稍后重试'))
    }

    return Promise.reject(new Error('网络错误，请检查网络连接'))
  }
)

export default request
