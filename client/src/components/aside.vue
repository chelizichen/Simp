<script lang="ts">
export default {
  name: "aside-component",
};
</script>
<script lang="ts" setup>
const props = defineProps<{
  serverList: any[];
}>();
const emits = defineEmits(["handleOpen"]);
function handleOpen(item: string) {
  emits("handleOpen", item);
}
</script>

<template>
  <div>
    <div class="app-bigger-size title">
      <el-icon style="color: rgb(207, 90, 124); font-size: 36px"><Help /></el-icon>
      Simp
    </div>
    <el-menu
      class="el-menu-vertical-demo"
      active-text-color="rgb(207, 15, 124)"
      style="border: none"
    >
      <el-menu-item
        v-for="(item, index) in props.serverList"
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
}
</style>
