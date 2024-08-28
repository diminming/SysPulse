<template>
    <a-layout-content :style="{ background: '#fff', padding: '24px', margin: 0, minHeight: '800px' }">
        <a-tabs v-model:activeKey="activeKey" type="card">
            <a-tab-pane key="overview" tab="性能概览" :force-render="false">
                <Overview></Overview>
            </a-tab-pane>
            <a-tab-pane key="topological" tab="资源拓扑" :force-render="false">
                <Topological></Topological>
            </a-tab-pane>
            <a-tab-pane key="profiling" tab="进程剖析" :force-render="false">
                <Profiling></Profiling>
            </a-tab-pane>
            <a-tab-pane key="traffic" tab="流量分析" :force-render="false">
                <Traffic></Traffic>
            </a-tab-pane>
            <a-tab-pane key="summary" tab="基本信息">Content of Tab Pane 3</a-tab-pane>
        </a-tabs>
    </a-layout-content>
</template>

<script lang="ts" setup>
import { ref } from 'vue';
import Overview from "@/views/linux/Overview.vue"
import Profiling from "@/views/linux/Profiling.vue"
import Topological from "@/views/linux/Topological.vue"
import Traffic from '@/views/linux/Traffic.vue';

const activeKey = ref('overview');
</script>

<style scoped></style>