<template>
  <div class="task-list">
    <div class="task-header">
      <h1>My Posts</h1>
      <button @click="showAddPostModal = true" class="btn btn-primary">
        Add New Post
      </button>
    </div>
    
    <div v-if="loading" class="loading-container">
      <div class="spinner"></div>
      <p>Loading tasks...</p>
    </div>
    
    <div v-else-if="error" class="error-message">
      <p>{{ error }}</p>
      <button @click="fetchPosts" class="btn btn-secondary">Try Again</button>
    </div>
    
    <div v-else-if="postsStore.posts.length === 0" class="empty-state">
      <div class="empty-icon" style="background-color: var(--secondary)">
        <i class="fas fa-tasks"></i>
      </div>
      <h3>No Tasks Yet</h3>
      <p>Create your first task to get started</p>
      <button @click="showAddPostModal = true" class="btn btn-primary">Add Post</button>
    </div>
    
    <div v-else class="task-grid">
      <!-- <p class="userInfo" v-if="authStore.user?.name"></p> -->
      <div v-for="post in postsStore.posts" :key="post.id" class="task-card">
      <div class="task-card-header">
        <h3>{{ post.post_text }}</h3>
        <div class="task-actions">
        <button @click="editPost(post)" class="btn-icon">
          <i class="fas fa-edit">Edit</i>
        </button>
        <button @click="confirmdeletePost(post)" class="btn-icon">
          <i class="fas fa-trash">Delete</i>
        </button>
        </div>
      </div>
      <p class="task-description">{{ post.post_image }}</p>
      <p v-if="userId == post.user_id"> {{ userName }}</p>
      </div>
    </div>
    
    <!-- Add/Edit Task Modal -->
    <div v-if="showAddPostModal || showeditPostModal" class="modal">
      <div class="modal-content">
        <div class="modal-header">
          <h2>{{ showeditPostModal ? 'Edit Post' : 'Add New Post' }}</h2>
          <button @click="closeModals" class="btn-close">&times;</button>
        </div>
        <div class="modal-body">
          <form @submit.prevent="submitPostForm">
            <div class="form-group">
              <label for="title">Post Text</label>
              <input 
                type="text" 
                id="Post Text" 
                v-model="postForm.post_text" 
                required
                placeholder="Enter post text"
              >
            </div>
            
            <div class="form-group">
              <label for="description">Post Image Link</label>
              <textarea 
                id="Post Image Link" 
                v-model="postForm.post_image" 
                rows="3"
                placeholder="Enter post image link"
              ></textarea>
            </div>
            
            <div class="form-actions">
              <button type="button" @click="closeModals" class="btn btn-secondary">
                Cancel
              </button>
              <button type="submit" class="btn btn-primary">
                {{ showeditPostModal ? 'Update Post' : 'Create Post' }}
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
          <p>Are you sure you want to delete the post "{{ postToDelete?.post_text }}"?</p>
          <p class="warning-text">This action cannot be undone.</p>
          
          <div class="form-actions">
            <button @click="showDeleteModal = false" class="btn btn-secondary">
              Cancel
            </button>
            <button @click="deletePost" class="btn btn-accent">
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
import { usePostStore } from '../stores/posts'

export default {
  name: 'PostListView',
  setup() {
    const postsStore = usePostStore()
    const loading = ref(false)
    const error = ref(null)
    
    const userId = localStorage.getItem('user_id')
    const userName = localStorage.getItem('user_name')
     
    console.log('userId:', userId)
    console.log('userName:', userName)

    // Task form and modals
    const postForm = reactive({
      post_text: '',
      post_image: ''
    })
    
    const showAddPostModal = ref(false)
    const showeditPostModal = ref(false)
    const showDeleteModal = ref(false)
    const postToEdit = ref(null)
    const postToDelete = ref(null)
    
    // Fetch tasks on component mount
    const fetchPosts = async () => {
      try {
        loading.value = true
        error.value = null
        return await postsStore.fetchPosts()
        // const data = await postsStore.fetchPosts()
        // console.log('Fetched posts:', data)

      } catch (err) {
        error.value = 'Failed to load tasks. Please try again.'
        console.error('Error fetching tasks:', err)
      } finally {
        loading.value = false
      }
    }
    
    onMounted(fetchPosts)
    
    // Reset form to default values
    const resetForm = () => {
      postForm.post_text = ''
      postForm.post_image = ''
      postToEdit.value = null
    }
    
    // Close all modals and reset form
    const closeModals = () => {
      showAddPostModal.value = false
      showeditPostModal.value = false
      showDeleteModal.value = false
      resetForm()
    }
    
    // Edit post
    const editPost = (post) => {
      postToEdit.value = post
      postForm.post_text = post.post_text
      postForm.post_image = post.post_image
      showeditPostModal.value = true
    }
    
    // Submit post form (create or update)
    const submitPostForm = async () => {
      try {
        loading.value = true
        error.value = null
        
        const postData = {
          post_text: postForm.post_text,
          post_image: postForm.post_image
        }
        console.log('Post data:', postToEdit)
        if (showeditPostModal.value && postToEdit.value) {
          await postsStore.updatePost(postToEdit.value.post_id, postData)
        } else {
          await postsStore.createPost(postData)
        }
        
        closeModals()
      } catch (err) {
        error.value = showeditPostModal.value 
          ? 'Failed to update post. Please try again.' 
          : 'Failed to create post. Please try again.'
        console.error('Error submitting post:', err)
      } finally {
        loading.value = false
      }
    }
    
    // Confirm delete task
    const confirmdeletePost = (post) => {
      postToDelete.value = post
      showDeleteModal.value = true
    }
    
    // Delete task
    const deletePost = async () => {
      if (!postToDelete.value) return
      
      try {
        loading.value = true
        error.value = null
        await postsStore.deletePost(postToDelete.value.post_id)
        showDeleteModal.value = false
        postToDelete.value = null
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
      userId,
      userName,
      postForm,
      showAddPostModal,
      showeditPostModal,
      showDeleteModal,
      postToDelete,
      postsStore,
      fetchPosts,
      editPost,
      submitPostForm,
      closeModals,
      confirmdeletePost,
      deletePost
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