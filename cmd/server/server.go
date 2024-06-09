package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/clegganator259/golearn/user"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	_ "github.com/mattn/go-sqlite3"
)


func main() {
	repo, err := user.NewSqliteRepo("file:db.sqlite?cache=shared")
	if err != nil {
		log.Fatal(err)
	}

    var (
        key   = []byte("super-secret-key")
        store = sessions.NewCookieStore(key)
    )

	router := mux.NewRouter()

	router.HandleFunc("/", logging(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
	}))


	router.HandleFunc("/users/{id:[0-9]+}/login", logging(Login(store)))
	router.HandleFunc("/users/{id:[0-9]+}/logout", logging(Logout(store)))
	router.HandleFunc("/users/{id:[0-9]+}/secret", logging(Secret(store)))
	router.HandleFunc("/users/{id:[0-9]+}", logging(GetUser(repo)))
	router.HandleFunc("/users/{username}/{password}", logging(CreateUser(repo)))

	static_server := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", static_server))

	http.ListenAndServe(":80", router)
}

func GetUser(repo user.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		reqId, err := strconv.Atoi(mux.Vars(req)["id"])
		if err != nil {
			log.Fatal(err)
		}
		user, err := repo.GetUserById(int64(reqId))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "Got user %d\nUsername: %s\nPassword: %s\nCreated At: %s", user.Id, user.Username, user.Password, user.CreatedAt.Format(time.RFC3339))
	}
}

func CreateUser(repo user.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		username := vars["username"]
		password := vars["password"]
        log.Printf("Creating User with username: %s and password: %s", username, password)
		user, err := repo.CreateUser(username, password)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "Created user with ID %d", user.Id)
	}
}

func Login(store sessions.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		session, _ := store.Get(req, "auth-cookie")

		session.Values["authenticated"] = true
		session.Save(req, w)
	}
}

func Logout(store sessions.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		session, _ := store.Get(req, "auth-cookie")

		session.Values["authenticated"] = false
		session.Save(req, w)
	}
}

func Secret(store sessions.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		session, _ := store.Get(req, "auth-cookie")
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Error(w, "Forbidden", http.StatusForbidden)
            return
		}

        fmt.Fprintf(w, "The cake is a lie")

	}
}

func logging(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path)
		f(w, r)
	}
}
