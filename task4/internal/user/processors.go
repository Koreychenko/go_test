package user

import "context"

type UserNotFoundError struct {
}

func (e *UserNotFoundError) Error() string {
	return "user not foud"
}

func NewUserNotFoundError() *UserNotFoundError {
	return &UserNotFoundError{}
}

type UserRepositoryInterface interface {
	Get(ctx context.Context, userId int) (*User, error)
	Delete(ctx context.Context, userId int) error
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
}

type CrudProcessor struct {
	repository UserRepositoryInterface
}

func (p *CrudProcessor) CreateUser(ctx context.Context, request CreateUserRequest) error {
	user := User{
		FirstName: request.FirstName,
		LastName:  request.LastName,
	}

	return p.repository.Create(ctx, &user)
}

func (p *CrudProcessor) UpdateUser(ctx context.Context, userId int, request UpdateUserRequest) error {
	user, err := p.repository.Get(ctx, userId)

	if err != nil {
		return nil
	}

	if user == nil {
		return NewUserNotFoundError()
	}

	if len(request.FirstName) > 0 {
		user.FirstName = request.FirstName
	}

	if len(request.LastName) > 0 {
		user.FirstName = request.LastName
	}

	return p.repository.Update(ctx, user)
}

func (p *CrudProcessor) GetUser(ctx context.Context, userId int) (*User, error) {
	user, err := p.repository.Get(ctx, userId)

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, NewUserNotFoundError()
	}

	return user, nil
}

func (p *CrudProcessor) DeleteUser(ctx context.Context, userId int) error {
	return p.repository.Delete(ctx, userId)
}

func NewUserProcessor(repository UserRepositoryInterface) *CrudProcessor {
	return &CrudProcessor{repository}
}
