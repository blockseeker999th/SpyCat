package deletetarget

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

type TargetDeleter interface {
	DeleteTarget(tId int, mId int) error
}

func New(log *slog.Logger, targetDeleter TargetDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.mission.DeleteTarget.New"

		log = logger.LogWith(log, op, r)

		param := chi.URLParam(r, "id")
		tId, err := strconv.Atoi(param)
		if err != nil {
			log.Error(customErr.ErrInvalidId.Error(), logger.Err(err))

			render.JSON(w, r, response.Response{
				Status: http.StatusBadRequest,
				Error:  customErr.ErrCreatingTarget.Error(),
			})

			return
		}

		param = chi.URLParam(r, "mission_id")
		mId, err := strconv.Atoi(param)
		if err != nil {
			log.Error(customErr.ErrInvalidId.Error(), logger.Err(err))

			render.JSON(w, r, response.Response{
				Status: http.StatusBadRequest,
				Error:  customErr.ErrCreatingTarget.Error(),
			})

			return
		}

		if err := targetDeleter.DeleteTarget(tId, mId); err != nil {
			log.Error(customErr.ErrDeletingTarget.Error(), logger.Err(err))

			render.JSON(w, r, response.Response{
				Status: http.StatusInternalServerError,
				Error:  customErr.ErrDeletingTarget.Error(),
			})

			return
		}

		log.Info("successfully deleted the target")

		render.JSON(w, r, response.Response{
			Status: http.StatusNoContent,
		})
	}
}
