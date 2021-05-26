package domain

import (
	"errors"
	"math"
)

type UsersToCompare map[string]*UserToCompare

type UserToCompare struct {
	UserToRecAvgRating  float64
	UserToCompAvgRating float64
	RatingsInCommon     []RatingInCommon
	PearsonCoefficient  float64
}

type RatingInCommon struct {
	UserToRecommendRating float64
	UserToCompareRating   float64
}

type UsersBasedRecommendation []UserBasedRecommendation

type UserBasedRecommendation struct {
	Movie
	Score float64
}

const requiredRatingsCompatibility int = 10

func (utc *UsersToCompare) FilterBasedOnRatingsCount() error {
	for i, userToCompare := range *utc {
		ratingsInCommonCount := len(userToCompare.RatingsInCommon)
		if ratingsInCommonCount < requiredRatingsCompatibility {
			delete(*utc, i)
		}
	}
	if len(*utc) < 1 {
		return errors.New("no similar users were found")
	}
	return nil
}

func (uc *UsersToCompare) CalculatePearson() error {
	for i, user := range *uc {
		userToRecAvgRating := user.UserToRecAvgRating
		userToCompareAvgRating := user.UserToCompAvgRating

		var nom float64
		n := len(user.RatingsInCommon) - 1
		var denomUserToRecommend float64
		var denomUserToCompare float64
		for _, rating := range user.RatingsInCommon {
			nom += (rating.UserToRecommendRating - userToRecAvgRating) * (rating.UserToCompareRating - userToCompareAvgRating)
			denomUserToRecommend += math.Pow(rating.UserToRecommendRating-userToRecAvgRating, 2)
			denomUserToCompare += math.Pow(rating.UserToCompareRating-userToCompareAvgRating, 2)
		}
		sX := denomUserToCompare / float64(n)
		sX = math.Sqrt(sX)
		sY := denomUserToRecommend / float64(n)
		sY = math.Sqrt(sY)
		denom := sX * sY * float64(n)

		if denom != 0 {
			pearsonCoefficient := nom / denom
			user.PearsonCoefficient = pearsonCoefficient
		} else {
			delete(*uc, i)
		}
	}
	return nil
}

func (uc *UsersToCompare) RemoveLowPearson(remainingUserIds *[]string) error {
	for _, u := range *remainingUserIds {
		_, ok := (*uc)[u]
		if !ok {
			delete(*uc, u)
		}
	}
	return nil
}
