package routes

import (
	"sodality/controllers"
	middlewares "sodality/handlers"

	"github.com/gorilla/mux"
)

// Routes -> define endpoints
func Routes() *mux.Router {
	router := mux.NewRouter()

	api := router.PathPrefix("/api/v1").Subrouter()

	// User API routes

	user := api.PathPrefix("/user").Subrouter()
	user.HandleFunc("/register", controllers.RegisterUser).Methods("POST")
	user.HandleFunc("/login", controllers.LoginUser).Methods("POST")
	// user.HandleFunc("/me", middlewares.IsAuthorized(controllers.GetMe)).Methods("GET")
	// user.HandleFunc("/me", middlewares.IsAuthorized(controllers.UpdateUser)).Methods("PUT")
	// user.HandleFunc("/{username}", controllers.GetUser).Methods("GET")

	// Content API routes

	content := api.PathPrefix("/content").Subrouter()
	content.HandleFunc("/post", middlewares.IsAuthorized(controllers.PostContent)).Methods("POST")

	// challenge.HandleFunc("/", middlewares.IsAuthorized(controllers.CreateChallenge)).Methods("POST")
	// challenge.HandleFunc("/user/{username}", controllers.GetChallenges).Methods("GET")
	// challenge.HandleFunc("/{id}", middlewares.IsAuthorized(controllers.GetChallenge)).Methods("GET")
	// challenge.HandleFunc("/{id}", middlewares.IsAuthorized(controllers.UpdateChallenge)).Methods("PUT")
	// challenge.HandleFunc("/{id}", middlewares.IsAuthorized(controllers.DeleteChallenge)).Methods("DELETE")
	// challenge.HandleFunc("/{id}/join/", middlewares.IsAuthorized(controllers.JoinChallenge)).Methods("POST")
	// challenge.HandleFunc("/{id}/unjoin/", middlewares.IsAuthorized(controllers.UnJoinChallenge)).Methods("POST")

	// api.HandleFunc("/person", controllers.CreatePersonEndpoint).Methods("POST")
	// api.HandleFunc("/people", middlewares.IsAuthorized(controllers.GetPeopleEndpoint)).Methods("GET")
	// api.HandleFunc("/person/{id}", controllers.GetPersonEndpoint).Methods("GET")
	// api.HandleFunc("/person/{id}", controllers.DeletePersonEndpoint).Methods("DELETE")
	// api.HandleFunc("/person/{id}", controllers.UpdatePersonEndpoint).Methods("PUT")

	// router.HandleFunc("/upload", controllers.UploadFileEndpoint).Methods("POST")
	// router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./uploaded/"))))
	return router
}
