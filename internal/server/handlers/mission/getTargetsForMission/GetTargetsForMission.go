package gettargetsformission

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

type MissionTargetsReceiver interface {
	ListTargetsForMission(missionID int) ([]models.Target, error)
}

func New(log *slog.Logger, missionTargetsReceiver MissionTargetsReceiver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.mission.GetTargetsForMission.New"

		log = logger.LogWith(log, op, r)

		param := chi.URLParam(r, "mission_id")
		missionID, err := strconv.Atoi(param)
		if err != nil {
			log.Error(customErr.ErrInvalidId.Error(), logger.Err(err))

			render.JSON(w, r, response.Response{
				Status: http.StatusBadRequest,
				Error:  customErr.ErrGettingTargetForMissionList.Error(),
			})

			return
		}

		targets, err := missionTargetsReceiver.ListTargetsForMission(missionID)
		if err != nil {
			log.Error(customErr.ErrGettingTargetForMissionList.Error(), logger.Err(err))

			render.JSON(w, r, response.Response{
				Status: http.StatusInternalServerError,
				Error:  customErr.ErrGettingTargetForMissionList.Error(),
			})

			return
		}

		log.Info("successfully added target to mission")

		render.JSON(w, r, response.Response{
			Status:  http.StatusOK,
			Payload: targets,
		})
	}
}
