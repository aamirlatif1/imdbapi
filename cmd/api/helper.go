package main

import (
	"encoding/json/v2"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := mux.Vars(r)

	id, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}
	return id, nil
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, values := range headers {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}
