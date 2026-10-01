package users_service

import (
	"context"
	"fmt"
	"study/internal/core/domain"
)

func (s *UsersService) GetUser(
	ctx context.Context,
	id int,
) (domain.User, error) {
	user, err := s.UsersRepository.GetUser(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("get users from repository:%w", err)
	}
	return user, nil
}
