package wire

import (
	"fmt"
	"github.com/hong-l1/project/wire/repository"
	"github.com/hong-l1/project/wire/repository/Dao"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(mysql.Open("dsn"))
	if err != nil {
		panic(err)
	}
	ud := Dao.NewUserDao(db)
	rd := repository.NewUserRepository(ud)
	fmt.Println(rd)
	InitRepository()
}
