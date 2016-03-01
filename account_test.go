package stormpath_test

import (
	"encoding/json"

	. "github.com/sappenin/stormpath-sdk-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Account", func() {
	Describe("JSON", func() {
		It("should marshal a minimum JSON with only the account required fields", func() {
			acc := NewAccount("test@test.org", "123", "test@test.org", "test", "test")

			jsonData, _ := json.Marshal(acc)

			Expect(string(jsonData)).To(Equal("{\"username\":\"test@test.org\",\"email\":\"test@test.org\",\"password\":\"123\",\"givenName\":\"test\",\"surname\":\"test\"}"))
		})
	})
	Describe("GetAccount", func() {
		It("should return an error if the account doesn't exists", func() {
			account, err := GetAccount(BaseURL+"/accounts/xxxxxx", MakeAccountCriteria())

			Expect(err).To(HaveOccurred())
			Expect(err.(Error).Status).To(Equal(404))
			Expect(account).To(BeNil())
		})
		It("should return the account for the given href", func() {
			newAccount := newTestAccount()
			app.RegisterAccount(newAccount)

			account, err := GetAccount(newAccount.Href, MakeAccountCriteria())

			Expect(err).NotTo(HaveOccurred())
			Expect(account).To(Equal(newAccount))
		})
	})
	Describe("VerifyEmailToken", func() {
		It("should return error if the token doesn't exists", func() {
			account, err := VerifyEmailToken("token")
			Expect(err).To(HaveOccurred())
			Expect(err.(Error).Status).To(Equal(404))
			Expect(account).To(BeNil())
		})
		It("should return an account if the token is valid", func() {
			directory := newTestDirectory()
			tenant.CreateDirectory(directory)

			policy, _ := directory.GetAccountCreationPolicy()
			policy.VerificationEmailStatus = Enabled
			policy.Update()

			account := newTestAccount()
			directory.RegisterAccount(account)
			a, err := VerifyEmailToken(GetToken(account.EmailVerificationToken.Href))

			Expect(err).NotTo(HaveOccurred())
			Expect(a.Href).To(Equal(account.Href))
		})
	})
	Describe("Update", func() {
		It("should update an existing account", func() {
			account := newTestAccount()
			app.RegisterAccount(account)

			account.GivenName = "julio"
			err := account.Update()

			Expect(err).NotTo(HaveOccurred())
			Expect(account.GivenName).To(Equal("julio"))
			Expect(account.CreatedAt).NotTo(BeNil())
			Expect(account.ModifiedAt).NotTo(BeNil())
		})
	})

	Describe("Delete", func() {
		It("should delete an existing account", func() {
			account := newTestAccount()
			app.RegisterAccount(account)

			err := account.Delete()

			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("AddToGroup", func() {
		It("should add an account to an existing group", func() {
			group := newTestGroup()
			app.CreateGroup(group)

			_, err := account.AddToGroup(group)
			gm, _ := account.GetGroupMemberships(MakeAccountCriteria().Offset(0).Limit(25))

			Expect(err).NotTo(HaveOccurred())
			Expect(gm.Items).To(HaveLen(1))
			account.RemoveFromGroup(group)
			group.Delete()
		})
	})

	Describe("RemoveFromGroup", func() {
		It("should remove an account from an existing group", func() {
			account := newTestAccount()
			app.RegisterAccount(account)

			var groupCountBefore int
			group := newTestGroup()
			app.CreateGroup(group)

			gm, _ := account.GetGroupMemberships(MakeAccountCriteria().Offset(0).Limit(25))
			groupCountBefore = len(gm.Items)

			account.AddToGroup(group)

			err := account.RemoveFromGroup(group)
			gm, _ = account.GetGroupMemberships(MakeAccountCriteria().Offset(0).Limit(25))

			Expect(err).NotTo(HaveOccurred())
			Expect(gm.Items).To(HaveLen(groupCountBefore))
			group.Delete()
		})
	})

	Describe("GetGroupMemberships", func() {
		It("should allow expanding the account", func() {
			acct := registerTestAccount()
			group := addAccountToGroup(acct)

			groupMemberships, err := acct.GetGroupMemberships(MakeGroupMemershipCriteria().WithAccount().Offset(0).Limit(25))

			Expect(err).NotTo(HaveOccurred())
			for _, gm := range groupMemberships.Items {
				Expect(gm.Account).To(BeEquivalentTo(*acct))
				Expect(gm.Group).NotTo(BeEquivalentTo(*group))
			}
		})
	})

	Describe("GetCustomData", func() {
		It("should retrieve an account custom data", func() {
			customData, err := account.GetCustomData()

			Expect(err).NotTo(HaveOccurred())
			Expect(customData).NotTo(BeEmpty())
		})
		It("should return error if account doesn't exists", func() {
			account := &Account{}
			account.Href = BaseURL + "/accounts/XXXX"

			customData, err := account.GetCustomData()
			Expect(err).To(HaveOccurred())
			Expect(err.(Error).Status).To(Equal(404))
			Expect(customData).To(BeNil())
		})
	})

	Describe("UpdateCustomData", func() {
		It("should set an account custom data", func() {
			customData, err := account.UpdateCustomData(map[string]interface{}{"custom": "data"})

			Expect(err).NotTo(HaveOccurred())
			Expect(customData["custom"]).To(Equal("data"))
		})

		It("should update an account custom data", func() {
			account.UpdateCustomData(map[string]interface{}{"custom": "data"})
			customData, err := account.UpdateCustomData(map[string]interface{}{"custom": "nodata"})

			Expect(err).NotTo(HaveOccurred())
			Expect(customData["custom"]).To(Equal("nodata"))
		})
	})
})
