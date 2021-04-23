package infrastructure

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/xhuliodo/couch-potatoes/clean-api/domain"
)

type Neo4jRepository struct {
	driver neo4j.Driver
}

func NewNeo4jRepository(driver neo4j.Driver) *Neo4jRepository {
	return &Neo4jRepository{driver}
}

// func executeQuery(){

// }

func (nr *Neo4jRepository) GetAllGenres() ([]domain.Genre, error) {
	session := nr.driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := "match (n:Genre) return n.genreId as Id, n.name as Name"
	parameters := map[string]interface{}{}
	res, err := session.Run(query, parameters)
	if err != nil {
		return nil, err
	}

	genres := []domain.Genre{}
	
	for res.Next() {
		rec:=res.Record().Get("Id")
		g:=domain.Genre{Id: }
		genres=append(genres, )
	}
	record, err := res.Single()
	if err != nil {
		return nil, err
	}
	empty := []domain.Genre{}
	return empty, nil
}
