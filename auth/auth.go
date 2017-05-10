package auth

import (
	"github.com/tomogoma/authms/auth/oauth/response"
	"github.com/tomogoma/authms/proto/authms"
	"github.com/tomogoma/go-commons/errors"
	"regexp"
	"strings"
	"github.com/dgrijalva/jwt-go"
	"github.com/tomogoma/authms/claim"
	"time"
)

type OAuthHandler interface {
	ValidateToken(appName, token string) (response.OAuth, error)
}

type TokenGenerator interface {
	Generate(claims jwt.Claims) (string, error)
	Validate(token string, claims jwt.Claims) (*jwt.Token, error)
}

type PasswordGenerator interface {
	SecureRandomString(length int) ([]byte, error)
}

type Logger interface {
	Info(interface{}, ...interface{})
	Error(interface{}, ...interface{}) error
}

type DBHelper interface {
	UserExists(u *authms.User) (int64, error)
	SaveUser(*authms.User) error
	SaveHistory(*authms.History) error
	GetByUserName(uname, pass string) (*authms.User, error)
	GetByAppUserID(appName, appUserID string) (*authms.User, error)
	GetByPhone(phone, pass string) (*authms.User, error)
	GetByEmail(email, pass string) (*authms.User, error)
	GetHistory(userID int64, offset, count int, accessType ...string) ([]*authms.History, error)
	UpdatePhone(userID int64, newPhone *authms.Value) error
	UpdateAppUserID(userID int64, new *authms.OAuth) error
	IsNotFoundError(err error) bool
	IsDuplicateError(err error) bool
}

type PhoneVerifier interface {
	SendSMSCode(toPhone string) (*authms.SMSVerificationStatus, error)
	VerifySMSCode(r *authms.SMSVerificationCodeRequest) (*authms.SMSVerificationStatus, error)
}

type Auth struct {
	dbHelper      DBHelper
	tokenG        TokenGenerator
	logger        Logger
	oAuthHandler  OAuthHandler
	phoneVerifier PhoneVerifier
	errors.ClErrCheck
	errors.AuthErrCheck
}

const (
	numPrevLogins        = 5
	AccessLogin          = "LOGIN"
	AccessRegistration   = "REGISTER"
	AccessUpdate         = "UPDATE"
	AccessVerification   = "VERIFICATION"
	AccessCodeValidation = "VERIFICATION_CODE_VALIDATION"
	numExp               = `[0-9]+`
	tokenValidity        = 8 * time.Hour
)

var ErrorNilTokenGenerator = errors.New("token generator was nil")
var ErrorNilLogger = errors.New("Logger was nil")
var ErrorNilDBHelper = errors.New("DBHelper was nil")
var ErrorNilOAuthHandler = errors.New("oauth handler was nil")
var ErrorNilPhoneVerifier = errors.New("PhoneVerifier was nil")
var rePhone = regexp.MustCompile(numExp)

func New(tg TokenGenerator, lg Logger, db DBHelper,
	oa OAuthHandler, pv PhoneVerifier) (*Auth, error) {
	if tg == nil {
		return nil, ErrorNilTokenGenerator
	}
	if lg == nil {
		return nil, ErrorNilLogger
	}
	if db == nil {
		return nil, ErrorNilDBHelper
	}
	if oa == nil {
		return nil, ErrorNilOAuthHandler
	}
	if pv == nil {
		return nil, ErrorNilPhoneVerifier
	}
	return &Auth{dbHelper: db, tokenG: tg, oAuthHandler: oa,
		logger:        lg, phoneVerifier: pv}, nil
}

func (a *Auth) Register(user *authms.User, devID, rIP string) error {
	defer func() {
		if user != nil {
			user.Password = ""
		}
	}()
	if err := validateUser(user); err != nil {
		return err
	}
	if devID == "" {
		return errors.NewClient("Dev ID was empty")
	}
	for appName, oa := range user.OAuths {
		if appName != oa.AppName {
			return errors.NewClient("an OAuth's key does not match its app name")
		}
		if err := a.validateOAuth(oa); err != nil {
			return err
		}
		user.OAuths[appName].Verified = true
	}
	if user.Phone != nil {
		user.Phone.Value = formatPhone(user.Phone.Value)
		user.Phone.Verified = false
	}
	if user.Email != nil {
		user.Email.Verified = false
	}
	if err := a.checkUserExists(user); err != nil {
		return err
	}
	err := a.dbHelper.SaveUser(user)
	if err != nil {
		return errors.Newf("error persisting user: %v", err)
	}
	go a.saveHistory(user, devID, AccessRegistration, rIP, nil)
	return nil
}

