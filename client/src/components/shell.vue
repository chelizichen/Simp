<script lang="ts">
export default {
  name: "shell-componetn",
};
</script>
<template>
  <el-dialog
    v-model="props.shellVisible"
    title="SimpShell"
    width="80%"
    @close="emits('closeDialog')"
  >
    <div style="display: flex">
      <div
        class="resu"
        style="
          background-color: black;
          height: 500px;
          padding: 5px 10px;
          overflow: scroll;
          width: 80%;
          border-right: 1px solid white;
        "
      >
        <div style="color: aliceblue">// SimpShell :: created By leeks</div>
        <div
          style="color: aliceblue; margin: 2px; white-space: pre"
          v-for="item in results"
          :key="item"
        >
          <span v-html="item"></span>
        </div>
      </div>
      <div
        class="resu"
        style="
          background-color: black;
          height: 500px;
          padding: 5px 10px;
          overflow: scroll;
          width: 20%;
        "
      >
        <div style="color: aliceblue">// Input Area</div>
        <div style="color: aliceblue; margin: 2px; white-space: pre">
          <el-input
            type="textarea"
            v-model="command"
            class="inp"
            rows="50"
            @keyup.enter="InputCommand"
          ></el-input>
        </div>
      </div>
    </div>
  </el-dialog>
</template>

<script setup lang="ts">
import { CirclePlus, Delete } from "@element-plus/icons-vue";
import { ElMessage, ElMessageBox } from "element-plus";
import { reactive, ref, watch, computed } from "vue";
import { NewEventSource } from "@/utils/shell";
import { useShellStore } from "@/stores/counter";

const props = defineProps<{
  shellVisible: boolean;
}>();
const emits = defineEmits(["closeDialog", "showReleaseDialog"]);
const store = useShellStore();
const results = computed(() => {
  return store.outputStack;
});
const command = ref();
function InputCommand() {
  store.pushStack(
    `<span style="color:red;font-size:24px;font-weight:700;font-family: serif;">${command.value}</span>`
  );
  NewEventSource(command.value);
  command.value = "";
}
</script>

<style scoped>
.inp >>> .el-textarea__inner {
  background-color: black;
  border: none;
  box-shadow: none;
  color: white;
}
</style>
