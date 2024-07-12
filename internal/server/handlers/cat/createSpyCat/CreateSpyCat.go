package createspycat

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/render"

	"github.com/blockseeker999th/SpyCat/internal/models"
	"github.com/blockseeker999th/SpyCat/internal/server/response"
	customErr "github.com/blockseeker999th/SpyCat/internal/utils/customErrors"
	"github.com/blockseeker999th/SpyCat/internal/utils/logger"
	"github.com/blockseeker999th/SpyCat/internal/validation"
)

type CatSaver interface {
	CreateSpyCat(cat *models.SpyCat) error
}

func New(log *slog.Logger, catSaver CatSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.cat.CreateSpyCat.New"

		log = logger.LogWith(log, op, r)

		var newCat models.SpyCat
		if err := render.DecodeJSON(r.Body, &newCat); err != nil {
			log.Error(customErr.ErrFailedToDecode.Error(), logger.Err(err))

			render.JSON(w, r, response.Response{
				Status: http.StatusBadRequest,
				Error:  customErr.ErrFailedToDecode.Error(),
			})

			return
		}

		if !validation.ValidateBreed(newCat.Breed) {
			log.Error(customErr.ErrInvalidBreed.Error())

			render.JSON(w, r, response.Response{
				Status: http.StatusBadRequest,
				Error:  customErr.ErrInvalidBreed.Error(),
			})

			return
		}

		if err := catSaver.CreateSpyCat(&newCat); err != nil {
			log.Error(customErr.ErrCreatingSpyCat.Error(), logger.Err(err))

			render.JSON(w, r, response.Response{
				Status: http.StatusInternalServerError,
				Error:  customErr.ErrCreatingSpyCat.Error(),
			})

			return
		}

		log.Info("spy cat successfully created", slog.Any("spy cat: ", newCat))

		render.JSON(w, r, response.Response{
			Status:  http.StatusCreated,
			Payload: newCat,
		})
	}
}
