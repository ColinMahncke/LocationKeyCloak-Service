package data

import "encoding/json"

type LocationEvent struct {
	EventType string
	Entity    json.RawMessage
}
type CreateEntity struct {
	Id        int             `json:"id"`
	Name      string          `json:"name"`
	Shortname string          `json:"shortname"`
	Geometry  json.RawMessage `json:"geometry"`
}
