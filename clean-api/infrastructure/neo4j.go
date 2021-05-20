package infrastructure

import (
	"errors"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/xhuliodo/couch-potatoes/clean-api/domain"
)

type Neo4jRepository struct {
	Driver neo4j.Driver
}

func NewNeo4jRepository(Driver neo4j.Driver) *Neo4jRepository {
	return &Neo4jRepository{Driver}
}

func (nr *Neo4jRepository) GetAllGenres() ([]domain.Genre, error) {
	session := nr.Driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := "match (g:Genre) return g.genreId as Id, g.name as Name"
	parameters := map[string]interface{}{}

	res, err := session.Run(query, parameters)
	if err != nil {
		return nil, err
	}

	genres := []domain.Genre{}

	for res.Next() {
		genre := res.Record()
		genreId, _ := genre.Get("Id")
		genreName, _ := genre.Get("Name")
		g := domain.Genre{Id: genreId.(string), Name: genreName.(string)}
		genres = append(genres, g)
	}
	return genres, nil
}

func (nr *Neo4jRepository) GetUserById(userId string) (domain.User, error) {
	session := nr.Driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `match (u:User) 
			where u.userId=$userId 
			return u.userId as userId`
	parameters := map[string]interface{}{"userId": userId}

	res, err := session.Run(query, parameters)
	if err != nil {
		return domain.User{}, err
	}

	record, err := res.Single()
	if err != nil {
		return domain.User{}, errors.New("user does not exist")
	}

	existingUserId, bool := record.Get("userId")
	if !bool {
		return domain.User{}, errors.New("user does not exist")
	}
	existingUser := domain.User{Id: existingUserId.(string)}

	return existingUser, nil
}

func (nr *Neo4jRepository) GetMovieById(movieId string) (domain.Movie, error) {
	session := nr.Driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `match (m:Movie) 
			where m.movieId=$movieId 
			return m.movieId as movieId`
	parameters := map[string]interface{}{"movieId": movieId}

	res, err := session.Run(query, parameters)
	if err != nil {
		return domain.Movie{}, err
	}

	record, _ := res.Single()
	existingMovieId, bool := record.Get("movieId")
	if !bool {
		return domain.Movie{}, errors.New("movie does not exist")
	}
	existingMovie := domain.Movie{Id: existingMovieId.(string)}

	return existingMovie, nil
}

func (nr *Neo4jRepository) SaveGenrePreferences(userId string, genres []domain.Genre) error {
	session := nr.Driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `match (u:User{userId:$userId}), (g:Genre)
			  where g.genreId in $genres
			  merge (u)-[:FAVORITE]->(g)
			  return distinct u.userId as userId`
	genresIdInterface := make([]interface{}, len(genres))
	for i, g := range genres {
		genresIdInterface[i] = g.Id
	}
	parameters := map[string]interface{}{"userId": userId, "genres": genresIdInterface}

	res, err := session.Run(query, parameters)
	if err != nil {
		return err
	}

	record, _ := res.Single()
	_, bool := record.Get("userId")
	if !bool {
		return errors.New("genre preferences did not get saved")
	}

	return nil
}

func (nr *Neo4jRepository) RateMovie(userId, movieId string, rating int) error {
	session := nr.Driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `match (u:User{userId:$userId}), (m:Movie{movieId:$movieId})
			merge (u)-[r:RATED]->(m) 
			on create set r.rating=toInteger($rating)
			on match set r.rating=toInteger($rating)
			return m.movieId as movieId`
	parameters := map[string]interface{}{"userId": userId, "movieId": movieId, "rating": rating}

	res, err := session.Run(query, parameters)
	if err != nil {
		return err
	}

	record, _ := res.Single()
	_, bool := record.Get("movieId")
	if !bool {
		return errors.New("rating was not successful")
	}

	return nil
}

