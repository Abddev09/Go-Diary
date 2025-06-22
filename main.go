package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB

func initDB() {
	var err error
	// PostgreSQL connection string - bu qatorni o'zgartiring
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// Default connection string
		dbURL = "postgresql://blog_user:0QyfSUPcO6kpq9ya5HkTLeWqz7mJaqwy@dpg-d1c1k83e5dus73f3h6eg-a.oregon-postgres.render.com/blogapi_2ge8"
		// Yoki agar parol bo'lsa:
		// dbURL = "postgres://postgres:yourpassword@localhost/blog_system?sslmode=disable"
		// Yoki yangi foydalanuvchi uchun:
		// dbURL = "postgres://bloguser:blogpass@localhost/blog_system?sslmode=disable"
	}

	log.Printf("Connecting to database: %s", dbURL)
	
	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Error pinging database:", err)
	}

	log.Println("Database connected successfully")
	
	// Jadvallar mavjudligini tekshirish
	testTables()
}

func testTables() {
	// Users jadvalini tekshirish
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM information_schema.tables WHERE table_name = 'users'").Scan(&count)
	if err != nil {
		log.Fatal("Error checking users table:", err)
	}
	if count == 0 {
		log.Fatal("Users table does not exist. Please run the SQL script first.")
	}
	
	// Posts jadvalini tekshirish
	err = db.QueryRow("SELECT COUNT(*) FROM information_schema.tables WHERE table_name = 'posts'").Scan(&count)
	if err != nil {
		log.Fatal("Error checking posts table:", err)
	}
	if count == 0 {
		log.Fatal("Posts table does not exist. Please run the SQL script first.")
	}
	
	log.Println("All tables exist")
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	initDB()
	defer db.Close()

	// uploads papkasini yaratamiz
	os.MkdirAll("uploads", os.ModePerm)

	r := mux.NewRouter()

	// CORS middleware
	r.Use(corsMiddleware)

	// Static fayllar uchun
	r.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads/"))))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	// Auth routes
	r.HandleFunc("/register", registerHandler)
	r.HandleFunc("/login", loginHandler)

	// Post routes
	r.HandleFunc("/posts", getAllPostsHandler).Methods("GET")
	r.HandleFunc("/posts", authMiddleware(createPostHandler)).Methods("POST")
	r.HandleFunc("/posts/mine", authMiddleware(getMyPostsHandler)).Methods("GET")
	r.HandleFunc("/posts/{id}", authMiddleware(updatePostHandler)).Methods("PUT")
	r.HandleFunc("/posts/{id}", authMiddleware(deletePostHandler)).Methods("DELETE")

	// HTML sahifalar
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/all-posts.html")
	})
	r.HandleFunc("/register.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/register.html")
	})
	r.HandleFunc("/login.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/login.html")
	})
	r.HandleFunc("/my-posts.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/my-posts.html")
	})
	r.HandleFunc("/all-posts.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/all-posts.html")
	})

	log.Println("Server starting on :9080")
	log.Fatal(http.ListenAndServe(":9080", r))
}
