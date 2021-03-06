package db

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/pkg/errors"
	"github.com/xhuliodo/couch-potatoes/clean-api/domain"
)

type Neo4jRepository struct {
	Driver neo4j.Driver
}

func NewNeo4jRepository(Driver neo4j.Driver) *Neo4jRepository {
	return &Neo4jRepository{Driver}
}

func (nr *Neo4jRepository) GetAllGenres() (
	genres []domain.Genre, err error,
) {
	session := nr.Driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := "match (g:Genre) return g.genreId as Id, g.name as Name"
	parameters := map[string]interface{}{}

	res, err := session.Run(query, parameters)
	if err != nil {
		cause := errors.New("db_connection")
		return genres, errors.Wrap(cause, err.Error())
	}

	for res.Next() {
		genre := res.Record()
		genreId, _ := genre.Get("Id")
		genreName, _ := genre.Get("Name")
		g := domain.Genre{Id: genreId.(string), Name: genreName.(string)}
		genres = append(genres, g)
	}

	return genres, nil
}

func (nr *Neo4jRepository) GetUserById(userId string) (
	existingUser domain.User, err error,
) {
	session := nr.Driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `match (u:User) 
			where u.userId=$userId 
			return u.userId as userId`
	parameters := map[string]interface{}{"userId": userId}

	res, err := session.Run(query, parameters)
	if err != nil {
		cause := errors.New("db_connection")
		return existingUser, errors.Wrap(cause, err.Error())
	}

	record, err := res.Single()
	if err != nil {
		cause := errors.New("not_found")
		return existingUser, errors.Wrap(cause, err.Error())
	}

	existingUserId, _ := record.Get("userId")
	existingUser = domain.User{Id: existingUserId.(string)}

	return existingUser, nil
}

func (nr *Neo4jRepository) GetMovieById(movieId string) (
	existingMovie domain.Movie, err error,
) {
	session := nr.Driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `match (m:Movie) 
			where m.movieId=$movieId 
			return m.movieId as movieId`
	parameters := map[string]interface{}{"movieId": movieId}

	res, err := session.Run(query, parameters)
	if err != nil {
		cause := errors.New("db_connection")
		return existingMovie, errors.Wrap(cause, err.Error())
	}

	record, err := res.Single()
	if err != nil {
		cause := errors.New("not_found")
		return existingMovie, errors.Wrap(cause, err.Error())
	}

	existingMovieId, _ := record.Get("movieId")
	existingMovie = domain.Movie{Id: existingMovieId.(string)}

	return existingMovie, nil
}

func (nr *Neo4jRepository) SaveGenrePreferences(userId string, genres []domain.Genre) error {
	session := nr.Driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `match (u:User{userId:$userId}), (g:Genre)
			  where g.genreId in $genres
			  merge (u)-[:FAVORITE]->(g)
			  return distinct u.userId as userId`
	genresIdInterface := make([]interface{}, len(genres))
	for i, g := range genres {
		genresIdInterface[i] = g.Id
	}
	parameters := map[string]interface{}{"userId": userId, "genres": genresIdInterface}

	if _, err := session.Run(query, parameters); err != nil {
		cause := errors.New("db_connection")
		return errors.Wrap(cause, err.Error())
	}

	return nil
}

func (nr *Neo4jRepository) RateMovie(
	userId, movieId string, rating int,
) error {
	session := nr.Driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `match (u:User{userId:$userId}), (m:Movie{movieId:$movieId})
			merge (u)-[r:RATED]->(m) 
			on create set r.rating=toInteger($rating)
			on match set r.rating=toInteger($rating)
			return m.movieId as movieId`
	parameters := map[string]interface{}{"userId": userId, "movieId": movieId, "rating": rating}

	if _, err := session.Run(query, parameters); err != nil {
		cause := errors.New("db_connection")
		return errors.Wrap(cause, err.Error())
	}

	return nil
}

func (nr *Neo4jRepository) RegisterNewUser(
	userId string,
) error {
	session := nr.Driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `merge (u:User{userId:$userId})
			on create set u.userId=$userId
			return u.userId as userId`
	parameters := map[string]interface{}{"userId": userId}

	if _, err := session.Run(query, parameters); err != nil {
		cause := errors.New("db_connection")
		return errors.Wrap(cause, err.Error())
	}

	return nil
}

