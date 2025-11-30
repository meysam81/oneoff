<template>
  <div class="projects-view">
    <n-space vertical :size="16">
      <n-card>
        <n-space justify="space-between">
          <n-space>
            <n-button type="primary" @click="showCreateModal = true">
              <template #icon>
                <n-icon><AddOutline /></n-icon>
              </template>
              Create Project
            </n-button>

            <n-button @click="fetchProjects">
              <template #icon>
                <n-icon><RefreshOutline /></n-icon>
              </template>
              Refresh
            </n-button>
          </n-space>

          <n-space>
            <n-input
              v-model:value="searchQuery"
              placeholder="Search projects..."
              clearable
              style="width: 300px"
            >
              <template #prefix>
                <n-icon><SearchOutline /></n-icon>
              </template>
            </n-input>

            <n-checkbox v-model:checked="showArchived">
              Show Archived
            </n-checkbox>
          </n-space>
        </n-space>
      </n-card>

      <n-card>
        <n-spin :show="loading">
          <n-table :bordered="false" :single-line="false">
            <thead>
              <tr>
                <th>Name</th>
                <th>Description</th>
                <th>Jobs</th>
                <th>Status</th>
                <th>Created</th>
                <th style="width: 120px">Actions</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="project in filteredProjects" :key="project.id">
                <td>
                  <n-space align="center">
                    <div
                      class="color-indicator"
                      :style="{ backgroundColor: project.color || '#6366f1' }"
                    ></div>
                    <n-text strong>{{ project.name }}</n-text>
                    <n-tag v-if="project.id === 'default'" type="info" size="small">
                      Default
                    </n-tag>
                  </n-space>
                </td>
                <td>
                  <n-text depth="3">{{ project.description || '-' }}</n-text>
                </td>
                <td>
                  <n-tag type="info">{{ project.job_count || 0 }} jobs</n-tag>
                </td>
                <td>
                  <n-tag :type="project.is_archived ? 'warning' : 'success'">
                    {{ project.is_archived ? 'Archived' : 'Active' }}
                  </n-tag>
                </td>
                <td>
                  <n-text depth="3">{{ formatDate(project.created_at) }}</n-text>
                </td>
                <td>
                  <n-space>
                    <n-button
                      quaternary
                      size="small"
                      @click="editProject(project)"
                    >
                      <template #icon>
                        <n-icon><CreateOutline /></n-icon>
                      </template>
                    </n-button>
                    <n-button
                      quaternary
                      size="small"
                      :disabled="project.id === 'default'"
                      @click="confirmDelete(project)"
                    >
                      <template #icon>
                        <n-icon><TrashOutline /></n-icon>
                      </template>
                    </n-button>
                  </n-space>
                </td>
              </tr>
              <tr v-if="filteredProjects.length === 0">
                <td colspan="6">
                  <n-empty description="No projects found" />
                </td>
              </tr>
            </tbody>
          </n-table>
        </n-spin>
      </n-card>
    </n-space>

    <!-- Create/Edit Modal -->
    <n-modal
      v-model:show="showCreateModal"
      preset="dialog"
      :title="editingProject ? 'Edit Project' : 'Create Project'"
      :positive-text="editingProject ? 'Update' : 'Create'"
      negative-text="Cancel"
      @positive-click="handleSubmit"
      @negative-click="closeModal"
    >
      <n-form ref="formRef" :model="formData" :rules="formRules">
        <n-form-item label="Name" path="name">
          <n-input
            v-model:value="formData.name"
            placeholder="Enter project name"
            :disabled="editingProject?.id === 'default'"
          />
        </n-form-item>

        <n-form-item label="Description" path="description">
          <n-input
            v-model:value="formData.description"
            type="textarea"
            placeholder="Enter project description"
            :rows="3"
          />
        </n-form-item>

        <n-form-item label="Color" path="color">
          <n-color-picker
            v-model:value="formData.color"
            :swatches="colorSwatches"
          />
        </n-form-item>

        <n-form-item v-if="editingProject" label="Status" path="is_archived">
          <n-switch
            v-model:value="formData.is_archived"
            :disabled="editingProject?.id === 'default'"
          >
            <template #checked>Archived</template>
            <template #unchecked>Active</template>
          </n-switch>
        </n-form-item>
      </n-form>
    </n-modal>

    <!-- Delete Confirmation Modal -->
    <n-modal
      v-model:show="showDeleteModal"
      preset="dialog"
      title="Delete Project"
      type="warning"
      positive-text="Delete"
      negative-text="Cancel"
      @positive-click="handleDelete"
      @negative-click="showDeleteModal = false"
    >
      <n-space vertical>
        <n-text>
          Are you sure you want to delete the project "{{ projectToDelete?.name }}"?
        </n-text>
        <n-alert type="warning" :show-icon="false">
          This action cannot be undone. All jobs in this project will be moved to
          the default project.
        </n-alert>
      </n-space>
    </n-modal>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, reactive } from "vue";
