import { defineStore } from "pinia";
import { ref, shallowRef } from "vue";
import { systemAPI, projectsAPI, tagsAPI } from "../utils/api";
import { Cache, dedupe } from "../utils/cache";

const CACHE_KEYS = {
  PROJECTS: "projects",
  TAGS: "tags",
  JOB_TYPES: "job_types",
  STATS: "stats",
  WORKER_STATUS: "worker_status",
};

const CACHE_TTL = {
  STATIC: 30 * 60 * 1000, // 30 minutes for static data
  DYNAMIC: 30 * 1000, // 30 seconds for dynamic data
};

export const useSystemStore = defineStore("system", () => {
  const stats = ref(null);
  const workerStatus = ref(null);
  const projects = shallowRef([]);
  const tags = shallowRef([]);
  const jobTypes = shallowRef([]);
  const loading = ref(false);
  const initialized = ref(false);

  const fetchStats = async (useCache = true) => {
    if (useCache) {
      const cached = Cache.get(CACHE_KEYS.STATS);
      if (cached) {
        stats.value = cached;
        return;
      }
    }

    try {
      const response = await dedupe("fetchStats", () => systemAPI.status());
      stats.value = response.data;
      Cache.set(CACHE_KEYS.STATS, response.data, CACHE_TTL.DYNAMIC);
    } catch (error) {
      console.error("Failed to fetch stats:", error);
    }
  };

  const fetchWorkerStatus = async (useCache = true) => {
    if (useCache) {
      const cached = Cache.get(CACHE_KEYS.WORKER_STATUS);
      if (cached) {
        workerStatus.value = cached;
        return;
      }
    }

    try {
      const response = await dedupe("fetchWorkerStatus", () =>
        systemAPI.workerStatus(),
      );
      workerStatus.value = response.data;
      Cache.set(CACHE_KEYS.WORKER_STATUS, response.data, CACHE_TTL.DYNAMIC);
    } catch (error) {
      console.error("Failed to fetch worker status:", error);
    }
  };

  const fetchProjects = async (useCache = true) => {
    if (useCache) {
      const cached = Cache.get(CACHE_KEYS.PROJECTS);
      if (cached) {
        projects.value = cached;
        return;
      }
    }

    try {
      const response = await dedupe("fetchProjects", () => projectsAPI.list());
      projects.value = response.data;
      Cache.set(CACHE_KEYS.PROJECTS, response.data, CACHE_TTL.STATIC);
    } catch (error) {
      console.error("Failed to fetch projects:", error);
    }
  };

  const fetchTags = async (useCache = true) => {
    if (useCache) {
      const cached = Cache.get(CACHE_KEYS.TAGS);
      if (cached) {
        tags.value = cached;
        return;
      }
    }

    try {
      const response = await dedupe("fetchTags", () => tagsAPI.list());
      tags.value = response.data;
      Cache.set(CACHE_KEYS.TAGS, response.data, CACHE_TTL.STATIC);
    } catch (error) {
      console.error("Failed to fetch tags:", error);
    }
  };

  const fetchJobTypes = async (useCache = true) => {
    if (useCache) {
      const cached = Cache.get(CACHE_KEYS.JOB_TYPES);
      if (cached) {
        jobTypes.value = cached;
        return;
      }
    }

    try {
      const response = await dedupe("fetchJobTypes", () =>
        systemAPI.jobTypes(),
      );
      jobTypes.value = response.data;
      Cache.set(CACHE_KEYS.JOB_TYPES, response.data, CACHE_TTL.STATIC);
    } catch (error) {
      console.error("Failed to fetch job types:", error);
    }
  };

  const initializeApp = async (force = false) => {
    if (initialized.value && !force) {
      return;
    }

    loading.value = true;
    try {
      await Promise.all([
        fetchProjects(),
        fetchTags(),
        fetchJobTypes(),
        fetchStats(),
        fetchWorkerStatus(),
      ]);
      initialized.value = true;
    } finally {
      loading.value = false;
    }
  };

  const invalidateCache = (keys = null) => {
    if (keys) {
      keys.forEach((key) => Cache.remove(key));
    } else {
      Object.values(CACHE_KEYS).forEach((key) => Cache.remove(key));
    }
  };

  const refreshData = async () => {
    invalidateCache();
    initialized.value = false;
    await initializeApp(true);
  };

  return {
    stats,
    workerStatus,
    projects,
    tags,
    jobTypes,
    loading,
    initialized,
    fetchStats,
    fetchWorkerStatus,
    fetchProjects,
    fetchTags,
    fetchJobTypes,
    initializeApp,
    invalidateCache,
    refreshData,
  };
});
