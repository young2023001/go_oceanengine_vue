import request from './request'

export interface Tenant {
  id: number
  name: string
  oauth_app_id: string
  oauth_secret: string
  created_at: string
}

export interface CreateTenantParams {
  name: string
  oauth_app_id: string
  oauth_secret: string
}

export interface OAuthURLResponse {
  url: string
}

export const tenantApi = {
  getList() {
    return request.get<Tenant[]>('/tenants')
  },

  create(data: CreateTenantParams) {
    return request.post<Tenant>('/tenants', data)
  },

  getOAuthURL(id: number) {
    return request.get<OAuthURLResponse>(`/tenants/${id}/oauth/url`)
  }
}
