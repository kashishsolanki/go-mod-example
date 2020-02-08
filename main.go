package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kashishsolanki/dt-hrms-golang/db"
)

var (
	ctx    context.Context
	secret = []byte(os.Getenv("SECRET_HRMS_KEY")) // []byte("SECRET_HRMS_KEY")
)

func main() {
	fmt.Println("DT HRMS system started")

	// Context to cancel functions
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(context.Background())

	db, err := db.Open()
	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	fmt.Println("Successfully connected to mysql database")

	router := Router()

	srv := &http.Server{
		Addr: ":8080",
		// Pass our instance of gorilla/mux in
		Handler: &MyServer{router},
	}

	// appengine.Main()

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C) or SIGKILL
	// SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	sigquit := make(chan os.Signal, 1)
	signal.Notify(sigquit, os.Kill, os.Interrupt, syscall.SIGTERM)
	sig := <-sigquit
	fmt.Printf("\ncaught sig: %+v\n", sig)
	fmt.Printf("Gracefully shutting down server...\n")
	cancel()

	srv.Shutdown(ctx)
}
