package token

import "github.com/shanelex111/go-common/third_party/geo"

type CacheToken struct {
	Account *CacheTokenAccount `json:"account"`
	Access  *CacheTokenAccess  `json:"access"`
	Device  *CacheTokenDevice  `json:"device"`
	Geo     *geo.GeoCity       `json:"geo"`
}

type CacheTokenAccount struct {
	ID uint `json:"id"`
}
type CacheTokenAccess struct {
	Token            string `json:"token"`
	ExpiredAt        int64  `json:"expired_at"`
	Refresh          string `json:"refresh"`
	RefreshExpiredAt int64  `json:"refresh_expired_at"`
}

type CacheTokenDevice struct {
	ID         uint   `json:"id"`
	Type       string `json:"type"`
	Model      string `json:"model"`
	AppVersion int    `json:"app_version"`
}
