package helpers

import (
	"errors"
	"io"
	"net/http"
)

func HttpPostHelper(url string, bodyOfReq io.Reader) ([]byte, error) {
	resp, err := http.Post(url, "application/json", bodyOfReq)

	if err != nil {
		return nil, errors.New("request could not be send")
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("answer could not be read")
	}

	return body, nil
}
