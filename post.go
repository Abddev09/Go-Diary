package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// Context dan user ma'lumotlarini olish
func getUserFromContext(r *http.Request) (int, string, error) {
	userID, ok := r.Context().Value(userIDKey).(int)
	if !ok {
		return 0, "", fmt.Errorf("user ID not found in context")
	}
	
	username, ok := r.Context().Value(usernameKey).(string)
	if !ok {
		return 0, "", fmt.Errorf("username not found in context")
	}
	
	return userID, username, nil
}

func getAllPostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rows, err := db.Query(`
		SELECT id, title, description, image_url, image_file, created_at, author_id, author_username 
		FROM posts 
		ORDER BY created_at DESC
	`)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(
			&post.ID, &post.Title, &post.Description, &post.ImageURL, 
			&post.ImageFile, &post.CreatedAt, &post.AuthorID, &post.AuthorUsername,
		)
		if err != nil {
			http.Error(w, "Error scanning posts", http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func getMyPostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, _, err := getUserFromContext(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	rows, err := db.Query(`
		SELECT id, title, description, image_url, image_file, created_at, author_id, author_username 
		FROM posts 
		WHERE author_id = $1 
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(
			&post.ID, &post.Title, &post.Description, &post.ImageURL, 
			&post.ImageFile, &post.CreatedAt, &post.AuthorID, &post.AuthorUsername,
		)
		if err != nil {
			http.Error(w, "Error scanning posts", http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func createPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, username, err := getUserFromContext(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Content-Type tekshiramiz
	contentType := r.Header.Get("Content-Type")
	
	var post Post
	post.AuthorID = userID
	post.AuthorUsername = username

	if strings.Contains(contentType, "multipart/form-data") {
		// Fayl upload qilish
		err := r.ParseMultipartForm(10 << 20) // 10 MB limit
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		post.Title = r.FormValue("title")
		post.Description = r.FormValue("description")
		imageURL := r.FormValue("image_url")
		if imageURL != "" {
			post.ImageURL = &imageURL
		}

		// Fayl upload
		file, handler, err := r.FormFile("image_file")
		if err == nil {
			defer file.Close()

			// uploads papkasini yaratamiz
			os.MkdirAll("uploads", os.ModePerm)

			// Fayl nomini yaratamiz
			filename := fmt.Sprintf("%d_%s", time.Now().Unix(), handler.Filename)
			filepath := filepath.Join("uploads", filename)

			// Faylni saqlaymiz
			dst, err := os.Create(filepath)
			if err != nil {
				http.Error(w, "Error saving file", http.StatusInternalServerError)
				return
			}
			defer dst.Close()

			_, err = io.Copy(dst, file)
			if err != nil {
				http.Error(w, "Error copying file", http.StatusInternalServerError)
				return
			}

			post.ImageFile = &filename
		}
	} else {
		// JSON ma'lumot
		var req PostRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		post.Title = req.Title
		post.Description = req.Description
		if req.ImageURL != "" {
			post.ImageURL = &req.ImageURL
		}
	}

	if post.Title == "" || post.Description == "" {
		http.Error(w, "Title and description required", http.StatusBadRequest)
		return
	}

	// Postni bazaga qo'shamiz
	var postID int
	err = db.QueryRow(`
		INSERT INTO posts (title, description, image_url, image_file, author_id, author_username) 
		VALUES ($1, $2, $3, $4, $5, $6) 
		RETURNING id
	`, post.Title, post.Description, post.ImageURL, post.ImageFile, post.AuthorID, post.AuthorUsername).Scan(&postID)

	if err != nil {
		http.Error(w, "Error creating post", http.StatusInternalServerError)
		return
	}

	post.ID = postID
	post.CreatedAt = time.Now()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

func updatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	postID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	userID, _, err := getUserFromContext(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Post egasini tekshiramiz
	var authorID int
	err = db.QueryRow("SELECT author_id FROM posts WHERE id = $1", postID).Scan(&authorID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if authorID != userID {
		http.Error(w, "You can only update your own posts", http.StatusForbidden)
		return
	}

	var req PostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Title == "" || req.Description == "" {
		http.Error(w, "Title and description required", http.StatusBadRequest)
		return
	}

	// Postni yangilaymiz
	var imageURL *string
	if req.ImageURL != "" {
		imageURL = &req.ImageURL
	}

	_, err = db.Exec(`
		UPDATE posts 
		SET title = $1, description = $2, image_url = $3 
		WHERE id = $4
	`, req.Title, req.Description, imageURL, postID)

	if err != nil {
		http.Error(w, "Error updating post", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Post updated successfully"})
}

func deletePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	postID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	userID, _, err := getUserFromContext(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Post egasini tekshiramiz
	var authorID int
	var imageFile *string
	err = db.QueryRow("SELECT author_id, image_file FROM posts WHERE id = $1", postID).Scan(&authorID, &imageFile)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if authorID != userID {
		http.Error(w, "You can only delete your own posts", http.StatusForbidden)
		return
	}

	// Postni o'chiramiz
	_, err = db.Exec("DELETE FROM posts WHERE id = $1", postID)
	if err != nil {
		http.Error(w, "Error deleting post", http.StatusInternalServerError)
		return
	}

	// Agar fayl bo'lsa, uni ham o'chiramiz
	if imageFile != nil && *imageFile != "" {
		filepath := filepath.Join("uploads", *imageFile)
		os.Remove(filepath)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Post deleted successfully"})
}
