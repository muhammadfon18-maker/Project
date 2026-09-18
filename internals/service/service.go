package service

import (
	"fmt"
	"project/helper"
	"project/internals/storage"

	"strings"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	users   []storage.User
	storage storage.Storager
}

type Storer interface {
	CreateUser(firstName, lastName, phone string, car []storage.CarDTO) (storage.User, error)
	GetUserById(id string) (storage.User, error)
	GetUsers() []storage.User
	UpdateUser(id, firstName, lastName, phone string, car []storage.CarDTO) (storage.User, error)
	DeleteUser(id string) error
}

func NewService(storage storage.Storager) (Storer, error) {

	users, err := storage.Load()
	if err != nil {
		return nil, err
	}

	return &Service{
		users:   users,
		storage: storage,
	}, nil

}

func (s *Service) CreateUser(
	firstName,
	lastName,
	phone string,
	car []storage.CarDTO) (storage.User, error) {
	trimF := strings.TrimSpace(firstName)

	if trimF == "" {
		return storage.User{}, helper.Invalid

	}
	trimL := strings.TrimSpace(lastName)

	if trimL == "" {
		return storage.User{}, helper.Invalid

	}
	trimP := strings.TrimSpace(phone)

	if trimP == "" {
		return storage.User{}, helper.Invalid

	}

	id_user := uuid.New().String()

	usersCar := make([]storage.Car, 0, len(car))

	for _, carDTO := range car {
		usersCar = append(usersCar, storage.Car{
			ID:      uuid.New().String(),
			Model:   carDTO.Model,
			CarType: carDTO.Car_type,
			UserId:  id_user,
		})
	}

	timeNow := time.Now()

	newUser := storage.User{

		ID:        id_user,
		FirstName: firstName,
		LastName:  lastName,
		Phone:     phone,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
		Car:       usersCar,
	}

	s.users = append(s.users, newUser)

	err := s.storage.Save(s.users)
	if err != nil {
		return storage.User{}, fmt.Errorf("failed to save user :%w ", err)
	}

	return newUser, nil

}
func (s *Service) GetUserById(id string) (storage.User, error) {

	uuid, err := uuid.Parse(id)
	if err != nil {
		return storage.User{}, helper.Invalid
	}

	for _, user := range s.users {
		if user.ID == uuid.String() {
			return user, nil
		}
	}

	return storage.User{}, helper.UserNotFound

}

func (s *Service) GetUsers() []storage.User {
	return s.users
}

func (s *Service) UpdateUser(
	id,
	firstName,
	lastName,
	phone string,
	cars []storage.CarDTO) (storage.User, error) {

	for i, user := range s.users {

		if user.ID == id {

			s.users[i].FirstName = firstName
			s.users[i].LastName = lastName
			s.users[i].Phone = phone
			s.users[i].UpdatedAt = time.Now()

			for _, carDTO := range cars {

				found := false

				for j, car := range s.users[i].Car {
					if car.ID == carDTO.CarDTO_ID {
						s.users[i].Car[j].Model = carDTO.Model
						s.users[i].Car[j].CarType = carDTO.Car_type

						found = true
						break
					}

				}

				if !found {
					return storage.User{}, fmt.Errorf(
						"%w : %s",
						helper.CarNotFound,
						carDTO.CarDTO_ID,
					)
				}

			}
			err := s.storage.Save(s.users)
			if err != nil {
				return storage.User{}, fmt.Errorf("failed to save user :%w ", err)

			}
			return s.users[i], nil

		}

	}
	return storage.User{},
		fmt.Errorf(
			"%w : %s",
			helper.UserNotFound,
			id)

}

func (s *Service) DeleteUser(id string) error {

	found := false

	for i, user := range s.users {
		if user.ID == id {
			found = true
			s.users = append(s.users[:i], s.users[i+1:]...)
			err := s.storage.Save(s.users)
			if err != nil {
				return err
			}
		}
	}

	if !found {
		return helper.UserNotFound
	}

	return nil

}

// save

/* UsetToSave := make([]storage.User, 0)

UsetToSave = append(UsetToSave, s.users...)

UsetToSave = append(UsetToSave, newUser)


	s.users = UsetToSave */
