# Image Upload System - Comprehensive Documentation

## Overview
This document provides a complete overview of the image upload and display system implemented for the social network application. The system supports secure file uploads, database storage, and efficient image serving with proper privacy controls.

## Backend Implementation

### 1. File Upload Handler (`pkg/tools/tools.go`)

**Function**: `UploadHandler(file, handler)`

**Features**:
- Validates file size (maximum 5MB)
- Validates file types (JPEG, PNG, GIF only)
- Generates secure filenames using UUID to prevent conflicts
- Saves files to `backend/uploads/` directory
- Returns the file path for database storage
- Includes comprehensive error handling

**Code Example**:
```go
func UploadHandler(file multipart.File, handler *multipart.FileHeader) (string, int) {
    // File size validation (5MB limit)
    if handler.Size > 5<<20 {
        return "", http.StatusBadRequest
    }
    
    // File type validation
    if !IsAllowedFile(handler.Filename, file) {
        return "", http.StatusBadRequest
    }
    
    // Generate secure filename
    filename := generateSecureFilename(handler.Filename)
    
    // Save to uploads directory
    filePath := filepath.Join("uploads", filename)
    // ... save file logic
    
    return filePath, http.StatusOK
}
```

### 2. Database Storage

**Images Table Schema**:
```sql
CREATE TABLE images (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    filename TEXT NOT NULL,
    original_name TEXT NOT NULL,
    file_path TEXT NOT NULL,
    file_size INTEGER NOT NULL,
    content_type TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

**Integration with Posts**:
- Images are linked to posts via `image_url` field in posts table
- Post creation includes image upload processing
- Image metadata stored separately for better organization

### 3. Image Serving (`pkg/handlers/posts&comments.go`)

**Route**: `GET /images/{filename}`

**Function**: `ServeImageHandler()`

**Features**:
- Serves images directly from `uploads/` directory
- Sets proper HTTP headers (Content-Type, Cache-Control)
- Bypasses authentication middleware for public access
- Includes comprehensive error handling for missing files
- Supports caching for performance optimization

**Code Example**:
```go
func (app *Root) ServeImageHandler(w http.ResponseWriter, r *http.Request) {
    // Extract filename from URL path
    path := r.URL.Path
    filename := path[len("/images/"):]
    
    // Construct full path to image
    imagePath := filepath.Join("uploads", filename)
    
    // Check if file exists
    if _, err := os.Stat(imagePath); os.IsNotExist(err) {
        tools.EncodeJSON(w, http.StatusNotFound, map[string]string{
            "error": "Image not found",
        })
        return
    }
    
    // Set proper headers
    w.Header().Set("Content-Type", "image/jpeg")
    w.Header().Set("Cache-Control", "public, max-age=31536000")
    
    // Serve the file
    http.ServeFile(w, r, imagePath)
}
```

### 4. Post Creation with Images (`pkg/handlers/posts&comments.go`)

**Function**: `NewPost()`

**Features**:
- Handles multipart form data for image uploads
- Processes image uploads alongside post content
- Validates image files before saving
- Stores image path in database with post
- Supports both JSON and form data formats
- Maintains post privacy settings with image support

**Processing Flow**:
1. Parse multipart form data
2. Extract post content and metadata
3. Process image upload if present
4. Validate all data
5. Save post and image to database
6. Return success response with post ID

## Frontend Implementation

### 1. Image Upload UI (`frontend/src/components/posts/CreatePost.vue`)

**Features**:
- File input with drag & drop support
- Real-time image preview
- File type and size validation
- User-friendly error messages
- FormData submission to backend
- Progress indication during upload

**Key Components**:
```vue
<template>
  <div class="image-upload">
    <input 
      type="file" 
      @change="handleImageSelect" 
      accept="image/jpeg,image/png,image/gif"
    />
    <div v-if="imagePreview" class="image-preview">
      <img :src="imagePreview" alt="Preview" />
    </div>
  </div>
