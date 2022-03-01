package models

import (
	"belli/onki-game-ideas-mongo-backend/responses"
	"time"
)

type IdeaCreate struct {
	Name        string
	WorkerID    string
	Date        time.Time
	Genre       string
	Mechanics   []string
	Description string
	Links       []responses.Link
	Files       []ParsedFile
	Paths       []string
}

type IdeaUpdate struct {
	ID          string
	Name        string
	WorkerID    string
	Date        time.Time
	Genre       string
	Mechanics   []string
	Description string
	Links       []responses.Link
}

type IdeaFilter struct {
	WorkerID  *string
	Name      *string
	Genre     *string
	Mechanics *[]string
	BeginDate *time.Time
	EndDate   *time.Time
	Condition *responses.IdeaCondition
	Limit     int
	Offset    int
}

type IdeaLightData struct {
	ID          string
	Worker      WorkerLightData
	Name        string
	Date        time.Time
	IsItNew     bool
	Description string
	FilePath    *string
	OverallRate int
}

type IdeaSpecData struct {
	ID            string
	Worker        WorkerLightData
	Name          string
	Date          time.Time
	Description   string
	Genre         string
	Mechanics     []string
	Links         []responses.Link
	FilePaths     []responses.Sketch
	CriteriaRates []responses.CriteriaRate
	OverallRate   int
}

type IdeaList struct {
	Total         int
	LastSubmitted time.Time
	Result        []IdeaLightData
}

type RateIdeaCritera struct {
	IdeaID     string
	CriteriaID string
	Rate       int
}

//CRITERIA

type CriteriaSpecData struct {
	ID   string
	Name string
}

type CriteriaUpdate struct {
	ID   string
	Name string
}

// GENRE

type GenreUpdate struct {
	OldGenre string
	NewGenre string
}

// MECHANIC

type MechanicUpdate struct {
	OldMech string
	NewMech string
}
