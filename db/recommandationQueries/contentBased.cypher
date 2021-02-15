// the query at this state and with my dataset all it does, it finds
// similair movies to one provided movie. This should be modified to take into
// account all movies rated with 1. In the end rated movies and movies in 
// the watchlist should be filtered out.

MATCH (m:Movie {title: "Inception"})-[:IN_GENRE|:ACTED_IN|:DIRECTED]-(t)<-[:IN_GENRE|:ACTED_IN|:DIRECTED]-(other:Movie)
WITH m, other, COUNT(t) AS intersection, COLLECT(t.name) AS i
MATCH (m)-[:IN_GENRE|:ACTED_IN|:DIRECTED]-(mt)
WITH m,other, intersection,i, COLLECT(mt.name) AS s1
MATCH (other)-[:IN_GENRE|:ACTED_IN|:DIRECTED]-(ot)
WITH m,other,intersection,i, s1, COLLECT(ot.name) AS s2

WITH m,other,intersection,s1,s2

WITH m,other,intersection,s1+[x IN s2 WHERE NOT x IN s1] AS union, s1, s2

RETURN m.title, other.title, s1,s2,((1.0*intersection)/SIZE(union)) AS jaccard ORDER BY jaccard DESC LIMIT 100