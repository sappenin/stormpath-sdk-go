package stormpath_test

import (
	. "github.com/sappenin/stormpath-sdk-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tenant", func() {
	Describe("CurrentTentant", func() {
		It("should retrive the current tenant", func() {
			tenant, err := CurrentTenant(ctx)

			Expect(err).NotTo(HaveOccurred())
			Expect(tenant.Href).NotTo(BeEmpty())
			Expect(tenant.Name).NotTo(BeEmpty())
			Expect(tenant.Key).NotTo(BeEmpty())
			Expect(tenant.Applications.Href).NotTo(BeEmpty())
			Expect(tenant.Directories.Href).NotTo(BeEmpty())
		})
	})

	Describe("CreateApplication", func() {
		It("should create a new application", func() {
			application := newTestApplication()
			err := tenant.CreateApplication(ctx, application)
			application.Purge(ctx)

			Expect(err).NotTo(HaveOccurred())
			Expect(application.Href).NotTo(BeEmpty())
		})
	})

	Describe("Custom Data", func() {
		It("UpdateCustomData should update the given tenant custom data", func() {
			customData := map[string]interface{}{
				"testIntField":    1,
				"testStringField": "test",
			}

			updatedCustomData, err := tenant.UpdateCustomData(ctx, customData)

			Expect(err).NotTo(HaveOccurred())
			Expect(updatedCustomData["testIntField"]).To(Equal(float64(1)))
			tenant.DeleteCustomData(ctx)
		})

		It("GetCustomData should update the given tenant custom data", func() {
			customData := map[string]interface{}{
				"testIntField":    1,
				"testStringField": "test",
			}

			tenant.UpdateCustomData(ctx, customData)

			customData, err := tenant.GetCustomData(ctx)

			Expect(err).NotTo(HaveOccurred())
			Expect(customData["testIntField"]).To(Equal(float64(1)))
		})

		It("DeleteCustomData should update the given tenant custom data", func() {
			customData := map[string]interface{}{
				"testIntField":    1,
				"testStringField": "test",
			}

			tenant.UpdateCustomData(ctx, customData)
			err := tenant.DeleteCustomData(ctx)

			customData, _ = tenant.GetCustomData(ctx)

			Expect(err).NotTo(HaveOccurred())
			Expect(customData).To(HaveLen(3))
		})

		//Describe("concurrent access", func() {
		//	It("should allow current access and be consistent at the end", func() {
		//		for i := 0; i < 8; i++ {
		//			go func() {
		//				defer GinkgoRecover()
		//
		//				customData := map[string]interface{}{
		//					"testIntField":    i,
		//					"testStringField": "test",
		//				}
		//
		//				data, err := tenant.UpdateCustomData(customData)
		//				Expect(err).NotTo(HaveOccurred())
		//				Expect(data["testIntField"]).To(Equal(float64(i)))
		//				tenant.DeleteCustomData()
		//			}()
		//		}
		//	})
		//})
	})

	Describe("CreateDirectory", func() {
		It("should create a new directory", func() {
			dir := newTestDirectory()
			err := tenant.CreateDirectory(ctx, dir)
			dir.Delete(ctx)

			Expect(err).NotTo(HaveOccurred())
			Expect(dir.Href).NotTo(BeEmpty())
		})
	})

	Describe("GetDirectories", func() {
		It("should retrive all the tenant directories", func() {
			tenant, _ := CurrentTenant(ctx)

			directories, err := tenant.GetDirectories(ctx, MakeDirectoriesCriteria())

			Expect(err).NotTo(HaveOccurred())
			Expect(directories.Href).NotTo(BeEmpty())
			Expect(directories.Offset).To(Equal(0))
			Expect(directories.Limit).To(Equal(25))
			Expect(directories.Items).NotTo(BeEmpty())
		})

		It("should retrive all the tenant directories by page", func() {
			tenant, _ := CurrentTenant(ctx)

			directories, err := tenant.GetDirectories(ctx, MakeDirectoriesCriteria().Limit(1))

			Expect(err).NotTo(HaveOccurred())
			Expect(directories.Href).NotTo(BeEmpty())
			Expect(directories.Offset).To(Equal(0))
			Expect(directories.Limit).To(Equal(1))
			Expect(directories.Items).To(HaveLen(1))
		})

		It("should retrive all the tenant directories by page and filter", func() {
			tenant, _ := CurrentTenant(ctx)

			directories, err := tenant.GetDirectories(ctx, MakeDirectoriesCriteria().NameEq("Stormpath Administrators"))

			Expect(err).NotTo(HaveOccurred())
			Expect(directories.Href).NotTo(BeEmpty())
			Expect(directories.Items).To(HaveLen(1))
		})

	})

	Describe("GetApplications", func() {
		It("should retrive all the tenant applications", func() {
			tenant, _ := CurrentTenant(ctx)

			apps, err := tenant.GetApplications(ctx, MakeApplicationCriteria())

			Expect(err).NotTo(HaveOccurred())
			Expect(apps.Href).NotTo(BeEmpty())
			Expect(apps.Offset).To(Equal(0))
			Expect(apps.Limit).To(Equal(25))
			Expect(apps.Items).NotTo(BeEmpty())
		})

		It("should retrive all the tenant applications by page", func() {
			tenant, _ := CurrentTenant(ctx)

			apps, err := tenant.GetApplications(ctx, MakeApplicationCriteria().Offset(0).Limit(1))

			Expect(err).NotTo(HaveOccurred())
			Expect(apps.Href).NotTo(BeEmpty())
			Expect(apps.Offset).To(Equal(0))
			Expect(apps.Limit).To(Equal(1))
			Expect(apps.Items).To(HaveLen(1))
		})

		It("should retrive all the tenant applications by page and filter", func() {
			tenant, _ := CurrentTenant(ctx)

			apps, err := tenant.GetApplications(ctx, MakeApplicationCriteria().NameEq("stormpath"))

			Expect(err).NotTo(HaveOccurred())
			Expect(apps.Href).NotTo(BeEmpty())
			Expect(apps.Items).To(HaveLen(1))
		})
	})
})
