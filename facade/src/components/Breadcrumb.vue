<template>
  <div style="background-color: white; padding: 0 0; margin-top: 0.5rem;">
    <a-space>

      <a-button type="text" @click="toggleCollapsed">
        <MenuUnfoldOutlined v-if="state.collapsed" />
        <MenuFoldOutlined v-else />
      </a-button>

      <a-breadcrumb style="margin: 1rem 0; height: 1rem;">
        <a-breadcrumb-item v-for="(item, index) in breadcrumbs" :key="index">
          <template v-if="item.isActive">
            <span>{{ item.title }}</span>
          </template>
          <template v-else>
            <RouterLink :to="item.url">
              <span>{{ item.title }}</span>
            </RouterLink>
          </template>
        </a-breadcrumb-item>
      </a-breadcrumb>

    </a-space>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from "vue";
import { useRouter } from "vue-router";
import { sidebarCollapsed } from "@/stores/sidebarStatus"

import { MenuFoldOutlined, MenuUnfoldOutlined } from '@ant-design/icons-vue';

const breadcrumbs = ref(),
  router = useRouter(),
  state = sidebarCollapsed(),
  getBreadcrumbs = () => {
    return router.currentRoute.value.matched.slice(1).map((route) => {
      return {
        isActive: route.path === router.currentRoute.value.fullPath,
        title: route.meta.text,
        url: `${router.options.history.base}${route.path}`,
      }
    })
  },
  toggleCollapsed = () => {
    state.toggle()
  }

router.afterEach(() => {
  breadcrumbs.value = getBreadcrumbs()
})

onMounted(() => {
  breadcrumbs.value = getBreadcrumbs()
})

</script>
<style scoped>
.nowPath {
  color: black;
}
</style>