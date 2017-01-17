package mongo

import (
	log "github.com/Sirupsen/logrus"
	"github.com/soprasteria/dad/server/types"
	"gopkg.in/mgo.v2"

	"strings"

	"github.com/spf13/viper"
)

//DadMongo containers all types of Mongo data ready to be used
type DadMongo struct {
	Users   types.UserRepo // Repo for accessing users methods
	Session *mgo.Session   // Cloned session
}

// Session stores mongo session
var session *mgo.Session

// Connect connects to mongodb
func Connect() {
	uri := viper.GetString("server.mongo.addr")
	if uri == "" {
		panic("Mongo url is empty. A Mongo database is required for Dad to work.")
	}
	s, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs: strings.Split(uri, ","),
	})

	if err != nil {
		log.WithError(err).Fatal("Can't connect to mongo")
	}
	s.SetSafe(&mgo.Safe{})
	log.Info("Connected to ", uri)
	session = s
}

// Get the connexion to mongodb
func Get() (*DadMongo, error) {
	username := viper.GetString("server.mongo.username")
	password := viper.GetString("server.mongo.password")
	s := session.Clone()
	s.SetSafe(&mgo.Safe{})
	database := s.DB("dad")
	if username != "" && password != "" {
		err := database.Login(username, password)
		if err != nil {
			return nil, err
		}
	}

	users := types.NewUserRepo(database)

	return &DadMongo{
		Users:   users,
		Session: s,
	}, nil
}
