package getspycatslist

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/render"

	"github.com/blockseeker999th/SpyCat/internal/models"
	"github.com/blockseeker999th/SpyCat/internal/server/response"
	customErr "github.com/blockseeker999th/SpyCat/internal/utils/customErrors"
	"github.com/blockseeker999th/SpyCat/internal/utils/logger"
)

type CatsListReceiver interface {
	ListSpyCats() ([]models.SpyCat, error)
}

func New(log *slog.Logger, catsReceiver CatsListReceiver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.cat.GetSpyCatsList.New"

		log = logger.LogWith(log, op, r)

		cats, err := catsReceiver.ListSpyCats()
		if err != nil {
			log.Error(customErr.ErrGettingSpyCatList.Error(), logger.Err(err))

			render.JSON(w, r, response.Response{
				Status: http.StatusInternalServerError,
				Error:  customErr.ErrGettingSpyCatList.Error(),
			})

			return
		}

		log.Info("spy cats list successfully received")

		render.JSON(w, r, response.Response{
			Status:  http.StatusOK,
			Payload: cats,
		})
	}
}
