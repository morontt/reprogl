package hashid

import (
	"errors"
	"github.com/speps/go-hashids/v2"
)

const (
	User Option = 1 << iota
	Commentator
	Male
	Female
)

type Option int

var hash *hashids.HashID

var WrongOptions = errors.New("hashid: wrong options")

func init() {
	var err error
	hash, err = hashids.NewWithData(&hashids.HashIDData{
		Salt:      "Thi5 is sa1t :)",
		Alphabet:  "1234567890ABCDEFGHJKLMNPQRSTUVWXYZ",
		MinLength: 6,
	})
	if err != nil {
		panic(err)
	}
}

func Decode(str string) (HashData, error) {
	var data HashData

	ids, err := decode(str)
	if err != nil {
		return data, err
	}
	if len(ids) != 2 {
		return data, errors.New("hashid: incorrect result length")
	}

	data.ID = ids[0]
	data.options = Option(ids[1])

	if !data.validOptions() {
		return data, WrongOptions
	}

	return data, nil
}

func decode(str string) ([]int, error) {
	return hash.DecodeWithError(str)
}
