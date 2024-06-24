package handler

import (
	"apk-sekolah/features/mapel"

	golangmodule "github.com/rahmathidayat72/golang-module"
)

type RequestMapel struct {
	ID     string `json:"id"`
	GuruID string `json:"guru_id"` // Not Null
	Mapel  string `json:"mapel"`   // Not Null
}

func FormatterRequest(Input RequestMapel) mapel.MapelCore {
	return mapel.MapelCore{
		ID:     golangmodule.GenerateUUIDV4(),
		GuruID: Input.GuruID,
		Mapel:  Input.Mapel,
	}
}
