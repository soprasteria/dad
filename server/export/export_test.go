package export

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/soprasteria/dad/server/types"
)

func TestGetServiceToIndicatorUsage(t *testing.T) {

	mapToSlice := func(usageIndicator map[string]types.UsageIndicator) []string {
		usageIndicatorSliced := make([]string, len(usageIndicator))
		i := 0
		for k := range usageIndicator {
			usageIndicatorSliced[i] = append(usageIndicatorSliced, k)
		}
		return usageIndicatorSliced
	}

	technicalServices := []string{"jenkins", "gitlabci", "tfs"}

	jenkinsIndicator := types.UsageIndicator{Service: "jenkins"}
	gitlabciIndicator := types.UsageIndicator{Service: "gitlabci"}
	tfsIndicator := types.UsageIndicator{Service: "tfs"}
	unknownIndicator := types.UsageIndicator{Service: "unknown"}

	Convey("Given 0 functional service and 0 indicator", t, func() {
		service := types.FunctionalService{}
		usageIndicators := []types.UsageIndicator{}
		Convey("When calling the getServiceToIndicatorUsage function", func() {
			usageIndicator := getServiceToIndicatorUsage(service, usageIndicators)
			Convey("Then the map is empty", func() {
				So(usageIndicator, ShouldBeEmpty)
			})
		})
	})

	Convey("Given 1 functional service and 0 indicator", t, func() {
		service := types.FunctionalService{Services: technicalServices}
		usageIndicators := []types.UsageIndicator{}
		Convey("When calling the getServiceToIndicatorUsage function", func() {
			usageIndicator := getServiceToIndicatorUsage(service, usageIndicators)
			Convey("Then the map is empty", func() {
				So(usageIndicator, ShouldBeEmpty)
			})
		})
	})

	Convey("Given 1 functional service and 1 indicator", t, func() {
		Convey("with indicator services doesn't matching services", func() {
			service := types.FunctionalService{Services: technicalServices}
			usageIndicators := []types.UsageIndicator{unknownIndicator}
			Convey("When calling the getServiceToIndicatorUsage function", func() {
				usageIndicator := getServiceToIndicatorUsage(service, usageIndicators)
				Convey("Then the map is empty", func() {
					So(usageIndicator, ShouldBeEmpty)
				})
			})
		})

		Convey("with indicator services matching services", func() {
			service := types.FunctionalService{Services: technicalServices}
			usageIndicators := []types.UsageIndicator{jenkinsIndicator}
			Convey("When calling the getServiceToIndicatorUsage function", func() {
				usageIndicator := getServiceToIndicatorUsage(service, usageIndicators)
				Convey("Then the map contains 1 indicator", func() {
					So(usageIndicator, ShouldHaveLength, 1)
				})
				Convey("and indicator status have been correctly assigned to the map", func() {
					So(usageIndicator, ShouldContainKey, jenkinsIndicator.Service)
					usageIndicatorSliced := mapToSlice(usageIndicator)
					So(jenkinsIndicator.Service, ShouldBeIn, usageIndicatorSliced)
				})
			})
		})
	})

	Convey("Given 1 functional service and 2 indicators", t, func() {
		Convey("with 1 indicator services doesn't matching services", func() {
			service := types.FunctionalService{Services: technicalServices}
			usageIndicators := []types.UsageIndicator{jenkinsIndicator, unknownIndicator}
			Convey("When calling the getServiceToIndicatorUsage function", func() {
				usageIndicator := getServiceToIndicatorUsage(service, usageIndicators)
				Convey("Then the map contains 1 indicator", func() {
					So(usageIndicator, ShouldHaveLength, 1)
				})
				Convey("and indicator status have been correctly assigned to the map", func() {
					So(usageIndicator, ShouldContainKey, jenkinsIndicator.Service)
					So(usageIndicator, ShouldNotContainKey, unknownIndicator.Service)
					usageIndicatorSliced := mapToSlice(usageIndicator)
					So(jenkinsIndicator.Service, ShouldBeIn, usageIndicatorSliced)
					So(unknownIndicator.Service, ShouldNotBeIn, usageIndicatorSliced)
				})
			})
		})

		Convey("with indicator services matching services", func() {
			service := types.FunctionalService{Services: technicalServices}
			usageIndicators := []types.UsageIndicator{jenkinsIndicator, gitlabciIndicator, tfsIndicator}
			Convey("When calling the getServiceToIndicatorUsage function", func() {
				usageIndicator := getServiceToIndicatorUsage(service, usageIndicators)
				Convey("Then the map contains 3 indicators", func() {
					So(usageIndicator, ShouldHaveLength, 3)
				})
				Convey("and indicator status have been correctly assigned to the map", func() {
					So(usageIndicator, ShouldContainKey, jenkinsIndicator.Service)
					So(usageIndicator, ShouldContainKey, gitlabciIndicator.Service)
					So(usageIndicator, ShouldContainKey, tfsIndicator.Service)
					usageIndicatorSliced := mapToSlice(usageIndicator)
					So(jenkinsIndicator.Service, ShouldBeIn, usageIndicatorSliced)
					So(gitlabciIndicator.Service, ShouldBeIn, usageIndicatorSliced)
					So(tfsIndicator.Service, ShouldBeIn, usageIndicatorSliced)
				})
			})
		})
	})
}

