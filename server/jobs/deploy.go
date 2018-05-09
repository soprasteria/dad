package jobs

import (
	"fmt"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/robfig/cron"
	"github.com/soprasteria/dad/server/docktor"
	"github.com/soprasteria/dad/server/mongo"
	"github.com/soprasteria/dad/server/types"
	"github.com/spf13/viper"
)

func getDocktorGroupData(project types.Project) (docktor.GroupDocktor, error) {
	docktorAPI, err := docktor.NewExternalAPI(
		viper.GetString("docktor.addr"),
		viper.GetString("docktor.user"),
		viper.GetString("docktor.password"),
	)
	if err != nil {
		log.WithFields(log.Fields{
			"address":  viper.GetString("docktor.addr"),
			"username": viper.GetString("docktor.user"),
		}).WithError(err).Error("Unable to connect to Docktor")
		return docktor.GroupDocktor{}, err
	}

	idDocktorGroup, err := docktorAPI.GetGroupIDFromURL(project.DocktorGroupURL)
	if err != nil {
		log.WithField("docktorGroupURL", project.DocktorGroupURL).WithError(err).Error("Error when parsing Docktor group ID")
		return docktor.GroupDocktor{}, err
	}

	docktorGroup, err := docktorAPI.GetGroup(idDocktorGroup)
	if err != nil {
		log.WithError(err).Error("Error when getting containers services")
		return docktor.GroupDocktor{}, err
	}

	return docktorGroup, nil
}

func isIsolatedNetwork(docktorGroupData docktor.GroupDocktor) bool {
	for _, container := range docktorGroupData.Containers {
		if container.ServiceTitle == "ISOLATED_NETWORK" {
			return true
		}
	}
	return false
}

func isCloud(docktorGroupData docktor.GroupDocktor) bool {
	for _, container := range docktorGroupData.Containers {
		if container.ServiceTitle == "CLOUD" {
			return true
		}
	}
	return false
}

func getAllFunctionalServicesDeployByProject(docktorGroupData docktor.GroupDocktor) ([]types.FunctionalService, error) {
	// Connect to mongo
	database, err := mongo.Get()
	if err != nil {
		log.WithError(err).Error("Unable to connect to the database")
		return nil, err
	}
	defer database.Session.Close()

	// Formatting to an array of services
	servicesDeployed := []string{}
	for _, container := range docktorGroupData.Containers {
		servicesDeployed = append(servicesDeployed, strings.ToLower(container.ServiceTitle))
	}

	// Find all deployed functional services
	functionalServicesDeployed, err := database.FunctionalServices.FindFunctionalServicesDeployByServices(servicesDeployed)
	if err != nil {
		log.WithError(err).Error("Error when getting functional services")
		return nil, err
	}

	return functionalServicesDeployed, nil
}

// constructFullMatrix updates the matrix of a project to make it exhaustive
func constructFullMatrix(project *types.Project, functionalServices []types.FunctionalService) {
	for _, functionalService := range functionalServices {
		found := false

		for key, matrixLine := range project.Matrix {
			// If the line is already present in the matrix, set its "Deployed" status
			// to yes and updates its progress if needed
			if matrixLine.Service == functionalService.ID {
				project.Matrix[key].Deployed = types.Deployed[0]
				if matrixLine.Progress < 1 {
					project.Matrix[key].Progress = 1
				}
				found = true
				break
			}
		}

		// If the line is not present, create it with default values
		if !found {
			project.Matrix = append(project.Matrix, types.MatrixLine{
				Service:  functionalService.ID,
				Deployed: types.Deployed[0],
				Progress: 1,
			})
		}
	}
}

// ExecuteDeploymentStatusAnalytics calculates whether a functional service are deployed or not for all projects.
// Each fonctional services has some container which provide this service, if a proper container is deployed for a project, this service should be at 20% of progression at least.
// Else, it should be at 0% or N/A if the administrator of the project has defined this fonctionnal service as N/A.
func ExecuteDeploymentStatusAnalytics() (string, error) {

	log.Info("Starting to compute deployment status analytics...")
	// Connect to mongo
	database, err := mongo.Get()
	if err != nil {
		log.WithError(err).Error("Unable to connect to the database. Analytics are stopped.")
		return "", err
	}
	defer database.Session.Close()

	// Get all the projects which have docktor group url
	projects, err := database.Projects.FindWithDocktorGroupURL()
	if err != nil {
		log.WithError(err).Error("Unable to find projets with docktor group url. Analytics are stopped.")
		return "", err
	}

	log.Infof("Found %s projects with a Docktor URL, target as potentially updatable with deployment status.")
	updatedProjects := 0
	projectsInError := []string{}
	for _, project := range projects {

		docktorGroupData, err := getDocktorGroupData(project)
		if err != nil {
			log.WithError(err).Error("Error while retrieving Docktor data")
			time.Sleep(1 * time.Second) // Let Docktor catch his breath when an error occurred.
			continue
		}

		// In the case of an isolated network or on the cloud, all services are declarative, so we don't check anything in deploy and progress status.
		if isIsolatedNetwork(docktorGroupData) || !isCloud(docktorGroupData) {
			continue
		}

		// check if declarative and default not deployed, unless we are in isolated network
		for key, MatrixLine := range project.Matrix {
			// get the functional service info
			functionalService, err := database.FunctionalServices.FindByID(MatrixLine.Service.Hex())
			if err != nil {
				continue
			}
			// check if declarative
			if !functionalService.DeclarativeDeployment {
				project.Matrix[key].Deployed = types.Deployed[-1]
			}
		}

		// Get all functional services deployed
		functionalServices, err := getAllFunctionalServicesDeployByProject(docktorGroupData)
		if err != nil {
			projectsInError = append(projectsInError, fmt.Sprintf("%v (docktor:%v)", project.Name, project.DocktorGroupName))
			log.WithError(err).WithField("project", project.ID).Warn("Error while retrieving functional services")
			continue
		}

		// Waiting a little for Docktor to accept new incoming request
		time.Sleep(50 * time.Millisecond)

		constructFullMatrix(&project, functionalServices)

		// Put all the no deployed services to a progress of 0
		// If the project is isolated or on the cloud, don't touch anything
		if !isIsolatedNetwork(docktorGroupData) || !isCloud(docktorGroupData) {
			for key, matrixLine := range project.Matrix {
				if matrixLine.Deployed == types.Deployed[-1] {
					project.Matrix[key].Progress = 0
				}
			}
		}

		// Save
		_, err = database.Projects.Save(project)
		if err != nil {
			projectsInError = append(projectsInError, fmt.Sprintf("%v (docktor:%v)", project.Name, project.DocktorGroupName))
			log.WithError(err).WithField("project", project.ID).Warn("Error when updating the project")
			continue
		}
		updatedProjects++
	}
	log.Info("Computing deployment status analytics is over")
	return fmt.Sprintf("%v projects updated, %v not updated because an error occurred. List of projects in error [%v]",
		updatedProjects, len(projectsInError), strings.Join(projectsInError, ",")), nil
}

// jobDeploy execute deployment statistics
func jobDeploy(scheduler cron.Schedule) {
	message, err := ExecuteDeploymentStatusAnalytics()
	if err != nil {
		log.WithError(err).Error("Could not execute deployment status analatics")
	} else {
		log.Info(message)
	}
	log.Infof("Deployment indicators will computed next at %s", scheduler.Next(time.Now()))
}
