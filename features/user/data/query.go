package data

import (
	"api/features/user"
	"errors"
	"log"

	"gorm.io/gorm"
)

type userQuery struct {
	db *gorm.DB
}

func New(db *gorm.DB) user.UserData {
	return &userQuery{
		db: db,
	}
}

func (uq *userQuery) Login(email string) (user.Core, error) {
	res := User{}

	if err := uq.db.Where("email = ?", email).First(&res).Error; err != nil {
		log.Println("login query error", err.Error())
		return user.Core{}, errors.New("data not found")
	}

	return ToCore(res), nil
}
func (uq *userQuery) Register(newUser user.Core) (user.Core, error) {
	cnv := CoreToData(newUser)
	err := uq.db.Create(&cnv).Error
	if err != nil {
		return user.Core{}, err
	}

	newUser.ID = cnv.ID

	return newUser, nil
}
func (uq *userQuery) Profile(id uint) (user.Core, error) {
	res := User{}
	if err := uq.db.Where("id = ?", id).First(&res).Error; err != nil {
		log.Println("Get By ID query error", err.Error())
		return user.Core{}, err
	}

	return ToCore(res), nil
}

func (uq *userQuery) Update(id uint, updateData user.Core) (user.Core, error) {
	cnv := CoreToData(updateData)

	qry := uq.db.Model(&User{}).Where("id = ?", id).Updates(&cnv)
	err := qry.Error

	if err != nil {
		log.Println("update query error ", err.Error())
		return user.Core{}, nil
	}

	return updateData, nil
}

func (uq *userQuery) Deactive(id uint) error {
	qry := uq.db.Delete(&User{}, id)

	affrows := qry.RowsAffected

	if affrows <= 0 {
		log.Println("no rows affected")
		return errors.New("tidak ada data yang dihapus")
	}

	err := qry.Error

	if err != nil {
		log.Println("delete query error")
		return errors.New("tidak bisa menghapus data")
	}

	return nil
}
