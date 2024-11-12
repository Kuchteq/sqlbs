package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

var bs_composite = map[string]func(bs Bs) (string, error){
	"num": func(bs Bs) (string, error) {
		arg_parsed := strings.Split(bs.Args[1], "to")
		min, _ := strconv.Atoi(arg_parsed[0])
		max, _ := strconv.Atoi(arg_parsed[1])
		return strconv.Itoa(rand.Intn(max-min) + min), nil
	},
	"rel": func(bs Bs) (string, error) {
		return strconv.Itoa(rand.Intn(int(bs.RelationshipTable.Volume))), nil
	},
	"val": func(bs Bs) (string, error) {
		return bs.Args[1], nil
	},
	"timestamp_epoch": func(bs Bs) (string, error) {
		return strconv.Itoa(rand.Intn(1731414911)), nil
	},
	"fullname": func(bs Bs) (string, error) {
		r_0 := rand.Intn(len(bsa["first_name"]))
		r_1 := rand.Intn(len(bsa["last_name"]))
		r := fmt.Sprintf("%s %s", bsa["first_name"][r_0], bsa["last_name"][r_1])
		return r, nil
	},
}
