package token

import "github.com/shanelex111/go-common/third_party/geo"

type CacheToken struct {
	Account *CacheTokenAccount `json:"account"`
	Access  *CacheTokenAccess  `json:"access"`
	Device  *CacheTokenDevice  `json:"device"`
	Geo     *geo.GeoCity       `json:"geo"`
}

type CacheTokenAccount struct {
	ID               uint   `json:"id"`
	Email            string `json:"email"`
	PhoneCountryCode string `json:"phone_country_code"`
	PhoneNumber      string `json:"phone_number"`
}
type CacheTokenAccess struct {
	Token            string `json:"token"`
	ExpiredAt        int64  `json:"expired_at"`
	Refresh          string `json:"refresh"`
	RefreshExpiredAt int64  `json:"refresh_expired_at"`
}

type CacheTokenDevice struct {
	DeviceID    string `json:"device_id"`
	DeviceModel string `json:"device_model"`
	DeviceType  string `json:"device_type"`
	AppVersion  int    `json:"app_version"`
	CreatedAt   int64  `json:"created_at"`
}
