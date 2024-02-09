package main

import (
	"encoding/json"
	"io"
	"net/http"
)

type (
	JsonObject = map[string]interface{}
	Request    JsonObject
)

func NewRequest(in *http.Request, body JsonObject) (out Request) {
	out = map[string]interface{}{
		"path":   in.URL.Path,
		"method": in.Method,
		"body":   body,
	}

	return
}

func ServeHTTP(w http.ResponseWriter, req *http.Request) {
	enc := json.NewEncoder(w)

	defer func() {
		r := recover()

		if r == nil {
			return
		}

		w.WriteHeader(http.StatusInternalServerError)

		enc.Encode(map[string]interface{}{"error": "internal server error"})

		panic(r)
	}()

	dec := json.NewDecoder(req.Body)

	w.Header().Set("Content-Type", "application/json")

	var err error

	var m Message

	err = dec.Decode(&m.Content)

	if err == io.EOF {
		err = nil
	}

	if err != nil {
		panic(err)
	}

	m.Content = NewRequest(req, m.Content)

	_, err = m.WriteToChrome()

	if err != nil {
		panic(err)
	}

	_, err = m.ReadFromChrome()

	if err != nil && err != io.EOF {
		panic(err)
	}

	res := m.Content

	w.WriteHeader(int(res["status"].(float64)))

	if b, ok := res["body"]; ok && len(b.(JsonObject)) > 0 {
		err = enc.Encode(b)

		if err != nil {
			panic(err)
		}
	}
}
