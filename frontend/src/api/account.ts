import request, { PageResponse } from './request'

export interface Account {
  id: number
  local_account_id: number
  name: string
  created_at: string
}

export interface AccountListParams {
  page?: number
  page_size?: number
}

export interface ImportAccountItem {
  local_account_id: number
  name: string
}

export interface ImportAccountsParams {
  items: ImportAccountItem[]
}

export const accountApi = {
  getList(params?: AccountListParams) {
    return request.get<PageResponse<Account>>('/accounts', params)
  },

  importAccounts(data: ImportAccountsParams) {
    return request.post<void>('/accounts/import', data)
  }
}
