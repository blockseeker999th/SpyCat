package createtarget

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

type TargetSaver interface {
	CreateTarget(id int, target *models.Target) error
}

func New(log *slog.Logger, targetSaver TargetSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.mission.CreateTarget.New"

		log = logger.LogWith(log, op, r)

		param := chi.URLParam(r, "mission_id")
		id, err := strconv.Atoi(param)
		if err != nil {
			log.Error(customErr.ErrInvalidId.Error(), logger.Err(err))

			render.JSON(w, r, response.Response{
				Status: http.StatusBadRequest,
				Error:  customErr.ErrCreatingTarget.Error(),
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

		if err := targetSaver.CreateTarget(id, &target); err != nil {
			log.Error(customErr.ErrCreatingTarget.Error(), logger.Err(err))

			render.JSON(w, r, response.Response{
				Status: http.StatusInternalServerError,
				Error:  customErr.ErrCreatingTarget.Error(),
			})

			return
		}

		log.Info("successfully created a targer")

		render.JSON(w, r, response.Response{
			Status:  http.StatusCreated,
			Payload: target,
		})
	}
}