import {
  AddOutline,
  RefreshOutline,
  SearchOutline,
  CreateOutline,
  TrashOutline,
} from "@vicons/ionicons5";
import { useMessage } from "naive-ui";
import { projectsAPI } from "../utils/api";
import { useSystemStore } from "../stores/system";
import { format } from "date-fns";

const message = useMessage();
const systemStore = useSystemStore();

const loading = ref(false);
const searchQuery = ref("");
const showArchived = ref(false);
const showCreateModal = ref(false);
const showDeleteModal = ref(false);
const editingProject = ref(null);
const projectToDelete = ref(null);
const formRef = ref(null);

const projects = computed(() => systemStore.projects);

const filteredProjects = computed(() => {
  let result = projects.value || [];

  // Filter by archived status
  if (!showArchived.value) {
    result = result.filter((p) => !p.is_archived);
  }

  // Filter by search query
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase();
    result = result.filter(
      (p) =>
        p.name.toLowerCase().includes(query) ||
        (p.description && p.description.toLowerCase().includes(query))
    );
  }

  return result;
});

const colorSwatches = [
  "#6366f1", // Indigo
  "#3b82f6", // Blue
  "#22c55e", // Green
  "#f59e0b", // Amber
  "#ef4444", // Red
  "#ec4899", // Pink
  "#8b5cf6", // Purple
  "#06b6d4", // Cyan
];

const formData = reactive({
  name: "",
  description: "",
  color: "#6366f1",
  is_archived: false,
});

const formRules = {
  name: [
    { required: true, message: "Name is required", trigger: "blur" },
    { min: 1, max: 100, message: "Name must be 1-100 characters", trigger: "blur" },
  ],
};

const formatDate = (date) => {
  if (!date) return "-";
  return format(new Date(date), "MMM d, yyyy HH:mm");
};

const fetchProjects = async () => {
  loading.value = true;
  try {
    await systemStore.fetchProjects(false);
  } catch (error) {
    message.error(`Failed to fetch projects: ${error?.message || 'Unknown error'}`);
  } finally {
    loading.value = false;
  }
};

const resetForm = () => {
  formData.name = "";
  formData.description = "";
  formData.color = "#6366f1";
  formData.is_archived = false;
  editingProject.value = null;
};

const closeModal = () => {
  showCreateModal.value = false;
  resetForm();
};

const editProject = (project) => {
  editingProject.value = project;
  formData.name = project.name;
  formData.description = project.description || "";
  formData.color = project.color || "#6366f1";
  formData.is_archived = project.is_archived || false;
  showCreateModal.value = true;
};

const handleSubmit = async () => {
  try {
    await formRef.value?.validate();
  } catch {
    return false;
  }

  try {
    if (editingProject.value) {
      // Update existing project
      await projectsAPI.update(editingProject.value.id, {
        name: formData.name,
        description: formData.description,
        color: formData.color,
        is_archived: formData.is_archived,
      });
      message.success("Project updated successfully");
    } else {
      // Create new project
      await projectsAPI.create({
        name: formData.name,
        description: formData.description,
        color: formData.color,
      });
      message.success("Project created successfully");
    }

    closeModal();
    await fetchProjects();
  } catch (error) {
    message.error(error.message || "Failed to save project");
    return false;
  }
};

const confirmDelete = (project) => {
  projectToDelete.value = project;
  showDeleteModal.value = true;
};

const handleDelete = async () => {
  if (!projectToDelete.value) return;

  try {
    await projectsAPI.delete(projectToDelete.value.id);
    message.success("Project deleted successfully");
    showDeleteModal.value = false;
    projectToDelete.value = null;
    await fetchProjects();
  } catch (error) {
    message.error(error.message || "Failed to delete project");
    return false;
  }
};

onMounted(async () => {
  await fetchProjects();
});
</script>

<style scoped>
.projects-view {
  padding: 0;
}

.color-indicator {
  width: 12px;
  height: 12px;
  border-radius: 3px;
}
</style>
