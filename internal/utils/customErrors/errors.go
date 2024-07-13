package customErrors

import "errors"

var (
	ErrFailedToDecode              = errors.New("failed to decode request")
	ErrInvalidBreed                = errors.New("invalid breed")
	ErrCreatingSpyCat              = errors.New("failed to create new spy cat")
	ErrGettingSpyCatList           = errors.New("failerd to get list of spy cats")
	ErrGettingSpyCat               = errors.New("failet to get spy cat")
	ErrUpdatingSpyCat              = errors.New("failed to update a spy cat")
	ErrDeletingSpyCat              = errors.New("failed to delete a spy cat")
	ErrCreatingMission             = errors.New("failed to create a mission")
	ErrGettingMissionsList         = errors.New("failed to get missions list")
	ErrGettingMission              = errors.New("failed to get mission")
	ErrUpdatingMission             = errors.New("failed to update the mission")
	ErrCompleteMission             = errors.New("failed to complete a mission")
	ErrAddingTargetToMission       = errors.New("failed to add target to mission")
	ErrGettingTargetForMissionList = errors.New("failed to get targets for mission list")
	ErrCreatingTarget              = errors.New("failed to create a target")
	ErrGettingTargetsList          = errors.New("failed to get targets list")
	ErrGettingTarget               = errors.New("failed to get a target")
	ErrUpdatingTarget              = errors.New("failed to update the target")
	ErrDeletingTarget              = errors.New("failed to delete the target")
	ErrDeletingMission             = errors.New("failed to delete a mission")
	ErrInvalidId                   = errors.New("invalid Id in parameters")
)
