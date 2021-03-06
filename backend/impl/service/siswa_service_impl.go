package service

import (
	"FinalProject/api/repository"
	"FinalProject/auth"
	"FinalProject/entity"
	"FinalProject/payload"
	"FinalProject/utility"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var secretKey = []byte("Final Project Beasiswa")

type siswaServiceImpl struct {
	siswaRepository repository.SiswaRepository
}

func NewSiswaServiceImpl(siswaRepository repository.SiswaRepository) *siswaServiceImpl {
	return &siswaServiceImpl{
		siswaRepository: siswaRepository,
	}
}

func (s *siswaServiceImpl) Login(request payload.LoginRequest) (*payload.LoginResponse, error) {
	isThere, err := s.siswaRepository.IsSiswaExistsByEmail(request.Email)
	if err != nil {
		return nil, err
	}

	if !isThere {
		return nil, utility.ErrNoDataFound
	}

	siswa, err := s.siswaRepository.Login(request.Email, request.Password)
	if err != nil {
		return nil, err
	}

	tokenString, err := auth.CreateJWTToken(siswa.Email, siswa.KategoriUser)
	if err != nil {
		return nil, err
	}

	return &payload.LoginResponse{
		Role:  siswa.KategoriUser,
		Token: tokenString,
	}, nil
}

func (s *siswaServiceImpl) GetListSiswa(request payload.ListSiswaRequest) (*payload.ListSiswaResponse, error) {
	totalSiswa, err := s.siswaRepository.GetTotalSiswa(request.Nama)
	if err != nil {
		return nil, err
	}

	nextPage, prevPage, totalPages := utility.GetPaginateURL("api/siswa", &request.Page, &request.Limit, totalSiswa)

	listSiswa, err := s.siswaRepository.GetListSiswa(request.Page, request.Limit, request.Nama)
	if err != nil {
		return nil, err
	}

	lenListSiswa := len(listSiswa)
	if lenListSiswa == 0 {
		return nil, utility.ErrNoDataFound
	}

	results := make([]payload.Siswa, 0)
	for i := 0; i < lenListSiswa; i++ {
		siswa := listSiswa[i]
		results = append(results, payload.Siswa{
			Id:                siswa.Id,
			IdUser:            siswa.IdUser,
			Nama:              siswa.Nama,
			NamaInstansi:      siswa.NamaInstansi,
			TingkatPendidikan: siswa.TingkatPendidikan,
			Alamat:            siswa.Alamat,
			NomorTelepon:      siswa.NomorTelepon,
			Email:             siswa.Email,
			TanggalLahir:      siswa.TanggalLahir,
			NomorRekening:     siswa.NomorRekening,
			NamaBank:          siswa.NamaBank,
		})
	}

	return &payload.ListSiswaResponse{
		Data: results,
		PaginateInfo: payload.PaginateInfo{
			NextPage:   nextPage,
			PrevPage:   prevPage,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *siswaServiceImpl) RegisterSiswa(request payload.Siswa) (*payload.LoginResponse, error) {
	siswa, err := s.siswaRepository.RegisterSiswa(&entity.SiswaDetail{
		IdUser:            request.IdUser,
		Nama:              request.Nama,
		NamaInstansi:      request.NamaInstansi,
		TingkatPendidikan: request.TingkatPendidikan,
		Alamat:            request.Alamat,
		NomorTelepon:      request.NomorTelepon,
		Email:             request.Email,
		TanggalLahir:      request.TanggalLahir,
		NomorRekening:     request.NomorRekening,
		NamaBank:          request.NamaBank,
	}, &entity.Siswa{
		Email:    request.Email,
		Password: request.Password,
	})

	if err != nil {
		return nil, err
	}

	claims := payload.Claims{
		Email: siswa.Email,
		Role:  "SISWA",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(3 * 60 * time.Minute).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenSting, err := token.SignedString(secretKey)
	if err != nil {
		return nil, err
	}

	return &payload.LoginResponse{
		Role:  "SISWA",
		Token: tokenSting,
	}, nil
}