func (nr *Neo4jRepository) GetGenrePreferencesCount(userId string) (
	genreCount uint, err error,
) {
	session := nr.Driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
	match (u:User{userId:$userId})-[:FAVORITE]->(g:Genre)
	return count(g) as genreCount
	`
	parameters := map[string]interface{}{"userId": userId}

	res, err := session.Run(query, parameters)
	if err != nil {
		cause := errors.New("db_connection")
		return genreCount, errors.Wrap(cause, err.Error())
	}

	rec, _ := res.Single()
	countInterface, _ := rec.Get("genreCount")
	count := countInterface.(int64)
	genreCount = uint(count)

	return genreCount, nil
}

func (nr *Neo4jRepository) GetAllRatingsForMoviesInGenre(userId string) (
	popularMovies domain.PopularMovies, err error,
) {
	session := nr.Driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
	match (u:User{userId:$userId})
	with u 
	match (:User)-[r:RATED]->(m:Movie)-[:IN_GENRE]->(g)<-[:FAVORITE]-(u)
	where not exists( (u)-[:RATED]->(m) )
	return m.movieId as Id, 
		count(distinct(g)) as GenreMultiplier, 
		count(r.rating) as RatingsCount,
        avg(r.rating) as AvgRating
	`
	parameters := map[string]interface{}{"userId": userId}

	res, err := session.Run(query, parameters)
	if err != nil {
		cause := errors.New("db_connection")
		return nil, errors.Wrap(cause, err.Error())
	}

	for res.Next() {
		movie := res.Record()
		movieId, _ := movie.Get("Id")
		genreMatchedInterface, _ := movie.Get("GenreMultiplier")
		genreMatchedInt64, _ := genreMatchedInterface.(int64)
		ratingsCountInterface, _ := movie.Get("RatingsCount")
		ratingCount := ratingsCountInterface.(int64)
		avgRatingInterface, _ := movie.Get("AvgRating")
		avgRating := avgRatingInterface.(float64)

		m := domain.PopularMovie{Movie: domain.Movie{
			Id: movieId.(string),
		},
			AvgRating:    avgRating,
			GenreMatched: uint(genreMatchedInt64),
			RatingsCount: uint(ratingCount),
		}

		popularMovies = append(popularMovies, m)
	}

	return popularMovies, nil
}

func (nr *Neo4jRepository) GetMoviesDetails(userIds []string) (
	domain.MoviesDetails, error,
) {
	moviesDetails := make(domain.MoviesDetails)

	session := nr.Driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
	match (m:Movie)
	where m.movieId in $movieIds
	return m.movieId as Id,
		m.title as Title, 
		m.releaseYear as ReleaseYear, 
		m.imdbLink as MoreInfoLink
	`
	movieIdsInterface := make([]interface{}, len(userIds))
	for i, u := range userIds {
		movieIdsInterface[i] = u
	}
	parameters := map[string]interface{}{"movieIds": movieIdsInterface}

	res, err := session.Run(query, parameters)
	if err != nil {
		cause := errors.New("db_connection")
		return moviesDetails, errors.Wrap(cause, err.Error())
	}

	for res.Next() {
		rec := res.Record()
		movieIdInterface, _ := rec.Get("Id")
		movieId, _ := movieIdInterface.(string)
		title, _ := rec.Get("Title")
		releaseYearInterface, _ := rec.Get("ReleaseYear")
		releaseYearInt64, _ := releaseYearInterface.(int64)
		moreInfoLink, _ := rec.Get("MoreInfoLink")
		moviesDetails[movieId] = domain.MovieDetails{
			Title:        title.(string),
			ReleaseYear:  int(releaseYearInt64),
			MoreInfoLink: moreInfoLink.(string),
		}
	}

	return moviesDetails, nil
}

func (nr *Neo4jRepository) GetUserRatingsCount(userId string) (
	userRatingsCount uint, err error,
) {
	session := nr.Driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
	match (u:User{userId: $userId})-[:RATED]->(m:Movie)
    return count(m) as RatedMoviesCount
	`
	parameters := map[string]interface{}{"userId": userId}

	res, err := session.Run(query, parameters)
	if err != nil {
		cause := errors.New("db_connection")
		return userRatingsCount, errors.Wrap(cause, err.Error())
	}

	rec, _ := res.Single()
	countInterface, _ := rec.Get("RatedMoviesCount")
	count := countInterface.(int64)
	userRatingsCount = uint(count)

	return userRatingsCount, nil
}

