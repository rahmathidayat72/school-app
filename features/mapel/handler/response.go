package handler

type (
	ResponseMapel struct {
		ID     string `json:"id"`      // Primary Key
		GuruID string `json:"guru_id"` // Not Null
		Mapel  string `json:"mapel"`   // Not Null
	}
)

func FormatterResponse(res ResponseMapel) ResponseMapel {
	return ResponseMapel{
		ID:     res.ID,
		GuruID: res.GuruID,
		Mapel:  res.Mapel,
	}
}
