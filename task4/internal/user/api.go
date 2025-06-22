package user

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"io"
	"main/internal/infra/db"
	"net/http"
	"strconv"
)

type Webserver interface {
	RegisterHandler(string, func(http.ResponseWriter, *http.Request))
}

type CreateUserProcessorInterface interface {
	CreateUser(request CreateUserRequest) error
}

type UpdateUserProcessorInterface interface {
	UpdateUser(userId int, request UpdateUserRequest) error
}

type DeleteUserProcessorInterface interface {
	DeleteUser(userId int) error
}

type GetUserProcessorInterface interface {
	GetUser(userId int) (*User, error)
}

type UserRepositoryInterface interface {
	Get(userId int) (*User, error)
	Delete(userId int) error
	Create(user *User) error
	Update(user *User) error
}

func RegisterHandlers(srv Webserver) {
	connection := db.GetConnection()
	repository := NewRepository(connection)
	userProcessor := NewUserProcessor(repository)

	srv.RegisterHandler("POST /user", createUserHandler(userProcessor))
	srv.RegisterHandler("PUT /user/{id}", updateUserHandler(userProcessor))
	srv.RegisterHandler("GET /user/{id}", getUserHandler(userProcessor))
	srv.RegisterHandler("DELETE /user/{id}", deleteUserHandler(userProcessor))
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
func createUserHandler(createUserProcessor CreateUserProcessorInterface) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		var createUserRequest CreateUserRequest

		err := decodeAndValidate[CreateUserRequest](request, &createUserRequest)

		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)

			return
		}

		err = createUserProcessor.CreateUser(createUserRequest)

		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)

			return
		}

		writer.WriteHeader(http.StatusCreated)
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
func updateUserHandler(updateUserProcessor UpdateUserProcessorInterface) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		userId, err := strconv.Atoi(request.PathValue("id"))

		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)

			return
		}

		var updateUserRequest UpdateUserRequest

		err = decodeAndValidate[UpdateUserRequest](request, &updateUserRequest)

		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)

			return
		}

		err = updateUserProcessor.UpdateUser(userId, updateUserRequest)

		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)

			return
		}

		writer.WriteHeader(http.StatusOK)
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
func getUserHandler(userProcessor GetUserProcessorInterface) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		userId, err := strconv.Atoi(request.PathValue("id"))

		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)

			return
		}

		user, err := userProcessor.GetUser(userId)

		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)

			return
		}

		if user == nil {
			writer.WriteHeader(http.StatusNotFound)

			return
		}

		encodedUser, err := json.Marshal(user)

		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)

			return
		}

		writer.WriteHeader(http.StatusOK)
		writer.Header().Set("Content-Type", "application/json")
		writer.Write(encodedUser)
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
func deleteUserHandler(userProcessor DeleteUserProcessorInterface) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		userId, err := strconv.Atoi(request.PathValue("id"))

		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)

			return
		}

		err = userProcessor.DeleteUser(userId)

		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)

			return
		}

		writer.WriteHeader(http.StatusOK)
	}
}

func decodeAndValidate[T any](request *http.Request, dest *T) error {
	bodyValue, err := io.ReadAll(request.Body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(bodyValue, &dest)

	if err != nil {
		return nil
	}

	vld := validator.New()
	err = vld.Struct(dest)

	if err != nil {
		return err
	}

	return nil
}