func (nr *Neo4jRepository) RegisterNewUser(user domain.User) error {
	session := nr.Driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `merge (u:User{userId:$userId})
			on create set u.userId=$userId
			return u.userId as userId`
	parameters := map[string]interface{}{"userId": user.Id}

	res, err := session.Run(query, parameters)
	if err != nil {
		return err
	}

	record, _ := res.Single()
	_, bool := record.Get("userId")
	if !bool {
		return errors.New("user was not registered successfully")
	}

	return nil
}

func (nr *Neo4jRepository) GetGenrePreferences(userId string) ([]domain.Genre, error) {
	session := nr.Driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
	match (u:User{userId:$userId})-[:FAVORITE]->(g:Genre)
	return g.genreId as Id
	`
	parameters := map[string]interface{}{"userId": userId}

	res, err := session.Run(query, parameters)
	if err != nil {
		return nil, err
	}

	genres := []domain.Genre{}

	// if _, err := res.Single(); err == nil {
	// 	return genres, errors.New("tuser has to give their genere preferences firs")
	// }
	for res.Next() {
		genre := res.Record()
		genreId, _ := genre.Get("Id")
		g := domain.Genre{Id: genreId.(string)}
		genres = append(genres, g)
	}

	if len(genres) == 0 {
		return genres, errors.New("no genre preference have been given")
	}

	return genres, nil
}

func (nr *Neo4jRepository) GetAllRatingsForMoviesInGenre(userId string, genres []domain.Genre) ([]domain.AggregateMovieRatings, error) {
	session := nr.Driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
	match (u:User{userId:$userId})
	with u 
	match (:User)-[r:RATED]->(m:Movie)-[:IN_GENRE]->(g)
	where not exists( (u)-[:RATED]->(m) ) and g.genreId in $genres
	return m.movieId as Id, 
		m.title as Title, 
		m.releaseYear as ReleaseYear, 
		m.imdbLink as MoreInfoLink, 
		count(distinct(g)) as GenreMultiplier, 
		collect(r.rating) as Ratings
	`
	genresIdInterface := make([]interface{}, len(genres))
	for i, g := range genres {
		genresIdInterface[i] = g.Id
	}
	parameters := map[string]interface{}{"userId": userId, "genres": genresIdInterface}

	res, err := session.Run(query, parameters)
	if err != nil {
		return nil, err
	}

	moviesAggregate := []domain.AggregateMovieRatings{}

	for res.Next() {
		movie := res.Record()
		movieId, _ := movie.Get("Id")
		title, _ := movie.Get("Title")
		releaseYearInterface, _ := movie.Get("ReleaseYear")
		releaseYearInt64, _ := releaseYearInterface.(int64)
		moreInfoLink, _ := movie.Get("MoreInfoLink")
		genreMultiplierInterface, _ := movie.Get("GenreMultiplier")
		genreMultiplierInt64, _ := genreMultiplierInterface.(int64)
		ratingsInterface, _ := movie.Get("Ratings")
		ratingsInterfaceSlice := ratingsInterface.([]interface{})

		m := domain.AggregateMovieRatings{Movie: domain.Movie{
			Id:           movieId.(string),
			Title:        title.(string),
			ReleaseYear:  int(releaseYearInt64),
			MoreInfoLink: moreInfoLink.(string),
		},
			GenreMatched: uint(genreMultiplierInt64),
			AllRatings:   convertRatingsInterfaceToFloatSlice(ratingsInterfaceSlice),
		}
		moviesAggregate = append(moviesAggregate, m)
	}

	if len(moviesAggregate) == 0 {
		return moviesAggregate, errors.New("there are no movies with ratings in the prefered genres")
	}
	return moviesAggregate, nil
}

func convertRatingsInterfaceToFloatSlice(ratingsInterfaceSlice []interface{}) []float32 {
	ratings := []float32{}
	for _, rating := range ratingsInterfaceSlice {
		r64 := rating.(float64)
		r := float32(r64)
		ratings = append(ratings, r)
	}
	return ratings
}
