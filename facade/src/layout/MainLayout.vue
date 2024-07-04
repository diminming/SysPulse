<template>
    <a-layout class="main_layout full_height">
        <a-layout-header class="headerWarp">
            <Navbar></Navbar>
        </a-layout-header>
        <a-layout-content class="tbody">
            <router-view></router-view>
        </a-layout-content>

    </a-layout>
</template>

<script setup lang="ts">
import Navbar from "@/components/Navbar.vue";
</script>

<style scoped>
.tbody {
    /* padding-top: 64px;
    overflow-y: hidden; */
}

.main_layout {
    display: flex;
    flex-direction: column;
}

.headerWarp {
    display: flex;
    justify-content: space-between;
    width: 100%;
    top: 0;
    left: 0;
    padding: 0 20px;
}
</style>