func TestFunctionBestIndicatorStatus(t *testing.T) {

	jenkinsIndicator := types.UsageIndicator{Service: "jenkins"}
	gitlabciIndicator := types.UsageIndicator{Service: "gitlabci"}

	Convey("Given 0 technical service and 0 indicator", t, func() {
		service := types.FunctionalService{}
		usageIndicators := []types.UsageIndicator{}
		Convey("When calling the bestIndicatorStatus function", func() {
			status := bestIndicatorStatus(service, usageIndicators)
			Convey("Then the result is nil", func() {
				So(status, ShouldBeNil)
			})
		})
	})

	Convey("Given 0 technical service and 1 indicator", t, func() {
		service := types.FunctionalService{}
		usageIndicators := []types.UsageIndicator{jenkinsIndicator}
		Convey("When calling the bestIndicatorStatus function", func() {
			status := bestIndicatorStatus(service, usageIndicators)
			Convey("Then the result is nil", func() {
				So(status, ShouldBeNil)
			})
		})
	})

	Convey("Given 1 technical service and 0 indicator", t, func() {
		service := types.FunctionalService{Services: []string{"jenkins"}}
		usageIndicators := []types.UsageIndicator{}
		Convey("When calling the bestIndicatorStatus function", func() {
			status := bestIndicatorStatus(service, usageIndicators)
			Convey("Then the result is nil", func() {
				So(status, ShouldBeNil)
			})
		})
	})

	Convey("Given 1 technical service and 1 indicator", t, func() {
		Convey("with the indicator service doesn't match the service", func() {
			service := types.FunctionalService{Services: []string{"jenkins"}}
			usageIndicators := []types.UsageIndicator{gitlabciIndicator}
			Convey("When calling the bestIndicatorStatus function", func() {
				status := bestIndicatorStatus(service, usageIndicators)
				Convey("Then the result is nil", func() {
					So(status, ShouldBeNil)
				})
			})
		})

		Convey("with the indicator service matching the service", func() {
			Convey("and the indicator equal to Empty", func() {
				service := types.FunctionalService{Services: []string{"jenkins"}}
				usageIndicators := []types.UsageIndicator{{Service: "jenkins", Status: "Empty"}}
				Convey("When calling the bestIndicatorStatus function", func() {
					status := bestIndicatorStatus(service, usageIndicators)
					Convey("Then the result is Empty", func() {
						So(status.String(), ShouldEqual, "Empty")
					})
				})
			})

			Convey("and the indicator equal to Undetermined", func() {
				service := types.FunctionalService{Services: []string{"jenkins"}}
				usageIndicators := []types.UsageIndicator{{Service: "jenkins", Status: "Undetermined"}}
				Convey("When calling the bestIndicatorStatus function", func() {
					status := bestIndicatorStatus(service, usageIndicators)
					Convey("Then the result is Undetermined", func() {
						So(status.String(), ShouldEqual, "Undetermined")
					})
				})
			})

			Convey("and the indicator equal to Inactive", func() {
				service := types.FunctionalService{Services: []string{"jenkins"}}
				usageIndicators := []types.UsageIndicator{{Service: "jenkins", Status: "Inactive"}}
				Convey("When calling the bestIndicatorStatus function", func() {
					status := bestIndicatorStatus(service, usageIndicators)
					Convey("Then the result is Inactive", func() {
						So(status.String(), ShouldEqual, "Inactive")
					})
				})
			})

			Convey("and the indicator equal to Active", func() {
				service := types.FunctionalService{Services: []string{"jenkins"}}
				usageIndicators := []types.UsageIndicator{{Service: "jenkins", Status: "Active"}}
				Convey("When calling the bestIndicatorStatus function", func() {
					status := bestIndicatorStatus(service, usageIndicators)
					Convey("Then the result is Active", func() {
						So(status.String(), ShouldEqual, "Active")
					})
				})
			})
		})
	})
	Convey("Given 2 technical services and 2 indicators", t, func() {
		Convey("with indicator services matching services", func() {
			Convey("and both indicators equal to Empty", func() {
				service := types.FunctionalService{Services: []string{"jenkins", "gitlabci"}}
				usageIndicators := []types.UsageIndicator{{Service: "jenkins", Status: "Empty"}, {Service: "gitlabci", Status: "Empty"}}
				Convey("When calling the bestIndicatorStatus function", func() {
					status := bestIndicatorStatus(service, usageIndicators)
					Convey("Then the result is Empty", func() {
						So(status.String(), ShouldEqual, "Empty")
					})
				})
			})

			Convey("and one indicator equal to Empty and the other equal to Undetermined", func() {
				service := types.FunctionalService{Services: []string{"jenkins", "gitlabci"}}
				usageIndicators := []types.UsageIndicator{{Service: "jenkins", Status: "Empty"}, {Service: "gitlabci", Status: "Undetermined"}}
				Convey("When calling the bestIndicatorStatus function", func() {
					status := bestIndicatorStatus(service, usageIndicators)
					Convey("Then the result is Undetermined", func() {
						So(status.String(), ShouldEqual, "Undetermined")
					})
				})
			})

			Convey("and one indicator equal to Empty and the other equal to Inactive", func() {
				service := types.FunctionalService{Services: []string{"jenkins", "gitlabci"}}
				usageIndicators := []types.UsageIndicator{{Service: "jenkins", Status: "Empty"}, {Service: "gitlabci", Status: "Inactive"}}
				Convey("When calling the bestIndicatorStatus function", func() {
					status := bestIndicatorStatus(service, usageIndicators)
					Convey("Then the result is Inactive", func() {
						So(status.String(), ShouldEqual, "Inactive")
					})
				})
			})

			Convey("and one indicator equal to Empty and the other equal to Active", func() {
				service := types.FunctionalService{Services: []string{"jenkins", "gitlabci"}}
				usageIndicators := []types.UsageIndicator{{Service: "jenkins", Status: "Empty"}, {Service: "gitlabci", Status: "Active"}}
				Convey("When calling the bestIndicatorStatus function", func() {
					status := bestIndicatorStatus(service, usageIndicators)
					Convey("Then the result is Active", func() {
						So(status.String(), ShouldEqual, "Active")
					})
				})
			})

			Convey("and both indicators equal to Undetermined", func() {
				service := types.FunctionalService{Services: []string{"jenkins", "gitlabci"}}
				usageIndicators := []types.UsageIndicator{{Service: "jenkins", Status: "Undetermined"}, {Service: "gitlabci", Status: "Undetermined"}}
				Convey("When calling the bestIndicatorStatus function", func() {
					status := bestIndicatorStatus(service, usageIndicators)
					Convey("Then the result is Undetermined", func() {
						So(status.String(), ShouldEqual, "Undetermined")
					})
				})
			})

			Convey("and one indicator equal to Undetermined and the other equal to Inactive", func() {
				service := types.FunctionalService{Services: []string{"jenkins", "gitlabci"}}
				usageIndicators := []types.UsageIndicator{{Service: "jenkins", Status: "Undetermined"}, {Service: "gitlabci", Status: "Inactive"}}
				Convey("When calling the bestIndicatorStatus function", func() {
					status := bestIndicatorStatus(service, usageIndicators)
					Convey("Then the result is Inactive", func() {
						So(status.String(), ShouldEqual, "Inactive")
					})
				})
			})

			Convey("and one indicator equal to Undetermined and the other equal to Active", func() {
				service := types.FunctionalService{Services: []string{"jenkins", "gitlabci"}}
				usageIndicators := []types.UsageIndicator{{Service: "jenkins", Status: "Undetermined"}, {Service: "gitlabci", Status: "Active"}}
				Convey("When calling the bestIndicatorStatus function", func() {
					status := bestIndicatorStatus(service, usageIndicators)
					Convey("Then the result is Active", func() {
						So(status.String(), ShouldEqual, "Active")
					})
				})
			})

			Convey("and both indicators equal to Inactive", func() {
				service := types.FunctionalService{Services: []string{"jenkins", "gitlabci"}}
				usageIndicators := []types.UsageIndicator{{Service: "jenkins", Status: "Inactive"}, {Service: "gitlabci", Status: "Inactive"}}
				Convey("When calling the bestIndicatorStatus function", func() {
					status := bestIndicatorStatus(service, usageIndicators)
					Convey("Then the result is Inactive", func() {
						So(status.String(), ShouldEqual, "Inactive")
					})
				})
			})

			Convey("and one indicator equal to Inactive and the other equal to Active", func() {
				service := types.FunctionalService{Services: []string{"jenkins", "gitlabci"}}
				usageIndicators := []types.UsageIndicator{{Service: "jenkins", Status: "Inactive"}, {Service: "gitlabci", Status: "Active"}}
				Convey("When calling the bestIndicatorStatus function", func() {
					status := bestIndicatorStatus(service, usageIndicators)
					Convey("Then the result is Active", func() {
						So(status.String(), ShouldEqual, "Active")
					})
				})
			})

			Convey("and both indicators equal to Active", func() {
				service := types.FunctionalService{Services: []string{"jenkins", "gitlabci"}}
				usageIndicators := []types.UsageIndicator{{Service: "jenkins", Status: "Active"}, {Service: "gitlabci", Status: "Active"}}
				Convey("When calling the bestIndicatorStatus function", func() {
					status := bestIndicatorStatus(service, usageIndicators)
					Convey("Then the result is Active", func() {
						So(status.String(), ShouldEqual, "Active")
					})
				})
			})
		})
	})

	Convey("Given 3 technical services and 3 indicators", t, func() {
		Convey("with indicator services matching services", func() {
			Convey("and indicators equal respectively to [Empty, Undetermined, Inactive]", func() {
				service := types.FunctionalService{Services: []string{"jenkins", "gitlabci", "tfs"}}
				usageIndicators := []types.UsageIndicator{{Service: "jenkins", Status: "Empty"}, {Service: "gitlabci", Status: "Undetermined"}, {Service: "tfs", Status: "Inactive"}}
				Convey("When calling the bestIndicatorStatus function", func() {
					status := bestIndicatorStatus(service, usageIndicators)
					Convey("Then the result is Inactive", func() {
						So(status.String(), ShouldEqual, "Inactive")
					})
				})
			})

			Convey("and indicators equal respectively to [Empty, Inactive, Undetermined]", func() {
				service := types.FunctionalService{Services: []string{"jenkins", "gitlabci", "tfs"}}
				usageIndicators := []types.UsageIndicator{{Service: "jenkins", Status: "Empty"}, {Service: "gitlabci", Status: "Inactive"}, {Service: "tfs", Status: "Undetermined"}}
				Convey("When calling the bestIndicatorStatus function", func() {
					status := bestIndicatorStatus(service, usageIndicators)
					Convey("Then the result is Inactive", func() {
						So(status.String(), ShouldEqual, "Inactive")
					})
				})
			})

			Convey("and indicators equal respectively to [Inactive, Empty, Undetermined]", func() {
				service := types.FunctionalService{Services: []string{"jenkins", "gitlabci", "tfs"}}
				usageIndicators := []types.UsageIndicator{{Service: "jenkins", Status: "Inactive"}, {Service: "gitlabci", Status: "Empty"}, {Service: "tfs", Status: "Undetermined"}}
				Convey("When calling the bestIndicatorStatus function", func() {
					status := bestIndicatorStatus(service, usageIndicators)
					Convey("Then the result is Inactive", func() {
						So(status.String(), ShouldEqual, "Inactive")
					})
				})
			})
		})
	})
}

