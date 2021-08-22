package utility

import (
	"math/rand"
	"motherbear/backend/constants"
	"net/url"
	"reflect"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
)

func IsAlphanumericString(val string) bool {

	if val == "" {
		return false
	}

	reg := regexp.MustCompile("^[a-zA-Z0-9_]*$")
	result := reg.FindStringIndex(val)
	if result == nil {
		return false
	}

	return true
}

func IsIdenticalSlice(var1, var2 []string) bool {
	if len(var1) != len(var2) {
		return false
	}

	for _, value1 := range var1 {
		check := false

		for _, value2 := range var2 {
			if value1 == value2 {
				check = true
				break
			}
		}

		if !check {
			return false
		}
	}
	return true
}

func IsExistValueInList(value, list interface{}) bool {
	switch reflect.TypeOf(list).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(list)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(value, s.Index(i).Interface()) == true {
				return true
			}
		}
	}

	return false
}

func ConvertInterfaceToStringSlice(list interface{}) []string {
	var stringSlice []string
	stringSlice = make([]string, 0)

	if reflect.Slice != reflect.TypeOf(list).Kind() {
		return nil
	}

	s := reflect.ValueOf(list)
	for i := 0; i < s.Len(); i++ {
		if !IsMatchTypeOf(s.Index(i).Interface(), "string") {
			return nil
		}

		stringSlice = append(stringSlice, s.Index(i).Interface().(string))
	}

	return stringSlice
}

func GetTypeOf(value interface{}) string {
	typeOf := reflect.TypeOf(value).String()

	return typeOf
}

func IsMatchTypeOf(value interface{}, typeOf string) bool {
	is := reflect.TypeOf(value).String() == typeOf

	return is
}

func GetOffsetListFromRequest(c *gin.Context) (int, int, error) {
	limit, err := strconv.Atoi(c.Query(constants.RequestQueryLimit))
	if err != nil {
		return -1, -1, err
	}

	offset, err := strconv.Atoi(c.Query(constants.RequestQueryOffset))
	if err != nil {
		return -1, -1, err
	}

	return offset, limit, nil
}

func GetOffsetListDateSearchFromRequest(c *gin.Context) (int, int, string, string, error) {
	limit, err := strconv.Atoi(c.Query(constants.RequestQueryLimit))
	if err != nil {
		return -1, -1, "", "", err
	}

	offset, err := strconv.Atoi(c.Query(constants.RequestQueryOffset))
	if err != nil {
		return -1, -1, "", "", err
	}

	from := c.Query(constants.RequestQueryFrom)
	to := c.Query(constants.RequestQueryTo)

	return offset, limit, from, to, nil
}

func generateRandString(n int) string {
	var letterRunes = []rune("abcdef0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func generateWalletID() string {
	return "hx" + generateRandString(40)
}

func generateBlockTxHash() string {
	return "0x" + generateRandString(68)
}

func AddHexHD(val string) string {
	h := val[:2]
	if h != "0x" {
		return "0x" + val
	}
	return val
}

func GetHostAndPortInURL(rawurl string) (string, error) {
	urlStruct, err := url.Parse(rawurl)
	if err != nil {
		return "", err
	}

	hostAndPort := urlStruct.Host

	return hostAndPort, err
}
