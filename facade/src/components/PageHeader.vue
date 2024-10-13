<template>
    <div style="background-color: white; margin-top:0.5rem">
        <a-page-header :title="title" :breadcrumb="{ routes, click }" @back="() => { }"
            style="border: 1px solid rgb(235, 237, 240)" />
    </div>

</template>

<script lang="ts" setup>
import { ref, onMounted, watch } from "vue";
import { useRouter } from "vue-router";

const router = useRouter()
const title = ref("")
const routes = ref([])
const getBreadcrumbs = () => {
    return router.currentRoute.value.matched.slice(1).map((route) => {
        return {
            //   isActive: route.path === router.currentRoute.value.fullPath,
            path: `${router.options.history.base}${route.path}`,
            breadcrumbName: route.meta.text,
        }
    })
}

router.afterEach(() => {
    const lst = getBreadcrumbs()
    routes.value = lst
    title.value = lst.slice(-1)[0].breadcrumbName
})

routes.value = getBreadcrumbs()

const click = (e: MouseEvent) => {
    console.log(e)
}

</script>

<style lang="css"></style>