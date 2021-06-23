package domain_test

import (
	"reflect"
	"testing"

	"github.com/xhuliodo/couch-potatoes/clean-api/domain"
)

func TestFilterBasedOnRatingCount(t *testing.T) {
	var r float64 = 0
	got := domain.UsersToCompare{
		"1": &domain.UserToCompare{
			RatingsInCommon: []domain.RatingInCommon{
				{r, r}, {r, r}, {r, r}, {r, r}, {r, r},
				{r, r}, {r, r}, {r, r}, {r, r}, {r, r},
				{r, r}, {r, r}, {r, r}, {r, r}, {r, r},
			},
		},
		"2": &domain.UserToCompare{
			RatingsInCommon: []domain.RatingInCommon{
				{r, r}, {r, r}, {r, r}, {r, r},
			},
		},
		"3": &domain.UserToCompare{
			RatingsInCommon: []domain.RatingInCommon{
				{r, r}, {r, r}, {r, r}, {r, r}, {r, r},
				{r, r}, {r, r}, {r, r}, {r, r}, {r, r},
				{r, r}, {r, r}, {r, r}, {r, r}, {r, r},
			},
		},
	}

	want := domain.UsersToCompare{
		"1": &domain.UserToCompare{
			RatingsInCommon: []domain.RatingInCommon{
				{r, r}, {r, r}, {r, r}, {r, r}, {r, r},
				{r, r}, {r, r}, {r, r}, {r, r}, {r, r},
				{r, r}, {r, r}, {r, r}, {r, r}, {r, r},
			},
		},
		"3": &domain.UserToCompare{
			RatingsInCommon: []domain.RatingInCommon{
				{r, r}, {r, r}, {r, r}, {r, r}, {r, r},
				{r, r}, {r, r}, {r, r}, {r, r}, {r, r},
				{r, r}, {r, r}, {r, r}, {r, r}, {r, r},
			},
		},
	}

	got.FilterBasedOnRatingsCount()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}

}

var userToRecAvgRating = 0.5

func TestCalculatePearson(t *testing.T) {
	usersToCompare := make(domain.UsersToCompare)
	populateWithStubUsersToCompare(&usersToCompare)

	usersToCompare.CalculatePearson()

	want := domain.UsersToCompare{
		"1": &domain.UserToCompare{
			PearsonCoefficient: 0.8528028654224418,
		},
		"2": &domain.UserToCompare{
			PearsonCoefficient: -0.33806170189140655,
		},
		"3": &domain.UserToCompare{
			PearsonCoefficient: 0.5619514869490164,
		},
	}

	for userId, user := range usersToCompare {
		if user.PearsonCoefficient != want[userId].PearsonCoefficient {
			t.Errorf("\n got: %+v,\n want: %+v", user.PearsonCoefficient, want[userId].PearsonCoefficient)
		}
	}

}

func populateWithStubUsersToCompare(utc *domain.UsersToCompare) {
	usersId := []string{"1", "2", "3"}
	userToCompDetails := []domain.UserToCompare{
		{UserToRecAvgRating: userToRecAvgRating,
			UserToCompAvgRating: 0.2,
			RatingsInCommon: []domain.RatingInCommon{
				{0.7, 0.5}, {0.5, 0.1}, {0.2, 0}, {0.4, 0}, {0.6, 0.4}, {0.6, 0.2},
				{0.7, 0.5}, {0.5, 0.1}, {0.2, 0}, {0.4, 0}, {0.6, 0.4}, {0.6, 0.2},
			},
		},
		{UserToRecAvgRating: userToRecAvgRating,
			UserToCompAvgRating: 0.8,
			RatingsInCommon: []domain.RatingInCommon{
				{0.5, 0.7}, {0.2, 1}, {0.7, 0.9}, {0.6, 0.8}, {0.5, 0.6},
				{0.5, 0.7}, {0.2, 1}, {0.7, 0.9}, {0.6, 0.8}, {0.5, 0.6},
			},
		},
		{UserToRecAvgRating: userToRecAvgRating,
			UserToCompAvgRating: 0.5,
			RatingsInCommon: []domain.RatingInCommon{
				{0.6, 0.5}, {0.5, 0.5}, {0.2, 0.2}, {0.6, 0.3}, {0.6, 1},
				{0.6, 0.5}, {0.5, 0.5}, {0.2, 0.2}, {0.6, 0.3}, {0.6, 1},
			},
		},
	}
	for i, userId := range usersId {
		(*utc)[userId] = &userToCompDetails[i]
	}
}
