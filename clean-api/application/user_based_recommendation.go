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

	usersSorted, _ := sortByPearsonDesc(&usersToCompare)
	// fmt.Println(usersSorted)
	end := nearestNeighborCircle
	sliceMaxLength := len(usersToCompare)
	if sliceMaxLength < int(nearestNeighborCircle) {
		end = uint(sliceMaxLength)
	}
	similairUser := usersSorted[:end]

	userIds := getIdsFromSimilairUser(&similairUser)
	usersToCompare.RemoveLowPearson(&userIds)

	moviesAndRatings, err := ubrs.repo.GetRatedMoviesForUsersYetToBeConsidered(userId, userIds)
	if err != nil {
		return emptyRec, err
	}

	recs := calculateScore(usersToCompare, moviesAndRatings)

	sort.SliceStable(recs, func(i, j int) bool {
		return recs[i].Score > recs[j].Score
	})

	// handle pagination
	length := len(recs)
	begin, end, err := handlePagination(uint(length), defaultSkip, limit)
	if err != nil {
		return emptyRec, errors.New("you're all caught up")
	}
	recs = recs[begin:end]

	return recs, nil
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
