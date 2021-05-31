package infrastructure

import (
	"errors"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/xhuliodo/couch-potatoes/clean-api/domain"
)

type Neo4jRepository struct {
	Driver neo4j.Driver
}

func NewNeo4jRepository(Driver neo4j.Driver) *Neo4jRepository {
	return &Neo4jRepository{Driver}
}

func (nr *Neo4jRepository) GetAllGenres() ([]domain.Genre, error) {
	session := nr.Driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := "match (g:Genre) return g.genreId as Id, g.name as Name"
	parameters := map[string]interface{}{}

	res, err := session.Run(query, parameters)
	if err != nil {
		return nil, err
	}

	genres := []domain.Genre{}

	for res.Next() {
		genre := res.Record()
		genreId, _ := genre.Get("Id")
		genreName, _ := genre.Get("Name")
		g := domain.Genre{Id: genreId.(string), Name: genreName.(string)}
		genres = append(genres, g)
	}
	return genres, nil
}

func (nr *Neo4jRepository) GetUserById(
	userId string,
) (domain.User, error) {
	session := nr.Driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `match (u:User) 
			where u.userId=$userId 
			return u.userId as userId`
	parameters := map[string]interface{}{"userId": userId}

	res, err := session.Run(query, parameters)
	if err != nil {
		return domain.User{}, err
	}

	record, err := res.Single()
	if err != nil {
		return domain.User{}, errors.New("user does not exist")
	}

	existingUserId, bool := record.Get("userId")
	if !bool {
		return domain.User{}, errors.New("user does not exist")
	}
	existingUser := domain.User{Id: existingUserId.(string)}

	return existingUser, nil
}

func (nr *Neo4jRepository) GetMovieById(
	movieId string,
) (domain.Movie, error) {
	session := nr.Driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `match (m:Movie) 
			where m.movieId=$movieId 
			return m.movieId as movieId`
	parameters := map[string]interface{}{"movieId": movieId}

	res, err := session.Run(query, parameters)
	if err != nil {
		return domain.Movie{}, err
	}

	record, err := res.Single()
	if err != nil {
		return domain.Movie{}, errors.New("movie does not exist")
	}

	existingMovieId, bool := record.Get("movieId")
	if !bool {
		return domain.Movie{}, errors.New("movie does not exist")
	}
	existingMovie := domain.Movie{Id: existingMovieId.(string)}

	return existingMovie, nil
}

func (nr *Neo4jRepository) SaveGenrePreferences(
	userId string, genres []domain.Genre,
) error {
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

	res, err := session.Run(query, parameters)
	if err != nil {
		return err
	}

	record, _ := res.Single()
	_, bool := record.Get("userId")
	if !bool {
		return errors.New("genre preferences did not get saved")
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

	res, err := session.Run(query, parameters)
	if err != nil {
		return err
	}

	record, _ := res.Single()
	_, bool := record.Get("movieId")
	if !bool {
		return errors.New("rating was not successful")
	}

	return nil
}

func (nr *Neo4jRepository) RegisterNewUser(
	user domain.User,
) error {
	session := nr.Driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `merge (u:User{userId:$userId})
			on create set u.userId=$userId
			return u.userId as userId`
	parameters := map[string]interface{}{"userId": user.Id}

	res, err := session.Run(query, parameters)
	if err != nil {
		return err
	}

	record, _ := res.Single()
	_, bool := record.Get("userId")
	if !bool {
		return errors.New("user was not registered successfully")
	}

	return nil
}

func (nr *Neo4jRepository) GetGenrePreferences(
	userId string,
) ([]domain.Genre, error) {
	session := nr.Driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
	match (u:User{userId:$userId})-[:FAVORITE]->(g:Genre)
	return g.genreId as Id
	`
	parameters := map[string]interface{}{"userId": userId}

	res, err := session.Run(query, parameters)
	if err != nil {
		return nil, err
	}

	genres := []domain.Genre{}

	// if _, err := res.Single(); err == nil {
	// 	return genres, errors.New("tuser has to give their genere preferences firs")
	// }
	for res.Next() {
		genre := res.Record()
		genreId, _ := genre.Get("Id")
		g := domain.Genre{Id: genreId.(string)}
		genres = append(genres, g)
	}

	if len(genres) == 0 {
		return genres, errors.New("no genre preference have been given")
	}

	return genres, nil
}

func (nr *Neo4jRepository) GetAllRatingsForMoviesInGenre(
	userId string,
) (domain.PopularMovies, error) {
	popularMovies := domain.PopularMovies{}
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
		return nil, err
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

	if len(popularMovies) == 0 {
		return popularMovies, errors.New("there are no movies with ratings in the prefered genres")
	}
	return popularMovies, nil
}

func (nr *Neo4jRepository) GetMoviesDetails(
	userIds []string,
) (domain.MoviesDetails, error) {
	emptyMovieDetails := domain.MoviesDetails{}

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
		return emptyMovieDetails, err
	}

	for res.Next() {
		rec := res.Record()
		movieIdInterface, _ := rec.Get("Id")
		movieId, _ := movieIdInterface.(string)
		title, _ := rec.Get("Title")
		releaseYearInterface, _ := rec.Get("ReleaseYear")
		releaseYearInt64, _ := releaseYearInterface.(int64)
		moreInfoLink, _ := rec.Get("MoreInfoLink")

		emptyMovieDetails[movieId] = domain.MovieDetails{
			Title:        title.(string),
			ReleaseYear:  int(releaseYearInt64),
			MoreInfoLink: moreInfoLink.(string),
		}
	}

	return emptyMovieDetails, nil
}

func (nr *Neo4jRepository) GetUserRatingsCount(
	userId string,
) (uint, error) {
	session := nr.Driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
	match (u:User{userId: $userId})-[:RATED]->(m:Movie)
    return count(m) as RatedMoviesCount
	`
	parameters := map[string]interface{}{"userId": userId}

	res, err := session.Run(query, parameters)
	if err != nil {
		return 0, err
	}

	record, _ := res.Single()
	ratedMoviesCountInterface, _ := record.Get("RatedMoviesCount")
	ratedMoviesCountInt64 := ratedMoviesCountInterface.(int64)

	return uint(ratedMoviesCountInt64), nil
}

func (nr *Neo4jRepository) GetSimilairUsersAndTheirAvgRating(
	userId string,
) (domain.UsersToCompare, error) {
	emptyUserToCompare := domain.UsersToCompare{}

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
		return emptyUserToCompare, err
	}

	usersToComp := domain.UsersToCompare{}

	for res.Next() {
		rec := res.Record()
		userToRecAvgRating, _ := rec.Get("UserToRecAvgRating")
		userToCompAvgRating, _ := rec.Get("UserToCompAvgRating")

		userToCompareId, _ := rec.Get("UserToCompareId")
		ratingsInCommonInterface, _ := rec.Get("RatingsInCommon")

		ratingsInCommonInterfaceSlice := ratingsInCommonInterface.([]interface{})

		usersToComp[userToCompareId.(string)] = &domain.UserToCompare{
			RatingsInCommon:     convertRatingsInCommonInterfaceSlice(ratingsInCommonInterfaceSlice),
			UserToRecAvgRating:  userToRecAvgRating.(float64),
			UserToCompAvgRating: userToCompAvgRating.(float64),
		}
	}

	if len(usersToComp) == 0 {
		return emptyUserToCompare, errors.New("there are no similiar user to you yet, keep rating some more")
	}

	return usersToComp, nil
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
	emptyScoringMovies := domain.ScoringMovies{}

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
		return emptyScoringMovies, err
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

		emptyScoringMovies[movieId] = &domain.Details{
			Movie: domain.Movie{
				Title:        movieTitle,
				ReleaseYear:  int(releaseYearInt64),
				MoreInfoLink: moreInfoLink.(string),
			},
		}

		ratings := convertRatedMoviesInterfaceSlice(userRatingCollectionInterfaceSlice)

		emptyScoringMovies[movieId].Ratings = ratings
	}

	return emptyScoringMovies, nil
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

