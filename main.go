package main

import (
	"encoding/csv"
	"io"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	ent := false
	if len(os.Args) > 1 && os.Args[1] == "ent" {
		ent = true // restrict to enterprise users
	}
	f, err := os.Open("/Users/sameer/Downloads/developer_survey_2019/survey_results_public.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	r := csv.NewReader(f)
	var (
		schema map[string]int
		counts = make(map[string]map[string]int)
	)
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if schema == nil {
			schema = make(map[string]int)
			for i, field := range rec {
				schema[field] = i
			}
			continue
		}
		if ent && rec[schema["OrgSize"]] != "10,000 or more employees" {
			continue
		}

		techs := make(map[string]bool)
		addTechs := func(key string) {
			for _, tech := range strings.Split(rec[schema[key]], ";") {
				techs[tech] = true
			}
		}
		addTechs("LanguageWorkedWith")
		addTechs("PlatformWorkedWith")
		addTechs("DevEnviron")
		for t1 := range techs {
			for t2 := range techs {
				if t2 < t1 {
					continue
				}
				if counts[t1] == nil {
					counts[t1] = make(map[string]int)
				}
				counts[t1][t2]++
			}
		}
	}

	var techs []string
	for t1 := range counts {
		techs = append(techs, t1)
	}
	sort.Strings(techs)

	w := csv.NewWriter(os.Stdout)
	w.Write(append([]string{"Tech"}, techs...))
	for _, t1 := range techs {
		rec := make([]string, len(techs)+1)
		rec[0] = t1
		for i, t2 := range techs {
			k1, k2 := t1, t2
			if k2 < k1 {
				k1, k2 = k2, k1
			}
			rec[i+1] = strconv.Itoa(counts[k1][k2])
		}
		w.Write(rec)
	}
}

func pct(a, b int) float64 {
	return math.Round(1000*float64(a)/float64(b)) / 10
}
