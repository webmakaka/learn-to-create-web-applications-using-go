package main

import (
	"fmt"
	"net/http"

	"./controllers"
	"./middleware"
	"./models"

	"github.com/gorilla/mux"
)

const (
	host     = "ec2-23-23-184-76.compute-1.amazonaws.com"
	port     = 5432
	user     = "hrgcmhzjkgllyf"
	password = "f867d132e78e27e50a27d0b7522dbf3f44dc835c903eb3040d74ecd5daf5c633"
	dbname   = "d61hvpjfrp6em7"
	sslmode  = "require"
)

func main() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)

	services, err := models.NewServices(psqlInfo)

	must(err)

	defer services.Close()
	// services.DestructiveReset()
	services.AutoMigrate()

	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(services.User)
	galleriesC := controllers.NewGalleries(services.Gallery)

	requireUserMw := middleware.RequireUser{
		UserService: services.User,
	}

	r := mux.NewRouter()
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")

	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")

	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")

	r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")

	r.Handle("/galleries/new", requireUserMw.Apply(galleriesC.New)).Methods("GET")
	r.HandleFunc("/galleries", requireUserMw.ApplyFn(galleriesC.Create)).Methods("POST")

	http.ListenAndServe(":3000", r)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
