# CookieCloudSDK for Go
项目主要目的是为了方便使用Go语言开发的服务与CookieCloud进行交互，目前只写了解密部分，加密部分还没写，后续会尝试完善。

## 使用方法
### 1. 导入包
```shell
    go get -u github.com/yazzyk/cookieCloudSDK
```
### 2. 使用
```go
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
```

本项目参考代码: [aienvoy](https://github.com/Vaayne/aienvoy)