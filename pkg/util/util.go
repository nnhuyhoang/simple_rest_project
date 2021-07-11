package util

import (
	"bytes"
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/consts"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/errors"
)

// ValidateEmail validate a string is email by regular expression
func ValidateEmail(email string) bool {
	Re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return Re.MatchString(email)
}

// ValidatePasswordString validate password
func ValidatePasswordString(val string) error {
	val = strings.TrimSpace(val)
	if val == "" {
		return errors.ErrInvalidPassword
	}

	if len(val) < 6 {
		return errors.ErrInvalidPassword
	}

	return nil
}

func ValidatePasswordMatch(password string, rePassword string) error {
	if password != rePassword {
		return errors.ErrPasswordNotMatch
	}
	return nil
}

func ValidatePhoneNumber(phoneNumber string) bool {
	Re := regexp.MustCompile(`^\d+$`)
	return len(phoneNumber) > 0 && Re.MatchString(phoneNumber)
}

func ParseDateRequest(dateString string, tz string) (time.Time, error) {
	timeSplitStr := strings.Split(dateString, "-")
	if len(timeSplitStr) != 3 {
		return time.Now(), errors.ErrInvalidDateFormat
	}
	var timeSplit []int
	for _, v := range timeSplitStr {
		timeValue, err := strconv.Atoi(v)
		if err != nil {
			return time.Now(), errors.ErrInvalidDateFormat
		}
		timeSplit = append(timeSplit, timeValue)
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		loc, _ = time.LoadLocation(consts.DefaultTimezone)
	}
	date := time.Date(timeSplit[0], time.Month(timeSplit[1]), timeSplit[2], 0, 0, 0, 0, loc)
	return date, nil
}

func SerializeStruct(msg map[string]interface{}) ([]byte, error) {
	var b bytes.Buffer
	encoder := json.NewEncoder(&b)
	err := encoder.Encode(msg)
	return b.Bytes(), err
}
