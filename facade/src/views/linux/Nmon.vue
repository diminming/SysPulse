<template>
    <a-row :span="24">
        <a-table :data-source="nmonLst" :columns="columns" size="small" :row-selection="rowSelection"
            :pagination="pagination" @change="onChange" class="data_lst">
            <template #bodyCell="{ column, record }">
                <template v-if="column.key === 'operation'">
                    <a-button size="small" primary type="link" @click="openNMONFile(record)">打开</a-button>
                </template>
                <template v-else-if="column.key === 'from'">
                    {{ timestamp2DateString(record.from) }}
                </template>
                <template v-else-if="column.key === 'to'">
                    {{ timestamp2DateString(record.to) }}
                </template>
            </template>
        </a-table>
    </a-row>
    <a-drawer height="90%" :title="title" placement="bottom" :open="open" @close="onClose">
        <a-layout>
            <a-layout-sider :width="380" style="background-color: #f0f2f5">
                <a-card style="width: 100%; height: 100%" size="small">
                    <a-input-search v-model:value="kwCategory" style="margin: 0 0.5 0.5 0.5" placeholder="Search" />
                    <a-tree :tree-data="categoryLst" @select="selectNode" style="height: 700px ;overflow-y: scroll;">
                        <template #title="{ title }">
                            <span v-if="title.indexOf(kwCategory) > -1">
                                {{ title.substring(0, title.indexOf(kwCategory)) }}
                                <span style="color: #f50">{{ kwCategory }}</span>
                                {{ title.substring(title.indexOf(kwCategory) + kwCategory.length) }}
                            </span>
                            <span v-else>{{ title }}</span>
                        </template>
                    </a-tree>
                </a-card>
            </a-layout-sider>
            <a-layout-content>
                <a-card style="width: 100%; height: 700px; margin: 0 0.5rem" title="NMON数据" size="small">
                    <div id="graph" style="width: 100%; height: 700px"></div>
                </a-card>
            </a-layout-content>
        </a-layout>
    </a-drawer>
</template>

<script lang="ts" setup>
import {timestamp2DateString} from "@/utils/common";
import { type TableColumnType, type TableProps } from 'ant-design-vue';
import { NMON } from "./nmon_api"
import { onMounted, reactive, ref } from 'vue';
import * as echarts from 'echarts';
import { useRoute } from 'vue-router';


const nmonLst = ref<NMON[]>([]);
const rowSelection: TableProps["rowSelection"] = {};

const pagination = reactive({
    page: 0,
    pageSize: 10,
    total: 0,
})

const columns: TableColumnType<NMON>[] = [
    {
        title: "主机名",
        dataIndex: "hostname",
        key: "hostname",
    },
    {
        title: "开始时间",
        dataIndex: "from",
        key: "from",
    },
    {
        title: "结束时间",
        dataIndex: "to",
        key: "to",
    },
    {
        title: "来源",
        dataIndex: "source",
        key: "source",
    },
    {
        title: "文件路径",
        dataIndex: "path",
        key: "path",
    },
    {
        title: "操作",
        dataIndex: "operation",
        key: "operation",
        width: 150
    }
];

const getPage = () => {
    NMON.getByPage({ "linuxId": linuxId.value, page: pagination.page, size: pagination.pageSize }).then((resp) => {
        const data = resp.data
        nmonLst.value = data.lst.map((item: any) => new NMON(
            item.id,
            item.hostname,
            item.from,
            item.to,
            item.source,
            item.path
        ));
        pagination.total = data.total
    });
}

const route = useRoute(), linuxId = ref(-1)

onMounted(() => {
    linuxId.value = parseInt(route.params.linuxId as string)
    getPage()
})

const onChange = (pg: { pageSize: number; current: number }) => {
    pagination.page = pg.current - 1;
    pagination.pageSize = pg.pageSize;
    getPage()
};

const categoryLst = ref([])

let currNMON = new NMON(-1)

const openNMONFile = (nmon: NMON) => {
    open.value = true
    title.value = "NMON文件详情: " + (nmon && nmon.path ? nmon.path : "")
    currNMON = nmon
    nmon.getCategories().then(resp => {
        const data = resp.data
        const categories = data[0]
        categoryLst.value = categories.map((item0: any) => {
            const id = item0['id']
            const name = item0['name']
            const fields = item0['fields'].map((item1: any) => {
                return {
                    key: id + "-" + item1,
                    title: item1
                }
            })
            return {
                key: id,
                title: name,
                children: fields
            }
        })
    })
}

const open = ref<boolean>(false);
const title = ref<String>("NMON文件详情: ")
const kwCategory = ref<String>("")

const onClose = () => {
    open.value = false;
};

const selectNode = (selectedKeys: string) => {
    const array = selectedKeys[0].split("-")
    const category = array[0]
    const field = array[1]
    currNMON.getNMONData(category, field).then(resp => {
        const data = resp['data']
        renderData(data)
    })
}

const draw = (series, title, legend) => {
    const dom = document.getElementById('graph')
    const chart = echarts.getInstanceByDom(dom as HTMLElement) || echarts.init(dom as HTMLElement)
    chart.setOption({
        title: {
            left: '30',
            text: title ? title : ""
        },
        legend: {
            type: "scroll",
            data: legend,
            orient: 'vertical',
            top: 30,
            right: 10,
            y: 'center',
            formatter: function (name) {
                return name.length > 10 ? name.substring(0, 10) + '...' : name
            },
            //开启tooltip
            tooltip: {
                show: true
            }
        },
        tooltip: {
            trigger: 'axis',
            axisPointer: {
                type: 'cross'
            },
            backgroundColor: 'rgba(255, 255, 255, 0.8)',
            showContent: true
        },
        toolbox: {
            feature: {
                dataZoom: {
                    yAxisIndex: 'none'
                },
                restore: {},
                saveAsImage: {}
            }
        },
        xAxis: {
            type: 'time',
            // boundaryGap: false
        },
        yAxis: {
            type: 'value',
            splitLine: {
                show: true
            }
        },
        dataZoom: [
            {
                type: 'inside',
                start: 0,
                end: 100
            },
            {
                start: 0,
                end: 100
            }
        ],
        series: series === undefined ? [] : series
    })

}

const renderData = (data: any, option: any) => {
    const chart = echarts.init(document.getElementById('graph'))
    const values = data;
    let series = [], legend = []
    chart.clear()
    Object.entries(values).forEach(([key, value]) => {
        legend.push(key)
        series.push({
            name: key,
            type: 'line',
            areaStyle: {},
            showSymbol: false,
            data: value
        })
    })
    draw(series, "", legend)
}

</script>

<style scoped lang="css">
.data_lst {
    width: 100%;
}
</style>