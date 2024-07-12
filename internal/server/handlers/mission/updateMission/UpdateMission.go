package updatemission

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

type MissionUpdater interface {
	UpdateMission(id int, mission *models.Mission) error
}

func New(log *slog.Logger, missionUpdater MissionUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.mission.UpdateMission.New"

		log = logger.LogWith(log, op, r)

		param := chi.URLParam(r, "id")

		id, err := strconv.Atoi(param)
		if err != nil {
			log.Error(customErr.ErrInvalidId.Error(), logger.Err(err))

			render.JSON(w, r, response.Response{
				Status: http.StatusBadRequest,
				Error:  customErr.ErrUpdatingMission.Error(),
			})

			return
		}

		var mission models.Mission
		if err := render.DecodeJSON(r.Body, &mission); err != nil {
			log.Error(customErr.ErrFailedToDecode.Error(), logger.Err(err))

			render.JSON(w, r, response.Response{
				Status: http.StatusBadRequest,
				Error:  customErr.ErrFailedToDecode.Error(),
			})

			return
		}

		if err := missionUpdater.UpdateMission(id, &mission); err != nil {
			log.Error(customErr.ErrUpdatingMission.Error(), logger.Err(err))

			render.JSON(w, r, response.Response{
				Status: http.StatusInternalServerError,
				Error:  customErr.ErrUpdatingMission.Error(),
			})

			return
		}

		log.Info("mission successfully updated")

		render.JSON(w, r, response.Response{
			Status: http.StatusNoContent,
		})
	}
}
