import { defineStore } from 'pinia'
import axios from 'axios'

export const useTasksStore = defineStore('tasks', {
  state: () => ({
    tasks: [],
    currentTask: {},
    loading: false,
    error: null
  }),

  getters: {
    getTasks: (state) => state.tasks,
    getCurrentTask: (state) => state.currentTask,
    getLoading: (state) => state.loading,
    getError: (state) => state.error
  },

  actions: {
    async fetchTasks() {
      try {
        this.loading = true
        this.error = null

        const response = await axios.get(`${import.meta.env.VITE_API_URL}/tasks`, {
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
        // console.log(response.data.tasks)
        if (Array.isArray(response.data.tasks) && response.data.tasks.length > 0) {
          this.tasks = response.data.tasks
        }

        return this.tasks
      } catch (error) {
        this.error = error.response?.data?.message || 'Failed to fetch tasks'
        throw error
      } finally {
        this.loading = false
      }
    },

    async fetchTaskById(id) {
      try {
        this.loading = true
        this.error = null

        const response = await axios.get(`${import.meta.env.VITE_API_URL}/tasks/${id}`, {
          headers: {
            Authorization: `Bearer ${localStorage.getItem('access_token')}`
          }
        })

        this.currentTask = response.data.task
        return this.currentTask
      } catch (error) {
        this.error = error.response?.data?.message || 'Failed to fetch task'
        throw error
      } finally {
        this.loading = false
      }
    },

    async createTask(taskData) {
      try {
        this.loading = true
        this.error = null

        const response = await axios.post(`${import.meta.env.VITE_API_URL}/tasks`, taskData, {
          headers: {
            Authorization: `Bearer ${localStorage.getItem('access_token')}`
          }
        })

        // Add the new task to the tasks array
        this.tasks.push(response.data.task)

        return response.data.tasks
      } catch (error) {
        this.error = error.response?.data?.message || 'Failed to create task'
        throw error
      } finally {
        this.loading = false
      }
    },

    async updateTask(id, taskData) {
      try {
        this.loading = true
        this.error = null

        const response = await axios.put(`${import.meta.env.VITE_API_URL}/tasks/${id}`, taskData, {
          headers: {
            Authorization: `Bearer ${localStorage.getItem('access_token')}`
          }
        })

        // Update the task in the tasks array
        const index = this.tasks.findIndex(task => task.id === id)
        if (index !== -1) {
          this.tasks[index] = response.data.task
        }

        // Update currentTask if it's the same task
        if (this.currentTask && this.currentTask.id === id) {
          this.currentTask = response.data.task
        }

        return response.data.task
      } catch (error) {
        this.error = error.response?.data?.message || 'Failed to update task'
        throw error
      } finally {
        this.loading = false
      }
    },

    async deleteTask(id) {
      try {
        this.loading = true
        this.error = null

        await axios.delete(`${import.meta.env.VITE_API_URL}/tasks/${id}`, {
          headers: {
            Authorization: `Bearer ${localStorage.getItem('access_token')}`
          }
        })

        // Remove the task from the tasks array
        this.tasks = this.tasks.filter(task => task.id !== id)

        // Clear currentTask if it's the same task
        if (this.currentTask && this.currentTask.id === id) {
          this.currentTask = null
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
