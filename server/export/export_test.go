package export

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/soprasteria/dad/server/types"
)

func TestFunctionBestIndicatorStatus(t *testing.T) {

	Convey("Given 0 technical service and 0 indicator", t, func() {
		Convey("When calling the bestIndicatorStatus function", func() {
			Convey("Then the result is Empty", func() {

				// Given
				services := []string{}
				servicesToMatch := []types.UsageIndicator{}

				// When
				status := bestIndicatorStatus(services, servicesToMatch)

				// Then
				So(status, ShouldEqual, "Empty")
			})
		})
	})

	Convey("Given 0 technical service and 1 indicator", t, func() {
		Convey("When calling the bestIndicatorStatus function", func() {
			Convey("Then the result is Empty", func() {

				// Given
				services := []string{}
				servicesToMatch := []types.UsageIndicator{{Service: "jenkins"}}

				// When
				status := bestIndicatorStatus(services, servicesToMatch)

				// Then
				So(status, ShouldEqual, "Empty")
			})
		})
	})

	Convey("Given 1 technical service and 0 indicator", t, func() {
		Convey("When calling the bestIndicatorStatus function", func() {
			Convey("Then the result is Empty", func() {

				// Given
				services := []string{"jenkins"}
				servicesToMatch := []types.UsageIndicator{}

				// When
				status := bestIndicatorStatus(services, servicesToMatch)

				// Then
				So(status, ShouldEqual, "Empty")
			})
		})
	})

	Convey("Given 1 technical service and 1 indicator", t, func() {
		Convey("with the indicator service doesn't match the service", func() {
			Convey("When calling the bestIndicatorStatus function", func() {
				Convey("Then the result is Empty", func() {

					// Given
					services := []string{"jenkins"}
					servicesToMatch := []types.UsageIndicator{{Service: "gitlabci"}}

					// When
					status := bestIndicatorStatus(services, servicesToMatch)

					// Then
					So(status, ShouldEqual, "Empty")
				})
			})
		})

		Convey("with the indicator service matching the service", func() {
			Convey("and the indicator equal to Empty", func() {
				Convey("When calling the bestIndicatorStatus function", func() {
					Convey("Then the result is Empty", func() {

						// Given
						services := []string{"jenkins"}
						servicesToMatch := []types.UsageIndicator{{Service: "jenkins", Status: "Empty"}}

						// When
						status := bestIndicatorStatus(services, servicesToMatch)

						// Then
						So(status, ShouldEqual, "Empty")
					})
				})
			})

			Convey("and the indicator equal to Undetermined", func() {
				Convey("When calling the bestIndicatorStatus function", func() {
					Convey("Then the result is Undetermined", func() {

						// Given
						services := []string{"jenkins"}
						servicesToMatch := []types.UsageIndicator{{Service: "jenkins", Status: "Undetermined"}}

						// When
						status := bestIndicatorStatus(services, servicesToMatch)

						// Then
						So(status, ShouldEqual, "Undetermined")
					})
				})
			})

			Convey("and the indicator equal to Inactive", func() {
				Convey("When calling the bestIndicatorStatus function", func() {
					Convey("Then the result is Inactive", func() {

						// Given
						services := []string{"jenkins"}
						servicesToMatch := []types.UsageIndicator{{Service: "jenkins", Status: "Inactive"}}

						// When
						status := bestIndicatorStatus(services, servicesToMatch)

						// Then
						So(status, ShouldEqual, "Inactive")
					})
				})
			})

			Convey("and the indicator equal to Active", func() {
				Convey("When calling the bestIndicatorStatus function", func() {
					Convey("Then the result is Active", func() {

						// Given
						services := []string{"jenkins"}
						servicesToMatch := []types.UsageIndicator{{Service: "jenkins", Status: "Active"}}

						// When
						status := bestIndicatorStatus(services, servicesToMatch)

						// Then
						So(status, ShouldEqual, "Active")
					})
				})
			})
		})
	})
	Convey("Given 2 technical services and 2 indicators", t, func() {
		Convey("with indicator services matching services", func() {
			Convey("and both indicators equal to Empty", func() {
				Convey("When calling the bestIndicatorStatus function", func() {
					Convey("Then the result is Empty", func() {

						// Given
						services := []string{"jenkins", "gitlabci"}
						servicesToMatch := []types.UsageIndicator{{Service: "jenkins", Status: "Empty"}, {Service: "gitlabci", Status: "Empty"}}

						// When
						status := bestIndicatorStatus(services, servicesToMatch)

						// Then
						So(status, ShouldEqual, "Empty")
					})
				})
			})

			Convey("and one indicator equal to Empty and the other equal to Undetermined", func() {
				Convey("When calling the bestIndicatorStatus function", func() {
					Convey("Then the result is Undetermined", func() {

						// Given
						services := []string{"jenkins", "gitlabci"}
						servicesToMatch := []types.UsageIndicator{{Service: "jenkins", Status: "Empty"}, {Service: "gitlabci", Status: "Undetermined"}}

						// When
						status := bestIndicatorStatus(services, servicesToMatch)

						// Then
						So(status, ShouldEqual, "Undetermined")
					})
				})
			})

			Convey("and one indicator equal to Empty and the other equal to Inactive", func() {
				Convey("When calling the bestIndicatorStatus function", func() {
					Convey("Then the result is Inactive", func() {

						// Given
						services := []string{"jenkins", "gitlabci"}
						servicesToMatch := []types.UsageIndicator{{Service: "jenkins", Status: "Empty"}, {Service: "gitlabci", Status: "Inactive"}}

						// When
						status := bestIndicatorStatus(services, servicesToMatch)

						// Then
						So(status, ShouldEqual, "Inactive")
					})
				})
			})

			Convey("and one indicator equal to Empty and the other equal to Active", func() {
				Convey("When calling the bestIndicatorStatus function", func() {
					Convey("Then the result is Active", func() {

						// Given
						services := []string{"jenkins", "gitlabci"}
						servicesToMatch := []types.UsageIndicator{{Service: "jenkins", Status: "Empty"}, {Service: "gitlabci", Status: "Active"}}

						// When
						status := bestIndicatorStatus(services, servicesToMatch)

						// Then
						So(status, ShouldEqual, "Active")
					})
				})
			})

			Convey("and both indicators equal to Undetermined", func() {
				Convey("When calling the bestIndicatorStatus function", func() {
					Convey("Then the result is Undetermined", func() {

						// Given
						services := []string{"jenkins", "gitlabci"}
						servicesToMatch := []types.UsageIndicator{{Service: "jenkins", Status: "Undetermined"}, {Service: "gitlabci", Status: "Undetermined"}}

						// When
						status := bestIndicatorStatus(services, servicesToMatch)

						// Then
						So(status, ShouldEqual, "Undetermined")
					})
				})
			})

			Convey("and one indicator equal to Undetermined and the other equal to Inactive", func() {
				Convey("When calling the bestIndicatorStatus function", func() {
					Convey("Then the result is Inactive", func() {

						// Given
						services := []string{"jenkins", "gitlabci"}
						servicesToMatch := []types.UsageIndicator{{Service: "jenkins", Status: "Undetermined"}, {Service: "gitlabci", Status: "Inactive"}}

						// When
						status := bestIndicatorStatus(services, servicesToMatch)

						// Then
						So(status, ShouldEqual, "Inactive")
					})
				})
			})

			Convey("and one indicator equal to Undetermined and the other equal to Active", func() {
				Convey("When calling the bestIndicatorStatus function", func() {
					Convey("Then the result is Active", func() {

						// Given
						services := []string{"jenkins", "gitlabci"}
						servicesToMatch := []types.UsageIndicator{{Service: "jenkins", Status: "Undetermined"}, {Service: "gitlabci", Status: "Active"}}

						// When
						status := bestIndicatorStatus(services, servicesToMatch)

						// Then
						So(status, ShouldEqual, "Active")
					})
				})
			})

			Convey("and both indicators equal to Inactive", func() {
				Convey("When calling the bestIndicatorStatus function", func() {
					Convey("Then the result is Inactive", func() {

						// Given
						services := []string{"jenkins", "gitlabci"}
						servicesToMatch := []types.UsageIndicator{{Service: "jenkins", Status: "Inactive"}, {Service: "gitlabci", Status: "Inactive"}}

						// When
						status := bestIndicatorStatus(services, servicesToMatch)

						// Then
						So(status, ShouldEqual, "Inactive")
					})
				})
			})

			Convey("and one indicator equal to Inactive and the other equal to Active", func() {
				Convey("When calling the bestIndicatorStatus function", func() {
					Convey("Then the result is Active", func() {

						// Given
						services := []string{"jenkins", "gitlabci"}
						servicesToMatch := []types.UsageIndicator{{Service: "jenkins", Status: "Inactive"}, {Service: "gitlabci", Status: "Active"}}

						// When
						status := bestIndicatorStatus(services, servicesToMatch)

						// Then
						So(status, ShouldEqual, "Active")
					})
				})
			})

			Convey("and both indicators equal to Active", func() {
				Convey("When calling the bestIndicatorStatus function", func() {
					Convey("Then the result is Active", func() {

						// Given
						services := []string{"jenkins", "gitlabci"}
						servicesToMatch := []types.UsageIndicator{{Service: "jenkins", Status: "Active"}, {Service: "gitlabci", Status: "Active"}}

						// When
						status := bestIndicatorStatus(services, servicesToMatch)

						// Then
						So(status, ShouldEqual, "Active")
					})
				})
			})
		})
	})
}
