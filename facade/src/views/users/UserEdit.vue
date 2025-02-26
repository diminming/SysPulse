<template>

    <a-layout-content :style="{ background: '#fff', padding: '24px', margin: 0, minHeight: '800px' }">
        <a-alert :message="alarmMsg" type="warning" v-if="showAlart" show-icon style="margin: 1rem;" />

        <a-form layout="horizontal" :model="user" :label-col="{ style: { width: '6rem' } }" class="edit-form">
            <a-form-item label="用户名">
                <a-input v-model:value="user.username" placeholder="请输入用户名" />
            </a-form-item>
            <a-form-item label="密码">
                <a-input-password v-model:value="user.password" placeholder="请输入密码" />
            </a-form-item>
            <a-form-item label="确认密码">
                <a-input-password v-model:value="confirPassword" placeholder="请再次输入密码" />
            </a-form-item>
            <a-form-item label="角色">
                <a-button type="primary" size="small" @click="showRoleSelecting = !showRoleSelecting">选择</a-button>
                <template v-for="(role, idx) in user.roleLst" :key="role.id">
                    <div style="display: inline-block; margin-right: 1rem;">
                        <a-button type="link" size="small"> {{ role.name }}</a-button>
                        <a-button danger size="small" shape="circle" :icon="h(CloseOutlined)" @click="removeRole(idx)"/>
                    </div>
                </template>
            </a-form-item>
            <a-form-item label="状态">
                <a-radio-group v-model:value="user.isActive" button-style="solid" size="small">
                    <a-radio-button :value="true">启用</a-radio-button>
                    <a-radio-button :value="false">禁用</a-radio-button>
                </a-radio-group>
            </a-form-item>
            <a-form-item>
                <div style="text-align: center;">
                    <a-button type="primary" @click="onSave">保存</a-button>
                    <a-button @click="onCancel" style="margin-left: 2rem;">取消</a-button>
                </div>
            </a-form-item>
        </a-form>
    </a-layout-content>
    <a-modal v-model:open="showRoleSelecting" title="选择用户角色" width="70%">
        <role-lst stage="selecting" @select="onSelect"></role-lst>
    </a-modal>
</template>

<script lang="ts" setup>
import { h, onMounted, reactive, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router'

import { notification } from 'ant-design-vue';

import { Role, User } from "@/views/users/user"
import RoleLst from "@/views/users/RoleLst.vue"
import { CloseOutlined } from '@ant-design/icons-vue';


const route = useRoute(),
    router = useRouter(),
    showRoleSelecting = ref(false)

const user = reactive<User>(new User({ id: - 1, isActive: false })),
    confirPassword = ref(""),
    showAlart = ref(false),
    alarmMsg = ref("")

function onCancel() {
    router.push("/main/user")
}

function onSave() {
    if (confirPassword.value === user.password) {
        user.save().then(() => {
            notification.success({
                message: '保存成功',
                description:
                    '记录已保存。',
                duration: 2,
                onClick: () => {
                    // console.log('Notification Clicked!');
                },
            });
            router.push("/main/user")
        })
    } else {
        showAlart.value = true
        alarmMsg.value = "密码和确认密码不一致，请重新输入"
        setTimeout(() => {
            showAlart.value = false
        }, 3000)
    }
}

const onSelect = (role: any) => {
    user.roleLst.push(new Role({id: role.id, name: role.name}))
}

function removeRole(idx: number) {
    if (user.roleLst && user.roleLst.length > idx) {
        user.roleLst.splice(idx, 1);
    }
}

onMounted(() => {
    if (route.query.userId) {
        let userId = parseInt(route.query.userId as string)
        let user = new User({ id: userId })

        user.load().then((resp) => {
            const data = resp.data
        })
    }

})

</script>
<style scoped>
.edit-form {
    width: 30rem;
}
</style>