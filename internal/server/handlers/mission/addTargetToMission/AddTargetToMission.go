package addtargettomission

import (
	"log/slog"
	"net/http"

	"github.com/blockseeker999th/SpyCat/internal/models"
	"github.com/blockseeker999th/SpyCat/internal/server/response"
	customErr "github.com/blockseeker999th/SpyCat/internal/utils/customErrors"
	"github.com/blockseeker999th/SpyCat/internal/utils/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type MissionTargeter interface {
	AddTargetToMission(missionID string, target *models.Target) error
}

func New(log *slog.Logger, missionTargeter MissionTargeter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.mission.AddTargetToMission.New"

		log = logger.LogWith(log, op, r)

		missionID := chi.URLParam(r, "mission_id")

		var target models.Target
		if err := render.DecodeJSON(r.Body, &target); err != nil {
			log.Error(customErr.ErrFailedToDecode.Error(), logger.Err(err))

			render.JSON(w, r, response.Response{
				Status: http.StatusBadRequest,
				Error:  customErr.ErrFailedToDecode.Error(),
			})

			return
		}

		if err := missionTargeter.AddTargetToMission(missionID, &target); err != nil {
			log.Error(customErr.ErrAddingTargetToMission.Error(), logger.Err(err))

			render.JSON(w, r, response.Response{
				Status: http.StatusInternalServerError,
				Error:  customErr.ErrAddingTargetToMission.Error(),
			})

			return
		}

		log.Info("successfully added target to mission")

		render.JSON(w, r, response.Response{
			Status:  http.StatusCreated,
			Payload: target,
		})
	}
}
