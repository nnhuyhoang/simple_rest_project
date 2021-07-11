package handler

import (
	"encoding/json"
	inerr "errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/consts"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/errors"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/model"
	"github.com/nnhuyhoang/simple_rest_project/backend/pkg/util"
	"gorm.io/gorm"
)

type loginData struct {
	UserId      int    `json:"userId"`
	Email       string `json:"email"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	FullName    string `json:"fullName"`
	RoleName    string `json:"roleName"`
	RoleCode    string `json:"roleCode"`
	AccessToken string `json:"accessToken"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type signupRequest struct {
	Firstname  string `json:"firstName" binding:"required"`
	LastName   string `json:"lastName"`
	Email      string `json:"email" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"rePassword" binding:"required"`
	Phone      string `json:"phone" binding:"required"`
	RoleCode   string `json:"roleCode" binding:"required"`
}

type loginResponse struct {
	Data loginData `json:"data"`
}

// LoginHandler --
// @Summary login
// @Description login
// @Tags Authentication
// @Accept	json
// @Produce  json
// @Security ApiKeyAuth
// @Param body body handler.loginRequest true "login request"
// @Success 200 {object} handler.loginResponse	"ok"
// @Failure 400 {object} model.Error
// @Router /auth/login [post]
func (h *Handler) LoginHandler(c *gin.Context) {
	var request loginRequest

	err := c.ShouldBindJSON(&request)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	loginResponse, err := h.doLoginWithMail(request)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, loginResponse)
}

// SignupHandler --
// @Summary signup
// @Description signup
// @Tags Authentication
// @Accept	json
// @Produce  json
// @Security ApiKeyAuth
// @Param body body handler.signupRequest true "signup request"
// @Success 200 {object} model.SuccessResponse	"ok"
// @Failure 400 {object} model.Error
// @Failure 404 {object} model.Error
// @Router /auth/signup [post]
func (h *Handler) SignupHandler(c *gin.Context) {
	var request signupRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	if !util.ValidateEmail(request.Email) {
		util.HandleError(c, errors.ErrInvalidEmail)
		return
	}
	if !util.ValidatePhoneNumber(request.Phone) {
		util.HandleError(c, errors.ErrInvalidPhoneNumber)
		return
	}

	if err = util.ValidatePasswordString(request.Password); err != nil {
		util.HandleError(c, err)
		return
	}

	if err = util.ValidatePasswordMatch(request.Password, request.RePassword); err != nil {
		util.HandleError(c, err)
		return
	}

	user, err := h.doSignUp(request)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	var userInter map[string]interface{}
	userRec, _ := json.Marshal(user)
	json.Unmarshal(userRec, &userInter)
	userSerial, err := util.SerializeStruct(userInter)
	if err != nil {
		util.HandleError(c, err)
		return
	}
	if user.Role.Code == consts.FieldTechnician {
		h.svs.EmailService.Send(consts.ExchangeNameSignUp, consts.WorkerRoutingKey, userSerial)
	} else {
		h.svs.EmailService.Send(consts.ExchangeNameSignUp, consts.ManagerRoutingKey, userSerial)
	}

	c.JSON(http.StatusOK, model.SuccessResponse{Code: http.StatusOK, Message: consts.SuccessMessage})
}

func (h *Handler) doSignUp(request signupRequest) (*model.User, error) {
	// _, err := h.repo.User.GetByEmail(h.store, request.Email)
	// if err == nil {
	// 	return nil, errors.ErrUserEmailAlreadyExisted
	// }

	// _, err := h.repo.User.GetByPhone(h.store, request.Phone)
	// if err == nil {
	// 	return nil, errors.ErrUserPhoneAlreadyExisted
	// }
	role, err := h.repo.Role.GetByRoleCode(h.store, request.RoleCode)
	if err != nil {
		return nil, errors.ErrUserRoleNotFound
	}

	hashedPassword, err := util.GenerateHashPassword(request.Password)
	if err != nil {
		return nil, errors.ErrInternalServerError
	}

	user := model.User{
		FirstName:      request.Firstname,
		LastName:       request.LastName,
		FullName:       request.Firstname + " " + request.LastName,
		Email:          request.Email,
		PhoneNumber:    request.Phone,
		RoleId:         int(role.Id),
		HashedPassword: hashedPassword,
	}

	insertedUser, err := h.repo.User.Create(h.store, user)
	if err != nil {
		return nil, errors.ErrInternalServerError
	}
	insertedUser.Role = role

	return insertedUser, nil
}

func (h *Handler) doLoginWithMail(request loginRequest) (*loginResponse, error) {
	if !util.ValidateEmail(request.Email) {
		h.log.Errorf("[handler.doLoginWithMail] ValidateEmail(), email=%s", request.Email)
		return nil, errors.ErrInvalidEmail
	}

	passwordErr := util.ValidatePasswordString(request.Password)
	if passwordErr != nil {
		h.log.Errorf("[handler.doLoginWithMail] ValidatePasswordString()")
		return nil, passwordErr
	}

	user, err := h.repo.User.GetByEmail(h.store, request.Email)
	if err != nil {
		if inerr.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.ErrUserNotFound
		}
		return nil, errors.ErrInternalServerError
	}

	if err = util.ValidatePassword(user.HashedPassword, request.Password); err != nil {
		h.log.Errorf("[handler.doLoginWithMail] ValidatePassword()", err)
		return nil, errors.ErrIncorrectEmailOrPassword
	}

	return h.buildTokenResponse(user, "", "")
}

//buildTokenResponse ...
func (h *Handler) buildTokenResponse(user *model.User, deviceID, userPodID string) (*loginResponse, error) {
	authenInfo := model.AuthenticationInfo{
		UserID:   int(user.Id),
		Email:    user.Email,
		RoleCode: user.Role.Code,
	}

	expireAt := h.cfg.GetTokenTTL().Unix()
	accesstoken, err := util.GenerateJWTToken(&authenInfo, expireAt, h.cfg.JWTSecret)
	if err != nil {
		return nil, errors.ErrGenJWTFailed
	}

	return &loginResponse{
		Data: loginData{
			UserId:      int(user.Id),
			Email:       user.Email,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			FullName:    user.FullName,
			RoleName:    user.Role.Name,
			RoleCode:    user.Role.Code,
			AccessToken: accesstoken,
		},
	}, nil
}
