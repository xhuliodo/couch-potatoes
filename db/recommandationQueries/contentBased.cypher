// the query at this state and with my dataset all it does, it finds
// similair movies to one provided movie. This should be modified to take into
// account all movies rated with 1. In the end rated movies and movies in 
// the watchlist should be filtered out.

match (u:User{userId:"google-oauth2|104772264931362464545"})-[r:RATED{rating:1}]->(m:Movie) 
with u, m
match (m)-[:IN_GENRE|:ACTED_IN|:DIRECTED]-(t)<-[:IN_GENRE|:ACTED_IN|:DIRECTED]-(other:Movie)
// it can be optimizes with index on actors id and collecting id instead of name
with u, m, other, count(t) as intersection, collect(t.name) as i
match (m)-[:IN_GENRE|:ACTED_IN|:DIRECTED]-(mt)
with u, m, other, intersection, i, collect(mt.name) as s1
match (other)-[:IN_GENRE|:ACTED_IN|:DIRECTED]-(ot)
with u, m,other,intersection,i, s1, collect(ot.name) as s2

with u, m,other,intersection,s1,s2

with u, m, other, intersection, s1+[x in s2 where not x in s1] as union, s1, s2

where not exists( (u)-[:RATED]->(m) ) and not exists ( (u)-[:WATCH_LATER]->(m) )

return m.title, other.title, s1,s2,((1.0*intersection)/size(union)) as jaccard 
order by jaccard desc 
limit 10