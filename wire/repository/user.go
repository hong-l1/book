package repository

import "github.com/hong-l1/project/wire/repository/Dao"

type UseRepository struct {
	dao *Dao.UserDao
}

func NewUserRepository(dao *Dao.UserDao) *UseRepository {
	return &UseRepository{
		dao: dao,
	}
}
