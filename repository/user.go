package repository

import (
	"a21hc3NpZ25tZW50/db/filebased"
	"a21hc3NpZ25tZW50/model"
)

type UserRepository interface {
	GetUserByEmail(email string) (model.User, error)
	CreateUser(user model.User) (model.User, error)
	GetUserTaskCategory() ([]model.UserTaskCategory, error)
	GetUsers() ([]model.User, error)
}

type userRepository struct {
	filebasedDb *filebased.Data
}

func NewUserRepo(filebasedDb *filebased.Data) *userRepository {
	return &userRepository{filebasedDb}
}

func (r *userRepository) GetUserByEmail(email string) (model.User, error) {
	user, err := r.filebasedDb.GetUserByEmail(email)

	if err != nil {
		return model.User{}, err
	}
	return user, nil // TODO: replace this
}

func (r *userRepository) CreateUser(user model.User) (model.User, error) {
	createdUser, err := r.filebasedDb.CreateUser(user)

	if err != nil {
		return model.User{}, err
	}

	return createdUser, nil
}

func (r *userRepository) GetUserTaskCategory() ([]model.UserTaskCategory, error) {
	var UserTaskCategory []model.UserTaskCategory
	UserTaskCategory, err := r.filebasedDb.GetUserTaskCategory()

	return UserTaskCategory, err // TODO: replace this
}

func (r *userRepository) GetUsers() ([]model.User, error) {
	return r.filebasedDb.GetUsers()
}
