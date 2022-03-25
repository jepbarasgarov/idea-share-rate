package errs

import (
	"belli/onki-game-ideas-mongo-backend/responses"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var ErrFileResponse map[string]responses.Language

func ReadErrorFile(source string) (err error) {

	var raw []byte
	raw, err = ioutil.ReadFile(source)
	if err != nil {
		wMsg := "error reading err-message from file, creating new sample"
		log.Warn(wMsg)

		err = createDefaultErrMessage(source)
		if err != nil {
			eMsg := "error creating err-message file"
			log.WithError(err).Error(eMsg)
			err = errors.Wrap(err, eMsg)
			return
		}

		raw, err = ioutil.ReadFile(source)
		if err != nil {
			eMsg := "error reading err-message from file"
			log.WithError(err).Error(eMsg)
			err = errors.Wrap(err, eMsg)
			return
		}
	}

	err = json.Unmarshal(raw, &ErrFileResponse)
	if err != nil {
		fmt.Println(err)
		return
	}

	return
}

type NotFoundKey string
type ConflictKey string
type ForbidenKey string
type NotAllowedKey string
type BadRequestKey string
type UnauthorizedKey string
type FileTooLargeKey string
type TooManyRequestKey string
type InternalErrorKey string
type LockedKey string

//keys
const (
	//common
	ERR_NA NotAllowedKey     = "ERR_NA"
	ERR_BR BadRequestKey     = "ERR_BR"
	ERR_UA UnauthorizedKey   = "ERR_UA"
	ERR_FL FileTooLargeKey   = "ERR_FL"
	ERR_MR TooManyRequestKey = "ERR_MR"
	ERR_IE InternalErrorKey  = "ERR_IE"

	//user
	ERR_NF_USER NotFoundKey = "ERR_NF_USER"

	ERR_UNIQUE_USER ConflictKey = "ERR_UNIQUE_USER"

	ERR_FB_USER         ForbidenKey = "ERR_FB_USER"
	ERR_FB_owndata_USER ForbidenKey = "ERR_FB_owndata_USER"
	ERR_FB_blocked_USER ForbidenKey = "ERR_FB_blocked_USER"
	ERR_FB_pwrd_USER    ForbidenKey = "ERR_FB_pwrd_USER"
	ERR_FB_ownpwrd_USER ForbidenKey = "ERR_FB_ownpwrd_USER"
	ERR_FB_delete_USER  ForbidenKey = "ERR_FB_delete_USER"

	//worker

	ERR_NF_WORKER       NotFoundKey = "ERR_NF_WORKER"
	ERR_UNIQUE_WORKER   ConflictKey = "ERR_UNIQUE_WORKER"
	ERR_WORKER_HAS_IDEA ForbidenKey = "ERR_WORKER_HAS_IDEA"

	//Idea

	ERR_NF_IDEA       NotFoundKey = "ERR_NF_IDEA"
	ERR_IP_RESTRICTED ConflictKey = "ERR_IP_RESTRICTED"

	ERR_NF_CRITERIA        NotFoundKey = "ERR_NF_CRITERIA"
	ERR_UNIQUE_CRITERIA    ConflictKey = "ERR_UNIQUE_CRITERIA"
	ERR_CRITERIA_HAS_RATES ForbidenKey = "ERR_CRITERIA_HAS_RATES"

	ERR_NF_GENRE       NotFoundKey = "ERR_NF_GENRE"
	ERR_UNIQUE_GENRE   ConflictKey = "ERR_UNIQUE_GENRE"
	ERR_GENRE_has_IDEA ForbidenKey = "ERR_GENRE_has_IDEA"

	ERR_NF_MECH       NotFoundKey = "ERR_NF_MECH"
	ERR_UNIQUE_MECH   ConflictKey = "ERR_UNIQUE_MECH"
	ERR_MECH_has_IDEA ForbidenKey = "ERR_MECH_has_IDEA"
)

//codes
const (
	ErrorCodeOK                  int = 200
	ErrorCodeTfaRequired         int = 250
	ErrorCodeBadRequest          int = 400
	ErrorCodeUnauthorized        int = 401
	ErrorCodeForbidden           int = 403
	ErrorCodeNotFound            int = 404
	ErrorMethodNotAllowed        int = 405
	ErrorCodeConflict            int = 409
	ErrorCodeExpired             int = 408
	ErrorCodeFileSizeTooLarge    int = 413
	ErrorCodeLocked              int = 423
	ErrorCodeTooManyRequests     int = 429
	ErrorCodeInternalServerError int = 500
)

//messages
const (
	ErrorMessageOK                  = "ok"
	ErrorMessageTfaRequired         = "tfa_required"
	ErrorMessageBadRequest          = "bad_request"
	ErrorMessageUnauthorized        = "unauthorized"
	ErrorMessageForbidden           = "forbidden"
	ErrorMessageNotFound            = "not_found"
	ErrorMessageMethodNotAllowed    = "method_not_allowed"
	ErrorMessageConflict            = "conflict"
	ErrorMessageExpired             = "expired"
	ErrorMessageFileSizeTooLarge    = "file_size_too_large"
	ErrorMessageTooManyRequests     = "otp_retry_limit_exceeded"
	ErrorMessageInternalServerError = "internal_server_error"
	ErrorMessageLocked              = "locked"
)

//structs
type errNotAllowed struct {
	key NotAllowedKey
}
type errBadRequest struct {
	key BadRequestKey
}
type errUnauthorized struct {
	key UnauthorizedKey
}
type errFileTooLarge struct {
	key FileTooLargeKey
}
type errTooManyRequest struct {
	key TooManyRequestKey
}
type errInternalError struct {
	key InternalErrorKey
}
type errNotFound struct {
	key NotFoundKey
}
type errForbidden struct {
	key ForbidenKey
}
type errConflict struct {
	key ConflictKey
}

type errLocked struct {
	key LockedKey
}

// Methods
func (x *errForbidden) Error() string {
	return string(x.key)
}
func (x *errConflict) Error() string {
	return string(x.key)
}
func (x *errNotFound) Error() string {
	return string(x.key)
}
func (x *errNotAllowed) Error() string {
	return string(x.key)
}
func (x *errBadRequest) Error() string {
	return string(x.key)
}
func (x *errUnauthorized) Error() string {
	return string(x.key)
}
func (x *errFileTooLarge) Error() string {
	return string(x.key)
}
func (x *errTooManyRequest) Error() string {
	return string(x.key)
}
func (x *errInternalError) Error() string {
	return string(x.key)
}

func (x *errLocked) Error() string {
	return string(x.key)
}

//functions
func NewHttpErrorNotFound(text NotFoundKey) error {
	return &errNotFound{text}
}
func NewHttpErrorForbidden(text ForbidenKey) error {
	return &errForbidden{text}
}
func NewHttpErrorConflict(text ConflictKey) error {
	return &errConflict{text}
}
func NewHttpErrorNotAllowed(text NotAllowedKey) error {
	return &errNotAllowed{text}
}
func NewHttpErrorBadRequest(text BadRequestKey) error {
	return &errBadRequest{text}
}
func NewHttpErrorUnauthorized(text UnauthorizedKey) error {
	return &errUnauthorized{text}
}
func NewHttpErrorFileTooLarge(text FileTooLargeKey) error {
	return &errFileTooLarge{text}
}
func NewHttpErrorTooManyRequest(text TooManyRequestKey) error {
	return &errTooManyRequest{text}
}
func NewHttpErrorInternalError(text InternalErrorKey) error {
	return &errInternalError{text}
}

func NewHttpErrorLocked(text LockedKey) error {
	return &errLocked{text}
}

func GetStatusCodeByError(err error, lang responses.Lang) (statusCode int, errMsg string, key string) {
	statusCode = ErrorCodeOK
	errMsg = ErrorMessageOK
	switch err.(type) {
	case *errNotFound:
		statusCode = ErrorCodeNotFound
		errMsg = ErrorMessageNotFound
		switch lang {
		case responses.Turkmen:
			key = ErrFileResponse[string(err.(*errNotFound).key)].TM
		case responses.Russian:
			key = ErrFileResponse[string(err.(*errNotFound).key)].RU
		default:
			key = ErrFileResponse[string(err.(*errNotFound).key)].EN

		}

	case *errForbidden:
		statusCode = ErrorCodeForbidden
		errMsg = ErrorMessageForbidden
		switch lang {
		case responses.Turkmen:
			key = ErrFileResponse[string(err.(*errForbidden).key)].TM
		case responses.Russian:
			key = ErrFileResponse[string(err.(*errForbidden).key)].RU
		default:
			key = ErrFileResponse[string(err.(*errForbidden).key)].EN
		}

	case *errConflict:
		statusCode = ErrorCodeConflict
		errMsg = ErrorMessageConflict
		switch lang {
		case responses.Turkmen:
			key = ErrFileResponse[string(err.(*errConflict).key)].TM
		case responses.Russian:
			key = ErrFileResponse[string(err.(*errConflict).key)].RU
		default:
			key = ErrFileResponse[string(err.(*errConflict).key)].EN
		}

	case *errNotAllowed:
		statusCode = ErrorMethodNotAllowed
		errMsg = ErrorMessageMethodNotAllowed
		switch lang {
		case responses.Turkmen:
			key = ErrFileResponse[string(err.(*errNotAllowed).key)].TM
		case responses.Russian:
			key = ErrFileResponse[string(err.(*errNotAllowed).key)].RU
		default:
			key = ErrFileResponse[string(err.(*errNotAllowed).key)].EN
		}

	case *errBadRequest:
		statusCode = ErrorCodeBadRequest
		errMsg = ErrorMessageBadRequest
		switch lang {
		case responses.Turkmen:
			key = ErrFileResponse[string(err.(*errBadRequest).key)].TM
		case responses.Russian:
			key = ErrFileResponse[string(err.(*errBadRequest).key)].RU
		default:
			key = ErrFileResponse[string(err.(*errBadRequest).key)].EN
		}

	case *errUnauthorized:
		statusCode = ErrorCodeUnauthorized
		errMsg = ErrorMessageUnauthorized
		switch lang {
		case responses.Turkmen:
			key = ErrFileResponse[string(err.(*errUnauthorized).key)].TM
		case responses.Russian:
			key = ErrFileResponse[string(err.(*errUnauthorized).key)].RU
		default:
			key = ErrFileResponse[string(err.(*errUnauthorized).key)].EN
		}

	case *errFileTooLarge:
		statusCode = ErrorCodeFileSizeTooLarge
		errMsg = ErrorMessageFileSizeTooLarge
		switch lang {
		case responses.Turkmen:
			key = ErrFileResponse[string(err.(*errFileTooLarge).key)].TM
		case responses.Russian:
			key = ErrFileResponse[string(err.(*errFileTooLarge).key)].RU
		default:
			key = ErrFileResponse[string(err.(*errFileTooLarge).key)].EN
		}

	case *errTooManyRequest:
		statusCode = ErrorCodeTooManyRequests
		errMsg = ErrorMessageTooManyRequests
		switch lang {
		case responses.Turkmen:
			key = ErrFileResponse[string(err.(*errTooManyRequest).key)].TM
		case responses.Russian:
			key = ErrFileResponse[string(err.(*errTooManyRequest).key)].RU
		default:
			key = ErrFileResponse[string(err.(*errTooManyRequest).key)].EN
		}

	case *errInternalError:
		statusCode = ErrorCodeInternalServerError
		errMsg = ErrorMessageInternalServerError
		switch lang {
		case responses.Turkmen:
			key = ErrFileResponse[string(err.(*errInternalError).key)].TM
		case responses.Russian:
			key = ErrFileResponse[string(err.(*errInternalError).key)].RU
		default:
			key = ErrFileResponse[string(err.(*errInternalError).key)].EN
		}

	case *errLocked:
		statusCode = ErrorCodeLocked
		errMsg = ErrorMessageLocked
		switch lang {
		case responses.Turkmen:
			key = ErrFileResponse[string(err.(*errLocked).key)].TM
		case responses.Russian:
			key = ErrFileResponse[string(err.(*errLocked).key)].RU
		default:
			key = ErrFileResponse[string(err.(*errLocked).key)].EN
		}

	}

	return
}

func SendResponse(w http.ResponseWriter, err error, data interface{}, clog *log.Entry, lang responses.Lang) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	statusCode, _, key := GetStatusCodeByError(err, lang)
	w.WriteHeader(statusCode)

	var resp responses.GeneralResponse
	if statusCode == ErrorCodeOK {
		resp.Success = true
		resp.Data = data
	} else {
		resp.Success = false
		resp.ErrMsg = key
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		clog.WithError(err).Error(fmt.Sprint(" data: ", resp))
	}
}

