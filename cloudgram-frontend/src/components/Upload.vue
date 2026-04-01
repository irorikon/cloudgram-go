<template>
  <div class="upload-modal-content">
    <div class="upload-header">
      <h3>上传文件</h3>
    </div>

    <div class="upload-content">
      <template v-if="!currentUploadingFile">
        <n-upload ref="uploadRef" v-model:file-list="fileListRef" :show-file-list="false" :multiple="true"
          :directory-dnd="true" :custom-request="handleCustomRequest" @finish="handleUploadFinish"
          @change="handleFileChange" accept="*">
          <n-upload-dragger>
            <div style="margin-bottom: 12px">
              <n-icon size="48" :depth="3">
                <CloudUploadOutline />
              </n-icon>
            </div>
            <n-text>
              将文件拖拽到该区域来上传
            </n-text>
            <n-p depth="3">
              或点击选择文件
            </n-p>
          </n-upload-dragger>
        </n-upload>
      </template>

      <!-- 自定义上传状态显示 -->
      <div v-else class="upload-status">
        <span class="uploading-filename">正在上传：{{ currentUploadingFile.name }}</span>
        <n-progress :percentage="currentUploadingFile.percentage || 0" :indicator-placement="'inside'"
          :show-indicator="true" :height="20" class="upload-progress" />
        <div v-if="remainingFilesCount > 0 && nextFileName" class="remaining-info">
          <span>剩余项目({{ remainingFilesCount }})</span>
          <span>{{ nextFileName }}</span>
        </div>
      </div>
    </div>

    <div class="upload-footer">
      <n-button @click="handleClose">关闭</n-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onUnmounted } from 'vue';
import { NUpload, NUploadDragger, NIcon, NText, NP, NButton, useMessage, NProgress } from 'naive-ui';
// 保留必要的类型导入
import type { UploadFileInfo, UploadCustomRequestOptions } from 'naive-ui';
import { CloudUploadOutline } from '@vicons/ionicons5';
import { useBreadcrumbStore } from '@/store/breadcrumb';
import { useChannelStore } from '@/store/channel';
import { uploadChunk, mergeFile, cleanupUploadSession } from '@/api/upload';
import { exists, createFile } from '@/api/file';

// 定义组件事件
interface Emits {
  (e: 'refresh-filelist'): void;
  (e: 'cancel'): void;
}

interface FolderInfo {
  id: string;
  name: string;
  parentId: string | null;
  isDir: boolean;
}

interface ParseFilePathResult {
  fileName: string;
  folderPath: string | null;
  folderNames: string[];
  isInFolder: boolean;
  originalFile: File | null;
}

const emit = defineEmits<Emits>();

// 定义常量
const CHUNK_SIZE = 10 * 1024 * 1024; // 10MB 分片大小

// 组件引用
const message = useMessage();
const breadcrumbStore = useBreadcrumbStore();

// 从 channel store 获取 channelId
const channelStore = useChannelStore();
const channel = channelStore.getChannel;
if (!channelStore.hasChannel || typeof channel !== 'object' || !('channelId' in channel)) {
  throw new Error('Channel not available');
}
const channelId = channel.channelId;

// 定义 messageId
const messageId = ref<number>(0);

// 添加文件列表状态管理（内部使用）
const fileListRef = ref<UploadFileInfo[]>([]);

// 当前正在上传的文件状态（用于自定义显示）
const currentUploadingFile = ref<UploadFileInfo | null>(null);

// 文件夹缓存列表
const folderCache = ref<FolderInfo[]>([]);

// 正在创建文件夹的 Promise 缓存，避免重复创建
const creatingFolders = new Map<string, Promise<FolderInfo>>();

// 上传模式
const uploadMode = ref<boolean>(false); // false: 文件上传 true: 文件夹上传

// 使用 Set 记录已创建的第一层目录
const createdFirstLevelFolders = new Set<string>();

// 计算剩余文件数量和下一个文件名
const remainingFilesCount = computed(() => {
  const uploadingIndex = fileListRef.value.findIndex(file => file.status === 'uploading');
  if (uploadingIndex === -1) return 0;
  return fileListRef.value.length - uploadingIndex - 1;
});

