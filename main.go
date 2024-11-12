package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strconv"
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
type Volume int
const (
	low = 100
	medium = 1000
	high = 10000
)
func StringToVolume(s string) (Volume, error) {
	switch s {
	case "low":
		return low, nil
	case "medium":
		return medium, nil
	case "high":
		return high, nil
	default:
		return -1, fmt.Errorf("unknown type: %s", s)
	}
}


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
func GetBsArgs(s string) []string {
	re := regexp.MustCompile(`--\s*bs:\s*(.+)`)
	match := re.FindStringSubmatch(s)
	if len(match) < 1 {
		return nil
	}
        args := strings.Split(match[1], ";")
	for i, part := range args {
		args[i] = strings.TrimSpace(part)
	}
	return args
}
func GetBs(bs Bs) (string, error) {
        c := bs.Args[0] // c as in collection or command
        if c == "num" {
                arg_parsed := strings.Split(bs.Args[1], "to")
                min, _ := strconv.Atoi(arg_parsed[0])
                max, _ := strconv.Atoi(arg_parsed[1])
                return strconv.Itoa(rand.Intn(max-min)+min), nil
        } else if c == "rel" {
                return strconv.Itoa(rand.Intn(int(bs.RelationshipTable.Volume))), nil
        } else if c == "strval" {
                return ApoQuote(bs.Args[1]), nil
        } else if c == "timestamp_epoch" {
                return strconv.Itoa(rand.Intn(1731414911)), nil
        } else if v, e := bsc[c]; e {
                r := rand.Intn(len(bsc[c]))
                return ApoQuote(v[r]), nil
        } else if v, e := bsc_composite[c]; e {
                return ApoQuote(v()), nil
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
	Fields []*Field
        Volume Volume
}

type Bs struct {
	Args []string
        RelationshipTable *Table
        RelationshipField *Field
}

func main() {
	f, _ := os.Open("schema.sql")
	sr := bufio.NewScanner(f)
	var tables = make(map[string]*Table)
        refmap := make(map[string]*Field)
        var ct *Table
	for sr.Scan() {
		l := sr.Text()
		s := strings.Fields(l)
		nl := len(s)
		bs := GetBsArgs(l)
		if len(bs) > 0 && nl > 2 && s[0] == "CREATE" && s[1] == "TABLE" {
                        // first argument in CREATE TABLE should be volume denominator
                        vol, _ := StringToVolume(bs[0])
			tables[s[2]] = &Table{
				Name: s[2],
                                Volume: vol,
			}
                        ct = tables[s[2]]
		} else if len(bs) > 0 && nl > 1 && bs[0] != "" {
                        // parsing regular field definition
			t, _ := StringToType(s[1])
                        f := Field{Name: s[0], Type: t, Bs: Bs{Args: bs}}
                        if (bs[0] == "rel") {
                                refmap[s[0]] = &f
                        }
                        ct.Fields = append(ct.Fields, &f)
		} else if nl > 4 && s[0] == "FOREIGN" && s[1] == "KEY" {
                        //minimal relationship would look like: FOREIGN KEY (rel_field) REFERENCES ref_table (ref_field)
                        rel_field := WithinParenthesis(s[2])
                        if _, e := refmap[rel_field]; !e {
                                continue
                        }
                        ref_table := s[4]
                        ref_field := WithinParenthesis(s[4])
                        refmap[rel_field].Bs.RelationshipTable = tables[ref_table]
                        refmap[rel_field].Bs.RelationshipField = findFieldByName(tables[ref_table].Fields,ref_field)
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
                for _ = range t.Volume {
                       for i, f := range t.Fields {
                                all_bs[i], e = GetBs(f.Bs)
                                if(e != nil) {
                                        fmt.Printf("*%s* collection is bs in table *%s* for field *%s*\n", f.Bs.Args, t.Name, f.Name)
                                        os.Exit(1)
                                }
                        }
                        fmt.Printf("INSERT INTO %s (%s) VALUES (%s);\n", t.Name, strings.Join(tf, ","), strings.Join(all_bs, ","))
                }
	}
}

