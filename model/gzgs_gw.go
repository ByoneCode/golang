package model

import (
	"errors"
	"golang-student/util"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var (
	gwUrl         string
	cookieStudent []*http.Cookie
	cookieJar     *cookiejar.Jar
)

// Score 成绩结构体
type Score struct {
	Semester      string `json:"semester"`
	Kcode         string `json:"kcode"`
	Number        string `json:"number"`
	Kname         string `json:"kname"`
	Ktype         string `json:"ktype"`
	Credit        string `json:"credit"`
	Grades        string `json:"grades"`
	FinalGrade    string `json:"final_grade"`
	OverallRating string `json:"overall_rating"`
	Performance   string `json:"performance"`
}

// ScoreData 成绩数据结构体
type ScoreData struct {
	Items []*Score `json:"items"`
}

// Login 登录结构体
type Login struct {
	Code  int    `json:"code"`
	Token string `json:"token"`
}

// 参数初始化
func init() {
	gwUrl = "http://jw.gzgs.org.cn"
	cookieStudent = nil
	cookieJar, _ = cookiejar.New(nil)
}

// getToken 获取动态值
func getToken(contents io.Reader) string {
	data, _ := ioutil.ReadAll(contents)
	regCode, err := regexp.Compile(`CryptoJS.SHA1\('(.*?)'`)
	if err != nil {
		log.Fatal("token解析错误")
	}
	// 提取到token
	token := regCode.FindStringSubmatch(string(data))[1]
	return token
}

// LoginStudent 模拟登录系统获取cookie
func LoginStudent(username, password string) (*Login, error) {
	var loginUrl = gwUrl + "/eams/login.action"
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Jar:     cookieJar,
		Timeout: time.Second * 3,
	}
	req, _ := http.NewRequest("GET", loginUrl, nil)
	res, _ := client.Do(req)
	token := getToken(res.Body)
	// 休眠500毫秒防止刷新过快
	time.Sleep(time.Millisecond * 500)
	// post数据
	postData := url.Values{}
	postData.Add("username", username)
	postData.Add("password", util.SHA1(token+password))
	postData.Add("encodedPassword", "")
	postData.Add("session_locale", "zh_CN")
	body := ioutil.NopCloser(strings.NewReader(postData.Encode()))
	requ, _ := http.NewRequest("POST", loginUrl, body)

	requ.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	requ.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36 Edg/91.0.864.59")
	requ.Header.Set("Host", "gw.gzgs.org.cn")
	requ.Header.Set("Origin", "http://jw.gzgs.org.cn")
	requ.Header.Set("Referer", "http://jw.gzgs.org.cn/eams/login.action")
	requ.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	requ.Header.Set("Proxy-Connection", "keep-alive")

	res, _ = client.Do(requ)
	logins, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	//获取登录结果
	msg := regexp.MustCompile(`(?sm)<div .+?>\n+\s+<span>(.*?)</span>\s+</div>`)
	var login Login
	if res.StatusCode == 302 {
		cookieStudent = cookieJar.Cookies(res.Request.URL)
		login = Login{
			Code:  res.StatusCode,
			Token: util.StrCookies(cookieStudent, ","),
		}
	} else {
		resMsg := msg.FindStringSubmatch(string(logins))
		return nil, errors.New(resMsg[1])
	}
	return &login, nil
}

// GetScore 获取成绩
func GetScore(token, semesterId string) (*ScoreData, error) {
	//fmt.Println(token)
	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: time.Second * 10,
	}
	score, _ := http.NewRequest(http.MethodGet, gwUrl+"/eams/teach/grade/course/person!search.action?semesterId="+semesterId+"&projectType=", nil)
	cookies := util.ParseCookie(token, ",")
	for key := range cookies {
		score.AddCookie(&http.Cookie{
			Name:  key,
			Value: cookies[key],
		})
	}
	//fmt.Println(score.Header)
	res, err := client.Do(score)
	if err != nil {
		panic(err)
	}

	h, _ := ioutil.ReadAll(res.Body)
	//fmt.Println(res.Status)
	if res.StatusCode != 302 {
		//正则取出表格主体
		regb := regexp.MustCompile(`(?ms)<tbody id="\w+">(.*?)</tbody>`)
		tbody := regb.FindAllStringSubmatch(string(h), -1)[0][0]
		//获取出每个tr
		regr := regexp.MustCompile(`(?sm)<tr>(.*?)</tr>`)
		tr := regr.FindAllString(tbody, -1)
		isData, _ := regexp.MatchString(`(?sm)<tr>\s+?</tr>`, tr[0])
		if isData {
			return nil, errors.New("未查询到数据")
		} else {
			scoreData := new(ScoreData)
			for _, v := range tr {
				// 取出每个数值
				regd := regexp.MustCompile(`(?sm)<td.*?>[\s]*(.*?)[\n\t\r]*</td>`)
				td := regd.FindAllStringSubmatch(v, -1)

				scoreData.Items = append(scoreData.Items, &Score{
					Semester:      td[0][1],
					Kcode:         td[1][1],
					Number:        td[2][1],
					Kname:         td[3][1],
					Ktype:         td[4][1],
					Credit:        td[5][1],
					Grades:        td[6][1],
					FinalGrade:    td[7][1],
					OverallRating: td[8][1],
					Performance:   td[9][1],
				})
			}
			res.Body.Close()
			return scoreData, nil
		}
	}
	return nil, errors.New("你的身份已过期请重新登录")
}
