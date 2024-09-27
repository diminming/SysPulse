import { createRouter, createWebHistory } from "vue-router";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      name: "Login",
      path: "/login",
      meta: {
        text: "用户登录",
      },
      component: () => import("@/views/login/Login.vue"),
    },
    {
      path: "/",
      name: "Main",
      redirect: "/main",
      meta: {
        text: import.meta.env.VITE_APP_NAME,
      },
      component: () => import("@/layout/MainLayout.vue"),
      children: [
        {
          name: "general",
          path: "/main",
          redirect: "/main/dashboard",
          component: () => import("@/layout/General.vue"),
          meta: {
            text: "常规",
          },
          children: [
            {
              path: "/main/dashboard",
              name: "dashboard",
              component: () => import("../views/dashboard/Dashboard.vue"),
              meta: {
                text: "仪表盘",
              },
            },
            {
              path: "/main/biz",
              name: "bizMgr",
              component: () => import("../views/biz/Business.vue"),
              meta: {
                text: "资源分组",
              },
              children: [],
            },
            {
              path: "/main/biz/:bizId/detail",
              name: "bizSystem",
              component: () => import("../views/biz/view/BizSystem.vue"),
              meta: {
                text: "系统详情",
              },
            },
            {
              path: "/main/biz/bizAdd",
              name: "bizAdd",
              component: () => import("../views/biz/view/BizAdd.vue"),
              meta: {
                text: "新增业务系统",
              },
            },
            {
              path: "/main/linux",
              name: "linux",
              component: () => import("../views/linux/Linux.vue"),
              meta: {
                text: "Linux资源列表",
              },
            },
            {
              path: "/main/linux/linuxSystem",
              name: "linuxSystem",
              component: () => import("../views/linux/view/LinuxSystem.vue"),
              meta: {
                text: "Linux资源详情",
              },
            },
            {
              path: "/main/linux/add",
              name: "linuxAdd",
              component: () => import("@/views/linux/LinuxEdit.vue"),
              meta: {
                text: "新增",
              },
            },
            {
              path: "/main/linux/edit",
              name: "linuxEdit",
              component: () => import("@/views/linux/LinuxEdit.vue"),
              meta: {
                text: "编辑",
              },
            },
            {
              path: "/main/linux/:linuxId/detail",
              name: "linuxDetail",
              component: () => import("../views/linux/Detail.vue"),
              meta: {
                text: "Linux详情",
              },
            },
            {
              path: "/main/database",
              name: "database",
              component: () => import("../views/database/Database.vue"),
              meta: {
                text: "数据库",
              },
            },
            {
              path: "/main/database/dbSystem",
              name: "dbSystem",
              component: () => import("../views/database/view/DBSystem.vue"),
              meta: {
                text: "数据库详情",
              },
            },
            {
              path: "/main/database/dbAdd",
              name: "dbAdd",
              component: () => import("../views/database/view/DBAdd.vue"),
              meta: {
                text: "新增",
              },
            },
            {
              path: "/main/cache",
              name: "cache",
              component: () => import("../views/cache/Cache.vue"),
              meta: {
                text: "缓存",
              },
            },
            {
              path: "/main/queue",
              name: "queue",
              component: () => import("../views/queue/Queue.vue"),
              meta: {
                text: "队列",
              },
            },
            {
              path: "/main/users",
              name: "users",
              component: () => import("../views/users/Users.vue"),
              meta: {
                text: "用户管理",
              },
            },
            {
              path: "/main/setting",
              name: "setting",
              component: () => import("../views/setting/Setting.vue"),
              meta: {
                text: "系统设置",
              },
            },
            {
              path: "/main/notification",
              name: "notification",
              component: () => import("../views/notification/AlarmLst.vue"),
              meta: {
                text: "消息中心",
              },
            },
          ],
        },
        {
          path: "/topo",
          name: "topo",
          component: () => import("@/layout/Topological.vue"),
          meta: {
            text: "拓扑图",
          },
        },
      ],
    },
  ],
});

export default router;
