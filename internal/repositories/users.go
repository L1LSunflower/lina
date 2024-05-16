package repositories

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"

	"github.com/L1LSunflower/lina/internal/entities"
	"github.com/L1LSunflower/lina/internal/tools"
)

type UsersRepository struct{}

func NewUsersRepository() Users {
	return &UsersRepository{}
}

func (u *UsersRepository) Add(ctx context.Context, user *entities.User) error {
	db, err := tools.DbFromCtx(ctx)
	if err != nil {
		return err
	}
	userArgs := pgx.NamedArgs{
		"id":         user.Id,
		"username":   user.Username,
		"first_name": user.FirstName,
	}
	fieldsInLine, namedArgs := FieldsAndArgsFromMap(userArgs, delimiter)
	dbCtx, cancelFunc := tools.CtxWithTimeout(context.Background(), defaultTimeout)
	defer cancelFunc()
	query := "insert into users (" + fieldsInLine + ") values (" + namedArgs + ")"
	if _, err = db.Pool.Exec(dbCtx, query, userArgs); err != nil {
		return fmt.Errorf("unable to insert row: %w", err)
	}
	return nil
}

func (u *UsersRepository) GetAll(ctx context.Context) ([]*entities.User, error) {
	db, err := tools.DbFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	query := "select id, username, first_name from users where blocked is false"
	dbCtx, cancelFunc := tools.CtxWithTimeout(context.Background(), defaultTimeout)
	defer cancelFunc()
	rows, err := db.Pool.Query(dbCtx, query)
	if err != nil {
		return nil, err
	}
	var users []*entities.User
	for rows.Next() {
		user := &entities.User{}
		if err = rows.Scan(&user.Id, &user.Username, &user.FirstName, &user.LastName); err != nil {
			log.Println("error scanning row:", err)
			continue
		}
		users = append(users, user)
	}
	if len(users) == 0 {
		return nil, fmt.Errorf("no users found")
	}
	return users, nil
}

func (u *UsersRepository) GetById(ctx context.Context, id int64) (*entities.User, error) {
	db, err := tools.DbFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	query := fmt.Sprintf("select id, username, first_name from users where id=@id", id)
	dbCtx, cancelFunc := tools.CtxWithTimeout(context.Background(), defaultTimeout)
	defer cancelFunc()
	user := &entities.User{}
	if err = db.Pool.QueryRow(dbCtx, query, pgx.NamedArgs{"id": id}).Scan(
		&user.Id,
		&user.Username,
		&user.FirstName,
	); err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}
	return user, nil
}

func (u *UsersRepository) Block(ctx context.Context, id string) error {
	db, err := tools.DbFromCtx(ctx)
	if err != nil {
		return err
	}
	query := "update users set blocked = true where id=@id"
	dbCtx, cancelFunc := tools.CtxWithTimeout(context.Background(), defaultTimeout)
	defer cancelFunc()
	_, err = db.Pool.Exec(dbCtx, query, pgx.NamedArgs{"id": id})
	return err
}

func (u *UsersRepository) Delete(ctx context.Context, id string) error {
	db, err := tools.DbFromCtx(ctx)
	if err != nil {
		return err
	}
	query := "delete from users where id=@id"
	dbCtx, cancelFunc := tools.CtxWithTimeout(context.Background(), defaultTimeout)
	defer cancelFunc()
	_, err = db.Pool.Exec(dbCtx, query, pgx.NamedArgs{"id": id})
	return err
}
