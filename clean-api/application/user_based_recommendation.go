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

func (ubrs UserBasedRecommendationService) GetUserBasedRecommendation(userId string, limit uint) (domain.UsersBasedRecommendation, error) {
	emptyRec := domain.UsersBasedRecommendation{}

	if _, err := ubrs.repo.GetUserById(userId); err != nil {
		return emptyRec, errors.New("a user with this identifier does not exist")
	}

	usersToCompare, err := ubrs.repo.GetSimilairUsersAndTheirAvgRating(userId)
	if err != nil {
		return emptyRec, err
	}

	if err := usersToCompare.FilterBasedOnRatingsCount(); err != nil {
		return emptyRec, err
	}

	if err := usersToCompare.CalculatePearson(); err != nil {
		return emptyRec, err
	}

	// for k, u := range usersToCompare {
	// 	fmt.Println(k, u.UserAvgRating, u.PearsonCoefficient)
	// }

	usersSorted, _ := sortByPearsonDesc(&usersToCompare)
	// fmt.Println(usersSorted)
	end := nearestNeighborCircle
	sliceMaxLength := len(usersToCompare)
	if sliceMaxLength < int(nearestNeighborCircle) {
		end = uint(sliceMaxLength)
	}
	similairUser := usersSorted[:end]

	userIds := getIdsFromSimilairUser(&similairUser)

	usersAndRatedMovies, err := ubrs.repo.GetRatedMoviesForUsers(userIds)
	if err != nil {
		return emptyRec, err
	}

	usersToCompare.RemoveLowPearson(&userIds)

	recommendationsWithNoMovieDetails := calculateScore(&usersAndRatedMovies, usersToCompare)

	sort.SliceStable(recommendationsWithNoMovieDetails, func(i, j int) bool {
		return recommendationsWithNoMovieDetails[i].Score > recommendationsWithNoMovieDetails[j].Score
	})
	recommendationsWithNoMovieDetails = recommendationsWithNoMovieDetails[:limit]

	return recommendationsWithNoMovieDetails, nil
}

func getIds(userComparison *domain.UsersToCompare) []string {
	usersIds := []string{}
	for key := range *userComparison {
		usersIds = append(usersIds, key)
	}
	return usersIds
}

type UserComparisonSortable struct {
	UserId             string
	PearsonCoefficient float64
}

type UsersComparisonSortable []UserComparisonSortable

func sortByPearsonDesc(usersToCompare *domain.UsersToCompare) (UsersComparisonSortable, error) {
	u := make(UsersComparisonSortable, len(*usersToCompare))

	i := 0
	for k, v := range *usersToCompare {
		u[i] = UserComparisonSortable{k, v.PearsonCoefficient}
		i++
	}

	sort.SliceStable(u, func(i, j int) bool {
		return u[i].PearsonCoefficient > u[j].PearsonCoefficient
	})

	return u, nil
}

func getIdsFromSimilairUser(users *UsersComparisonSortable) []string {
	usersIds := []string{}
	for _, u := range *users {
		usersIds = append(usersIds, u.UserId)
	}
	return usersIds
}

func calculateScore(users *[]domain.User, otherUsers domain.UsersToCompare) domain.UsersBasedRecommendation {
	recs := domain.UsersBasedRecommendation{}
	for _, user := range *users {
		pearson := otherUsers[user.Id].PearsonCoefficient
		for _, movie := range user.RatedMovies {
			score := pearson * movie.Rating
			newRec := domain.UserBasedRecommendation{
				Movie: movie.Movie,
				Score: score,
			}
			recs = append(recs, newRec)
		}
	}

	return recs
}
