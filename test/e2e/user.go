package e2e

// User
type User struct {
	Id       uint64
	Username string
	Email    string
	Password string
}

// UserService
type UserService struct {
	store  map[uint64]*User
	nextId uint64
}

func NewUserService() *UserService {
	return &UserService{
		store:  make(map[uint64]*User),
		nextId: 1,
	}
}

func (u *UserService) CreateUser(username string, email string, password string) (*User, error) {
	user := User{
		Id:       u.nextId,
		Username: username,
		Email:    email,
		Password: password,
	}

	u.store[user.Id] = &user
	u.nextId++

	return &user, nil
}

func (u *UserService) FindUserByEmail(email string) *User {
	for _, user := range u.store {
		if user.Email == email {
			return user
		}
	}

	return nil
}

// UserController
type UserController struct {
	UserService *UserService
}

func NewUserController(userService *UserService) *UserController {
	return &UserController{
		UserService: userService,
	}
}
