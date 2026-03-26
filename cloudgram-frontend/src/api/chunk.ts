import { request } from '@/utils/request'



/**
 * 查询分片下载链接
 */
export function queryChunk(telegram_file_id: string) {
  return request.post('chunk/query', { telegram_file_id })
}

/**
 * 获取文件分片列表
 */
export function getFileChunks(file_id: string) {
  return request.post('chunk/list', { file_id: file_id })
}

/**
 * 代理下载分片（解决CORS问题）
 * @param telegram_file_id 分片telegramFileID
 * @returns Blob 数据
 */
export function getChunkProxyDownload(telegram_file_id: string) {
  return request.get<Blob>(`download/chunk/${telegram_file_id}`, {}, { responseType: 'blob' })
}

