package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/SawitProRecruitment/UserService/cons"
	"github.com/SawitProRecruitment/UserService/core"
	"github.com/SawitProRecruitment/UserService/core/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

var _ core.UserRepo = (*Repository)(nil)

type Repository struct {
	sawitDB *sqlx.DB
}

type NewRepoOptions struct {
	Username    string
	Password    string
	Host        string
	Port        int
	Database    string
	SSLMode     bool
	MaxIdleConn int
	MaxOpenConn int
	MaxIdleTime time.Duration
}

func New(ctx context.Context, opts NewRepoOptions) (*Repository, error) {
	sslModeStr := "disable"
	if opts.SSLMode {
		sslModeStr = "enable"
	}

	addr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		opts.Username,
		opts.Password,
		opts.Host,
		opts.Port,
		opts.Database,
		sslModeStr,
	)

	var pg, err = sqlx.ConnectContext(ctx, "postgres", addr)
	if err != nil {
		return nil, err
	}
	pg.SetMaxIdleConns(opts.MaxIdleConn)
	pg.SetMaxOpenConns(opts.MaxOpenConn)
	pg.SetConnMaxIdleTime(opts.MaxIdleTime)

	return &Repository{
		sawitDB: pg,
	}, nil
}

func (r *Repository) CreateUser(data *domain.User) (*domain.User, error) {
	q := `
		INSERT INTO users (full_name, phone_number, password) 
		VALUES (:full_name, :phone_number, crypt(:password, gen_salt('bf')))
		RETURNING id, full_name, phone_number;
	`

	arg := User{
		FullName:    data.FullName,
		PhoneNumber: data.PhoneNumber,
		Password:    data.Password,
	}

	nstmt, err := r.sawitDB.PrepareNamed(q)
	if err != nil {
		return nil, err
	}

	rows, err := nstmt.Queryx(&arg)
	if err != nil {
		return nil, err
	}

	u := User{}
	for rows.Next() {
		err = rows.StructScan(&u)
		if err != nil {
			return nil, err
		}
	}

	return &domain.User{
		ID:          u.ID.String(),
		FullName:    u.FullName,
		PhoneNumber: u.PhoneNumber,
	}, nil
}

func (r *Repository) GetUserByID(id string) (*domain.User, error) {
	q := `
		SELECT id, full_name, phone_number, created_at, updated_at 
		FROM users 
		WHERE id = :id AND is_active = TRUE;
	`

	validID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	arg := User{ID: validID}
	resQuery := User{}

	nstmt, err := r.sawitDB.PrepareNamed(q)
	if err != nil {
		return nil, err
	}

	rows, err := nstmt.Queryx(&arg)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.StructScan(&resQuery)
		if err != nil {
			return nil, err
		}
	}

	res := &domain.User{
		ID:          resQuery.ID.String(),
		FullName:    resQuery.FullName,
		PhoneNumber: resQuery.PhoneNumber,
	}

	return res, nil
}

func (r *Repository) PatchUserByID(id string, data *domain.User) (*domain.User, error) {
	qt := `
		UPDATE users SET %s 
		WHERE id = :id
		AND is_active = true
		RETURNING id, full_name, phone_number;
	`

	validID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	params := []string{}
	if data.FullName != "" {
		params = append(params, "full_name= :full_name")
	}

	if data.PhoneNumber != "" {
		params = append(params, "phone_number= :phone_number")
	}

	paramsStr := strings.Join(params, ",")
	q := fmt.Sprintf(qt, paramsStr)

	queryArg := UserPatchArg{
		ID:          validID,
		FullName:    data.FullName,
		PhoneNumber: data.PhoneNumber,
	}
	rows, err := r.sawitDB.NamedQuery(q, &queryArg)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, cons.ErrDataConflict
		}

		return nil, err
	}

	u := User{}
	for rows.Next() {
		err = rows.StructScan(&u)
		if err != nil {
			return nil, err
		}
	}

	res := &domain.User{
		ID:          u.ID.String(),
		FullName:    u.FullName,
		PhoneNumber: u.PhoneNumber,
	}

	return res, nil
}

func (r *Repository) Login(phone string, pass string) (*domain.User, error) {
	q := `
		SELECT id, full_name, phone_number, login_count
		FROM users
		WHERE phone_number = :phone_number
		AND password = crypt(:password, password)
		AND is_active = true;
	`

	arg := &UserLoginArg{
		PhoneNumber: phone,
		Password:    pass,
	}

	nstmt, err := r.sawitDB.PrepareNamed(q)
	if err != nil {
		return nil, err
	}

	rows, err := nstmt.Queryx(&arg)
	if err != nil {
		return nil, err
	}

	u := User{}
	for rows.Next() {
		err = rows.StructScan(&u)
		if err != nil {
			return nil, err
		}
	}

	if u.ID == uuid.Nil {
		return nil, errors.New("phone and password do not match")
	}

	q = `
		UPDATE users SET login_count = login_count + 1
		WHERE id = :id;
	`

	nstmt, err = r.sawitDB.PrepareNamed(q)
	if err != nil {
		return nil, err
	}

	rows, err = nstmt.Queryx(&u)

	if err != nil {
		return nil, err
	}

	res := &domain.User{
		ID:          u.ID.String(),
		FullName:    u.FullName,
		PhoneNumber: u.PhoneNumber,
		LoginCount:  u.LoginCount,
	}
	return res, nil
}
