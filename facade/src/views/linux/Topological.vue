<template>
    <a-row>
        <a-col :span="12">
            <div class="topo"></div>
        </a-col>
        <a-col :span="12">right</a-col>
    </a-row>
</template>
<script lang="ts" setup>
import { useRoute } from 'vue-router';
import { onMounted, reactive, ref } from 'vue';
import { Linux } from '@/views/linux/linux'

const linux = ref<Linux>(new Linux(-1))

onMounted(() => {
    let route = useRoute()
    let l = new Linux(parseInt(route.query.linuxId as string))
    linux.value = l
    l.RenderTopological()
})
</script>
<style lang="css" scoped>
.topo {
    height: 600px;
}
</style>