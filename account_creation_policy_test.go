// +build all_tests
package stormpath_test

import (
	. "github.com/sappenin/stormpath-sdk-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AccountCreationPolicy", func() {

	Describe("GetVerificationEmailTemplates", func() {
		It("should return an error if the policy doesn't exists", func() {
			policy := &AccountCreationPolicy{VerificationEmailTemplates: &EmailTemplates{}}
			policy.VerificationEmailTemplates.Href = "https://api.stormpath.com/v1/accountCreationPolicies/xxxx/verificationEmailTemplates"

			templates, err := policy.GetVerificationEmailTemplates(ctx)

			Expect(err).To(HaveOccurred())
			Expect(templates).To(BeNil())
		})

		It("should return the default verification email templates collection", func() {
			directory := newTestDirectory()
			tenant.CreateDirectory(ctx, directory)

			policy, _ := directory.GetAccountCreationPolicy(ctx)

			templates, err := policy.GetVerificationEmailTemplates(ctx)

			Expect(err).NotTo(HaveOccurred())
			Expect(templates.Items).To(HaveLen(1))
		})
	})
	Describe("GetVerificationSuccessEmailTemplates", func() {
		It("should return an error if the policy doesn't exists", func() {
			policy := &AccountCreationPolicy{VerificationSuccessEmailTemplates: &EmailTemplates{}}
			policy.VerificationSuccessEmailTemplates.Href = "https://api.stormpath.com/v1/accountCreationPolicies/xxxx/verificationEmailTemplates"

			templates, err := policy.GetVerificationSuccessEmailTemplates(ctx)

			Expect(err).To(HaveOccurred())
			Expect(templates).To(BeNil())
		})

		It("should return the default verification success email templates collection", func() {
			directory := newTestDirectory()
			tenant.CreateDirectory(ctx, directory)

			policy, _ := directory.GetAccountCreationPolicy(ctx)

			templates, err := policy.GetVerificationSuccessEmailTemplates(ctx)

			Expect(err).NotTo(HaveOccurred())
			Expect(templates.Items).To(HaveLen(1))
		})
	})
	Describe("GetWelcomeEmailTemplates", func() {
		It("should return an error if the policy doesn't exists", func() {
			policy := &AccountCreationPolicy{WelcomeEmailTemplates: &EmailTemplates{}}
			policy.WelcomeEmailTemplates.Href = "https://api.stormpath.com/v1/accountCreationPolicies/xxxx/verificationEmailTemplates"

			templates, err := policy.GetWelcomeEmailTemplates(ctx)

			Expect(err).To(HaveOccurred())
			Expect(templates).To(BeNil())
		})

		It("should return the default welcome email templates collection", func() {
			directory := newTestDirectory()
			tenant.CreateDirectory(ctx, directory)

			policy, _ := directory.GetAccountCreationPolicy(ctx)

			templates, err := policy.GetWelcomeEmailTemplates(ctx)

			Expect(err).NotTo(HaveOccurred())
			Expect(templates.Items).To(HaveLen(1))
		})
	})
	Describe("Update", func() {
		It("should update a given account creation policy", func() {
			directory := newTestDirectory()
			tenant.CreateDirectory(ctx, directory)

			policy, _ := directory.GetAccountCreationPolicy(ctx)
			policy.VerificationEmailStatus = Enabled
			err := policy.Update(ctx)

			Expect(err).NotTo(HaveOccurred())
			Expect(policy.VerificationEmailStatus).To(Equal(Enabled))
		})
		It("should return error not found if the policy doesn't exists", func() {
			policy := AccountCreationPolicy{}
			policy.Href = BaseURL + "accountCreationPolicies/XXXX"

			err := policy.Update(ctx)

			Expect(err).To(HaveOccurred())
			Expect(err.(Error).Status).To(Equal(404))
		})
	})
	Describe("Refresh", func() {
		It("should refresh a given account creation policy", func() {
			directory := newTestDirectory()
			tenant.CreateDirectory(ctx, directory)

			policy, _ := directory.GetAccountCreationPolicy(ctx)
			policy.VerificationEmailStatus = Enabled
			err := policy.Refresh(ctx)

			Expect(err).NotTo(HaveOccurred())
			Expect(policy.VerificationEmailStatus).To(Equal(Disabled))
		})
		It("should return error not found if the policy doesn't exists", func() {
			policy := AccountCreationPolicy{}
			policy.Href = BaseURL + "accountCreationPolicies/XXXX"

			err := policy.Refresh(ctx)

			Expect(err).To(HaveOccurred())
			Expect(err.(Error).Status).To(Equal(404))
		})
	})
})
