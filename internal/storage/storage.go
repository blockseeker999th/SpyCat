package storage

import (
	"database/sql"
	"errors"

	"github.com/blockseeker999th/SpyCat/internal/models"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) CreateSpyCat(cat *models.SpyCat) error {
	err := s.db.QueryRow(
		"INSERT INTO spy_cats (name, years_of_experience, breed, salary) VALUES ($1, $2, $3, $4) RETURNING id",
		cat.Name, cat.YearsOfExperience, cat.Breed, cat.Salary,
	).Scan(&cat.ID)
	return err
}

func (s *Storage) ListSpyCats() ([]models.SpyCat, error) {
	rows, err := s.db.Query("SELECT id, name, years_of_experience, breed, salary FROM spy_cats")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []models.SpyCat
	for rows.Next() {
		var cat models.SpyCat
		if err := rows.Scan(&cat.ID, &cat.Name, &cat.YearsOfExperience, &cat.Breed, &cat.Salary); err != nil {
			return nil, err
		}
		cats = append(cats, cat)
	}
	return cats, nil
}

func (s *Storage) GetSpyCat(id int) (models.SpyCat, error) {
	var cat models.SpyCat
	err := s.db.QueryRow(
		"SELECT id, name, years_of_experience, breed, salary FROM spy_cats WHERE id=$1",
		id,
	).Scan(&cat.ID, &cat.Name, &cat.YearsOfExperience, &cat.Breed, &cat.Salary)
	if err != nil {
		return cat, err
	}

	var missionID int
	err = s.db.QueryRow("SELECT id FROM missions WHERE spy_cat_id=$1 AND status='in_progress'", id).Scan(&missionID)
	if err == nil && missionID > 0 {

		mission, err := s.GetMission(missionID)
		if err == nil {
			cat.Mission = &mission
		}
	}

	return cat, nil
}

func (s *Storage) UpdateSpyCat(id int, salary *float64) error {
	_, err := s.db.Exec(
		"UPDATE spy_cats SET salary=$1 WHERE id=$2",
		&salary, id,
	)
	return err
}

func (s *Storage) DeleteSpyCat(id int) error {
	_, err := s.db.Exec("DELETE FROM spy_cats WHERE id=$1", id)
	return err
}

