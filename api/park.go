
package api

import (
"encoding/json"
"io/ioutil"
"net/http"
)

// Park type with Name, Author and Id
type Park struct {
	Name       string `json:"name"`
	Address      string `json:"address"`
	Id        string `json:"Id"`
}

var parks = map[string]Park{
	"0345391802": Park{Name: "Boulder Dog Park 1", Address: "101 Raod", Id: "0001"},
	"0000000000": Park{Name: "Boulder Dog Park 2", Address: "202 Road", Id: "0002"},
}

// ToJSON to be used for marshalling of Park type
func (p Park) ToJSON() []byte {
	ToJSON, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	return ToJSON
}

// FromJSON to be used for unmarshalling of Park type
func FromJSON(data []byte) Park {
	park := Park{}
	err := json.Unmarshal(data, &park)
	if err != nil {
		panic(err)
	}
	return park
}

// AllParks returns a slice of all parks
func AllParks() []Park {
	values := make([]Park, len(parks))
	idx := 0
	for _, park := range parks {
		values[idx] = park
		idx++
	}
	return values
}

// ParksHandleFunc to be used as http.HandleFunc for Park API
func ParksHandleFunc(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
	case http.MethodGet:
		parks := AllParks()
		writeJSON(w, parks)
	case http.MethodPost:
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		park := FromJSON(body)
		id, created := CreatePark(park)
		if created {
			w.Header().Add("Location", "/api/parks/"+id)
			w.WriteHeader(http.StatusCreated)
		} else {
			w.WriteHeader(http.StatusConflict)
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unsupported request method."))
	}
}

// ParkHandleFunc to be used as http.HandleFunc for Park API
func ParkHandleFunc(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/api/parks/"):]

	switch method := r.Method; method {
	case http.MethodGet:
		park, found := GetPark(id)
		if found {
			writeJSON(w, park)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	case http.MethodPut:
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		park := FromJSON(body)
		exists := UpdatePark(id, park)
		if exists {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	case http.MethodDelete:
		DeletePark(id)
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unsupported request method."))
	}
}

func writeJSON(w http.ResponseWriter, i interface{}) {
	b, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.Write(b)
}

// GetPark returns the park for a given Id
func GetPark(id string) (Park, bool) {
	park, found := parks[id]
	return park, found
}

// CreatePark creates a new Park if it does not exist
func CreatePark(park Park) (string, bool) {
	_, exists := parks[park.Id]
	if exists {
		return "", false
	}
	parks[park.Id] = park
	return park.Id, true
}

// UpdatePark updates an existing park
func UpdatePark(id string, park Park) bool {
	_, exists := parks[id]
	if exists {
		parks[id] = park
	}
	return exists
}

// DeletePark removes a park from the map by Id key
func DeletePark(id string) {
	delete(parks, id)
}

