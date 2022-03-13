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
	AllFiles    []SketchSturctInIdea
}

type IdeaUpdate struct {
	ID          primitive.ObjectID
	Name        string
	Worker      WorkerBsonModelInIdea
	Date        time.Time
	Genre       string
	Mechanics   []string
	Description string
	Links       []responses.Link
}

type IdeaFilter struct {
	UserID    primitive.ObjectID
	WorkerID  *primitive.ObjectID
	Name      *string
	Genre     *string
	Mechanics *[]string
	BeginDate *primitive.DateTime
	EndDate   *primitive.DateTime
	Condition *responses.IdeaCondition
	Limit     int
	Offset    int
}

type IdeaLightData struct {
	ID          primitive.ObjectID    `bson:"_id" json:"id"`
	Worker      WorkerBsonModelInIdea `bson:"worker" json:"worker"`
	Name        string                `bson:"name" json:"game_name"`
	Date        time.Time             `bson:"date" json:"date"`
	IsItNew     bool                  `bson:"is_it_new" json:"is_it_new"`
	Description string                `bson:"description" json:"description"`
	FilePath    string                `bson:"path" json:"file_path"`
	OverallRate int                   `bson:"rate" json:"rate"`
	AvgRate     *float64              `bson:"avg" json:"avg,omitempty"`
	CreatesTS   time.Time             `bson:"create_ts"`
}

type IdeaSpecData struct {
	ID            primitive.ObjectID       `bson:"_id" json:"id"`
	Worker        WorkerBsonModelInIdea    `bson:"worker" json:"worker"`
	Name          string                   `bson:"name" json:"name"`
	Date          time.Time                `bson:"date" json:"date"`
	Description   string                   `bson:"description" json:"description"`
	Genre         string                   `bson:"genre" json:"genre"`
	Mechanics     []string                 `bson:"mechanics" json:"mechanics"`
	Links         []responses.Link         `bson:"links" json:"links"`
	FilePaths     []responses.Sketch       `bson:"files" json:"file_paths"`
	CriteriaRates []responses.CriteriaRate `bson:"rates" json:"criteria_rate"`
	OverallRate   int                      `bson:"rate" json:"rate"`
}

type IdeaList struct {
	Total         int             `json:"total"`
	LastSubmitted time.Time       `json:"last_submitted"`
	Result        []IdeaLightData `json:"result"`
}

type RateIdeaCritera struct {
	IdeaID primitive.ObjectID
	Rating RatingStructInIdea
}

//CRITERIA

type CriteriaSpecData struct {
	ID   primitive.ObjectID `bson:"_id"`
	Name string             `bson:"name"`
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

type RatingStructInIdea struct {
	CriteriaID    primitive.ObjectID `bson:"criteria_id,omitempty"`
	CrieteriaName string             `bson:"criteria_name,omitempty"`
	UserID        primitive.ObjectID `bson:"user_id,omitempty"`
	Rate          int                `bson:"rate,omitempty"`
}

type ArrayOfRatesIdea struct {
	Rates []RatingStructInIdea `bson:"rates"`
}

type SketchSturctInIdea struct {
	SketchID primitive.ObjectID `bson:"sketch_id"`
	FileName string             `bson:"name"`
	Path     string             `bson:"file_path"`
}
