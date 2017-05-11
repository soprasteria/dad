package types

import (
	"errors"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// UsageIndicator represents a Sopra Steria entity
type UsageIndicator struct {
	ID bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	// Name of the Docktor Group, owner of the service instance which generated this indicator
	DocktorGroup string `bson:"docktorGroup" json:"docktorGroup,omitempty"`
	// Name of the service generating the indicator. e.g. jenkins
	Service string `bson:"service,omitempty" json:"service,omitempty"`
	// Name of the instance of service generating the indicator. e.g. GROUP1-jenkins
	ServiceInstanceName string `bson:"serviceInstance" json:"serviceInstance,omitempty"`
	// Indicator status of the service instance. e.g. active, inactive, undetermined...
	Status string `bson:"status,omitempty" json:"status,omitempty"`
	// Date when the indicator was last updated
	Updated time.Time `bson:"updated,omitempty" json:"updated,omitempty"`
}

// UsageIndicatorRepo wraps all requests to database for accessing usage indicators
type UsageIndicatorRepo struct {
	database *mgo.Database
}

// NewUsageIndicatorRepo creates a new usage indicators repo from database
// This UsageIndicatorRepo is wrapping all requests with database
func NewUsageIndicatorRepo(database *mgo.Database) UsageIndicatorRepo {
	return UsageIndicatorRepo{database: database}
}

func (r *UsageIndicatorRepo) col() *mgo.Collection {
	return r.database.C("usageIndicators")
}

func (r *UsageIndicatorRepo) isInitialized() bool {
	return r.database != nil
}

// CreateIndexes creates Index
func (r *UsageIndicatorRepo) CreateIndexes() error {
	if !r.isInitialized() {
		return ErrDatabaseNotInitialized
	}
	return r.col().EnsureIndex(mgo.Index{
		Key:      []string{"docktorGroup", "service"},
		Unique:   true,
		DropDups: true,
	})
}

// FindAll get all usage indicators from the database
func (r *UsageIndicatorRepo) FindAll() ([]UsageIndicator, error) {
	if !r.isInitialized() {
		return []UsageIndicator{}, ErrDatabaseNotInitialized
	}
	usageIndicators := []UsageIndicator{}
	err := r.col().Find(bson.M{}).All(&usageIndicators)
	if err != nil {
		return []UsageIndicator{}, errors.New("Can't retrieve all usage indicators")
	}
	return usageIndicators, nil
}

// FindAllFromGroup get all usage indicators with a given Docktor group
func (r *UsageIndicatorRepo) FindAllFromGroup(docktorGroup string) ([]UsageIndicator, error) {
	if !r.isInitialized() {
		return []UsageIndicator{}, ErrDatabaseNotInitialized
	}
	usageIndicators := []UsageIndicator{}

	err := r.col().Find(bson.M{"docktorGroup": docktorGroup}).All(&usageIndicators)
	if err != nil {
		return []UsageIndicator{}, errors.New("Can't retrieve any usage indicators")
	}
	return usageIndicators, nil
}

// BulkImportUsageIndicatorsResults is the result of a bulk import of usage indicators
type BulkImportUsageIndicatorsResults struct {
	All      int                `json:"all"`      // Number of usage indicators to import
	Imported int                `json:"imported"` // Number of usage indicators actually imported (inserted or updated)
	InError  int                `json:"inError"`  // Number of usage indicators not imported because an error happened
	Errors   []IndicatorInError `json:"errors"`   // Occurred errors and its details
}

// IndicatorInError represents an indicator that could not be imported because an error occurred
type IndicatorInError struct {
	Indicator UsageIndicator `json:"indicator"`
	Message   string         `json:"message"`
	Index     int            `json:"index"` // Index of usage indicator in error, in original slice
}

// BulkImport imports a list of indicators usages
// It updates existing indicators (with given service and Docktor group name), or create new ones.
func (r *UsageIndicatorRepo) BulkImport(usageIndicators []UsageIndicator) (BulkImportUsageIndicatorsResults, error) {
	if !r.isInitialized() {
		return BulkImportUsageIndicatorsResults{}, ErrDatabaseNotInitialized
	}

	errs := []IndicatorInError{}

	// Bulk upsert documents with given service and Docktor group name
	// A single operation is skipped when an error occurred while processing
	b := r.col().Bulk()
	b.Unordered()
	for i, indicator := range usageIndicators {
		// Handle business errors
		if indicator.DocktorGroup == "" || indicator.Service == "" || indicator.Status == "" {
			errs = append(errs, IndicatorInError{
				Indicator: indicator,
				Message:   "Docktor group name, status and service are mandatory and should not be empty. Skipped",
				Index:     i,
			})
			continue
		}
		indicator.Updated = time.Now()
		b.Upsert(bson.M{"docktorGroup": indicator.DocktorGroup, "service": indicator.Service}, indicator)
	}
	_, err := b.Run()

	// Handles technical errors, not previously handled by business checking
	if err != nil {
		ecases := err.(*mgo.BulkError).Cases()
		for _, c := range ecases {
			indicatorInError := UsageIndicator{}
			if c.Index >= 0 && c.Index < len(usageIndicators) {
				indicatorInError = usageIndicators[c.Index]
			}
			errs = append(errs, IndicatorInError{
				Indicator: indicatorInError,
				Message:   c.Err.Error(),
				Index:     c.Index,
			})
		}
	}

	return BulkImportUsageIndicatorsResults{
		All:      len(usageIndicators),
		Imported: len(usageIndicators) - len(errs),
		InError:  len(errs),
		Errors:   errs,
	}, nil
}