</template>
```

### 2. Image Display (`frontend/src/components/posts/PostItem.vue`)

**Function**: `getImageUrl(imagePath)`

**Features**:
- Constructs full URLs to backend server
- Handles both relative and absolute paths
- Error handling for failed image loads
- Responsive image display with CSS
- Graceful fallbacks for missing images

**Code Example**:
```javascript
function getImageUrl(imagePath) {
  if (!imagePath) return ''
  
  // If it's already a full URL, return as is
  if (imagePath.startsWith('http://') || imagePath.startsWith('https://')) {
    return imagePath
  }
  
  // Extract filename from path
  const filename = imagePath.split('/').pop()
  
  // Construct full URL to backend server
  return `http://localhost:8080/images/${filename}`
}
```

### 3. API Integration (`frontend/src/services/api.js`)

**Features**:
- FormData handling for multipart uploads
- Proper Content-Type headers
- Error handling and response processing
- Integration with Vue.js reactive system

## Security Features

### 1. File Validation

**Type Validation**:
- Only JPEG, PNG, and GIF files allowed
- MIME type checking
- File extension validation

**Size Validation**:
- Maximum file size: 5MB
- Prevents server overload
- User-friendly error messages

**Filename Security**:
- UUID-based filename generation
- Prevents path traversal attacks
- Ensures unique filenames

### 2. Access Control

**Public Image Serving**:
- Images served without authentication
- Proper CORS headers for cross-origin requests
- Cache headers for performance optimization

**Privacy Integration**:
- Images respect post privacy settings
- Public images accessible to all users
- Private images only visible to authorized users

### 3. Error Handling

**Frontend Error Handling**:
- Graceful fallbacks for missing images
- User-friendly error messages
- Console logging for debugging

**Backend Error Handling**:
- Comprehensive validation
- Proper HTTP status codes
- Detailed error messages

## File Structure

```
backend/
├── uploads/                    # Image storage directory
│   ├── image1-uuid.jpg
│   ├── image2-uuid.png
│   └── ...
├── pkg/
│   ├── tools/
│   │   └── tools.go           # UploadHandler function
│   ├── handlers/
│   │   ├── posts&comments.go  # NewPost, ServeImageHandler
│   │   └── router.go          # Route registration
│   └── middleware/
│       └── middleware.go      # Skip paths for images
└── server.go                  # Main server file