func (a *Auth) UpdatePhone(user *authms.User, token, devID, rIP string) error {
	if user == nil {
		return errors.NewClient("user was empty")
	}
	clm := claim.Auth{}
	_, err := a.tokenG.Validate(token, &clm)
	defer func() {
		go a.saveHistory(user, devID, AccessUpdate, rIP, err)
	}()
	if err != nil || clm.UsrID != user.ID {
		return errors.NewAuthf("invalid token: %v", err)
	}
	if !hasValue(user.Phone) {
		return errors.NewClient("phone was invalid")
	}
	if devID == "" {
		return errors.NewClient("device ID was empty")
	}
	user.Phone.Value = formatPhone(user.Phone.Value)
	user.Phone.Verified = false
	if err := a.checkUserExists(user); err != nil {
		return err
	}
	err = a.dbHelper.UpdatePhone(user.ID, user.Phone)
	if err != nil {
		return errors.Newf("error persisting phone update: %v", err)
	}
	return nil
}

func (a *Auth) VerifyPhone(req *authms.SMSVerificationRequest, rIP string) (*authms.SMSVerificationStatus, error) {
	if req == nil {
		return nil, errors.NewClient("SMSVerificationRequest was empty")
	}
	clm := claim.Auth{}
	_, err := a.tokenG.Validate(req.Token, &clm)
	defer func() {
		go a.saveHistory(&authms.User{ID: req.UserID}, req.DeviceID,
			AccessVerification, rIP, err)
	}()
	if err != nil || clm.UsrID != req.UserID {
		return nil, errors.NewAuthf("invalid token: %v", err)
	}
	req.Phone = formatPhone(req.Phone)
	testExistsUsr := &authms.User{
		ID:    req.UserID,
		Phone: &authms.Value{Value: req.Phone},
	}
	if err := a.checkUserExists(testExistsUsr); err != nil {
		return nil, err
	}
	vs, err := a.phoneVerifier.SendSMSCode(req.Phone)
	if err != nil {
		return nil, err
	}
	return vs, nil
}

func (a *Auth) VerifyPhoneCode(req *authms.SMSVerificationCodeRequest, rIP string) (*authms.SMSVerificationStatus, error) {
	if req == nil {
		return nil, errors.NewClient("SMSVerificationRequest was empty")
	}
	clm := claim.Auth{}
	_, err := a.tokenG.Validate(req.Token, &clm)
	defer func() {
		go a.saveHistory(&authms.User{ID: req.UserID}, req.DeviceID,
			AccessCodeValidation, rIP, err)
	}()
	if err != nil || clm.UsrID != req.UserID {
		return nil, errors.NewAuthf("invalid token: %v", err)
	}
	vs, err := a.phoneVerifier.VerifySMSCode(req)
	if err != nil {
		return nil, errors.NewClientf("%v", err)
	}
	phone := &authms.Value{
		Verified: vs.Verified,
		Value:    formatPhone(vs.Phone),
	}
	testExistsUsr := &authms.User{
		ID:    req.UserID,
		Phone: phone,
	}
	if err := a.checkUserExists(testExistsUsr); err != nil {
		return nil, err
	}
	err = a.dbHelper.UpdatePhone(req.UserID, phone)
	if err != nil {
		return nil, errors.Newf("error persisting phone update: %v", err)
	}
	return vs, nil
}

func (a *Auth) UpdateOAuth(user *authms.User, appName, token, devID, rIP string) error {
	if user == nil {
		return errors.NewClient("user was empty")
	}
	clm := claim.Auth{}
	_, err := a.tokenG.Validate(token, &clm)
	defer func() {
		go a.saveHistory(user, devID, AccessUpdate, rIP, err)
	}()
	if err != nil || clm.UsrID != user.ID {
		return errors.NewAuthf("invalid token: %v", err)
	}
	if user.OAuths == nil || user.OAuths[appName] == nil {
		return errors.NewClient("OAuth was not provided")
	}
	oa := user.OAuths[appName]
	if err := a.validateOAuth(oa); err != nil {
		return err
	}
	if devID == "" {
		return errors.NewClient("device ID was empty")
	}
	oa.Verified = true
	if err := a.checkUserExists(user); err != nil {
		return err
	}
	err = a.dbHelper.UpdateAppUserID(user.ID, oa)
	if err != nil {
		return errors.Newf("error persisting OAuth changes: %v", err)
	}
	return nil
}

func (a *Auth) UpdatePassword(userID int64, oldPass, newPass, devID, rIP string) error {
	return errors.New("not yet implemented")
}

func (a *Auth) LoginUserName(uName, pass, devID, rIP string) (*authms.User, error) {
	if devID == "" {
		return nil, errors.NewClient("Dev ID was empty")
	}
	usr, err := a.dbHelper.GetByUserName(uName, pass)
	if err = a.processLoginResults(usr, devID, rIP, err); err != nil {
		return nil, err
	}
	return usr, nil
}

func (a *Auth) LoginPhone(phone, pass, devID, rIP string) (*authms.User, error) {
	if devID == "" {
		return nil, errors.NewClient("Dev ID was empty")
	}
	phone = formatPhone(phone)
	usr, err := a.dbHelper.GetByPhone(phone, pass)
	if err = a.processLoginResults(usr, devID, rIP, err); err != nil {
		return nil, err
	}
	return usr, nil
}

