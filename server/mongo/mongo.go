package mongo

import (
	log "github.com/Sirupsen/logrus"
	"github.com/soprasteria/dad/server/types"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"

	"strings"
	"time"
)

const mongoTimeout = 10 * time.Second

//DadMongo containers all types of Mongo data ready to be used
type DadMongo struct {
	Users              types.UserRepo              // Repo for accessing users methods
	Entities           types.EntityRepo            // Repo for accessing entities methods
	FunctionalServices types.FunctionalServiceRepo // Repo for accessing functional services methods
	Projects           types.ProjectRepo           // Repo for accessing projects methods
	Technologies       types.TechnologyRepo        // Repo for accessing technologies methods
	UsageIndicators    types.UsageIndicatorRepo    // Repo for accessing usage indicators methods
	Languages          types.LanguageRepo          // Repo for accessing languages methods
	Session            *mgo.Session                // Cloned session
	collections        []types.IsCollection        // Cache for listing all collections. Useful when doing operations on all collections at once (e.g. index creation at startup)
}

// CreateIndexes creates all indexes for every collections if needed
func (dm *DadMongo) CreateIndexes() {
	if dm.collections != nil {
		for _, db := range dm.collections {
			if dbWithIndex, ok := db.(types.IsCollectionWithIndexes); ok {
				err := dbWithIndex.CreateIndexes()
				if err != nil {
					log.WithError(err).Error("Cannot create index")
				}
			}
		}
	}
}

// Session stores mongo session
var session *mgo.Session

// Connect connects to mongodb
func Connect() {
	// Check availability of Mongo
	uri := viper.GetString("server.mongo.addr")
	if uri == "" {
		panic("Mongo url is empty. A Mongo database is required for Dad to work.")
	}
	s, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:   strings.Split(uri, ","),
		Timeout: mongoTimeout,
	})

	if err != nil {
		log.WithError(err).Fatal("Can't connect to mongo")
	}
	s.SetSafe(&mgo.Safe{})
	log.Info("Connected to ", uri)
	session = s

	// Create needed indexes at startup
	dadConn, err := Get()
	if err != nil {
		log.WithError(err).Fatal("Can't get mongo session for index creation")
	}
	dadConn.CreateIndexes()
	dadConn.Session.Close()
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

	collections := []types.IsCollection{}
	users := types.NewUserRepo(database)
	entities := types.NewEntityRepo(database)
	functionalServices := types.NewFunctionalServiceRepo(database)
	usageIndicators := types.NewUsageIndicatorRepo(database)
	projects := types.NewProjectRepo(database)
	technologies := types.NewTechnologyRepo(database)
	languages := types.NewLanguageRepo(database)

	collections = append(collections, &users)
	collections = append(collections, &entities)
	collections = append(collections, &functionalServices)
	collections = append(collections, &usageIndicators)
	collections = append(collections, &projects)
	collections = append(collections, &technologies)
	collections = append(collections, &languages)

	return &DadMongo{
		Users:              users,
		Entities:           entities,
		FunctionalServices: functionalServices,
		UsageIndicators:    usageIndicators,
		Projects:           projects,
		Technologies:       technologies,
		Languages:          languages,
		Session:            s,
		collections:        collections,
	}, nil
}
