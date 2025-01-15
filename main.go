package main

import (
	"log"
	"net/http"
)

func main() {

	if err := openDB(); err != nil {
		log.Printf("ERROR connecting to database %v", err)
	}
	defer closeDB()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static")))) // Serve static files (css, img)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { 
		http.ServeFile(w, r, "web/templates/index.html")
	})

	addr := ":8080"
	log.Printf("Server is running on http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
