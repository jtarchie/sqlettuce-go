package db_test

import (
	"context"

	"github.com/jtarchie/sqlettus/db"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Strings", func() {
	var client *db.Client

	BeforeEach(func() {
		var err error

		client, err = db.NewClient(":memory:?cache=shared&mode=memory")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		client.Close()
	})

	When("Set", func() {
		It("sets a value", func() {
			err := client.Set(context.TODO(), "key", "value")
			Expect(err).NotTo(HaveOccurred())

			value, found, err := client.Get(context.TODO(), "key")
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(value).To(Equal("value"))
		})
	})

	When("Get", func() {
		It("return nil with non existent keys", func() {
			value, found, err := client.Get(context.TODO(), "key")
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeFalse())
			Expect(value).To(Equal(""))
		})
	})

	When("Delete", func() {
		It("can delete a value", func() {
			err := client.Set(context.TODO(), "key", "value")
			Expect(err).NotTo(HaveOccurred())

			value, found, err := client.Delete(context.TODO(), "key")
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(value).To(Equal("value"))

			value, found, err = client.Get(context.TODO(), "key")
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeFalse())
			Expect(value).To(Equal(""))
		})

		It("does not fail on missing value", func() {
			value, found, err := client.Delete(context.TODO(), "key")
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeFalse())
			Expect(value).To(Equal(""))
		})
	})

	When("Append", func() {
		It("appends values to a string", func() {
			length, err := client.Append(context.TODO(), "name", "Hello")
			Expect(err).NotTo(HaveOccurred())
			Expect(length).To(BeEquivalentTo(5))

			value, found, err := client.Get(context.TODO(), "name")
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(value).To(Equal("Hello"))

			length, err = client.Append(context.TODO(), "name", " World")
			Expect(err).NotTo(HaveOccurred())
			Expect(length).To(BeEquivalentTo(11))

			value, found, err = client.Get(context.TODO(), "name")
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(value).To(Equal("Hello World"))
		})
	})

	When("Substr", func() {
		It("handles start and end, using negative indexes, too", func() {
			err := client.Set(context.TODO(), "key", "This is a string")
			Expect(err).NotTo(HaveOccurred())

			value, err := client.Substr(context.TODO(), "key", 0, 3)
			Expect(err).NotTo(HaveOccurred())
			Expect(value).To(Equal("This"))

			value, err = client.Substr(context.TODO(), "key", -3, -1)
			Expect(err).NotTo(HaveOccurred())
			Expect(value).To(Equal("ing"))

			value, err = client.Substr(context.TODO(), "key", 0, -1)
			Expect(err).NotTo(HaveOccurred())
			Expect(value).To(Equal("This is a string"))

			value, err = client.Substr(context.TODO(), "key", 10, 100)
			Expect(err).NotTo(HaveOccurred())
			Expect(value).To(Equal("string"))

			value, err = client.Substr(context.TODO(), "nokey", 0, 1)
			Expect(err).NotTo(HaveOccurred())
			Expect(value).To(Equal(""))
		})
	})
})
