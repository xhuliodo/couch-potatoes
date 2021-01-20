// This query is done in the registration process, 
// when we know nothing about the user but the genres
// he likes.

match (m:Movie)-[:IN_GENRE]->(g:Genre) 
where g.name in ['${genre_1}','${genre_2}','${genre_3}', ...] 
with  distinct(m) as movie 
match (:User)-[r:RATED]->(movie) 
with movie, count(r.rating) as howMany,avg(r.rating) as reviews 
return movie 
order by howMany desc, reviews desc 
limit ${howeverManyMoviesYouWantToShow}