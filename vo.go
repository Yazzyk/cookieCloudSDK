package cookieCloudSDK

// 解密前的数据结构
type CookieCloudRs struct {
	Encrypted string `json:"encrypted"`
}

// 解密后的数据结构
type CookieCloudDecodeData struct {
	CookieData       map[string][]CookieData     `json:"cookie_data"`
	LocalStorageData map[string]LocalStorageData `json:"local_storage_data"`
}

type CookieData struct {
	Domain         string  `json:"domain"`
	ExpirationDate float64 `json:"expirationDate"`
	HostOnly       bool    `json:"hostOnly"`
	HttpOnly       bool    `json:"httpOnly"`
	Name           string  `json:"name"`
	Path           string  `json:"path"`
	SameSite       string  `json:"sameSite"`
	Secure         bool    `json:"secure"`
	Session        bool    `json:"session"`
	StoreId        string  `json:"storeId"`
	Value          string  `json:"value"`
}

type LocalStorageData struct {
	IsWhitelist string `json:"isWhitelist"`
}
