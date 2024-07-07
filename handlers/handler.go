package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Onyekachukwu-Nweke/hng-backend-stage-2/config"
	"github.com/Onyekachukwu-Nweke/hng-backend-stage-2/database"
	"github.com/Onyekachukwu-Nweke/hng-backend-stage-2/middleware"
	"github.com/Onyekachukwu-Nweke/hng-backend-stage-2/repositories"
	"github.com/Onyekachukwu-Nweke/hng-backend-stage-2/services"
	"github.com/gorilla/mux"
)

type Handler struct {
	Server *http.Server
	Router *mux.Router
}

func NewHandler() *Handler {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)
	database.Connect(connectionString)

	// Initialize repositories
	userRepo := repositories.NewUserRepository(database.GetDB())
	orgRepo := repositories.NewOrganisationRepository(database.GetDB())

	// Initialize services
	userService := services.NewUserService(userRepo, orgRepo)
	orgService := services.NewOrganisationService(orgRepo)

	// Initialize handlers
	authHandler := NewAuthHandler(userService)
	organisationHandler := NewOrganisationHandler(orgService)

	// Initialize router
	router := mux.NewRouter()

	// Auth routes
	router.HandleFunc("/auth/register", authHandler.Register).Methods("POST")
	router.HandleFunc("/auth/login", authHandler.Login).Methods("POST")

	// Protected routes
	protectedRouter := router.PathPrefix("/api").Subrouter()
	protectedRouter.Use(middleware.JWTMiddleware)
	protectedRouter.HandleFunc("/users/{id}", organisationHandler.GetUser).Methods("GET")
	protectedRouter.HandleFunc("/organisations", organisationHandler.GetOrganisations).Methods("GET")
	protectedRouter.HandleFunc("/organisations/{orgId}", organisationHandler.GetOrganisation).Methods("GET")
	protectedRouter.HandleFunc("/organisations", organisationHandler.Create).Methods("POST")
	protectedRouter.HandleFunc("/organisations/{orgId}/users", organisationHandler.AddUserToOrganisation).Methods("POST")

	server := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: router,
	}

	return &Handler{
		Server: server,
		Router: router,
	}
}

func (h *Handler) Serve() error {
	go func() {
		if err := h.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v\n", h.Server.Addr, err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := h.Server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")

	return nil
}
