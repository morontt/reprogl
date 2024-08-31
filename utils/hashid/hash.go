package hashid

import (
	"errors"
	"fmt"

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

func Encode(id int, options Option) (hashString string) {
	hashString, err := encode(id, options)
	if err != nil {
		panic(err)
	}

	return
}

func Decode(str string, isAvatar bool) (HashData, error) {
	var data = HashData{Hash: str}

	ids, err := decode(str)
	if err != nil {
		return data, err
	}
	if len(ids) != 2 {
		return data, errors.New(fmt.Sprintf("hashid: incorrect result length: %d", len(ids)))
	}

	data.ID = ids[0]
	data.options = Option(ids[1])

	if isAvatar && !data.validOptions() {
		return data, WrongOptions
	}

	return data, nil
}

func decode(str string) ([]int, error) {
	return hash.DecodeWithError(str)
}

func encode(id int, options Option) (string, error) {
	return hash.EncodeInt64([]int64{int64(id), int64(options)})
}
