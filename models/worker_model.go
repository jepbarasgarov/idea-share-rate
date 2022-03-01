package models

type WorkerCreate struct {
	Firstname string
	Lastname  string
	Position  string
}

type WorkerUpdate struct {
	ID        string
	Firstname string
	Lastname  string
	Position  string
}

type WorkerLightData struct {
	ID        string
	Firstname string
	Lastname  string
	Position  string
}
