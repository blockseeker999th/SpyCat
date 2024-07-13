package getmission

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/blockseeker999th/SpyCat/internal/models"
	"github.com/blockseeker999th/SpyCat/internal/server/response"
	customErr "github.com/blockseeker999th/SpyCat/internal/utils/customErrors"
	"github.com/blockseeker999th/SpyCat/internal/utils/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type MissionReceiver interface {
	GetMission(id int) (models.Mission, error)
}

func New(log *slog.Logger, missionReceiver MissionReceiver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.mission.GetMission.New"

		log = logger.LogWith(log, op, r)

		param := chi.URLParam(r, "id")

		id, err := strconv.Atoi(param)
		if err != nil {
			log.Error(customErr.ErrInvalidId.Error(), logger.Err(err))

			render.JSON(w, r, response.Response{
				Status: http.StatusBadRequest,
				Error:  customErr.ErrGettingMission.Error(),
			})

			return
		}

		mission, err := missionReceiver.GetMission(id)
		if err != nil {
			log.Error(customErr.ErrGettingMission.Error(), logger.Err(err))

			render.JSON(w, r, response.Response{
				Status: http.StatusNotFound,
				Error:  customErr.ErrGettingMission.Error(),
			})

			return
		}

		log.Info("mission successfully received")

		render.JSON(w, r, response.Response{
			Status:  http.StatusOK,
			Payload: mission,
		})
	}
}
