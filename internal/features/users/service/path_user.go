package users_service

import (
	"context"
	"fmt"
	"study/internal/core/domain"
)

func (s *UsersService) PathUser(
	ctx context.Context,
	id int,
	patch domain.UserPatch,
) (domain.User, error) {
	user, err := s.UsersRepository.GetUser(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("get user:%w", err)
	}
	if err := user.ApplyPatch(patch); err != nil {
		return domain.User{}, fmt.Errorf("apply user patch:%w", err)
	}
	patchedUser, err := s.UsersRepository.PatchUser(ctx, id, user)
	if err != nil {
		return domain.User{}, fmt.Errorf("patch user:%w", err)
	}
	return patchedUser, err
}
