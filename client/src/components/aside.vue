<script lang="ts">
export default {
  name: "aside-component",
};
</script>
<script lang="ts" setup>
import { Search } from "@element-plus/icons-vue";
import { computed, ref } from "vue";
const props = defineProps<{
  serverList: any[];
}>();
const emits = defineEmits(["handleOpen"]);
function handleOpen(item: string) {
  emits("handleOpen", item);
}
const keyword = ref("");
const serverList = computed(() => {
  return props.serverList
    .filter((v) => v.match(keyword.value))
    .reduce(
      (pre: [string[], string[]], curr: string) => {
        if (curr.startsWith("Simp")) {
          pre[0].push(curr);
        } else {
          pre[1].push(curr);
        }
        return pre;
      },
      [[], []]
    );
});
function toGit() {
  window.open("https://github.com/chelizichen/Simp");
}
</script>

<template>
  <div>
    <div class="app-bigger-size title" @click="toGit()">
      <el-icon style="color: rgb(207, 90, 124); font-size: 36px"><Help /></el-icon>
      Simp
    </div>
    <el-menu
      class="el-menu-vertical-demo"
      active-text-color="rgb(207, 15, 124)"
      style="border: none"
      :default-openeds="['2']"
    >
      <el-menu-item class="app-text-center">
        <el-icon class="app-not-show"><Search /></el-icon>
        <el-input v-model="keyword"></el-input>
      </el-menu-item>
      <el-sub-menu index="1">
        <template #title>
          <el-icon><Setting /></el-icon>
          <span>Controller</span>
        </template>
        <el-menu-item
          v-for="(item, index) in serverList[0]"
          class="app-text-center"
          :index="item"
          :key="index"
          @click="handleOpen(item)"
        >
          <el-icon class="app-not-show">
            <TrendCharts />
          </el-icon>
          <template #title>{{ item }}</template>
        </el-menu-item>
      </el-sub-menu>
      <el-sub-menu index="2">
        <template #title>
          <el-icon> <Operation /></el-icon>
          <span>SubServer</span>
        </template>
        <el-menu-item
          v-for="(item, index) in serverList[1]"
          class="app-text-center"
          :index="item"
          :key="index"
          @click="handleOpen(item)"
        >
          <el-icon class="app-not-show">
            <Menu />
          </el-icon>
          <template #title>{{ item }}</template>
        </el-menu-item>
      </el-sub-menu>
    </el-menu>
  </div>
</template>

<style>
.title {
  color: rgb(207, 15, 124);
  text-align: center;
  display: flex;
  align-items: center;
  font-size: 30px;
  width: 200px;
  justify-content: center;
  cursor: pointer;
}
</style>
