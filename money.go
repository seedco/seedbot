package main

import (
	"strconv"
	"strings"
)

func CentsToDollarStringWithCommas(cents int64) string {
	value := CentsToDollarString(cents)
	decimalSplit := strings.Split(value, ".")
	if len(decimalSplit) < 1 {
		return value
	}
	digits := []string{}
	for len(decimalSplit[0]) > 0 {
		remainder := len(decimalSplit[0]) % 3
		if remainder == 0 && len(decimalSplit[0]) >= 3 {
			remainder += 3
		}
		digits = append(digits, decimalSplit[0][:remainder])
		decimalSplit[0] = decimalSplit[0][remainder:]
	}
	decimalSplit[0] = strings.Join(digits, ",")
	return strings.Join(decimalSplit, ".")
}

func CentsToDollarString(cents int64) string {
	str := strconv.FormatInt(cents, 10)
	sign := (str[0] != '-')
	str = strings.Trim(str, "-")
	switch len(str) {
	case 1:
		str = "0.0" + str
	case 2:
		str = "0." + str
	default:
		str = str[:len(str)-2] + "." + str[len(str)-2:]
	}
	if sign {
		return str
	}
	return "-" + str
}
