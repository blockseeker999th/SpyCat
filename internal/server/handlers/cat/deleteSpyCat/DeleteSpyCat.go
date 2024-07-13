package deletespycat

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

type CatRemover interface {
	DeleteSpyCat(id int) error
}

func New(log *slog.Logger, catRemover CatRemover) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.cat.DeleteSpyCat.New"

		log = logger.LogWith(log, op, r)

		param := chi.URLParam(r, "id")

		id, err := strconv.Atoi(param)
		if err != nil {
			log.Error(customErr.ErrInvalidId.Error(), logger.Err(err))

			render.JSON(w, r, response.Response{
				Status: http.StatusBadRequest,
				Error:  customErr.ErrDeletingSpyCat.Error(),
			})
			return
		}

		if err := catRemover.DeleteSpyCat(id); err != nil {
			log.Error(customErr.ErrDeletingSpyCat.Error(), logger.Err(err))

			render.JSON(w, r, response.Response{
				Status: http.StatusInternalServerError,
				Error:  customErr.ErrDeletingSpyCat.Error(),
			})

			return
		}

		log.Info("spy cat successfully deleted")

		render.JSON(w, r, response.Response{
			Status: http.StatusNoContent,
		})
	}
}
