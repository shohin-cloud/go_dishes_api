package model

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type Permissions []string

func (p Permissions) Include(code string) bool {
	for i := range p {
		if code == p[i] {
			return true
		}
	}
	return false
}

type PermissionModel struct {
	DB *sql.DB
}

func (m PermissionModel) GetAllForMember(memberID int64) (Permissions, error) {
	query := `
	SELECT permissions.code
	FROM permissions
	INNER JOIN members_permissions ON members_permissions.permission_id = permissions.id
	INNER JOIN members ON members_permissions.member_id = members.id
	WHERE members.id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, memberID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions Permissions

	for rows.Next() {
		var permission string
		err := rows.Scan(&permission)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}

func (m PermissionModel) AddForMember(memberID int64, codes ...string) error {
	query := `
	INSERT INTO members_permissions
	SELECT $1, permissions.id FROM permissions WHERE permissions.code = ANY($2)`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, memberID, pq.Array(codes))
	return err
}
