<template>
    <a-layout class="full_height">
        <a-layout-sider theme="light" collapsible v-model:collapsed="state.collapsed">
            <Sidebar />
        </a-layout-sider>
        <a-layout-content class="content">
            <Breadcrumb/>
            <!-- <PageHeader/> -->
            <div class="workspace">
                <router-view></router-view>
            </div>
        </a-layout-content>
    </a-layout>
</template>
<script setup lang="ts">
import Sidebar from "@/components/Sidebar.vue";
import PageHeader from "@/components/PageHeader.vue";
import Breadcrumb from "@/components/Breadcrumb.vue"
import { sidebarCollapsed } from '@/stores/sidebarStatus';

const state = sidebarCollapsed()
</script>

<style scoped>
.content {
    display: flex;
    flex-direction: column;
    padding: 0 1rem;
}
.workspace {
    margin-top: 0.5rem;
    overflow-y: auto;
}
</style>