func (a *Auth) LoginEmail(email, pass, devID, rIP string) (*authms.User, error) {
	if devID == "" {
		return nil, errors.NewClient("Dev ID was empty")
	}
	usr, err := a.dbHelper.GetByEmail(email, pass)
	if err = a.processLoginResults(usr, devID, rIP, err); err != nil {
		return nil, err
	}
	return usr, nil
}

func (a *Auth) LoginOAuth(app *authms.OAuth, devID, rIP string) (*authms.User, error) {
	if devID == "" {
		return nil, errors.NewClient("Dev ID was empty")
	}
	if err := a.validateOAuth(app); err != nil {
		return nil, err
	}
	usr, err := a.dbHelper.GetByAppUserID(app.AppName, app.AppUserID)
	if err = a.processLoginResults(usr, devID, rIP, err); err != nil {
		return nil, err
	}
	return usr, nil
}

func (a *Auth) checkUserExists(user *authms.User) error {
	existUsrID, err := a.dbHelper.UserExists(user)
	if err != nil {
		return errors.Newf("Error checking if user exists")
	}
	if existUsrID >= 1 && user.ID != existUsrID {
		return errors.NewClient("A user with some of the provided" +
			" credentials already exists")
	}
	return nil
}

func (a *Auth) processLoginResults(usr *authms.User, devID, rIP string, loginErr error) error {
	if a.dbHelper.IsNotFoundError(loginErr) {
		loginErr = errors.NewAuth("invalid credentials")
	}
	defer func() {
		go a.saveHistory(usr, devID, AccessLogin, rIP, loginErr)
	}()
	if loginErr != nil {
		return loginErr
	}
	clm := claim.NewAuth(usr.ID, devID, tokenValidity)
	tkn, loginErr := a.tokenG.Generate(clm)
	if loginErr != nil {
		loginErr = errors.Newf("error generating token: %v", loginErr)
		return loginErr
	}
	usr.Token = tkn
	prevLogins, loginErr := a.dbHelper.GetHistory(usr.ID, 0, numPrevLogins,
		AccessLogin)
	if loginErr != nil {
		loginErr = errors.Newf("error fetching login history: %v", loginErr)
		return loginErr
	}
	usr.LoginHistory = prevLogins
	return nil
}

func (a *Auth) validateOAuth(claimOA *authms.OAuth) error {
	if claimOA == nil {
		return nil
	}
	oa, err := a.oAuthHandler.ValidateToken(claimOA.AppName, claimOA.AppToken)
	if err != nil {
		return errors.Newf("error validating OAuth token: %v", err)
	}
	if !oa.IsValid() || oa.UserID() != claimOA.AppUserID {
		return errors.NewAuth("OAuth token was not valid")
	}
	return nil
}

func (a *Auth) saveHistory(user *authms.User, devID, accType, rIP string, err error) {
	if user == nil || user.ID < 1 {
		return
	}
	accSuccessful := true
	if err != nil {
		accSuccessful = false
	}
	h := &authms.History{UserID: user.ID, AccessType: accType,
		SuccessStatus:       accSuccessful, IpAddress: rIP, DevID: devID}
	err = a.dbHelper.SaveHistory(h)
	if err != nil {
		a.logger.Error("unable to save auth history entry (' %+v '): %s", h, err)
	}
}

func formatPhone(phone string) string {
	parts := rePhone.FindAllString(phone, -1)
	formatted := ""
	if strings.HasPrefix(phone, "+") {
		formatted = "+"
	}
	for _, part := range parts {
		formatted = formatted + part
	}
	return formatted
}

func havePasswordComboAuth(u *authms.User) bool {
	return hasValue(u.Phone) || hasValue(u.Email) || u.UserName != ""
}

func validateUser(u *authms.User) error {
	if u == nil {
		return errors.NewClient("user was not provided")
	}
	hasPhone := hasValue(u.Phone)
	hasMail := hasValue(u.Email)
	if u.UserName == "" && !hasPhone && !hasMail && u.OAuths == nil {
		return errors.NewClient("A user must have at least one" +
			" identifier (UserName, Phone, Email, OAuthApp")
	}
	if havePasswordComboAuth(u) && u.Password == "" {
		return errors.NewClient("password was not provided")
	}
	for name, oa := range u.OAuths {
		if name != oa.AppName {
			return errors.NewClient("an OAuth key did not match AppName")
		}
		if err := validateOAuth(oa); err != nil {
			return err
		}
	}
	return nil
}

func validateOAuth(oa *authms.OAuth) error {
	if oa == nil {
		return errors.NewClient("OAuth was not provided")
	}
	if oa.AppName == "" {
		return errors.NewClient("AppName was not provided")
	}
	if oa.AppToken == "" {
		return errors.NewClient("AppToken was not provided")
	}
	if oa.AppUserID == "" {
		return errors.NewClient("AppUserID was not provided")
	}
	return nil
}

func hasValue(v *authms.Value) bool {
	return v != nil && v.Value != ""
}
