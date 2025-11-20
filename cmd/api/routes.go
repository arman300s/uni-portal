package main

import (
	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/arman300s/uni-portal/internal/http/controllers"
	"github.com/arman300s/uni-portal/pkg/middleware"
)

// RouteDeps groups all controllers required by the router.
type RouteDeps struct {
	Auth         *controllers.AuthController
	User         *controllers.UserController
	AdminSubject *controllers.AdminSubjectController
	Student      *controllers.StudentController
	Teacher      *controllers.TeacherController
}

func SetupRoutes(r *mux.Router, deps RouteDeps) {
	// Swagger docs
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Public routes
	r.HandleFunc("/signup", deps.Auth.Signup).Methods("POST")
	r.HandleFunc("/login", deps.Auth.Login).Methods("POST")
	r.Handle("/me", middleware.JWTAuth(http.HandlerFunc(deps.User.Me))).Methods("GET")

	// Admin routes
	admin := r.PathPrefix("/admin").Subrouter()
	admin.Use(middleware.JWTAuth)
	admin.Use(middleware.LoadUserMiddleware)
	admin.Use(middleware.RequireRole("admin"))

	// User management
	admin.HandleFunc("/users", deps.User.ListUsers).Methods("GET")
	admin.HandleFunc("/users/{id}", deps.User.GetUser).Methods("GET")
	admin.HandleFunc("/users/{id}", deps.User.UpdateUser).Methods("PUT")
	admin.HandleFunc("/users/{id}", deps.User.DeleteUser).Methods("DELETE")
	admin.HandleFunc("/users/create", deps.User.CreateUser).Methods("POST")

	// Subject management
	admin.HandleFunc("/subjects", deps.AdminSubject.ListSubjects).Methods("GET")
	admin.HandleFunc("/subjects/{id}", deps.AdminSubject.GetSubject).Methods("GET")
	admin.HandleFunc("/subjects", deps.AdminSubject.CreateSubject).Methods("POST")
	admin.HandleFunc("/subjects/{id}", deps.AdminSubject.UpdateSubject).Methods("PUT")
	admin.HandleFunc("/subjects/{id}", deps.AdminSubject.DeleteSubject).Methods("DELETE")

	// Student routes
	student := r.PathPrefix("/student").Subrouter()
	student.Use(middleware.JWTAuth)
	student.Use(middleware.LoadUserMiddleware)
	student.Use(middleware.RequireRole("student"))
	student.HandleFunc("/subjects", deps.Student.ListSubjects).Methods("GET")

	// Teacher routes
	teacher := r.PathPrefix("/teacher").Subrouter()
	teacher.Use(middleware.JWTAuth)
	teacher.Use(middleware.LoadUserMiddleware)
	teacher.Use(middleware.RequireRole("teacher"))
	teacher.HandleFunc("/subjects", deps.Teacher.ListMySubjects).Methods("GET")
}
