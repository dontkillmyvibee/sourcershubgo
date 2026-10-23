package postgres

import "github.com/dontkillmyvibee/sourcershubgo/internal/core/repository/postgres/pool"

type UsersRepository struct {
	pool pool.Pool
}

func NewUsersRepository(pool pool.Pool) *UsersRepository {
	return &UsersRepository{
		pool: pool,
	}
}
