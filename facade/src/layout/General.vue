<template>
    <a-layout class="full_height">
        <a-layout-sider theme="light">
            <Sidebar />
        </a-layout-sider>
        <a-layout-content class="content">
            <Breadcrumb />
            <div class="workspace">
                <router-view></router-view>
            </div>
        </a-layout-content>
    </a-layout>
</template>
<script setup lang="ts">
import Sidebar from "@/components/Sidebar.vue";
import Breadcrumb from "@/components/Breadcrumb.vue"

</script>

<style scoped>
.content {
    display: flex;
    flex-direction: column;
    padding: 0 1rem;
}
.workspace {
    overflow-y: auto;
}
</style>