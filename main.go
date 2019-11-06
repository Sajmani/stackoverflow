package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"os"
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
		schema                      map[string]int
		all, gcp, aws, azure, cloud langCounts
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
		all.N++
		var Go, Java bool
		langs := strings.Split(rec[schema["LanguageWorkedWith"]], ";")
		for _, lang := range langs {
			switch lang {
			case "Go":
				all.Go++
				Go = true
			case "Java":
				all.Java++
				Java = true
			}
		}
		plats := strings.Split(rec[schema["PlatformWorkedWith"]], ";")
		var Cloud bool
		for _, plat := range plats {
			switch plat {
			case "Google Cloud Platform":
				gcp.N++
				Cloud = true
				if Go {
					gcp.Go++
				}
				if Java {
					gcp.Java++
				}
			case "AWS":
				aws.N++
				Cloud = true
				if Go {
					aws.Go++
				}
				if Java {
					aws.Java++
				}
			case "Microsoft Azure":
				azure.N++
				Cloud = true
				if Go {
					azure.Go++
				}
				if Java {
					azure.Java++
				}
			}
		}
		if Cloud {
			cloud.N++
			if Go {
				cloud.Go++
			}
			if Java {
				cloud.Java++
			}
		}
	}
	fmt.Printf("all %s\n", all.String(all.N))
	fmt.Printf("gcp %s\n", gcp.String(all.N))
	fmt.Printf("aws %s\n", aws.String(all.N))
	fmt.Printf("azure %s\n", azure.String(all.N))
	fmt.Printf("cloud %s\n", cloud.String(all.N))

	fmt.Printf("P(GCP|Go) = P(GCP,Go)/P(Go) = %.1f/%.1f%% = %.1f%% vs P(GCP|Java) = %.1f%% vs P(GCP) = %.1f%%\n",
		pct(gcp.Go, all.N), pct(all.Go, all.N), pct(gcp.Go, all.Go), pct(gcp.Java, all.Java), pct(gcp.N, all.N))
	fmt.Printf("P(AWS|Go) = P(AWS,Go)/P(Go) = %.1f/%.1f%% = %.1f%% vs P(AWS|Java) = %.1f%% vs P(AWS) = %.1f%%\n",
		pct(aws.Go, all.N), pct(all.Go, all.N), pct(aws.Go, all.Go), pct(aws.Java, all.Java), pct(aws.N, all.N))
	fmt.Printf("P(Azure|Go) = P(Azure,Go)/P(Go) = %.1f/%.1f%% = %.1f%% vs P(Azure|Java) = %.1f%% vs P(Azure) = %.1f%%\n",
		pct(azure.Go, all.N), pct(all.Go, all.N), pct(azure.Go, all.Go), pct(azure.Java, all.Java), pct(azure.N, all.N))
	fmt.Printf("P(Cloud|Go) = P(Cloud,Go)/P(Go) = %.1f/%.1f%% = %.1f%% vs P(Cloud|Java) = %.1f%% vs P(Cloud) = %.1f%%\n",
		pct(cloud.Go, all.N), pct(all.Go, all.N), pct(cloud.Go, all.Go), pct(cloud.Java, all.Java), pct(cloud.N, all.N))
}

type langCounts struct {
	N    int
	Go   int
	Java int
}

func (c langCounts) String(n int) string {
	return fmt.Sprintf("Total %d (%.1f%%)  Go %d (%.1f%% of Total, %.1f%% of all)  Java %d (%.1f%% of Total, %.1f%% of all)",
		c.N, pct(c.N, n),
		c.Go, pct(c.Go, c.N), pct(c.Go, n),
		c.Java, pct(c.Java, c.N), pct(c.Java, n))
}

func pct(a, b int) float64 {
	return math.Round(1000*float64(a)/float64(b)) / 10
}
