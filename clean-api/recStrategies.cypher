// genre based rec
match (u:User{userId:$cypherParams.userId})-[:FAVORITE]->(g:Genre)
with g, u match (m:Movie)-[:IN_GENRE]->(g)
where not exists( (u)-[:RATED]->(m) ) and not exists ( (u)-[:WATCH_LATER]->(m) )
with  distinct(m) as movie
match (:User)-[r:RATED]->(movie)
with movie, count(r.rating) as howMany, avg(r.rating) as reviews
return movie
order by howMany desc, reviews desc
skip toInteger($skip)
limit toInteger($limit)

// user based rec
match (u1:User {userId:$cypherParams.userId})-[r:RATED]->(m:Movie)
with u1, avg(r.rating) AS u1_mean
match (u1)-[r1:RATED]->(m:Movie)<-[r2:RATED]-(u2)
with u1, u1_mean, u2, collect({r1: r1, r2: r2}) as ratings
where size(ratings) > toInteger($minimumRatings)
match (u2)-[r:RATED]->(m:Movie)
with u1, u1_mean, u2, avg(r.rating) as u2_mean, ratings unwind ratings as r
with sum( (r.r1.rating-u1_mean) * (r.r2.rating-u2_mean) ) as nom, sqrt( sum( (r.r1.rating - u1_mean)^2) * sum( (r.r2.rating - u2_mean) ^2)) as denom, u1, u2
where denom <> 0 with u1, u2, nom/denom as pearson
order by pearson desc
limit toInteger($peopleToCompare)
match (u2)-[r:RATED]->(m:Movie)
where not exists( (u1)-[:RATED]->(m) ) and not exists ( (u1)-[:WATCH_LATER]->(m) )
return m, sum( pearson * r.rating) as score
order by score desc
limit toInteger($moviesToRecommend)

// content based rec
match (u:User {userId:$cypherParams.userId})-[r:RATED{rating:1}]->(m:Movie)
with u, m
match (m)<-[:ACTED_IN|:DIRECTED|:WROTE]-(t)-[:ACTED_IN|:DIRECTED|:WROTE]->(other:Movie)
where not exists( (u)-[:RATED]->(other) ) and not exists ( (u)-[:WATCH_LATER]->(other) )
with m, other, count(t) as intersection
match (m)<-[:ACTED_IN|:DIRECTED|:WROTE]-(mt)
with m, other, intersection, collect(mt.name) as s1
match (other)<-[:ACTED_IN|:DIRECTED|:WROTE]-(ot)
with other, intersection, s1, collect(ot.name) as s2
with other, intersection,s1,s2
with other, intersection, s1+[x in s2 where not x in s1] as union, s1, s2
with other, s1,s2,((1.0*intersection)/size(union)) as jaccard
order by jaccard desc
return distinct other
limit toInteger($moviesToRecommend)