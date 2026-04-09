import { request } from '@/utils/request'



/**
 * 根据文件ID获取文件详情
 */
export function getFileById(file_id: string) {
  return request.get(`file/detail/${file_id}`)
}

/**
 * 获取指定父目录下的文件列表
 */
export function getFilesByParentId(parent_id: string | null = null) {
  return request.post('file/list', { parent_id })
}

/**
 * 获取指定父目录下的文件夹列表
 */
export function getFoldersByParentId(parent_id: string | null = null) {
  return request.post('file/dir', { parent_id })
}

/**
 * 检查文件或目录是否已存在
 */
export function exists(name: string, parentId: string | null = null) {
  return request.post('file/exists', { file_name: name, parent_id: parentId })
}

/**
 * 创建文件夹
 */
export function createFile(file_name: string, parent_id: string | null = null) {
  return request.post('file/create', {
    file_name,
    parent_id,
    isDir: true
  })
}

/**
 * 重命名文件
 */
export function renameFile(file_id: string, newName: string) {
  return request.post('file/update', {
    active: 'rename',
    id: file_id,
    new_name: newName
  })
}

/**
 * 移动文件
 */
export function moveFile(fileIds: string[], newParentId: string) {
  return request.post('file/update', {
    active: 'move',
    ids: fileIds,
    parent_id: newParentId
  })
}


/**
 * 删除文件
 */
export function deleteFile(file_id: string, delete_with_tg: boolean = false, recursive: boolean = false) {
  return request.post('file/delete', {
    id: file_id,
    recursive,
    delete_with_tg
  })
}
