<script setup lang="ts">
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const route = useRoute()

const activeMenu = () => {
  if (route.path.startsWith('/config')) return '/config'
  return '/tasks'
}

const goHome = () => router.push('/')
</script>

<template>
  <el-container style="min-height: 100vh">
    <el-aside width="220px" style="background-color: #304156">
      <div
        style="height: 60px; display: flex; align-items: center; justify-content: center; color: #fff; font-size: 16px; font-weight: bold; cursor: pointer; border-bottom: 1px solid rgba(255,255,255,0.1)"
        @click="goHome"
      >
        Web 网页变更监听
      </div>
      <el-menu
        :default-active="activeMenu()"
        background-color="#304156"
        text-color="#bfcbd9"
        active-text-color="#409EFF"
        router
        style="border-right: none"
      >
        <el-menu-item index="/tasks">
          <el-icon><List /></el-icon>
          <span>监听列表</span>
        </el-menu-item>
        <el-menu-item index="/config">
          <el-icon><Setting /></el-icon>
          <span>系统配置</span>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <el-container>
      <el-main style="background: #f0f2f5; padding: 20px">
        <router-view />
      </el-main>
      <el-footer style="text-align: center; font-size: 12px; color: #999; height: 40px">
        Web 网页变更监听与推送系统 ©2026
      </el-footer>
    </el-container>
  </el-container>
</template>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}
body {
  font-family: 'Helvetica Neue', Helvetica, 'PingFang SC', 'Microsoft YaHei', Arial, sans-serif;
}
</style>
