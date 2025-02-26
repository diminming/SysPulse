<template>

    <a-layout-content :style="{ background: '#fff', padding: '24px', margin: 0, minHeight: '800px' }">
        <a-form layout="horizontal" :model="permission" :label-col="{ style: { width: '6rem' } }" class="edit-form">
            <a-form-item label="标识符">
                <a-input v-model:value="permission.identity" placeholder="请输入标识符" />
            </a-form-item>
            <a-form-item label="名称">
                <a-input v-model:value="permission.name" placeholder="请输入权限名称" />
            </a-form-item>
            <a-form-item label="HTTP Method">
                <a-radio-group v-model:value="permission.method" button-style="solid" size="small" style="width: 40rem;">
                    <a-radio-button value="get">GET</a-radio-button>
                    <a-radio-button value="post">POST</a-radio-button>
                    <a-radio-button value="put">PUT</a-radio-button>
                    <a-radio-button value="delete">DELETE</a-radio-button>
                    <a-radio-button value="patch">PATCH</a-radio-button>
                    <a-radio-button value="head">HEAD</a-radio-button>
                    <a-radio-button value="connect">CONNECT</a-radio-button>
                    <a-radio-button value="options">OPTIONS</a-radio-button>
                    <a-radio-button value="trace">TRACE</a-radio-button>
                </a-radio-group>
            </a-form-item>
            <a-form-item label="URL">
                <a-input v-model:value="permission.url" placeholder="请输入URL" />
            </a-form-item>
            <a-form-item>
                <div style="text-align: center;">
                    <a-button type="primary" @click="onSave">保存</a-button>
                    <a-button @click="onCancel" style="margin-left: 2rem;">取消</a-button>
                </div>
            </a-form-item>
        </a-form>
    </a-layout-content>
</template>

<script lang="ts" setup>
import { onMounted, reactive, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router'

import { notification } from 'ant-design-vue';

import { Permission, User } from "@/views/users/user"


const route = useRoute(),
    router = useRouter()

const permission = reactive<Permission>(new Permission({ id: - 1, method: "get"}))

function onCancel() {
    router.push("/main/permission")
}

function onSave() {
    permission.save().then(() => {
        notification.success({
            message: '保存成功',
            description:
                '记录已保存。',
            duration: 2,
            onClick: () => {
                // console.log('Notification Clicked!');
            },
        });
        router.push("/main/permission")
    })
}


onMounted(() => {
    if (route.query.userId) {
        let userId = parseInt(route.query.userId as string)
        let linux0 = new User({ id: userId })

        linux0.load().then((resp) => {
            const data = resp.data
            // linux0.agent_conn = data['agent_conn']
            // linux0.hostname = data['hostname']
            // linux0.linux_id = data['linux_id']

            // const bizData = data['biz']
            // const biz0 = new BizObj(bizData['id'], bizData['bizName'])
            // linux0.biz = biz0
            // Object.assign(biz, biz0)

            // Object.assign(user, linux0)
        })
    }

})

</script>
<style scoped>
.edit-form {
    width: 30rem;
}
</style>