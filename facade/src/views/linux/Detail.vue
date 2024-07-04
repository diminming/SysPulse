<template>
    <a-layout-content :style="{ background: '#fff', padding: '24px', margin: 0, minHeight: '800px' }">
        <a-tabs v-model:activeKey="activeKey">
            <a-tab-pane key="overview" tab="概览">
                <Overview></Overview>
            </a-tab-pane>
            <a-tab-pane key="profiling" tab="进程剖析" force-render>
                <Profiling></Profiling>
            </a-tab-pane>
            <a-tab-pane key="summary" tab="信息">Content of Tab Pane 3</a-tab-pane>
        </a-tabs>
    </a-layout-content>
</template>

<script lang="ts" setup>
import { ref } from 'vue';
import type { Linux } from "@/views/linux/linux"
import Overview from "@/views/linux/Overview.vue"
import Profiling from "@/views/linux/Profiling.vue"

const activeKey = ref('overview');
</script>

<style scoped></style>