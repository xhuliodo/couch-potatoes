package domain

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
	Id   string
	Name string
}

type Cast struct {
	Id   string
	Name string
	Role string
}

type PopulatiryScoredMovie struct {
	Movie
	AvgRating    float32
	CountRatings uint
}

type AggregateMovieRatings struct {
	Movie      Movie
	AllRatings []float32
}