const nextFileName = computed(() => {
  const uploadingIndex = fileListRef.value.findIndex(file => file.status === 'uploading');
  if (uploadingIndex === -1 || uploadingIndex + 1 >= fileListRef.value.length) return '';
  const nextFile = fileListRef.value[uploadingIndex + 1];
  return nextFile ? nextFile.name : '';
});

// 串行上传相关状态
const uploadQueue = ref<UploadCustomRequestOptions[]>([]);
const isUploading = ref(false);

// 解析文件路径，提取文件夹信息
const parseFilePath = (uploadFileInfo: UploadFileInfo): ParseFilePathResult => {
  // 优先使用 fullPath（如果存在）
  if (uploadFileInfo.fullPath && uploadFileInfo.fullPath.includes('/')) {
    const path = uploadFileInfo.fullPath;
    const segments = path.split('/');
    const fileName = segments.pop() || uploadFileInfo.name;
    const folderPath = segments.join('/');

    return {
      fileName,
      folderPath,
      folderNames: segments,
      isInFolder: segments.length > 0,
      originalFile: uploadFileInfo.file || null
    };
  }

  // 其次使用 file 对象的 webkitRelativePath（目录上传）
  if (uploadFileInfo.file?.webkitRelativePath && uploadFileInfo.file.webkitRelativePath.includes('/')) {
    const path = uploadFileInfo.file.webkitRelativePath;
    const segments = path.split('/');
    const fileName = segments.pop() || uploadFileInfo.name;
    const folderPath = segments.join('/');

    return {
      fileName,
      folderPath,
      folderNames: segments,
      isInFolder: segments.length > 0,
      originalFile: uploadFileInfo.file || null
    };
  }

  // 如果文件名包含路径分隔符
  if (uploadFileInfo.name.includes('/')) {
    const segments = uploadFileInfo.name.split('/');
    const fileName = segments.pop() || uploadFileInfo.name;
    const folderPath = segments.join('/');

    return {
      fileName,
      folderPath: folderPath || null,
      folderNames: segments,
      isInFolder: segments.length > 0,
      originalFile: uploadFileInfo.file || null
    };
  }

  // 无文件夹结构
  return {
    fileName: uploadFileInfo.name,
    folderPath: null,
    folderNames: [],
    isInFolder: false,
    originalFile: uploadFileInfo.file || null
  };
};

// 在缓存中查找文件夹
const findFolderInCache = (name: string, parentId: string): FolderInfo | undefined => {
  return folderCache.value.find(folder =>
    folder.name === name && folder.parentId === parentId
  );
};

// 查找或创建文件夹
const findOrCreateFolder = async (folderName: string, parentId: string): Promise<FolderInfo> => {
  const cacheKey = `${parentId}_${folderName}`;

  // 检查是否已经在创建中
  if (creatingFolders.has(cacheKey)) {
    return creatingFolders.get(cacheKey)!;
  }

  // 检查是否已存在
  const existingFolder = findFolderInCache(folderName, parentId);
  if (existingFolder) {
    return existingFolder;
  }

  // 创建文件夹
  const createPromise = (async (): Promise<FolderInfo> => {
    try {
      // 首先检查文件夹是否已存在
      const existsResponse = await exists(folderName, parentId);

      if (existsResponse.exists && existsResponse.file) {
        // 文件夹已存在，从响应中获取ID
        const folderInfo: FolderInfo = {
          id: existsResponse.file.id,
          name: folderName,
          parentId,
          isDir: true
        };

        folderCache.value.push(folderInfo);
        return folderInfo;
      }

      // 创建文件夹
      const response = await createFile(folderName, parentId);

      const folderInfo: FolderInfo = {
        id: response.id,
        name: response.name,
        parentId: response.parent_id,
        isDir: response.is_dir
      };

      folderCache.value.push(folderInfo);
      return folderInfo;
    } catch (error) {
      creatingFolders.delete(cacheKey);
      message.error(`创建文件夹 "${folderName}" 失败: ${error instanceof Error ? error.message : '未知错误'}`);
      throw error;
    }
  })();

  creatingFolders.set(cacheKey, createPromise);
  return createPromise;
};

