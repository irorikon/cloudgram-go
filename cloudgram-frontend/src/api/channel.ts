// src/api/auth.ts
import { request } from '@/utils/request'

export function listChannels() {
    return request.get('channel/list')
}

export function getOneChannel() {
    return request.get('channel/get')
}

export function createChannel(channelId: number, channelName: string) {
    return request.post('channel/create', { channel_id: channelId, channel_name: channelName })
}

export function deleteChannel(channelId: number) {
    return request.post('channel/delete', { channel_id: channelId })
}

export function updateChannel(channelId: number, channelName: string) {
    return request.post('channel/update', { channel_id: channelId, channel_name: channelName })
}

export function checkChannel(channelId: number) {
    return request.post('channel/check', { channel_id: channelId })
}