package application

import (
	"errors"
	"sort"

	"github.com/xhuliodo/couch-potatoes/clean-api/domain"
)

type UserBasedRecommendationService struct {
	repo domain.Repository
}

func NewUserBasedRecommendationService(repo domain.Repository) UserBasedRecommendationService {
	return UserBasedRecommendationService{repo}
}

const (
	nearestNeighborCircle uint = 25
)

func (ubrs UserBasedRecommendationService) GetUserBasedRecommendation(userId string, limit uint) ([]domain.UserBasedRecommendation, error) {
	emptyRec := []domain.UserBasedRecommendation{}

	if _, err := ubrs.repo.GetUserById(userId); err != nil {
		return emptyRec, errors.New("a user with this identifier does not exist")
	}

	usersToCompare, err := ubrs.repo.GetAvgRatingAndCollectSimilairUsers(userId)
	if err != nil {
		return emptyRec, err
	}

	if err := usersToCompare.FilterBasedOnRatingsCount(); err != nil {
		return emptyRec, err
	}

	usersLeftIds := getIds(usersToCompare)

	if err := ubrs.repo.GetUsersAvgRating(usersLeftIds, &usersToCompare); err != nil {
		return emptyRec, err
	}

	if err := usersToCompare.CalculatePearson(); err != nil {
		return emptyRec, err
	}

	usersSorted, _ := SortByPearsonDesc(usersToCompare)

	similairUser := usersSorted[:nearestNeighborCircle]

	userIds := getIdsFromUsersSorted(similairUser)

	usersAndRatedMovies, err := ubrs.repo.GetRatedMoviesForUsers(userIds)
	if err != nil {
		return emptyRec, err
	}

	recommendationsWithNoMovieDetails := calculateScore(usersAndRatedMovies, usersToCompare)

	sort.SliceStable(recommendationsWithNoMovieDetails, func(i, j int) bool {
		return recommendationsWithNoMovieDetails[i].Score > recommendationsWithNoMovieDetails[j].Score
	})

	rwnd := recommendationsWithNoMovieDetails[:limit]

	

	return emptyRec, nil
}

func calculateScore(users []domain.User, otherUsers domain.UserComparison) []domain.UserBasedRecommendation {
	recs := []domain.UserBasedRecommendation{}
	for _, user := range users {
		for _, movie := range user.RatedMovies {
			pearson := otherUsers.UsersToCompare[user.Id].PearsonCoefficient
			score := pearson * movie.Rating
			newRec := domain.UserBasedRecommendation{
				Movie:                  movie.Movie,
				CorrelationCoefficient: pearson,
				Score:                  score,
			}
			recs = append(recs, newRec)
		}
	}

	return recs
}

func getIds(userComparison domain.UserComparison) []string {
	usersIds := []string{}
	for key := range userComparison.UsersToCompare {
		usersIds = append(usersIds, key)
	}
	return usersIds
}

type UserComparisonSortable struct {
	UserId             string
	PearsonCoefficient float64
}

type UsersComparisonSortable []UserComparisonSortable

func SortByPearsonDesc(usersToCompare domain.UserComparison) (UsersComparisonSortable, error) {
	u := make(UsersComparisonSortable, len(usersToCompare.UsersToCompare))

	i := 0
	for k, v := range usersToCompare.UsersToCompare {
		u[i] = UserComparisonSortable{k, v.PearsonCoefficient}
		i++
	}

	sort.SliceStable(u, func(i, j int) bool {
		return u[i].PearsonCoefficient > u[j].PearsonCoefficient
	})

	return u, nil
}

func getIdsFromUsersSorted(users UsersComparisonSortable) []string {
	usersIds := []string{}
	for _, u := range users {
		usersIds = append(usersIds, u.UserId)
	}
	return usersIds
}
