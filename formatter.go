package linkcxo

import "strconv"

type Formatter struct {
}

func (f Formatter) ToInt(s string) int64 {
	r, _ := strconv.ParseInt(s, 10, 64)
	return r
}

func (f Formatter) ToInts(s []string) []int64 {
	res := []int64{}
	for _, si := range s {
		res = append(res, f.ToInt(si))
	}
	return res
}

func (f Formatter) ToString(s int64) string {
	return strconv.FormatInt(s, 10)

}

func (f Formatter) ToStrings(s []int64) []string {
	res := []string{}
	for _, si := range s {
		res = append(res, f.ToString(si))
	}
	return res
}
