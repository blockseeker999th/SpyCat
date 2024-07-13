package gettarget

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

type TargetReceiver interface {
	GetTarget(id int) (models.Target, error)
}

func New(log *slog.Logger, targetReceiver TargetReceiver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.mission.GetTarget.New"

		log = logger.LogWith(log, op, r)

		param := chi.URLParam(r, "id")
		id, err := strconv.Atoi(param)
		if err != nil {
			log.Error(customErr.ErrInvalidId.Error(), logger.Err(err))

			render.JSON(w, r, response.Response{
				Status: http.StatusBadRequest,
				Error:  customErr.ErrGettingTarget.Error(),
			})

			return
		}

		target, err := targetReceiver.GetTarget(id)
		if err != nil {
			log.Error(customErr.ErrGettingTarget.Error(), logger.Err(err))

			render.JSON(w, r, response.Response{
				Status: http.StatusNotFound,
				Error:  customErr.ErrGettingTarget.Error(),
			})

			return
		}

		log.Info("successfully received target")

		render.JSON(w, r, response.Response{
			Status:  http.StatusOK,
			Payload: target,
		})
	}
}
