package testutils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func ReadJSONTestFixture(fixturePath string, v interface{}) {
	data, err := ioutil.ReadFile(fixturePath)
	if err != nil {
		panic(fmt.Sprintf("error reading the test data file : %v", err))
	}

	if err = json.Unmarshal(data, v); err != nil {
		panic(fmt.Sprintf("error parsing the test data json: %v", err))
	}
}

func ReadJSONTestFixtureAsString(fixturePath string) string {
	data, err := ioutil.ReadFile(fixturePath)
	if err != nil {
		panic(fmt.Sprintf("error reading the test data file : %v", err))
	}

	return string(data)
}
