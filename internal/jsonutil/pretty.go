package jsonutil

import (
	"encoding/json"
	"fmt"
)

func PrettyPrint(in interface{}) {
	s, err := MarshalPretty(in)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", s)
}

func MarshalPretty(in interface{}) ([]byte, error) {
	pp, err := json.MarshalIndent(in, "", "    ")
	if err != nil {
		return nil, err
	}
	return pp, nil
}
