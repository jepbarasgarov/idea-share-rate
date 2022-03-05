package models

import (
	"belli/onki-game-ideas-mongo-backend/responses"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IdeaCreate struct {
	Name        string
	Worker      WorkerBsonModelInIdea
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
	IdeaID primitive.ObjectID
	Rating RatingStructInIdea
}

//CRITERIA

type CriteriaSpecData struct {
	ID   string
	Name string
}

type CriteriaUpdate struct {
	ID   primitive.ObjectID
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

////////////////////////////////////////////////////////////////////////BSON///////////////////////////////////////////////////////////////

type RatingStructInIdea struct {
	CriteriaID    primitive.ObjectID `bson:"criteria_id"`
	CrieteriaName string             `bson:"criteria_name"`
	UserID        primitive.ObjectID `bson:"user_id"`
	Rate          int                `bson:"rate"`
}

type ArrayOfRatesIdea struct {
	Rates []RatingStructInIdea `bson:"rates"`
}
