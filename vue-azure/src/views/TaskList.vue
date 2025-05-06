<template>
  <div class="task-list">
    <div class="task-header">
      <h1>My Tasks</h1>
      <button @click="showAddTaskModal = true" class="btn btn-primary">
        Add New Task
      </button>
    </div>
    
    <div v-if="loading" class="loading-container">
      <div class="spinner"></div>
      <p>Loading tasks...</p>
    </div>
    
    <div v-else-if="error" class="error-message">
      <p>{{ error }}</p>
      <button @click="fetchTasks" class="btn btn-secondary">Try Again</button>
    </div>
    
    <div v-else-if="tasksStore.tasks.length === 0" class="empty-state">
      <div class="empty-icon" style="background-color: var(--secondary)">
        <i class="fas fa-tasks"></i>
      </div>
      <h3>No Tasks Yet</h3>
      <p>Create your first task to get started</p>
      <button @click="showAddTaskModal = true" class="btn btn-primary">Add Task</button>
    </div>
    
    <div v-else class="task-grid">
      <div v-for="task in tasksStore.tasks" :key="task.id" class="task-card">
        <div class="task-card-header">
          <h3>{{ task.title }}</h3>
          <div class="task-actions">
            <button @click="editTask(task)" class="btn-icon">
              <i class="fas fa-edit">Edit</i>
            </button>
            <button @click="confirmDeleteTask(task)" class="btn-icon">
              <i class="fas fa-trash">Delete</i>
            </button>
          </div>
        </div>
        <p class="task-description">{{ task.description }}</p>
      </div>
    </div>
    
    <!-- Add/Edit Task Modal -->
    <div v-if="showAddTaskModal || showEditTaskModal" class="modal">
      <div class="modal-content">
        <div class="modal-header">
          <h2>{{ showEditTaskModal ? 'Edit Task' : 'Add New Task' }}</h2>
          <button @click="closeModals" class="btn-close">&times;</button>
        </div>
        <div class="modal-body">
          <form @submit.prevent="submitTaskForm">
            <div class="form-group">
              <label for="title">Title</label>
              <input 
                type="text" 
                id="title" 
                v-model="taskForm.title" 
                required
                placeholder="Enter task title"
              >
            </div>
            
            <div class="form-group">
              <label for="description">Description</label>
              <textarea 
                id="description" 
                v-model="taskForm.description" 
                rows="3"
                placeholder="Enter task description"
              ></textarea>
            </div>
            
            <div class="form-actions">
              <button type="button" @click="closeModals" class="btn btn-secondary">
                Cancel
              </button>
              <button type="submit" class="btn btn-primary">
                {{ showEditTaskModal ? 'Update Task' : 'Create Task' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
    
    <!-- Delete Confirmation Modal -->
    <div v-if="showDeleteModal" class="modal">
      <div class="modal-content">
        <div class="modal-header">
          <h2>Confirm Delete</h2>
          <button @click="showDeleteModal = false" class="btn-close">&times;</button>
        </div>
        <div class="modal-body">
          <p>Are you sure you want to delete the task "{{ taskToDelete?.title }}"?</p>
          <p class="warning-text">This action cannot be undone.</p>
          
          <div class="form-actions">
            <button @click="showDeleteModal = false" class="btn btn-secondary">
              Cancel
            </button>
            <button @click="deleteTask" class="btn btn-accent">
              Delete
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, reactive, onMounted } from 'vue'
import { useTasksStore } from '../stores/tasks'

export default {
  name: 'TaskListView',
  setup() {
    const tasksStore = useTasksStore()
    const loading = ref(false)
    const error = ref(null)
    
    // Task form and modals
    const taskForm = reactive({
      title: '',
      description: ''
    })
    
    const showAddTaskModal = ref(false)
    const showEditTaskModal = ref(false)
    const showDeleteModal = ref(false)
    const taskToEdit = ref(null)
    const taskToDelete = ref(null)
    
    // Fetch tasks on component mount
    const fetchTasks = async () => {
      try {
        loading.value = true
        error.value = null
        await tasksStore.fetchTasks()
      } catch (err) {
        error.value = 'Failed to load tasks. Please try again.'
        console.error('Error fetching tasks:', err)
      } finally {
        loading.value = false
      }
    }
    
    onMounted(fetchTasks)
    
    // Reset form to default values
    const resetForm = () => {
      taskForm.title = ''
      taskForm.description = ''
      taskToEdit.value = null
    }
    
    // Close all modals and reset form
    const closeModals = () => {
      showAddTaskModal.value = false
      showEditTaskModal.value = false
      showDeleteModal.value = false
      resetForm()
    }
    
    // Edit task
    const editTask = (task) => {
      taskToEdit.value = task
      taskForm.title = task.title
      taskForm.description = task.description || ''
      showEditTaskModal.value = true
    }
    
    // Submit task form (create or update)
    const submitTaskForm = async () => {
      try {
        loading.value = true
        error.value = null
        
        const taskData = {
          title: taskForm.title,
          description: taskForm.description
        }
        
        if (showEditTaskModal.value && taskToEdit.value) {
          await tasksStore.updateTask(taskToEdit.value.id, taskData)
        } else {
          await tasksStore.createTask(taskData)
        }
        
        closeModals()
      } catch (err) {
        error.value = showEditTaskModal.value 
          ? 'Failed to update task. Please try again.' 
          : 'Failed to create task. Please try again.'
        console.error('Error submitting task:', err)
      } finally {
        loading.value = false
      }
    }
    
    // Confirm delete task
    const confirmDeleteTask = (task) => {
      taskToDelete.value = task
      showDeleteModal.value = true
    }
    
    // Delete task
    const deleteTask = async () => {
      if (!taskToDelete.value) return
      
      try {
        loading.value = true
        error.value = null
        await tasksStore.deleteTask(taskToDelete.value.id)
        showDeleteModal.value = false
        taskToDelete.value = null
      } catch (err) {
        error.value = 'Failed to delete task. Please try again.'
        console.error('Error deleting task:', err)
      } finally {
        loading.value = false
      }
    }
    
    return {
      loading,
      error,
      taskForm,
      showAddTaskModal,
      showEditTaskModal,
      showDeleteModal,
      taskToDelete,
      tasksStore,
      fetchTasks,
      editTask,
      submitTaskForm,
      closeModals,
      confirmDeleteTask,
      deleteTask
    }
  }
}
</script>

<style scoped lang="scss">
.task-list {
  position: relative;
}

.task-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
  
  h1 {
    color: var(--primary);
    margin: 0;
  }
}

