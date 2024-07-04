<template>
  <a-menu v-model:selectedKeys="selectedKeys2" v-model:openKeys="openKeys" mode="inline" >
      <a-menu-item v-for="item in MenuLst" :key="item.name">
        <RouterLink :to="item.href">
          <icon :component="item.icon" />
          {{ item.text }}
        </RouterLink>
      </a-menu-item>
    </a-menu>
</template>

<script lang="ts" setup>
import { ref, reactive } from "vue";
import MenuLst from "./MenuLst"
import Icon from '@ant-design/icons-vue';
const selectedKeys2 = ref<string[]>(["1"]);
const openKeys = ref<string[]>(["sub1"]);
</script>

<style scoped>
#components-layout-demo-top-side-2 .logo {
  float: left;
  width: 120px;
  height: 31px;
  margin: 16px 24px 16px 0;
  background: rgba(255, 255, 255, 0.3);
}

.ant-row-rtl #components-layout-demo-top-side-2 .logo {
  float: right;
  margin: 16px 0 16px 24px;
}

.site-layout-background {
  background: #fff;
}
</style>
