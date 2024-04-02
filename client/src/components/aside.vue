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
    .reduce((pre: string[], curr: string) => {
      if (curr.startsWith("Simp")) {
        pre.unshift(curr);
      } else {
        pre.push(curr);
      }
      return pre;
    }, []);
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
    >
      <el-menu-item class="app-text-center">
        <el-icon class="app-not-show"><Search /></el-icon>
        <el-input v-model="keyword"></el-input>
      </el-menu-item>
      <el-menu-item
        v-for="(item, index) in serverList"
        class="app-text-center"
        :index="item"
        :key="index"
        @click="handleOpen(item)"
      >
        <el-icon class="app-not-show"><Menu /></el-icon>
        <template #title>{{ item }}</template>
      </el-menu-item>
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