// 递归创建嵌套文件夹
const createNestedFolders = async (folderPath: string, baseParentId: string): Promise<string> => {
  if (!folderPath) {
    return baseParentId;
  }

  const pathSegments = folderPath.split('/').filter(Boolean);
  let currentParentId = baseParentId;
  let lastFolder: FolderInfo | null = null;

  for (const segment of pathSegments) {
    // 生成第一层目录的唯一标识
    const firstLevelKey = `${baseParentId}/${segment}`;
    // 查找当前层级的文件夹
    const existingFolder = findFolderInCache(segment, currentParentId);

    if (existingFolder) {
      lastFolder = existingFolder;
      currentParentId = existingFolder.id;
    } else {
      // 创建文件夹
      const newFolder = await findOrCreateFolder(segment, currentParentId);
      lastFolder = newFolder;
      currentParentId = newFolder.id;
      // 只有在创建第一层目录且尚未创建过时才刷新
      if (pathSegments.indexOf(segment) === 0 && !createdFirstLevelFolders.has(firstLevelKey)) {
        createdFirstLevelFolders.add(firstLevelKey);
        console.log(`Created folder: ${newFolder.name}`);
        emit('refresh-filelist');
      }
    }
  }

  return lastFolder?.id || baseParentId;
};

// 处理文件变化事件
const handleFileChange = (data: { file: UploadFileInfo; fileList: UploadFileInfo[] }) => {
  // 创建新的文件列表，确保每个文件都有正确的状态
  const newFileList: UploadFileInfo[] = data.fileList.map(newFile => {
    // 查找已存在的文件状态
    const existingFile = fileListRef.value.find(f => f.id === newFile.id);
    if (existingFile && (existingFile.status === 'uploading' || existingFile.status === 'finished' || existingFile.status === 'error')) {
      // 如果文件正在上传或已完成，保留现有状态
      return existingFile;
    }
    // 新文件或 pending 状态的文件
    return {
      ...newFile,
      status: 'pending',
      percentage: 0
    };
  });
  // 更新文件列表
  fileListRef.value = newFileList;

  // 更新当前正在上传的文件
  const uploadingFile = fileListRef.value.find(file => file.status === 'uploading');
  currentUploadingFile.value = uploadingFile || null;
};

// 自定义上传请求处理 - 修改为串行上传
const handleCustomRequest = async (options: UploadCustomRequestOptions) => {
  // 将上传任务添加到队列
  uploadQueue.value.push(options);

  // 如果没有正在上传，则开始处理队列
  if (!isUploading.value) {
    processUploadQueue();
  }
};

