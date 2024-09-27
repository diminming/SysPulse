<template>
    <a-layout-content :style="{
        background: '#fff',
        padding: '24px',
        margin: 0,
        minHeight: '800px',
    }">
        <div class="opBar" v-if="stage !== 'dashboard'">
            <label class="search_label">时间范围</label>
            <a-range-picker style="width: 400px" show-time :format="dateTimeFormat" :presets="rangePresets"
                v-model:value="dateTimeRange" @change="onRangeChange" size="small" />
        </div>
        <div>
            <a-table :data-source="alarmData" :columns="columns" size="small" :row-selection="rowSelection"
                :pagination="pgSetting" @change="onChange">
                <template #bodyCell="{ text, column, record }">
                    <template v-if="column.key === 'linux'">
                        {{ record.linux.hostname }}
                    </template>
                    <template v-else-if="column.key === 'title'">
                        <a @click="showAlarmDetail(record)">
                            {{ record.trigger }}
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
                            <a-tag color="green">已确认</a-tag>
                        </span>
                        <span v-else style="font-weight: bold;;color: red;">
                            <a-tag color="red">未确认</a-tag>
                        </span>
                    </template>
                </template>
            </a-table>
        </div>
    </a-layout-content>
    <a-modal v-model:open="isShowDetail" title="消息详情" @ok="isShowDetail = false" width="1000px">
        <a-descriptions bordered :column="2">
            <a-descriptions-item label="消息标题" :span="2">
                <span style="">{{ alarm?.trigger }}</span>
            </a-descriptions-item>
            <a-descriptions-item label="消息时间">{{ dayjs(alarm?.timestamp).format("YYYY/MM/DD HH:mm:ss") }}</a-descriptions-item>
            <a-descriptions-item label="记录时间">{{ dayjs(alarm?.createTimestamp).format("YYYY/MM/DD HH:mm:ss") }}</a-descriptions-item>
            <a-descriptions-item label="消息主体" :span="2">{{ alarm?.linux?.hostname  }}</a-descriptions-item>
            <a-descriptions-item label="消息内容" :span="2">
                <span style="">{{ JSON.stringify(alarm?.perfData) }}</span>
            </a-descriptions-item>
        </a-descriptions>
    </a-modal>
</template>
<script lang="ts" setup>
import { computed, onMounted, reactive, defineProps, ref } from 'vue';
import { Alarm } from '.';
import dayjs from 'dayjs';
import { Linux } from '../linux/api';

const props = defineProps({
  stage: String
})

const pagination = reactive({
    page: 0,
    pageSize: props["stage"] == "dashboard" ? 200 : 20,
    total: 0,
})

const isShowDetail = ref(false)

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
        title: "消息时间",
        dataIndex: "timestamp",
        key: "timestamp",
    }, {
        title: "消息标题",
        dataIndex: "title",
        key: "title",
    }, {
        title: "消息主体",
        dataIndex: "linux",
        key: "linux",
    }, {
        title: "已确认",
        dataIndex: "ack",
        key: "ack",
    }, {
        title: "记录时间",
        dataIndex: "createTimestamp",
        key: "createTimestamp",
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
        a.perfData = data['perfData']
        isShowDetail.value = true
    })
    alarm.value = a
}

onMounted(() => {
    Alarm.loadPage(pagination).then((resp) => {
        alarmData.value = resp["data"]["lst"];
        pagination.total = resp["data"]["total"];
    })
})
</script>
<style lang="css" scoped>
.opBar {
    display: flex;
    margin-bottom: 1rem;
    width: 25%;
    justify-content: space-between;
}

.op-item {
    margin-right: 1rem;
}

.search_label {
    width: 5rem;
    font-weight: bold;
}
</style>