import { createRouter, createWebHistory } from "vue-router";

const routes = [
  {
    path: "/",
    name: "Dashboard",
    component: () =>
      import(/* webpackChunkName: "dashboard" */ "./views/Dashboard.vue"),
    meta: { keepAlive: true },
  },
  {
    path: "/jobs",
    name: "Jobs",
    component: () => import(/* webpackChunkName: "jobs" */ "./views/Jobs.vue"),
    meta: { keepAlive: true },
  },
  {
    path: "/jobs/:id",
    name: "JobDetails",
    component: () =>
      import(/* webpackChunkName: "job-details" */ "./views/JobDetails.vue"),
    meta: { keepAlive: false },
  },
  {
    path: "/executions",
    name: "Executions",
    component: () =>
      import(/* webpackChunkName: "executions" */ "./views/Executions.vue"),
    meta: { keepAlive: true },
  },
  {
    path: "/projects",
    name: "Projects",
    component: () =>
      import(/* webpackChunkName: "projects" */ "./views/Projects.vue"),
    meta: { keepAlive: true },
  },
  {
    path: "/settings",
    name: "Settings",
    component: () =>
      import(/* webpackChunkName: "settings" */ "./views/Settings.vue"),
    meta: { keepAlive: true },
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition;
    } else {
      return { top: 0 };
    }
  },
});

// Prefetch strategy: prefetch likely next routes
router.beforeEach((to, from, next) => {
  // Prefetch Jobs view when on Dashboard
  if (to.name === "Dashboard") {
    import("./views/Jobs.vue");
  }
  next();
});

export default router;
