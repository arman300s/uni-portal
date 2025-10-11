package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	_ "github.com/arman300s/uni-portal/cmd/api/docs"
	"github.com/arman300s/uni-portal/internal/models"
	"github.com/arman300s/uni-portal/internal/seeder"
	"github.com/arman300s/uni-portal/pkg/db"
)

func main() {
	db.Connect()
	if err := db.DB.AutoMigrate(&models.Role{}, &models.User{}, &models.Subject{}); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	seeder.SeedRoles(db.DB)
	seeder.SeedAdmin(db.DB)
	seeder.SeedTeachers(db.DB)
	seeder.SeedSubjects(db.DB)
	r := mux.NewRouter()
	SetupRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8079"
	}

	log.Printf("🚀 Server started on port %s\n", port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, r))
}
