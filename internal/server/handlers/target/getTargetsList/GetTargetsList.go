package gettargetslist

import (
	"log/slog"
	"net/http"

	"github.com/blockseeker999th/SpyCat/internal/models"
	"github.com/blockseeker999th/SpyCat/internal/server/response"
	customErr "github.com/blockseeker999th/SpyCat/internal/utils/customErrors"
	"github.com/blockseeker999th/SpyCat/internal/utils/logger"
	"github.com/go-chi/render"
)

type TargetsReceiver interface {
	ListTargets() ([]models.Target, error)
}

func New(log *slog.Logger, targetsReceiver TargetsReceiver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.mission.GetTargetsList.New"

		log = logger.LogWith(log, op, r)

		targets, err := targetsReceiver.ListTargets()
		if err != nil {
			log.Error(customErr.ErrGettingTargetsList.Error(), logger.Err(err))

			render.JSON(w, r, response.Response{
				Status: http.StatusInternalServerError,
				Error:  customErr.ErrGettingTargetsList.Error(),
			})

			return
		}

		log.Info("successfully received targets list")

		render.JSON(w, r, response.Response{
			Status:  http.StatusOK,
			Payload: targets,
		})
	}
}
