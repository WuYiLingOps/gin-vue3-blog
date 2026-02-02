/*
 * @ProjectName: go-vue3-blog
 * @FileName: blog.ts
 * @CreateTime: 2026-02-02 11:54:38
 * @SystemUser: Administrator
 * @Author: 無以菱
 * @Contact: huangjing510@126.com
 * @Description: 博客状态管理，管理分类、标签、文章等博客相关状态
 */

// 博客状态管理

import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Post, Category, Tag } from '@/types/blog'
import { getCategories, getTags } from '@/api'

export const useBlogStore = defineStore(
  'blog',
  () => {
    // 状态
    const categories = ref<Category[]>([])
    const tags = ref<Tag[]>([])
    const currentPost = ref<Post | null>(null)
    const loading = ref(false)

    // 获取分类列表
    async function fetchCategories() {
      try {
        loading.value = true
        const res = await getCategories()
        if (res.data) {
          categories.value = res.data
        }
      } finally {
        loading.value = false
      }
    }

    // 获取标签列表
    async function fetchTags() {
      try {
        loading.value = true
        const res = await getTags()
        if (res.data) {
          tags.value = res.data
        }
      } finally {
        loading.value = false
      }
    }

    // 设置当前文章
    function setCurrentPost(post: Post | null) {
      currentPost.value = post
    }

    // 初始化
    async function init() {
      await Promise.all([fetchCategories(), fetchTags()])
    }

    return {
      categories,
      tags,
      currentPost,
      loading,
      fetchCategories,
      fetchTags,
      setCurrentPost,
      init
    }
  },
  {
    // 配置持久化
    persist: {
      key: 'blog-data',
      storage: localStorage,
      pick: ['categories', 'tags']
    }
  }
)

