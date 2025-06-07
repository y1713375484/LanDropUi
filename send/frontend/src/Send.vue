<script setup>
import { ref,computed } from 'vue';
import {ChooseFile,Send,Msgalert} from "../wailsjs/go/main/App.js";
import {EventsOn} from "../wailsjs/runtime/runtime.js";


const filePathList = ref({});
const isLoading = ref(false); // 新增：加载状态变量
// 新增：计算属性，用于获取 filePathList 的键值长度
const fileCount = computed(() => Object.keys(filePathList.value).length);
//文件传输百分比
EventsOn('percent', (data) => {
  console.log(data)
  const fileUUid=data.fileUUID
  filePathList.value[fileUUid]["percent"]=data.percent
});



function handleFileChange(){
  ChooseFile(filePathList.value).then(res=>{
    filePathList.value=res
  })
}

async function sendFile(){
  if (filePathList.value.length === 0) {
    Msgalert("请选择要发送的文件！")
    return;
  }
  isLoading.value = true; // 设置加载状态为 true
  try {
    const res = await Send(filePathList.value); // 使用 await 等待方法执行完毕
    console.log(res)
    Msgalert(res)
  } catch (error) {

  } finally {
    isLoading.value = false; // 无论成功或失败，都设置加载状态为 false
  }
}

</script>
<template>
  <div class="size-revise-container">

    <div class="section">
      <div class="input-group">
        <button class="btn primary-btn" @click="handleFileChange">选择要发送的文件</button>
      </div>
      <div class="file-list-card">
        <h3>已选择文件 ({{ fileCount }})</h3>
        <ul class="file-list">
          <li v-if="fileCount === 0" class="placeholder-item">
            请选择要发送的文件
          </li>
          <li v-for="(item, index) in filePathList" :key="index" class="file-list-item">
            <span class="file-path-text">{{item.filePath}}</span>
            <van-progress :percentage="item.percent" stroke-width="8" />
          </li>
        </ul>
      </div>
    </div>
    <div class="action-section">
      <button class="btn process-btn" @click="sendFile" :disabled="isLoading">
        {{ isLoading ? '发送中...' : '发送文件' }} <!-- 根据加载状态显示不同文本 -->
      </button>
    </div>

  </div>
</template>

<style scoped>
.size-revise-container {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
  padding: 2rem;
  background-color: #f9fbfd;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
  max-width: 800px;
  margin: 2rem auto;
}

.header {
  text-align: center;
  margin-bottom: 1rem;
}

.header h2 {
  font-size: 2rem;
  color: #333;
  margin-bottom: 0.5rem;
}

.tip {
  color: #e67e22;
  font-weight: bold;
}

.section {
  background-color: #fff;
  padding: 1.5rem;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.03);
}

.input-group {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-bottom: 1rem;
}

.dimension-input-group {
  gap: 0.8rem;
}

.dimension-input {
  padding: 0.6rem 1rem;
  border: 1px solid #ccc;
  border-radius: 5px;
  font-size: 1rem;
  width: 100%;
  max-width: 250px;
  box-sizing: border-box;
}

.dimension-input:focus {
  outline: none;
  border-color: #3498db;
  box-shadow: 0 0 0 3px rgba(52, 152, 219, 0.2);
}

.btn {
  padding: 0.8rem 1.5rem;
  border: none;
  border-radius: 5px;
  cursor: pointer;
  font-size: 1rem;
  transition: background-color 0.3s ease, transform 0.2s ease;
}

.btn:hover {
  transform: translateY(-2px);
}

.primary-btn {
  background-color: #3498db;
  color: white;
}

.primary-btn:hover {
  background-color: #2980b9;
}

.output-path {
  color: #c0392b;
  font-weight: bold;
  flex-grow: 1;
  word-break: break-all;
}

.file-list-card {
  margin-top: 1rem;
  border: 1px solid #eee;
  border-radius: 5px;
  padding: 1rem;
  max-height: 250px;
  overflow-y: auto;
}

.file-list-card h3 {
  font-size: 1.2rem;
  color: #555;
  margin-bottom: 0.8rem;
  border-bottom: 1px solid #eee;
  padding-bottom: 0.5rem;
}

.file-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.file-list-item {
  background-color: #f0f4f7;
  padding: 0.6rem 1rem;
  margin-bottom: 0.5rem;
  border-radius: 4px;
  color: #444;
  font-size: 0.9rem;
  word-break: break-all;
  display: flex; /* 使用 Flexbox */
  flex-direction: column; /* 垂直排列 */
  align-items: flex-start; /* 左对齐 */
}

.file-path-text {
  font-size: 1rem; /* 增大文件路径字体大小 */
  color: #333; /* 更改文件路径颜色 */
  margin-bottom: 0.5rem; /* 增加文件路径和进度条之间的间距 */
}

.van-progress {
  width: 100%; /* 进度条宽度占满 */
  margin-top: 0.5rem; /* 增加进度条顶部外边距 */
}

.placeholder-item {
  color: #888;
  text-align: center;
  padding: 1rem;
}

.action-section {
  text-align: center;
  margin-top: 1.5rem;
}

.process-btn {
  background-color: #27ae60;
  color: white;
  padding: 1rem 2.5rem;
  font-size: 1.2rem;
  font-weight: bold;
}

.process-btn:hover {
  background-color: #219653;
}

.process-btn:disabled { /* 新增：禁用状态样式 */
  background-color: #cccccc;
  cursor: not-allowed;
}
</style>