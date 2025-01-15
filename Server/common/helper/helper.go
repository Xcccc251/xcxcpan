package Server_Helper

import (
	"crypto/md5"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/jordan-wright/email"
	uuid "github.com/satori/go.uuid"
	"math/rand"
	"net/smtp"
	"os"
	"strconv"
	"time"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var src = rand.NewSource(time.Now().UnixNano())

const (
	// 6 bits to represent a letter index
	letterIdBits = 6
	// All 1-bits as many as letterIdBits
	letterIdMask = 1<<letterIdBits - 1
	letterIdMax  = 63 / letterIdBits
)

type UserClaims struct {
	UserId int `json:id`
	jwt.StandardClaims
}

var myKey = []byte("gin-gorm-take_out")

func GetMd5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

func GenerateToken(id int) (string, error) {
	UserClaim := &UserClaims{
		UserId: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 30).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim)
	tokenString, err := token.SignedString(myKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func AnalysisToken(tokenString string) (*UserClaims, error) {
	userClaim := new(UserClaims)
	claims, err := jwt.ParseWithClaims(tokenString, userClaim, func(token *jwt.Token) (interface{}, error) {
		return myKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !claims.Valid {
		return nil, err
	}
	return userClaim, err

}

// 用正则表达式校验邮箱格式
func IsValidEmail(email string) bool {
	// 定义邮箱的正则表达式模式
	emailRegex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	// 编译正则表达式
	re := regexp.MustCompile(emailRegex)
	// 返回匹配结果
	return re.MatchString(email)
}
func GetUUID() string {
	return uuid.NewV4().String()
}

func SendCode(toUserEmail, code string) error {
	e := email.NewEmail()
	e.From = "Welcome <18516267632@163.com>"
	e.To = []string{toUserEmail}
	e.Subject = "邮箱验证码"
	e.HTML = []byte("您的验证码是：<b>" + code + "</b>，验证码有效期15分钟")
	return e.SendWithTLS("smtp.163.com:465",
		smtp.PlainAuth("", "18516267632@163.com", "DUdBaN4pAtAXFB9u", "smtp.163.com"),
		&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.163.com"})
}

func GetRandomStr(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdMax letters!
	for i, cache, remain := n-1, src.Int63(), letterIdMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdMax
		}
		if idx := int(cache & letterIdMask); idx < len(letters) {
			b[i] = letters[idx]
			i--
		}
		cache >>= letterIdBits
		remain--
	}
	return string(b)
}

func GetRandNumber() string {
	rand.Seed(time.Now().UnixNano())
	s := ""
	for i := 0; i < 6; i++ {
		s += strconv.Itoa(rand.Intn(10))
	}
	return s
}

func CodeSave(code []byte) (string, error) {
	dirName := "code/" + GetUUID()
	path := dirName + "/Run_user.go"
	err := os.Mkdir(dirName, 0777)
	if err != nil {
		return "", err
	}
	file, err := os.Create(path)
	if err != nil {
		return "", err
	}
	file.Write(code)
	defer file.Close()
	return path, nil

}

var result struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

func MapToJsonString(m map[string]string) (string, error) {
	jsonByte, err := json.Marshal(m)
	if err != nil {
		fmt.Printf("Marshal with error: %+v\n", err)
		return "", nil
	}

	return string(jsonByte), nil

}

func GetAddressByIp(ip string) string {
	resp, err := http.Get(fmt.Sprintf("https://api.ipify.org?format=json&ip=%s", ip))
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	return string(body)

}

func StripMarkdown(markdown string) string {
	// 去掉标题
	re := regexp.MustCompile(`(?m)^\s*#.*$`)
	markdown = re.ReplaceAllString(markdown, "")

	// 去掉加粗
	re = regexp.MustCompile(`\*\*(.*?)\*\*`)
	markdown = re.ReplaceAllString(markdown, "$1")

	// 去掉斜体
	re = regexp.MustCompile(`\*(.*?)\*`)
	markdown = re.ReplaceAllString(markdown, "$1")

	// 去掉行内代码
	re = regexp.MustCompile("`([^`]*)`")
	markdown = re.ReplaceAllString(markdown, "$1")

	// 去掉删除线
	re = regexp.MustCompile(`~~(.*?)~~`)
	markdown = re.ReplaceAllString(markdown, "$1")

	// 去掉链接
	re = regexp.MustCompile(`\[(.*?)\]\(.*?\)`)
	markdown = re.ReplaceAllString(markdown, "$1")

	// 去掉图片
	re = regexp.MustCompile(`!\[.*?\]\(.*?\)`)
	markdown = re.ReplaceAllString(markdown, "")

	// 去掉引用
	re = regexp.MustCompile(`>\s?`)
	markdown = re.ReplaceAllString(markdown, "")

	// 去掉无序列表
	re = regexp.MustCompile(`(?m)^\s*[-*+]\s+`)
	markdown = re.ReplaceAllString(markdown, "")

	// 去掉有序列表
	re = regexp.MustCompile(`(?m)^\s*\d+\.\s+`)
	markdown = re.ReplaceAllString(markdown, "")

	// 去掉换行符
	markdown = strings.ReplaceAll(markdown, "\n", " ")

	return markdown
}

func Contains(arr []string, value string) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}
func DiffOfIntArray(array1, array2 []int) []int {
	// 创建一个 map 来存储 array2 的元素
	elements := make(map[int]bool)
	for _, v := range array2 {
		elements[v] = true
	}

	// 遍历 array1，找出不在 array2 中的元素
	var diff []int
	for _, v := range array1 {
		if !elements[v] {
			diff = append(diff, v)
		}
	}

	return diff
}

func GetShuffle[T any](arr []T) []T {
	newarr := arr
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(newarr), func(i, j int) {
		newarr[i], newarr[j] = newarr[j], newarr[i]
	})
	return newarr

}
