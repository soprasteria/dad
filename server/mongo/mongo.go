package mongo

import (
	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2"

	"strings"

	"github.com/spf13/viper"
)

// Session stores mongo session
var session *mgo.Session

// Connect connects to mongodb
func Connect() {
	uri := viper.GetString("server.mongo.addr")
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
func Get() (*mgo.Database, error) {
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
	return database, nil
}
