package user

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"io"
	"main/internal/infra/db"
	"main/internal/webserver"
	"net/http"
	"strconv"
)

type Webserver interface {
	RegisterHandler(string, func(http.ResponseWriter, *http.Request))
}

type CreateUserProcessorInterface interface {
	CreateUser(ctx context.Context, request CreateUserRequest) error
}

type UpdateUserProcessorInterface interface {
	UpdateUser(ctx context.Context, userId int, request UpdateUserRequest) error
}

type DeleteUserProcessorInterface interface {
	DeleteUser(ctx context.Context, userId int) error
}

type GetUserProcessorInterface interface {
	GetUser(ctx context.Context, userId int) (*User, error)
}

func RegisterHandlers(srv Webserver) {
	connection := db.GetConnection()
	repository := NewRepository(connection)
	userProcessor := NewUserProcessor(repository)

	srv.RegisterHandler("POST /user", webserver.WrapHandler(createUserHandler(userProcessor)))
	srv.RegisterHandler("PUT /user/{id}", webserver.WrapHandler(updateUserHandler(userProcessor)))
	srv.RegisterHandler("GET /user/{id}", webserver.WrapHandler(getUserHandler(userProcessor)))
	srv.RegisterHandler("DELETE /user/{id}", webserver.WrapHandler(deleteUserHandler(userProcessor)))
}

// @Summary      Create a new user
// @Description  Create a new user with specified First and Last names
// @Tags         user
// @Param request body CreateUserRequest true "Create user request"
// @Accept       json
// @Success      201
// @Failure      400
// @Failure      500
// @Router       /user [post]
func createUserHandler(createUserProcessor CreateUserProcessorInterface) func(*http.Request) ([]byte, error) {
	return func(request *http.Request) ([]byte, error) {
		var createUserRequest CreateUserRequest

		err := decodeAndValidate[CreateUserRequest](request, &createUserRequest)

		if err != nil {
			return nil, err
		}

		err = createUserProcessor.CreateUser(request.Context(), createUserRequest)

		if err != nil {
			return nil, err
		}

		return nil, nil
	}
}

// @Summary      Update a user
// @Description  Update a Last name or a First name of the specific user
// @Tags         user
// @Accept       json
// @Param request body UpdateUserRequest true "Update user request"
// @Param        id   path      int  true  "User ID"
// @Success      200
// @Failure      400
// @Failure      500
// @Router       /user/{id} [put]
func updateUserHandler(updateUserProcessor UpdateUserProcessorInterface) func(*http.Request) ([]byte, error) {
	return func(request *http.Request) ([]byte, error) {
		userId, err := strconv.Atoi(request.PathValue("id"))

		if err != nil {
			return nil, webserver.NewInvalidRequestError(err)
		}

		var updateUserRequest UpdateUserRequest

		err = decodeAndValidate[UpdateUserRequest](request, &updateUserRequest)

		if err != nil {
			return nil, err
		}

		err = updateUserProcessor.UpdateUser(request.Context(), userId, updateUserRequest)

		if err != nil {
			if errors.Is(err, &UserNotFoundError{}) {
				return nil, webserver.NewNotFoundError(err)
			}

			return nil, err
		}

		return nil, nil
	}
}

// @Summary      Find User By Id
// @Description  Get User Object By Id
// @Tags         user
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200 {object} User
// @Failure      400
// @Failure      404
// @Failure      500
// @Router       /user/{id} [get]
func getUserHandler(userProcessor GetUserProcessorInterface) func(*http.Request) ([]byte, error) {
	return func(request *http.Request) ([]byte, error) {
		userId, err := strconv.Atoi(request.PathValue("id"))

		if err != nil {
			return nil, webserver.NewInvalidRequestError(err)
		}

		user, err := userProcessor.GetUser(request.Context(), userId)

		if err != nil {
			if errors.Is(err, &UserNotFoundError{}) {
				return nil, webserver.NewNotFoundError(err)
			}

			return nil, err
		}

		encodedUser, err := json.Marshal(user)

		if err != nil {
			return nil, err
		}

		return encodedUser, nil
	}
}

// @Summary      Delete User By Id
// @Description  Delete User By Id
// @Tags         user
// @Param        id   path      int  true  "User ID"
// @Success      200
// @Failure      400
// @Failure      500
// @Router       /user/{id} [delete]
func deleteUserHandler(userProcessor DeleteUserProcessorInterface) func(*http.Request) ([]byte, error) {
	return func(request *http.Request) ([]byte, error) {
		userId, err := strconv.Atoi(request.PathValue("id"))

		if err != nil {
			return nil, webserver.NewInvalidRequestError(err)
		}

		err = userProcessor.DeleteUser(request.Context(), userId)

		if err != nil {
			return nil, err
		}

		return nil, nil
	}
}

func decodeAndValidate[T any](request *http.Request, dest *T) error {
	bodyValue, err := io.ReadAll(request.Body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(bodyValue, &dest)

	if err != nil {
		return webserver.NewInvalidRequestError(err)
	}

	vld := validator.New()
	err = vld.Struct(dest)

	if err != nil {
		return webserver.NewValidationError(err)
	}

	return nil
}
