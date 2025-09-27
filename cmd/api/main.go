package main

import (
	"log"
	"net/http"
	"os"

	docs "github.com/arman300s/uni-portal/cmd/api/docs"
	"github.com/arman300s/uni-portal/internal/api"
	"github.com/arman300s/uni-portal/internal/models"
	"github.com/arman300s/uni-portal/pkg/db"
	"github.com/arman300s/uni-portal/pkg/middleware"

	_ "github.com/arman300s/uni-portal/cmd/api/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	db.Connect()
	if err := db.DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong 🏓"))
	})

	mux.HandleFunc("/signup", api.SignupHandler)
	mux.HandleFunc("/login", api.LoginHandler)
	mux.Handle("/me", middleware.JWTAuth(http.HandlerFunc(api.MeHandler)))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8079"
	}

	docs.SwaggerInfo.Host = ""
	docs.SwaggerInfo.BasePath = "/"

	mux.Handle("/swagger/", httpSwagger.WrapHandler)
	mux.HandleFunc("/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/index.html", http.StatusMovedPermanently)
	})

	log.Printf("🚀 Server started on port %s\n", port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, mux))
}
