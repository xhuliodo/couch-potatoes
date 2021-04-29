package infrastructure

import (
	"github.com/google/uuid"
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

	query := "match (n:Genre) return n.genreId as Id, n.name as Name"
	parameters := map[string]interface{}{}
	res, err := session.Run(query, parameters)
	if err != nil {
		return nil, err
	}

	genres := []domain.Genre{}

	for res.Next() {
		genre := res.Record()
		genreId, _ := genre.Get("Id")
		genreUuid,_:=uuid.Parse(genreId.(string))
		genreName, _ := genre.Get("Name")
		g := domain.Genre{Id: genreUuid, Name: genreName.(string)}
		genres = append(genres, g)
	}
	return genres, nil
}
