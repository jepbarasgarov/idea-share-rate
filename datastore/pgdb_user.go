package datastore

import (
	"belli/onki-game-ideas-mongo-backend/models"
	"belli/onki-game-ideas-mongo-backend/responses"
	"context"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	sqlUserGetByUsername   = `SELECT id, firstname, lastname , role, status, password FROM tbl_user  WHERE username = $1`
	sqlUserGetByID         = `SELECT username, firstname, lastname, role, status FROM tbl_user WHERE id = $1`
	sqlUserGetPasswordByID = `SELECT password FROM tbl_user WHERE id=$1`
	sqlUserCreate          = `INSERT INTO tbl_user(username, password, firstname, lastname, status, role) VALUES($1, $2, $3, $4, $5, $6) RETURNING id`
	sqlSelectDirectorIDs   = `SELECT ARRAY(SELECT id FROM tbl_user WHERE role = $1 AND status = $2)`
	sqlSelectDirectorList  = `SELECT us.id, us.firstname, us.lastname, position.name, department.id, department.name FROM tbl_user us
	INNER JOIN tbl_department department ON us.department_id = department.id
	INNER JOIN tbl_position position ON us.position_id = position.id
	WHERE us.role = $1 AND us.status = $2`
	sqlSelectStaffDepList = `SELECT us.id, us.firstname, us.lastname, position.name, department.id, department.name FROM tbl_user us
	INNER JOIN tbl_department department ON us.department_id = department.id
	INNER JOIN tbl_position position ON us.position_id = position.id
	WHERE us.role = $1 AND us.status = $2 AND department.id = $3`
	sqlSelectHeadDepList = `SELECT us.id, us.firstname, us.lastname, position.name, department.id, department.name FROM tbl_user us
	INNER JOIN tbl_department department ON us.department_id = department.id
	INNER JOIN tbl_position position ON us.position_id = position.id
	WHERE us.role = $1 AND us.status = $2`
	sqlGetUserList                    = `SELECT id, username, firstname, lastname, role, status FROM tbl_user WHERE id <> $1 LIMIT $2 OFFSET $3`
	sqlCountUsers                     = `SELECT COUNT(*) FROM tbl_user WHERE  id <> $1`
	sqlRemoveHeadOfDepartment         = `UPDATE tbl_user SET role = $1 WHERE department_id = $2 AND role = $3 AND status <> $4`
	sqlUpdateUserOwnData              = `UPDATE tbl_user SET phone = $1, email = $2, update_ts = $3 WHERE id = $4`
	sqlUpdateUserPassword             = `UPDATE tbl_user SET password = $1, update_ts = $2 WHERE id = $3`
	sqlGetLastDocListGotTimeOfUser    = `SELECT last_doc_get FROM tbl_user WHERE id = $1`
	sqlUpdateLastDocListGotTimeOfUser = `UPDATE tbl_user SET last_doc_get = $1 WHERE id = $2`
	sqlUserAutocompleteList           = `SELECT id, role, firstname, lastname from tbl_user WHERE status <> $1`
	sqlUpdateUser                     = `UPDATE tbl_user SET username = $1, firstname = $2, lastname = $3,role = $4, status = $5 WHERE id = $6`
	sqlDeleteUser                     = `DELETE FROM tbl_user WHERE id = $1`
	sqlSetTwoFAUser                   = `UPDATE tbl_user SET two_fa_key = $1 WHERE id = $2`
)

//////////////////////////////////////////////////////////////////////MONGO/////////////////////////////////////////////////////////////////////////////////

func (d *MgAccess) UserGetByUsername(
	ctx context.Context,
	username string,
) (item *models.UserSpecDataBson, err error) {
	clog := log.WithFields(log.Fields{
		"method": "PgAccess.UserGetByUsername",
	})

	defer func() {
		if err != nil {
			item = nil
		}
	}()
	item = &models.UserSpecDataBson{}

	db := d.client.Database("idea-share")
	coll := db.Collection("user")

	var u models.UserSpecDataBson
	err = coll.FindOne(ctx, bson.M{"username": username}).Decode(&u)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
			item = nil
			return
		}
		eMsg := "Error in FindOne"
		clog.WithError(err).Error(eMsg)
	}
	item = &u

	return
}

func (d *MgAccess) UserCreate(
	ctx context.Context,
	user *models.UserCreate,
) (item *models.UserSpecDataBson, err error) {
	clog := log.WithFields(log.Fields{
		"method": "MgAccess.UserCreate",
	})

	defer func() {
		if err != nil {
			item = nil
		}
	}()

	db := d.client.Database("idea-share")
	coll := db.Collection("user")

	row, err := coll.InsertOne(ctx, bson.D{
		{Key: "firstname", Value: user.Firstname},
		{Key: "lastname", Value: user.Lastname},
		{Key: "username", Value: user.Username},
		{Key: "role", Value: user.Role},
		{Key: "status", Value: responses.Active},
		{Key: "password", Value: user.Password},
		{Key: "create_ts", Value: time.Now().UTC()},
		{Key: "update_ts", Value: time.Now().UTC()},
	})
	if err != nil {
		eMsg := "An error occurred on Insert one"
		clog.WithError(err).Error(eMsg)
		return
	}

	item = &models.UserSpecDataBson{
		ID:        row.InsertedID.(primitive.ObjectID),
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Username:  user.Username,
		Role:      user.Role,
		Status:    responses.Active,
	}

	return
}

