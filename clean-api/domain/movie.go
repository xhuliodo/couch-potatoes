package domain

type Movie struct {
	Id           int
	Title        string
	ReleaseYear  int
	Poster       string
	MoreInfoLink string
	Genres       []Genre
	Writer       []Writer
	Director     []Director
	Actor        []Actor
}

type Genre struct {
	Id   int
	Name string
}

type Actor struct {
	Id   int
	Name string
}

type Director struct {
	Id   int
	Name string
}

type Writer struct {
	Id   int
	Name string
}

type MovieWatcher struct {
	Id             int
	Name           string
	RatedMovies    []RatedMovie
	Watchlist      []Movie
	FavoriteGenres []Genre
}

type RatedMovie struct {
	Movie
	Rating float32
}

func (mw *MovieWatcher) AddFavoriteGenre(g Genre) {
	mw.FavoriteGenres = append(mw.FavoriteGenres, g)
}
