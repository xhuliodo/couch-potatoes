package infrastructure

import (
	"errors"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/xhuliodo/couch-potatoes/clean-api/application"
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

func (nr *Neo4jRepository) GetUserById(userId string) (application.User, error) {
	session := nr.Driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `match (u:User) 
			where u.userId=$userId 
			return u.userId as userId`
	parameters := map[string]interface{}{"userId": userId}

	res, err := session.Run(query, parameters)
	if err != nil {
		return application.User{}, err
	}

	record, _ := res.Single()
	existingUserId, bool := record.Get("userId")
	if !bool {
		return application.User{}, errors.New("user does not exist")
	}
	existingUser := application.User{Id: existingUserId.(string)}

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
	existingMovie := domain.Movie{Id: domain.MovieID(existingMovieId.(string))}

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
