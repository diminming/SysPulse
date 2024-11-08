<template>
    <a-row :span="24">
        <div class="toolbar">
            <a-button type="primary" @click="openCreateJobDialog">创建任务</a-button>
        </div>
    </a-row>
    <a-row :span="24">
        <a-table :data-source="jobLst" :columns="columns" size="small" :row-selection="rowSelection"
            :pagination="pagination" @change="onChange" class="job_lst">
            <template #bodyCell="{ text, column, record }">

                <template v-if="column.key === 'job_name'">
                    <a-button type="link" @click="getResult(record)">{{ record['job_name'] }}</a-button>
                </template>

                <template v-else-if="column.key === 'status'">
                    {{ (record as Job).getJobStatusTxt() }}
                </template>

                <template v-else-if="column.key === 'direction'">
                    {{
                        record.direction.map(item => {
                            if (item === "in")
                                return "入站"
                            else if (item === "out")
                                return "出站"
                        }).join(", ")
                    }}
                </template>

                <template v-else-if="column.key === 'operation'">
                    <span>
                        <a-popconfirm title="是否确认删除该记录?" ok-text="确认" cancel-text="取消" @confirm="onDelete(record.id)">
                            <a-button danger type="link">删除</a-button>
                        </a-popconfirm>
                    </span>
                </template>

            </template>
        </a-table>
    </a-row>

    <a-row style="margin-top: 2rem;">
        <a-col :span="24">
            <a-row :gutter="16">
                <a-col :span="6">
                    <a-statistic title="TCP包" :value="countTcp" style="" />
                </a-col>
                <a-col :span="6">
                    <a-statistic title="TCP连接" :value="countConn" style="" />
                </a-col>
                <a-col :span="6">
                    <a-statistic title="HTTP请求" :value="countReq" style="" />
                </a-col>
                <a-col :span="6">
                    <a-statistic title="URL" :value="countUrl" style="" />
                </a-col>
            </a-row>

            <a-row :gutter="16" style="margin-top: 2rem;">
                <a-col :span="8">
                    <a-card class="bar_graph" title="响应状态分析图" size="small">
                        <div class="graph status">

                        </div>
                    </a-card>
                </a-col>
                <a-col :span="8">
                    <a-card class="bar_graph" title="响应时间分析图" size="small">
                        <div class="graph time">

                        </div>
                    </a-card>
                </a-col>
                <a-col :span="8">
                    <a-card class="bar_graph" title="吞吐量分析图" size="small">
                        <div class="graph throughputStat">

                        </div>
                    </a-card>
                </a-col>
            </a-row>

            <a-row :gutter="16" style="margin-top: 2rem;">
                <a-col :span="24">
                    <a-table :columns="columns2" :data-source="connDataLst" size="small"></a-table>
                </a-col>
            </a-row>

        </a-col>
    </a-row>

    <a-modal v-model:open="showCreateJobDialog" title="创建流量分析任务" @ok="handleOk">
        <a-form ref="createJobForm" :model="jobObj" :label-col="labelCol" :wrapper-col="wrapperCol">

            <a-form-item ref="name" label="任务标识" name="job_name">
                <a-input v-model:value="jobObj.job_name" />
            </a-form-item>

            <a-form-item ref="ifName" label="网卡名称" name="ifName">
                <a-select v-model:value="jobObj.ifName" :options="ifLst" @change="changeIPLst"></a-select>
            </a-form-item>

            <a-form-item ref="ipAddr" label="IP地址" name="ifName">
                <a-select v-model:value="jobObj.ipAddr" :options="ipLst"></a-select>
            </a-form-item>

            <a-form-item label="端口" name="port">
                <a-input-number v-model:value="jobObj.port" :min="1" :max="65535" :step="1" />
            </a-form-item>

            <a-form-item label="流量方向" name="direction">
                <a-checkbox-group v-model:value="jobObj.direction">
                    <a-checkbox value="in" name="direction">入站</a-checkbox>
                    <a-checkbox value="out" name="direction">出站</a-checkbox>
                </a-checkbox-group>
            </a-form-item>

            <a-form-item label="抓包数量" name="count">
                <a-input-number v-model:value="jobObj.count" :min="100" :max="10000" :step="10" />
            </a-form-item>

        </a-form>
    </a-modal>
</template>

