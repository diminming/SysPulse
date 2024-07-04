<template>
  <a-layout-content :style="{
    background: '#fff',
    padding: '24px',
    margin: 0,
    minHeight: '800px',
  }">
    <a-descriptions title="Linux资源详情" layout="vertical" bordered>
      <a-descriptions-item label="主机名称">{{  }}</a-descriptions-item>
      <a-descriptions-item label="ID">{{  }}</a-descriptions-item>
      <a-descriptions-item label="主机标识">{{  }}</a-descriptions-item>
      <a-descriptions-item label="运行状态" :span="3">
        <a-badge status="processing" text="Running" />
      </a-descriptions-item>
      <a-descriptions-item label="创建时间">{{  }}</a-descriptions-item>
      <a-descriptions-item label="更新时间" :span="2">{{  }}</a-descriptions-item>
      <a-descriptions-item label="系统名称" :span="3">{{  }}</a-descriptions-item>

      <a-descriptions-item label="描述信息">
        Data disk type: MongoDB
        <br />
        Database version: 3.4
        <br />
        Package: dds.mongo.mid
        <br />
        Storage space: 10 GB
        <br />
        Replication factor: 3
        <br />
        Region: East China 1
        <br />
      </a-descriptions-item>
    </a-descriptions>
  </a-layout-content>
</template>

<script setup lang="ts">
import { getCurrentInstance, reactive, onMounted } from 'vue'

const linuxInfo = reactive({ obj: {} })
onMounted(() => {
  console.log('点击数据的',linuxInfo.obj)
});
</script>

<style scoped></style>