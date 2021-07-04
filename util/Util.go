package util

import (
	"crypto/sha1"
	"encoding/hex"
	"net/http"
	"strings"
)

// SHA1 sha1加密
func SHA1(s string) string {
	o := sha1.New()
	o.Write([]byte(s))
	return hex.EncodeToString(o.Sum(nil))
}

// StrCookies 将cookies转换为字符串
func StrCookies(cookies []*http.Cookie, sep string) string {
	var build strings.Builder
	for k, cookie := range cookies {
		s := cookie.String()
		if k == 0 {
			build.WriteString(s)
		} else {
			build.WriteString(sep + s)
		}
	}
	return build.String()
}

// ParseCookie 将字符串cookie解析成map格式
func ParseCookie(cookies string, sep string) map[string]string {
	mapCookie := make(map[string]string)
	cookie := strings.Split(cookies, sep)
	for _, v := range cookie {
		p := strings.Split(v, "=")
		mapCookie[p[0]] = p[1]
	}
	return mapCookie
}
