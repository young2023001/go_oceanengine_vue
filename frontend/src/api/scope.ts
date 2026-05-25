import request from './request'

export interface UserScope {
  account_ids: number[]
}

export const scopeApi = {
  getScope(userId: number) {
    return request.get<UserScope>(`/system/users/${userId}/scope`)
  },
  setScope(userId: number, data: { account_ids: number[] }) {
    return request.post<void>(`/system/users/${userId}/scope`, data)
  }
}
