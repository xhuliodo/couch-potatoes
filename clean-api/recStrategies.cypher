// DONE
// IMPORTANT: content based rec

// step: get users all liked movies
match (u:User {userId:$cypherParams.userId})-[r:RATED{rating:1}]->(m:Movie)
with u, m
// step: compare these movies with all other movies that the user has yet
//       to consider, and count how many things in common they share
//       (actors, directors, writers)
match (m)<-[:ACTED_IN|:DIRECTED|:WROTE]-(t)-[:ACTED_IN|:DIRECTED|:WROTE]->(other:Movie)
where not exists( (u)-[:RATED]->(other) ) and not exists ( (u)-[:WATCH_LATER]->(other) )
with m, other, count(t) as intersection
// step: get all details for liked movies
match (m)<-[:ACTED_IN|:DIRECTED|:WROTE]-(mt)
with m, other, intersection, collect(mt.name) as s1
// step: get all details for rec movies
match (other)<-[:ACTED_IN|:DIRECTED|:WROTE]-(ot)
with other, intersection, s1, collect(ot.name) as s2
with other, intersection,s1,s2
// step: calculate the things in common
with other, intersection, s1+[x in s2 where not x in s1] as union, s1, s2
// step: calculate the jaccard and order by and limit 
with other, s1,s2,((1.0*intersection)/size(union)) as jaccard
order by jaccard desc
return distinct other
limit toInteger($moviesToRecommend)

// DONE
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

// transformed query

match (u:User{userId:"google-oauth2|104772264931362464545"})
with u 
match (:User)-[r:RATED]->(m:Movie)-[:IN_GENRE]->(g)
where not exists( (u)-[:RATED]->(m) ) and g.genreId in ["c4f88090-9166-4ebf-920b-ff9a34872b84", "acffe5b6-d327-43f6-b5ca-0a86f6780629"]
return m as movie, count(distinct(g)), collect(r.rating)

// DONE
// IMPORTANT: user based rec

//  step: get all user's rated movies, and returns
//        - avg of the user's all ratings
//        - all ratings in common with other users, filtering only the users that have 
//          over a certain number
match (u1:User {userId:$cypherParams.userId})-[r1:RATED]->(m:Movie)<-[r2:RATED]-(u2)
with u1, avg(r1.rating) AS u1_mean, u2, collect({r1: r1, r2: r2}) as ratings
where size(ratings) > toInteger($minimumRatings)
// step: get user's that have passed the first step ratings
//      - calculate their avg ratings across all movies, not just the ones in common
match (u2)-[r:RATED]->(m:Movie)
with u1, u1_mean, u2, avg(r.rating) as u2_mean, ratings unwind ratings as r


//      - calculate nom and demon of Pearson correlation
with sum( (r.r1.rating-u1_mean) * (r.r2.rating-u2_mean) ) as nom, sqrt( sum( (r.r1.rating - u1_mean)^2) * sum( (r.r2.rating - u2_mean) ^2)) as denom, u1, u2
//      - filter out where Pearson correlation is 0
where denom <> 0 
//      - calculate Pearson correlation
with u1, u2, nom/denom as pearson
//      - order by person and limit to the top $peopleToCompare

order by pearson desc
limit toInteger($peopleToCompare)
//  step: now check all ratings of the people that passed
match (u2)-[r:RATED]->(m:Movie)
where not exists( (u1)-[:RATED]->(m) ) and not exists ( (u1)-[:WATCH_LATER]->(m) )
//        calculate ratings with pearson as a multiplier and return the highest score
return m, sum( pearson * r.rating) as score
order by score desc
limit toInteger($moviesToRecommend)