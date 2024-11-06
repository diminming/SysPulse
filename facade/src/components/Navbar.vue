<template>
  <div class="headerLeft">
    <img class="headerImg" :src="ICON_PATH" />
    <a-menu class="headerLeft" v-model:selectedKeys="selectedKeys1" theme="dark" mode="horizontal"
      :style="{ lineHeight: '64px' }">
      <a-menu-item class="headerLeftName"> {{ APP_NAME }} </a-menu-item>
      <a-menu-item v-for="(item, index) in titleList" :key="index" style="font-weight: bold">
        <RouterLink :to="item.href">{{ item.title }}</RouterLink>
      </a-menu-item>
    </a-menu>
  </div>
  <div class="headerRight">
    <span class="headerRightUser">欢迎 {{ username }}</span>
    <a-avatar>
      <template #icon>
        <UserOutlined />
      </template>
    </a-avatar>
  </div>
</template>
 
<script lang="ts" setup>
import { ref, reactive, onMounted } from "vue";
import { UserOutlined } from "@ant-design/icons-vue";
const APP_NAME = import.meta.env.VITE_APP_NAME
const ICON_PATH = import.meta.env.VITE_ICON_PATH
let username = ref("")

onMounted(() => {
  let cache = localStorage.getItem("curr_user")
  if (cache) {
    let user = JSON.parse(cache)
    user['username'] ? username.value = user['username'] : ""
  }
})

const selectedKeys1 = ref<string[]>(["2"]);
const titleList = reactive([
  {
    title: "常规",
    "href": "/main"
  },
  // {
  //   title: "拓扑",
  //   href: "/topo"
  // }
]);
</script>
 
<style>
.headerLeftName {
  font-weight: bold;
  font-size: 20px;
  margin-right: 20px !important;
}

.headerLeft {
  display: flex;
}

.headerImg {
  width: 60px;
  margin-right: 10px;
  object-fit: cover;
}

.headerRightUser {
  color: white;
  font-weight: bold;
  margin-right: 20px;
}

.headerRight {
  display: flex;
  align-items: center;
}
</style>