package responses

import (
	"time"
)

type Link struct {
	Label string `json:"label"`
	URL   string `json:"url"`
}

type Sketch struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	FilePath string `json:"file_path"`
}

type CriteriaRate struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Rate int    `json:"rate"`
}

type IdeaLightData struct {
	ID          string          `json:"id"`
	Worker      WorkerLightData `json:"worker"`
	Name        string          `json:"game_name"`
	Date        time.Time       `json:"date"`
	IsItNew     bool            `json:"is_it_new"`
	Description string          `json:"desription"`
	FilePath    *string         `json:"file_path"`
	OverallRate int             `json:"rate"`
}

type IdeaSpecData struct {
	ID            string          `json:"id"`
	Worker        WorkerLightData `json:"worker"`
	Name          string          `json:"name"`
	Date          time.Time       `json:"date"`
	Description   string          `json:"description"`
	Genre         string          `json:"genre"`
	Mechanics     []string        `json:"mechanics"`
	Links         []Link          `json:"links"`
	FilePaths     []Sketch        `json:"file_paths"`
	CriteriaRates []CriteriaRate  `json:"criteria_rate"`
	OverallRate   int             `json:"rate"`
}

type IdeaList struct {
	Total         int             `json:"total"`
	LastSubmitted time.Time       `json:"last_submitted"`
	Result        []IdeaLightData `json:"result"`
}

type OverAllRate struct {
	Rate int `json:"overall_rate"`
}

type IdeaCondition string

const (
	RatedIdea    IdeaCondition = "RATED"
	NotRatedIdea IdeaCondition = "NOT RATED"
)

//CRITERIA

type CriteriaSpecData struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
