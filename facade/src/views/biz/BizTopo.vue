<template>
    <div id="biz_topo" class="biz_topo"></div>
    <a-button type="primary" v-if="!isOpen" class="btn_setting" shape="circle" :icon="h(FormOutlined)"
        @click="isOpen = true" />
    <a-drawer v-model:open="isOpen" title="图遍历设置" placement="right" size="large">
        <template #extra>
            <a-button style="margin-right: 8px" @click="isOpen = false">取消</a-button>
            <a-button type="primary" @click="renderGraph">查询</a-button>
        </template>
        <a-form :model="graphSetting">
            <a-form-item label="最小深度">
                <a-input-number v-model:value="graphSetting.min" :min="1" :max="5"></a-input-number>
            </a-form-item>
            <a-form-item label="最大深度">
                <a-input-number v-model:value="graphSetting.max" :min="1" :max="5"></a-input-number>
            </a-form-item>

        </a-form>
    </a-drawer>
</template>

<script lang="ts" setup>
import * as echarts from "echarts/core";
import { CustomChart, GraphChart } from 'echarts/charts';
import { CanvasRenderer } from "echarts/renderers"; // 引入渲染器
import { TooltipComponent, LegendComponent } from "echarts/components"; // 引入所需组件

import { h, reactive, ref } from 'vue';
import { FormOutlined } from '@ant-design/icons-vue';

import { pathDecomposition, renderTopoGraph } from "@/utils/graphutils";
import { BizUtil } from "@/views/biz/BizAPI"
import { onMounted } from "vue";
import { useRoute, } from "vue-router";

const route = useRoute(),
    isOpen = ref(false),
    graphSetting = reactive<{
        min: number,
        max: number,
        vset: string[],
        eset: string[]
    }>({
        min: 0,
        max: 4,
        vset: ['linux', 'process'],
        eset: ['belong', 'deployment', 'connection']
    })
echarts.use([CustomChart, CanvasRenderer, GraphChart, TooltipComponent, LegendComponent]);

const renderGraph = () => {
    let id = parseInt(route.params.bizId as string)
    const biz = new BizUtil(id)

    const dom = document.getElementById("biz_topo"),
        chart = echarts.getInstanceByDom(dom as HTMLElement) || echarts.init(dom as HTMLElement)
    chart.clear()
    chart.showLoading();

    biz.loadTopo(graphSetting).then((resp) => {

        const pathLst = resp["data"]

        let { verteies, edges } = pathDecomposition(pathLst)
        const nodeLst1 = [...verteies.values()]

        renderTopoGraph(chart, nodeLst1, edges)
        chart.hideLoading()
        isOpen.value = false

    })
}

onMounted(() => {
    renderGraph()
})

</script>

<style lang="css">
.biz_topo {
    width: 100%;
    height: 768px;
}

.btn_setting {
    position: absolute;
    bottom: 100px;
    right: 100px;
}

.graph_setting {
    width: 100rem;
}
</style>