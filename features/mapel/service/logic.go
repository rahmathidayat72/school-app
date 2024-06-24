package service

import (
	"apk-sekolah/features/mapel"
	"errors"
	"fmt"

	golangmodule "github.com/rahmathidayat72/golang-module"
)

type mapelService struct {
	mapelData mapel.DataMapelInterface
}

func NewServiceMapel(repo mapel.DataMapelInterface) mapel.ServiceMapelInterface {
	return &mapelService{
		mapelData: repo,
	}
}

// Insert implements mapel.ServiceMapelInterface.
func (s *mapelService) Insert(insert mapel.MapelCore) error {
	// panic("unimplemented")
	err := golangmodule.ValidasiRequired(insert.GuruID, insert.Mapel)
	if err != nil {
		return fmt.Errorf("validation error: %v", err)

	}

	if err := s.mapelData.Insert(insert); err != nil {
		return fmt.Errorf("validation error: %v", err)

	}
	return nil
}

// GetAll implements mapel.ServiceMapelInterface.
func (s *mapelService) GetAll() ([]mapel.MapelCore, error) {
	// panic("unimplemented")
	result, err := s.mapelData.SelectAll()
	if err != nil {
		// Manajemen kesalahan di sini, Anda bisa memilih untuk melakukan logging atau
		// mengembalikan kesalahan kepada pemanggil fungsi.
		// Contoh: log.Println("Error in GetAll:", err)
		return nil, err
	}
	return result, nil
}

// Update implements mapel.ServiceMapelInterface.
func (s *mapelService) Update(insert mapel.MapelCore, id string) error {
	// panic("unimplemented")
	if id == "" {
		return errors.New("validation error. invalid ID")

	}
	err := s.mapelData.Update(insert, id)
	if err != nil {
		return err
	}
	return nil
}

// GetByGuruID implements mapel.ServiceMapelInterface.
func (s *mapelService) GetByGuruID(mapelList *[]mapel.MapelCore, guruID string) ([]mapel.MapelCore, error) {
	// panic("unimplemented")
	result, err := s.mapelData.GetByGuruID(mapelList, guruID)
	if err != nil {
		// Manajemen kesalahan di sini, Anda bisa memilih untuk melakukan logging atau
		// mengembalikan kesalahan kepada pemanggil fungsi.
		// Contoh: log.Println("Error in GetByRole:", err)
		return nil, err
	}
	return result, nil
}

// SearchMapel implements mapel.ServiceMapelInterface.
func (s *mapelService) SearchMapel(mapelList *[]mapel.MapelCore, searchParam string) ([]mapel.MapelCore, error) {
	// panic("unimplemented")
	result, err := s.mapelData.SearchMapel(mapelList, searchParam)
	if err != nil {
		// Manajemen kesalahan di sini, Anda bisa memilih untuk melakukan logging atau
		// mengembalikan kesalahan kepada pemanggil fungsi.
		// Contoh: log.Println("Error in SearchUsers:", err)
		return nil, err
	}
	return result, nil
}

func (s *mapelService) Delete(id string) error {
	// panic("unimplemented")
	if id == "" {
		return errors.New("validation error. invalid id")
	}

	err := s.mapelData.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
