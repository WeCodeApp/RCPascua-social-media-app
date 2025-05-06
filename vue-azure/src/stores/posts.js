import { defineStore } from 'pinia'
import axios from 'axios'

export const usePostStore = defineStore('posts', {
  state: () => ({
    posts: [],
    currentPost: {},
    loading: false,
    error: null
  }),

  getters: {
    getPosts: (state) => state.posts,
    getCurrentPost: (state) => state.currentPost,
    getLoading: (state) => state.loading,
    getError: (state) => state.error
  },

  actions: {
    async fetchPosts() {
      try {
        this.loading = true
        this.error = null

        const response = await axios.get(`${import.meta.env.VITE_API_URL}/posts`, {
          headers: {
            Authorization: `Bearer ${localStorage.getItem('access_token')}`
          }
        })

        // const response2 = await axios.get(`${import.meta.env.VITE_API_URL}/posts`, {
        //   headers: {
        //     Authorization: `Bearer ${localStorage.getItem('access_token')}`
        //   }
        // })

        // console.log(response2.data.posts)
        // console.log(response.data)
        console.log(response.data.posts)
        // console.log(response.data.current_page)
        // console.log(response.data.posts[0].post_text)
        if (Array.isArray(response.data.posts) && response.data.posts.length > 0) {
          this.posts = response.data.posts
        }

        return this.posts
      } catch (error) {
        this.error = error.response?.data?.message || 'Failed to fetch tasks'
        throw error
      } finally {
        this.loading = false
      }
    },

    async fetchPostsById(id) {
      try {
        
        this.loading = true
        this.error = null

        const response = await axios.get(`${import.meta.env.VITE_API_URL}/posts/${id}`, {
          headers: {
            Authorization: `Bearer ${localStorage.getItem('access_token')}`
          }
        })

        this.currentPosts = response.data.post
        return this.currentPost
      } catch (error) {
        this.error = error.response?.data?.message || 'Failed to fetch task'
        throw error
      } finally {
        this.loading = false
      }
    },

    async createPost(postsData) {
      try {
        this.loading = true
        this.error = null

        const response = await axios.post(`${import.meta.env.VITE_API_URL}/posts`, postsData, {
          headers: {
            Authorization: `Bearer ${localStorage.getItem('access_token')}`
          }
        })

        // Add the new task to the tasks array
        this.posts.push(response.data.post)

        return response.data.posts
      } catch (error) {
        this.error = error.response?.data?.message || 'Failed to create task'
        throw error
      } finally {
        this.loading = false
      }
    },

    async updatePost(id, postData) {
      try {

      
        this.loading = true
        this.error = null

        const response = await axios.put(`${import.meta.env.VITE_API_URL}/posts/${id}`, postData, {
          headers: {
            Authorization: `Bearer ${localStorage.getItem('access_token')}`
          }
        })

        // Update the task in the tasks array
        const index = this.posts.findIndex(post => post.post_id === id)
        if (index !== -1) {
          this.posts[index] = response.data.post
        }

        // Update currentTask if it's the same task
        if (this.currentPost && this.currentPost.post_id === id) {
          this.currentPost = response.data.post
        }

        return response.data.posts
      } catch (error) {
        this.error = error.response?.data?.message || 'Failed to update task'
        throw error
      } finally {
        this.loading = false
      }
    },

    async deletePost(id) {
      try {
        this.loading = true
        this.error = null

        await axios.delete(`${import.meta.env.VITE_API_URL}/posts/${id}`, {
          headers: {
            Authorization: `Bearer ${localStorage.getItem('access_token')}`
          }
        })

        // Remove the task from the tasks array
        this.posts = this.posts.filter(post => post.post_id !== id)

        // Clear currentTask if it's the same post
        if (this.currentPost && this.currentPost.id === id) {
          this.currentPost = null
        }

        return true
      } catch (error) {
        this.error = error.response?.data?.message || 'Failed to delete task'
        throw error
      } finally {
        this.loading = false
      }
    }
  }
})
