//go:build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/hong-l1/project/wire/repository"
	"github.com/hong-l1/project/wire/repository/Dao"
)

func InitRepository() *repository.UseRepository {
	wire.Build(repository.NewUserRepository,
		Dao.NewUserDao, InitDB)
	return new(repository.UseRepository)
}
