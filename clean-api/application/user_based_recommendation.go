package application

import (
	"sort"

	"github.com/pkg/errors"
	"github.com/xhuliodo/couch-potatoes/clean-api/domain"
)

type UserBasedRecommendationService struct {
	repo domain.Repository
}

func NewUserBasedRecommendationService(repo domain.Repository) UserBasedRecommendationService {
	return UserBasedRecommendationService{repo}
}

const (
	nearestNeighborCircle uint   = 25
	noSimilarUsersYet     string = "there are no similiar user to you yet, keep rating some more"
)

func (ubrs UserBasedRecommendationService) GetUserBasedRecommendation(userId string, limit uint) (
	recs domain.UsersBasedRecommendation, err error,
) {
	if _, err := ubrs.repo.GetUserById(userId); err != nil {
		errStack := errors.Wrap(err, "a user with this identifier does not exist")
		return recs, errStack
	}

	usersToCompare, err := ubrs.repo.GetSimilairUsersAndTheirAvgRating(userId)
	if err != nil {
		return recs, err
	}

	if len(usersToCompare) < 1 {
		cause := errors.New("not_found")
		return recs, errors.Wrap(cause, noSimilarUsersYet)
	}

	if err := usersToCompare.FilterBasedOnRatingsCount(); err != nil {
		return recs, err
	}

	if err := usersToCompare.CalculatePearson(); err != nil {
		return recs, err
	}

	usersSorted := sortByPearsonDesc(&usersToCompare)
	end := nearestNeighborCircle
	sliceMaxLength := len(usersToCompare)
	if sliceMaxLength < int(nearestNeighborCircle) {
		end = uint(sliceMaxLength)
	}
	similairUser := usersSorted[:end]

	userIds := getIdsFromSimilairUser(&similairUser)
	if err := usersToCompare.RemoveLowPearson(&userIds); err != nil {
		return recs, err
	}

	moviesAndRatings, err := ubrs.repo.GetRatedMoviesForUsersYetToBeConsidered(userId, userIds)
	if err != nil {
		return recs, errors.Wrap(err, "could not get rated movies for users yet to be considered")

	}

	recs = calculateScore(usersToCompare, moviesAndRatings)

	sort.SliceStable(recs, func(i, j int) bool {
		return recs[i].Score > recs[j].Score
	})

	// handle pagination
	length := len(recs)
	begin, end, err := handlePagination(uint(length), defaultSkip, limit)
	if err != nil {
		return recs, err
	}
	recs = recs[begin:end]

	return recs, nil
}

type UserComparisonSortable struct {
	UserId             string
	PearsonCoefficient float64
}

type UsersComparisonSortable []UserComparisonSortable

func sortByPearsonDesc(usersToCompare *domain.UsersToCompare) UsersComparisonSortable {
	u := make(UsersComparisonSortable, len(*usersToCompare))

	i := 0
	for k, v := range *usersToCompare {
		u[i] = UserComparisonSortable{k, v.PearsonCoefficient}
		i++
	}

	sort.SliceStable(u, func(i, j int) bool {
		return u[i].PearsonCoefficient > u[j].PearsonCoefficient
	})

	return u
}

func getIdsFromSimilairUser(users *UsersComparisonSortable) []string {
	usersIds := []string{}
	for _, u := range *users {
		usersIds = append(usersIds, u.UserId)
	}
	return usersIds
}

func calculateScore(users domain.UsersToCompare, moviesAndRatings domain.ScoringMovies) domain.UsersBasedRecommendation {
	moviesScored := []domain.UserBasedRecommendation{}

	// calculate score for a movie by adding up all
	// pearson * rating for each user
	for movieId, movie := range moviesAndRatings {
		var score float64
		for userId, rating := range movie.Ratings {
			score += users[userId].PearsonCoefficient * rating
		}
		// append it into a slice that can be sorted out
		moviesScored = append(moviesScored, domain.UserBasedRecommendation{
			Movie: domain.Movie{
				Id:           movieId,
				Title:        movie.Title,
				ReleaseYear:  movie.ReleaseYear,
				MoreInfoLink: movie.MoreInfoLink,
			},
			Score: score,
		})
	}

	return moviesScored
}
