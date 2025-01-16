package main

import (
	"database/sql"
	"log"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func main() {

	if err := openDB(); err != nil {
		log.Printf("ERROR connecting to database %v", err)
	}
	defer closeDB()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static")))) 


	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { 

		if r.Method != http.MethodPost {
			http.ServeFile(w, r, "web/templates/index.html")
			return
		} 


		email := r.FormValue("email")
    	password := r.FormValue("password")

		var hashedPassword string
		err := DB.QueryRow("SELECT password_hash FROM users WHERE email = $1", email).Scan(&hashedPassword)
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		} else if err != nil {
			http.Error(w, "Failed to fetch user", http.StatusInternalServerError)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		if err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		http.Redirect(w, r, "/main", http.StatusSeeOther)
	})	

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.ServeFile(w, r, "web/templates/register.html")
			return
		} 
		

		email := strings.TrimSpace(r.FormValue("email"))
		password := strings.TrimSpace(r.FormValue("password"))
		if email == "" || password == "" {
			http.Error(w, "Email and password cannot be empty", http.StatusBadRequest)
			return
		}

		var exists bool
		err := DB.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)", email).Scan(&exists)
		if err != nil {
			log.Printf("Error checking email existence: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if exists {
			http.Error(w, "Email already registered", http.StatusConflict)
			return
		}
			

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}

		_, err = DB.Exec("INSERT INTO users (email, password_hash) VALUES ($1, $2)", email, hashedPassword)
		if err != nil {
			log.Printf("Error inserting new user: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/main", http.StatusSeeOther)
	})
	
	http.HandleFunc("/main", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/templates/main.html")
	})	


	addr := ":8080"
	log.Printf("Server is running on http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
