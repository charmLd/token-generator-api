package usecases

import (
	"context"

	"fmt"

	log "github.com/sirupsen/logrus"

	"time"

	"github.com/google/uuid"

	"golang.org/x/crypto/bcrypt"

	"github.com/charmLd/token-generator-api/domain/boundary/adapters"
	"github.com/charmLd/token-generator-api/domain/boundary/repositories"
	"github.com/charmLd/token-generator-api/domain/entities"
	appErr "github.com/charmLd/token-generator-api/domain/error"
)

// Password hashes support in Auth service
const (
	Bcrypt = "bcrypt"
	Md5    = "md5"
)

// AuthUseCase implements the base type for usecase
type AuthUseCase struct {
	Config AuthConfig

	UserRepository repositories.UserRepositoryInterface

	TokenRepository repositories.TokenRepositoryInterface

	TokenAdapter adapters.TokenAdapterInterface
}

type AuthConfig struct {
	LoginTokenExpiry     int
	GeneratedTokenExpiry int
	Location             *time.Location
}

type EmailLoginRequest struct {
	Email    string
	Password string
}

// Authenticate authenticates a user
func (au *AuthUseCase) AuthenticateByEmail(ctx context.Context, loginReq EmailLoginRequest) (tokenRes entities.Token, err error) {
	//Testing
	//admin user
	//passHas,_:=au.generateHashPassword(ctx,"Test@123AbC")
	//fmt.Println("Pass:",passHas)

	//passHas, _ := au.generateHashPassword(ctx, "Abc@123cAb")
	//fmt.Println("Pass:", passHas)
	// Fetch user information
	user, err := au.UserRepository.GetUserByEmail(ctx, loginReq.Email)

	if err != nil {
		log.Error(ctx, "user db query failed", fmt.Sprintf("user authentication failed for email %v", loginReq.Email))
		if err.Error() == appErr.InvalidUserStr {
			err = au.throwUserNotExistError(ctx)
			return
		}
		return
	}

	if user.UserID == 0 {
		log.Error(ctx, "user does not exist", fmt.Sprintf("user authentication failed for %v", user.UserID))
		err = au.throwUserNotExistError(ctx)
		return
	}
	// password hasher : Bcrypt
	//salt is a plain test
	fmt.Println("user ready")
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(loginReq.Password+user.Salt))

	if err != nil {

		log.Error(ctx, err, fmt.Sprintf("user authentication failed for %v", user.UserID))
		// Incorrect password
		err = au.throwPasswordError(ctx)
		return
	}
	fmt.Println("pass correct")

	// Set JWT claims
	jwt := entities.JWTClaims{}
	jwt.UserID = user.UserID
	jwt.TokenID = uuid.New().String()
	jwt.UserRole = user.Role

	//Updating last login - should be non-blocking
	go func() {
		err = au.UserRepository.UpdateLastLogin(ctx, fmt.Sprint(user.UserID))
		if err != nil {
			log.Error(ctx, "updating last login failed", user.UserID)
		}
	}()

	// Generate JWT (logintoken)
	authToken, err := au.TokenAdapter.GenerateLoginToken(ctx, jwt)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("user authentication failed for %v", user.UserID))
		err = au.throwTokenIssueError(ctx)
		return
	}

	// Generate token entity
	tokenRes.Token = jwt.TokenID

	tokenRes.GeneratedToken = authToken

	tokenRes.Expiry = time.Now().Local().Add(time.Second * time.Duration(au.Config.LoginTokenExpiry))
	tokenRes.UserId = fmt.Sprint(user.UserID)

	// Persist token
	err = au.TokenRepository.CreateNewToken(ctx, tokenRes)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("user authentication failed for %v", user.UserID))
		return
	}

	log.Trace(ctx, fmt.Sprintf("user authentication is successful for user_id %v and token_id %v", user.UserID, tokenRes.Token))
	return

}

/*
//generateHashPassword create user passwords
func (*AuthUseCase) generateHashPassword(ctx context.Context, pass string) (passHash string, err error) {
	// Transform password to hash
	passHashBytes, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		log.Error(ctx, "", err, "generating hash password failed")
		return
	}
	passHash = string(passHashBytes)
	return
}
*/
