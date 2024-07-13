package models

type SpyCat struct {
	ID                int      `json:"id"`
	Name              string   `json:"name"`
	YearsOfExperience int      `json:"years_of_experience"`
	Breed             string   `json:"breed"`
	Salary            float64  `json:"salary"`
	Mission           *Mission `json:"mission,omitempty"`
}

type Mission struct {
	ID       int      `json:"id"`
	SpyCatID int      `json:"spy_cat_id"`
	Status   string   `json:"status"`
	Targets  []Target `json:"targets"`
}

type Target struct {
	ID        int    `json:"id"`
	MissionID int    `json:"mission_id"`
	Name      string `json:"name"`
	Country   string `json:"country"`
	Notes     string `json:"notes"`
	Completed bool   `json:"completed"`
}