// 处理上传队列
const processUploadQueue = async () => {
  if (uploadQueue.value.length === 0 || isUploading.value) {
    return;
  }

  isUploading.value = true;
  const options = uploadQueue.value.shift()!;
  const { file, onFinish, onError } = options;

  if (!file.file) {
    onError();
    isUploading.value = false;
    // 继续处理下一个文件
    await processUploadQueue();
    return;
  }

  try {
    // 找到 fileListRef 中对应的文件对象
    const existingFileIndex = fileListRef.value.findIndex(f => f.id === file.id);
    let targetFile: UploadFileInfo;

    if (existingFileIndex !== -1) {
      const foundFile = fileListRef.value[existingFileIndex];
      if (foundFile) {
        targetFile = foundFile;
      } else {
        // 如果没找到（理论上不应该发生，但为了类型安全）
        targetFile = {
          id: file.id,
          name: file.name,
          status: 'uploading',
          percentage: 0,
          file: file.file
        };
        fileListRef.value.push(targetFile);
      }
    } else {
      // 如果没找到，创建一个新的文件信息对象
      targetFile = {
        id: file.id,
        name: file.name,
        status: 'uploading',
        percentage: 0,
        file: file.file
      };
      fileListRef.value.push(targetFile);
    }

    // 更新当前正在上传的文件
    currentUploadingFile.value = targetFile;

    // 解析文件路径
    const { fileName, folderPath, originalFile } = parseFilePath(file);

    if (!originalFile) {
      throw new Error('文件不存在');
    }

    // 获取基础父目录ID
    const baseParentId = breadcrumbStore.lastCrumbTypeDir()?.id || '';

    // 获取目标文件夹ID
    let targetFolderId: string;
    if (folderPath) {
      // 更新上传模式为文件夹上传
      uploadMode.value = true;
      // 有文件夹路径，递归创建文件夹
      try {
        targetFolderId = await createNestedFolders(folderPath, baseParentId);
        // 更新文件名为实际的文件名（去掉路径）
        targetFile.name = fileName;
        file.name = fileName;
      } catch (error) {
        console.error('创建文件夹失败，将使用父目录上传:', error);
        targetFolderId = baseParentId;
      }
    } else {
      // 更新上传模式为文件上传
      uploadMode.value = false;
      // 无文件夹路径，使用当前父目录
      targetFolderId = baseParentId;
    }

    // 首先检查文件是否已存在
    const fileExistsResponse = await exists(fileName, targetFolderId);
    if (fileExistsResponse.exists) {
      // 文件已存在，跳过上传，设置状态为 finished
      targetFile.status = 'finished';
      targetFile.percentage = 100;
      message.warning(`文件 "${fileName}" 已存在，跳过上传`);
      onFinish();
    } else {
      // 文件不存在，执行正常上传流程
      targetFile.status = 'uploading';
      targetFile.percentage = 0;

      const fileSize = originalFile.size;
      const totalChunks = fileSize <= CHUNK_SIZE ? 1 : Math.ceil(fileSize / CHUNK_SIZE);

      // 前端自己生成 uploadId，使用兼容性更好的方法
      const uploadId = generateUploadId();

      // 将 uploadId 存储到文件对象中
      (file as any).uploadId = uploadId;

      await processFileUploadWithId(originalFile, uploadId, targetFolderId, fileName, targetFile);

      // 设置文件状态为 finished
      targetFile.status = 'finished';
      targetFile.percentage = 100;
      onFinish();
    }
  } catch (error) {
    console.error('Upload error:', error);
    // 找到对应的文件对象并设置状态为 error
    const errorFileIndex = fileListRef.value.findIndex(f => f.id === file.id);
    if (errorFileIndex !== -1) {
      const errorFile = fileListRef.value[errorFileIndex];
      if (errorFile) {
        errorFile.status = 'error';
        // 显示错误文件名
        const { fileName } = parseFilePath(file);
        errorFile.name = fileName;
      }
    } else {
      // 如果没找到，添加一个错误状态的文件
      const { fileName } = parseFilePath(file);
      fileListRef.value.push({
        id: file.id,
        name: fileName,
        status: 'error',
        percentage: 0,
        file: file.file
      });
    }
    onError();
  } finally {
    isUploading.value = false;
    // 更新当前正在上传的文件（如果还有文件在上传）
    const stillUploading = fileListRef.value.find(f => f.status === 'uploading');
    currentUploadingFile.value = stillUploading || null;
    // 继续处理队列中的下一个文件
    await processUploadQueue();
  }
};

