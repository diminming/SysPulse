<template>
    <a-layout-content :style="{
        background: '#fff',
        padding: '24px',
        margin: 0,
        minHeight: '800px',
    }">
        <div class="opBar" v-if="stage !== 'dashboard'">
            <a-form ref="formRef" class="searchForm" :model="formState" size="small" @finish="onFinish">
                <a-row :gutter="12">
                    <a-col :span="6">
                        <a-form-item name="timeRange" label="告警时间">
                            <a-range-picker v-model:value="formState['timeRange']" show-time
                                format="YYYY-MM-DD HH:mm:ss" value-format="YYYY-MM-DD HH:mm:ss" />
                        </a-form-item>
                    </a-col>
                    <a-col :span="6">
                        <a-form-item name="" label="告警对象">
                            <a-tag closable v-for="linux in formState.linuxLst" :key="linux.id" color="blue">{{
                                linux.hostname }}</a-tag>
                            <a-button @click="showLinuxLst = !showLinuxLst">选择</a-button>
                        </a-form-item>
                    </a-col>
                    <a-col :span="4">
                        <a-form-item name="" label="告警状态">
                            <a-select v-model:value="formState['status']">
                                <a-select-option value="inactive">已恢复</a-select-option>
                                <a-select-option value="active">告警中</a-select-option>
                            </a-select>
                        </a-form-item>
                    </a-col>

                    <a-col :span="1">
                        <a-button type="primary" html-type="submit">检索</a-button>
                    </a-col>
                </a-row>
            </a-form>
        </div>
        <div>
            <a-table :data-source="alarmData" :columns="columns" size="small" :pagination="pgSetting">
                <template #bodyCell="{ text, column, record }">
                    <template v-if="column.key === 'linux'">
                        {{ record.linux.hostname }}
                    </template>
                    <template v-else-if="column.key === 'msg'">
                        <a @click="showAlarmDetail(record)">
                            {{ record.msg }}
                        </a>
                    </template>
                    <template v-else-if="column.key === 'timestamp'">
                        {{ dayjs(record.timestamp).format("YYYY/MM/DD HH:mm:ss") }}
                    </template>
                    <template v-else-if="column.key === 'createTimestamp'">
                        {{ dayjs(record.createTimestamp).format("YYYY/MM/DD HH:mm:ss") }}
                    </template>
                    <template v-else-if="column.key === 'ack'">
                        <span v-if="record.ack === true" style="font-weight: bold;color: green;">
                            <a-tag color="green">已恢复</a-tag>
                        </span>
                        <span v-else style="font-weight: bold;;color: red;">
                            <a-popconfirm title="您正在手动关闭一个告警，是否确认？" ok-text="确认" cancel-text="取消"
                                @confirm="disableAlarm(record.id)">
                                <a-tag color="red">告警中</a-tag>
                            </a-popconfirm>
                        </span>
                    </template>
                </template>
            </a-table>
        </div>
    </a-layout-content>
    <a-modal v-model:open="isShowDetail" title="消息详情" @ok="isShowDetail = false" width="1000px">
        <a-descriptions bordered :column="2">
            <a-descriptions-item label="消息" :span="2">
                <span style="">{{ alarm?.msg }}</span>
            </a-descriptions-item>
            <a-descriptions-item label="Trigger ID">{{ alarm?.triggerId }}</a-descriptions-item>
            <a-descriptions-item label="Trigger">{{ alarm?.trigger }}</a-descriptions-item>
            <a-descriptions-item label="消息时间">
                {{ dayjs(alarm?.timestamp).format("YYYY/MM/DD HH:mm:ss") }}</a-descriptions-item>
            <a-descriptions-item label="记录时间">{{ dayjs(alarm?.createTimestamp).format("YYYY/MM/DD HH:mm:ss")
                }}</a-descriptions-item>
            <a-descriptions-item label="产生对象">{{ alarm?.linux?.hostname }}</a-descriptions-item>
            <a-descriptions-item label="消息状态">
                <span v-if="alarm?.ack === true" style="font-weight: bold;color: green;">
                    <a-tag color="green">已恢复</a-tag>
                </span>
                <span v-else style="font-weight: bold;;color: red;">
                    <a-tag color="red">生效中</a-tag>
                </span>
            </a-descriptions-item>
            <!-- <a-descriptions-item label="" :span="2">
                <span style="">{{ alarm?.msg }}</span>
            </a-descriptions-item> -->
        </a-descriptions>
    </a-modal>
    <a-modal v-model:open="showLinuxLst" title="选择Linux对象" @ok="showLinuxLst = false" width="1000px">
        <LinuxLst stage="select" @select="onSelectLinux"></LinuxLst>
    </a-modal>
