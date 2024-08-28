<template>
    <a-row>
        <a-col :span="12">
            <div class="topo" style="height: 600px;"></div>
        </a-col>
        <a-col :span="12">
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
                                    <a @click="showNodeDetail(record)">详情</a>
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
            <a-descriptions-item label="PID">{{ activeNode.pid}}</a-descriptions-item>
            <a-descriptions-item label="名称">{{ activeNode.name}}</a-descriptions-item>
            <a-descriptions-item label="执行程序">{{ activeNode.exec}}</a-descriptions-item>
            <a-descriptions-item label="启动时间">{{ activeNode.timestamp}}</a-descriptions-item>
        </a-descriptions>
    </a-modal>
    <a-modal v-model:open="edgeDetailIsOpen" title="Basic Modal" @ok="handleOk">
        <a-descriptions bordered title="Custom Size" size="small">
            <template #extra>
                <a-button type="primary">Edit</a-button>
            </template>
            <a-descriptions-item label="Product">Cloud Database</a-descriptions-item>
            <a-descriptions-item label="Billing">Prepaid</a-descriptions-item>
            <a-descriptions-item label="Time">18:00:00</a-descriptions-item>
            <a-descriptions-item label="Amount">$80.00</a-descriptions-item>
            <a-descriptions-item label="Discount">$20.00</a-descriptions-item>
            <a-descriptions-item label="Official">$60.00</a-descriptions-item>
            <a-descriptions-item label="Config Info">
                Data disk type: MongoDB
                <br />
                Database version: 3.4
                <br />
                Package: dds.mongo.mid
                <br />
                Storage space: 10 GB
                <br />
                Replication factor: 3
                <br />
                Region: East China 1
                <br />
            </a-descriptions-item>
        </a-descriptions>
    </a-modal>

</template>
<script lang="ts" setup>
import { useRoute } from 'vue-router';
import { onMounted, ref } from 'vue';
import { Linux } from '@/views/linux/api'

const linux = ref<Linux>(new Linux(-1))
const activeKey = ref('nodes'),
    nodeDetailIsOpen = ref<Boolean>(false),
    edgeDetailIsOpen = ref<Boolean>(false),
    activeNode = ref({}),
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
 }

const showNodeDetail = (record: any) => {
    const detail = record["_detail"]
    activeNode.value = detail
    nodeDetailIsOpen.value = true
}

onMounted(() => {
    let route = useRoute()
    let l = new Linux(parseInt(route.params.linuxId as string))
    l.RenderTopological((chart: any, verteies: any, links: any) => {
        const nodeLst1 = [...verteies.values()]
        l.RenderTopoGraph(chart, nodeLst1, links)
        nodeLst.value = nodeLst1.map((item, idx) => {
            return {
                "idx": idx + 1,
                "name": item["name"] || item['value']["name"],
                "_detail": item['_detail']
            }
        })
        edgeLst.value = links.map((link: any, idx: Number) => {
            return {
                "from": link.source,
                "to": link.target,
                "type": link['_detail']['type'],
                "_detail": link['_detail']
            }
        })
    })

    linux.value = l
})
</script>
<style lang="css" scoped>
.topo {
    height: 600px;
}
</style>