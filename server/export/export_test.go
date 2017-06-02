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
			usageIndicatorSliced[i] = k
			i++
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
