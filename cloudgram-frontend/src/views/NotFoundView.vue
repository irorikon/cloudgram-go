<template>
  <div class="not-found-container">
    <n-result status="404" title="404 资源不存在" description="生活总归带点荒谬">
      <template #footer>
        <n-button @click="goHome">返回首页</n-button>
      </template>
    </n-result>
  </div>
</template>

<script lang="ts">
import { NResult, NButton } from 'naive-ui'
import { useRouter } from 'vue-router'

export default {
  name: 'NotFoundView',
  components: {
    NResult,
    NButton
  },
  setup() {
    const router = useRouter()
    
    const goHome = () => {
      router.push('/home')
    }
    
    return {
      goHome
    }
  }
}
</script>

<style scoped>
.not-found-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  padding: 20px;
  box-sizing: border-box;
}
</style>