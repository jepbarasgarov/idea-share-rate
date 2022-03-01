package helpers

import (
	"belli/onki-game-ideas-mongo-backend/errs"
	"belli/onki-game-ideas-mongo-backend/models"
	"belli/onki-game-ideas-mongo-backend/responses"
	"context"
	"errors"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"math/rand"
	"mime"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
	_ "golang.org/x/image/tiff"
)

const (
	layoutISO  = "2006-01-02"
	layoutISO1 = "2006-01-02 15:04:05"
	chars      = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func IsPasswordValid(pwd string) bool {
	var (
		hasMinLen = false
		hasNumber = false
	)

	if len(pwd) > 5 {
		hasMinLen = true
	}

	for _, char := range pwd {
		switch {
		case unicode.IsNumber(char):
			hasNumber = true
		}
	}

	return hasMinLen && hasNumber
}

func IsEmailValid(email string) bool {

	var (
		hasMinlen         = false
		hasMail           = false
		hasDot            = false
		lessThanMaxlentgh = false
	)

	if len(email) > 5 {
		hasMinlen = true
	}

	if len(email) <= 64 {
		lessThanMaxlentgh = true
	}

	for _, char := range email {
		switch char {
		case '@':
			hasMail = true
		case '.':
			hasDot = true
		}
	}

	return hasMinlen && hasMail && hasDot && lessThanMaxlentgh
}

func IsPhoneValid(phone string) bool {
	k := len(phone)
	if k != 8 {
		return false
	}
	if phone[:1] != "6" {
		return false
	}

	_, err := strconv.Atoi(phone)

	return err == nil
}

func ConvertStringToUserRole(str string) (r responses.UserRole, err error) {
	roles := [2]responses.UserRole{
		responses.UserRoleAdmin,
		responses.UserRoleUser,
	}

	str = strings.ToUpper(str)

	for _, role := range roles {
		if string(role) == str {
			r = role
			return
		}
	}

	err = errors.New("invalid role")

	return
}

func ConvertStringToIdeaCondition(str string) (r responses.IdeaCondition, err error) {
	conds := [2]responses.IdeaCondition{
		responses.RatedIdea,
		responses.NotRatedIdea,
	}

	str = strings.ToUpper(str)

	for _, cond := range conds {
		if string(cond) == str {
			r = cond
			return
		}
	}

	err = errors.New("invalid condition")

	return
}

func ConvertStringToUserStatus(str string) (status responses.UserStatus, err error) {

	statusArray := [2]responses.UserStatus{
		responses.Active,
		responses.Blocked,
	}

	str = strings.ToUpper(str)

	for _, sts := range statusArray {
		if string(sts) == str {
			status = sts
			return
		}
	}

	err = errors.New("invalid status")

	return
}

func ConvertStringToUUID(str string) (id uuid.UUID, err error) {
	id, err = uuid.FromString(str)
	return
}

func GetExtFromFileName(str string) (ext string, err error) {

	l := len(str)

	for i := l - 1; i >= 0; i-- {
		if str[i] == '.' {
			ext = str[i+1 : l]
			break
		}
	}

	if len(ext) == 0 {
		err = errors.New("not extension")
	}

	return
}

// TODO su funcksiya .tiff,   .xlsx, .docx  okap bilenok tazeden duzet
func GetContentType(filename string) string {
	ext := filepath.Ext(filename)
	contentType := mime.TypeByExtension(ext)
	return contentType
}

func CheckFileExt(ext string, extensions []string) (err error) {
	for _, extension := range extensions {
		if extension == strings.ToLower(ext) {
			return nil
		}
	}

	err = errors.New("invalid file ext")

	return
}

func CheckFileContentType(contentType string, contentTypes []string) bool {
	for _, cType := range contentTypes {
		if cType == contentType {
			return true
		}
	}
	return false
}

func GenerateFileName() (name string, err error) {

	fileName, err := uuid.NewV4()
	if err != nil {
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	return fileName.String(), err
}

func GenerateSessionID() (sessionID string, err error) {

	fileName, err := uuid.NewV4()
	if err != nil {
		err = errs.NewHttpErrorInternalError(errs.ERR_IE)
		return
	}

	return fileName.String(), err
}

func VerifyMinLen(m map[string]string) (names string, err error) {

	for k, v := range m {
		if len(v) == 0 {
			err = errors.New("len == 0")
			names += k
			names += " "
		}
	}

	return
}

func InTimeSpan(start, end, check time.Time) bool {
	return check.After(start) && check.Before(end)
}

func GetTimeNowString() string {
	return time.Now().Format("2006-01-02 15:04:05.000000")
}

func ChangeStringToDate(x string) (t time.Time, err error) {
	t, err = time.Parse(layoutISO, x)
	return
}

func ChangeStringToDateWithHour(x string) (t time.Time, err error) {
	t, err = time.Parse(layoutISO1, x)
	return
}

func DiffrenceBetweenDates(x time.Time) int64 {
	date := x
	CurrentTime := time.Now()
	StringOfDate := CurrentTime.Format("2006-01-02 15:04:05")
	t, _ := time.Parse(layoutISO1, StringOfDate)
	k := t.Unix() - date.Unix()
	return k
}

func DiffrenceBetweenTwoDates(first, second time.Time) int64 {
	date := first
	date2 := second
	k := date.Unix() - date2.Unix()
	return k
}

func RandInt(min int, max int) int {
	return min + rand.Intn(max-min+1)
}

func Min(a, b int64) int64 {
	if a > b {
		return b
	} else {
		return a
	}
}

func GetIP(r *http.Request) (string, error) {
	//Get IP from the X-REAL-IP header
	ip := r.Header.Get("X-REAL-IP")
	netIP := net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}

	//Get IP from X-FORWARDED-FOR header
	ips := r.Header.Get("X-FORWARDED-FOR")
	splitIps := strings.Split(ips, ",")
	for _, ip := range splitIps {
		netIP := net.ParseIP(ip)
		if netIP != nil {
			return ip, nil
		}
	}

	//Get IP from RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}
	netIP = net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}
	return "", fmt.Errorf("No valid ip found")
}

