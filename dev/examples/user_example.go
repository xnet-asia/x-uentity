package examples

import (
	"github.com/xnetltd/x-uentity/domain"
	"github.com/xnetltd/x-uentity/repositories"
	"github.com/xnetltd/x-uentity/usecases"
)

type User struct {
	domain.BaseEntity
	Name  string
	Email string
}

type UserUsecase struct {
	*usecases.BaseUsecase[User]
}

func NewUserUsecase(repo repositories.Repository[User]) *UserUsecase {
	return &UserUsecase{
		BaseUsecase: usecases.NewBaseUsecase(repo),
	}
}

func (u *UserUsecase) GetUserByEmail(email string) ([]User, error) {
	return u.GetAll(func(user User) bool {
		return user.Email == email
	})
}

func Example() {
	// Initialize repository
	userRepo := repositories.NewInMemoryRepository[User]()

	// Initialize usecase
	userUsecase := NewUserUsecase(userRepo)

	// Create user
	user := User{
		BaseEntity: domain.BaseEntity{ID: "user-1"},
		Name:       "John Doe",
		Email:      "john@example.com",
	}

	userRepo.Create(user.ID, user)

	// Query users
	users, _ := userUsecase.GetAll(func(u User) bool {
		return u.Name == "John Doe"
	})

	println("Users found:", len(users))
}
