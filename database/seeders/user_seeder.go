package seeders

import (
	"user-service/constants"
	"user-service/domain/models"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func RunUserSeeder(db *gorm.DB) {
	password, _ := bcrypt.GenerateFromPassword([]byte("1234567890"), bcrypt.DefaultCost)
	user := models.User{
		UUID:     uuid.New(),
		Name:     "Administrator",
		Email:    "admin@clswork.com",
		Password: string(password),
		Phone:    "+6282219193211",
		RoleID:   constants.Admin,
	}

	err := db.FirstOrCreate(&user, models.User{Email: user.Email}).Error
	if err != nil {
		logrus.Errorf("failed to seed user: %v", err)
		panic(err)
	}

	logrus.Infof("user %s successfully seeded", user.Name)
}
