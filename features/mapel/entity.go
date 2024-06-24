package mapel

import "time"

type (
	MapelCore struct {
		ID       string     `json:"id"`                  // Primary Key
		GuruID   string     `json:"guru_id"`             // Not Null
		Mapel    string     `json:"mapel"`               // Not Null
		CreateAd time.Time  `json:"create_ad"`           // Not Null, Default: CURRENT_TIMESTAMP
		UpdateAd *time.Time `json:"update_ad,omitempty"` // Nullable
		DeleteAd *time.Time `json:"delete_ad,omitempty"` // Nullable
	}
	DataMapelInterface interface {
		Insert(insert MapelCore) error
		SelectAll() ([]MapelCore, error)
		Update(insert MapelCore, id string) error
		Delete(id string) error
		GetByGuruID(mapelList *[]MapelCore, guruID string) ([]MapelCore, error)      // Tambahkan metode baru
		SearchMapel(mapelList *[]MapelCore, searchParam string) ([]MapelCore, error) // Tambahkan metode baru
	}

	ServiceMapelInterface interface {
		Insert(insert MapelCore) error
		GetAll() ([]MapelCore, error)
		Update(insert MapelCore, id string) error
		Delete(id string) error
		GetByGuruID(mapelList *[]MapelCore, guruID string) ([]MapelCore, error)     // Tambahkan metode baru
		SearchMapel(userList *[]MapelCore, searchParam string) ([]MapelCore, error) // Tambahkan metode baru
	}
)
