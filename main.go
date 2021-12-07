package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

// 2021 columns:
// ResponseId,MainBranch,Employment,Country,US_State,UK_Country,EdLevel,Age1stCode,LearnCode,YearsCode,YearsCodePro,DevType,OrgSize,
// Currency,CompTotal,CompFreq,LanguageHaveWorkedWith,LanguageWantToWorkWith,DatabaseHaveWorkedWith,DatabaseWantToWorkWith,
// PlatformHaveWorkedWith,PlatformWantToWorkWith,WebframeHaveWorkedWith,WebframeWantToWorkWith,MiscTechHaveWorkedWith,MiscTechWantToWorkWith,
// ToolsTechHaveWorkedWith,ToolsTechWantToWorkWith,NEWCollabToolsHaveWorkedWith,NEWCollabToolsWantToWorkWith,OpSys,NEWStuck,NEWSOSites,
// SOVisitFreq,SOAccount,SOPartFreq,SOComm,NEWOtherComms,Age,Gender,Trans,Sexuality,Ethnicity,Accessibility,MentalHealth,
// SurveyLength,SurveyEase,ConvertedCompYearly

func main() {
	f, err := os.Open("/Users/sameer/Downloads/stack-overflow-developer-survey-2021/survey_results_public.csv")
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
		techSet := make(map[string]bool)
		addTechs := func(key string) {
			for _, tech := range strings.Split(rec[schema[key]], ";") {
				techSet[tech] = true
			}
		}
		addTechs("LanguageHaveWorkedWith")
		addTechs("DatabaseHaveWorkedWith")
		addTechs("PlatformHaveWorkedWith")
		addTechs("WebframeHaveWorkedWith")
		addTechs("MiscTechHaveWorkedWith")
		addTechs("ToolsTechHaveWorkedWith")
		addTechs("NEWCollabToolsHaveWorkedWith")
		if techSet["AWS"] || techSet["Microsoft Azure"] || techSet["Google Cloud Platform"] {
			techSet["ANY CLOUD"] = true
		}
		if rec[schema["OrgSize"]] == "10,000 or more employees" {
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