// 生成唯一上传ID的兼容性方法
const generateUploadId = (): string => {
  // 方法1: 使用 crypto.randomUUID (现代浏览器)
  if (typeof crypto !== 'undefined' && typeof crypto.randomUUID === 'function') {
    return crypto.randomUUID();
  }

  // 方法2: 使用 crypto.getRandomValues (较老的浏览器)
  if (typeof crypto !== 'undefined' && typeof crypto.getRandomValues === 'function') {
    const array = new Uint8Array(16);
    crypto.getRandomValues(array);
    return Array.from(array, byte => byte.toString(16).padStart(2, '0')).join('');
  }

  // 方法3: 使用 Date.now() + Math.random() 作为最后的备选方案
  return `upload-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
};

// 处理单个文件上传（带 uploadId）
const processFileUploadWithId = async (
  file: File,
  uploadId: string,
  parentId: string,
  fileName: string,
  uploadFileInfo?: UploadFileInfo
) => {
  const fileSize = file.size;

  try {
    if (fileSize <= CHUNK_SIZE) {
      // 小文件直接上传
      await uploadSmallFileWithId(file, uploadId, parentId, fileName, uploadFileInfo);
    } else {
      // 大文件分片上传
      await uploadLargeFileWithId(file, uploadId, parentId, fileName, uploadFileInfo);
    }
  } catch (error) {
    message.error(`上传文件 "${fileName}" 失败: ${error instanceof Error ? error.message : '未知错误'}`);
    throw error;
  }
};

// 上传小文件（不分片）- 带 uploadId
const uploadSmallFileWithId = async (
  file: File,
  uploadId: string,
  parentId: string,
  fileName: string,
  uploadFileInfo?: UploadFileInfo
) => {
  try {
    // 更新进度为 0%
    if (uploadFileInfo) {
      uploadFileInfo.percentage = 0;
    }

    // 上传文件分片
    const response = await uploadChunk(uploadId, 1, 0, file.size, file, channelId);
    if (response !== null) {
      messageId.value = response.telegram_msg_id;
    }

    // 更新进度为 100%
    if (uploadFileInfo) {
      uploadFileInfo.percentage = 100;
    }

    // 合并文件
    await mergeFile(
      uploadId,
      fileName, // 使用解析后的文件名
      parentId, // 使用目标文件夹ID
      file.size,
      file.type || 'application/octet-stream',
      1,
      channelId,
      messageId.value
    );

    message.success(`文件 "${fileName}" 上传成功`);
  } catch (uploadError) {
    // 上传失败时清理上传会话
    const isCleanedUp = (file as any).isCleanedUp;
    if (!isCleanedUp) {
      try {
        await cleanupUploadSession(uploadId, channelId);
        (file as any).isCleanedUp = true;
      } catch (cleanupError) {
        console.error('Failed to cleanup upload session:', cleanupError);
      }
    }
    throw uploadError;
  }
};

// 上传大文件（分片）- 带 uploadId
const uploadLargeFileWithId = async (
  file: File,
  uploadId: string,
  parentId: string,
  fileName: string,
  uploadFileInfo?: UploadFileInfo
) => {
  try {
    const totalChunks = Math.ceil(file.size / CHUNK_SIZE);

    // 串行上传分片 - 逐个上传，失败时立即停止
    for (let chunkIndex = 0; chunkIndex < totalChunks; chunkIndex++) {
      const start = chunkIndex * CHUNK_SIZE;
      const end = Math.min(start + CHUNK_SIZE, file.size);
      const chunk = file.slice(start, end);

      const response = await uploadChunk(
        uploadId,
        totalChunks,
        chunkIndex,
        chunk.size,
        new File([chunk], `${fileName}.${String(chunkIndex + 1).padStart(3, '0')}`, { type: file.type }),
        channelId
      );
      if (response !== null) {
        messageId.value = response.telegram_msg_id;
      }

      // 更新进度
      if (uploadFileInfo) {
        uploadFileInfo.percentage = Math.min(Math.floor(((chunkIndex + 1) / totalChunks) * 100), 100);
      }
    }

    // 合并文件
    await mergeFile(
      uploadId,
      fileName, // 使用解析后的文件名
      parentId, // 使用目标文件夹ID
      file.size,
      file.type || 'application/octet-stream',
      totalChunks,
      channelId,
      messageId.value
    );

    message.success(`文件 "${fileName}" 上传成功`);
  } catch (uploadError) {
    // 上传失败时清理上传会话
    const isCleanedUp = (file as any).isCleanedUp;
    if (!isCleanedUp) {
      try {
        await cleanupUploadSession(uploadId, channelId);
        (file as any).isCleanedUp = true;
      } catch (cleanupError) {
        console.error('Failed to cleanup upload session:', cleanupError);
      }
    }
    throw uploadError;
  }
};

// 处理上传完成
const handleUploadFinish = () => {
  // 上传完成后可以触发文件列表刷新等操作
  if (!uploadMode.value) {
    emit('refresh-filelist');
  }
};

// 处理关闭按钮点击（原取消按钮）
const handleClose = () => {
  emit('cancel');
};

// 清理创建文件夹的缓存
onUnmounted(() => {
  creatingFolders.clear();
});
</script>

<style scoped>
.upload-modal-content {
  display: flex;
  flex-direction: column;
  width: 100%;
  max-width: 100%;
  min-width: 400px;
}

.upload-header {
  text-align: center;
  margin-bottom: 16px;
}

.upload-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
}

.upload-content {
  flex: 1;
  margin-bottom: 16px;
  width: 100%;
}

.upload-status {
  padding: 16px;
  background-color: var(--n-color-bg-soft);
  border-radius: 8px;
  font-size: 14px;
}

.uploading-filename {
  display: block;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
  max-width: 100%;
}

.upload-progress {
  margin: 12px 0;
}

.remaining-info {
  display: flex;
  justify-content: space-between;
  margin-top: 8px;
  color: var(--n-text-color-secondary);
  font-size: 14px;
}

.remaining-info span:first-child {
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
  max-width: 45%;
}

.remaining-info span:last-child {
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
  max-width: 45%;
  text-align: right;
}

.upload-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>