<script lang="ts" setup>
import { onMounted, reactive, ref } from 'vue';
import { Job, TrafficAnalyzationJob, GetTrafficJobLst } from '@/views/linux/api';
import { useRoute } from 'vue-router';
import { Linux } from './api';
import type { TableColumnType, TableProps } from "ant-design-vue";
import echarts from '@/utils/echarts';
import { notification } from 'ant-design-vue';

const countTcp = ref<Number>(0),
    countConn = ref<Number>(0),
    countUrl = ref<Number>(0),
    countReq = ref<Number>(0),
    connDataLst = ref([])

const labelCol = { span: 5 };
const wrapperCol = { span: 13 };

const rowSelection: TableProps["rowSelection"] = {};

const pagination = reactive({
    page: 0,
    pageSize: 20,
    total: 0,
})

const columns: TableColumnType<Job>[] = [
    {
        title: "任务名称",
        dataIndex: "job_name",
        key: "job_name",
    },
    {
        title: "任务状态",
        dataIndex: "status",
        key: "status",
    },
    {
        title: "网卡名称",
        dataIndex: "ifName",
        key: "ifName",
    },
    {
        title: "IP地址",
        dataIndex: "ipAddr",
        key: "ipAddr",
    },
    {
        title: "端口",
        dataIndex: "port",
        key: "port",
    }, {
        title: "流量方向",
        dataIndex: "direction",
        key: "direction",
    }, {
        title: "抓包数量",
        dataIndex: "count",
        key: "count",
    },
    {
        title: "操作",
        dataIndex: "operation",
        key: "operation",
    },
];

const jobObj: TrafficAnalyzationJob = reactive(new TrafficAnalyzationJob(-1));

const showCreateJobDialog = ref<Boolean>(false)
const ifLst = ref<[]>(), ipLst = ref<[]>(), jobLst = ref<[]>()

const openCreateJobDialog = () => {
    showCreateJobDialog.value = true
}

const handleOk = () => {
    jobObj.createJob().then((resp => {
        notification.success({
            message: '创建成功',
            description:
                '流量分析任务已创建。',
            duration: 2,
            onClick: () => {
                // console.log('Notification Clicked!');
            },
        });
        showCreateJobDialog.value = false
        getJobList()
    }))
};

const mapping4IfIp = new Map();

const changeIPLst = () => {
    jobObj.ipAddr = undefined
    ipLst.value = mapping4IfIp.get(jobObj.ifName).map(item => {
        return {
            "value": item.split("/")[0]
        }
    })
}

const getJobList = () => {
    GetTrafficJobLst(pagination.page, pagination.pageSize, jobObj.linux_id).then((resp) => {
        jobLst.value = resp.data.lst.map((item: any) => {
            const job = new TrafficAnalyzationJob(item.id)
            job.port = item.port
            job.count = item.count
            job.job_name = item.job_name
            job.direction = item.direction
            job.ifName = item.ifName
            job.ipAddr = item.ipAddr
            job.status = item.status
            return job
        })
    })
}

onMounted(() => {
    let route = useRoute()
    const linux = new Linux(parseInt(route.params.linuxId as string))
    linux.GetInterfaceLst().then((resp: any) => {
        ifLst.value = resp['data'][0]['if_lst'].map((item: any) => {
            return {
                "value": item['name']
            }
        })
        resp['data'][0]['if_lst'].forEach(element => {
            if (mapping4IfIp.has(element['name']))
                mapping4IfIp.set(element['name'], mapping4IfIp.get(element['name']).concat(element['addrs'].map(item => item['addr'])))
            else
                mapping4IfIp.set(element['name'], element['addrs'].map(item => item['addr']))
        });
        console.log(mapping4IfIp)
    })
    jobObj.linux_id = parseInt(route.params.linuxId as string)

    getJobList()

})

const columns2: TableColumnType[] = [
    {
        title: "URL",
        dataIndex: "url",
        key: "url",
    },
    {
        title: "HTTP Method",
        dataIndex: "method",
        key: "method",
    },
    {
        title: "响应码",
        dataIndex: "status",
        key: "status",
    },
    {
        title: "耗时(毫秒)",
        dataIndex: "time",
        key: "time",
    },
    {
        title: "吞吐(KB)",
        dataIndex: "throughput",
        key: "throughput",
    }
];

