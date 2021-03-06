package repository

import (
	"FinalProject/auth"
	"FinalProject/entity"
	"FinalProject/utility"
	"database/sql"
	"strings"
)

type mitraRepositoryImpl struct {
	db *sql.DB
}

func NewMitraRepositoryImpl(db *sql.DB) *mitraRepositoryImpl {
	return &mitraRepositoryImpl{
		db: db,
	}
}

func (m *mitraRepositoryImpl) IsMitraExistsByEmail(email string) (bool, error) {
	count := 0

	query := `
	SELECT
		COUNT(id)
	FROM
		fp_user
	WHERE
		email = ?
	`
	row := m.db.QueryRow(query, email)
	if err := row.Scan(&count); err != nil {
		return false, err
	}

	if count != 1 {
		return false, utility.ErrNoDataFound
	}

	return true, nil
}

func (m *mitraRepositoryImpl) Login(username string, password string) (*entity.Mitra, error) {
	query := `
	SELECT
		email, password
	FROM
		fp_user
	WHERE
		email = ? AND kategori_user = "MITRA"
	`
	row := m.db.QueryRow(query, username)

	currentEmail := ""
	hashedPassword := ""
	if err := row.Scan(
		&currentEmail,
		&hashedPassword,
	); err != nil {
		return nil, err
	}

	passwordMatch, err := auth.ComparePassword(hashedPassword, password)
	if err != nil {
		return nil, err
	}

	if !passwordMatch || strings.Compare(currentEmail, username) != 0 {
		return nil, utility.ErrUnauthorized
	}
	
	query = `
	SELECT
		id, email, password, kategori_user
	FROM
		fp_user
	WHERE
		email = ? AND kategori_user = "MITRA"
	`

	mitra := entity.Mitra{}

	row = m.db.QueryRow(query, username)
	if err := row.Scan(
		&mitra.Id,
		&mitra.Email,
		&mitra.Password,
		&mitra.KategoriUser,
	); err != nil {
		return nil, err
	}

	if mitra == (entity.Mitra{}) {
		return nil, utility.ErrNoDataFound
	}

	return &mitra, nil
}

func (m *mitraRepositoryImpl) RegisterMitra(mitra *entity.MitraDetail, user *entity.Mitra) (*entity.MitraDetail, error) {
	passsword, err := auth.CreatePassword(user.Password)
	if err != nil {
		return nil, err
	}

	queryUser := `
	INSERT INTO
		fp_user (email, password, kategori_user)
	VALUES
		(?, ?, "MITRA")
	`
	result, err := m.db.Exec(queryUser, user.Email, passsword)
	if err != nil {
		return nil, err
	}
	userID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	mitra.IdUser = int(userID)

	query := `
	INSERT INTO fp_mitra (id_user, nama, about, nomor_pic, nama_pic, situs, alamat)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err = m.db.Exec(query, mitra.IdUser, mitra.Nama, mitra.About, mitra.NomorPic, mitra.NamaPic, mitra.Situs, mitra.Alamat)
	if err != nil {
		return nil, err
	}

	return mitra, nil
}
