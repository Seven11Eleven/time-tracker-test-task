package database

import (
	"fmt"
	
    "context"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/Seven11Eleven/time-tracker-test-task/internal/models"
)

type UserRepository struct {
    db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
    user.CreatedAt = time.Now()
    user.UpdatedAt = user.CreatedAt
    query := `INSERT INTO users (passport_number, surname, name, patronymic, address, created_at, updated_at) 
              VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
    err := r.db.QueryRow(ctx, query, user.PassportNumber, user.Surname, user.Name, user.Patronymic, user.Address, user.CreatedAt, user.UpdatedAt).Scan(&user.ID)
    return err
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int) (*models.User, error) {
    user := &models.User{}
    query := `SELECT id, passport_number, surname, name, patronymic, address, created_at, updated_at FROM users WHERE id=$1`
    err := r.db.QueryRow(ctx, query, id).Scan(&user.ID, &user.PassportNumber, &user.Surname, &user.Name, &user.Patronymic, &user.Address, &user.CreatedAt, &user.UpdatedAt)
    return user, err
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *models.User) error {
    user.UpdatedAt = time.Now()
    query := `UPDATE users SET passport_number=$1, surname=$2, name=$3, patronymic=$4, address=$5, updated_at=$6 WHERE id=$7`
    _, err := r.db.Exec(ctx, query, user.PassportNumber, user.Surname, user.Name, user.Patronymic, user.Address, user.UpdatedAt, user.ID)
    return err
}

func (r *UserRepository) DeleteUser(ctx context.Context, id int) error {
    query := `DELETE FROM users WHERE id=$1`
    _, err := r.db.Exec(ctx, query, id)
    return err
}

func (r *UserRepository) GetUsers(ctx context.Context, filter map[string]interface{}, limit, offset int) ([]models.User, error) {
    var argID int = 1
    query := "SELECT id, passport_number, surname, name, patronymic, address, created_at, updated_at FROM users WHERE true"
    args := []interface{}{}
    
    for key, val := range filter{
        query += fmt.Sprintf(" AND %s = $%d", key, argID)
        args = append(args, val)
        argID++ 
    }

    query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argID, argID+1)
    args = append(args, limit, offset)
    rows, err := r.db.Query(ctx, query, args...)
    if err != nil{
        return nil, err
    }
    defer rows.Close()

    var users []models.User
    for rows.Next(){
        var user models.User
        if err := rows.Scan(&user.ID, &user.PassportNumber,&user.Surname, &user.Name, &user.Patronymic, &user.Address, &user.CreatedAt, &user.UpdatedAt); err != nil {
            return nil, err
        }

        users = append(users, user)
    }
    if rows.Err() != nil{
        return nil, rows.Err()
    }

    return users, nil
}