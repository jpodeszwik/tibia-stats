package domain

import "time"

type Death struct {
	CharacterName string
	Guild         string
	Time          time.Time
	Reason        string
}