func GetPowerOfNumber(number int, exponent int) (result int64) {
	result = 1
	for i := 0; i < exponent; i++ {
		result = result * int64(number)
	}
	return
}

func GetRequestLang(r *http.Request) (lang responses.Lang) {
	l := r.Header.Get("X-Lang")
	switch l {
	case "tk":
		lang = responses.Turkmen
	case "ru":
		lang = responses.Russian
	default:
		lang = responses.English
	}
	return lang
}

func GenRandomString(length int) (string, error) {
	bytes := make([]byte, length)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	for i, b := range bytes {
		bytes[i] = chars[b%byte(len(chars))]
	}

	return string(bytes), nil
}

/// FILE related functions

func ProcessFile(ctx context.Context, staticDir string, pFile *models.ParsedFile) (string, error) {
	clog := log.WithFields(log.Fields{
		"method": "ProcessFile",
	})

	innerDir, filename, err := createDir(staticDir, "sketchs")
	if err != nil {
		clog.Error("an error occurred on createDir function")
		return "", err
	}

	path, err := saveFile(pFile, staticDir, filepath.Join("sketchs", innerDir), filename)
	if err != nil {
		clog.Error("an error occurred on saveFile function")
		return "", err
	}

	return path, nil
}

func saveFile(pFile *models.ParsedFile, staticDir, dir, filename string) (string, error) {
	clog := log.WithFields(log.Fields{
		"method": "saveFile",
	})

	ext := getExtByContentType(pFile.ContentType)
	path := filepath.Join(dir, fmt.Sprintf("%s%s", filename, ext))
	fullpath := filepath.Join(staticDir, path)

	f, err := os.OpenFile(fullpath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		clog.Error("an error occurred on os.OpenFile function")
		return "", err
	}
	defer f.Close()

	_, err = pFile.File.Seek(0, 0)
	if err != nil {
		clog.Error("an error occurred on pFile.File.Seek function")
		return "", err
	}

	fBytes, err := ioutil.ReadAll(pFile.File)
	if err != nil {
		clog.Error("an error occurred on ioutil.ReadAll function")
		return "", err
	}

	f.Write(fBytes)
	return path, nil
}

func createDir(staticDir, outerDir string) (string, string, error) {
	clog := log.WithFields(log.Fields{
		"method": "createDir",
	})

	filename, err := genNewFilename()
	if err != nil {
		clog.Error("an error occurred on genNewFilename function")
		return "", "", err
	}

	chunks := strings.Split(filename, "-")
	innerDir := filepath.Join(chunks[0][:2], chunks[1][:2], chunks[2][:2])

	fullpath := filepath.Join(staticDir, outerDir, innerDir)

	err = os.MkdirAll(fullpath, 0755)
	if err != nil {
		clog.Error("an error occurred on os.MkdirAll function")
		return "", "", err
	}

	return innerDir, filename, nil
}

func genNewFilename() (string, error) {
	clog := log.WithFields(log.Fields{
		"method": "genNewFilename",
	})
	filename, err := uuid.NewV4()
	if err != nil {
		clog.Error("an error occurred on uuid.NewV4 function")
		return "", err
	}
	return filename.String(), nil
}

func getExtByContentType(contentType string) string {
	// TODO: fix this sample of code
	switch contentType {
	case "image/jpeg":
		return ".jpeg"
	case "image/png":
		return ".png"
	case "image/tiff":
		return ".tiff"
	case "application/pdf":
		return ".pdf"
	case "application/msword":
		return ".doc"
	case "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
		return ".docx"
	case "application/vnd.ms-excel":
		return ".xls"
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		return ".xlsx"
	}

	return ""
}
