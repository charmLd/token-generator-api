package repositories

import (
	"context"

	"fmt"

	"time"

	"github.com/charmLd/token-generator-api/domain/boundary/adapters"
	"github.com/charmLd/token-generator-api/domain/entities"
	log "github.com/sirupsen/logrus"
)

type UserRepository struct {
	DBAdapter   adapters.DBAdapterInterface
	transaction adapters.TransactionInterface
	Location    *time.Location
}

func NewUserRepository(dbAdapter adapters.DBAdapterInterface, loc *time.Location) *UserRepository {

	return &UserRepository{
		DBAdapter: dbAdapter,
		Location:  loc,
	}
}

func (usr *UserRepository) UpdateLastLogin(ctx context.Context, userID string) (err error) {

	query := `UPDATE users SET last_login = ? WHERE users.user_id = ?;`

	statement, err := usr.DBAdapter.Prepare(ctx, query)
	if statement != nil {
		defer statement.Close()
	}
	//defer statement.Close()
	if err != nil {

		log.Error(ctx, err, "preparing last login query failed")
		return
	}
	fmt.Println("login valuee: ", time.Now().In(usr.Location).Format("2006-01-02 15:04:05"))
	result, err := statement.Exec(time.Now().In(usr.Location).Format("2006-01-02 15:04:05"), userID)
	if err != nil {

		log.Error(ctx, err, "executing last login query failed")
		return
	}
	_, err = result.RowsAffected()
	if err != nil {

		log.Error(ctx, err, "reading last login query result failed")
		errRollback := usr.transaction.Rollback()
		if errRollback != nil {
			log.Error(ctx, err, fmt.Sprintf("before rollback, authToken "))
			return
		}

		return
	}

	return
}

func (usr *UserRepository) GetLastLoginTime(ctx context.Context, userID string) (lastlogin time.Time, err error) {

	query := `SELECT last_login from users WHERE user_id = ?;`

	statement, err := usr.DBAdapter.Prepare(ctx, query)
	if statement != nil {
		defer statement.Close()
	}
	if err != nil {

		log.Error(ctx, "", err, "preparing query  failed")
		return
	}

	rows, err := statement.Query(userID)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {

		log.Error(ctx, "", err, "query statement failed ")
		return
	}

	for rows.Next() {

		err = rows.Scan(
			&lastlogin,
		)

		if err != nil {

			log.Error(ctx, "", err, "reading results failed ")
			return
		}

	}

	log.Trace(ctx, "", "getting last login is successful")

	return
}

func (usr *UserRepository) GetUserByEmail(ctx context.Context, email string) (user entities.User, err error) {

	userQuery := `
		SELECT 
			users.user_id, users.hashed_pass,users.salt,users.role
			
		FROM users
		
		WHERE 
			users.email = ? 
			
	`

	statement, err := usr.DBAdapter.Prepare(ctx, userQuery)
	if statement != nil {
		defer statement.Close()
	}
	if err != nil {

		log.Error(ctx, err, fmt.Sprintf("preparing query to fetch user by email failed for %v", email))
		return
	}

	rows, err := statement.Query(email)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {

		log.Error(ctx, err, fmt.Sprintf("executing get user by email query failed for %v", email))
		return
	}

	for rows.Next() {

		err = rows.Scan(
			&user.UserID,
			&user.HashedPassword,
			&user.Salt,
			&user.Role,
		)

		if err != nil {

			log.Error(ctx, err, fmt.Sprintf("reading results failed for user %v", email))
			return
		}
	}

	return user, nil
}