</template>
<script lang="ts" setup>
import { computed, onMounted, reactive, defineProps, ref } from 'vue';
import { Alarm } from '.';
import dayjs from 'dayjs';
import type { FormInstance } from 'ant-design-vue';
import { Linux } from '../linux/api';
import LinuxLst from "@/views/linux/LinuxLst.vue"

const props = defineProps({
    stage: String
})

const showLinuxLst = ref(false)

const formRef = ref<FormInstance>()
const pagination = reactive({
    page: 0,
    pageSize: props["stage"] == "dashboard" ? 200 : 15,
    total: 0,
})

const formState = reactive<{
    timeRange: number[],
    status: String,
    linuxLst: Linux[]
}>({
    "timeRange": [0, 0],
    status: "",
    "linuxLst": []
}), onSelectLinux = (linuxLst: Linux[]) => {
    console.log(linuxLst)
    formState.linuxLst = linuxLst
}
const onFinish = () => {
    console.log(formState)
    Alarm.loadPage({
        page: 0,
        pageSize: pagination.pageSize,
        from: (() => formState.timeRange[0] !== 0 ? dayjs(formState.timeRange[0], "YYYY/MM/DD HH:mm:ss").valueOf() : undefined)(),
        util: (() => formState.timeRange[1] !== 0 ? dayjs(formState.timeRange[1], "YYYY/MM/DD HH:mm:ss").valueOf() : undefined)(),
        status: (() => formState.status === "" ? undefined : formState.status)(),
        target: (() =>
            formState.linuxLst.length === 0 ?
                undefined :
                formState.linuxLst.map(linux => {
                    return linux.id
                }).join(",")
        )()
    }).then((resp) => {
        console.log(resp)
        alarmData.value = resp["data"]["lst"];
        pagination.total = resp["data"]["total"];
    })
}

const isShowDetail = ref(false)

const disableAlarm = (alarmId: number) => {
    new Alarm(alarmId).disable().then(() => {
        Alarm.loadPage({
            page: pagination.page,
            pageSize: pagination.pageSize
        }).then((resp) => {
            alarmData.value = resp["data"]["lst"];
            pagination.total = resp["data"]["total"];
        })
    })
}

const pgSetting = computed(() => ({
    total: pagination.total,
    current: pagination.page + 1,
    pageSize: pagination.pageSize,
    showQuickJumper: true,
    showSizeChanger: true,
    showTotal: (total: number, range: number) => {
        return range + ", 共" + total;
    },
}));

const alarmData = ref([])
const alarm = ref<Alarm>()

const columns = [
    {
        title: "告警时间",
        dataIndex: "timestamp",
        key: "timestamp",
        width: 160
    }, {
        title: "告警对象",
        dataIndex: "linux",
        key: "linux",
        width: 160
    }, {
        title: "告警内容",
        dataIndex: "msg",
        key: "msg",
    }, {
        title: "告警状态",
        dataIndex: "ack",
        key: "ack",
        width: 100
    }
]

const showAlarmDetail = (record: any) => {
    const a = new Alarm(record['id'])
    a.load().then(resp => {
        const data = resp["data"]
        a.ack = data['ack']
        a.timestamp = data['timestamp']
        a.createTimestamp = data['createTimestamp']
        a.linux = new Linux(data["linux"]["id"], data["linux"]["hostname"])
        a.trigger = data['trigger']
        a.triggerId = data['triggerId']
        a.perfData = data['perfData']
        a.msg = data["msg"]
        isShowDetail.value = true
    })
    alarm.value = a
}

onMounted(() => {
    Alarm.loadPage({
        page: pagination.page,
        pageSize: pagination.pageSize
    }).then((resp) => {
        alarmData.value = resp["data"]["lst"];
        pagination.total = resp["data"]["total"];
    })
})
</script>
<style lang="css" scoped>
.opBar {
    display: flex;
    margin-bottom: 1rem;
    width: 100%;
    /* justify-content: space-between; */
}

.searchForm {
    width: 100%;
}

.op-item {
    margin-right: 1rem;
}

.search_label {
    width: 5rem;
    font-weight: bold;
}
</style>