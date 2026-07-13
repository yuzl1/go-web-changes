import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      redirect: '/tasks',
    },
    {
      path: '/tasks',
      name: 'TaskList',
      component: () => import('../views/TaskList.vue'),
    },
    {
      path: '/tasks/new',
      name: 'TaskCreate',
      component: () => import('../views/TaskForm.vue'),
    },
    {
      path: '/tasks/:id/edit',
      name: 'TaskEdit',
      component: () => import('../views/TaskForm.vue'),
      props: true,
    },
    {
      path: '/tasks/:id/records',
      name: 'RecordList',
      component: () => import('../views/RecordList.vue'),
      props: true,
    },
    {
      path: '/config',
      name: 'SysConfig',
      component: () => import('../views/SysConfig.vue'),
    },
  ],
})

export default router
