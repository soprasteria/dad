package export

import (
	"sort"
	"strconv"
	"testing"

	"gopkg.in/mgo.v2/bson"

	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/soprasteria/dad/server/types"
	"github.com/tealeg/xlsx"
)

func TestGetServiceToIndicatorUsage(t *testing.T) {

	mapToSlice := func(usageIndicator map[string]types.UsageIndicator) []string {
		usageIndicatorSliced := make([]string, len(usageIndicator))

		for k := range usageIndicator {
			usageIndicatorSliced = append(usageIndicatorSliced, k)
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

	Convey("Given functional services with technical service 'Pipeline' initialized, with headers and 2 projects (1 initialized and 1 not initialized)", t, func() {

		technicalServicePipeline := []string{"jenkins"}
		services := []types.FunctionalService{
			{
				ID:       bson.ObjectId("1"),
				Name:     "service1",
				Package:  "pkg1",
				Position: 10},
			{
				ID:       bson.ObjectId("2"),
				Name:     "service2",
				Package:  "pkg1",
				Position: 20},
			{
				ID:       bson.ObjectId("3"),
				Name:     "service3",
				Package:  "pkg1",
				Position: 30},
			{
				ID:       bson.ObjectId("4"),
				Name:     "service4",
				Package:  "pkg2",
				Position: 10,
				Services: technicalServicePipeline},
			{
				ID:       bson.ObjectId("5"),
				Name:     "service5",
				Package:  "pkg2",
				Position: 20},
			{
				ID:       bson.ObjectId("6"),
				Name:     "service6",
				Package:  "pkg2",
				Position: 30},
			{
				ID:       bson.ObjectId("7"),
				Name:     "service7",
				Package:  "pkg2",
				Position: 40},
			{
				ID:       bson.ObjectId("8"),
				Name:     "service8",
				Package:  "pkg2",
				Position: 50},
			{
				ID:       bson.ObjectId("9"),
				Name:     "service9",
				Package:  "pkg2",
				Position: 60},
			{
				ID:       bson.ObjectId("10"),
				Name:     "service10",
				Package:  "pkg3",
				Position: 10},
			{
				ID:       bson.ObjectId("11"),
				Name:     "service11",
				Package:  "pkg3",
				Position: 20},
			{
				ID:       bson.ObjectId("12"),
				Name:     "service12",
				Package:  "pkg3",
				Position: 30},
			{
				ID:       bson.ObjectId("13"),
				Name:     "service13",
				Package:  "pkg3",
				Position: 40},
			{
				ID:       bson.ObjectId("14"),
				Name:     "service14",
				Package:  "pkg3",
				Position: 50},
			{
				ID:       bson.ObjectId("15"),
				Name:     "service15",
				Package:  "pkg3",
				Position: 60},
			{
				ID:       bson.ObjectId("16"),
				Name:     "service16",
				Package:  "pkg3",
				Position: 70},
			{
				ID:       bson.ObjectId("17"),
				Name:     "service17",
				Package:  "pkg3",
				Position: 80},
			{
				ID:       bson.ObjectId("18"),
				Name:     "service18",
				Package:  "pkg4",
				Position: 10},
			{
				ID:       bson.ObjectId("19"),
				Name:     "service19",
				Package:  "pkg4",
				Position: 20},
			{
				ID:       bson.ObjectId("20"),
				Name:     "service20",
				Package:  "pkg4",
				Position: 30},
			{
				ID:       bson.ObjectId("21"),
				Name:     "service21",
				Package:  "pkg4",
				Position: 40},
			{
				ID:       bson.ObjectId("22"),
				Name:     "service22",
				Package:  "pkg5",
				Position: 10},
			{
				ID:       bson.ObjectId("23"),
				Name:     "service23",
				Package:  "pkg5",
				Position: 20},
			{
				ID:       bson.ObjectId("24"),
				Name:     "service24",
				Package:  "pkg6",
				Position: 10},
			{
				ID:       bson.ObjectId("25"),
				Name:     "service25",
				Package:  "pkg6",
				Position: 20},
			{
				ID:       bson.ObjectId("26"),
				Name:     "service26",
				Package:  "pkg6",
				Position: 30}}

		servicesMap := make(map[string][]types.FunctionalService)
		for _, service := range services {
			servicesMap[service.Package] = append(servicesMap[service.Package], service)
		}

		headers := generateExportHeaders()

		dateCreated := time.Now()
		dateUpdated := time.Now()
		tmp := time.Now()
		dueDate := &tmp

		initializedProject := ProjectDataExport{}
		initializedProject.Name = "name"
		initializedProject.Description = "decription"
		initializedProject.BusinessUnit = "businessUnit"
		initializedProject.ServiceCenter = "serviceCenter"
		initializedProject.Domain = []string{"domain1", "domain2"}
		initializedProject.Client = "client"
		initializedProject.ProjectManager = "pm"
		initializedProject.Deputies = []string{"deputy1", "deputy2", "deputy3"}
		initializedProject.DocktorGroupName = "docktorGroupName"
		initializedProject.DocktorGroupURL = "docktorGroupURL"
		initializedProject.Technologies = []string{"Go", "JS"}
		initializedProject.Mode = "mode"
		initializedProject.VersionControlSystem = "versionControlSystem"
		initializedProject.DeliverablesInVersionControl = true
		initializedProject.SourceCodeInVersionControl = false
		initializedProject.SpecificationsInVersionControl = false
		initializedProject.Created = dateCreated
		initializedProject.Updated = dateUpdated

		initServiceDataExport := make(map[int][]ServiceDataExport)
		initServiceDataExport[0] = []ServiceDataExport{
			{
				Progress:  "00%",
				Goal:      "00%",
				Priority:  "P0",
				DueDate:   dueDate,
				Indicator: "N/A",
				Comment:   "service1"},
			{
				Progress:  "00%",
				Goal:      "00%",
				Priority:  "P0",
				DueDate:   dueDate,
				Indicator: "N/A",
				Comment:   "service2"},
			{
				Progress:  "00%",
				Goal:      "00%",
				Priority:  "P0",
				DueDate:   dueDate,
				Indicator: "N/A",
				Comment:   "service3"}}
		initServiceDataExport[1] = []ServiceDataExport{

			{
				Progress:  "20%",
				Goal:      "20%",
				Priority:  "P1",
				DueDate:   dueDate,
				Indicator: "Active",
				Comment:   "service4"},
			{
				Progress:  "20%",
				Goal:      "20%",
				Priority:  "P1",
				DueDate:   dueDate,
				Indicator: "N/A",
				Comment:   "service5"},
			{
				Progress:  "20%",
				Goal:      "20%",
				Priority:  "P1",
				DueDate:   dueDate,
				Indicator: "N/A",
				Comment:   "service6"},
			{
				Progress:  "20%",
				Goal:      "20%",
				Priority:  "P1",
				DueDate:   dueDate,
				Indicator: "N/A",
				Comment:   "service7"},
			{
				Progress:  "20%",
				Goal:      "20%",
				Priority:  "P1",
				DueDate:   dueDate,
				Indicator: "N/A",
				Comment:   "service8"},
			{
				Progress:  "20%",
				Goal:      "20%",
				Priority:  "P1",
				DueDate:   dueDate,
				Indicator: "N/A",
				Comment:   "service9"}}

		initServiceDataExport[2] = []ServiceDataExport{
			{
				Progress:  "40%",
				Goal:      "40%",
				Priority:  "P2",
				DueDate:   dueDate,
				Indicator: "N/A",
				Comment:   "service10"},
			{
				Progress:  "40%",
				Goal:      "40%",
				Priority:  "P2",
				DueDate:   dueDate,
				Indicator: "N/A",
				Comment:   "service11"},
			{
				Progress:  "40%",
				Goal:      "40%",
				Priority:  "P2",
				DueDate:   dueDate,
				Indicator: "N/A",
				Comment:   "service12"},
			{
				Progress:  "40%",
				Goal:      "40%",
				Priority:  "P2",
				DueDate:   dueDate,
				Indicator: "N/A",
				Comment:   "service13"},
			{
				Progress:  "40%",
				Goal:      "40%",
				Priority:  "P2",
				DueDate:   dueDate,
				Indicator: "N/A",
				Comment:   "service14"},
			{
				Progress:  "40%",
				Goal:      "40%",
				Priority:  "P2",
				DueDate:   dueDate,
				Indicator: "N/A",
				Comment:   "service15"},
			{
				Progress:  "40%",
				Goal:      "40%",
				Priority:  "P2",
				DueDate:   dueDate,
				Indicator: "N/A",
				Comment:   "service16"},
			{
				Progress:  "40%",
				Goal:      "40%",
				Priority:  "P2",
				DueDate:   dueDate,
				Indicator: "N/A",
				Comment:   "service17"}}

		initServiceDataExport[3] = []ServiceDataExport{
			{
				Progress:  "60%",
				Goal:      "60%",
				Priority:  "P0",
				DueDate:   dueDate,
				Indicator: "N/A",
				Comment:   "service18"},
			{
				Progress:  "60%",
				Goal:      "60%",
				Priority:  "P0",
				DueDate:   dueDate,
				Indicator: "N/A",
				Comment:   "service19"},
			{
				Progress:  "60%",
				Goal:      "60%",
				Priority:  "P0",
				DueDate:   dueDate,
				Indicator: "N/A",
				Comment:   "service20"},
			{
				Progress:  "60%",
				Goal:      "60%",
				Priority:  "P0",
				DueDate:   dueDate,
				Indicator: "N/A",
				Comment:   "service21"}}

		initServiceDataExport[4] = []ServiceDataExport{
			{
				Progress:  "80%",
				Goal:      "80%",
				Priority:  "P1",
				DueDate:   dueDate,
				Indicator: "N/A",
				Comment:   "service22"},
			{
				Progress:  "80%",
				Goal:      "80%",
				Priority:  "P1",
				DueDate:   dueDate,
				Indicator: "N/A",
				Comment:   "service23"}}

		initServiceDataExport[5] = []ServiceDataExport{
			{
				Progress:  "100%",
				Goal:      "100%",
				Priority:  "P2",
				DueDate:   dueDate,
				Indicator: "N/A",
				Comment:   "service24"},
			{
				Progress:  "100%",
				Goal:      "100%",
				Priority:  "P2",
				DueDate:   dueDate,
				Indicator: "N/A",
				Comment:   "service25"},
			{
				Progress:  "100%",
				Goal:      "100%",
				Priority:  "P2",
				DueDate:   dueDate,
				Indicator: "N/A",
				Comment:   "service26"}}

		initializedProject.ServicesDataExport = initServiceDataExport

		emptyProject := ProjectDataExport{}
		emptyProject.Name = "empty"
		emptyProject.DeliverablesInVersionControl = false
		emptyProject.SourceCodeInVersionControl = false
		emptyProject.SpecificationsInVersionControl = false
		emptyProject.Created = dateCreated
		emptyProject.Updated = dateUpdated

		emptyServiceDataExport := make(map[int][]ServiceDataExport)
		emptyServiceDataExport[0] = []ServiceDataExport{
			{
				Progress:  "N/A",
				Goal:      "N/A",
				Priority:  "N/A",
				DueDate:   nil,
				Indicator: "N/A"},
			{
				Progress:  "N/A",
				Goal:      "N/A",
				Priority:  "N/A",
				DueDate:   nil,
				Indicator: "N/A"},
			{
				Progress:  "N/A",
				Goal:      "N/A",
				Priority:  "N/A",
				DueDate:   nil,
				Indicator: "N/A"}}

		emptyServiceDataExport[1] = []ServiceDataExport{
			{
				Progress:  "N/A",
				Goal:      "N/A",
				Priority:  "N/A",
				DueDate:   nil,
				Indicator: "N/A"},
			{
				Progress:  "N/A",
				Goal:      "N/A",
				Priority:  "N/A",
				DueDate:   nil,
				Indicator: "N/A"},
			{
				Progress:  "N/A",
				Goal:      "N/A",
				Priority:  "N/A",
				DueDate:   nil,
				Indicator: "N/A"},
			{
				Progress:  "N/A",
				Goal:      "N/A",
				Priority:  "N/A",
				DueDate:   nil,
				Indicator: "N/A"},
			{
				Progress:  "N/A",
				Goal:      "N/A",
				Priority:  "N/A",
				DueDate:   nil,
				Indicator: "N/A"},
			{
				Progress:  "N/A",
				Goal:      "N/A",
				Priority:  "N/A",
				DueDate:   nil,
				Indicator: "N/A"}}

		emptyServiceDataExport[2] = []ServiceDataExport{
			{
				Progress:  "N/A",
				Goal:      "N/A",
				Priority:  "N/A",
				DueDate:   nil,
				Indicator: "N/A"},
			{
				Progress:  "N/A",
				Goal:      "N/A",
				Priority:  "N/A",
				DueDate:   nil,
				Indicator: "N/A"},
			{
				Progress:  "N/A",
				Goal:      "N/A",
				Priority:  "N/A",
				DueDate:   nil,
				Indicator: "N/A"},
			{
				Progress:  "N/A",
				Goal:      "N/A",
				Priority:  "N/A",
				DueDate:   nil,
				Indicator: "N/A"},
			{
				Progress:  "N/A",
				Goal:      "N/A",
				Priority:  "N/A",
				DueDate:   nil,
				Indicator: "N/A"},
			{
				Progress:  "N/A",
				Goal:      "N/A",
				Priority:  "N/A",
				DueDate:   nil,
				Indicator: "N/A"},
			{
				Progress:  "N/A",
				Goal:      "N/A",
				Priority:  "N/A",
				DueDate:   nil,
				Indicator: "N/A"},
			{
				Progress:  "N/A",
				Goal:      "N/A",
				Priority:  "N/A",
				DueDate:   nil,
				Indicator: "N/A"}}

		emptyServiceDataExport[3] = []ServiceDataExport{
			{
				Progress:  "N/A",
				Goal:      "N/A",
				Priority:  "N/A",
				DueDate:   nil,
				Indicator: "N/A"},
			{
				Progress:  "N/A",
				Goal:      "N/A",
				Priority:  "N/A",
				DueDate:   nil,
				Indicator: "N/A"},
			{
				Progress:  "N/A",
				Goal:      "N/A",
				Priority:  "N/A",
				DueDate:   nil,
				Indicator: "N/A"},
			{
				Progress:  "N/A",
				Goal:      "N/A",
				Priority:  "N/A",
				DueDate:   nil,
				Indicator: "N/A"}}

		emptyServiceDataExport[4] = []ServiceDataExport{
			{
				Progress:  "N/A",
				Goal:      "N/A",
				Priority:  "N/A",
				DueDate:   nil,
				Indicator: "N/A"},
			{
				Progress:  "N/A",
				Goal:      "N/A",
				Priority:  "N/A",
				DueDate:   nil,
				Indicator: "N/A"}}

		emptyServiceDataExport[5] = []ServiceDataExport{
			{
				Progress:  "N/A",
				Goal:      "N/A",
				Priority:  "N/A",
				DueDate:   nil,
				Indicator: "N/A"},
			{
				Progress:  "N/A",
				Goal:      "N/A",
				Priority:  "N/A",
				DueDate:   nil,
				Indicator: "N/A"},
			{
				Progress:  "N/A",
				Goal:      "N/A",
				Priority:  "N/A",
				DueDate:   nil,
				Indicator: "N/A"}}

		emptyProject.ServicesDataExport = emptyServiceDataExport

		Convey("with 1 project with his infos initialized", func() {

			dataExportWithProjectFullFilled := make(map[string]ProjectDataExport)
			dataExportWithProjectFullFilled["name"] = initializedProject

			Convey("When calling the generateXlsx function", func() {

				xlsxExport, _ := generateXlsx(servicesMap, headers, dataExportWithProjectFullFilled)

				Convey("Then the xlsx file should contains 4 rows (3 headers + 1 project)", func() {

					file, _ := xlsx.OpenReaderAt(xlsxExport, int64(xlsxExport.Len()))

					So(file, ShouldNotBeNil)
					sheet := file.Sheet["Deployment plan"]
					So(sheet, ShouldNotBeNil)
					So(sheet.Rows, ShouldHaveLength, 4)

					Convey("and all infos should correspond to those used to generate the xlsx file", func() {

						// project row
						cells := sheet.Rows[3].Cells
						So(cells, ShouldNotBeNil)
						So(cells, ShouldHaveLength, len(headers.MatrixHeaders)+len(services)*len(headers.ServicesInfosHeaders))

						// 18 cells for Matrix
						So(cells[0].Value, ShouldEqual, dataExportWithProjectFullFilled["name"].Name)
						So(cells[1].Value, ShouldEqual, dataExportWithProjectFullFilled["name"].Description)
						So(cells[2].Value, ShouldEqual, dataExportWithProjectFullFilled["name"].BusinessUnit)
						So(cells[3].Value, ShouldEqual, dataExportWithProjectFullFilled["name"].ServiceCenter)
						for _, domain := range dataExportWithProjectFullFilled["name"].Domain {
							So(cells[4].Value, ShouldContainSubstring, domain)
						}
						So(cells[5].Value, ShouldEqual, dataExportWithProjectFullFilled["name"].Client)
						So(cells[6].Value, ShouldEqual, dataExportWithProjectFullFilled["name"].ProjectManager)
						for _, deputy := range dataExportWithProjectFullFilled["name"].Deputies {
							So(cells[7].Value, ShouldContainSubstring, deputy)
						}
						So(cells[8].Value, ShouldEqual, dataExportWithProjectFullFilled["name"].DocktorGroupName)
						So(cells[9].Value, ShouldEqual, dataExportWithProjectFullFilled["name"].DocktorGroupURL)
						for _, technology := range dataExportWithProjectFullFilled["name"].Technologies {
							So(cells[10].Value, ShouldContainSubstring, technology)
						}
						So(cells[11].Value, ShouldEqual, dataExportWithProjectFullFilled["name"].Mode)
						So(cells[12].Value, ShouldEqual, dataExportWithProjectFullFilled["name"].VersionControlSystem)
						cellDeliverablesInVersionControl, _ := strconv.ParseBool(cells[13].Value)
						cellSourceCodeInVersionControl, _ := strconv.ParseBool(cells[14].Value)
						cellSpecificationsInVersionControl, _ := strconv.ParseBool(cells[15].Value)
						So(cellDeliverablesInVersionControl, ShouldEqual, dataExportWithProjectFullFilled["name"].DeliverablesInVersionControl)
						So(cellSourceCodeInVersionControl, ShouldEqual, dataExportWithProjectFullFilled["name"].SourceCodeInVersionControl)
						So(cellSpecificationsInVersionControl, ShouldEqual, dataExportWithProjectFullFilled["name"].SpecificationsInVersionControl)
						// can't verify dateCreated and dateUpdated (cell[16] and cell[17])

						servicesDataExportSortedKeys := make([]int, 0)
						for k := range dataExportWithProjectFullFilled["name"].ServicesDataExport {
							servicesDataExportSortedKeys = append(servicesDataExportSortedKeys, k)
						}
						sort.Ints(servicesDataExportSortedKeys)

						indexMatrix := len(headers.MatrixHeaders)
						for _, sortedKey := range servicesDataExportSortedKeys {

							serviceDataExport := dataExportWithProjectFullFilled["name"].ServicesDataExport[sortedKey]

							for _, service := range serviceDataExport {
								So(cells[indexMatrix].Value, ShouldEqual, service.Progress)
								So(cells[indexMatrix+1].Value, ShouldEqual, service.Goal)
								So(cells[indexMatrix+2].Value, ShouldEqual, service.Priority)
								if service.DueDate == nil {
									So(cells[indexMatrix+3].Value, ShouldEqual, "N/A")
								} else {
									// can't verify dueDate
								}
								So(cells[indexMatrix+4].Value, ShouldEqual, service.Indicator)
								So(cells[indexMatrix+5].Value, ShouldEqual, service.Comment)
								indexMatrix = indexMatrix + len(headers.ServicesInfosHeaders)
							}
						}
					})
				})
			})
		})

		Convey("with 1 project with his infos not initialized", func() {

			dataExportWithProjectEmpty := make(map[string]ProjectDataExport)
			dataExportWithProjectEmpty["empty"] = emptyProject

			Convey("When calling the generateXlsx function", func() {

				xlsxExport, _ := generateXlsx(servicesMap, headers, dataExportWithProjectEmpty)

				Convey("Then the xlsx file should contains 4 rows (3 headers + 1 project)", func() {

					file, _ := xlsx.OpenReaderAt(xlsxExport, int64(xlsxExport.Len()))

					So(file, ShouldNotBeNil)
					sheet := file.Sheet["Deployment plan"]
					So(sheet, ShouldNotBeNil)
					So(sheet.Rows, ShouldHaveLength, 4)

					Convey("and all infos should correspond to those used to generate the xlsx file", func() {

						// project row
						cells := sheet.Rows[3].Cells
						So(cells, ShouldNotBeNil)
						So(cells, ShouldHaveLength, len(headers.MatrixHeaders)+len(services)*len(headers.ServicesInfosHeaders))

						// 18 cells for Matrix
						So(cells[0].Value, ShouldEqual, dataExportWithProjectEmpty["empty"].Name)
						So(cells[1].Value, ShouldBeBlank)
						So(cells[2].Value, ShouldBeBlank)
						So(cells[3].Value, ShouldBeBlank)
						So(cells[4].Value, ShouldBeBlank)
						So(cells[5].Value, ShouldBeBlank)
						So(cells[6].Value, ShouldBeBlank)
						So(cells[7].Value, ShouldBeBlank)
						So(cells[8].Value, ShouldBeBlank)
						So(cells[9].Value, ShouldBeBlank)
						So(cells[10].Value, ShouldBeBlank)
						So(cells[11].Value, ShouldBeBlank)
						So(cells[12].Value, ShouldBeBlank)
						cellDeliverablesInVersionControl, _ := strconv.ParseBool(cells[13].Value)
						cellSourceCodeInVersionControl, _ := strconv.ParseBool(cells[14].Value)
						cellSpecificationsInVersionControl, _ := strconv.ParseBool(cells[15].Value)
						So(cellDeliverablesInVersionControl, ShouldEqual, dataExportWithProjectEmpty["empty"].DeliverablesInVersionControl)
						So(cellSourceCodeInVersionControl, ShouldEqual, dataExportWithProjectEmpty["empty"].SourceCodeInVersionControl)
						So(cellSpecificationsInVersionControl, ShouldEqual, dataExportWithProjectEmpty["empty"].SpecificationsInVersionControl)
						// can't verify dateCreated and dateUpdated (cell[16] and cell[17])

						servicesDataExportSortedKeys := make([]int, 0)
						for k := range dataExportWithProjectEmpty["empty"].ServicesDataExport {
							servicesDataExportSortedKeys = append(servicesDataExportSortedKeys, k)
						}
						sort.Ints(servicesDataExportSortedKeys)

						indexMatrix := len(headers.MatrixHeaders)
						for _, sortedKey := range servicesDataExportSortedKeys {

							serviceDataExport := dataExportWithProjectEmpty["empty"].ServicesDataExport[sortedKey]

							for _, service := range serviceDataExport {
								So(service, ShouldNotBeNil)
								So(cells[indexMatrix].Value, ShouldEqual, "N/A")
								So(cells[indexMatrix+1].Value, ShouldEqual, "N/A")
								So(cells[indexMatrix+2].Value, ShouldEqual, "N/A")
								So(cells[indexMatrix+3].Value, ShouldEqual, "N/A")
								So(cells[indexMatrix+4].Value, ShouldEqual, "N/A")
								So(cells[indexMatrix+5].Value, ShouldBeBlank)
								indexMatrix = indexMatrix + len(headers.ServicesInfosHeaders)
							}
						}
					})
				})
			})
		})

		Convey("with 2 project (1 initialized and 1 not initialized)", func() {

			dataExport := make(map[string]ProjectDataExport)
			dataExport["name"] = initializedProject
			dataExport["empty"] = emptyProject

			Convey("When calling the generateXlsx function", func() {

				xlsxExport, _ := generateXlsx(servicesMap, headers, dataExport)

				Convey("Then the xlsx file should contains 4 rows (3 headers + 2 projects)", func() {

					file, _ := xlsx.OpenReaderAt(xlsxExport, int64(xlsxExport.Len()))

					So(file, ShouldNotBeNil)
					sheet := file.Sheet["Deployment plan"]
					So(sheet, ShouldNotBeNil)
					So(sheet.Rows, ShouldHaveLength, 5)

					Convey("and all infos should correspond to those used to generate the xlsx file", func() {

						// project empty row
						cells := sheet.Rows[3].Cells
						So(cells, ShouldNotBeNil)
						So(cells, ShouldHaveLength, len(headers.MatrixHeaders)+len(services)*len(headers.ServicesInfosHeaders))

						// 18 cells for Matrix
						So(cells[0].Value, ShouldEqual, dataExport["empty"].Name)
						So(cells[1].Value, ShouldBeBlank)
						So(cells[2].Value, ShouldBeBlank)
						So(cells[3].Value, ShouldBeBlank)
						So(cells[4].Value, ShouldBeBlank)
						So(cells[5].Value, ShouldBeBlank)
						So(cells[6].Value, ShouldBeBlank)
						So(cells[7].Value, ShouldBeBlank)
						So(cells[8].Value, ShouldBeBlank)
						So(cells[9].Value, ShouldBeBlank)
						So(cells[10].Value, ShouldBeBlank)
						So(cells[11].Value, ShouldBeBlank)
						So(cells[12].Value, ShouldBeBlank)
						cellDeliverablesInVersionControl, _ := strconv.ParseBool(cells[13].Value)
						cellSourceCodeInVersionControl, _ := strconv.ParseBool(cells[14].Value)
						cellSpecificationsInVersionControl, _ := strconv.ParseBool(cells[15].Value)
						So(cellDeliverablesInVersionControl, ShouldEqual, dataExport["empty"].DeliverablesInVersionControl)
						So(cellSourceCodeInVersionControl, ShouldEqual, dataExport["empty"].SourceCodeInVersionControl)
						So(cellSpecificationsInVersionControl, ShouldEqual, dataExport["empty"].SpecificationsInVersionControl)
						// can't verify dateCreated and dateUpdated (cell[16] and cell[17])

						servicesDataExportSortedKeys := make([]int, 0)
						for k := range dataExport["empty"].ServicesDataExport {
							servicesDataExportSortedKeys = append(servicesDataExportSortedKeys, k)
						}
						sort.Ints(servicesDataExportSortedKeys)

						indexMatrix := len(headers.MatrixHeaders)
						for _, sortedKey := range servicesDataExportSortedKeys {

							serviceDataExport := dataExport["empty"].ServicesDataExport[sortedKey]

							for _, service := range serviceDataExport {
								So(service, ShouldNotBeNil)
								So(cells[indexMatrix].Value, ShouldEqual, "N/A")
								So(cells[indexMatrix+1].Value, ShouldEqual, "N/A")
								So(cells[indexMatrix+2].Value, ShouldEqual, "N/A")
								So(cells[indexMatrix+3].Value, ShouldEqual, "N/A")
								So(cells[indexMatrix+4].Value, ShouldEqual, "N/A")
								So(cells[indexMatrix+5].Value, ShouldBeBlank)
								indexMatrix = indexMatrix + len(headers.ServicesInfosHeaders)
							}
						}

						// project name row
						cells = sheet.Rows[4].Cells
						So(cells, ShouldNotBeNil)
						So(cells, ShouldHaveLength, len(headers.MatrixHeaders)+len(services)*len(headers.ServicesInfosHeaders))

						// 18 cells for Matrix
						So(cells[0].Value, ShouldEqual, dataExport["name"].Name)
						So(cells[1].Value, ShouldEqual, dataExport["name"].Description)
						So(cells[2].Value, ShouldEqual, dataExport["name"].BusinessUnit)
						So(cells[3].Value, ShouldEqual, dataExport["name"].ServiceCenter)
						for _, domain := range dataExport["name"].Domain {
							So(cells[4].Value, ShouldContainSubstring, domain)
						}
						So(cells[5].Value, ShouldEqual, dataExport["name"].Client)
						So(cells[6].Value, ShouldEqual, dataExport["name"].ProjectManager)
						for _, deputy := range dataExport["name"].Deputies {
							So(cells[7].Value, ShouldContainSubstring, deputy)
						}
						So(cells[8].Value, ShouldEqual, dataExport["name"].DocktorGroupName)
						So(cells[9].Value, ShouldEqual, dataExport["name"].DocktorGroupURL)
						for _, technology := range dataExport["name"].Technologies {
							So(cells[10].Value, ShouldContainSubstring, technology)
						}
						So(cells[11].Value, ShouldEqual, dataExport["name"].Mode)
						So(cells[12].Value, ShouldEqual, dataExport["name"].VersionControlSystem)
						cellDeliverablesInVersionControl, _ = strconv.ParseBool(cells[13].Value)
						cellSourceCodeInVersionControl, _ = strconv.ParseBool(cells[14].Value)
						cellSpecificationsInVersionControl, _ = strconv.ParseBool(cells[15].Value)
						So(cellDeliverablesInVersionControl, ShouldEqual, dataExport["name"].DeliverablesInVersionControl)
						So(cellSourceCodeInVersionControl, ShouldEqual, dataExport["name"].SourceCodeInVersionControl)
						So(cellSpecificationsInVersionControl, ShouldEqual, dataExport["name"].SpecificationsInVersionControl)
						// can't verify dateCreated and dateUpdated (cell[16] and cell[17])

						servicesDataExportSortedKeys = make([]int, 0)
						for k := range dataExport["name"].ServicesDataExport {
							servicesDataExportSortedKeys = append(servicesDataExportSortedKeys, k)
						}
						sort.Ints(servicesDataExportSortedKeys)

						indexMatrix = len(headers.MatrixHeaders)
						for _, sortedKey := range servicesDataExportSortedKeys {

							serviceDataExport := dataExport["name"].ServicesDataExport[sortedKey]

							for _, service := range serviceDataExport {
								So(cells[indexMatrix].Value, ShouldEqual, service.Progress)
								So(cells[indexMatrix+1].Value, ShouldEqual, service.Goal)
								So(cells[indexMatrix+2].Value, ShouldEqual, service.Priority)
								if service.DueDate == nil {
									So(cells[indexMatrix+3].Value, ShouldEqual, "N/A")
								} else {
									// can't verify dueDate
								}
								So(cells[indexMatrix+4].Value, ShouldEqual, service.Indicator)
								So(cells[indexMatrix+5].Value, ShouldEqual, service.Comment)
								indexMatrix = indexMatrix + len(headers.ServicesInfosHeaders)
							}
						}
					})
				})
			})
		})
	})
}
