package inmemory

import "fmt"

var AdsKey = "ads"

type UserAds string

func (a UserAds) GetKey() string {
	return fmt.Sprintf("ads_user_id_%v", a)
}