func TestGenerateXlsx(t *testing.T) {

	/*Convey("Given functional services with technical service 'Pipeline' initialized", t, func() {

	technicalServicePipeline := []string{"jenkins"}

	services := []types.FunctionalService{
		{
			ID:       bson.ObjectId("1"),
			Name:     "Collecte des exigences",
			Package:  "1. Requirement Management",
			Position: 10},
		{
			ID:       bson.ObjectId("2"),
			Name:     "Gestion du patrimoine des exigences",
			Package:  "1. Requirement Management",
			Position: 20},
		{
			ID:       bson.ObjectId("3"),
			Name:     "Visualisation de la traçabilité bout en bout des exigences",
			Package:  "1. Requirement Management",
			Position: 30},
		{
			ID:       bson.ObjectId("4"),
			Name:     "Pipeline d'intégration continue",
			Package:  "2. Build",
			Position: 10,
			Services: []string{"jenkins"}},
		{
			ID:       bson.ObjectId("5"),
			Name:     "Automatisation des tests unitaires",
			Package:  "2. Build",
			Position: 20},
		{
			ID:       bson.ObjectId("6"),
			Name:     "Gestion des dépendances et des artéfacts",
			Package:  "2. Build",
			Position: 30},
		{
			ID:       bson.ObjectId("7"),
			Name:     "Gestion des revues de code",
			Package:  "2. Build",
			Position: 40},
		{
			ID:       bson.ObjectId("8"),
			Name:     "Gestion de la qualimétrie",
			Package:  "2. Build",
			Position: 50},
		{
			ID:       bson.ObjectId("9"),
			Name:     "Packaging et déploiement automatique de la livraison",
			Package:  "2. Build",
			Position: 60},
		{
			ID:       bson.ObjectId("10"),
			Name:     "Gestion du patrimoine de tests",
			Package:  "3. Acceptance",
			Position: 10},
		{
			ID:       bson.ObjectId("11"),
			Name:     "Modèles de stratégie de tests",
			Package:  "3. Acceptance",
			Position: 20},
		{
			ID:       bson.ObjectId("12"),
			Name:     "Automatisation des tests de non régression",
			Package:  "3. Acceptance",
			Position: 30},
		{
			ID:       bson.ObjectId("13"),
			Name:     "Tests de sécurité applicative",
			Package:  "3. Acceptance",
			Position: 40},
		{
			ID:       bson.ObjectId("14"),
			Name:     "Simulation des flux pour l'automatisation des tests",
			Package:  "3. Acceptance",
			Position: 50},
		{
			ID:       bson.ObjectId("15"),
			Name:     "Dashboard d'avancement et bilan de tests",
			Package:  "3. Acceptance",
			Position: 60},
		{
			ID:       bson.ObjectId("16"),
			Name:     "Contrôle des licences",
			Package:  "3. Acceptance",
			Position: 70},
		{
			ID:       bson.ObjectId("17"),
			Name:     "Anonymisation des données sensibles d'une base",
			Package:  "3. Acceptance",
			Position: 80},
		{
			ID:       bson.ObjectId("18"),
			Name:     "Modèles de stratégie de tests",
			Package:  "4. Performance",
			Position: 10},
		{
			ID:       bson.ObjectId("19"),
			Name:     "Bilan de tests",
			Package:  "4. Performance",
			Position: 20},
		{
			ID:       bson.ObjectId("20"),
			Name:     "Automatisation des tests de performances",
			Package:  "4. Performance",
			Position: 30},
		{
			ID:       bson.ObjectId("21"),
			Name:     "Dashboard des performances",
			Package:  "4. Performance",
			Position: 40},
		{
			ID:       bson.ObjectId("22"),
			Name:     "Déploiement et installation automatique",
			Package:  "5. OPS",
			Position: 10},
		{
			ID:       bson.ObjectId("23"),
			Name:     "Analyse automatique des logs",
			Package:  "5. OPS",
			Position: 20},
		{
			ID:       bson.ObjectId("24"),
			Name:     "Vision de l'avancement global et suivi des environnements",
			Package:  "6. Monitoring",
			Position: 10},
		{
			ID:       bson.ObjectId("25"),
			Name:     "Suivi des versions et des demandes (Release Monitoring)",
			Package:  "6. Monitoring",
			Position: 20},
		{
			ID:       bson.ObjectId("26"),
			Name:     "Tweeter",
			Package:  "6. Monitoring",
			Position: 30}}

	Convey("When calling the generateXlsx function", func() {

		generateXlsx
		status := bestIndicatorStatus(service, usageIndicators)
		Convey("Then the result is Inactive", func() {
			So(status.String(), ShouldEqual, "Inactive")
		})
	})*/

	Convey("Given headers", t, func() {
		headers := generateExportHeaders()

		Convey("with 1 project with all his infos initialized", func() {

			Convey("When calling the generateXlsx function", func() {
			})
		})
	})
}
