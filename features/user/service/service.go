package service

import (
	"apk-sekolah/features/user"
	"errors"
	"fmt"

	golangmodule "github.com/rahmathidayat72/golang-module"
	"gorm.io/gorm"
)

type userService struct {
	userData user.DataInterface
}

func NewServiceUser(repo user.DataInterface) user.ServiceInterface {
	return &userService{
		userData: repo,
	}
}

// Insert implements user.ServiceInterface.
func (s *userService) Insert(insert user.UserCore) error {
	// Validasi data yang wajib diisi
	if err := golangmodule.ValidasiRequired(insert.Nama, insert.Email, insert.Password); err != nil {
		return fmt.Errorf("validation error: %v", err)
	}

	// Validasi email input
	if err := golangmodule.ValidasiEmail(insert.Email); err != nil {
		return fmt.Errorf("validation error: %v", err)
	}

	// Validasi password input
	if err := golangmodule.InputCombinationPassword(insert.Password); err != nil {
		return fmt.Errorf("validation error: Password must meet the following criteria - minimum length of 8 characters, first character must be uppercase, and it must consist of a combination of letters, numbers, and special characters,validation error: %v", err)
	}

	// Set default role if not provided
	if insert.Role == "" {
		insert.Role = "user"
	} else {
		// Validasi role input
		if insert.Role != "admin" && insert.Role != "user" {
			return errors.New("validation error: Role can only be set to 'user' or 'admin'")
		}
	}

	// Memasukkan data ke dalam database
	err := s.userData.Insert(insert)
	if err != nil {
		return fmt.Errorf("validation error: %v", err)
	}

	return nil
}

// GetAll implements user.ServiceInterface.
func (s *userService) GetAll() ([]user.UserCore, error) {
	// panic("unimplemented")
	result, err := s.userData.SelectAll()
	if err != nil {
		// Manajemen kesalahan di sini, Anda bisa memilih untuk melakukan logging atau
		// mengembalikan kesalahan kepada pemanggil fungsi.
		// Contoh: log.Println("Error in GetAll:", err)
		return nil, err
	}
	return result, nil
}

// SelectById implements user.ServiceInterface.
func (s *userService) SelectById(id string) (user.UserCore, error) {
	if id == "" {
		return user.UserCore{}, errors.New("id cannot be empty")
	}

	result, err := s.userData.SelectById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user.UserCore{}, fmt.Errorf("user with id %s not found", id)
		}
		return user.UserCore{}, fmt.Errorf("error in SelectById (userData.SelectById): %s", err)
	}

	if result.ID == "" {
		// Jika ID adalah string kosong, ini menunjukkan bahwa data tidak ditemukan
		return user.UserCore{}, fmt.Errorf("user with id %s not found", id)
	}

	return result, nil
}

// GetByRole implements user.ServiceInterface.
func (s *userService) GetByRole(userList *[]user.UserCore, role string) ([]user.UserCore, error) {

	result, err := s.userData.GetByRole(userList, role)
	if err != nil {
		// Manajemen kesalahan di sini, Anda bisa memilih untuk melakukan logging atau
		// mengembalikan kesalahan kepada pemanggil fungsi.
		// Contoh: log.Println("Error in GetByRole:", err)
		return nil, err
	}
	return result, nil
}

// SearchUsers implements user.ServiceInterface.
func (s *userService) SearchUsers(userList *[]user.UserCore, searchParam string) ([]user.UserCore, error) {
	result, err := s.userData.SearchUsers(userList, searchParam)
	if err != nil {
		// Manajemen kesalahan di sini, Anda bisa memilih untuk melakukan logging atau
		// mengembalikan kesalahan kepada pemanggil fungsi.
		// Contoh: log.Println("Error in SearchUsers:", err)
		return nil, err
	}
	return result, nil
}

// SelectByName implements user.ServiceInterface.
func (s *userService) DetailByName(nama string) (user.UserCore, error) {
	if nama == "" {
		return user.UserCore{}, errors.New("name cannot be empty")
	}
	result, err := s.userData.DetailByName(nama)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user.UserCore{}, fmt.Errorf("user with name %s not found", nama)
		}
		return user.UserCore{}, fmt.Errorf("error in DetailByName (userData.DetailByName): %s", err)
	}

	if result.Nama == "" {
		// Jika ID adalah string kosong, ini menunjukkan bahwa data tidak ditemukan
		return user.UserCore{}, fmt.Errorf("user with name %s not found", nama)
	}
	return result, nil
}

// Update implements user.ServiceInterface.
func (s *userService) Update(insert user.UserCore, id string) error {
	if id == "" {
		return errors.New("validation error. invalid ID")
	}

	err := s.userData.Update(insert, id)
	if err != nil {
		return err
	}
	return nil
}

// Delete implements user.ServiceInterface.
func (s *userService) Delete(id string) error {
	// panic("unimplemented")
	if id == "" {
		return errors.New("validation error. invalid id")
	}
	err := s.userData.Delete(id)
	return err
}
