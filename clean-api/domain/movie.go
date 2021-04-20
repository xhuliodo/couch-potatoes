package domain

import (
	"github.com/google/uuid"
)

type MovieID string

type Movie struct {
	Id             MovieID
	Title          string
	ReleaseYear    int
	Poster         string
	MoreInfoLink   string
	Genres         []Genre
	PeopleInvolved []Cast
}

type Genre struct {
	Id   uuid.UUID
	Name string
}

type Cast struct {
	Id   uuid.UUID
	Name string
	Role string
}
