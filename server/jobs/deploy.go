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

func getAllFunctionalServicesDeployByProject(project types.Project) ([]types.FunctionalService, error) {

	// Connect to docktor api
	docktorAPI, err := docktor.NewExternalAPI(
		viper.GetString("docktor.addr"),
		viper.GetString("docktor.user"),
		viper.GetString("docktor.password"),
	)
	if err != nil {
		log.WithError(err).Error("Unable to connect to docktor")
		return nil, err
	}

	// Get the docktor group id
	idDocktorGroup, err := docktorAPI.GetGroupIDFromURL(project.DocktorGroupURL)
	if err != nil {
		log.WithError(err).Error("Error when parse group id")
		return nil, err
	}

	// Connect to docktor to get group info here Containers.ServiceTitle
	docktorGroup, err := docktorAPI.GetGroup(idDocktorGroup)
	if err != nil {
		log.WithError(err).Error("Error when getting containers services")
		return nil, err
	}

	log.Infof("Services availables : %s", docktorGroup.Containers)

	// Connect to mongo
	database, err := mongo.Get()
	if err != nil {
		log.WithError(err).Error("Unable to connect to the database")
		return nil, err
	}
	defer database.Session.Close()

	// Formatting to an array of services
	servicesDeployed := []string{}
	for _, container := range docktorGroup.Containers {
		servicesDeployed = append(servicesDeployed, strings.ToLower(container.ServiceTitle))
	}

	// Find all deployed functional services
	functionalServicesDeployed, err := database.FunctionalServices.FindFunctionalServicesDeployByServices(servicesDeployed)
	if err != nil {
		log.WithError(err).Error("Error when getting functional services")
		return nil, err
	}

	return functionalServicesDeployed, err
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

		// check if declarative and default not deployed
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
		functionalServices, err := getAllFunctionalServicesDeployByProject(project)
		if err != nil {
			// When an error occurred, it's often because of
			projectsInError = append(projectsInError, fmt.Sprintf("%v (docktor:%v)", project.Name, project.DocktorGroupName))
			log.WithError(err).WithField("project", project.ID).Warn("Error when getting all functional services")
			time.Sleep(1 * time.Second) // Let Docktor catch his breath when an error occurred.
			continue
		}

		// Waiting a little for Docktor to accept new incoming request
		time.Sleep(50 * time.Millisecond)

	OUTER:
		for _, functionalService := range functionalServices {
			log.Infof("%s - Functional service available : %s", project.Name, functionalService.Name)
			// Check if the matrix already exist
			for key, matrixLine := range project.Matrix {
				if matrixLine.Service == functionalService.ID {
					// Found
					project.Matrix[key].Deployed = types.Deployed[0]
					if matrixLine.Progress < 1 {
						project.Matrix[key].Progress = 1
					}
					continue OUTER
				}
			}
			// Not found
			project.Matrix = append(project.Matrix, types.MatrixLine{
				Service:  functionalService.ID,
				Deployed: types.Deployed[0],
				Progress: 1,
			})
		}

		// Put all the no deployed services to a progress of 0
		for key, matrixLine := range project.Matrix {
			if matrixLine.Deployed == types.Deployed[-1] {
				project.Matrix[key].Progress = 0
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
	log.Infof("Computing deployment status analytics is over.")
	return fmt.Sprintf("%v projects updated, %v not updated because an error occurred. List of projects in error [%v].",
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
