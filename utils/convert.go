package utils

import (
	"strconv"
	"strings"
)

func ConvertHex(v string) uint64 {
	cleaned := strings.Replace(v, "0x", "", -1)
	res, _ := strconv.ParseUint(cleaned, 16, 64)
	return uint64(res)
}

func ConvertBool(v string) bool {
	if v == "off" || v == "false" {
		return false
	}
	return true
}

func ConverArray(v string) []string {
	return strings.Split(v, " ")
}

func ConvertInt(v string) int {
	if i, err := strconv.Atoi(v); err != nil {
		return 0
	} else {
		return i
	}
}
