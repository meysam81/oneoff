import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'Dashboard',
    component: () => import('./views/Dashboard.vue'),
  },
  {
    path: '/jobs',
    name: 'Jobs',
    component: () => import('./views/Jobs.vue'),
  },
  {
    path: '/jobs/:id',
    name: 'JobDetails',
    component: () => import('./views/JobDetails.vue'),
  },
  {
    path: '/executions',
    name: 'Executions',
    component: () => import('./views/Executions.vue'),
  },
  {
    path: '/projects',
    name: 'Projects',
    component: () => import('./views/Projects.vue'),
  },
  {
    path: '/settings',
    name: 'Settings',
    component: () => import('./views/Settings.vue'),
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
