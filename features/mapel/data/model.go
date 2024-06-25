package data

import (
	"apk-sekolah/features/mapel"
	"time"
)

type Mapel struct {
	ID       string     `json:"id"`      // Primary Key
	GuruID   string     `json:"guru_id"` // Not Null
	NamaGuru string     `json:"nama_guru"`
	Mapel    string     `json:"mapel"`               // Not Null
	CreateAd time.Time  `json:"create_ad"`           // Not Null, Default: CURRENT_TIMESTAMP
	UpdateAd *time.Time `json:"update_ad,omitempty"` // Nullable
	DeleteAd *time.Time `json:"delete_ad,omitempty"` // Nullable
}

func (u *Mapel) TableName() string {
	return "mapel"
}

func FormatterRequest(req mapel.MapelCore) Mapel {
	return Mapel{
		ID:       req.ID,
		GuruID:   req.GuruID,
		NamaGuru: req.NamaGuru,
		Mapel:    req.Mapel,
	}
}

func FormatterResponse(res Mapel) mapel.MapelCore {
	return mapel.MapelCore{
		ID:       res.ID,
		GuruID:   res.GuruID,
		NamaGuru: res.NamaGuru,
		Mapel:    res.Mapel,
	}
}
