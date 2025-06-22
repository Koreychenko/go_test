package user

type UserProcessor struct {
	repository UserRepositoryInterface
}

func (p *UserProcessor) CreateUser(request CreateUserRequest) error {
	user := User{
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Password:  request.Password,
	}

	return p.repository.Create(&user)
}

func (p *UserProcessor) UpdateUser(userId int, request UpdateUserRequest) error {
	user, err := p.repository.Get(userId)

	if err != nil {
		return nil
	}
	if len(request.FirstName) > 0 {
		user.FirstName = request.FirstName
	}

	if len(request.LastName) > 0 {
		user.FirstName = request.LastName
	}

	return p.repository.Update(user)
}

func (p *UserProcessor) GetUser(userId int) (*User, error) {
	user, err := p.repository.Get(userId)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (p *UserProcessor) DeleteUser(userId int) error {
	return p.repository.Delete(userId)
}

func NewUserProcessor(repository UserRepositoryInterface) *UserProcessor {
	return &UserProcessor{repository}
}
