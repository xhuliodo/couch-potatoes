package domain

type Movie struct {
	Id             string
	Title          string
	ReleaseYear    int
	Poster         string
	MoreInfoLink   string
	Genres         []Genre
	PeopleInvolved []Cast
}

type Genre struct {
	Id   string
	Name string
}

type Cast struct {
	Id   string
	Name string
	Role string
}
