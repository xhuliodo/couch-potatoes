// the query at this state and with my dataset all it does, it finds
// similair movies to one provided movie. This should be modified to take into
// account all movies rated with 1. In the end rated movies and movies in 
// the watchlist should be filtered out.

match (u:User{userId:"google-oauth2|104772264931362464545"})-[r:RATED{rating:1}]->(m:Movie) 
with u, m
// important to limit the movies to the last 10 rated movies,
// bcs taste changes over time and also it's too taxing on db,
// (without it, the more ratings a user has given the slower the response)
order by r desc
limit 10

match (m)-[:IN_GENRE|:ACTED_IN|:DIRECTED]-(t)<-[:IN_GENRE|:ACTED_IN|:DIRECTED]-(other:Movie)
where not exists( (u)-[:RATED]->(other) ) and not exists ( (u)-[:WATCH_LATER]->(other) )
with m, other, count(t) as intersection

// collect all ref points for liked movies
match (m)-[:IN_GENRE|:ACTED_IN|:DIRECTED]-(mt)
with m, other, intersection, collect(mt.name) as s1

// collect all ref points for suggestions
match (other)-[:IN_GENRE|:ACTED_IN|:DIRECTED]-(ot)
with m, other, intersection, s1, collect(ot.name) as s2

with m, other, intersection,s1,s2

with m, other, intersection, s1+[x in s2 where not x in s1] as union, s1, s2

with other, s1,s2,((1.0*intersection)/size(union)) as jaccard 

return other
order by jaccard desc 
limit 5