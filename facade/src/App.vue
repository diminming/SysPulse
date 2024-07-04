<!--<script lang="ts" setup>-->
<!--import { ref, reactive } from 'vue';-->
<!--import { MailOutlined, AppstoreOutlined, SettingOutlined } from '@ant-design/icons-vue';-->

<!--const selectedKeys1 = ref<string[]>(['2']);-->
<!--const selectedKeys2 = ref<string[]>(['1']);-->
<!--const openKeys = ref<string[]>(['sub1']);-->

<!--const titleList = reactive({ list: [{ title: '常规' }, { title: '拓扑' }] })-->
<!--</script>-->

<!--<template>-->
<!--  <a-layout class="full_height">-->
<!--    <a-layout-header class="header">-->
<!--      <div class="logo" />-->
<!--      <a-menu v-model:selectedKeys="selectedKeys1" theme="dark" mode="horizontal" :style="{ lineHeight: '64px' }">-->
<!--        &lt;!&ndash; <a-menu-item key="1">nav 11</a-menu-item>-->
<!--        <a-menu-item key="2">nav 2</a-menu-item>-->
<!--        <a-menu-item key="3">nav 3</a-menu-item> &ndash;&gt;-->

<!--        <a-menu-item v-for="(item, index) in titleList.list" :key="index">{{ item.title }}</a-menu-item>-->
<!--      </a-menu>-->
<!--    </a-layout-header>-->
<!--    <a-layout>-->
<!--      <a-layout-sider width="200" style="background: #fff">-->
<!--        <a-menu v-model:selectedKeys="selectedKeys2" v-model:openKeys="openKeys" mode="inline"-->
<!--          :style="{ height: '100%', borderRight: 0 }">-->
<!--          <a-menu-item key="dashboard">-->
<!--            Dashboard-->
<!--          </a-menu-item>-->
<!--          <a-sub-menu key="sub1">-->
<!--            <template #title>-->
<!--              <span>-->
<!--                <user-outlined />-->
<!--                Dashboard-->
<!--              </span>-->
<!--            </template>-->
<!--            <a-menu-item key="1">option1</a-menu-item>-->
<!--            <a-menu-item key="2">option2</a-menu-item>-->
<!--            <a-menu-item key="3">option3</a-menu-item>-->
<!--            <a-menu-item key="4">option4</a-menu-item>-->
<!--          </a-sub-menu>-->
<!--          <a-sub-menu key="sub2">-->
<!--            <template #title>-->
<!--              <span>-->
<!--                <laptop-outlined />-->
<!--                分组管理-->
<!--              </span>-->
<!--            </template>-->
<!--            <a-menu-item key="5">option5</a-menu-item>-->
<!--            <a-menu-item key="6">option6</a-menu-item>-->
<!--            <a-menu-item key="7">option7</a-menu-item>-->
<!--            <a-menu-item key="8">option8</a-menu-item>-->
<!--          </a-sub-menu>-->
<!--          <a-sub-menu key="sub3">-->
<!--            <template #title>-->
<!--              <span>-->
<!--                <notification-outlined />-->
<!--                Linux-->
<!--              </span>-->
<!--            </template>-->
<!--            <a-menu-item key="9">option9</a-menu-item>-->
<!--            <a-menu-item key="10">option10</a-menu-item>-->
<!--            <a-menu-item key="11">option11</a-menu-item>-->
<!--            <a-menu-item key="12">option12</a-menu-item>-->
<!--          </a-sub-menu>-->
<!--        </a-menu>-->
<!--      </a-layout-sider>-->
<!--      <a-layout style="padding: 0 24px 24px">-->
<!--        <a-breadcrumb style="margin: 16px 0">-->
<!--          <a-breadcrumb-item>Home</a-breadcrumb-item>-->
<!--          <a-breadcrumb-item>List</a-breadcrumb-item>-->
<!--          <a-breadcrumb-item>App</a-breadcrumb-item>-->
<!--        </a-breadcrumb>-->
<!--        <a-layout-content :style="{ background: '#fff', padding: '24px', margin: 0, minHeight: '280px' }">-->
<!--          Content-->
<!--        </a-layout-content>-->
<!--      </a-layout>-->
<!--    </a-layout>-->
<!--  </a-layout>-->
<!--</template>-->

<!--<style scoped>-->
<!--#components-layout-demo-top-side-2 .logo {-->
<!--  float: left;-->
<!--  width: 120px;-->
<!--  height: 31px;-->
<!--  margin: 16px 24px 16px 0;-->
<!--  background: rgba(255, 255, 255, 0.3);-->
<!--}-->

<!--.ant-row-rtl #components-layout-demo-top-side-2 .logo {-->
<!--  float: right;-->
<!--  margin: 16px 0 16px 24px;-->
<!--}-->

<!--.site-layout-background {-->
<!--  background: #fff;-->
<!--}-->
<!--</style>-->


<template>
  <router-view></router-view>
</template>

<script setup lang="ts">
</script>
<style scoped>
</style>