func createDefaultErrMessage(source string) (err error) {
	c := map[string]string{
		"ERR_NA": "method not allowed",
		"ERR_BR": "Forma nädogry!",
		"ERR_UA": "Ulgama birikdirilmediňiz!",
		"ERR_FL": "Faýlyň göwrümi uly!",
		"ERR_MR": "too many request",
		"ERR_IE": "Sistemanyň näsazlygy!",

		"ERR_NF_USER": "Ulanyjy tapylmady!",

		"ERR_UNIQUE_USER": "Ulanyjy ady eýýäm ulanyşda!",

		"ERR_FB_USER":         "Ulanyjy bu hereket gadagan!",
		"ERR_FB_delete_USER":  "Ulanyja özüni pozmak bolmaýar!",
		"ERR_FB_ADMIN":        "Admin öz rolyny üýtgedip bilmeýär!",
		"ERR_FB_pwrd_USER":    "Açar sözi gabat gelmeýär!",
		"ERR_FB_ownpwrd_USER": "Öz açar sözüňi bu ýerde üýtgetmek bolmaýar!",
		"ERR_FB_blocked_USER": "Ulanyjy bloklanan!",
	}

	b, err := json.MarshalIndent(c, "", "    ")

	if err != nil {
		eMsg := "error marshall err-message file"
		log.WithError(err).Error(eMsg)
		err = errors.Wrap(err, eMsg)
		return
	}

	err = ioutil.WriteFile(source, b, 0644)
	if err != nil {
		eMsg := "error creating err-message file"
		log.WithError(err).Error(eMsg)
		err = errors.Wrap(err, eMsg)
		return
	}

	return
}
