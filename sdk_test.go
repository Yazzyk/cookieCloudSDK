package cookieCloudSDK

import "testing"

func TestNewCookieCloudSDK(t *testing.T) {
	sdk, err := NewCookieCloudSDK("", "", "")
	if err != nil {
		t.Error(err)
		return
	}
	data, err := sdk.GetCookie()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(data)
}
