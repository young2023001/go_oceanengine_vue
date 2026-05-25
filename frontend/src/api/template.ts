import request from './request'

export interface ProjectTemplate {
  id: number
  name: string
  config: Record<string, unknown>
  use_count: number
  created_by: number
  created_at: string
}

export interface PromotionTemplate {
  id: number
  name: string
  config: Record<string, unknown>
  use_count: number
  created_by: number
  created_at: string
}

export interface CreateTemplateParams {
  name: string
  config: Record<string, unknown>
}

export const templateApi = {
  createProject(data: CreateTemplateParams) {
    return request.post<ProjectTemplate>('/templates/projects', data)
  },
  listProjects() {
    return request.get<ProjectTemplate[]>('/templates/projects')
  },
  getProject(id: number) {
    return request.get<ProjectTemplate>(`/templates/projects/${id}`)
  },
  updateProject(id: number, data: CreateTemplateParams) {
    return request.put<void>(`/templates/projects/${id}`, data)
  },
  deleteProject(id: number) {
    return request.delete<void>(`/templates/projects/${id}`)
  },
  createPromotion(data: CreateTemplateParams) {
    return request.post<PromotionTemplate>('/templates/promotions', data)
  },
  listPromotions() {
    return request.get<PromotionTemplate[]>('/templates/promotions')
  },
  deletePromotion(id: number) {
    return request.delete<void>(`/templates/promotions/${id}`)
  }
}
