package getmissionslist

import (
	"log/slog"
	"net/http"

	"github.com/blockseeker999th/SpyCat/internal/models"
	"github.com/blockseeker999th/SpyCat/internal/server/response"
	customErr "github.com/blockseeker999th/SpyCat/internal/utils/customErrors"
	"github.com/blockseeker999th/SpyCat/internal/utils/logger"
	"github.com/go-chi/render"
)

type MissionsReceiver interface {
	ListMissions() ([]models.Mission, error)
}

func New(log *slog.Logger, missionsReceiver MissionsReceiver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.mission.GetMissionsList.New"

		log = logger.LogWith(log, op, r)

		missions, err := missionsReceiver.ListMissions()
		if err != nil {
			log.Error(customErr.ErrGettingMissionsList.Error(), logger.Err(err))

			render.JSON(w, r, response.Response{
				Status: http.StatusInternalServerError,
				Error:  customErr.ErrGettingMissionsList.Error(),
			})

			return
		}

		log.Info("missions list successfully received")

		render.JSON(w, r, response.Response{
			Status:  http.StatusOK,
			Payload: missions,
		})
	}
}