func (s *Storage) CreateMission(mission *models.Mission) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	var count int
	err = s.db.QueryRow(
		"SELECT COUNT(*) FROM missions WHERE spy_cat_id=$1 AND status='in_progress'",
		mission.SpyCatID,
	).Scan(&count)
	if err != nil {
		tx.Rollback()
		return err
	}
	if count > 0 {
		tx.Rollback()
		return errors.New("the spy cat already has an ongoing mission")
	}

	err = tx.QueryRow(
		"INSERT INTO missions (spy_cat_id, status) VALUES ($1, $2) RETURNING id",
		mission.SpyCatID, mission.Status,
	).Scan(&mission.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if len(mission.Targets) < 1 || len(mission.Targets) > 3 {
		tx.Rollback()
		return errors.New("a mission must have between 1 and 3 targets")
	}

	for idx, target := range mission.Targets {
		err := tx.QueryRow(
			"INSERT INTO targets (mission_id, name, country, notes, completed) VALUES ($1, $2, $3, $4, $5) RETURNING id",
			mission.ID, target.Name, target.Country, target.Notes, target.Completed,
		).Scan(&mission.Targets[idx].ID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (s *Storage) ListMissions() ([]models.Mission, error) {
	rows, err := s.db.Query("SELECT id, spy_cat_id, status FROM missions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var missions []models.Mission
	for rows.Next() {
		var mission models.Mission
		if err := rows.Scan(&mission.ID, &mission.SpyCatID, &mission.Status); err != nil {
			return nil, err
		}

		targets, err := s.ListTargetsForMission(mission.ID)
		if err != nil {
			return nil, err
		}
		mission.Targets = targets

		missions = append(missions, mission)
	}

	return missions, nil
}

func (s *Storage) GetMission(id int) (models.Mission, error) {
	var mission models.Mission
	err := s.db.QueryRow(
		"SELECT id, spy_cat_id, status FROM missions WHERE id=$1",
		id,
	).Scan(&mission.ID, &mission.SpyCatID, &mission.Status)
	if err != nil {
		return mission, err
	}

	targets, err := s.ListTargetsForMission(mission.ID)
	if err != nil {
		return mission, err
	}
	mission.Targets = targets

	return mission, nil
}

func (s *Storage) UpdateMission(id int, mission *models.Mission) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	var targetCount int
	err = tx.QueryRow("SELECT COUNT(*) FROM targets WHERE mission_id=$1", id).Scan(&targetCount)
	if err != nil {
		tx.Rollback()
		return err
	}

	if targetCount < 1 || targetCount > 3 {
		tx.Rollback()
		return errors.New("a mission must have between 1 and 3 targets")
	}

	_, err = tx.Exec(
		"UPDATE missions SET spy_cat_id=$1, status=$2 WHERE id=$3",
		mission.SpyCatID, mission.Status, id,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, target := range mission.Targets {
		_, err := tx.Exec(
			"UPDATE targets SET name=$1, country=$2, notes=$3, completed=$4 WHERE id=$5 AND mission_id=$6",
			target.Name, target.Country, target.Notes, target.Completed, target.ID, id,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	var completedTargetsCount int
	err = s.db.QueryRow(
		"SELECT COUNT(*) FROM targets WHERE mission_id=$1 AND completed=false",
		id,
	).Scan(&completedTargetsCount)
	if err != nil {
		tx.Rollback()
		return err
	}

	if completedTargetsCount == 0 {
		_, err := tx.Exec(
			"UPDATE missions SET status='completed' WHERE id=$1",
			id,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (s *Storage) ListTargetsForMission(missionID int) ([]models.Target, error) {
	rows, err := s.db.Query("SELECT id, mission_id, name, country, notes, completed FROM targets WHERE mission_id=$1", missionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var targets []models.Target
	for rows.Next() {
		var target models.Target
		if err := rows.Scan(&target.ID, &target.MissionID, &target.Name, &target.Country, &target.Notes, &target.Completed); err != nil {
			return nil, err
		}
		targets = append(targets, target)
	}

	return targets, nil
}

func (s *Storage) CreateTarget(id int, target *models.Target) error {
	var missionStatus string
	err := s.db.QueryRow(
		"SELECT status FROM missions WHERE id=$1",
		id,
	).Scan(&missionStatus)
	if err != nil {
		return err
	}

	if missionStatus == "completed" {
		return errors.New("cannot create more targets after the missions was completed")
	}

	var targetCount int
	err = s.db.QueryRow("SELECT COUNT(*) FROM targets WHERE mission_id=$1", id).Scan(&targetCount)
	if err != nil {
		return err
	}

	if targetCount < 1 || targetCount >= 3 {
		return errors.New("a mission must have between 1 and 3 targets")
	}

	_, err = s.db.Exec(
		"INSERT INTO targets (mission_id, name, country, notes, completed) VALUES ($1, $2, $3, $4, $5)",
		id, target.Name, target.Country, target.Notes, target.Completed,
	)
	return err
}

func (s *Storage) ListTargets() ([]models.Target, error) {
	rows, err := s.db.Query("SELECT id, mission_id, name, notes, completed FROM targets")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var targets []models.Target
	for rows.Next() {
		var target models.Target
		if err := rows.Scan(&target.ID, &target.MissionID, &target.Name, &target.Notes, &target.Completed); err != nil {
			return nil, err
		}
		targets = append(targets, target)
	}

	return targets, nil
}

func (s *Storage) GetTarget(id int) (models.Target, error) {
	var target models.Target
	err := s.db.QueryRow(
		"SELECT id, mission_id, name, notes, completed FROM targets WHERE id=$1",
		id,
	).Scan(&target.ID, &target.MissionID, &target.Name, &target.Notes, &target.Completed)
	if err != nil {
		return target, err
	}
	return target, nil
}

func (s *Storage) UpdateTarget(tId int, mId int, target *models.Target) error {
	var currentCompleted bool
	err := s.db.QueryRow(
		"SELECT completed FROM targets WHERE id=$1",
		tId,
	).Scan(&currentCompleted)
	if err != nil {
		return err
	}

	var currentNotes string
	err = s.db.QueryRow(
		"SELECT notes FROM targets WHERE id=$1",
		tId,
	).Scan(&currentNotes)
	if err != nil {
		return err
	}

	if currentCompleted && target.Notes != currentNotes {
		return errors.New("cannot update notes after target is completed")
	}
	_, err = s.db.Exec(
		"UPDATE targets SET name=$1, country=$2, notes=$3, completed=$4 WHERE id=$5",
		target.Name, target.Country, target.Notes, target.Completed, tId,
	)
	if err != nil {
		return err
	}

	var completedTargetsCount int

	err = s.db.QueryRow(
		"SELECT COUNT(*) FROM targets WHERE mission_id=$1 AND completed=false",
		mId,
	).Scan(&completedTargetsCount)
	if err != nil {
		return err
	}

	if completedTargetsCount == 0 {
		_, err := s.db.Exec(
			"UPDATE missions SET status='completed' WHERE id=$1",
			mId,
		)
		if err != nil {
			return err
		}
	}

	return err
}

func (s *Storage) DeleteTarget(tId int, mId int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	var currentCompleted bool
	err = tx.QueryRow(
		"SELECT completed FROM targets WHERE id=$1",
		tId,
	).Scan(&currentCompleted)
	if err != nil {
		tx.Rollback()
		return err
	}

	if currentCompleted {
		tx.Rollback()
		return errors.New("cannot delete target after it is completed")
	}

	var targetCount int
	err = tx.QueryRow("SELECT COUNT(*) FROM targets WHERE mission_id=$1", mId).Scan(&targetCount)
	if err != nil {
		tx.Rollback()
		return err
	}

	if targetCount <= 1 || targetCount > 3 {
		tx.Rollback()
		return errors.New("a mission must have between 1 and 3 targets")
	}

	_, err = s.db.Exec("DELETE FROM targets WHERE id=$1", tId)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Storage) MarkMissionAsCompleted(missionID string) error {
	_, err := s.db.Exec(
		"UPDATE missions SET status='completed' WHERE id=$1",
		missionID,
	)
	return err
}

func (s *Storage) DeleteMission(id int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	var status string

	err = tx.QueryRow("SELECT status FROM missions WHERE id=$1", id).Scan(&status)
	if err != nil {
		tx.Rollback()
		return err
	}

	if status != "completed" {
		tx.Rollback()
		return errors.New("mission must have 'completed' status to be deleted")
	}

	_, err = tx.Exec("DELETE FROM targets WHERE mission_id=$1", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("DELETE FROM missions WHERE id=$1", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
