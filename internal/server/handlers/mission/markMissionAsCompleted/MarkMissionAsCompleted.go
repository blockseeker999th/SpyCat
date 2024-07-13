package markmissionascompleted

import (
	"log/slog"
	"net/http"

	"github.com/blockseeker999th/SpyCat/internal/server/response"
	customErr "github.com/blockseeker999th/SpyCat/internal/utils/customErrors"
	"github.com/blockseeker999th/SpyCat/internal/utils/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type MissionCompleter interface {
	MarkMissionAsCompleted(id string) error
}

func New(log *slog.Logger, missionCompleter MissionCompleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.mission.MarkMissionAsCompleted.New"

		log = logger.LogWith(log, op, r)

		id := chi.URLParam(r, "id")

		if err := missionCompleter.MarkMissionAsCompleted(id); err != nil {
			log.Error(customErr.ErrCompleteMission.Error(), logger.Err(err))

			render.JSON(w, r, response.Response{
				Status: http.StatusInternalServerError,
				Error:  customErr.ErrCompleteMission.Error(),
			})

			return
		}

		log.Info("mission marked as completed")

		render.JSON(w, r, response.Response{
			Status: http.StatusNoContent,
		})
	}
}
