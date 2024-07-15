package updatespycat

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"github.com/blockseeker999th/SpyCat/internal/models"
	"github.com/blockseeker999th/SpyCat/internal/server/response"
	customErr "github.com/blockseeker999th/SpyCat/internal/utils/customErrors"
	"github.com/blockseeker999th/SpyCat/internal/utils/logger"
)

type CatUpdater interface {
	UpdateSpyCat(id int, salary *float64) error
}

func New(log *slog.Logger, catUpdater CatUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.cat.UpdateSpyCat.New"

		log = logger.LogWith(log, op, r)

		param := chi.URLParam(r, "id")

		id, err := strconv.Atoi(param)
		if err != nil {
			log.Error(customErr.ErrInvalidId.Error(), logger.Err(err))

			render.JSON(w, r, response.Response{
				Status: http.StatusBadRequest,
				Error:  customErr.ErrUpdatingSpyCat.Error(),
			})
			return
		}

		var cat models.SpyCat
		if err := render.DecodeJSON(r.Body, &cat); err != nil {
			log.Error(customErr.ErrFailedToDecode.Error(), logger.Err(err))

			render.JSON(w, r, response.Response{
				Status: http.StatusBadRequest,
				Error:  customErr.ErrFailedToDecode.Error(),
			})

			return
		}

		if err := catUpdater.UpdateSpyCat(id, &cat.Salary); err != nil {
			render.JSON(w, r, response.Response{
				Status: http.StatusInternalServerError,
				Error:  customErr.ErrUpdatingSpyCat.Error(),
			})

			return
		}

		log.Info("spy cat salary successfully updated")

		render.JSON(w, r, response.Response{
			Status:  http.StatusNoContent,
			Payload: cat.Salary,
		})
	}
}
