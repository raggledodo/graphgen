package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/robfig/cron"

	log "github.com/sirupsen/logrus"
)

const (
	restAPI  = `REST API: 
GET '/'               => REST API message
GET '/getlog'         => Get log level
GET '/setlog/<level>' => Set log level
GET '/stop'           => Stop cronjob`
)

type Controller struct {
	router *mux.Router
	server *http.Server

	mgr *GRPCServer
	job *cron.Cron
}

func NewController(port int) *Controller {
	router := mux.NewRouter().StrictSlash(true)
	server := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: router}
	ctrl := &Controller{
		router: router,
		server: server,
	}

	router.Path("/").
	HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintf(w, restAPI)
	})
	router.Path("/getlog").
	HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, log.GetLevel().String())
	})
	router.Path("/setlog").Queries("level", "{level}").
	HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		level := r.FormValue("level")
		logLevel, err := log.ParseLevel(level)
		if err != nil {
			fmt.Fprintf(w, "Log level parsing error: %v", err)
		} else {
			log.SetLevel(logLevel)
			fmt.Fprintf(w, "Set log level to %s", level)
		}
	})
	router.Path("/stop").
	HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		ctrl.Stop()
	})

	return ctrl
}

func (ctrl Controller) Start(wg sync.WaitGroup) {
	defer wg.Done()
	log.Info("Controller server running")
	err := ctrl.server.ListenAndServe();
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("Controller server error: %v", err)
	}
	log.Info("Controller server stopped")
}

func (ctrl Controller) Stop() {
	log.Info("Controller server stopping")
	ctrl.server.Close()

	if ctrl.job != nil {
		ctrl.job.Stop()
	}
	if ctrl.mgr != nil {
		ctrl.mgr.Stop()
	}
}
