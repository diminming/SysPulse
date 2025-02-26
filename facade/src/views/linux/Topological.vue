<template>
    <a-row>
        <a-col :span="12">
            <div class="topo" style="height: 600px;"></div>
            <a-collapse v-model:activeKey="showGraphSetting">
                <a-collapse-panel key="1" header="拓扑图设置">
                    <a-form size="small">
                        <a-form-item label="显示所有进程">
                            <a-switch size="small" v-model:checked="showAll" @change="renderGraph()" />
                        </a-form-item>
                    </a-form>
                </a-collapse-panel>
            </a-collapse>
        </a-col>
        <a-col :span="12" style="padding: 1rem;">
            <a-tabs v-model:activeKey="activeKey">
                <a-tab-pane key="nodes" tab="节点" :force-render="false">
                    <a-table size="small" :dataSource="nodeLst" :columns="nodeColumns">
                        <template #bodyCell="{ column, record }">
                            <template v-if="column.key === 'action'">
                                <span>
                                    <a @click="showNodeDetail(record)">详情</a>
                                </span>
                            </template>
                        </template>
                    </a-table>
                </a-tab-pane>
                <a-tab-pane key="links" tab="关系" :force-render="false">
                    <a-table size="small" :dataSource="edgeLst" :columns="edgeColumns">
                        <template #bodyCell="{ column, record }">
                            <template v-if="column.key === 'type'">
                                <template v-if="record['type'] === 'depl'">
                                    <b>部署</b>
                                </template>
                                <template v-else-if="record['type'] === 'conn_tcp'">
                                    <b>TCP连接</b>
                                </template>
                            </template>
                            <template v-if="column.key === 'action'">
                                <span>
                                    <a @click="showEdgeDetail(record)">详情</a>
                                </span>
                            </template>
                        </template>
                    </a-table>
                </a-tab-pane>
            </a-tabs>
        </a-col>
    </a-row>
    <a-modal v-model:open="nodeDetailIsOpen" title="节点详情" @ok="closeNodeDetail">
        <a-descriptions bordered size="small" :column=1>
            <a-descriptions-item label="PID">{{ activeNode.pid }}</a-descriptions-item>
            <a-descriptions-item label="名称">{{ activeNode.name }}</a-descriptions-item>
            <a-descriptions-item label="执行程序">{{ activeNode.exec }}</a-descriptions-item>
            <a-descriptions-item label="启动时间">{{ timestamp2DateString(activeNode.timestamp) }}</a-descriptions-item>
        </a-descriptions>
    </a-modal>
</template>
<script lang="ts" setup>
import { timestamp2DateString } from '@/utils/common';
import { useRoute } from 'vue-router';
import { onMounted, ref } from 'vue';
import { Linux } from '@/views/linux/api'

const linux = ref<Linux>(new Linux(-1)),
    activeKey = ref('nodes'),
    route = useRoute(),
    l = new Linux(parseInt(route.params.linuxId as string)),
    showGraphSetting = ref<Boolean>(false),
    showAll = ref<Boolean>(false),
    nodeDetailIsOpen = ref<Boolean>(false),
    edgeDetailIsOpen = ref<Boolean>(false),
    activeNode = ref({}),
    activeEdge = ref({}),
    nodeLst = ref([]),
    edgeLst = ref([]),
    nodeColumns = [{
        title: 'IDX',
        dataIndex: 'idx',
        key: 'idx',
    }, {
        title: 'Name',
        dataIndex: 'name',
        key: 'name',
    }, {
        title: '',
        key: 'action',
    },],
    edgeColumns = [{
        title: '源',
        dataIndex: 'from',
        key: 'from',
    }, {
        title: "类型",
        dataIndex: "type",
        key: 'type'
    }, {
        title: '目标',
        dataIndex: 'to',
        key: 'to',
    }, {
        title: '',
        key: 'action',
    }]

const closeNodeDetail = () => {
    nodeDetailIsOpen.value = false
    activeNode.value = {}
}, closeEdgeDetail = () => {
    edgeDetailIsOpen.value = false
}

const showNodeDetail = (record: any) => {
    const detail = record["_detail"]
    activeNode.value = detail
    nodeDetailIsOpen.value = true
}, showEdgeDetail = (record: any) => {
    const detail = record["_detail"]
    activeEdge.value = detail
    edgeDetailIsOpen.value = true
}

const renderGraph = () => {
    
    l.RenderTopological(showAll.value).then(({ verteies, edges }) => {
        nodeLst.value = verteies.map((item, idx) => {
            return {
                "idx": idx + 1,
                "name": item["name"] || item['value']["name"],
                "_detail": item['detail']
            }
        })
        edgeLst.value = edges.map((link: any, idx: Number) => {
            return {
                "from": link.source,
                "to": link.target,
                "type": link['_detail']['type'],
                "_detail": link['_detail']
            }
        })
    })
    linux.value = l
}

onMounted(() => {
    renderGraph()
})
</script>
<style lang="css" scoped>
.topo {
    height: 600px;
}
</style>