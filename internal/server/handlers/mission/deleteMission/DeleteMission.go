package deletemission

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/blockseeker999th/SpyCat/internal/server/response"
	customErr "github.com/blockseeker999th/SpyCat/internal/utils/customErrors"
	"github.com/blockseeker999th/SpyCat/internal/utils/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type MissionDeleter interface {
	DeleteMission(id int) error
}

func New(log *slog.Logger, missionDeleter MissionDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.mission.DeleteMission.New"

		log = logger.LogWith(log, op, r)

		param := chi.URLParam(r, "id")
		id, err := strconv.Atoi(param)
		if err != nil {
			log.Error(customErr.ErrInvalidId.Error(), logger.Err(err))

			render.JSON(w, r, response.Response{
				Status: http.StatusBadRequest,
				Error:  customErr.ErrDeletingMission.Error(),
			})
			return
		}

		if err = missionDeleter.DeleteMission(id); err != nil {
			log.Error(customErr.ErrDeletingMission.Error(), logger.Err(err))

			render.JSON(w, r, response.Response{
				Status: http.StatusInternalServerError,
				Error:  customErr.ErrDeletingMission.Error(),
			})
			return
		}

		log.Info("missing successfully deleted")

		render.JSON(w, r, response.Response{
			Status: http.StatusNoContent,
		})
	}
}
