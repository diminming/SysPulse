<template>
  <div class="tile">
    <div class="tile-icon" v-bind:style="{ backgroundColor: color }">
      <icon class="title-icon-img" :component="icon" />
    </div>
    <div class="text">
      <div class="num" :style="{ color: color }">{{ value }}</div>
      <div class="title">{{ title }}</div>
    </div>
    <div class="handler">
      <span> 点击查看-> </span>
    </div>
  </div>
</template>
  
<script setup lang="ts">
import Icon from "@ant-design/icons-vue";
import { BorderOutlined } from "@ant-design/icons-vue";

defineProps({
  color: {
    type: String,
    default: "#ebedef",
  },
  icon: {
    default: BorderOutlined,
  },
  title: {
    type: String,
    default: "<标题>",
  },
  value: {
    type: Number,
    default: 123456789,
  },
});
</script>

<style  scoped>
.tile {
  width: 290px;
  height: 135px;
  border: 1px solid #e4e4e4;
  border-radius: 10px;
  position: relative;
}

.tile-icon {
  width: 60px;
  height: 60px;
  position: absolute;
  top: 20px;
  left: 20px;
  font-size: 48px;
  /* 弹性盒 */
  display: flex;
  /*水平居中 */
  justify-content: center; 
  /*竖直居中 */
  align-items: center;
  color: #fff;
}

.tile-icon > span {
  width: 48px;
  height: 48px;
}


:deep(.anticon) {
  width: 48px;
  height: 48px;
}

.tile .handler {
  position: absolute;
  width: 100%;
  text-align: right;
  color: #2c384a;
  padding-right: 10px;
  right: 0;
  bottom: 0;
  height: 32px;
  line-height: 32px;
  font-size: 18px;
  background-color: #e4e4e4;
}

.tile .handler span {
  cursor: pointer;
  opacity: 68%;
}

.tile .text {
  position: absolute;
  top: 20px;
  left: 100px;
}

.title {
  font-size: 20px;
  margin: 0.5em 0.5em 0 0;
  color: #2c384a;
  opacity: 68%;
  height: 30px;
}

.num {
  font-size: 22px;
  font-weight: bold;
  height: 30px;
}
</style>