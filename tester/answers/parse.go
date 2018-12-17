package answers

import (
	"bufio"
	"encoding/json"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
)

// Parse parses answers from a file.
func Parse(
	appEndpoint *url.URL,
	dataPath string,
) ([]*Answer, error) {

	file, err := os.Open(dataPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	answers := []*Answer{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		answer, err := parseAnswer(appEndpoint, scanner.Text())
		if err != nil {
			return nil, err
		}

		answers = append(answers, answer)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	log.Printf("%s: parsed %d answers", dataPath, len(answers))

	return answers, nil
}

func parseAnswer(baseURL *url.URL, line string) (*Answer, error) {
	ans := &Answer{}
	var err error
	for i, s := range strings.Split(line, "	") {
		switch i {
		case 0:
			ans.Method = s
		case 1:
			ans.URL, err = url.Parse(baseURL.String() + s)
		case 2:
			ans.StatusCode, err = strconv.Atoi(s)
		case 3:
			err = json.Unmarshal([]byte(s), &ans.Response)
		}
		if err != nil {
			break
		}
	}
	if err != nil {
		return nil, err
	}
	return ans, nil
}
