package domain

import (
	"errors"
	"fmt"
	"math"
)

type UserToRecommend struct {
	UserId        string
	UserAvgRating float64
}

type UsersToCompare map[string]*UserToCompare

type UserToCompare struct {
	UserAvgRating      float64
	RatingsInCommon    []RatingInCommon
	PearsonCoefficient float64
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

func (uc *UsersToCompare) CalculatePearson(userToRecommend *UserToRecommend) error {
	userAvgRating := userToRecommend.UserAvgRating
	fmt.Println("u1_mean", userAvgRating)
	for i, user := range *uc {
		userToCompareAvgRating := user.UserAvgRating
		var nom float64
		var denomUserToRecommend float64
		var denomUserToCompare float64
		for _, rating := range user.RatingsInCommon {
			nom += (rating.UserToRecommendRating - userAvgRating) * (rating.UserToCompareRating - userToCompareAvgRating)
			denomUserToRecommend += math.Pow(rating.UserToRecommendRating-userAvgRating, 2)
			denomUserToCompare += math.Pow(rating.UserToCompareRating-userToCompareAvgRating, 2)
		}
		denom := math.Sqrt(denomUserToCompare * denomUserToCompare)
		if denom != 0 {
			pearsonCoefficient := nom / denom
			user.PearsonCoefficient = pearsonCoefficient
			fmt.Println("denom", denom)
			fmt.Println("corresponding pearson", pearsonCoefficient)
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
