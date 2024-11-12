package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strings"
)
//go:generate go run gen/bs_gen.go
type Type int

const (
	TEXT = iota
	NUMERIC
	INTEGER
	REAL
	BLOB
)

func StringToType(s string) (Type, error) {
	switch strings.ToUpper(s) {
	case "TEXT":
		return TEXT, nil
	case "NUMERIC":
		return NUMERIC, nil
	case "INTEGER":
		return INTEGER, nil
	case "REAL":
		return REAL, nil
	case "BLOB":
		return BLOB, nil
	default:
		return -1, fmt.Errorf("unknown type: %s", s)
	}
}
func GetBsArgs(s string) string {
	re := regexp.MustCompile(`--\s*bs:\s*(.+)`)
	match := re.FindStringSubmatch(s)
	if len(match) < 1 {
		return ""
	}
	return match[1]
}
func GetBs(bs Bs) (string, error) {
        if bs.CollectionName == "number" {
                return string(rand.Intn(1000)), nil
        } else if v, e := bsc[bs.CollectionName]; e {
                r := rand.Intn(len(bsc[bs.CollectionName]))
                return "'"+v[r]+"'", nil
        } else if v, e := bsc_composite[bs.CollectionName]; e {
                return "'"+v()+"'", nil
        }

        return "", errors.New("Non existent")
}

type Field struct {
	Name string
	Type Type
	Bs   Bs
}

type Table struct {
	Name   string
	Fields []Field
}

type Bs struct {
	CollectionName string
}

func main() {

	f, _ := os.Open("schema.sql")

	sr := bufio.NewScanner(f)
	var tables []Table
	var ti int = -1
	for sr.Scan() {
		l := sr.Text()
		s := strings.Fields(l)
		nl := len(s)
		bs := GetBsArgs(l)
		if nl > 2 && s[0] == "CREATE" && s[1] == "TABLE" {
			tables = append(tables, Table{
				Name: s[2],
			})
			ti++
		} else if nl > 1 && bs != "" {
			t, _ := StringToType(s[1])
			tables[ti].Fields = append(tables[ti].Fields, Field{Name: s[0], Type: t, Bs: Bs{CollectionName: bs}})
		}
	}

	for _, t := range tables {
                if (len(t.Fields) < 1) {
                        continue
                }
		tf := make([]string, len(t.Fields))
		for i, f := range t.Fields {
			tf[i] = f.Name
		}

		all_bs := make([]string, len(t.Fields))
                var e error
		for i, f := range t.Fields {
			all_bs[i], e = GetBs(f.Bs)
                        if(e != nil) {
                                fmt.Printf("*%s* collection is bs in table *%s* for field *%s*\n", f.Bs.CollectionName, t.Name, f.Name)
                                os.Exit(1)
                        }
		}
		fmt.Printf("INSERT INTO %s (%s) VALUES (%s);\n", t.Name, strings.Join(tf, ","), strings.Join(all_bs, ","))
	}
}
