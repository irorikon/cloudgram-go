<template>
  <div class="file-list" @contextmenu.prevent="handleGlobalContextMenu">
    <n-data-table :key="'file-table-' + tableRenderKey" ref="dataTableRef" :columns="columns" :data="files"
      :row-key="getRowKey" :bordered="false" :loading="loading" size="medium"
      @update:checked-row-keys="handleSelectionChange" :style="{ minHeight: minHeight }" />
  </div>

  <!-- 右键菜单 -->
  <n-dropdown ref="contextMenuRef" :options="contextMenuOptions" :show="showContextMenu" :x="contextMenuX"
    :y="contextMenuY" @select="handleContextMenuItemSelect" @clickoutside="hideContextMenu" placement="bottom-start" />
</template>

<script setup lang="ts">
import { ref, computed, h, onMounted, watch, nextTick } from 'vue';
import { NIcon, NDataTable, NDropdown } from 'naive-ui';
import {
  AddOutline,
  PencilOutline,
  TrashOutline,
  DownloadOutline,
  MoveOutline,
  CloudUploadOutline,
  RefreshOutline,
  CheckboxOutline
} from '@vicons/ionicons5';
import { useBreadcrumbStore } from '@/store/breadcrumb';
import { getIcon } from '@/utils/mimetype';
import type { FileItem } from '@/types/file';

// 定义组件属性
interface Props {
  files: FileItem[];
  loading: boolean;
}

const props = defineProps<Props>();

// 定义事件
interface Emits {
  (e: 'folder-click', file: FileItem): void;
  (e: 'file-operation', operation: string, file: FileItem | null): void;
  (e: 'selection-operation', operation: string, selectedFiles: FileItem[]): void;
  (e: 'file-click', file: FileItem): void;
}

const emit = defineEmits<Emits>();

// Naive UI hooks
const breadcrumbStore = useBreadcrumbStore();

// 响应式数据
const showContextMenu = ref(false);
const contextMenuX = ref(0);
const contextMenuY = ref(0);
const minHeight = ref('300px');
const currentFile = ref<FileItem | null>(null);
const selectedCount = ref(0);
const selectedFiles = ref<FileItem[]>([]);
const showSelection = ref(false);
// 表格实例
const dataTableRef = ref();
// 选中的行 keys
const selectedRowKeys = ref<(string)[]>([]);
// 表格重新渲染的 key
const tableRenderKey = ref(0);

// 表格列定义
const columns = computed(() => [
  ...(showSelection.value ? [{
    type: 'selection' as const,
    width: 36,
    align: 'center' as const
  }] : []),
  {
    title: selectedCount.value > 0 ? `已选择${selectedCount.value}项` : '文件名',
    key: 'name',
    align: 'left' as const,
    // ellipsis: true,
    render: (row: FileItem) => {
      return h('div', {
        class: 'file-item',
        onContextmenu: (event: MouseEvent) => handleRowContextMenu(event, row),
        onClick: (event: MouseEvent) => handleRowClick(row, event)
      }, [
        h(NIcon, {
          component: getIcon(row.is_dir, row.name, row.mime_type),
          class: ['file-list-icon', { 'folder-icon': row.is_dir }],
          color: "#1890ff",
          size: '1.5rem'
        }),
        h('span', {
          class: 'file-list-text',
          title: row.name,
        }, row.name)
      ]);
    }
  },
  {
    title: '大小',
    key: 'size',
    width: 120,
    align: 'right' as const,
    render: (row: FileItem) => {
      if (row.is_dir) {
        return '-';
      }
      return formatFileSize(row.size);
    }
  },
  {
    title: '修改时间',
    key: 'updated_at',
    width: 200,
    align: 'right' as const,
    render: (row: FileItem) => formatDate(row.UpdatedAt)
  },
]);

// 获取行键
const getRowKey = (row: FileItem) => row.id;

// 格式化日期
const formatDate = (date: string): string => {
  if (!date) return '-';

  const d = new Date(date);
  if (isNaN(d.getTime())) return '-';

  const year = d.getFullYear();
  const month = String(d.getMonth() + 1).padStart(2, '0');
  const day = String(d.getDate()).padStart(2, '0');
  const hours = String(d.getHours()).padStart(2, '0');
  const minutes = String(d.getMinutes()).padStart(2, '0');
  const seconds = String(d.getSeconds()).padStart(2, '0');

  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
};

