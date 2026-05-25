import request, { PageResponse } from './request'

export interface Group {
  id: number
  name: string
  type: 'franchisee' | 'region' | 'custom'
  parent_id: number
  sort: number
  created_at: string
}

export interface CreateGroupParams {
  name: string
  type: 'franchisee' | 'region' | 'custom'
  parent_id?: number
  sort?: number
}

export interface UpdateGroupParams {
  name?: string
  sort?: number
}

export interface MembersParams {
  account_ids: number[]
}

export const groupApi = {
  create(data: CreateGroupParams) {
    return request.post<Group>('/groups', data)
  },
  getList(params?: { type?: string }) {
    return request.get<Group[]>('/groups', params)
  },
  update(id: number, data: UpdateGroupParams) {
    return request.put<void>(`/groups/${id}`, data)
  },
  remove(id: number) {
    return request.delete<void>(`/groups/${id}`)
  },
  addMembers(id: number, data: MembersParams) {
    return request.post<void>(`/groups/${id}/members`, data)
  },
  removeMembers(id: number, data: MembersParams) {
    return request.post<void>(`/groups/${id}/members/remove`, data)
  }
}
