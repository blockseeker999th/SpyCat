package updatetarget

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

type TargetUpdater interface {
	UpdateTarget(tId int, mId int, target *models.Target) error
}

func New(log *slog.Logger, targetUpdater TargetUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.mission.UpdateTarget.New"

		log = logger.LogWith(log, op, r)

		targetId := chi.URLParam(r, "id")
		missionId := chi.URLParam(r, "mission_id")
		tId, err := strconv.Atoi(targetId)
		if err != nil {
			log.Error(customErr.ErrInvalidId.Error(), logger.Err(err))

			render.JSON(w, r, response.Response{
				Status: http.StatusBadRequest,
				Error:  customErr.ErrUpdatingTarget.Error(),
			})

			return
		}

		mId, err := strconv.Atoi(missionId)
		if err != nil {
			log.Error(customErr.ErrInvalidId.Error(), logger.Err(err))

			render.JSON(w, r, response.Response{
				Status: http.StatusBadRequest,
				Error:  customErr.ErrUpdatingTarget.Error(),
			})

			return
		}

		var target models.Target
		if err := render.DecodeJSON(r.Body, &target); err != nil {
			log.Error(customErr.ErrFailedToDecode.Error(), logger.Err(err))

			render.JSON(w, r, response.Response{
				Status: http.StatusBadRequest,
				Error:  customErr.ErrFailedToDecode.Error(),
			})

			return
		}

		if err := targetUpdater.UpdateTarget(tId, mId, &target); err != nil {
			log.Error(customErr.ErrUpdatingTarget.Error(), logger.Err(err))

			render.JSON(w, r, response.Response{
				Status: http.StatusInternalServerError,
				Error:  customErr.ErrUpdatingTarget.Error(),
			})

			return
		}

		log.Info("successfully received the target")

		render.JSON(w, r, response.Response{
			Status: http.StatusNoContent,
		})
	}
}
