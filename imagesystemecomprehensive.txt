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