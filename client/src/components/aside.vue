<script lang="ts">
export default {
  name: 'aside-component'
}
</script>
<script lang="ts" setup>
const props = defineProps<{
  serverList: any[]
}>()
const emits = defineEmits(['handleOpen'])
function handleOpen(item: string) {
  emits('handleOpen', item)
}
</script>

<template>
  <div>
    <h1
      style="color: rgb(207, 15, 124); text-align: center; font-family: fantasy"
      class="app-bigger-size"
    >
      <el-icon style="color: rgb(207, 90, 124); font-size: 36px"><Help /></el-icon>
      Simp
    </h1>
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
