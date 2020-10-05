package utils

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
)

func TestMakeBody(pathJSON string) io.ReadCloser {
	if len(pathJSON) <= 0 {
		return nil
	}

	jsonFile, err := os.Open(pathJSON)
	if err != nil {
		panic(err)
	}

	b, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}

	return ioutil.NopCloser(bytes.NewReader(b))
}
