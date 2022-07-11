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
}

type EmailLoginRequest struct {
	Email    string
	Password string
}

// Authenticate authenticates a user
func (au *AuthUseCase) AuthenticateByEmail(ctx context.Context, loginReq EmailLoginRequest) (tokenRes entities.Token, err error) {
	//Testing
	/*passHas,_:=au.generateHashPassword(ctx,"Test@123AbC")
	fmt.Println("Pass:",passHas)*/
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
	//todo if user login are exceeded its concurrent limit to access the web app it should remove  tokens as concurrent login increaese

	/*

		// If multiple login is not allowed for appId, revoke the other token and issue a new token.
		if app.ConcurrentLoginLimit == 1 {
			// Revoke all other tokens
			err = au.RevokeApp(ctx, user.UserID, app.AppID)
			if err != nil {
				log.Error(ctx, fmt.Sprintf("user authentication failed for %v", user.UserID))
				return
			}
			log.DebugContext(ctx, fmt.Sprintf("revoked the existing token as concurrent login limit is 1 for %v", loginReq.AppName))
		} else if app.ConcurrentLoginLimit > 1 {
			// Get tokens issued for the same app
			var tokens []entities.Token
			tokens, err = au.TokenRepository.GetIssuedTokens(ctx, user.UserID, app.AppID)
			if err != nil {
				log.Error(ctx, fmt.Sprintf("user authentication failed for %v", user.UserID))
				return
			}

			// If multiple login is allowed, but concurrent_login_limit has not been exceeded, issue a new token.
			if app.ConcurrentLoginLimit <= len(tokens) {
				//revoking the oldest token
				//no validation required for the array since len(tokens) >= concurrentLimit > 1
				oldestToken := tokens[0]
				for _, token := range tokens {
					// oldest token needs to be to_be_revoked=0 since unless the concurrent limit will be exceeded
					if token.CreatedAt.Unix() < oldestToken.CreatedAt.Unix() && !token.ToBeRevoked {
						oldestToken = token
					}
				}

				err = au.Revoke(ctx, oldestToken.ID, user.UserID)
				if err != nil {
					log.Error(ctx, fmt.Sprintf("revoking oldest token failed for %v", oldestToken.ID), fmt.Sprintf("user authentication failed for %v", user.UserID))
					return
				}
				log.DebugContext(ctx, fmt.Sprintf("oldest token which was created at %v was revoked for the user %v", oldestToken.CreatedAt, user.UserID))
			}
		} else {
			log.Error(ctx, fmt.Sprintf("concurrent login limit for app %v is 0 or invalid", app.AppID), fmt.Sprintf("user authentication failed for %v", user.UserID))
			err = au.throwConcurrentLoginLimitError(ctx)
			return
		}
	*/
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
	tokenRes.ID = jwt.TokenID

	tokenRes.GeneratedToken = authToken

	tokenRes.Expiry = time.Now().Local().Add(time.Second * time.Duration(au.Config.LoginTokenExpiry))
	tokenRes.UserId = fmt.Sprint(user.UserID)

	// Persist token
	err = au.TokenRepository.CreateNewToken(ctx, tokenRes)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("user authentication failed for %v", user.UserID))
		return
	}

	log.Trace(ctx, fmt.Sprintf("user authentication is successful for user_id %v and token_id %v", user.UserID, tokenRes.ID))
	return

}

func (au *AuthUseCase) throwPasswordError(ctx context.Context) error {
	return (&appErr.Error{}).New(ctx, appErr.DOMAIN, "incorrect password",
		"CORE-2004",
		"incorrect password")
}

func (au *AuthUseCase) throwUserNotExistError(ctx context.Context) error {
	return (&appErr.Error{}).New(ctx, appErr.DOMAIN, "invalid user",
		"CORE-2006",
		"invalid user")
}

func (au *AuthUseCase) throwTokenIssueError(ctx context.Context) error {
	return (&appErr.Error{}).New(ctx, appErr.DOMAIN, "token Issue Error",
		"CORE-2012",
		"could not issue a token")
}

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