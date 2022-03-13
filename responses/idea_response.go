package responses

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Link struct {
	Label string `bson:"label" json:"label"`
	URL   string `bson:"url" json:"url"`
}

type Sketch struct {
	ID       string `bson:"sketch_id" json:"id"`
	Name     string `bson:"name" json:"name"`
	FilePath string `bson:"file_path" json:"file_path"`
}

type CriteriaRate struct {
	ID     primitive.ObjectID  `bson:"criteria_id" json:"id"`
	Name   string              `bson:"criteria_name" json:"name"`
	Rate   int                 `bson:"rate" json:"rate"`
	UserID *primitive.ObjectID `bson:"user_id,omitempty" json:"user_id,omitempty"`
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
	ID   primitive.ObjectID `json:"id"`
	Name string             `json:"name"`
}
