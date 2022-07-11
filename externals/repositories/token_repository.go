package repositories

import (
	"context"
	"errors"
	"fmt"
	"github.com/charmLd/token-generator-api/domain/boundary/repositories"

	"time"

	"github.com/charmLd/token-generator-api/domain/boundary/adapters"
	"github.com/charmLd/token-generator-api/domain/entities"
	log "github.com/sirupsen/logrus"
)

type TokenRepository struct {
	DBAdapter   adapters.DBAdapterInterface
	transaction adapters.TransactionInterface
	Location    *time.Location
}

func NewTokenRepository(dbAdapter adapters.DBAdapterInterface, loc *time.Location) repositories.TokenRepositoryInterface {

	return &TokenRepository{
		DBAdapter: dbAdapter,
		Location:  loc,
	}
}

// Revoke invalidated a token
func (tr *TokenRepository) Revoke(ctx context.Context, userId string) (err error) {

	tr.transaction, err = tr.DBAdapter.BeginTransaction(ctx)
	if err != nil {
		log.Error(ctx, err, "begin transaction failed for revoke query")
		return
	}
	row, err := tr.revokeToken(ctx, userId)
	if err != nil {
		log.Error(ctx, "revoking query failed")
		return
	}
	err = tr.transaction.Commit()
	if err != nil {
		log.Error(ctx, err, "committing revoke query failed")
		return
	}
	log.Debug(ctx, fmt.Sprintf("user : %v, revoked, rows affected : %v", userId, row))
	return
}

func (tr *TokenRepository) InsertUniqueToken(ctx context.Context, token entities.Token) (err error) {
	tr.transaction, err = tr.DBAdapter.BeginTransaction(ctx)
	if err != nil {
		log.Error(ctx, err, "beginning transaction failed")
		return
	}

	err = tr.insertToken(ctx, token)
	if err != nil {
		log.Error(ctx, "creating new token failed")
		return
	}

	err = tr.transaction.Commit()
	if err != nil {
		log.Error(ctx, err, "committing create new token transaction failed")
		return
	}

	return
}

func (tr *TokenRepository) insertToken(ctx context.Context, token entities.Token) (err error) {

	query := `
		INSERT INTO generated_tokens
		(gen_token_id,user_id,token,created_at,expiry,is_blacklisted)
		VALUES(?,?, ?, ?, ?, ?);
	`

	statement, err := tr.transaction.Prepare(query)
	if statement != nil {
		defer statement.Close()
	}
	if err != nil {
		log.Error(ctx, err, "preparing add token with transaction failed")
		errRollback := tr.transaction.Rollback()
		if errRollback != nil {
			log.Error(ctx, err, fmt.Sprintf("before rollback, authToken "))
		}
		return
	}

	_, err = statement.Exec(
		token.Token,
		token.UserId,
		token.GeneratedToken,
		token.CreatedAt.In(tr.Location).Format("2006-01-02 15:04:05"),
		token.Expiry.In(tr.Location).Format("2006-01-02 15:04:05"),
		token.IsBlacklisted,
	)
	if err != nil {
		log.Error(ctx, err, "executing add token query failed")
		errRollback := tr.transaction.Rollback()
		if errRollback != nil {
			log.Error(ctx, err, fmt.Sprintf("before rollback, authToken "))
		}
		return
	}

	return
}

func (tr *TokenRepository) revokeToken(ctx context.Context, userId string) (tokensRevoked int, err error) {

	query := `
		UPDATE generated_tokens
		SET is_blacklisted=?
		WHERE user_id = ?;
	`

	statement, err := tr.transaction.Prepare(query)
	if statement != nil {
		defer statement.Close()
	}
	if err != nil {
		log.Error(ctx, err, "preparing revoke user token query failed")
		return
	}

	if userId == "" {
		log.Error(ctx, "no user id provided")
		return 0, errors.New("user ID is invalid :" + userId)
	}
	result, err := statement.Exec(1, userId)
	if err != nil {
		errRollback := tr.transaction.Rollback()
		if errRollback != nil {
			log.Error(ctx, err, fmt.Sprintf("before rollback, authToken "))
			return
		}
		log.Error(ctx, err, "executing revoke token query failed")
		return
	}
	rowsAff, err := result.RowsAffected()
	if err != nil {
		errRollback := tr.transaction.Rollback()
		if errRollback != nil {
			log.Error(ctx, err, fmt.Sprintf("before rollback, authToken "))
			return
		}
		log.Error(ctx, err, "reading revoke token query result failed")
		return
	}

	tokensRevoked = int(rowsAff)
	return
}

