<!--
 * 项目名称：blog-frontend
 * 文件名称：DiffViewer.vue
 * 创建时间：2026-04-29 15:00:00
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：文本差异对比展示组件，统一使用行内高亮模式
 -->
<template>
  <div class="diff-viewer">
    <div class="inline-diff">
      <div class="diff-content">
        <span v-for="(part, index) in diffParts" :key="index" :class="getDiffClass(part.type)">{{
          part.text
        }}</span>
      </div>
      <div class="diff-legend">
        <span class="legend-item">
          <span class="legend-color delete"></span>
          <span class="legend-text">删除</span>
        </span>
        <span class="legend-item">
          <span class="legend-color insert"></span>
          <span class="legend-text">新增</span>
        </span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { computeDiff, getDiffWithContext, type DiffResult } from '@/utils/diff'

interface Props {
  oldValue: string
  newValue: string
  field?: string
}

const props = defineProps<Props>()

// 计算差异
const diffParts = computed<DiffResult[]>(() => {
  const diffs = computeDiff(props.oldValue, props.newValue)
  return getDiffWithContext(diffs, 100)
})

// 获取差异样式类
function getDiffClass(type: string): string {
  switch (type) {
    case 'delete':
      return 'diff-delete'
    case 'insert':
      return 'diff-insert'
    default:
      return 'diff-equal'
  }
}
</script>

<style scoped lang="scss">
.diff-viewer {
  width: 100%;

  .inline-diff {
    .diff-content {
      padding: 16px;
      background: var(--card-color);
      border-radius: 4px;
      border: 1px solid var(--border-color);
      line-height: 1.8;
      word-break: break-word;
      white-space: pre-wrap;
      font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
      font-size: 14px;

      .diff-equal {
        color: var(--text-color-1);
      }

      .diff-delete {
        background-color: #ffe6e6;
        color: #d32f2f;
        text-decoration: line-through;
        padding: 2px 4px;
        border-radius: 2px;
      }

      .diff-insert {
        background-color: #e6ffe6;
        color: #2e7d32;
        padding: 2px 4px;
        border-radius: 2px;
      }
    }

    .diff-legend {
      display: flex;
      gap: 16px;
      margin-top: 12px;
      padding: 8px 12px;
      background: var(--card-color);
      border-radius: 4px;
      font-size: 12px;

      .legend-item {
        display: flex;
        align-items: center;
        gap: 6px;

        .legend-color {
          width: 16px;
          height: 16px;
          border-radius: 2px;

          &.delete {
            background-color: #ffe6e6;
            border: 1px solid #d32f2f;
          }

          &.insert {
            background-color: #e6ffe6;
            border: 1px solid #2e7d32;
          }
        }

        .legend-text {
          color: var(--text-color-2);
        }
      }
    }
  }
}
</style>
