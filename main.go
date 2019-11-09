package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func main() {
	path := "/Users/sameer/Downloads/developer_survey_2019/survey_results_public.csv"
	f, err := os.Open(path)
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
		keys := make(map[string]keyset)
		keys["2019"] = keyset{lang: "LanguageWorkedWith", plat: "PlatformWorkedWith", ed: "DevEnviron", orgSize: "OrgSize"}
		keys["2018"] = keyset{lang: "LanguageWorkedWith", plat: "PlatformWorkedWith", ed: "IDE", orgSize: "CompanySize"}
		keys["2017"] = keyset{lang: "HaveWorkedLanguage", plat: "HaveWorkedPlatform", ed: "IDE", orgSize: "CompanySize"}
		year := yearFromPath(path)
		techSet := make(map[string]bool)
		addTechs := func(key string) {
			for _, tech := range strings.Split(rec[schema[key]], ";") {
				tech = strings.TrimSpace(tech)
				techSet[tech] = true
			}
		}
		addTechs(keys[year].lang)
		addTechs(keys[year].plat)
		addTechs(keys[year].ed)
		if techSet["AWS"] || techSet["Microsoft Azure"] || techSet["Google Cloud Platform"] || techSet["Amazon Web Services (AWS)"] || techSet["Google Cloud Platform/App Engine"] || techSet["Azure"] {
			techSet["ANY CLOUD"] = true
		}
		if rec[schema[keys[year].orgSize]] == "10,000 or more employees" {
			techSet["ANY ENTERPRISE"] = true
		}
		techSet["ANY"] = true
		var techs []string
		for tech := range techSet {
			techs = append(techs, tech)
		}
		sort.Strings(techs)

		inc := func(k1, k2 string) {
			if k2 < k1 {
				k1, k2 = k2, k1
			}
			if counts[k1] == nil {
				counts[k1] = make(map[string]int)
			}
			counts[k1][k2]++
		}
		for i, t1 := range techs {
			for _, t2 := range techs[i:] {
				inc(t1, t2)
			}
		}
	}

	techSet := make(map[string]bool)
	for t1, m := range counts {
		techSet[t1] = true
		for t2 := range m {
			techSet[t2] = true
		}
	}
	var techs []string
	for tech := range techSet {
		techs = append(techs, tech)
	}
	sort.Strings(techs)

	w := csv.NewWriter(os.Stdout)
	defer w.Flush()
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

type keyset struct {
	lang    string
	plat    string
	ed      string
	orgSize string
}

func yearFromPath(path string) string {
	// expects the default SO path structure (.../developer_survey_YYYY/survey_results_public.csv)
	re := regexp.MustCompile(`developer_survey_(\d+)`)
	matches := re.FindStringSubmatch(path)
	if matches == nil {
		return "2019"
	}
	return matches[1]
}