func (nr *Neo4jRepository) GetSimilairUsersAndTheirAvgRating(userId string) (
	domain.UsersToCompare, error,
) {
	usersToCompare := make(domain.UsersToCompare)

	session := nr.Driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
	match (u1:User {userId:$userId})-[r1:RATED]->(m:Movie)<-[r2:RATED]-(u2)
	return u1.userId as UserToRecId, 
		u2.userId as UserToCompareId, 
		avg(r1.rating) AS UserToRecAvgRating, 
		avg(r2.rating) AS UserToCompAvgRating, 
		collect([toFloat(r1.rating), toFloat(r2.rating)]) as RatingsInCommon
	`
	parameters := map[string]interface{}{"userId": userId}

	res, err := session.Run(query, parameters)
	if err != nil {
		cause := errors.New("db_connection")
		return usersToCompare, errors.Wrap(cause, err.Error())
	}

	for res.Next() {
		rec := res.Record()
		userToRecAvgRating, _ := rec.Get("UserToRecAvgRating")
		userToCompAvgRating, _ := rec.Get("UserToCompAvgRating")

		userToCompareId, _ := rec.Get("UserToCompareId")
		ratingsInCommonInterface, _ := rec.Get("RatingsInCommon")

		ratingsInCommonInterfaceSlice := ratingsInCommonInterface.([]interface{})

		usersToCompare[userToCompareId.(string)] = &domain.UserToCompare{
			RatingsInCommon:     convertRatingsInCommonInterfaceSlice(ratingsInCommonInterfaceSlice),
			UserToRecAvgRating:  userToRecAvgRating.(float64),
			UserToCompAvgRating: userToCompAvgRating.(float64),
		}
	}

	return usersToCompare, nil
}

func convertRatingsInCommonInterfaceSlice(ratingsInterfaceSlice []interface{}) []domain.RatingInCommon {
	ratingsInCommon := []domain.RatingInCommon{}
	for _, rating := range ratingsInterfaceSlice {
		r := rating.([]interface{})
		ratingInCommon := domain.RatingInCommon{
			UserToRecommendRating: r[0].(float64),
			UserToCompareRating:   r[1].(float64),
		}

		ratingsInCommon = append(ratingsInCommon, ratingInCommon)
	}
	return ratingsInCommon
}

func (nr *Neo4jRepository) GetRatedMoviesForUsersYetToBeConsidered(
	userId string,
	userIds []string,
) (domain.ScoringMovies, error) {
	scoringMovies := make(domain.ScoringMovies)

	session := nr.Driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
	match (userToRec:User{userId:$userToRecId})
    with userToRec
	match (u:User)-[r:RATED]->(m:Movie)
	where u.userId in $userIds
	and not exists ((userToRec)-[:RATED]->(m))
	and not exists ((userToRec)-[:WATCH_LATER]->(m))
	return m.movieId as MovieId,
	   m.title as MovieTitle,
	   m.releaseYear as ReleaseYear, 
	   m.imdbLink as MoreInfoLink,
	   collect([u.userId, r.rating]) as UserRatingCollection
	`
	userIdsInterface := make([]interface{}, len(userIds))
	for i, u := range userIds {
		userIdsInterface[i] = u
	}
	parameters := map[string]interface{}{"userToRecId": userId, "userIds": userIdsInterface}

	res, err := session.Run(query, parameters)
	if err != nil {
		cause := errors.New("db_connection")
		return scoringMovies, errors.Wrap(cause, err.Error())
	}

	for res.Next() {
		rec := res.Record()
		movieIdInterface, _ := rec.Get("MovieId")
		movieId := movieIdInterface.(string)
		movieTitleInterface, _ := rec.Get("MovieTitle")
		movieTitle := movieTitleInterface.(string)
		releaseYearInterface, _ := rec.Get("ReleaseYear")
		releaseYearInt64, _ := releaseYearInterface.(int64)
		moreInfoLink, _ := rec.Get("MoreInfoLink")
		userRatingCollectionInterface, _ := rec.Get("UserRatingCollection")
		userRatingCollectionInterfaceSlice := userRatingCollectionInterface.([]interface{})

		scoringMovies[movieId] = &domain.Details{
			Movie: domain.Movie{
				Title:        movieTitle,
				ReleaseYear:  int(releaseYearInt64),
				MoreInfoLink: moreInfoLink.(string),
			},
		}

		ratings := convertRatedMoviesInterfaceSlice(userRatingCollectionInterfaceSlice)

		scoringMovies[movieId].Ratings = ratings
	}

	return scoringMovies, nil
}

