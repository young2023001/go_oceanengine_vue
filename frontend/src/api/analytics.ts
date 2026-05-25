import request, { PageResponse } from './request'

export interface OverviewData {
  date: string
  impressions: number
  clicks: number
  cost: number
  conversions: number
}

export interface TrendData {
  date: string
  value: number
}

export interface RankItem {
  id: number
  name: string
  value: number
}

export interface CompareData {
  current: OverviewData
  previous: OverviewData
}

export interface DetailItem {
  id: number
  name: string
  impressions: number
  clicks: number
  cost: number
  conversions: number
}

export interface ExportTask {
  id: number
  title: string
  status: string
  url?: string
  created_at: string
}

export interface TrendParams {
  start_date: string
  end_date: string
}

export interface RankParams {
  start_date: string
  end_date: string
  order_by?: string
  limit?: number
}

export interface CompareParams {
  date: string
  compare_type?: string
}

export interface DetailParams {
  start_date: string
  end_date: string
  page?: number
  page_size?: number
}

export interface CreateExportParams {
  title: string
  start_date: string
  end_date: string
}

export const analyticsApi = {
  getOverview(params?: { date?: string }) {
    return request.get<OverviewData>('/analytics/overview', params)
  },

  getTrend(params: TrendParams) {
    return request.get<TrendData[]>('/analytics/trend', params)
  },

  getRank(params: RankParams) {
    return request.get<RankItem[]>('/analytics/rank', params)
  },

  getCompare(params: CompareParams) {
    return request.get<CompareData>('/analytics/compare', params)
  },

  getDetail(params: DetailParams) {
    return request.get<PageResponse<DetailItem>>('/analytics/detail', params)
  },

  createExport(data: CreateExportParams) {
    return request.post<ExportTask>('/analytics/export', data)
  },

  getExports() {
    return request.get<ExportTask[]>('/analytics/exports')
  }
}