// 格式化文件大小
const formatFileSize = (size: number): string => {
  if (size === 0) return '0 B';

  const units = ['B', 'KB', 'MB', 'GB', 'TB'];
  let unitIndex = 0;
  let fileSize = size;

  while (fileSize >= 1024 && unitIndex < units.length - 1) {
    fileSize /= 1024;
    unitIndex++;
  }

  return unitIndex === 0
    ? `${fileSize} ${units[unitIndex]}`
    : `${fileSize.toFixed(2)} ${units[unitIndex]}`;
};

// 处理选择变化
const handleSelectionChange = (keys: any[]) => {
  selectedRowKeys.value = keys;
  selectedFiles.value = props.files.filter(file => keys.includes(file.id));
  selectedCount.value = keys.length;
};

// 重置表格选择状态
const resetTableSelection = async () => {
  // 1. 清空选中的 keys
  selectedRowKeys.value = [];

  // 2. 更新内部状态
  selectedFiles.value = [];
  selectedCount.value = 0;

  // 3. 等待 DOM 更新
  await nextTick();

  // 4. 强制表格重新渲染
  tableRenderKey.value += 1;
};

// 监听文件列表变化
watch(() => props.files, async (newFiles, oldFiles) => {
  // 简单但有效的比较
  if (newFiles.length !== oldFiles.length ||
    JSON.stringify(newFiles.map(f => f.id)) !== JSON.stringify(oldFiles.map(f => f.id))) {
    await resetTableSelection();
  }
}, { deep: true });

// 处理文件点击
const handleRowClick = (row: FileItem, event: MouseEvent) => {
  // 如果是文件夹，则导航到文件夹路由
  if (row.is_dir) {
    breadcrumbStore.addCrumb({
      id: row.id,
      isDir: row.is_dir,
      mimeType: row.mime_type,
      name: row.name,
    });
    emit('folder-click', row);
  } else {
    // 否则，触发文件点击事件用于预览
    emit('file-click', row);
  }
};

// 计算并设置容器高度以填充可用空间
const updateContainerHeight = () => {
  // 获取视口高度
  const viewportHeight = window.innerHeight;

  // 计算已知组件占用的高度
  const headerHeight = 64;
  const breadcrumbHeight = 40;
  const footerHeight = 64;
  const margins = 40;
  const filelistheader = 50;

  // 计算文件列表可用高度
  const availableHeight = Math.max(
    viewportHeight - headerHeight - breadcrumbHeight - footerHeight - margins - filelistheader,
    300
  );
  // 设置容器高度
  minHeight.value = `${availableHeight}px`;
};

// 显示右键菜单
const showFileContextMenu = (event: MouseEvent, file: FileItem) => {
  currentFile.value = file;
  contextMenuX.value = event.clientX;
  contextMenuY.value = event.clientY;
  showContextMenu.value = true;
};

// 行右键菜单处理
const handleRowContextMenu = (event: MouseEvent, file: FileItem) => {
  event.preventDefault();
  event.stopPropagation();
  showFileContextMenu(event, file);
};

// 全局右键菜单处理（空白区域）
const handleGlobalContextMenu = (event: MouseEvent) => {
  currentFile.value = null;
  contextMenuX.value = event.clientX;
  contextMenuY.value = event.clientY;
  showContextMenu.value = true;
};

