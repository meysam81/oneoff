import { defineStore } from "pinia";
import { ref, computed, shallowRef } from "vue";
import { jobsAPI } from "../utils/api";
import { Cache, dedupe } from "../utils/cache";

const CACHE_KEY_PREFIX = "jobs_";
const JOB_CACHE_TTL = 30 * 1000; // 30 seconds

export const useJobsStore = defineStore("jobs", () => {
  const jobs = shallowRef([]);
  const currentJob = ref(null);
  const loading = ref(false);
  const filter = ref({
    project_id: "",
    status: "",
    type: "",
    search: "",
    tag_ids: [],
    sort_by: "scheduled_at",
    sort_order: "asc",
    limit: 50,
    offset: 0,
  });
  const total = ref(0);

  const getCacheKey = (params) => {
    return CACHE_KEY_PREFIX + JSON.stringify(params);
  };

  const fetchJobs = async (params = {}, useCache = true) => {
    const queryParams = { ...filter.value, ...params };
    const cacheKey = getCacheKey(queryParams);

    if (useCache) {
      const cached = Cache.get(cacheKey);
      if (cached) {
        jobs.value = cached.data;
        total.value = cached.total || 0;
        return;
      }
    }

    loading.value = true;
    try {
      const response = await dedupe(cacheKey, () => jobsAPI.list(queryParams));
      jobs.value = response.data;
      total.value = response.total || 0;
      Cache.set(cacheKey, response, JOB_CACHE_TTL);
    } catch (error) {
      console.error("Failed to fetch jobs:", error);
      throw error;
    } finally {
      loading.value = false;
    }
  };

  const fetchJob = async (id) => {
    loading.value = true;
    try {
      const response = await jobsAPI.get(id);
      currentJob.value = response.data;
      return response.data;
    } catch (error) {
      console.error("Failed to fetch job:", error);
      throw error;
    } finally {
      loading.value = false;
    }
  };

  const invalidateJobsCache = () => {
    Object.keys(localStorage)
      .filter((key) => key.startsWith("oneoff_cache_" + CACHE_KEY_PREFIX))
      .forEach((key) => localStorage.removeItem(key));
  };

  const createJob = async (data) => {
    loading.value = true;
    try {
      const response = await jobsAPI.create(data);
      invalidateJobsCache();
      await fetchJobs({}, false);
      return response.data;
    } catch (error) {
      console.error("Failed to create job:", error);
      throw error;
    } finally {
      loading.value = false;
    }
  };

  const updateJob = async (id, data) => {
    loading.value = true;
    try {
      const response = await jobsAPI.update(id, data);
      invalidateJobsCache();
      await fetchJobs({}, false);
      return response.data;
    } catch (error) {
      console.error("Failed to update job:", error);
      throw error;
    } finally {
      loading.value = false;
    }
  };

  const deleteJob = async (id) => {
    loading.value = true;
    try {
      await jobsAPI.delete(id);
      invalidateJobsCache();
      await fetchJobs({}, false);
    } catch (error) {
      console.error("Failed to delete job:", error);
      throw error;
    } finally {
      loading.value = false;
    }
  };

  const executeJob = async (id) => {
    loading.value = true;
    try {
      await jobsAPI.execute(id);
      invalidateJobsCache();
      await fetchJobs({}, false);
    } catch (error) {
      console.error("Failed to execute job:", error);
      throw error;
    } finally {
      loading.value = false;
    }
  };

  const cancelJob = async (id) => {
    loading.value = true;
    try {
      await jobsAPI.cancel(id);
      invalidateJobsCache();
      await fetchJobs({}, false);
    } catch (error) {
      console.error("Failed to cancel job:", error);
      throw error;
    } finally {
      loading.value = false;
    }
  };

  const scheduledJobs = computed(() =>
    jobs.value.filter((job) => job.status === "scheduled"),
  );

  const runningJobs = computed(() =>
    jobs.value.filter((job) => job.status === "running"),
  );

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
  };
});
