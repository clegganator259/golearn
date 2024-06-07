package main

import (
    "github.com/clegganator259/golearn/user"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"strconv"
	"time"
)

func main() {
	repo, err := NewSqliteRepo("file:db.sqlite?cache=shared")
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
	})

	router.HandleFunc("/users/{username}/{password}", func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		username := vars["username"]
		password := vars["password"]
		user, err := repo.createUser(username, password)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "Created user with ID %d", user.id)
	})

	router.HandleFunc("/users/{id}", func(w http.ResponseWriter, req *http.Request) {
		reqId, err := strconv.Atoi(mux.Vars(req)["id"])
		if err != nil {
			log.Fatal(err)
		}
		user, err := repo.getUserById(int64(reqId))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "Got user %d\nUsername: %s\nPassword: %s\nCreated At: %s", user.id, user.username, user.password, user.createdAt.Format(time.RFC3339))
	})

	static_server := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", static_server))

	http.ListenAndServe(":80", router)
}

