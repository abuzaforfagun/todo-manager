package auth_repository

import (
	"log"
	"restful-service/db"
	"restful-service/models"
)

func Register(userName string, password string) error {
	db := db.Get()

	sql, err := db.Prepare("INSERT INTO UserLogin (Username, Password) VALUES (?,?)")
	if err != nil {
		log.Println("Unable to prepare the sql", err)
		return err
	}

	_, err = sql.Exec(userName, password)
	if err != nil {
		log.Println("Unable to insert", err)
	}
	return err
}

func HasUser(username string) (hasUser bool, err error) {
	db := db.Get()

	sql, err := db.Prepare("SELECT EXISTS (SELECT 1 FROM UserLogin WHERE UserName = ?)")
	if err != nil {
		log.Println("Unable to prepare the sql", err)
		return false, err
	}
	userRow := sql.QueryRow(username)

	err = userRow.Scan(&hasUser)

	if err != nil {
		return false, err
	}

	return hasUser, err
}

func GetUser(username string) (user models.Credentials, err error) {
	db := db.Get()

	sql, err := db.Prepare("SELECT UserName, Password FROM UserLogin WHERE UserName = ?")
	if err != nil {
		log.Println("Unable to prepare the sql", err)
		return models.Credentials{}, err
	}
	userRow := sql.QueryRow(username)

	err = userRow.Scan(&user.Username, &user.Password)

	if err != nil {
		return models.Credentials{}, err
	}

	return user, err
}
