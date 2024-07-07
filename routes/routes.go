package routes

import (
	"fmt"
	"net/http"

	"github.com/Onyekachukwu-Nweke/hng-backend-stage-2/handlers"
	"github.com/Onyekachukwu-Nweke/hng-backend-stage-2/middleware"
	"github.com/gorilla/mux"
)

func InitRoutes(authHandler *handlers.AuthHandler, orgHandler *handlers.OrganisationHandler) *mux.Router {
    // Initialize router
	router := mux.NewRouter()

	// Auth routes
	router.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World")
	})
	router.HandleFunc("/auth/register", authHandler.Register).Methods("POST")
	router.HandleFunc("/auth/login", authHandler.Login).Methods("POST")

	// Protected routes
	protectedRouter := router.PathPrefix("/api").Subrouter()
	protectedRouter.Use(middleware.JWTMiddleware)
	protectedRouter.HandleFunc("/users/{id}", orgHandler.GetUser).Methods("GET")
	protectedRouter.HandleFunc("/organisations", orgHandler.GetOrganisations).Methods("GET")
	protectedRouter.HandleFunc("/organisations/{orgId}", orgHandler.GetOrganisation).Methods("GET")
	protectedRouter.HandleFunc("/organisations", orgHandler.Create).Methods("POST")
	protectedRouter.HandleFunc("/organisations/{orgId}/users", orgHandler.AddUserToOrganisation).Methods("POST")

	return router
}