// 右键菜单选项
const contextMenuOptions = computed(() => {
  const options = [];

  const hasCurrentFile = currentFile.value

  if (selectedCount.value > 0) {
    options.push(
      {
        label: '移动',
        key: 'move',
        icon: () => h(NIcon, {
          component: MoveOutline,
          size: 24,
          color: '#ffb224'
        }),
      },
      {
        label: '多选',
        key: 'selection',
        icon: () => h(NIcon, {
          component: CheckboxOutline,
          size: 24,
          color: '#1890ff'
        }),
      }
    )
  } else if (hasCurrentFile) {
    options.push(
      {
        label: '重命名',
        key: 'rename',
        icon: () => h(NIcon, {
          component: PencilOutline,
          size: 24,
          color: '#6e56cf'
        }),
      },
      {
        label: '删除',
        key: 'delete',
        icon: () => h(NIcon, {
          component: TrashOutline,
          size: 24,
          color: '#e5484d'
        }),
      },
      {
        label: '移动',
        key: 'move',
        icon: () => h(NIcon, {
          component: MoveOutline,
          size: 24,
          color: '#ffb224'
        }),
      },
      {
        label: '下载',
        key: 'download',
        icon: () => h(NIcon, {
          component: DownloadOutline,
          size: 24,
          color: '#05a2c2'
        }),
      },
      {
        label: '多选',
        key: 'selection',
        icon: () => h(NIcon, {
          component: CheckboxOutline,
          size: 24,
          color: '#1890ff'
        }),
      }
    );
  } else {
    options.push(
      {
        label: '新建目录',
        key: 'new-folder',
        icon: () => h(NIcon, {
          component: AddOutline,
          size: 24,
          color: '#30a46c'
        }),
      },
      {
        label: '上传',
        key: 'upload',
        icon: () => h(NIcon, {
          component: CloudUploadOutline,
          size: 24,
          color: '#1890ff'
        }),
      },
      {
        label: '刷新',
        key: 'refresh',
        icon: () => h(NIcon, {
          component: RefreshOutline,
          size: 24,
          color: '#05a2c2'
        }),
      }
    );
  }

  return options;
});

// 处理右键菜单项选择
const handleContextMenuItemSelect = (key: string) => {
  if (key === 'selection') {
    showSelection.value = !showSelection.value;
  } else {
    if (showSelection.value && selectedRowKeys.value.length > 0) {
      emit('selection-operation', key, selectedFiles.value)
    } else {
      emit('file-operation', key, currentFile.value)
    }
  }
  hideContextMenu();
};

// 隐藏右键菜单
const hideContextMenu = () => {
  showContextMenu.value = false;
  currentFile.value = null;
};

// 暴露方法
defineExpose({
  refresh: () => { },
  getSelectedFiles: () => selectedFiles.value
});

onMounted(() => {
  updateContainerHeight()
})
</script>

<style scoped>
.file-list {
  width: 100%;
  min-height: 300px;
  background-color: var(--n-color);
  border-radius: var(--n-border-radius);
  overflow: hidden;
  border: 1px solid var(--n-border-color);
}

/* 强制 DataTable 使用固定布局 */
:deep(.n-data-table-table) {
  table-layout: fixed !important;
}

:deep(.n-data-table) {
  --n-border-color: var(--n-border-color);
  --n-color: var(--n-color);
  --n-th-color: var(--n-th-color);
  --n-th-color-hover: var(--n-th-color-hover);
  --n-th-color-striped: var(--n-th-color-striped);
  --n-td-color: var(--n-td-color);
  --n-td-color-hover: var(--n-td-color-hover);
  --n-td-color-striped: var(--n-td-color-striped);
}

:deep(.n-data-table-thead) {
  border-bottom: 2px solid var(--n-border-color);
}

:deep(.n-data-table-th) {
  font-weight: 600;
  color: var(--n-text-color-1);
  padding: 12px 8px;
}

:deep(.n-data-table-td) {
  padding: 10px 16px;
  /* 增加左右内边距到16px，为放大效果提供足够空间 */
  border-bottom: 1px solid var(--n-border-color);
}

:deep(.n-data-table-tr:last-child .n-data-table-td) {
  border-bottom: none;
}

:deep(.file-item) {
  display: flex;
  align-items: center;
  gap: var(--gap-medium);
  width: 100%;
  cursor: pointer;
  user-select: none;
  transition: background-color 0.2s ease;
}

/* 文件名文本截断样式 */
:deep(.file-list-text) {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis !important;
  white-space: nowrap;
}

/* 鼠标悬停放大吸附效果 - 只应用于数据行，不包括表头 */
:deep(.n-data-table-tbody .n-data-table-tr) {
  transition: transform 0.2s ease, box-shadow 0.2s ease;
  will-change: transform;
}

:deep(.n-data-table-tbody .n-data-table-tr:hover) {
  transform: scale(1.02);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  z-index: 1;
}
</style>