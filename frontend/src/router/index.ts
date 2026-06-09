import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { useUserStore } from '@/stores/user'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    component: () => import('@/layouts/DefaultLayout.vue'),
    redirect: '/dashboard',
    meta: { requiresAuth: true },
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/Index.vue'),
        meta: { title: '仪表盘', icon: 'Monitor' },
      },
      {
        path: 'users',
        name: 'Users',
        component: () => import('@/views/user/Index.vue'),
        meta: { title: '用户管理', icon: 'User', requiresAdmin: true },
      },
    ],
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/Index.vue'),
    meta: { hidden: true },
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/error/404.vue'),
    meta: { hidden: true },
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

const whiteList = ['/login']

router.beforeEach(async (to, _from, next) => {
  const userStore = useUserStore()

  if (whiteList.includes(to.path)) {
    if (userStore.isLoggedIn) {
      next('/dashboard')
    } else {
      next()
    }
    return
  }

  if (!userStore.isLoggedIn) {
    next('/login')
    return
  }

  if (!userStore.userInfo) {
    await userStore.fetchUserInfo()
  }

  if (to.meta.requiresAdmin && !userStore.isAdmin) {
    next('/dashboard')
    return
  }

  next()
})

export default router
