import { defineStore } from "pinia";
import { ref } from "vue";
import { systemAPI, projectsAPI, tagsAPI } from "../utils/api";

export const useSystemStore = defineStore("system", () => {
  const stats = ref(null);
  const workerStatus = ref(null);
  const projects = ref([]);
  const tags = ref([]);
  const jobTypes = ref([]);
  const loading = ref(false);

  const fetchStats = async () => {
    try {
      const response = await systemAPI.status();
      stats.value = response.data;
    } catch (error) {
      console.error("Failed to fetch stats:", error);
    }
  };

  const fetchWorkerStatus = async () => {
    try {
      const response = await systemAPI.workerStatus();
      workerStatus.value = response.data;
    } catch (error) {
      console.error("Failed to fetch worker status:", error);
    }
  };

  const fetchProjects = async () => {
    try {
      const response = await projectsAPI.list();
      projects.value = response.data;
    } catch (error) {
      console.error("Failed to fetch projects:", error);
    }
  };

  const fetchTags = async () => {
    try {
      const response = await tagsAPI.list();
      tags.value = response.data;
    } catch (error) {
      console.error("Failed to fetch tags:", error);
    }
  };

  const fetchJobTypes = async () => {
    try {
      const response = await systemAPI.jobTypes();
      jobTypes.value = response.data;
    } catch (error) {
      console.error("Failed to fetch job types:", error);
    }
  };

  const initializeApp = async () => {
    loading.value = true;
    try {
      await Promise.all([
        fetchProjects(),
        fetchTags(),
        fetchJobTypes(),
        fetchStats(),
        fetchWorkerStatus(),
      ]);
    } finally {
      loading.value = false;
    }
  };

  return {
    stats,
    workerStatus,
    projects,
    tags,
    jobTypes,
    loading,
    fetchStats,
    fetchWorkerStatus,
    fetchProjects,
    fetchTags,
    fetchJobTypes,
    initializeApp,
  };
});
