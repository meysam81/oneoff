import ky from 'ky'

const api = ky.create({
  prefixUrl: '/api',
  timeout: 30000,
  retry: {
    limit: 2,
    methods: ['get', 'post', 'patch', 'delete'],
    statusCodes: [408, 413, 429, 500, 502, 503, 504],
  },
  hooks: {
    beforeError: [
      async (error) => {
        const { response } = error
        if (response && response.body) {
          const body = await response.json()
          error.message = body.error || error.message
        }
        return error
      },
    ],
  },
})

// Jobs API
export const jobsAPI = {
  list: (params) => api.get('jobs', { searchParams: params }).json(),
  get: (id) => api.get(`jobs/${id}`).json(),
  create: (data) => api.post('jobs', { json: data }).json(),
  update: (id, data) => api.patch(`jobs/${id}`, { json: data }).json(),
  delete: (id) => api.delete(`jobs/${id}`),
  execute: (id) => api.post(`jobs/${id}/execute`).json(),
  clone: (id, scheduledAt) => api.post(`jobs/${id}/clone`, { json: { scheduled_at: scheduledAt } }).json(),
  cancel: (id) => api.post(`jobs/${id}/cancel`).json(),
}

// Executions API
export const executionsAPI = {
  list: (params) => api.get('executions', { searchParams: params }).json(),
  get: (id) => api.get(`executions/${id}`).json(),
}

// Projects API
export const projectsAPI = {
  list: (includeArchived = false) => api.get('projects', { searchParams: { include_archived: includeArchived } }).json(),
  get: (id) => api.get(`projects/${id}`).json(),
  create: (data) => api.post('projects', { json: data }).json(),
  update: (id, data) => api.patch(`projects/${id}`, { json: data }).json(),
  delete: (id) => api.delete(`projects/${id}`),
}

// Tags API
export const tagsAPI = {
  list: () => api.get('tags').json(),
  get: (id) => api.get(`tags/${id}`).json(),
  create: (data) => api.post('tags', { json: data }).json(),
  update: (id, data) => api.patch(`tags/${id}`, { json: data }).json(),
  delete: (id) => api.delete(`tags/${id}`),
}

// System API
export const systemAPI = {
  status: () => api.get('system/status').json(),
  config: () => api.get('system/config').json(),
  updateConfig: (key, value) => api.patch('system/config', { json: { key, value } }).json(),
  workerStatus: () => api.get('workers/status').json(),
  jobTypes: () => api.get('job-types').json(),
}

export default api
