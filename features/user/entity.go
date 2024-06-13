package user

import (
	"time"
)

type (
	UserCore struct {
		ID       string    `json:"id"`
		Nama     string    `json:"nama"`
		Email    string    `json:"email"`
		Password string    `json:"password"`
		Telepon  string    `json:"telepon"`
		Alamat   string    `json:"alamat"`
		Role     string    `json:"role"`
		CreateAt time.Time `json:"create_at"`
	}
)

type DataInterface interface {
	Insert(insert UserCore) error
	SelectAll() ([]UserCore, error)
	SelectById(id string) (UserCore, error)
	DetailByName(nama string) (UserCore, error)
	Update(insert UserCore, id string) error
	Delete(id string) error
	GetByRole(userList *[]UserCore, role string) ([]UserCore, error)          // Tambahkan metode baru
	SearchUsers(userList *[]UserCore, searchParam string) ([]UserCore, error)  // Tambahkan metode baru
}

type ServiceInterface interface {
	Insert(insert UserCore) error
	GetAll() ([]UserCore, error)
	SelectById(id string) (UserCore, error)
	DetailByName(nama string) (UserCore, error)
	Update(insert UserCore, id string) error
	Delete(id string) error
	GetByRole(userList *[]UserCore, role string) ([]UserCore, error)     // Tambahkan metode baru
	SearchUsers(userList *[]UserCore, searchParam string) ([]UserCore, error) // Tambahkan metode baru
}
