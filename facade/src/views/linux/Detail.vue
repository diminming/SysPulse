<template>
    <a-layout-content :style="{ background: '#fff', padding: '24px', margin: 0, minHeight: '800px' }">
        <a-tabs v-model:activeKey="activeKey" type="card">
            <a-tab-pane key="overview" tab="性能概览" :force-render="false">
                <Overview></Overview>
            </a-tab-pane>
            <a-tab-pane key="topological" tab="资源拓扑" :force-render="false">
                <Topological></Topological>
            </a-tab-pane>
            <template v-if="SHOW_TRFFIC">
                <a-tab-pane key="profiling" tab="进程剖析" :force-render="false">
                    <Profiling></Profiling>
                </a-tab-pane>
            </template>
            <template v-if="SHOW_PROFILING">
                <a-tab-pane key="traffic" tab="流量分析" :force-render="false">
                    <Traffic></Traffic>
                </a-tab-pane>
            </template>
            <template v-if="SHOW_NMON">
                <a-tab-pane key="nmon" tab="NMON" :force-render="false">
                    <Nmon></Nmon>
                </a-tab-pane>
            </template>
            <a-tab-pane key="summary" tab="基本信息">
                <Description></Description>
            </a-tab-pane>
        </a-tabs>
    </a-layout-content>
</template>

<script lang="ts" setup>
import { ref, onMounted, onBeforeMount } from 'vue';
import { useRoute, useRouter } from "vue-router"
import Overview from "@/views/linux/Overview.vue"
import Profiling from "@/views/linux/Profiling.vue"
import Topological from "@/views/linux/Topological.vue"
import Traffic from '@/views/linux/Traffic.vue';
import Description from '@/views/linux/Description.vue';
import { Linux } from './api';
import Nmon from '@/views/linux/Nmon.vue';

const activeKey = ref('overview');

const SHOW_PROFILING = import.meta.env.VITE_SHOW_LINUX_PROFILING === 'true'
const SHOW_TRFFIC = import.meta.env.VITE_SHOW_LINUX_TRAFFIC === 'true'
const SHOW_NMON = import.meta.env.VITE_SHOW_LINUX_NMON === 'true'

</script>

<style scoped></style>