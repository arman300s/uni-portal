package main

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/arman300s/uni-portal/internal/api"
	"github.com/arman300s/uni-portal/pkg/middleware"
)

func SetupRoutes(r *mux.Router) {
	// Public routes
	r.HandleFunc("/signup", api.SignupHandler).Methods("POST")
	r.HandleFunc("/login", api.LoginHandler).Methods("POST")
	r.Handle("/me", middleware.JWTAuth(http.HandlerFunc(api.MeHandler))).Methods("GET")

	// Admin routes
	admin := r.PathPrefix("/admin").Subrouter()
	admin.Use(middleware.JWTAuth)
	admin.Use(middleware.LoadUserMiddleware)
	admin.Use(middleware.RequireRole("admin"))

	// User management
	admin.HandleFunc("/users", api.AdminListUsersHandler).Methods("GET")
	admin.HandleFunc("/users/{id}", api.AdminGetUserHandler).Methods("GET")
	admin.HandleFunc("/users/{id}", api.AdminUpdateUserHandler).Methods("PUT")
	admin.HandleFunc("/users/{id}", api.AdminDeleteUserHandler).Methods("DELETE")
	admin.HandleFunc("/users/create", api.AdminCreateUserHandler).Methods("POST")

	// Subject management
	admin.HandleFunc("/subjects", api.AdminListSubjectsHandler).Methods("GET")
	admin.HandleFunc("/subjects/{id}", api.AdminGetSubjectHandler).Methods("GET")
	admin.HandleFunc("/subjects", api.AdminCreateSubjectHandler).Methods("POST")
	admin.HandleFunc("/subjects/{id}", api.AdminUpdateSubjectHandler).Methods("PUT")
	admin.HandleFunc("/subjects/{id}", api.AdminDeleteSubjectHandler).Methods("DELETE")

	// Student routes
	student := r.PathPrefix("/student").Subrouter()
	student.Use(middleware.JWTAuth)
	student.Use(middleware.LoadUserMiddleware)
	student.Use(middleware.RequireRole("student"))
	student.HandleFunc("/subjects", api.StudentListSubjectsHandler).Methods("GET")

	// Teacher routes
	teacher := r.PathPrefix("/teacher").Subrouter()
	teacher.Use(middleware.JWTAuth)
	teacher.Use(middleware.LoadUserMiddleware)
	teacher.Use(middleware.RequireRole("teacher"))
	teacher.HandleFunc("/subjects", api.TeacherListMySubjectsHandler).Methods("GET")
}
