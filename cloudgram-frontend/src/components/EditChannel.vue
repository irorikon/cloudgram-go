<template>
    <div class="edit-channel-dialog">
        <!-- 频道表格 -->
        <n-data-table :columns="columns" :data="channels" :row-key="row => row.channel_id" :bordered="false"
            size="small" :pagination="false" max-height="300" class="channel-table" :row-props="rowProps" />

        <div v-if="channels.length === 0" class="empty-state">
            <n-text type="secondary">暂无频道</n-text>
        </div>

        <!-- 编辑/新增表单 -->
        <div class="channel-form">
            <n-form :model="form" :rules="rules" ref="formRef" label-placement="top" :label-width="80"
                class="centered-form">
                <n-form-item label="频道ID" path="channelId" class="form-item">
                    <n-input v-model:value="form.channelId" placeholder="请输入频道ID" :disabled="!!editingChannelId"
                        clearable />
                </n-form-item>

                <n-form-item label="频道名称" path="channelName" class="form-item">
                    <n-input v-model:value="form.channelName" placeholder="请输入频道名称" clearable />
                </n-form-item>
            </n-form>

            <div class="form-actions">
                <n-button size="medium" @click="handleCancel" class="form-btn">
                    清除
                </n-button>
                <n-button type="primary" size="medium" :loading="submitLoading" :disabled="!isFormValid"
                    @click="handleSubmit" class="form-btn form-btn--submit">
                    {{ editingChannelId ? '更新' : '新增' }}
                </n-button>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, h } from 'vue'
import {
    NButton,
    NInput,
    NForm,
    NFormItem,
    NText,
    NDataTable,
    useMessage,
    NIcon
} from 'naive-ui'
import type { FormInst, FormRules, DataTableColumns } from 'naive-ui'
import { listChannels, createChannel, deleteChannel, updateChannel, checkChannel } from '@/api/channel'
import type { RowData } from 'naive-ui/es/data-table/src/interface'
import { CheckmarkCircleOutline, CloseCircleOutline } from '@vicons/ionicons5'

// 定义频道类型
interface ChannelRecord {
    channel_id: number
    name: string
    limited: boolean
    message_id: number
    CreatedAt: string
    UpdatedAt: string
}

const message = useMessage()

// 状态管理
const channels = ref<ChannelRecord[]>([])
const selectedChannelId = ref<number | null>(null)
const editingChannelId = ref<number | null>(null)
const submitLoading = ref(false)
const deletingChannelId = ref<number | null>(null)
const formRef = ref<FormInst | null>(null)

// 表单数据
const form = ref({
    channelId: '',
    channelName: ''
})

// 表格列定义 - 移除固定宽度，让列自适应
const columns: DataTableColumns<ChannelRecord> = [
    {
        title: '频道名称',
        key: 'name',
        ellipsis: {  // 在列级别也可以设置
            tooltip: true
        },
        render: (row: ChannelRecord) => {
            return h('div', {
                style: {
                    display: 'flex',
                    alignItems: 'center',  // 垂直居中
                    gap: '4px'  // 图标和文字间距
                }
            }, [
                h(NIcon, {
                    component: getIcon(row.limited),
                    color: getColor(row.limited),
                    size: '1.5rem'
                }),
                h('span', {
                    style: {
                        overflow: 'hidden',
                        textOverflow: 'ellipsis',
                        whiteSpace: 'nowrap',
                        flex: 1,
                        minWidth: 0
                    }
                }, row.name)
            ]);
        }
    },
    {
        title: '频道 ID',
        key: 'channel_id',
        ellipsis: {  // 在列级别也可以设置
            tooltip: true
        },
    },
    {
        title: '操作',
        key: 'actions',
        render(row) {
            return h(NButton, {
                size: 'small',
                type: 'error',
                ghost: true,
                disabled: deletingChannelId.value === row.channel_id,
                onClick: (e: Event) => {
                    e.stopPropagation()
                    handleDeleteChannel(row.channel_id)
                }
            }, {
                default: () => deletingChannelId.value === row.channel_id ? '删除中...' : '删除'
            })
        }
    }
]

const getIcon = (limited: boolean) => {
    console.log('getIcon', limited)
    return limited ? CloseCircleOutline : CheckmarkCircleOutline
}

const getColor = (limited: boolean) => {
    return limited ? 'red' : 'green'
}

// 表单验证规则
const rules: FormRules = {
    channelId: [
        {
            required: true,
            message: '请输入频道ID',
            trigger: 'blur'
        },
        {
            validator: (rule, value) => {
                if (!value.trim()) return true
                const num = Number(value)
                return !isNaN(num) && Number.isInteger(num)
            },
            message: '频道ID必须是整数',
            trigger: 'blur'
        }
    ],
    channelName: [
        { required: true, message: '请输入频道名称', trigger: 'blur' },
        { max: 255, message: '频道名称不能超过255个字符', trigger: 'blur' }
    ]
}

// 计算属性
const isFormValid = computed(() => {
    const id = form.value.channelId.trim()
    const name = form.value.channelName.trim()

    if (!id || !name) return false

    const numId = Number(id)
    return !isNaN(numId) && Number.isInteger(numId) && name.length > 0
})

