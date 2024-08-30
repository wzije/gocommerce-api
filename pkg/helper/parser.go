package helper

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"net/mail"
	"regexp"
	"strconv"
	"strings"
)

func NormalizeString(str string) string {
	return strings.TrimSpace(strings.ToLower(str))
}

func IsEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(strings.ReplaceAll(snake, " ", ""))
}

func SnakeToCamel(inputUnderScoreStr string) (camelCase string) {
	//snake_case to camelCase

	isToUpper := false

	for k, v := range inputUnderScoreStr {
		if k == 0 {
			camelCase = strings.ToUpper(string(inputUnderScoreStr[0]))
		} else {
			if isToUpper {
				camelCase += strings.ToUpper(string(v))
				isToUpper = false
			} else {
				if v == '_' {
					isToUpper = true
				} else {
					camelCase += string(v)
				}
			}
		}
	}
	return

}

func ToJsonString(object interface{}) string {
	jsonData, err := json.Marshal(object)

	if err != nil {
		log.Fatal(err)
	}

	return string(jsonData)
}

func AnyToInt(data interface{}) (*int, error) {
	var value int
	switch data.(type) {
	case int:
		value = data.(int)
	case float64:
		value = int(data.(float64))
	case string:
		value, _ = strconv.Atoi(data.(string))
	default:
		return nil, errors.New("conversion to integer failed, invalid value type")
	}

	return &value, nil
}

func AnyToBase64(data interface{}) string {

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	return base64.StdEncoding.EncodeToString(jsonData)
}

func StringToUint64(value string) uint64 {
	result, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	return result
}
