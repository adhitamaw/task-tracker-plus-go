package service

import (
	"a21hc3NpZ25tZW50/model"
	repo "a21hc3NpZ25tZW50/repository"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// Untuk debug: ambil semua user
func (s *userService) GetUsers() ([]model.User, error) {
	return s.userRepo.GetUsers()
}

type UserService interface {
	Register(user *model.User) (model.User, error)
	Login(user *model.User) (token *string, err error)
	GetUserByEmail(email string) (model.User, error)
	GetUserTaskCategory() ([]model.UserTaskCategory, error)
	GetUsers() ([]model.User, error) // debug: ambil semua user
}

type userService struct {
	userRepo     repo.UserRepository
	sessionsRepo repo.SessionRepository
}

func NewUserService(userRepository repo.UserRepository, sessionsRepo repo.SessionRepository) UserService {
	return &userService{userRepository, sessionsRepo}
}

func (s *userService) Register(user *model.User) (model.User, error) {
	dbUser, err := s.userRepo.GetUserByEmail(user.Email)
	if err != nil {
		return *user, err
	}

	if dbUser.Email != "" || dbUser.ID != 0 {
		return *user, errors.New("email already exists")
	}

	// Hash password dengan bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return *user, err
	}
	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now()

	fmt.Printf("DEBUG REGISTER: Akan menyimpan user dengan email: %s, fullname: %s\n", user.Email, user.Fullname)
	newUser, err := s.userRepo.CreateUser(*user)
	if err != nil {
		fmt.Printf("DEBUG REGISTER: Gagal menyimpan user: %v\n", err)
		return *user, err
	}
	fmt.Printf("DEBUG REGISTER: User berhasil disimpan: %+v\n", newUser)

	return newUser, nil
}

func (s *userService) Login(user *model.User) (token *string, err error) {
	dbUser, err := s.userRepo.GetUserByEmail(user.Email)
	if err != nil {
		return nil, err
	}

	if dbUser.Email == "" || dbUser.ID == 0 {
		return nil, errors.New("user not found")
	}

	// Verifikasi password menggunakan bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		return nil, errors.New("wrong email or password")
	}

	claims := &model.Claims{
		ID:             dbUser.ID,
		Email:          dbUser.Email,
		StandardClaims: jwt.StandardClaims{}, // tanpa ExpiresAt, token tidak pernah expired
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := t.SignedString(model.JwtKey)
	if err != nil {
		return nil, err
	}

	session := model.Session{
		Token: tokenString,
		Email: user.Email,
		// Expiry dikosongkan karena token tidak pernah expired
	}

	_, err = s.sessionsRepo.SessionAvailEmail(session.Email)
	if err != nil {
		err = s.sessionsRepo.AddSessions(session)
		if err != nil {
			return nil, err
		}
	} else {
		err = s.sessionsRepo.UpdateSessions(session)
		if err != nil {
			return nil, err
		}
	}

	return &tokenString, nil
}

func (s *userService) GetUserByEmail(email string) (model.User, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (s *userService) GetUserTaskCategory() ([]model.UserTaskCategory, error) {
	// Cek koneksi ke repo
	users, err := s.userRepo.GetUsers()
	if err != nil {
		return nil, err
	}

	// Log jumlah user untuk debug
	fmt.Printf("DEBUG - Total users found: %d\n", len(users))

	userTaskCategories, err := s.userRepo.GetUserTaskCategory()
	if err != nil {
		return nil, err
	}

	// Log hasil
	fmt.Printf("DEBUG - Total UserTaskCategory entries: %d\n", len(userTaskCategories))

	return userTaskCategories, nil
}