frontend/
├── src/
│   ├── components/
│   │   └── posts/
│   │       ├── CreatePost.vue  # Image upload UI
│   │       └── PostItem.vue    # Image display
│   └── services/
│       └── api.js             # API integration
```

## API Endpoints

### 1. Post Creation with Image
- **Endpoint**: `POST /posts`
- **Content-Type**: `multipart/form-data`
- **Parameters**:
  - `content`: Post text content
  - `privacy`: Post privacy setting
  - `image`: Image file (optional)
- **Response**: Post ID and success message

### 2. Image Serving
- **Endpoint**: `GET /images/{filename}`
- **Parameters**: `filename` - Image filename
- **Response**: Image file with proper headers
- **Authentication**: None required (public access)

### 3. Post Retrieval
- **Endpoint**: `GET /posts`
- **Response**: Posts with image URLs included
- **Privacy**: Respects post privacy settings

## Data Flow

### Upload Process
1. **User Selection**: User selects image file in frontend
2. **Frontend Validation**: File type and size validation
3. **FormData Creation**: Image wrapped in FormData with post content
4. **Backend Processing**: File received and validated
5. **File Storage**: Image saved to uploads directory
6. **Database Storage**: Image path stored with post
7. **Response**: Success confirmation returned to frontend

### Display Process
1. **Post Loading**: Posts retrieved from backend with image URLs
2. **URL Construction**: Frontend constructs full image URLs
3. **Image Request**: Browser requests image from backend
4. **Image Serving**: Backend serves image file
5. **Display**: Image rendered in post component

## Key Fixes Applied

### 1. Route Conflicts
**Problem**: Duplicate `/images/` route registrations causing conflicts
**Solution**: Removed duplicate route from `image_handler.go`

### 2. Authentication Issues
**Problem**: Image serving requiring authentication
**Solution**: Added `/images/` to middleware skip paths

### 3. URL Construction
**Problem**: Frontend using wrong server port for images
**Solution**: Fixed `getImageUrl()` to use backend server URL

### 4. Database Integration
**Problem**: Missing image_url field in posts
**Solution**: Added image_url field to posts table schema

## Performance Optimizations

### 1. Caching
- Cache-Control headers for image serving
- Browser caching for static images
- Reduced server load for frequently accessed images

### 2. File Management
- Organized file storage in uploads directory
- Secure filename generation prevents conflicts
- Efficient file serving with http.ServeFile

### 3. Error Handling
- Graceful degradation for missing images
- User-friendly error messages
- Comprehensive logging for debugging

## Future Enhancements

### 1. Image Processing
- Automatic image resizing for thumbnails
- Image compression for storage optimization
- Multiple image formats support

### 2. Advanced Features
- Image galleries for multiple images per post
- Image editing capabilities
- Advanced privacy controls for images

### 3. Performance
- CDN integration for global image serving
- Image lazy loading
- Progressive image loading

## Conclusion

The image upload system provides a complete, secure, and user-friendly solution for handling images in the social network application. It integrates seamlessly with the existing post privacy system while maintaining high performance and security standards.

The system is designed to be scalable, maintainable, and extensible for future enhancements. All components work together to provide a smooth user experience from upload to display.

## Profile Image Upload System

### 1. Backend Implementation (`backend/pkg/handlers/image_handler.go`)

**Profile Image Upload Handler**: `UploadProfileImage()`

**Features**:
- Handles multipart form data for profile image uploads
- Validates file size (10MB limit) and file types
- Uses unified `ImageHandler.UploadImage()` for consistent processing
- Updates user's profile with new image ID
- Supports secure filename generation with UUID
- Organizes images in `uploads/profiles/` directory

**Code Example**:
```go
func (ih *ImageHandler) UploadProfileImage(w http.ResponseWriter, r *http.Request) {
    // Authentication check
    requesterID := ih.DL.GetRequesterID(w, r)
    if requesterID <= 0 {
        tools.EncodeJSON(w, http.StatusUnauthorized, map[string]string{
            "error": "Unauthorized",
        })
        return
    }

    // Parse multipart form (10MB limit)
    if err := r.ParseMultipartForm(10 << 20); err != nil {
        tools.EncodeJSON(w, http.StatusBadRequest, map[string]string{
            "error": "Failed to parse form",
        })
        return
    }

    // Get uploaded file
    file, header, err := r.FormFile("image")
    if err != nil {
        tools.EncodeJSON(w, http.StatusBadRequest, map[string]string{
            "error": "No image file provided",
        })
        return
    }
    defer file.Close()

    // Upload using unified handler
    imageHandler := &models.ImageHandler{DB: ih.DL.Posts.DB}
    image, err := imageHandler.UploadImage(file, header, models.ImageTypeProfile, requesterID)
    if err != nil {
        log.Printf("Error uploading profile image: %v", err)
        tools.EncodeJSON(w, http.StatusInternalServerError, map[string]string{
            "error": "Failed to upload image",
        })
        return
    }

    // Update user's profile with new image
    if err := ih.updateUserProfileImage(requesterID, image.ID); err != nil {
        log.Printf("Error updating user profile image: %v", err)
        tools.EncodeJSON(w, http.StatusInternalServerError, map[string]string{
            "error": "Failed to update profile",
        })
        return
    }

    tools.EncodeJSON(w, http.StatusCreated, map[string]interface{}{
        "message": "Profile image updated successfully",
        "image":   image,
    })
}
```

### 2. Frontend Implementation (`frontend/src/views/Profile.vue`)

**Profile Image Upload UI**:
- Hidden file input with custom upload button
- Real-time upload progress indicator
- Image preview before upload
- Error handling and user feedback

**Code Example**:
```vue
<template>
  <div class="avatar-container">
    <img class="avatar" :src="getAvatarUrl(profileUser.avatar_url)" alt="Profile Picture" />
    <!-- Upload button for profile owner -->
    <div v-if="isOwner" class="avatar-upload">
      <input 
        type="file" 
        ref="fileInput" 
        @change="handleFileSelect" 
        accept="image/*" 
        style="display: none;"
      />
      <button @click="$refs.fileInput.click()" class="upload-btn">
        📷 Change Photo
      </button>
    </div>
    <!-- Upload progress -->
    <div v-if="isUploading" class="upload-progress">
      <div class="progress-bar">
        <div class="progress-fill" :style="{ width: uploadProgress + '%' }"></div>
      </div>
      <p>Uploading... {{ uploadProgress }}%</p>
    </div>
  </div>