func (t *TokenRepository) GetAllTokenForFilter(ctx context.Context, fetchDetailsFilters entities.TokenDetailsReqParam) (tokenDetailsArray []entities.Token, err error) {

	var userQuery string

	if fetchDetailsFilters.Balcklisted.IsOK && fetchDetailsFilters.Balcklisted.Value != "" {

		userQuery = `
		SELECT gen_token_id,user_id,token,created_at,expiry,is_blacklisted
		FROM generated_tokens WHERE is_blacklisted=` + fetchDetailsFilters.Balcklisted.Value + ` AND user_id=
		 ` + fetchDetailsFilters.UserId + ` ORDER BY user_id DESC  ;
`
	} else {
		userQuery = `
		SELECT gen_token_id,user_id,token,created_at,expiry,is_blacklisted
		FROM generated_tokens WHERE user_id=
		 ` + fetchDetailsFilters.UserId + ` ORDER BY user_id DESC  ;
`
	}

	statement, err := t.DBAdapter.Prepare(ctx, userQuery)
	if statement != nil {
		defer statement.Close()
	}
	if err != nil {

		log.Error(ctx, "", err, fmt.Sprintf("preparing query to fetch app login failed"))
		return
	}

	rows, err := statement.Query()
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {

		log.Error(ctx, "", err, fmt.Sprintf("executing get app login query failed "))
		return
	}

	for rows.Next() {
		tokendetail := entities.Token{}

		err = rows.Scan(
			&tokendetail.Token,
			&tokendetail.UserId,
			&tokendetail.GeneratedToken,
			&tokendetail.CreatedAt,
			&tokendetail.Expiry,
			&tokendetail.IsBlacklisted,
		)

		if err != nil {

			log.Error(ctx, "", err, fmt.Sprintf("reading results failed "))
			return
		}
		tokenDetailsArray = append(tokenDetailsArray, tokendetail)
	}

	log.Trace(ctx, "", "GetAllTokenForFilter is successful")
	return tokenDetailsArray, nil
}
func (tr *TokenRepository) CreateNewToken(ctx context.Context, token entities.Token) (err error) {
	tr.transaction, err = tr.DBAdapter.BeginTransaction(ctx)
	if err != nil {
		log.Error(ctx, err, "beginning transaction failed")
		return
	}

	err = tr.addToken(ctx, token)
	if err != nil {
		log.Error(ctx, "creating new token failed")
		return
	}

	err = tr.transaction.Commit()
	if err != nil {
		log.Error(ctx, err, "committing create new token transaction failed")
		return
	}

	return
}
func (tr *TokenRepository) addToken(ctx context.Context, token entities.Token) (err error) {

	query := `
		INSERT INTO auth_tokens
		(token_id, user_id,  is_blacklisted, auth_token, expiry)
		VALUES(UUID_TO_BIN(?), ?, ?, ?, ?);
	`

	statement, err := tr.transaction.Prepare(query)
	if statement != nil {
		defer statement.Close()
	}
	if err != nil {
		log.Error(ctx, err, "preparing add token with transaction failed")
		errRollback := tr.transaction.Rollback()
		if errRollback != nil {
			log.Error(ctx, err, fmt.Sprintf("before rollback, authToken "))
		}
		return
	}

	_, err = statement.Exec(
		token.Token,
		token.UserId,

		token.IsBlacklisted,
		token.GeneratedToken,
		token.Expiry.In(tr.Location).Format("2006-01-02 15:04:05"),
	)
	if err != nil {
		log.Error(ctx, err, "executing add token query failed")
		errRollback := tr.transaction.Rollback()
		if errRollback != nil {
			log.Error(ctx, err, fmt.Sprintf("before rollback, authToken "))
		}
		return
	}

	return
}
func (t *TokenRepository) FetchTokenInfo(ctx context.Context, tokenDetails entities.ValidateRequest) (tokendetail entities.Token, err error) {

	var userQuery string

	userQuery = `
		SELECT user_id,token,created_at,expiry,is_blacklisted
		FROM generated_tokens WHERE user_id=? AND  gen_token_id =? ;
`

	statement, err := t.DBAdapter.Prepare(ctx, userQuery)
	if statement != nil {
		defer statement.Close()
	}
	if err != nil {

		log.Error(ctx, "", err, fmt.Sprintf("preparing query to fetch app login failed"))
		return
	}

	rows, err := statement.Query(tokenDetails.UserId, tokenDetails.InviteToken)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {

		log.Error(ctx, "", err, fmt.Sprintf("executing get app login query failed "))
		return
	}

	for rows.Next() {

		err = rows.Scan(
			&tokendetail.UserId,
			&tokendetail.GeneratedToken,
			&tokendetail.CreatedAt,
			&tokendetail.Expiry,
			&tokendetail.IsBlacklisted,
		)

		if err != nil {

			log.Error(ctx, "", err, fmt.Sprintf("reading results failed "))
			return
		}
		break
	}

	fmt.Println(ctx, " ", "fetch token details is successful")
	return tokendetail, nil
}