func (d *MgAccess) UserUpdate(
	ctx context.Context,
	user *models.UserUpdate,
) (item *models.UserSpecDataBson, err error) {
	clog := log.WithFields(log.Fields{
		"method": "PgAccess.UserUpdate",
	})

	db := d.client.Database("idea-share")
	workerColl := db.Collection("user")

	filterUser := bson.M{"_id": user.ID}

	updateUser := bson.M{"$set": bson.M{
		"firstname": user.Firstname,
		"lastname":  user.Lastname,
		"username":  user.Username,
		"role":      user.Role,
		"status":    user.Status,
		"update_ts": time.Now().UTC(),
	}}

	_, err = workerColl.UpdateOne(ctx, filterUser, updateUser)
	if err != nil {
		eMsg := "error in Updating user"
		clog.WithError(err).Error(eMsg)
		err = errors.Wrap(err, eMsg)
		return
	}

	item = &models.UserSpecDataBson{
		ID:        user.ID,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Username:  user.Username,
		Role:      user.Role,
		Status:    user.Status,
	}

	return
}

func (d *MgAccess) UserUpdateOwnPassword(
	ctx context.Context,
	cu *responses.ActionInfo,
	newPassword string,
) (err error) {
	clog := log.WithFields(log.Fields{
		"method": "PgAccess.UserGetByID",
	})

	db := d.client.Database("idea-share")
	workerColl := db.Collection("user")

	filterUser := bson.M{"_id": cu.ID}

	updateUser := bson.M{"$set": bson.M{
		"password": newPassword,
	}}

	_, err = workerColl.UpdateOne(ctx, filterUser, updateUser)
	if err != nil {
		eMsg := "error in Updating user's password"
		clog.WithError(err).Error(eMsg)
		err = errors.Wrap(err, eMsg)
		return
	}

	return
}

func (d *MgAccess) AdminUpdatePassword(
	ctx context.Context,
	cu *responses.ActionInfo,
	userid primitive.ObjectID,
	newPassword string,
) (err error) {
	clog := log.WithFields(log.Fields{
		"method": "PgAccess.AdminUpdatePassword",
	})

	db := d.client.Database("idea-share")
	workerColl := db.Collection("user")

	filterUser := bson.M{"_id": userid}

	updateUser := bson.M{"$set": bson.M{
		"password": newPassword,
	}}

	_, err = workerColl.UpdateOne(ctx, filterUser, updateUser)
	if err != nil {
		eMsg := "error in Updating user's password"
		clog.WithError(err).Error(eMsg)
		err = errors.Wrap(err, eMsg)
		return
	}

	return
}

func (d *MgAccess) UserDelete(
	ctx context.Context,
	id primitive.ObjectID,
) (err error) {
	clog := log.WithFields(log.Fields{
		"method": "PgAccess.UserList",
	})

	db := d.client.Database("idea-share")
	coll := db.Collection("user")

	_, err = coll.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		eMsg := "Error in delete user"
		clog.WithError(err).Error(eMsg)
		return
	}

	return

}

//GET

func (d *MgAccess) UserGetByID(
	ctx context.Context,
	id primitive.ObjectID,
) (item *models.UserSpecDataBson, err error) {
	clog := log.WithFields(log.Fields{
		"method": "PgAccess.UserGetByID",
	})

	defer func() {
		if err != nil {
			item = nil
		}
	}()

	db := d.client.Database("idea-share")
	coll := db.Collection("user")

	var u models.UserSpecDataBson
	err = coll.FindOne(ctx, bson.M{"_id": id}).Decode(&u)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
			item = nil
			return
		}
		eMsg := "Error in Find user with ID"
		clog.WithError(err).Error(eMsg)
		return
	}

	u.HashedPassword = ""

	item = &u

	return
}

func (d *MgAccess) UserGetPasswordByID(
	ctx context.Context,
	id primitive.ObjectID,
) (pwdHash string, err error) {

	clog := log.WithFields(log.Fields{
		"method": "PgAccess.UserGetPasswordByID",
	})

	db := d.client.Database("idea-share")
	coll := db.Collection("user")

	options := options.FindOne()
	options.Projection = bson.M{"_id": 0, "password": 1}

	var u bson.M
	err = coll.FindOne(ctx, bson.M{"_id": id}, options).Decode(&u)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = nil
			pwdHash = ""
			return
		}
		eMsg := "Error in Find user's password with ID"
		clog.WithError(err).Error(eMsg)
		return
	}

	pwdHash = u["password"].(string)

	return
}

func (d *MgAccess) UserAutocompleteList(
	ctx context.Context,
) (item *[]models.UserLightDataBson, err error) {
	clog := log.WithFields(log.Fields{
		"method": "PgAccess.UserAutocompleteList",
	})

	defer func() {
		if err != nil {
			item = nil
		}
	}()

	db := d.client.Database("idea-share")
	coll := db.Collection("user")

	options := options.Find()
	options.Projection = bson.M{"firstname": 1, "lastname": 1, "role": 1}

	cursor, err := coll.Find(ctx, bson.M{}, options)
	if err != nil {
		eMsg := "Error in Find user list"
		clog.WithError(err).Error(eMsg)
		return
	}
	var users []models.UserLightDataBson

	if err = cursor.All(ctx, &users); err != nil {
		eMsg := "Error in reading cursor"
		clog.WithError(err).Error(eMsg)
		return
	}

	item = &users

	return
}
