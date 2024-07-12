package createmission

import (
	"log/slog"
	"net/http"

	"github.com/blockseeker999th/SpyCat/internal/models"
	"github.com/blockseeker999th/SpyCat/internal/server/response"
	customErr "github.com/blockseeker999th/SpyCat/internal/utils/customErrors"
	"github.com/blockseeker999th/SpyCat/internal/utils/logger"
	"github.com/go-chi/render"
)

type MissionSaver interface {
	CreateMission(mission *models.Mission) error
}

func New(log *slog.Logger, missionSaver MissionSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.mission.CreateMission.New"

		log = logger.LogWith(log, op, r)

		var mission models.Mission
		if err := render.DecodeJSON(r.Body, &mission); err != nil {
			log.Error(customErr.ErrFailedToDecode.Error(), logger.Err(err))

			render.JSON(w, r, response.Response{
				Status: http.StatusBadRequest,
				Error:  customErr.ErrFailedToDecode.Error(),
			})

			return
		}

		if err := missionSaver.CreateMission(&mission); err != nil {
			log.Error(customErr.ErrCreatingMission.Error(), logger.Err(err))

			render.JSON(w, r, response.Response{
				Status: http.StatusInternalServerError,
				Error:  customErr.ErrCreatingMission.Error(),
			})

			return
		}

		log.Info("missing successfully created")

		render.JSON(w, r, response.Response{
			Status:  http.StatusCreated,
			Payload: mission,
		})
	}
}
