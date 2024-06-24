package data

import (
	"apk-sekolah/features/user"
	"apk-sekolah/helpers"
	"fmt"

	"errors"
	"log"

	"gorm.io/gorm"
)

type UserQuery struct {
	db *gorm.DB
}

func NewDataUser(db *gorm.DB) user.DataInterface {
	return &UserQuery{
		db: db,
	}
}

// Insert implements user.DataInterface.
func (r *UserQuery) Insert(insert user.UserCore) error {
	userInput := FormatterRequest(insert)
	userInput.Password = helpers.HashPassword(userInput.Password)

	result := r.db.Exec(`
	INSERT INTO school."user" (id,nama, email, password, telepon, alamat, role)
	VALUES (?,?, ?, ?, ?, ?, ?)
	`, userInput.ID, userInput.Nama, userInput.Email, userInput.Password, userInput.Telepon, userInput.Alamat, userInput.Role)

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

func (r *UserQuery) SelectAll() ([]user.UserCore, error) {
	var dataUser []User

	// Menggunakan r.db.Raw untuk query raw SQL dengan GORM
	tx := r.db.Raw(`SELECT id, nama, email, telepon, alamat FROM school."user" WHERE "delete_ad" IS NULL`).Scan(&dataUser)
	if tx.Error != nil {
		log.Printf("Error executing SELECT query: %v", tx.Error)
		return nil, tx.Error
	}

	var coreUser []user.UserCore

	for _, v := range dataUser {
		var user = FormatterResponse(v)
		coreUser = append(coreUser, user)
	}

	// Logging informasi sukses setelah loop
	log.Printf("Successfully fetched %d users from database", len(coreUser))
	fmt.Println("2", coreUser)
	return coreUser, nil
}

func (r *UserQuery) SelectById(id string) (user.UserCore, error) {
	var dataUser User

	tx := r.db.Raw(`
				SELECT * FROM school."user" u WHERE u."id" = ? AND u."delete_ad" IS NULL
				`, id).Scan(&dataUser)

	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return user.UserCore{}, errors.New("User id not found")
		}

		log.Printf("Error executing SELECT query: %v", tx.Error)
		return user.UserCore{}, tx.Error
	}

	userResponse := FormatterResponse(dataUser)
	return userResponse, nil
}

// SelectByName implements user.DataInterface.
func (r *UserQuery) DetailByName(nama string) (user.UserCore, error) {
	var dataUser User
	tx := r.db.Raw(`
	SELECT * FROM school."user" u WHERE u."nama" = ? AND u."delete_ad" IS NULL
	`, nama).Scan(&dataUser)

	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return user.UserCore{}, errors.New("User id not found")
		}

		log.Printf("Error executing SELECT query: %v", tx.Error)
		return user.UserCore{}, tx.Error
	}

	userResponse := FormatterResponse(dataUser)
	return userResponse, nil
}

func (r *UserQuery) Update(insert user.UserCore, id string) error {
	// Checking if the user exists
	var count int64
	r.db.Raw(`SELECT COUNT(*) FROM school."user" u WHERE "id" =?;`, id).Count(&count)
	if count == 0 {
		log.Printf("User with ID %s not found", id)
		return errors.New("User id not found")
	}

	// Constructing the SQL query for update
	query := `
	UPDATE school."user" u
	SET
	nama = COALESCE(NULLIF(?,''), nama),
	email = COALESCE(NULLIF(?,''), email),
	password = CASE WHEN ? != '' THEN ? ELSE password END,
	telepon = COALESCE(NULLIF(?,''), telepon),
	alamat = COALESCE(NULLIF(?,''), alamat)
	WHERE id = ?
    `
	hashedPassword := helpers.HashPassword(insert.Password)

	// Executing the update query
	result := r.db.Exec(query,
		insert.Nama,
		insert.Email,
		insert.Password, hashedPassword,
		insert.Telepon,
		insert.Alamat,
		id,
	)

	if result.Error != nil {
		log.Printf("Error updating user: %v", result.Error)
		return result.Error
	}

	return nil
}

func (r *UserQuery) Delete(id string) error {
	// Checking if the user exists
	var count int64
	r.db.Raw(`SELECT COUNT(*) FROM school."user" u WHERE "id" = ? AND "delete_ad" IS NULL`, id).Count(&count)
	if count == 0 {
		log.Printf("User with ID %s not found or already deleted", id)
		return errors.New("User not found or already deleted")
	}

	// Soft delete by updating the "deleted_at" column
	query := `UPDATE school."user" u SET "delete_ad" = NOW() WHERE "id" = ?`
	result := r.db.Exec(query, id)

	if result.Error != nil {
		log.Printf("Failed to soft delete user: %v", result.Error)
		return errors.New("failed to soft delete user")
	}

	return nil
}

// GetByRole implements user.DataInterface.
func (r *UserQuery) GetByRole(userList *[]user.UserCore, role string) ([]user.UserCore, error) {
	var dataUser []User
	tx := r.db.Raw(`SELECT * FROM school."user" u WHERE LOWER(u."role") = LOWER(?) AND "delete_ad" IS NULL`,
		role).Scan(&dataUser)
	if tx.Error != nil {
		log.Printf("Error executing SELECT query: %v", tx.Error)
		return nil, tx.Error
	}

	var coreUser []user.UserCore

	for _, v := range dataUser {
		var user = FormatterResponse(v)
		coreUser = append(coreUser, user)
	}

	// Logging informasi sukses setelah loop
	log.Printf("Successfully fetched %d users from database", len(coreUser))

	return coreUser, nil
}

// SearchUsers implements user.DataInterface.
func (r *UserQuery) SearchUsers(userList *[]user.UserCore, searchParam string) ([]user.UserCore, error) {
	var dataUser []User
	tx := r.db.Raw(`
	SELECT * FROM school."user" u
	WHERE u."nama" LIKE ?  OR u."email" LIKE ? OR u."alamat" LIKE ? AND "delete_ad" IS NULL
	`, "%"+searchParam+"%", "%"+searchParam+"%", "%"+searchParam+"%").Scan(&dataUser)
	if tx.Error != nil {
		log.Printf("Error executing SELECT query: %v", tx.Error)
		return nil, tx.Error
	}

	var coreUser []user.UserCore

	for _, v := range dataUser {
		var user = FormatterResponse(v)
		coreUser = append(coreUser, user)
	}

	// Logging informasi sukses setelah loop
	log.Printf("Successfully fetched %d users from database", len(coreUser))

	return coreUser, nil
}