func (nr *Neo4jRepository) GetAllLikedMovies(userId string) (domain.UsersLikedMovies, error) {
	emptyLikedMovies := domain.UsersLikedMovies{}

	session := nr.Driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
	match (u:User {userId:$userId})-[r:RATED{rating:1}]->(m:Movie)
	return m.movieId as MovieId
	`
	parameters := map[string]interface{}{"userId": userId}

	res, err := session.Run(query, parameters)
	if err != nil {
		return emptyLikedMovies, err
	}

	for res.Next() {
		rec := res.Record()
		movieIdInterface, _ := rec.Get("MovieId")
		movieId := movieIdInterface.(string)

		emptyLikedMovies[movieId] = domain.UsersLikedMovie{AllCast: map[string]bool{}}
	}

	return emptyLikedMovies, nil
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
		return err
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
		return emptySimilarMovies, err
	}

	for res.Next() {
		rec := res.Record()
		movieIdInterface, _ := rec.Get("MovieId")
		movieId := movieIdInterface.(string)

		emptySimilarMovies[movieId] = domain.SimilarMovieToLikedOnes{}
	}

	return emptySimilarMovies, nil
}

func (nr *Neo4jRepository) AddToWatchlist(userId, movieId string, timeOfAdding int64) error {
	session := nr.Driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	query := `
	match (u:User{userId:$userId}), (m:Movie{movieId:$movieId})
    merge (u)-[w:WATCH_LATER]->(m) on create set w.createdAt=$timeOfAdding
	return m.movieId as MovieId
	`
	parameters := map[string]interface{}{"userId": userId, "movieId": movieId, "timeOfAdding": timeOfAdding}
	res, err := session.Run(query, parameters)
	if err != nil {
		return err
	}

	if _, err := res.Single(); err != nil {
		return err
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
	res, err := session.Run(query, parameters)
	if err != nil {
		return err
	}

	if _, err := res.Single(); err != nil {
		return err
	}

	return nil
}

func (nr *Neo4jRepository) GetWatchlist(userId string) (domain.UserWatchlist, error) {
	emptyWatchlist := domain.UserWatchlist{}

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
		return emptyWatchlist, err
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

		emptyWatchlist = append(emptyWatchlist, watchlistItem)
	}

	if len(emptyWatchlist) == 0 {
		return emptyWatchlist, errors.New("there are no more movies in your watchlist")
	}

	return emptyWatchlist, nil
}
