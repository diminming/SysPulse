<template>
    <a-row :gutter="[8, 8]">
        <a-col :span="12">
            <a-row :gutter="[8, 8]">
                <a-card title="进程列表" size="small" style="width: 100%;">
                    <template #extra>
                        <a-button type="link" @click="getProcessLst">刷新</a-button>
                    </template>
                    <a-table :dataSource="procData" :columns="ProcColumnLst" size="small">
                        <template
                            #customFilterDropdown="{ setSelectedKeys, selectedKeys, confirm, clearFilters, column }">
                            <div style="padding: 8px">
                                <a-input ref="searchInput" :placeholder="`Search ${column.dataIndex}`"
                                    :value="selectedKeys[0]" style="width: 188px; margin-bottom: 8px; display: block"
                                    @change="(e: any) => setSelectedKeys(e.target.value ? [e.target.value] : [])"
                                    @pressEnter="handleSearch(selectedKeys, confirm, column.dataIndex)" />
                                <a-button type="primary" size="small" style="width: 90px; margin-right: 8px"
                                    @click="handleSearch(selectedKeys, confirm, column.dataIndex)">
                                    <template #icon>
                                        <SearchOutlined />
                                    </template>
                                    Search
                                </a-button>
                                <a-button size="small" style="width: 90px" @click="handleReset(clearFilters)">
                                    Reset
                                </a-button>
                            </div>
                        </template>
                        <template #bodyCell="{ text, column, record }">
                            <template v-if="column.key === 'name'">
                                <a @click="chooseProccess(record)">
                                    {{ record.name }}
                                </a>
                            </template>
                        </template>
                    </a-table>
                </a-card>
            </a-row>
            <a-row :gutter="[8, 8]">
                <a-card style="width: 100%; margin-top: 1rem;" title="分析任务" size="small">
                    <template #extra>
                        <a-button type="link" @click="showCreateJobDialog = true">新建</a-button>
                        <a-button type="link" @click="refreshJobLst">刷新</a-button>
                    </template>
                    <a-table :dataSource="jobData" :columns="JobColumnsLst" size="small">
                        <template #bodyCell="{ text, column, record }">
                            <template v-if="column.key === 'job_name'">
                                <a @click="getJobResult(record)">
                                    {{ record.job_name }}
                                </a>
                            </template>
                            <template v-if="column.key === 'status'">
                                {{ (record as Job).getStatusTxt() }}
                            </template>
                            <template v-if="column.key === 'startup_time'">
                                {{ (record as Job).getStartupTimeTxt() }}
                            </template>
                            <template v-if="column.key === 'immediately'">
                                {{ (record as Job).getImmediatelyTxt() }}
                            </template>
                            <template v-if="column.key === 'type'">
                                {{ (record as Job).getTypeTxt() }}
                            </template>
                            <template v-if="column.key === 'create_timestamp'">
                                {{ (record as Job).getCreateTimestampTxt() }}
                            </template>
                        </template>
                    </a-table>
                </a-card>
            </a-row>
        </a-col>
        <a-col :span="12">
            <div class="graph flame"></div>
        </a-col>
    </a-row>
    <a-modal v-model:open="showCreateJobDialog" title="创建剖析任务" :confirm-loading="confirmLoading" @ok="onCreateJob">
        <a-form>
            <a-form-item label="目标进程">
                <a-select v-model:value="selectProc" style="width: 220px" disabled>
                </a-select>
            </a-form-item>
            <a-form-item label="任务名称">
                <a-input v-model:value="job.job_name" />
            </a-form-item>
            <a-form-item label="任务类型">
                <a-select v-model:value="job.type" style="width: 120px">
                    <a-select-option value="oncpu">On CPU</a-select-option>
                    <a-select-option value="offcpu">Off CPU</a-select-option>
                </a-select>
            </a-form-item>
            <a-form-item label="开始时间">
                <a-date-picker :show-time="{ format: 'HH:mm' }" placeholder="Select Time" :disabled="job.immediately"
                    @ok="confirmStartupTime" />
                <a-checkbox style="margin-left: 1rem;" v-model:checked="job.immediately">即刻运行</a-checkbox>
            </a-form-item>
            <a-form-item label="采集时长">
                <a-input-number id="inputNumber" v-model:value="job.duration" :min="1" :max="10">
                    <template #addonAfter>
                        分钟
                    </template>
                </a-input-number>
            </a-form-item>
        </a-form>
    </a-modal>
</template>