.task-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 1.5rem;
}

.task-card {
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  padding: 1.5rem;
  transition: transform 0.2s, box-shadow 0.2s;
  
  &:hover {
    transform: translateY(-3px);
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
  }
}

.task-card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 1rem;
  
  h3 {
    margin: 0;
    color: var(--primary);
    font-size: 1.2rem;
  }
}

.task-actions {
  display: flex;
  gap: 0.5rem;
}

.btn-icon {
  background: none;
  border: none;
  color: var(--gray);
  cursor: pointer;
  font-size: 1rem;
  padding: 0.25rem;
  transition: color 0.2s;
  
  &:hover {
    color: var(--primary);
  }
}

.task-description {
  color: var(--dark);
  margin-bottom: 1.5rem;
  line-height: 1.5;
  overflow: hidden;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
}

.task-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.btn-text {
  background: none;
  border: none;
  color: var(--primary);
  cursor: pointer;
  font-size: 0.9rem;
  padding: 0;
  text-decoration: none;
  
  &:hover {
    text-decoration: underline;
  }
}

.empty-state {
  text-align: center;
  padding: 3rem 1rem;
  
  h3 {
    color: var(--primary);
    margin: 1rem 0 0.5rem;
  }
  
  p {
    color: var(--gray);
    margin-bottom: 1.5rem;
  }
}

.empty-icon {
  width: 70px;
  height: 70px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto;
  color: white;
  font-size: 1.8rem;
}

.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 3rem 0;
  
  .spinner {
    width: 40px;
    height: 40px;
    border: 4px solid rgba(0, 0, 0, 0.1);
    border-radius: 50%;
    border-top-color: var(--primary);
    animation: spin 1s ease-in-out infinite;
    margin-bottom: 1rem;
  }
  
  @keyframes spin {
    to { transform: rotate(360deg); }
  }
  
  p {
    color: var(--gray);
  }
}

.error-message {
  background-color: rgba(207, 8, 34, 0.1);
  border-left: 4px solid var(--accent);
  padding: 1.5rem;
  margin: 1rem 0;
  
  p {
    color: var(--accent);
    margin-bottom: 1rem;
  }
}

.modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background-color: white;
  border-radius: 8px;
  width: 100%;
  max-width: 500px;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem;
  border-bottom: 1px solid #eee;
  
  h2 {
    margin: 0;
    color: var(--primary);
    font-size: 1.5rem;
  }
}

.btn-close {
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  color: var(--gray);
  
  &:hover {
    color: var(--accent);
  }
}

.modal-body {
  padding: 1.5rem;
}

.form-group {
  margin-bottom: 1.5rem;
  
  label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: 500;
    color: var(--dark);
  }
  
  input, textarea, select {
    width: 100%;
    padding: 0.75rem;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-size: 1rem;
    
    &:focus {
      outline: none;
      border-color: var(--primary);
      box-shadow: 0 0 0 2px rgba(1, 53, 170, 0.2);
    }
  }
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  margin-top: 2rem;
}

.warning-text {
  color: var(--accent);
  font-style: italic;
}
</style>