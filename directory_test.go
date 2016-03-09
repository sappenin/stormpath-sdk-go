// +build all_tests
package stormpath_test

import (
	"encoding/json"

	. "github.com/sappenin/stormpath-sdk-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Directory", func() {
	Describe("JSON", func() {
		It("should marshal a minimum JSON with only the directory name", func() {
			directory := NewDirectory("name")

			jsonData, _ := json.Marshal(directory)

			Expect(string(jsonData)).To(Equal("{\"name\":\"name\"}"))
		})
	})

	Describe("Delete", func() {
		It("should delete an existing directory", func() {
			directory := newTestDirectory()

			tenant.CreateDirectory(ctx, directory)
			err := directory.Delete(ctx)

			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("GetAccountCreationPolicy", func() {
		It("should retrive the directory account creation policy", func() {
			directory := newTestDirectory()
			tenant.CreateDirectory(ctx, directory)

			policy, err := directory.GetAccountCreationPolicy(ctx)

			Expect(err).NotTo(HaveOccurred())
			Expect(policy).To(Equal(directory.AccountCreationPolicy))
			Expect(policy.VerificationEmailStatus).To(Equal("DISABLED"))
			Expect(policy.VerificationSuccessEmailStatus).To(Equal("DISABLED"))
			Expect(policy.WelcomeEmailStatus).To(Equal("DISABLED"))
		})
	})

	Describe("GetGroups", func() {
		It("should retrive all directory groups", func() {
			directory := newTestDirectory()
			tenant.CreateDirectory(ctx, directory)

			groups, err := directory.GetGroups(ctx, MakeGroupCriteria())

			Expect(err).NotTo(HaveOccurred())
			Expect(groups.Href).NotTo(BeEmpty())
			Expect(groups.Offset).To(Equal(0))
			Expect(groups.Limit).To(Equal(25))
			Expect(groups.Items).To(BeEmpty())
			directory.Delete(ctx)
		})
	})

	Describe("GetAccounts", func() {
		It("should retrieve all directory accounts", func() {
			directory := newTestDirectory()
			tenant.CreateDirectory(ctx, directory)

			accounts, err := directory.GetAccounts(ctx, MakeAccountCriteria())

			Expect(err).NotTo(HaveOccurred())
			Expect(accounts.Href).NotTo(BeEmpty())
			Expect(accounts.Offset).To(Equal(0))
			Expect(accounts.Limit).To(Equal(25))
			Expect(accounts.Items).To(BeEmpty())
			directory.Delete(ctx)
		})
	})

	Describe("CreateGroup", func() {
		It("should create new group", func() {
			directory := newTestDirectory()
			tenant.CreateDirectory(ctx, directory)

			group := NewGroup("new-group")
			err := directory.CreateGroup(ctx, group)

			Expect(err).NotTo(HaveOccurred())
			Expect(group.Href).NotTo(BeEmpty())
			directory.Delete(ctx)
		})
	})

	Describe("RegisterAccount", func() {
		It("should create a new accout for the group", func() {
			directory := newTestDirectory()
			tenant.CreateDirectory(ctx, directory)

			account := newTestAccount()
			err := directory.RegisterAccount(ctx, account)
			Expect(err).NotTo(HaveOccurred())
			Expect(account.Href).NotTo(BeEmpty())
			directory.Delete(ctx)
		})
	})
})
