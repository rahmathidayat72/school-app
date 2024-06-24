package data

import (
	"apk-sekolah/features/mapel"
	"errors"
	"fmt"
	"log"

	"gorm.io/gorm"
)

type MapelQuery struct {
	db *gorm.DB
}

func NewDataMApel(db *gorm.DB) mapel.DataMapelInterface {
	return &MapelQuery{
		db: db,
	}
}

// Insert implements mapel.DataMapelInterface.
func (r *MapelQuery) Insert(insert mapel.MapelCore) error {
	// panic("unimplemented")
	mapelInput := FormatterRequest(insert)

	result := r.db.Exec(`INSERT INTO school."mapel" (id,guru_id, mapel )
	VALUES (?,?, ?)`, mapelInput.ID, mapelInput.GuruID, mapelInput.Mapel)

	if result.Error != nil {
		log.Printf("Error inserting user: %v", result.Error)
		return result.Error
	}

	rowsAffected := result.RowsAffected
	if rowsAffected == 0 {
		err := errors.New("failed to insert, row affected is 0")
		log.Printf("Error inserting user: %v", err)
		return err
	}

	return nil
}

// SelectAll implements mapel.DataMapelInterface.
func (r *MapelQuery) SelectAll() ([]mapel.MapelCore, error) {
	var dataMapel []Mapel

	if tx := r.db.Raw(`SELECT id,guru_id,mapel FROM school."mapel" WHERE "delete_ad" IS NULL`).Scan(&dataMapel); tx.Error != nil {
		log.Printf("Error executing SELECT query: %v", tx.Error)
		return nil, tx.Error
	}

	var coreMapel []mapel.MapelCore

	for _, v := range dataMapel {
		var mapel = FormatterResponse(v)
		coreMapel = append(coreMapel, mapel)
	}
	log.Printf("Successfully fetched %d users from database", len(coreMapel))
	fmt.Println("1", coreMapel)
	return coreMapel, nil
}

// Update implements mapel.DataMapelInterface.
func (r *MapelQuery) Update(insert mapel.MapelCore, id string) error {
	// panic("unimplemented")
	var count int64
	tx := r.db.Raw(`SELECT COUNT(*) FROM school."mapel" u WHERE "id" =?`, id).Count(&count)
	if tx.Error != nil {
		if count == 0 {
			log.Printf("Mapel with ID %s not found", id)
			return errors.New("Mapel id not found")
		}
	}
	query := `
	UPDATE school."mapel"
	SET
	guru_id = COALESCE(NULLIF(?, ''), guru_id),
	mapel = COALESCE(NULLIF(?, ''), mapel)
	WHERE id = ?
    `

	result := r.db.Exec(query, insert.GuruID,
		insert.Mapel,
		id)

	if result.Error != nil {
		log.Printf("Error updating user: %v", result.Error)
		return result.Error
	}

	return nil
}

// Delete implements mapel.DataMapelInterface.
func (r *MapelQuery) Delete(id string) error {
	// panic("unimplemented")
	var count int64
	tx := r.db.Raw(`SELECT COUNT(*) FROM school."mapel" u WHERE "id" =?`, id).Count(&count)
	if tx.Error != nil {
		if count == 0 {
			log.Printf("Mapel with ID %s not found", id)
			return errors.New("Mapel id not found")
		}
	}

	query := `UPDATE school."mapel"  SET "delete_ad" = NOW() WHERE "id" = ?`
	result := r.db.Exec(query, id)

	if result.Error != nil {
		log.Printf("Failed to soft delete mapel: %v", result.Error)
		return errors.New("failed to soft delete mapel")
	}
	return nil
}

// GetByGuruID implements mapel.DataMapelInterface.
func (r *MapelQuery) GetByGuruID(mapelList *[]mapel.MapelCore, guruID string) ([]mapel.MapelCore, error) {
	// panic("unimplemented")
	var dataMapel []Mapel
	querry := `SELECT * FROM school.mapel WHERE guru_id = ? AND delete_ad IS NULL`
	if tx := r.db.Raw(querry, guruID).Scan(&dataMapel); tx.Error != nil {
		log.Printf("Error executing SELECT query: %v", tx.Error)
		return nil, tx.Error
	}
	var mapelCore []mapel.MapelCore
	for _, v := range dataMapel {
		var mapel = FormatterResponse(v)
		mapelCore = append(mapelCore, mapel)
	}
	// Logging informasi sukses setelah loop
	log.Printf("Successfully fetched %d users from database", len(mapelCore))

	return mapelCore, nil
}

// SearchMapel implements mapel.DataMapelInterface.
func (r *MapelQuery) SearchMapel(mapelList *[]mapel.MapelCore, searchParam string) ([]mapel.MapelCore, error) {
	// panic("unimplemented")
	var dataMapel []Mapel
	querry := `
	SELECT * FROM school."mapel" m
	WHERE m."mapel" LIKE ?  OR m."id" LIKE ? AND "delete_ad" IS NULL
	`
	if tx := r.db.Raw(querry, "%"+searchParam+"%", "%"+searchParam+"%").Scan(&dataMapel); tx.Error != nil {
		log.Printf("Error executing SELECT query: %v", tx.Error)
		return nil, tx.Error
	}

	var coreMapel []mapel.MapelCore

	for _, v := range dataMapel {
		var mapel = FormatterResponse(v)
		coreMapel = append(coreMapel, mapel)
	}
	// Logging informasi sukses setelah loop
	log.Printf("Successfully fetched %d mapel from database", len(coreMapel))

	return coreMapel, nil
}
