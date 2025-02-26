<template>
    <a-layout-content :style="{ background: '#fff', padding: '24px', margin: 0, minHeight: '800px' }">
        <a-form layout="horizontal" :model="menu" :label-col="{ style: { width: '6rem' } }" class="edit-form">
            <a-form-item label="标识符">
                <a-input v-model:value="menu.identity" placeholder="请输入菜单标识符" />
            </a-form-item>
            <a-form-item label="菜单标题">
                <a-input v-model:value="menu.title" placeholder="请输入菜单标题" />
            </a-form-item>
            <a-form-item label="URL">
                <a-input v-model:value="menu.url" placeholder="请输入URL" />
            </a-form-item>
            <a-form-item label="菜单类型">
                <a-radio-group v-model:value="menu.type" button-style="solid" size="small">
                    <a-radio-button value="group">菜单分组</a-radio-button>
                    <a-radio-button value="item">菜单项</a-radio-button>
                </a-radio-group>
            </a-form-item>
            <a-form-item label="菜单索引">
                <a-input-number v-model:value="menu.index" placeholder="请设置菜单位置" />
            </a-form-item>
            <a-form-item label="上级菜单">
                <template v-if="menu.parentTitle">
                    <a-button type="link" size="small">{{ menu.parentTitle }}</a-button>
                </template>
                <a-button size="small" type="primary" @click="isShow2 = !isShow2">选择</a-button>
            </a-form-item>
            <a-form-item label="权限列表">
                <a-button type="primary" size="small" @click="isShow = !isShow">选择</a-button><br/>
                <template v-for="(permission, idx ) in menu.permissionLst" :key="permission.id">
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
        <permission-lst stage="selecting" @select="onSelectPermission"></permission-lst>
    </a-modal>
    <a-modal v-model:open="isShow2" title="选择上级菜单" width="70%" @ok="isShow2 = !isShow2">
        <menu-lst stage="selecting" @select="onSelectParentMenu"></menu-lst>
    </a-modal>
</template>

<script lang="ts" setup>
import { onMounted, reactive, ref, h } from 'vue';
import { CloseOutlined } from '@ant-design/icons-vue';
import { useRouter } from 'vue-router'

// import { message } from 'ant-design-vue';
import { notification } from 'ant-design-vue';

// import { Business as BizObj, Linux } from "@/views/linux/api"
import PermissionLst from './PermissionLst.vue';
import MenuLst from './MenuLst.vue';
import { Permission, Menu } from "@/views/users/user"


const router = useRouter()

const menu = reactive<Menu>(new Menu({ id: - 1 })),
    isShow = ref(false),
    isShow2 = ref(false)

function onCancel() {
    router.push("/main/menu")
}

function onSave() {
    menu.save().then(() => {
        notification.success({
            message: '保存成功',
            description:
                '记录已保存。',
            duration: 2,
            onClick: () => {
                // console.log('Notification Clicked!');
            },
        });
        router.push("/main/menu")
    })
}

const onSelectParentMenu = (m: Menu) => {
    menu.parentId = m.id
    menu.parentTitle = m.title || ""
    isShow2.value = false
}

const onSelectPermission = (p: Permission) => {
    menu.permissionLst.push(p)
}

function removePermission(idx: number) {
    if (menu.permissionLst && menu.permissionLst.length > idx) {
        menu.permissionLst.splice(idx, 1);
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