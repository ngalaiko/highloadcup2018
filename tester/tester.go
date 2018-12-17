package tester

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/ngalayko/highloadcup/tester/answers"
)

// Tester runs tests for the app.
type Tester struct {
	answers []*answers.Answer
}

// New is a tester constructor.
func New(
	u string,
	dataPath string,
) (*Tester, error) {
	endpoint, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	aa, err := answers.Parse(endpoint, dataPath)
	if err != nil {
		return nil, err
	}

	return &Tester{
		answers: aa,
	}, nil
}

var digits = regexp.MustCompile(`\d+`)

// Run runs tests.
func (t *Tester) Run(
	healthcheckEndpoint string,
	logPath string,
) error {
	<-waitHealthcheck(healthcheckEndpoint)

	log.Printf("starting tests...")

	reports := make(map[string]*answers.Report, len(t.answers))

	successful := 0
	for _, a := range t.answers {
		ok, report, err := a.Test()
		if err != nil {
			return err
		}

		if ok {
			successful++
			continue
		}

		reports[report.URL.Path] = report
	}

	for path, report := range reports {
		path = strings.Replace(path, "/", "_", -1)
		path = digits.ReplaceAllString(path, "")
		if err := writeToFile(
			fmt.Sprintf("%s/%s.log", logPath, path),
			report,
		); err != nil {
			return err
		}
	}

	log.Printf("logs are saved to: %s", logPath)
	log.Printf("successful: %d", successful)
	log.Printf("failed: %d", len(t.answers)-successful)
	log.Printf("success rate: %.2f%%", float64(successful)/float64(len(t.answers))*100)
	return nil
}

func writeToFile(filePath string, data interface{}) error {
	reportsData, err := json.MarshalIndent(data, "", "	")
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprint(file, string(reportsData))
	return err
}

func waitHealthcheck(url string) <-chan bool {
	log.Printf("waiting for status code 200 from: %s", url)

	respChan := make(chan bool)

	go func() {
		for range time.Tick(time.Second) {
			resp, _ := http.Get(url)
			if resp != nil && resp.StatusCode == 200 {
				respChan <- true
			}
		}
	}()

	return respChan
}
