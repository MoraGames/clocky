package planetScale

import (
	"os"

	"github.com/MoraGames/clockyuwu/model"
	"github.com/MoraGames/clockyuwu/pkg/util"
	"github.com/MoraGames/clockyuwu/repo"
	"github.com/jmoiron/sqlx"
)

// Check if the repo implements the interface
var _ repo.UserRepoer = new(UserRepo)

// mock.UserRepo
type UserRepo struct {
	tcp               string
	dbname            string
	tls               bool
	interpolateParams bool
	dns               string
}

// Return a new UserRepo
func NewUserRepo(tcp, dbname string, tls, interpolateParams bool) *UserRepo {
	return &UserRepo{
		tcp:               tcp,
		dbname:            dbname,
		tls:               tls,
		interpolateParams: interpolateParams,
		dns:               os.Getenv("PLANETSCALEDB_USERNAME") + ":" + os.Getenv("PLANETSCALEDB_PASSWORD") + "@tcp(" + tcp + ")" + "/" + dbname + "?tls=" + util.BtoS(tls) + "&interpolateParams=" + util.BtoS(interpolateParams),
	}
}

func (ur *UserRepo) Create(user *model.User) error {
	// Open a connection to PlanetScale
	db, err := sqlx.Open("mysql", ur.dns)
	if err != nil {
		return err
	}
	defer db.Close()

	// Insert the user into the database
	_, err = db.Queryx("INSERT INTO Users (telegram_id) VALUES (?);", user.TelegramUser.ID)
	return err
}

func (ur *UserRepo) Get(id int64) (*model.User, error) {
	// Open a connection to PlanetScale
	db, err := sqlx.Open("mysql", ur.dns)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Query the database
	rows, err := db.Queryx("SELECT * FROM Users AS u WHERE u.telegram_id=?;", id)

	var user model.User
	for rows.Next() {
		err = rows.StructScan(&user)
		if err != nil {
			return nil, err
		}
	}

	return &user, nil
}

func (ur *UserRepo) GetAll() ([]*model.User, error) {
	// Open a connection to PlanetScale
	db, err := sqlx.Open("mysql", ur.dns)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Query the database
	rows, err := db.Queryx("SELECT * FROM Users;")

	var users []*model.User
	for rows.Next() {
		var user model.User
		err = rows.StructScan(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

func (ur *UserRepo) Update(id int64, user *model.User) error {
	// Open a connection to PlanetScale
	db, err := sqlx.Open("mysql", ur.dns)
	if err != nil {
		return err
	}
	defer db.Close()

	// Query the database
	_, err = db.Queryx("UPDATE Users AS u SET u.telegram_id=? WHERE u.telegram_id=?;", user.TelegramUser.ID, id)

	return err
}

func (ur *UserRepo) Delete(id int64) error {
	// Open a connection to PlanetScale
	db, err := sqlx.Open("mysql", ur.dns)
	if err != nil {
		return err
	}
	defer db.Close()

	// Query the database
	_, err = db.Queryx("DELETE FROM Users AS u WHERE u.telegram_id=?;", id)

	return err
}
