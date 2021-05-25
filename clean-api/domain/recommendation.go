package domain

import (
	"errors"
	"math"
)

type UserComparison struct {
	UserId         string
	UserAvgRating  float64
	UsersToCompare map[string]UserToCompare
}

type UserToCompare struct {
	UserId             string
	UserAvgRating      float64
	RatingsInCommon    []RatingInCommon
	PearsonCoefficient float64
}

type RatingInCommon struct {
	UserRating          float64
	UserToCompareRating float64
}

type UserBasedRecommendation struct {
	Movie
	CorrelationCoefficient float64
	Score                  float64
}

const requiredRatingsCompatibility uint = 10

func (uc *UserComparison) FilterBasedOnRatingsCount() error {
	usersToCompare := uc.UsersToCompare
	for i, userToCompare := range usersToCompare {
		ratingsInCommonCount := len(userToCompare.RatingsInCommon)
		if ratingsInCommonCount < int(requiredRatingsCompatibility) {
			delete(uc.UsersToCompare, i)
		}
	}
	if len(uc.UsersToCompare) < 1 {
		return errors.New("no similar users were found")
	}
	return nil
}

func (uc *UserComparison) CalculatePearson() error {
	userAvgRating := uc.UserAvgRating
	for i, user := range uc.UsersToCompare {
		userToCompareAvgRating := user.UserAvgRating
		var nom float64
		var denomUser float64
		var denomUserToCompare float64
		for _, rating := range user.RatingsInCommon {
			nom += (rating.UserRating - userAvgRating) * (rating.UserToCompareRating - userToCompareAvgRating)
			denomUser += math.Pow((rating.UserRating - userAvgRating), 2)
			denomUserToCompare += math.Pow((rating.UserToCompareRating - userToCompareAvgRating), 2)
		}
		denom := math.Sqrt(denomUser * denomUserToCompare)
		if denom != 0 {
			pearsonCoefficient := nom / denom
			user.PearsonCoefficient = pearsonCoefficient
		} else {
			delete(uc.UsersToCompare, i)
		}
	}
	return nil
}