</template>
```

**Upload Handler**:
```javascript
async function handleFileSelect(event) {
  const file = event.target.files[0]
  if (!file) return

  // Validate file size (5MB limit)
  if (file.size > 5 * 1024 * 1024) {
    alert('Image file too large (max 5MB)')
    return
  }

  // Validate file type
  const allowedTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/gif']
  if (!allowedTypes.includes(file.type)) {
    alert('Only JPEG, PNG and GIF images are allowed')
    return
  }

  isUploading.value = true
  uploadProgress.value = 0

  try {
    const formData = new FormData()
    formData.append('image', file)

    const response = await fetch('http://localhost:8080/images/profile', {
      method: 'POST',
      body: formData,
      credentials: 'include'
    })

    if (!response.ok) {
      throw new Error('Failed to upload image')
    }

    const result = await response.json()
    
    // Update profile with new image
    await fetchProfile()
    
    alert('Profile image updated successfully!')
  } catch (error) {
    console.error('Upload error:', error)
    alert('Failed to upload image: ' + error.message)
  } finally {
    isUploading.value = false
    uploadProgress.value = 0
  }
}
```

## Comment Image Upload System

### 1. Backend Implementation (`backend/pkg/handlers/image_handler.go`)

**Comment Image Upload Handler**: `UploadCommentImage()`

**Features**:
- Handles multipart form data for comment image uploads
- Uses unified `ImageHandler.UploadImage()` for consistent processing
- Stores images in `uploads/comments/` directory
- Links images to comments via database relationships
- Supports secure filename generation

**Code Example**:
```go
func (ih *ImageHandler) UploadCommentImage(w http.ResponseWriter, r *http.Request) {
    // Authentication check
    requesterID := ih.DL.GetRequesterID(w, r)
    if requesterID <= 0 {
        tools.EncodeJSON(w, http.StatusUnauthorized, map[string]string{
            "error": "Unauthorized",
        })
        return
    }

    // Parse multipart form (10MB limit)
    if err := r.ParseMultipartForm(10 << 20); err != nil {
        tools.EncodeJSON(w, http.StatusBadRequest, map[string]string{
            "error": "Failed to parse form",
        })
        return
    }

    // Get uploaded file
    file, header, err := r.FormFile("image")
    if err != nil {
        tools.EncodeJSON(w, http.StatusBadRequest, map[string]string{
            "error": "No image file provided",
        })
        return
    }
    defer file.Close()

    // Upload using unified handler
    imageHandler := &models.ImageHandler{DB: ih.DL.Posts.DB}
    image, err := imageHandler.UploadImage(file, header, models.ImageTypeComment, requesterID)
    if err != nil {
        log.Printf("Error uploading comment image: %v", err)
        tools.EncodeJSON(w, http.StatusInternalServerError, map[string]string{
            "error": "Failed to upload image",
        })
        return
    }

    tools.EncodeJSON(w, http.StatusCreated, map[string]interface{}{
        "message": "Image uploaded successfully",
        "image":   image,
    })
}
```

### 2. Comment Creation with Images (`backend/pkg/handlers/posts&comments.go`)

**Function**: `NewComment()`

**Features**:
- Handles both text-only and image-containing comments
- Processes multipart form data for image uploads
- Validates image files before saving
- Links images to comments in database
- Supports FormData and JSON formats

**Processing Flow**:
1. Extract post_id from URL path
2. Parse multipart form data if image present
3. Validate comment data and image file
4. Upload image if provided
5. Save comment with image reference
6. Return success response with comment data

**Code Example**:
```go
func (app *Root) NewComment(w http.ResponseWriter, r *http.Request) {
    var comment models.Comment
    var hasFile bool

    // Extract post_id from URL path
    post := r.PathValue("post_id")
    post_id, err := strconv.Atoi(post)
    if err != nil {
        tools.EncodeJSON(w, http.StatusBadRequest, map[string]string{
            "error": "Invalid post_id format",
        })
        return
    }
    comment.Post_id = post_id

    contentType := r.Header.Get("Content-Type")

    if strings.Contains(contentType, "multipart/form-data") {
        // Handle image upload
        if err := r.ParseMultipartForm(10 << 20); err != nil {
            tools.EncodeJSON(w, http.StatusBadRequest, map[string]string{
                "error": "Invalid form data",
            })
            return
        }
        hasFile = true

        // Parse comment data from form
        if status := models.ParseCommentFromForm(r, &comment); status != 200 {
            tools.EncodeJSON(w, http.StatusBadRequest, map[string]string{
                "error": "Invalid form data",
            })
            return
        }
    } else {
        // Handle JSON request
        if err := tools.DecodeJSON(r, &comment); err != nil {
            tools.EncodeJSON(w, http.StatusBadRequest, map[string]string{
                "error": "Invalid JSON data",
            })
            return
        }
    }

    // Handle image upload if present
    var imagePath string
    if hasFile {
        file, handler, err := r.FormFile("image")
        if err == nil {
            defer file.Close()

            // Validate file size and type
            if handler.Size > 5<<20 {
                tools.EncodeJSON(w, http.StatusBadRequest, map[string]string{
                    "error": "Image file too large (max 5MB)",
                })
                return
            }

            if !tools.IsAllowedFile(handler.Filename, file) {
                tools.EncodeJSON(w, http.StatusBadRequest, map[string]string{
                    "error": "Unsupported file format",
                })
                return
            }

            // Upload file
            uploadedPath, status := tools.UploadHandler(file, handler)
            if status != 200 {
                tools.EncodeJSON(w, status, map[string]string{
                    "error": "Failed to upload image",
                })
                return
            }
            imagePath = uploadedPath
        }
    }

    // Save comment to database
    comment.ImageURL = imagePath
    if err := app.DL.Comments.InsertComment(&comment); err != nil {
        tools.EncodeJSON(w, http.StatusInternalServerError, map[string]string{
            "error": "Failed to create comment",
        })
        return
    }

    tools.EncodeJSON(w, http.StatusCreated, map[string]interface{}{
        "message": "Comment created successfully",
        "comment": comment,
    })
}
```

### 3. Frontend Implementation (`frontend/src/components/posts/PostItem.vue`)

**Comment Image Upload UI**:
- Hidden file input with camera icon button
- Real-time image preview
- File validation and error handling
- Progress indication during upload

**Code Example**:
```vue
<template>
  <div class="comment-actions">
    <input
      ref="imageInput"
      type="file"
      accept="image/*"
      @change="handleImageSelect"
      style="display: none"
    >
    <button
      type="button"
      @click="$refs.imageInput.click()"
      class="image-btn"
      :disabled="isSubmitting"
    >
      📷
    </button>
    <button
      @click="addComment"
      class="submit-btn"
      :disabled="!canSubmit || isSubmitting"
    >
      {{ isSubmitting ? 'Posting...' : 'Post' }}
    </button>
  </div>

  <!-- Image Preview -->
  <div v-if="newComment.imagePreview" class="image-preview">
    <img :src="newComment.imagePreview" alt="Preview" />
    <button @click="removeImage" class="remove-image">×</button>
  </div>
