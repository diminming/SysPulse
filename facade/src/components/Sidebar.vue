<template>
  <a-menu 
    v-model:selectedKeys="selectedKeys2" 
    v-model:openKeys="openKeys" 
    mode="inline" 
    :items="items"
    @click="handleClick"
  >
  </a-menu>
</template>

<script lang="ts" setup>
import { reactive, ref, watch, VueElement, h } from 'vue';
import {useRouter} from "vue-router"
import type { MenuProps, ItemType } from 'ant-design-vue';
import MenuLst from "./MenuLst"

const selectedKeys2 = ref<string[]>(["1"]);
const openKeys = ref<string[]>(["sub1"]);
const router = useRouter()

function getItem(label: VueElement | string, key: string, icon?: any, children?: ItemType[], type?: 'group', href?: string): ItemType {
  let subMenu = undefined

  if(children) {
    subMenu = []
    for(let i = 0; i < children.length; i++) {
      const child = children[i]
      subMenu.push(getItem(child.text, child.name, h(child.icon), child.children, child.type, child.href))
    }
  }
  
  return {
    key,
    icon,
    label,
    href,
    type,
    children: subMenu
  } as ItemType;
}

const handleClick: MenuProps['onClick'] = e => {
  const item = e.item
  router.push(item.href)
};

const items: ItemType[] = reactive(
  MenuLst.map(item => {
    return getItem(
      item.text, 
      item.name, 
      h(item.icon), 
      item.children, 
      item.type, 
      item.href
    )
  })
)

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