// 获取频道列表
const fetchChannels = async () => {
    try {
        const channelsData = await listChannels()
        channels.value = Array.isArray(channelsData) ? channelsData : []

        // 如果有频道且没有选中的频道，自动选择第一个
        if (channels.value.length > 0 && selectedChannelId.value === null) {
            const firstChannel = channels.value[0]
            if (firstChannel) {
                selectChannel(firstChannel.channel_id)
            }
        }
    } catch (error) {
        message.error('获取频道列表失败')
        channels.value = []
    }
}

function rowProps(row: RowData) {
    return {
        style: 'cursor: pointer;',
        onClick: () => {
            selectChannel(row.channel_id)
        }
    }
}
// 选择频道
const selectChannel = (channelId: number) => {
    selectedChannelId.value = channelId
    editingChannelId.value = channelId
    console.log('selectChannel', channelId)
    const channel = channels.value.find(c => c.channel_id === channelId)
    if (channel) {
        form.value = {
            channelId: channel.channel_id.toString(),
            channelName: channel.name
        }
    }
}

// 新增频道模式
const newChannelMode = () => {
    selectedChannelId.value = null
    editingChannelId.value = null
    form.value = {
        channelId: '',
        channelName: ''
    }
}

// 提交表单
const handleSubmit = async () => {
    if (!formRef.value) return

    try {
        await formRef.value.validate()
        submitLoading.value = true

        const channelIdNum = Number(form.value.channelId)

        if (isNaN(channelIdNum) || !Number.isInteger(channelIdNum)) {
            message.error('频道ID格式不正确，必须是整数')
            return
        }

        if (editingChannelId.value !== null) {
            // 更新逻辑
            await updateChannel(editingChannelId.value, form.value.channelName)
            message.success('频道更新成功')
        } else {
            // 新增频道
            const available = await checkChannel(channelIdNum)
            if (!available.channel_status) {
                message.error('频道不存在或已关闭')
            } else {
                await createChannel(channelIdNum, form.value.channelName)
                message.success('频道新增成功')
            }
        }

        // 重新获取列表并重置表单
        await fetchChannels()
        newChannelMode()
    } catch (error) {
        message.error(editingChannelId.value !== null ? '更新频道失败' : '新增频道失败')
    } finally {
        submitLoading.value = false
    }
}

// 删除频道
const handleDeleteChannel = async (channelId: number) => {
    if (window.confirm('确定要删除该频道吗？')) {
        try {
            deletingChannelId.value = channelId
            await deleteChannel(channelId)
            message.success('频道删除成功')

            // 重新获取列表
            await fetchChannels()

            // 如果删除的是当前选中的频道，切换到新增模式
            if (selectedChannelId.value === channelId) {
                newChannelMode()
            }
        } catch (error) {
            message.error('删除频道失败')
        } finally {
            deletingChannelId.value = null
        }
    }
}

// 取消操作
const handleCancel = () => {
    newChannelMode()
}

// 初始化
onMounted(() => {
    fetchChannels()
})
</script>

<style scoped>
.edit-channel-dialog {
    width: 100%;
    padding: 0 20px;
    box-sizing: border-box;
}

.channel-table {
    margin-bottom: 20px;
    border: 1px solid var(--border-color);
    border-radius: 6px;
}

.empty-state {
    text-align: center;
    padding: 24px;
    margin-bottom: 20px;
}

.channel-form {
    border: 1px solid var(--border-color);
    border-radius: 6px;
    padding: 20px;
}

/* 表单居中 */
.centered-form {
    display: flex;
    justify-content: center;
    flex-wrap: wrap;
    gap: 16px;
    margin-bottom: 20px;
}

.form-item {
    min-width: 200px;
    flex: 1;
}

/* 确保两个输入框宽度完全一致 */
:deep(.form-item .n-input) {
    width: 100%;
}

/* 移除可能导致宽度差异的默认样式 */
:deep(.n-form-item-blank) {
    flex: 1;
    min-width: 0;
}

.form-actions {
    display: flex;
    justify-content: center;
    gap: 12px;
    margin-top: 20px;
}

.form-btn {
    min-width: 80px;
}

/* 响应式设计 */
@media (max-width: 768px) {
    .edit-channel-dialog {
        padding: 0 16px;
    }

    .channel-table {
        font-size: 14px;
    }

    .channel-form {
        padding: 16px;
    }

    /* 小屏幕时表单项垂直排列 */
    .centered-form {
        flex-direction: column;
        align-items: stretch;
        gap: 16px;
    }

    .form-item {
        min-width: auto;
        flex: none;
    }

    .form-actions {
        flex-direction: column;
        gap: 8px;
    }

    .form-btn {
        width: 100%;
        min-width: auto;
    }
}

@media (min-width: 769px) {
    .edit-channel-dialog {
        padding: 0 30px;
    }

    /* 大屏幕时确保表单不会过宽 */
    .centered-form {
        max-width: 600px;
        margin: 0 auto 20px auto;
    }

    /* 确保两个表单项等宽 */
    .centered-form .form-item {
        flex: 1;
    }

    /* 强制输入框容器等宽 */
    :deep(.centered-form .n-form-item) {
        display: flex;
        flex-direction: column;
        flex: 1;
    }
}
</style>