package handler

type (
	ResponseMapel struct {
		ID       string `json:"id"`      // Primary Key
		GuruID   string `json:"guru_id"` // Not Null
		NamaGuru string `json:"nama_guru"`
		Mapel    string `json:"mapel"` // Not Null
	}
)

func FormatterResponse(res ResponseMapel) ResponseMapel {
	return ResponseMapel{
		ID:       res.ID,
		GuruID:   res.GuruID,
		NamaGuru: res.NamaGuru,
		Mapel:    res.Mapel,
	}
}
