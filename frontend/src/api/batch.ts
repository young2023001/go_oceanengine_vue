import request, { PageResponse } from './request'

export interface BatchTask {
  id: number
  task_type: string
  title: string
  config: object
  account_ids: number[]
  status: string
  created_at: string
}

export interface BatchTaskListParams {
  page?: number
  page_size?: number
}

export interface CreateBatchTaskParams {
  task_type: string
  title: string
  config: object
  account_ids: number[]
}

export const batchApi = {
  createTask(data: CreateBatchTaskParams) {
    return request.post<BatchTask>('/batch/projects', data)
  },

  getList(params?: BatchTaskListParams) {
    return request.get<PageResponse<BatchTask>>('/batch/tasks', params)
  },

  getDetail(id: number) {
    return request.get<BatchTask>(`/batch/tasks/${id}`)
  },

  cancel(id: number) {
    return request.post<void>(`/batch/tasks/${id}/cancel`)
  },

  retry(id: number) {
    return request.post<void>(`/batch/tasks/${id}/retry`)
  }
}