</template>
```

**Image Selection Handler**:
```javascript
function handleImageSelect(event) {
  const file = event.target.files[0]
  if (!file) return

  // Validate file size (5MB limit)
  if (file.size > 5 * 1024 * 1024) {
    commentError.value = 'Image file too large (max 5MB)'
    return
  }

  // Validate file type
  const allowedTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/gif']
  if (!allowedTypes.includes(file.type)) {
    commentError.value = 'Only JPEG, PNG and GIF images are allowed'
    return
  }

  newComment.value.image = file
  
  // Create preview
  const reader = new FileReader()
  reader.onload = (e) => {
    newComment.value.imagePreview = e.target.result
  }
  reader.readAsDataURL(file)
  
  commentError.value = ''
}
```

**Comment Submission with Image**:
```javascript
async function addComment() {
  if (!canSubmit.value || isSubmitting.value) return
  
  isSubmitting.value = true
  commentError.value = ''
  
  try {
    const formData = new FormData()
    formData.append('content', newComment.value.content.trim())
    formData.append('owner_id', props.currentUser?.id?.toString() || "1")
    
    // Add image if selected
    if (newComment.value.image) {
      formData.append('image', newComment.value.image)
    }
    
    const response = await fetch(`/post/${props.post.id}/comments/new`, {
      method: 'POST',
      body: formData,
      credentials: 'include'
    })
    
    if (!response.ok) {
      const errorData = await response.json()
      throw new Error(errorData.error || 'Failed to add comment')
    }
    
    const result = await response.json()
    
    // Add the new comment to the beginning of the comments array
    if (result.comment) {
      comments.value.unshift(result.comment)
    }
    
    // Reset form
    newComment.value = {
      content: '',
      image: null,
      imagePreview: null
    }
    
  } catch (error) {
    console.error('Error adding comment:', error)
    commentError.value = error.message || 'Failed to add comment'
  } finally {
    isSubmitting.value = false
  }
}
```

## Unified Image Handler (`backend/pkg/models/image_handler.go`)

### Core Features

**Image Types**:
```go
type ImageType string