func convertRatedMoviesInterfaceSlice(urcis []interface{}) domain.Rating {
	ratings := domain.Rating{}
	for _, ur := range urcis {
		r := ur.([]interface{})
		userId := r[0].(string)
		rating := r[1].(float64)
		ratings[userId] = rating
	}
	return ratings
}

func (nr *Neo4jRepository) GetAllLikedMovies(userId string) (
	domain.UsersLikedMovies, error,
) {
	likedMovies := make(domain.UsersLikedMovies)

	session := nr.Driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
	match (u:User {userId:$userId})-[r:RATED{rating:1}]->(m:Movie)
	return m.movieId as MovieId
	`
	parameters := map[string]interface{}{"userId": userId}

	res, err := session.Run(query, parameters)
	if err != nil {
		cause := errors.New("db_connection")
		return likedMovies, errors.Wrap(cause, err.Error())
	}

	for res.Next() {
		rec := res.Record()
		movieIdInterface, _ := rec.Get("MovieId")
		movieId := movieIdInterface.(string)

		likedMovies[movieId] = domain.UsersLikedMovie{AllCast: map[string]bool{}}
	}

	return likedMovies, nil
}

func (nr *Neo4jRepository) GetMoviesCasts(movieIds []string, movies domain.MoviesWithoutCastDetails) error {
	session := nr.Driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
	match (m)<-[:ACTED_IN|:DIRECTED|:WROTE]-(cast)
	where m.movieId in $movieIds
	return m.movieId as MovieId, 
		collect(cast.castId) as CastDetails
	`
	movieIdsInterface := make([]interface{}, len(movieIds))
	for i, m := range movieIds {
		movieIdsInterface[i] = m
	}
	parameters := map[string]interface{}{"movieIds": movieIdsInterface}

	res, err := session.Run(query, parameters)
	if err != nil {
		cause := errors.New("db_connection")
		return errors.Wrap(cause, err.Error())
	}

	for res.Next() {
		rec := res.Record()
		movieIdInterface, _ := rec.Get("MovieId")
		movieId := movieIdInterface.(string)
		castDetailsInterface, _ := rec.Get("CastDetails")
		castDetailsInterfaceSlice := castDetailsInterface.([]interface{})
		castDetailsSlice := getCastDetailsSlice(castDetailsInterfaceSlice)

		movies.PopulateWithCast(movieId, castDetailsSlice)
	}

	return nil
}

func getCastDetailsSlice(castDetailsInterfaceSlice []interface{}) []string {
	castSlice := []string{}
	for _, castDetailsInterface := range castDetailsInterfaceSlice {
		castDetails := castDetailsInterface.(string)
		castSlice = append(castSlice, castDetails)
	}
	return castSlice
}

