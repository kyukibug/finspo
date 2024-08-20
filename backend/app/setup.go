package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"com.fukubox/config"
	"com.fukubox/database" // Import the package that contains the StartDB function
	"com.fukubox/router"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func SetupAndRunApp() error {
	// load env
	err := config.LoadENV()
	if err != nil {
		return err
	}

	// start database
	err = database.StartDB()
	if err != nil {
		return err
	}

	fmt.Println("Database connected")
	defer database.CloseDB()

	r := chi.NewRouter()
	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	router.SetupAuthenticatedRoutes(r)

	// get the port and start
	port := os.Getenv("PORT")

	err = http.ListenAndServe(":"+port, r)

	if err != nil {
		log.Printf("Failed to launch api server:%+v\n", err)
	}

	return nil
}
