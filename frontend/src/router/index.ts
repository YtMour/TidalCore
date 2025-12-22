import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/store/user'

const router = createRouter({
  history: createWebHistory(),
  scrollBehavior(_to, _from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    }
    return { top: 0, behavior: 'smooth' }
  },
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('@/views/Home.vue'),
      meta: { title: 'TidalCore - 开源盆底肌训练平台' }
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/Login.vue'),
      meta: { guest: true, title: '登录 - TidalCore' }
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('@/views/Register.vue'),
      meta: { guest: true, title: '注册 - TidalCore' }
    },
    {
      path: '/train',
      name: 'train',
      component: () => import('@/views/Train.vue'),
      meta: { title: '训练 - TidalCore' }
    },
    {
      path: '/dashboard',
      name: 'dashboard',
      component: () => import('@/views/Dashboard.vue'),
      meta: { requiresAuth: true, title: '用户中心 - TidalCore' }
    },
    {
      path: '/leaderboard',
      name: 'leaderboard',
      component: () => import('@/views/Leaderboard.vue'),
      meta: { title: '排行榜 - TidalCore' }
    },
    {
      path: '/about',
      name: 'about',
      component: () => import('@/views/About.vue'),
      meta: { title: '关于凯格尔运动 - TidalCore' }
    },
    {
      path: '/admin',
      component: () => import('@/views/admin/AdminLayout.vue'),
      meta: { requiresAuth: true, requiresAdmin: true },
      children: [
        {
          path: '',
          redirect: '/admin/dashboard'
        },
        {
          path: 'dashboard',
          name: 'admin-dashboard',
          component: () => import('@/views/admin/Dashboard.vue'),
          meta: { title: '仪表盘 - 管理后台 - TidalCore' }
        },
        {
          path: 'users',
          name: 'admin-users',
          component: () => import('@/views/admin/Users.vue'),
          meta: { title: '用户管理 - 管理后台 - TidalCore' }
        },
        {
          path: 'checkins',
          name: 'admin-checkins',
          component: () => import('@/views/admin/Checkins.vue'),
          meta: { title: '打卡数据 - 管理后台 - TidalCore' }
        },
        {
          path: 'settings',
          name: 'admin-settings',
          component: () => import('@/views/admin/Settings.vue'),
          meta: { title: '系统设置 - 管理后台 - TidalCore' }
        }
      ]
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'not-found',
      redirect: '/'
    }
  ]
})

router.beforeEach((to, _from, next) => {
  const userStore = useUserStore()

  // Update page title
  if (to.meta.title) {
    document.title = to.meta.title as string
  }

  if (to.meta.requiresAuth && !userStore.isLoggedIn) {
    next({ name: 'login', query: { redirect: to.fullPath } })
  } else if (to.meta.guest && userStore.isLoggedIn) {
    next({ name: 'home' })
  } else {
    next()
  }
})

export default router