func (nr *Neo4jRepository) GetSimilarMoviesToAlreadyLikedOnes(userId string, movieIds []string) (domain.SimilarMoviesToLikedOnes, error) {
	emptySimilarMovies := domain.SimilarMoviesToLikedOnes{}

	session := nr.Driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
	match (u:User{userId:$userId})
	match (m:Movie)<-[:ACTED_IN|:DIRECTED|:WROTE]-(t)-[:ACTED_IN|:DIRECTED|:WROTE]->(other:Movie)
	where m.movieId in $movieIds
    and not exists( (u)-[:RATED]->(other) ) 
    and not exists ( (u)-[:WATCH_LATER]->(other) )
	return other.movieId as MovieId
	`
	movieIdsInterface := make([]interface{}, len(movieIds))
	for i, m := range movieIds {
		movieIdsInterface[i] = m
	}
	parameters := map[string]interface{}{"userId": userId, "movieIds": movieIdsInterface}

	res, err := session.Run(query, parameters)
	if err != nil {
		cause := errors.New("db_connection")
		return emptySimilarMovies, errors.Wrap(cause, err.Error())
	}

	for res.Next() {
		rec := res.Record()
		movieIdInterface, _ := rec.Get("MovieId")
		movieId := movieIdInterface.(string)

		emptySimilarMovies[movieId] = domain.SimilarMovieToLikedOnes{}
	}

	return emptySimilarMovies, nil
}

func (nr *Neo4jRepository) AddToWatchlist(
	userId, movieId string, timeOfAdding int64,
) error {
	session := nr.Driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
	match (u:User{userId:$userId}), (m:Movie{movieId:$movieId})
    merge (u)-[w:WATCH_LATER]->(m) on create set w.createdAt=$timeOfAdding
	return m.movieId as MovieId
	`
	parameters := map[string]interface{}{"userId": userId, "movieId": movieId, "timeOfAdding": timeOfAdding}
	if _, err := session.Run(query, parameters); err != nil {
		cause := errors.New("db_connection")
		return errors.Wrap(cause, err.Error())
	}

	return nil
}

func (nr *Neo4jRepository) RemoveFromWatchlist(userId, movieId string) error {
	session := nr.Driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
	match (u:User {userId:$userId})-[w:WATCH_LATER]->(m:Movie{movieId:$movieId})
    delete w
	return m.movieId as MovieId
	`
	parameters := map[string]interface{}{"userId": userId, "movieId": movieId}
	if _, err := session.Run(query, parameters); err != nil {
		cause := errors.New("db_connection")
		return errors.Wrap(cause, err.Error())
	}

	return nil
}

func (nr *Neo4jRepository) GetWatchlist(userId string) (
	watchlist domain.UserWatchlist, err error,
) {
	session := nr.Driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
	match (u:User{userId:$userId})-[w:WATCH_LATER]->(m:Movie)
	where  not exists((u)-[:RATED]->(m))
	return m.movieId as MovieId,
		w.createdAt as TimeAdded
	`
	parameters := map[string]interface{}{"userId": userId}

	res, err := session.Run(query, parameters)
	if err != nil {
		cause := errors.New("db_connection")
		return watchlist, errors.Wrap(cause, err.Error())
	}

	for res.Next() {
		rec := res.Record()
		movieIdInterface, _ := rec.Get("MovieId")
		movieId := movieIdInterface.(string)
		timeAddedInterface, _ := rec.Get("TimeAdded")
		timeAdded := timeAddedInterface.(int64)

		watchlistItem := domain.Watchlist{
			Movie: domain.Movie{
				Id: movieId,
			},
			TimeAdded: timeAdded,
		}

		watchlist = append(watchlist, watchlistItem)
	}

	return watchlist, nil
}

func (nr *Neo4jRepository) GetWatchlistHistory(userId string) (
	watchlistHistory domain.UserWatchlist, err error,
) {
	session := nr.Driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
	match (u:User {userId:$userId})-[w:WATCH_LATER]->(m:Movie)<-[r:RATED]-(u)
	return m.movieId as MovieId,
		w.createdAt as TimeAdded,
		r.rating as Rating
	`
	parameters := map[string]interface{}{"userId": userId}

	res, err := session.Run(query, parameters)
	if err != nil {
		cause := errors.New("db_connection")
		return watchlistHistory, errors.Wrap(cause, err.Error())
	}

	for res.Next() {
		rec := res.Record()
		movieIdInterface, _ := rec.Get("MovieId")
		movieId := movieIdInterface.(string)
		timeAddedInterface, _ := rec.Get("TimeAdded")
		timeAdded := timeAddedInterface.(int64)
		ratingInterface, _ := rec.Get("Rating")
		rating := ratingInterface.(int64)

		watchlistItem := domain.Watchlist{
			Movie: domain.Movie{
				Id: movieId,
			},
			TimeAdded: timeAdded,
			Rating:    float64(rating),
		}

		watchlistHistory = append(watchlistHistory, watchlistItem)
	}

	return watchlistHistory, nil
}

func (nr *Neo4jRepository) Healthcheck() error {
	err := nr.Driver.VerifyConnectivity()
	if err != nil {
		cause := errors.New("db_connection")
		return errors.Wrap(cause, err.Error())
	}
	return nil
}
