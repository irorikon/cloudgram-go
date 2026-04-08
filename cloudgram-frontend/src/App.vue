<template>
  <n-config-provider :theme="currentTheme" :theme-overrides="themeOverrides">
    <n-global-style />
    <div class="app-container">
      <n-message-provider>
        <n-dialog-provider>
          <n-notification-provider>
            <n-loading-bar-provider>
              <RouterView />
            </n-loading-bar-provider>
          </n-notification-provider>
        </n-dialog-provider>
      </n-message-provider>
    </div>
  </n-config-provider>
</template>
<script setup lang="ts">
import { ref, watch, onMounted, computed } from 'vue';
import { RouterView } from 'vue-router'
import {
  NMessageProvider,
  NDialogProvider,
  NNotificationProvider,
  NLoadingBarProvider,
  NConfigProvider,
  NGlobalStyle,
  darkTheme,
  useOsTheme
} from 'naive-ui';
import { useThemeStore } from '@/store/theme';

// 主题配置 - 使用Naive UI默认值，只覆盖必要的部分
const currentTheme = ref<any>(null);
const themeOverrides = {
  common: {
    primaryColor: '#1890ff',
    borderRadius: '8px',
  }
};

// 获取主题store实例
const themeStore = useThemeStore();

// 使用 NaiveUI 提供的 useOsTheme 获取系统主题
const osTheme = useOsTheme();

// 计算当前主题，优先使用用户设置，如果没有用户设置则使用系统主题
const computedIsDark = computed(() => {
  if (themeStore.userDefined) {
    return themeStore.isDark;
  }
  return osTheme.value === 'dark';
});

// 监听计算后的主题值并更新当前主题
watch(computedIsDark, (isDark) => {
  currentTheme.value = isDark ? darkTheme : null;
  // 如果需要，也可以同步到 store
  themeStore.isDark = isDark;
}, { immediate: true });

// 监听系统主题变化，如果用户没有自定义主题，则跟随系统主题变化
watch(osTheme, (newOsTheme) => {
  if (!themeStore.userDefined) {
    const isDark = newOsTheme === 'dark';
    currentTheme.value = isDark ? darkTheme : null;
    themeStore.isDark = isDark;
  }
});

// 初始化主题逻辑
const initializeTheme = () => {
  if (!themeStore.userDefined) {
    const isDark = osTheme.value === 'dark';
    if (isDark !== themeStore.isDark) {
      themeStore.isDark = isDark;
    }
  }
};

// 初始化主题 - 在应用挂载时执行
onMounted(() => {
  initializeTheme();
});
</script>

<style scoped>
.app-container {
  width: 100%;
  min-height: 100vh;
  background-color: var(--n-color-body);
}
</style>

<!-- 全局样式定义 -->
<style>
/* 统一的UI显示参数 */
:root {
  /* 图标相关 */
  --icon-size-small: 16px;
  --icon-size-medium: 20px;
  --icon-size-large: 24px;
  --icon-size-xlarge: 28px;

  /* 文字相关 */
  --text-size-small: 12px;
  --text-size-medium: 14px;
  --text-size-large: 16px;
  --text-size-xlarge: 18px;

  /* 颜色相关 - 使用Naive UI的主题变量 */
  --icon-color-primary: var(--n-color-primary);
  --icon-color-secondary: var(--n-text-color-2);
  --icon-color-tertiary: var(--n-text-color-3);

  /* 间距相关 */
  --gap-small: 4px;
  --gap-medium: 6px;
  --gap-large: 8px;
  --gap-xlarge: 12px;
}
</style>