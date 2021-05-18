// genre based rec

// step: get user favorite genres
match (u:User{userId:"google-oauth2|104772264931362464545"})-[:FAVORITE]->(g:Genre)
// step: get all distinct movies that belong in 
//       those genres that the user has yet to rate.
with g, u 
match (m:Movie)-[:IN_GENRE]->(g)
where not exists( (u)-[:RATED]->(m) ) 
with distinct(m) as movie
// step: get ratings of those movies and return 
//       movies, how many ratings it has gotten
//       and an average per movie.
//       order by howMany and reviews DESC while
//       skipping by x and limiting by x
match (:User)-[r:RATED]->(movie)
with movie, count(r.rating) as howMany, avg(r.rating) as reviews
return movie
order by howMany desc, reviews desc
skip toInteger($skip)
limit toInteger($limit)

// ! if for some reason you might want to remove the multiplier
// for genres you can replace the above query with this one
// should be even faster since you're replacing two queries
// with one
match (:User)-[r:RATED]->(m:Movie)-[:IN_GENRE]->(g)
where not exists( (u)-[:RATED]->(m) ) 
return distinct(m) as movie, collect(r.rating) as ratings


// IMPORTANT: user based rec

//  step: get all user's rated movies, return user and
//       average of all rated movies
match (u1:User {userId:$cypherParams.userId})-[r:RATED]->(m:Movie)
with u1, avg(r.rating) AS u1_mean
//  step: 
match (u1)-[r1:RATED]->(m:Movie)<-[r2:RATED]-(u2)
with u1, u1_mean, u2, collect({r1: r1, r2: r2}) as ratings
where size(ratings) > toInteger($minimumRatings)
// step: 
match (u2)-[r:RATED]->(m:Movie)
with u1, u1_mean, u2, avg(r.rating) as u2_mean, ratings unwind ratings as r
with sum( (r.r1.rating-u1_mean) * (r.r2.rating-u2_mean) ) as nom, sqrt( sum( (r.r1.rating - u1_mean)^2) * sum( (r.r2.rating - u2_mean) ^2)) as denom, u1, u2
where denom <> 0 with u1, u2, nom/denom as pearson
order by pearson desc
limit toInteger($peopleToCompare)
//  step: 
match (u2)-[r:RATED]->(m:Movie)
where not exists( (u1)-[:RATED]->(m) ) and not exists ( (u1)-[:WATCH_LATER]->(m) )
return m, sum( pearson * r.rating) as score
order by score desc
limit toInteger($moviesToRecommend)


// content based rec

// step:
match (u:User {userId:$cypherParams.userId})-[r:RATED{rating:1}]->(m:Movie)
with u, m
// step: 
match (m)<-[:ACTED_IN|:DIRECTED|:WROTE]-(t)-[:ACTED_IN|:DIRECTED|:WROTE]->(other:Movie)
where not exists( (u)-[:RATED]->(other) ) and not exists ( (u)-[:WATCH_LATER]->(other) )
with m, other, count(t) as intersection
// step:
match (m)<-[:ACTED_IN|:DIRECTED|:WROTE]-(mt)
with m, other, intersection, collect(mt.name) as s1
// step:
match (other)<-[:ACTED_IN|:DIRECTED|:WROTE]-(ot)
with other, intersection, s1, collect(ot.name) as s2
with other, intersection,s1,s2
with other, intersection, s1+[x in s2 where not x in s1] as union, s1, s2
with other, s1,s2,((1.0*intersection)/size(union)) as jaccard
order by jaccard desc
return distinct other
limit toInteger($moviesToRecommend)