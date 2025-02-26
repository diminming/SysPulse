<template>

    <a-layout-content :style="{ background: '#fff', padding: '24px', margin: 0, minHeight: '800px' }">
        <a-form layout="horizontal" :model="role" :label-col="{ style: { width: '6rem' } }" class="edit-form">
            <a-form-item label="标识符">
                <a-input v-model:value="role.identity" placeholder="请输入标识符" />
            </a-form-item>
            <a-form-item label="角色名称">
                <a-input v-model:value="role.name" placeholder="请输入角色名称" />
            </a-form-item>
            <a-form-item label="权限列表">
                <a-button type="primary" size="small" @click="isShow = !isShow">选择</a-button><br/>
                <template v-for="(permission, idx ) in role.permissionLst" :key="permission.id">
                    <div style="display: inline-block; margin-right: 1rem;">
                        <a-button type="link" size="small"> {{ permission.name }}</a-button>
                        <a-button danger size="small" shape="circle" :icon="h(CloseOutlined)" @click="removePermission(idx)"/>
                    </div>
                </template>
            </a-form-item>
            <a-form-item>
                <div style="text-align: center;">
                    <a-button type="primary" @click="onSave">保存</a-button>
                    <a-button @click="onCancel" style="margin-left: 2rem;">取消</a-button>
                </div>
            </a-form-item>
        </a-form>
    </a-layout-content>
    <a-modal v-model:open="isShow" title="选择需要的权限" width="70%" @ok="isShow = !isShow">
        <permission-lst stage="selecting" @select="onSelect"></permission-lst>
    </a-modal>
</template>

<script lang="ts" setup>
import { onMounted, reactive, ref, h } from 'vue';
import { CloseOutlined } from '@ant-design/icons-vue';
import { useRoute, useRouter } from 'vue-router'

// import { message } from 'ant-design-vue';
import { notification } from 'ant-design-vue';

// import { Business as BizObj, Linux } from "@/views/linux/api"
import PermissionLst from './PermissionLst.vue';
import { Permission, Role, User } from "@/views/users/user"


const route = useRoute(),
    router = useRouter()

const role = reactive<Role>(new Role({ id: - 1 })),
    isShow = ref(false),
    parentTitle = ref("")

function onCancel() {
    router.push("/main/role")
}

function onSave() {
    role.save().then(() => {
        notification.success({
            message: '保存成功',
            description:
                '记录已保存。',
            duration: 2,
            onClick: () => {
                // console.log('Notification Clicked!');
            },
        });
        router.push("/main/role")
    })
}

const onSelect = (p: Permission) => {
    role.permissionLst.push(p)
}

function removePermission(idx: number) {
    if (role.permissionLst && role.permissionLst.length > idx) {
        role.permissionLst.splice(idx, 1);
    }
}

onMounted(() => {

})

</script>
<style scoped>
.edit-form {
    width: 30rem;
}
</style>