<script lang="ts" setup>
import { onMounted, reactive, ref } from 'vue';
import { useRoute } from 'vue-router';
import { Linux, Job } from '@/views/linux/linux'
import type { SelectProps } from 'ant-design-vue/es/vc-select/Select';
import { SearchOutlined } from '@ant-design/icons-vue';
import { Dayjs } from 'dayjs';

const linux = ref<Linux>(new Linux(-1)),
    job = ref<Job>(new Job(-1)),
    showCreateJobDialog = ref<boolean>(false),
    confirmLoading = ref<boolean>(false),
    procData = ref([]),
    jobData = ref([]),
    searchInput = ref(),
    JobColumnsLst = [{
        title: "任务名称",
        dataIndex: "job_name",
        key: "job_name"
    }, {
        title: "任务状态",
        dataIndex: "status",
        key: "status"
    }, {
        title: "任务类型",
        dataIndex: "type",
        key: "type"
    }, {
        title: "启动时间",
        dataIndex: "startup_time",
        key: "startup_time"
    }, {
        title: "立即启动",
        dataIndex: "immediately",
        key: "immediately"
    }, {
        title: "任务时长(分钟)",
        dataIndex: "duration",
        key: "duration"
    }, {
        title: "创建时间",
        dataIndex: "create_timestamp",
        key: "create_timestamp"
    }],
    ProcColumnLst = [{
        title: '进程号',
        dataIndex: 'pid',
        key: 'pid',
    },
    {
        title: '进程名',
        dataIndex: 'name',
        key: 'name',
        customFilterDropdown: true,
        onFilter: (value: any, record: any) => record.name.toString().toLowerCase().includes(value.toLowerCase()),
        onFilterDropdownOpenChange: (visible: any) => {
            if (visible) {
                setTimeout(() => {
                    searchInput.value.focus();
                }, 100);
            }
        },
    },
    {
        title: '父进程',
        dataIndex: 'ppid',
        key: 'ppid',
    },
    {
        title: '命令行',
        dataIndex: 'exec',
        key: 'exec',
    }],
    selectProc = ref<SelectProps['options']>();

const getJobResult = (record: any) => {
    const job = new Job(record.id)
    job.RenderStackTraceChart()
}

const state = reactive({
    searchText: '',
    searchedColumn: '',
});

const handleSearch = (selectedKeys: any, confirm: any, dataIndex: any) => {
    confirm();
    state.searchText = selectedKeys[0];
    state.searchedColumn = dataIndex;
};

const handleReset = (clearFilters: any) => {
    clearFilters({ confirm: true });
    state.searchText = '';
};

const confirmStartupTime = (value: Dayjs) => {
    job.value.startup_time = value.valueOf()
};

function onCreateJob() {
    if (selectProc.value) {
        job.value.pid = selectProc.value[0]["value"]
        job.value.linux_id = linux.value.id
        job.value.category = "proc_profiling"
        confirmLoading.value = true
        job.value.createJob().then(resp => {
            confirmLoading.value = false
            showCreateJobDialog.value = false
            if(resp['status'] === 200) {
                debugger
                if(selectProc.value && selectProc.value.length > 0)
                    getJobLst({"pid": selectProc.value[0]["value"]})
            }
        })
    } else {
        alert("请在进程列表中选择要分析的进程...")
    }
}

function chooseProccess(record: any) {
    selectProc.value = [{
        value: record["pid"],
        label: record["name"]
    }]
    getJobLst(record);
}

function refreshJobLst() {
    if(selectProc.value) {
        const record = {
            pid: selectProc.value[0]["value"]
        }
        getJobLst(record)
    }
}

function getJobLst(record: any) {
    linux.value.GetAnalyzationJobLst(record['pid']).then((resp) => {
        const data = resp.data
        if (data !== null) {
            jobData.value = data.map((item: any) => {
                let job = new Job(item["id"])
                job.status = item['status']
                job.job_name = item['job_name']
                job.type = item['type']
                job.immediately = (item['immediately'] === 1)
                job.duration = item['duration']
                job.startup_time = item["startup_time"]
                job.create_timestamp = item['create_timestamp']
                return job
            })
        } else {
            jobData.value = []
        }

    })
}

function getProcessLst() {
    linux.value.GetProcessLst().then((resp => {
        const data = resp.data
        data.reverse()
        procData.value = data
    }))
}

onMounted(() => {
    let route = useRoute()
    let l = new Linux(parseInt(route.query.linuxId as string))
    linux.value = l
    getProcessLst()
    job.value.immediately = true
})


</script>

<style scoped>
.flame {
    height: 100%;
}
</style>