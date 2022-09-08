package models

type Partner struct {
	Id                 int16
	Name               string
	Lat                float32
	Lng                float32
	Radius             float32
	Rating             float32 `minimum:"0" maximum:"10" default:"0"`
	FlooringExperience string  `enums:"carpet,tiles,wood"`
}
