package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	_ "github.com/arman300s/uni-portal/cmd/api/docs"
	"github.com/arman300s/uni-portal/internal/core/repositories"
	"github.com/arman300s/uni-portal/internal/core/services"
	"github.com/arman300s/uni-portal/internal/http/controllers"
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

	userRepo := repositories.NewUserRepository(db.DB)
	roleRepo := repositories.NewRoleRepository(db.DB)
	subjectRepo := repositories.NewSubjectRepository(db.DB)

	authService := services.NewAuthService(userRepo, roleRepo)
	userService := services.NewUserService(userRepo, roleRepo)
	subjectService := services.NewSubjectService(subjectRepo, userRepo)

	routeDeps := RouteDeps{
		Auth:         controllers.NewAuthController(authService),
		User:         controllers.NewUserController(userService),
		AdminSubject: controllers.NewAdminSubjectController(subjectService),
		Student:      controllers.NewStudentController(subjectService),
		Teacher:      controllers.NewTeacherController(subjectService),
	}

	r := mux.NewRouter()
	SetupRoutes(r, routeDeps)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8079"
	}

	log.Printf("🚀 Server started on port %s\n", port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, r))
}
