package payload

type Siswa struct {
	Id                int    `json:"id"`
	IdUser            int    `json:"idUser"`
	Nama              string `json:"nama"`
	NamaInstansi      string `json:"namaInstansi"`
	TingkatPendidikan string `json:"tingkatPendidikan"`
	Alamat            string `json:"alamat"`
	NomorTelepon      string `json:"nomorTelepon"`
	Email             string `json:"email"`
	TanggalLahir      string `json:"tanggalLahir"`
	NomorRekening     string `json:"nomorRekening"`
	NamaBank          string `json:"namaBank"`
	Password          string `json:"password"`
}

type Mitra struct {
	Id       int    `json:"id"`
	Nama     string `json:"nama"`
	Profile  string `json:"profile"`
	Pic      string `json:"pic"`
	Situs    string `json:"situs"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ListSiswaRequest struct {
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
	Nama  string `json:"nama"`
}

type ListSiswaResponse struct {
	Data         []Siswa      `json:"data"`
	PaginateInfo PaginateInfo `json:"paginateInfo"`
}
