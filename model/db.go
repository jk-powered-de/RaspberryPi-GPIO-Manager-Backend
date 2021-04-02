package model

type Database struct {
	Name string `json:"DB_NAME"`
	User string `json:"DB_USER"`
	Pass string `json:"DB_PASS"`
}
