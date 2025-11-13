import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { jobsAPI } from '../utils/api'

export const useJobsStore = defineStore('jobs', () => {
  const jobs = ref([])
  const currentJob = ref(null)
  const loading = ref(false)
  const filter = ref({
    project_id: '',
    status: '',
    type: '',
    search: '',
    tag_ids: [],
    sort_by: 'scheduled_at',
    sort_order: 'asc',
    limit: 50,
    offset: 0,
  })
  const total = ref(0)

  const fetchJobs = async (params = {}) => {
    loading.value = true
    try {
      const response = await jobsAPI.list({ ...filter.value, ...params })
      jobs.value = response.data
      total.value = response.total || 0
    } catch (error) {
      console.error('Failed to fetch jobs:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  const fetchJob = async (id) => {
    loading.value = true
    try {
      const response = await jobsAPI.get(id)
      currentJob.value = response.data
      return response.data
    } catch (error) {
      console.error('Failed to fetch job:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  const createJob = async (data) => {
    loading.value = true
    try {
      const response = await jobsAPI.create(data)
      await fetchJobs()
      return response.data
    } catch (error) {
      console.error('Failed to create job:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  const updateJob = async (id, data) => {
    loading.value = true
    try {
      const response = await jobsAPI.update(id, data)
      await fetchJobs()
      return response.data
    } catch (error) {
      console.error('Failed to update job:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  const deleteJob = async (id) => {
    loading.value = true
    try {
      await jobsAPI.delete(id)
      await fetchJobs()
    } catch (error) {
      console.error('Failed to delete job:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  const executeJob = async (id) => {
    loading.value = true
    try {
      await jobsAPI.execute(id)
      await fetchJobs()
    } catch (error) {
      console.error('Failed to execute job:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  const cancelJob = async (id) => {
    loading.value = true
    try {
      await jobsAPI.cancel(id)
      await fetchJobs()
    } catch (error) {
      console.error('Failed to cancel job:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  const scheduledJobs = computed(() =>
    jobs.value.filter(job => job.status === 'scheduled')
  )

  const runningJobs = computed(() =>
    jobs.value.filter(job => job.status === 'running')
  )

  return {
    jobs,
    currentJob,
    loading,
    filter,
    total,
    scheduledJobs,
    runningJobs,
    fetchJobs,
    fetchJob,
    createJob,
    updateJob,
    deleteJob,
    executeJob,
    cancelJob,
  }
})
