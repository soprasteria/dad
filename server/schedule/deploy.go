package schedule

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/robfig/cron"
	"github.com/soprasteria/dad/server/controllers"
	"github.com/soprasteria/dad/server/mongo"
	"github.com/soprasteria/dad/server/types"
)

// jobDeploy is a Schedule a task that check which container are deployed for each functional services.
// Each fonctional services has some container which provide this service, if a proper container is deployed for a project, this service should be at 20% of progression at least.
// Else, it should be at 0% or N/A if the administrator of the project has defined this fonctionnal service as N/A.
func jobDeploy(scheduler cron.Schedule) {

	// Connect to mongo
	database, err := mongo.Get()
	if err != nil {
		log.WithError(err).Error("Unable to connect to the database")
	}

	// Get all the projects which have docktor group url
	projects, err := database.Projects.FindWithDocktorGroupURL()
	if err != nil {
		log.WithError(err).Error("Unable to find projets with docktor group url")
	}

	for _, project := range projects {

		// Default not deployed
		for keyMatrix := range project.Matrix {
			project.Matrix[keyMatrix].Deployed = "no"
		}

		// Get all functionnal services deployed
		servicesC := controllers.FunctionalServices{}
		functionalServices, err := servicesC.GetAllFunctionnalServicesDeployByProject(project)
		if err != nil {
			log.WithError(err).Error("Error when getting all functionnal services")
			continue
		}

	OUTER:
		for _, functionalService := range functionalServices {
			log.Infof("%s - Functional service available : %s", project.Name, functionalService.Name)
			// Check if the matrix already exist
			for key, matrixLine := range project.Matrix {
				if matrixLine.Service == functionalService.ID {
					// Found
					project.Matrix[key].Service = functionalService.ID
					project.Matrix[key].Deployed = "yes"
					if matrixLine.Progress < 1 {
						project.Matrix[key].Progress = 1
					}
					continue OUTER
				}
			}
			// Not found
			project.Matrix = append(project.Matrix, types.MatrixLine{
				Service:  functionalService.ID,
				Deployed: "yes",
				Progress: 1,
			})
		}

		// Save
		_, err = database.Projects.Save(project)
		if err != nil {
			log.WithError(err).Error("Error when updating the project")
			continue
		}
	}
	log.Infof("Deployment indicators will computed next at %s", scheduler.Next(time.Now()))
}
