package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"

	"github.com/blockseeker999th/SpyCat/internal/config"
	"github.com/blockseeker999th/SpyCat/internal/db"
	createspycat "github.com/blockseeker999th/SpyCat/internal/server/handlers/cat/createSpyCat"
	deletespycat "github.com/blockseeker999th/SpyCat/internal/server/handlers/cat/deleteSpyCat"
	getspycat "github.com/blockseeker999th/SpyCat/internal/server/handlers/cat/getSpyCat"
	getspycatslist "github.com/blockseeker999th/SpyCat/internal/server/handlers/cat/getSpyCatsList"
	updatespycat "github.com/blockseeker999th/SpyCat/internal/server/handlers/cat/updateSpyCat"
	createmission "github.com/blockseeker999th/SpyCat/internal/server/handlers/mission/createMission"
	deletemission "github.com/blockseeker999th/SpyCat/internal/server/handlers/mission/deleteMission"
	getmission "github.com/blockseeker999th/SpyCat/internal/server/handlers/mission/getMission"
	getmissionslist "github.com/blockseeker999th/SpyCat/internal/server/handlers/mission/getMissionsList"
	gettargetsformission "github.com/blockseeker999th/SpyCat/internal/server/handlers/mission/getTargetsForMission"
	markmissionascompleted "github.com/blockseeker999th/SpyCat/internal/server/handlers/mission/markMissionAsCompleted"
	updatemission "github.com/blockseeker999th/SpyCat/internal/server/handlers/mission/updateMission"
	createtarget "github.com/blockseeker999th/SpyCat/internal/server/handlers/target/createTarget"
	deletetarget "github.com/blockseeker999th/SpyCat/internal/server/handlers/target/deleteTarget"
	updatetarget "github.com/blockseeker999th/SpyCat/internal/server/handlers/target/updateTarget"
	mwLogger "github.com/blockseeker999th/SpyCat/internal/server/middleware/logger"
	"github.com/blockseeker999th/SpyCat/internal/storage"
	"github.com/blockseeker999th/SpyCat/internal/utils/logger"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Errorf("error loading env variables: %s", err.Error())
	}

	cfg := config.LoadConfig()

	log := logger.SetupLogger(os.Getenv("ENV"))
	log.Info("starting SCA", slog.String("env", os.Getenv("ENV")))

	db, err := db.ConnectDB(cfg)
	if err != nil {
		log.Error("failed connect to DB", logger.Err(err))
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Error("Error closing DB: ", logger.Err(err))
		}
	}()

	st := storage.NewStorage(db)

	router := chi.NewRouter()

	router.Use(mwLogger.New(log))
	router.Use(middleware.RequestID)

	router.Route("/spycats", func(r chi.Router) {
		r.Post("/", createspycat.New(log, st))
		r.Get("/", getspycatslist.New(log, st))
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", getspycat.New(log, st))
			r.Put("/", updatespycat.New(log, st))
			r.Delete("/", deletespycat.New(log, st))
		})
	})

	router.Route("/missions", func(r chi.Router) {
		r.Post("/", createmission.New(log, st))
		r.Get("/", getmissionslist.New(log, st))
		r.Get("/{id}", getmission.New(log, st))
		r.Put("/{id}", updatemission.New(log, st))
		r.Patch("/{id}/complete", markmissionascompleted.New(log, st))
		r.Delete("/{id}", deletemission.New(log, st))

		r.Route("/{mission_id}/targets", func(r chi.Router) {
			r.Post("/", createtarget.New(log, st))
			r.Get("/", gettargetsformission.New(log, st))
			r.Put("/{id}", updatetarget.New(log, st))
			r.Delete("/{id}", deletetarget.New(log, st))
		})
	})

	log.Info("starting server")

	srv := http.Server{
		Addr:         cfg.HTTPServer.Address,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		Handler:      router,
	}

	srv.ListenAndServe()

	log.Error("server stopped")
}
