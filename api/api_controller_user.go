package api

import (
	"belli/onki-game-ideas-mongo-backend/config"
	"belli/onki-game-ideas-mongo-backend/errs"
	"belli/onki-game-ideas-mongo-backend/helpers"
	"belli/onki-game-ideas-mongo-backend/models"
	"belli/onki-game-ideas-mongo-backend/responses"
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

////////////////////////////////////////////////////////////////////////////////MONGO////////////////////////////////////////////////////////////////////////////////////

func (api *APIController) UserLogin(
	ctx context.Context,
	username string,
	password string,
) (item *models.UserSpecDataBson, tokens *models.Tokens, err error) {
	clog := log.WithContext(ctx).WithFields(log.Fields{
		"method":   "api.UserLogin",
		"username": username,
	})

	user, err := api.access.UserGetByUsername(ctx, username)
	if err != nil {
		eMsg := "error in api.access.UserGetByUsername"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}
	if user == nil || user.Status == responses.Blocked {
		eMsg := "user not found"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorUnauthorized(errs.ERR_UA)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		eMsg := "error in bcrypt.CompareHashAndPassword"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorUnauthorized(errs.ERR_UA)
		return
	}

	tokens = &models.Tokens{}

	//----------AccessToken-------------------------
	tokens.AccessToken, err = helpers.GenerateAccessToken(
		username,
		config.Conf.Jwt.AccessTokenExpiry,
		config.Conf.Jwt.Secret,
	)
	if err != nil {
		eMsg := "Error at GenerateAccessToken"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	//----------RefreshToken------------------------
	tokens.RefreshToken, err = helpers.GenerateRefreshToken()
	if err != nil {
		eMsg := "Error at GenerateRefreshToken"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	item = &models.UserSpecDataBson{
		ID:        user.ID,
		Username:  user.Username,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Role:      user.Role,
	}

	err = api.cache.TokenSetWithExpiry(
		ctx,
		user.ID.Hex(),
		tokens.RefreshToken,
		time.Duration(config.Conf.Jwt.RefreshTokenExpiry)*time.Second,
	)
	if err != nil {
		eMsg := "error in api.cache.TokenSetWithExpiry"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	return
}

func (api *APIController) UserGiveAccessToken(
	ctx context.Context,
	username string,
	refresh string,
) (item *models.Tokens, err error) {

	clog := log.WithContext(ctx).WithFields(log.Fields{
		"method": "api.UserGiveAccessToken",
	})

	Tokens := models.Tokens{}

	if username == "<nil>" {
		eMsg := "owner of request can't be determined"
		clog.Error(eMsg)
		err = errs.NewHttpErrorUnauthorized(errs.ERR_UA)
		return
	}
	user, err := api.access.UserGetByUsername(ctx, username)
	if err != nil {
		eMsg := "error in api.access.UserGetByUsername"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	if user == nil || user.Status == responses.Blocked {
		eMsg := "owner of request can't be determined"
		clog.Error(eMsg)
		err = errs.NewHttpErrorUnauthorized(errs.ERR_UA)
		return
	}

	userID, err := api.cache.GetUserIDByToken(ctx, refresh)
	if err != nil {
		eMsg := "can't get refreshToken from redis"
		clog.Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	if userID == nil {
		eMsg := "refresh token has expired"
		clog.Error(eMsg)
		err = errs.NewHttpErrorUnauthorized(errs.ERR_UA)
		return
	}

	primitiveObjectIDForm, _ := primitive.ObjectIDFromHex(*userID)

	if user.ID != primitiveObjectIDForm {
		eMsg := "refresh token is not yours"
		clog.Error(eMsg)
		err = errs.NewHttpErrorUnauthorized(errs.ERR_UA)
		return
	}

	if user.Status == responses.Blocked {
		eMsg := "user has been blocked"
		clog.Error(eMsg)
		err = errs.NewHttpErrorUnauthorized(errs.ERR_UA)
		return
	}

	//----------AccessToken-------------------------
	Tokens.AccessToken, err = helpers.GenerateAccessToken(
		username,
		config.Conf.Jwt.AccessTokenExpiry,
		config.Conf.Jwt.Secret,
	)
	if err != nil {
		eMsg := "Error at GenerateAccessToken"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	//----------RefreshToken------------------------
	Tokens.RefreshToken, err = helpers.GenerateRefreshToken()
	if err != nil {
		eMsg := "Error at GenerateRefreshToken"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	err = api.cache.TokenSetWithExpiry(
		ctx,
		user.ID.Hex(),
		Tokens.RefreshToken,
		time.Duration(config.Conf.Jwt.RefreshTokenExpiry)*time.Second,
	)
	if err != nil {
		eMsg := "error in api.cache.TokenSetWithExpiry"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}
	_ = api.cache.DeletePreviousRefreshToken(ctx, refresh)

	item = &Tokens

	return
}

func (api *APIController) UserCreate(
	ctx context.Context,
	user *models.UserCreate,
	password string,
	cu *responses.ActionInfo,
) (item *models.UserSpecDataBson, err error) {

	clog := log.WithContext(ctx).WithFields(log.Fields{
		"method":   "api.UserCreate",
		"username": cu.Username,
	})

	usr, err := api.access.UserGetByUsername(ctx, user.Username)
	if err != nil {
		eMsg := "error in api.access.UserGetByUsername"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	if usr != nil {
		eMsg := "Username is in use"
		clog.Error(eMsg)
		err = errs.NewHttpErrorConflict(errs.ERR_UNIQUE_USER)
		return
	}

	pwdHashBytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		eMsg := "error in bcrypt.GenerateFromPassword"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}
	user.Password = string(pwdHashBytes)

	item, err = api.access.UserCreate(ctx, user)
	if err != nil {
		eMsg := "error in api.access.UserCreate"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	return
}

func (api *APIController) UserUpdate(
	ctx context.Context,
	user *models.UserUpdate,
) (item *models.UserSpecDataBson, err error) {

	clog := log.WithContext(ctx).WithFields(log.Fields{
		"method":   "api.UserUpdate",
		"username": user.Username,
	})

	oldUserData, err := api.access.UserGetByID(ctx, user.ID)
	if err != nil {
		eMsg := "error in api.access.UserGetByID"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	if oldUserData == nil {
		eMsg := "user not found"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorNotFound(errs.ERR_NF_USER)
		return
	}

	usr, err := api.access.UserGetByUsername(ctx, user.Username)
	if err != nil {
		eMsg := "error in api.access.UserGetByUsername"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	if usr != nil && usr.ID != user.ID {
		eMsg := "Username is in use"
		clog.Error(eMsg)
		err = errs.NewHttpErrorConflict(errs.ERR_UNIQUE_USER)
		return
	}

	item, err = api.access.UserUpdate(ctx, user)
	if err != nil {
		eMsg := "error in api.access.UserUpdate"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	return
}

func (api *APIController) UserUpdateOwnPassword(
	ctx context.Context,
	cu *responses.ActionInfo,
	oldPassword, newPassword string,
) (err error) {
	clog := log.WithContext(ctx).WithFields(log.Fields{
		"method":   "api.UserUpdateOwnPassword",
		"username": cu.Username,
	})

	pwdHash, err := api.access.UserGetPasswordByID(ctx, cu.ID)
	if err != nil {
		eMsg := "error in api.access.UserGetPasswordByID"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	if pwdHash == "" {
		eMsg := "password is not found"
		clog.Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(pwdHash), []byte(oldPassword))
	if err != nil {
		eMsg := "old password is not correct"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorForbidden(errs.ERR_FB_pwrd_USER)
		return
	}

	newPwrdHashByte, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		eMsg := "error in bcrypt.GenerateFromPassword"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	newPwrdHash := string(newPwrdHashByte)

	err = api.access.UserUpdateOwnPassword(ctx, cu, newPwrdHash)
	if err != nil {
		eMsg := "error in api.access.UserUpdateOwnPassword"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	return nil
}

func (api *APIController) AdminUpdatePassword(
	ctx context.Context,
	cu *responses.ActionInfo,
	id primitive.ObjectID,
	password string,
) (err error) {
	clog := log.WithContext(ctx).WithFields(log.Fields{
		"method":   "api.AdminUpdatePassword",
		"username": cu.Username,
	})

	user, err := api.access.UserGetByID(ctx, id)
	if err != nil {
		eMsg := "error in api.access.UserGetByID"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}
	if user == nil {
		eMsg := "user is not found"
		clog.Error(eMsg)
		err = errs.NewHttpErrorNotFound(errs.ERR_NF_USER)
		return
	}
	if user.ID == cu.ID {
		eMsg := "you can't update own password"
		clog.Error(eMsg)
		err = errs.NewHttpErrorForbidden(errs.ERR_FB_ownpwrd_USER)
		return
	}
	pwdHashBytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		eMsg := "error in bcrypt.GenerateFromPassword"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	pwdHash := string(pwdHashBytes)

	err = api.access.AdminUpdatePassword(ctx, cu, id, pwdHash)
	if err != nil {
		eMsg := "error in api.access.AdminUpdatePassword"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	return nil
}

func (api *APIController) UserDelete(
	ctx context.Context,
	cu *responses.ActionInfo,
	id primitive.ObjectID,
) (err error) {
	clog := log.WithContext(ctx).WithFields(log.Fields{
		"method":   "api.UserDelete",
		"username": cu.Username,
	})

	user, err := api.access.UserGetByID(ctx, id)
	if err != nil {
		eMsg := "error in api.access.UserGetByID"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}
	if user == nil {
		eMsg := "user is not found"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorNotFound(errs.ERR_NF_USER)
		return
	}
	if cu.ID == id {
		eMsg := "you can't delete yourself"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorForbidden(errs.ERR_FB_delete_USER)
		return
	}
	err = api.access.UserDelete(ctx, id)
	if err != nil {
		eMsg := "error in api.access.UserDelete"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	return nil
}

//GET

func (api *APIController) UserGet(
	ctx context.Context,
	cu *responses.ActionInfo,
	ID primitive.ObjectID,
) (item *models.UserSpecDataBson, err error) {
	clog := log.WithContext(ctx).WithFields(log.Fields{
		"method":   "api.UserGet",
		"username": cu.Username,
	})
	item, err = api.access.UserGetByID(ctx, ID)
	if err != nil {
		eMsg := "error in api.access.UserGetByID"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	if item == nil {
		eMsg := "user is not found"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorNotFound(errs.ERR_NF_USER)
		return
	}
	return
}

func (api *APIController) UserAutocompleteList(
	ctx context.Context,
	cu *responses.ActionInfo,
) (item *[]models.UserLightDataBson, err error) {
	clog := log.WithContext(ctx).WithFields(log.Fields{
		"method":   "api.UserAutocompleteList",
		"username": cu.Username,
	})

	item, err = api.access.UserAutocompleteList(ctx)
	if err != nil {
		eMsg := "error in api.access.UserAutocompleteList"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}
	return
}

func (api *APIController) UserGetByUsername(ctx context.Context, username string) (item *models.UserSpecDataBson, err error) {

	clog := log.WithContext(ctx).WithFields(log.Fields{
		"method":   "api.UserGetByUsername",
		"username": username,
	})

	item, err = api.access.UserGetByUsername(ctx, username)
	if err != nil {
		eMsg := "error in api.access.UserGetByUsername"
		clog.WithError(err).Error(eMsg)
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	return
}
