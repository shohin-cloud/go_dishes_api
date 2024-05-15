package model

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"time"

	"github.com/shohin-cloud/dishes-api/pkg/dishes/validator"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

var AnonymousMember = &Member{}

type Member struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	Activated bool      `json:"activated"`
	Version   int       `json:"-"`
}

func (m *Member) IsAnonymous() bool {
	return m == AnonymousMember
}

type MemberModel struct {
	DB *sql.DB
}

type password struct {
	plaintext *string
	hash      []byte
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}
	p.plaintext = &plaintextPassword
	p.hash = hash
	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func ValidateMember(v *validator.Validator, member *Member) {
	v.Check(member.Name != "", "name", "must be provided")
	v.Check(len(member.Name) <= 500, "name", "must not be more than 500 bytes long")

	ValidateEmail(v, member.Email)

	if member.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *member.Password.plaintext)
	}

	if member.Password.hash == nil {
		panic("missing password hash for member")
	}
}

func (m MemberModel) Insert(member *Member) error {
	query := `
		INSERT INTO members (name, email, password_hash, activated)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, version`

	args := []interface{}{member.Name, member.Email, member.Password.hash, member.Activated}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&member.ID, &member.CreatedAt, &member.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "members_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}

func (m MemberModel) GetByEmail(email string) (*Member, error) {
	query := `
		SELECT id, created_at, name, email, password_hash, activated, version
		FROM members
		WHERE email = $1`

	var member Member

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&member.ID,
		&member.CreatedAt,
		&member.Name,
		&member.Email,
		&member.Password.hash,
		&member.Activated,
		&member.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &member, nil
}

func (m MemberModel) Update(member *Member) error {
	query := `
		UPDATE members
		SET name = $1, email = $2, password_hash = $3, activated = $4, version = version + 1
		WHERE id = $5 AND version = $6
		RETURNING version`

	args := []interface{}{
		member.Name,
		member.Email,
		member.Password.hash,
		member.Activated,
		member.ID,
		member.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&member.Version)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "members_email_key"`:
			return ErrDuplicateEmail
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (m MemberModel) GetForToken(tokenScope, tokenPlaintext string) (*Member, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))

	query := `
		SELECT members.id, members.created_at, members.name, members.email, members.password_hash, members.activated, members.version
		FROM members
		INNER JOIN tokens
		ON members.id = tokens.member_id
		WHERE tokens.hash = $1
		AND tokens.scope = $2
		AND tokens.expiry > $3`

	args := []interface{}{tokenHash[:], tokenScope, time.Now()}

	var member Member

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&member.ID,
		&member.CreatedAt,
		&member.Name,
		&member.Email,
		&member.Password.hash,
		&member.Activated,
		&member.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &member, nil
}