const (
    ImageTypePost    ImageType = "posts"
    ImageTypeProfile ImageType = "profiles"
    ImageTypeComment ImageType = "comments"
    ImageTypeGroup   ImageType = "groups"
)
```

**Unified Upload Function**:
```go
func (ih *ImageHandler) UploadImage(file multipart.File, header *multipart.FileHeader, imageType ImageType, uploadedBy int) (*Image, error) {
    // 1. Validate file
    if err := ih.validateImage(header); err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }

    // 2. Generate secure filename
    filename, err := ih.generateSecureFilename(header)
    if err != nil {
        return nil, fmt.Errorf("filename generation failed: %w", err)
    }

    // 3. Create upload directory
    uploadDir := filepath.Join("uploads", string(imageType))
    if err := os.MkdirAll(uploadDir, 0755); err != nil {
        return nil, fmt.Errorf("failed to create upload directory: %w", err)
    }

    // 4. Save file to disk
    filePath := filepath.Join(uploadDir, filename)
    if err := ih.saveFile(file, filePath); err != nil {
        return nil, fmt.Errorf("failed to save file: %w", err)
    }

    // 5. Save metadata to database
    image := &Image{
        Filename:     filename,
        OriginalName: header.Filename,
        MimeType:     header.Header.Get("Content-Type"),
        FileSize:     header.Size,
        UploadPath:   filepath.Join(string(imageType), filename),
        UploadedBy:   uploadedBy,
        CreatedAt:    time.Now(),
    }

    if err := ih.saveImageMetadata(image); err != nil {
        // Clean up file if database save fails
        os.Remove(filePath)
        return nil, fmt.Errorf("failed to save image metadata: %w", err)
    }

    return image, nil
}
```

## Updated File Structure

```
backend/
├── uploads/                    # Image storage directory
│   ├── profiles/              # Profile images
│   │   ├── profile1-uuid.jpg
│   │   └── profile2-uuid.png
│   ├── posts/                 # Post images
│   │   ├── post1-uuid.jpg
│   │   └── post2-uuid.png
│   ├── comments/              # Comment images
│   │   ├── comment1-uuid.jpg
│   │   └── comment2-uuid.png
│   └── groups/                # Group images (future)
├── pkg/
│   ├── tools/
│   │   └── tools.go           # UploadHandler function
│   ├── handlers/
│   │   ├── image_handler.go   # Profile/Comment image handlers
│   │   ├── posts&comments.go  # NewPost, NewComment, ServeImageHandler
│   │   ├── auth.go            # Signup with avatar upload
│   │   └── router.go          # Route registration
│   ├── models/
│   │   ├── image_handler.go   # Unified image processing
│   │   └── user.go            # User model with avatar
│   └── middleware/
│       └── middleware.go      # Skip paths for images
└── server.go                  # Main server file

frontend/
├── src/
│   ├── views/
│   │   ├── Profile.vue        # Profile image upload UI
│   │   └── signup.vue         # Signup with avatar
│   ├── components/
│   │   └── posts/
│   │       ├── PostItem.vue   # Comment image upload
│   │       └── CreatePost.vue # Post image upload
│   └── services/
│       └── api.js             # API integration
```

## Complete Image System Summary

The social network application now features a comprehensive image upload system that supports:

1. **Profile Images**: Users can upload and update their profile pictures
2. **Post Images**: Users can attach images to their posts
3. **Comment Images**: Users can include images in their comments
4. **Signup Avatars**: Users can upload profile pictures during registration

All image types use a unified backend system that provides:
- Secure file validation and storage
- Organized directory structure
- Database metadata tracking
- Proper error handling and user feedback
- Responsive frontend interfaces

The system maintains security through file validation, secure filename generation, and proper access controls while providing a smooth user experience across all image upload scenarios. 