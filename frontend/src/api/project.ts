import request, { PageResponse } from './request'

export interface Project {
  id: number
  name: string
  status: string
  created_at: string
}

export interface ProjectListParams {
  status?: string
  page?: number
  page_size?: number
}

export interface Promotion {
  id: number
  project_id: number
  name: string
  status: string
  created_at: string
}

export interface PromotionListParams {
  project_id?: number
  page?: number
  page_size?: number
}

export const projectApi = {
  getList(params?: ProjectListParams) {
    return request.get<PageResponse<Project>>('/projects', params)
  },

  getDetail(id: number) {
    return request.get<Project>(`/projects/${id}`)
  },

  updateStatus(id: number, status: string) {
    return request.put<void>(`/projects/${id}/status`, { status })
  },

  getPromotions(params?: PromotionListParams) {
    return request.get<PageResponse<Promotion>>('/promotions', params)
  }
}
