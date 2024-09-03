<template>
    <div style="text-align: right; margin-top: 1rem;">
        <a-range-picker v-model:value="value2" show-time />
    </div>
    <div style="position: absolute;bottom: 5rem;right: 5rem;">
        <a-button shape="circle" :icon="h(CopyOutlined)" />
    </div>
    <div>
        <a-slider v-model:value="value2" range :disabled="false" />
    </div>
    <div class="graph full_height">
        <canvas id="topo" width="1280" height="960"></canvas>
    </div>
    
</template>

<script setup lang="ts">
import * as echarts from 'echarts';
import { onMounted } from 'vue';
import request from '@/utils/request';
import data from "@/layout/demo"
import { h } from 'vue';
import { SearchOutlined, CopyOutlined } from '@ant-design/icons-vue';

onMounted(() => {
    let container: any = document.querySelector(".graph")
    let chart_dom: any = document.getElementById('topo')

    let chart = echarts.init(chart_dom, undefined, {
        width: container.offsetWidth,
        height: container.offsetHeight
    });
    let option;

    chart.showLoading();

    let webkitDep = data
    chart.hideLoading();
    chart.setOption(
        (option = {
            legend: {
                data: ['Linux', 'Database', 'Middleware', 'Cache', 'Other']
            },
            series: [
                {
                    type: 'graph',
                    layout: 'force',
                    animation: false,
                    label: {
                        position: 'right',
                        formatter: '{b}'
                    },
                    draggable: true,
                    data: webkitDep.nodes.map(function (node, idx) {
                        node.id = idx;
                        return node;
                    }),
                    categories: webkitDep.categories,
                    force: {
                        edgeLength: 5,
                        repulsion: 20,
                        gravity: 0.2
                    },
                    edges: webkitDep.links
                }
            ]
        }),
    true
    );


option && chart.setOption(option);
})

</script>
<style scoped>
.graph {
    width: 100%;
    text-align: center;
}
</style>