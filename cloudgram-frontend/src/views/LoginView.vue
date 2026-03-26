<template>
  <div class="login-view">
    <login-page
      :loading="loading"
      @login="handleLogin"
    />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useMessage, useLoadingBar } from 'naive-ui';
import LoginPage from '../components/LoginPage.vue';
import { login } from '@/api/auth';
import { useUserStore } from '@/store/user';
import type { LoginParams } from '@/types/login';

// 状态管理
const loading = ref(false);
const message = useMessage();
const loadingBar = useLoadingBar();
const router = useRouter();
const userStore = useUserStore();

// 处理登录
const handleLogin = async (data: LoginParams) => {
  loading.value = true;
  try {
    loadingBar.start();
    
    // 调用实际的登录 API
    const response = await login(data);
    
    // 登录成功后设置用户信息到 store
    userStore.setUsername(data.username);
    if (response.token) {
      userStore.setToken(response.token);
    }
    
    loadingBar.finish();
    message.success(`登录成功，欢迎回来！`);
    router.replace({ path: '/', query: { from: 'login' } });
  } catch (error: any) {
    loadingBar.error();
    // 错误处理已经在 api/auth.ts 中处理
    if (error.code === 401) {
      message.error('用户名或密码错误');
    } else if (error.code === 429) {
      message.error('请求过于频繁，请稍后重试');
    } else if (error.code === 500) {
      message.error('服务器错误，请稍后重试');
    } else if (error.message?.includes('网络错误') || error.message?.includes('Network')) {
      message.error('网络连接失败，请检查网络');
    } else if (error.message?.includes('请求超时') || error.message?.includes('timeout')) {
      message.error('请求超时，请稍后重试');
    } else {
      message.error('登录失败，请检查账号密码');
    }
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
.login-view {
  width: 100%;
  min-height: 100vh;
}
</style>