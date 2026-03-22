package handler

import (
	"net/http"

	"task-manager/internal/model"
	"task-manager/internal/service"

	"github.com/gin-gonic/gin"
)

// Register godoc
// @Summary Регистрация пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param input body model.RegisterInput true "Данные регистрации"
// @Success 201 {object} model.User
// @Failure 400 {object} model.ErrorResponse
// @Router /api/auth/register [post]
func Register(c *gin.Context) {
	var input model.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error: "Некорректные данные: " + err.Error(),
			Code:  http.StatusBadRequest,
		})
		return
	}

	user, err := service.Register(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error: err.Error(),
			Code:  http.StatusBadRequest,
		})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// Login godoc
// @Summary Вход в систему
// @Tags auth
// @Accept json
// @Produce json
// @Param input body model.LoginInput true "Данные для входа"
// @Success 200 {object} model.TokenResponse
// @Failure 401 {object} model.ErrorResponse
// @Router /api/auth/login [post]
func Login(c *gin.Context) {
	var input model.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error: "Некорректные данные: " + err.Error(),
			Code:  http.StatusBadRequest,
		})
		return
	}

	token, err := service.Login(input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{
			Error: err.Error(),
			Code:  http.StatusUnauthorized,
		})
		return
	}

	c.JSON(http.StatusOK, model.TokenResponse{Token: token})
}