function RenderStatusStatChart(statusMap: Map<number, number>) {
    const xData = Array.from(statusMap.keys()).sort((a, b) => {
        return a - b
    });
    let dom = document.querySelector(".graph.status"),
        chart = echarts.init(dom),
        option = {
            xAxis: {
                type: 'value'
            },
            yAxis: {
                type: 'category',
                data: xData
            },
            legend: {
                show: true
            },
            series: [
                {
                    data: xData.map(status => {
                        return statusMap.get(status)
                    }),
                    type: 'bar'
                }
            ]
        };
    chart.setOption(option)
}

function RenderRespTimeStatChart(dataLst) {
    let dom = document.querySelector(".graph.time"),
        chart = echarts.init(dom),
        option = {
            tooltip: {
                trigger: 'axis',
                axisPointer: {
                    type: 'shadow'
                }
            },
            xAxis: {
                type: 'value'
            },
            yAxis: {
                inverse: true,
                type: 'category',
                axisLabel: {
                    show: false
                },
                data: dataLst.map(item => {
                    return item["url"]
                }),
            },
            series: [
                {
                    data: dataLst.map(item => {
                        return item["time"]
                    }),
                    label: {
                        show: true,
                        position: "insideLeft",
                        align: "left",
                        fontSize: 10,
                        formatter: (params) => {
                            let name = params.name
                            if (name.length > 25) {
                                return params.value + "ms  " + name.slice(0, 20) + "..." + name.slice(-5)
                            }
                            return params.value + "ms  " + name
                        },
                        rich: {
                            name: {}
                        }
                    },
                    type: 'bar'
                }
            ]
        };
    chart.setOption(option)
}

function RenderThroughputStatChart(dataLst) {
    console.log(dataLst)
    let dom = document.querySelector(".graph.throughputStat"),
        chart = echarts.init(dom),
        option = {
            tooltip: {
                trigger: 'axis',
                axisPointer: {
                    type: 'shadow'
                }
            },
            xAxis: {
                type: 'value'
            },
            yAxis: {
                inverse: true,
                type: 'category',
                axisLabel: {
                    show: false
                },
                data: dataLst.map(item => {
                    return item["url"]
                }),
            },
            series: [
                {
                    data: dataLst.map(item => {
                        return item["throughput"]
                    }),
                    label: {
                        show: true,
                        position: "insideLeft",
                        align: "left",
                        fontSize: 10,
                        formatter: (params) => {
                            let name = params.name
                            if (name.length > 25) {
                                return params.value + "kb  " + name.slice(0, 20) + "..." + name.slice(-5)
                            }
                            return params.value + "kb  " + name
                        },
                        rich: {
                            name: {}
                        }
                    },
                    type: 'bar'
                }
            ]
        };
    chart.setOption(option)
}

const getResult = (job: Job) => {
    job.getResult().then(resp => {
        const data = resp['data']
        countTcp.value = data["count_packet"]
        countConn.value = data["count_conn"]
        const connLst = data['lst']
        let reqCount = 0
        const urlSet = new Set(),
            statusMap = new Map(),
            connDataLst1 = []

        connLst.forEach(el => {
            const req = el['req']
            const status = el['status']
            const reqInfo = req.split(" "),
                statusInfo = (status && status !== "" ? status.split(" ") : ["", "", ""]),
                code = statusInfo[1]
            if (req && req !== "") {
                reqCount += 1
                urlSet.add(reqInfo[1])
            }
            if (status && status !== "") {
                if (statusMap.has(code)) {
                    statusMap.set(code, statusMap.get(code) + 1)
                } else {
                    statusMap.set(code, 1)
                }
            }
            connDataLst1.push({
                "url": reqInfo[1],
                "method": reqInfo[0],
                "status": code,
                "time": parseInt(el['time']),
                "throughput": (parseInt(el['throughput']) / 1024).toFixed(3),
            })
        });
        connDataLst.value = connDataLst1
        countReq.value = reqCount
        countUrl.value = urlSet.size

        RenderStatusStatChart(statusMap)
        RenderRespTimeStatChart(connDataLst1.sort((a, b) => {
            return b["time"] - a["time"]
        }).slice(0, 10).map(item => {
            return {
                "url": item['url'],
                "time": item["time"]
            }
        }))
        RenderThroughputStatChart(connDataLst1.sort((a, b) => {
            return b["throughput"] - a["throughput"]
        }).slice(0, 10).map(item => {
            return {
                "url": item['url'],
                "throughput": item["throughput"]
            }
        }))
    }).catch((resp) => {
        console.log(resp)
    })
}
</script>
<style lang="css" scoped>
.job_lst {
    width: 100%;
}

.graph {
    height: 20rem;
}
</style>