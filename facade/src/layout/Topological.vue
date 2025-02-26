<template>
    <div class="graph full_height">
        <canvas id="topo" width="1280" height="960"></canvas>
        <div style="width: 60px; height: 60px; position:absolute; bottom: 5rem; right: 5rem;">
            <a-button shape="circle" size="large" :icon="h(ToolOutlined)" @click="settingVisible = true" />
        </div>
    </div>
    <a-drawer title="拓扑图查询选项" placement="right" v-model:visible="settingVisible" width="30%">
        <a-form>
            <a-form-item label="节点数量">
                <a-input-number v-model:value="setting.limit" :min="1" :max="1000" />
            </a-form-item>
            <a-form-item>
                <a-button type="primary" @click="renderGraph">确认</a-button>
            </a-form-item>
        </a-form>
    </a-drawer>
</template>

<script setup lang="ts">
import * as echarts from 'echarts';
import { onMounted, reactive, ref, h } from 'vue';
import { ToolOutlined } from '@ant-design/icons-vue';
import request from '@/utils/request';
import type { EChartsType } from 'echarts';

let chart: EChartsType;

const setting = reactive({
    limit: 100,
    start: {},
    depth: 4
})

const settingVisible = ref(false)

let option = {
    legend: [
        {
            data: ["Linux", "Process", "Business"]
        }
    ],
    series: [
        {
            type: 'graph',
            layout: 'force',
            animation: true,
            emphasis: {
                focus: 'adjacency',
                label: {
                    position: 'right',
                    show: true
                }
            },
            edgeSymbol: ['circle', 'arrow'],
            edgeSymbolSize: [4, 8],
            roam: true,
            label: {
                show: true,
                position: 'right',
                formatter: '{b}'
            },
            labelLayout: {
                hideOverlap: true
            },
            draggable: true,
            data: <any>[],
            categories: [{
                "name": "Linux",
                "base": "Linux",
                "keyword": {}
            }, {
                "name": "Process",
                "base": "Process",
                "keyword": {}
            }, {
                "name": "Business",
                "base": "Business",
                "keyword": {}
            }],
            force: {
                initLayout: "circular",
                edgeLength: 200,
                repulsion: 100,
                gravity: 0.1,
                friction: 0.1
            },
            lineStyle: {
                width: 0.8,
                curveness: 0.3,
                opacity: 0.7
            },
            edges: <any>[]
        }
    ]
}

const start = reactive<{
    id: string,
    key: string,
    name: string,
    identity: string
}>({ id: "", key: "", name: "", identity: "" })

const init = (): EChartsType => {
    let container: any = document.querySelector(".graph")
    let chart_dom: any = document.getElementById('topo')
    chart_dom.width = container.offsetWidth
    chart_dom.height = container.offsetHeight
    const chart = echarts.init(chart_dom, undefined, {
        width: container.offsetWidth,
        height: container.offsetHeight
    });
    chart.on('click', (params) => {
        const data: any = params.data
        if (data == null) return
        start.id = data._id
        start.key = data.key
        start.name = data.name
        start.identity = data.identity
        setting.start = data
        getHostGraph().then((data: any) => {

            const nodeMap = new Map()
            const edges = new Array()

            data.forEach((item: any) => {
                let vertex = item.vertex
                let node;
                if (vertex._id.startsWith("host/")) {
                    node = {
                        id: vertex._id,
                        _id: vertex._id,
                        key: vertex.key,
                        name: vertex.name,
                        identity: vertex.identity,
                        category: 0,
                        "symbol": "roundRect",
                        "symbolSize": [40, 30]
                    }
                } else if (vertex._id.startsWith("process/")) {
                    node = {
                        id: vertex._id,
                        _id: vertex._id,
                        key: vertex._key,
                        name: vertex.info.name,
                        identity: vertex.pid,
                        category: 1,
                        "symbolSize": 20
                    }
                } else if (vertex._id.startsWith("business/")) {
                    node = {
                        id: vertex._id,
                        _id: vertex._id,
                        key: vertex._key,
                        name: vertex.bizName,
                        identity: vertex.bizId,
                        category: 2,
                        "symbol": "triangle",
                        "symbolSize": 40
                    }
                } else {
                    console.log(vertex)
                }
                nodeMap.set(node?.id, node)
            })

            const nodes = Array.from(nodeMap.values())

            data.forEach((item: any) => {
                if (item.edge == null) {
                    return
                }

                edges.push({
                    source: nodes.findIndex((node: any) => {
                        return node._id == item.edge._from
                    }),
                    target: nodes.findIndex((node: any) => {
                        return node._id == item.edge._to
                    })
                })
            })

            option.series[0].data = nodes
            option.series[0].edges = edges
            chart.clear()
            chart.setOption(option);
        }).catch((error) => {
            console.error(error)
        })
    })
    return chart
}

const getHostGraph = () => {
    return new Promise((resolve, reject) => {
        request({
            url: '/linuxGraph',
            method: 'get',
            params: {
                "limit": setting.limit,
                "start": start.id,
                "depth": setting.depth
            }
        }).then(resp => {
            const data = resp.data;
            resolve(data);
        }).catch(error => {
            reject(error);
        });
    });
}

const renderGraph = () => {
    chart.showLoading();
    getHostGraph().then((data: any) => {
        chart.hideLoading();
        const nodeMap = new Map<String, any>()
        const edges = new Array()

        data.forEach((item: any) => {
            let from = item.from
            let to = item.to
            nodeMap.has(from._id) || nodeMap.set(from._id, {
                "name": from.name,
                "key": from.key,
                "_id": from._id,
                "identity": from.identity,
                "category": 0,
                "symbol": "roundRect",
                "symbolSize": [40, 30]
            })
            nodeMap.has(to._id) || nodeMap.set(to._id, {
                "name": to.name,
                "key": to.key,
                "_id": to._id,
                "identity": to.identity,
                "category": 0,
                "symbol": "roundRect",
                "symbolSize": [40, 30]
            })
        });

        const nodes = Array.from(nodeMap.values())
        data.forEach((item: any) => {
            edges.push({
                source: nodes.findIndex((node: any) => node._id == item.from._id),
                target: nodes.findIndex((node: any) => node._id == item.to._id)
            })
        });
        option.series[0].data = nodes
        option.series[0].edges = edges

        chart.setOption(option);
    }).catch((error) => {
        console.error(error);
        chart.hideLoading();
    });
}


onMounted(() => {

    chart = init()
    renderGraph()

})

</script>
<style scoped>
.graph {
    width: 100%;
    text-align: center;
}
</style>