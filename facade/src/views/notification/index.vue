<template>
    <a-layout-content :style="{
        background: '#fff',
        padding: '24px',
        margin: 0,
        minHeight: '800px',
    }">
        <div class="opBar">
            <label class="search_label">时间范围</label>
            <a-range-picker style="width: 400px" show-time :format="dateTimeFormat" :presets="rangePresets"
                v-model:value="dateTimeRange" @change="onRangeChange" size="small" />
        </div>
        <div>
            <a-table :data-source="tabData" :columns="columns" size="small" :row-selection="rowSelection"
                :pagination="pgSetting" @change="onChange">
            </a-table>
        </div>
    </a-layout-content>
</template>
<script lang="ts" setup>
import { computed, onMounted, reactive } from 'vue';
import { Alarm } from '.';

const pagination = reactive({
    page: 0,
    pageSize: 20,
    total: 0,
})

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

const columns = [
    {
        title: "消息时间",
        dataIndex: "timestamp",
        key: "timestamp",
    }, {
        title: "消息主体",
        dataIndex: "linux",
        key: "timestamp",
    }, {
        title: "已确认",
        dataIndex: "timestamp",
        key: "timestamp",
    }, {
        title: "生成时间",
        dataIndex: "timestamp",
        key: "timestamp",
    }
]

onMounted(() => {
    Alarm.loadPage(pagination).then((resp) => {
        console.log(resp)
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