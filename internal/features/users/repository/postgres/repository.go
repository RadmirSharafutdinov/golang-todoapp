package users_postgres_repository

import core_postgres_poll "study/internal/core/repository/postgres/pool"

type UsersRepository struct {
	pool core_postgres_poll.Pool
}

func NewUsersRepository(
	pool core_postgres_poll.Pool,
) *UsersRepository {
	return &UsersRepository{
		pool: pool,
	}
}
