// CREATE CONSTRAINT ON (m:Movie) ASSERT m.movieId IS UNIQUE;
// CREATE CONSTRAINT ON (g:Genre) ASSERT g.name IS UNIQUE;
// CREATE CONSTRAINT ON (u:User) ASSERT u.userId IS UNIQUE;

// The moment you set a unique constraint on a label of a node
// you have already put a index on in, that's how it keeps track
// whether a value is unique or not
// CREATE INDEX ON :Movie(movieId);
// CREATE INDEX ON :User(userId);

// Load movies
LOAD CSV WITH HEADERS FROM "file:///movies.csv" AS row
MERGE (m:Movie {movieId: row.movieId})
ON CREATE SET m.title = row.title, m.releaseYear = toInteger(row.year), m.imdbLink=row.imdbUrl, m.imdbId = row.imdbId
WITH *
UNWIND split(row.genres, "|") AS genre
MERGE (g:Genre {name: genre})
ON CREATE SET g.genreId = apoc.create.uuid()
MERGE (m)-[:IN_GENRE]->(g); 

// Load users / ratings
LOAD CSV WITH HEADERS FROM "file:///ratings.csv" AS row
MERGE (u:User {userId: row.userId})
WITH *
MATCH (m:Movie {movieId: row.movieId})
CREATE (u)-[r:RATED]->(m)
SET r.rating = toFloat(row.rating);

// Load directors
LOAD CSV WITH HEADERS FROM "file:///directors.csv" AS row
MATCH (m:Movie {imdbId: row.imdbId})
WITH *
UNWIND split(row.directorName, "|") AS director
MERGE (c:Cast {name: director})
ON CREATE SET c.castId = apoc.create.uuid()
CREATE (c)-[:DIRECTED]->(m);

// Load writers
LOAD CSV WITH HEADERS FROM "file:///writers.csv" AS row
MATCH (m:Movie {imdbId: row.imdbId})
WITH *
UNWIND split(row.writerName, "|") AS writer
MERGE (c:Cast {name: writer})
ON CREATE SET c.castId = apoc.create.uuid()
CREATE (c)-[:WROTE]->(m);