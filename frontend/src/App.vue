<template>
  <router-view />
</template>

<script setup>
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html, body {
  margin: 0;
  padding: 0;
  width: 100%;
  min-height: 100%;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, 'Microsoft YaHei', sans-serif;
  font-size: 14px;
  line-height: 1.6;
  color: #333;
}

#app {
  width: 100%;
  min-height: 100vh;
  margin: 0;
  padding: 0;
}

/* 确保Element Plus组件在web端显示正常 */
.el-form-item {
  margin-bottom: 22px;
}

.el-input, .el-select {
  width: 100%;
}

.el-button {
  font-size: 14px;
  padding: 12px 20px;
}
</style>
