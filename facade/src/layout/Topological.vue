<template>
    <div class="graph full_height">
        <canvas id="topo" width="1280" height="960"></canvas>
    </div>
</template>

<script setup lang="ts">
import * as echarts from 'echarts';
import { onMounted } from 'vue';
import request from '@/utils/request';
import data from "@/layout/demo"

onMounted(() => {
    let container: any = document.querySelector(".graph")
    let chart_dom: any = document.getElementById('topo')
    
    let chart = echarts.init(chart_dom, undefined, {
        width: container.offsetWidth,
        height: container.offsetHeight
    });
    let option;

    chart.showLoading();

    let json = data
    chart.hideLoading();
    chart.setOption(
        (option = {
            title: {
                text: 'NPM Dependencies'
            },
            animationDurationUpdate: 1500,
            // animationEasingUpdate: 'quinticInOut',
            series: [
                {
                    type: 'graph',
                    layout: 'none',
                    // progressiveThreshold: 700,
                    data: json.nodes.map(function (node) {
                        return {
                            x: node.x,
                            y: node.y,
                            id: node.id,
                            name: node.label,
                            symbolSize: node.size,
                            itemStyle: {
                                color: node.color
                            }
                        };
                    }),
                    edges: json.edges.map(function (edge) {
                        return {
                            source: edge.sourceID,
                            target: edge.targetID
                        };
                    }),
                    emphasis: {
                        focus: 'adjacency',
                        label: {
                            position: 'right',
                            show: true
                        }
                    },
                    roam: true,
                    lineStyle: {
                        width: 0.5,
                        curveness: 0.3,
                        opacity: 0.7
                    }
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