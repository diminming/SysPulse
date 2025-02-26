<template>
    <a-layout-content :style="{ background: '#fff', padding: '24px', margin: 0, minHeight: '840px' }">
        <a-space :size="200">
            <a-statistic title="业务系统" :value="biz.name" style="margin-right: 50px" />
            <a-statistic title="业务标识" :value="biz.idenetity" style="margin-right: 50px" />
            <a-statistic title="实例数量" :value="instCount" style="margin-right: 50px" />
        </a-space>

        <br>
        <br>
        <br>

        <a-tabs v-model:activeKey="activeKey">
            <a-tab-pane key="topo" tab="系统拓扑">
                <biz-topo></biz-topo>
            </a-tab-pane>

            <a-tab-pane key="alarm" tab="告警统计">
                <alarm-lst :biz="biz"></alarm-lst>
            </a-tab-pane>
        </a-tabs>
    </a-layout-content>
</template>

<script lang="ts" setup>
import { onMounted, reactive, ref } from 'vue';
import BizTopo from './BizTopo.vue';
import AlarmLst from '../notification/AlarmLst.vue';
import { BizUtil } from './BizAPI';
import { useRoute } from 'vue-router';

const activeKey = ref("topo")
const instCount = ref(0)
const route = useRoute()
const biz = reactive<BizUtil>(new BizUtil(-1))

onMounted(()=>{
    let id = parseInt(route.params.bizId as string)
    biz.id = id
    biz.load().then(resp => {
        const data = resp.data
        biz.name = data.bizName
        biz.desc = data.bizDesc
        biz.idenetity = data.bizId
    })
    biz.countInst().then(resp=>{
        instCount.value = resp.data
    })
})

</script>

<style lang="css"></style>