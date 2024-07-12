package getspycat

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

type CatReceiver interface {
	GetSpyCat(id int) (models.SpyCat, error)
}

func New(log *slog.Logger, catReceiver CatReceiver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.cat.GetSpyCat.New"

		log = logger.LogWith(log, op, r)

		param := chi.URLParam(r, "id")

		id, err := strconv.Atoi(param)
		if err != nil {
			log.Error(customErr.ErrInvalidId.Error(), logger.Err(err))

			render.JSON(w, r, response.Response{
				Status: http.StatusBadRequest,
				Error:  customErr.ErrGettingSpyCat.Error(),
			})
			return
		}
		cat, err := catReceiver.GetSpyCat(id)
		if err != nil {
			log.Error(customErr.ErrGettingSpyCat.Error(), logger.Err(err))

			render.JSON(w, r, response.Response{
				Status: http.StatusNotFound,
				Error:  customErr.ErrGettingSpyCat.Error(),
			})

			return
		}

		log.Info("spy cat successfully received")

		render.JSON(w, r, response.Response{
			Status:  http.StatusOK,
			Payload: cat,
		})
	}
}
