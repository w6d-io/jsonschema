package jsonschema_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/w6d-io/x/errorx"

	"github.com/w6d-io/jsonschema"
	"gopkg.in/yaml.v3"

	_ "embed"
)

//go:embed testdata/test1_schema.json
var schema1 string

//go:embed testdata/test2_schema.json
var schema2 string

var test1 = `listen: ":8080"
get:
  url: "unit-test:latest"
  timeout: 600
list:
  - "1"
  - "2"
`

var test2 = `{
  "listen": ":8080",
  "get": {
    "url": "unit-test:latest",
    "timeout": 600
  },
  "list": [
    "1",
    "2"
  ]
}`

var _ = Describe("", func() {
	Context("", func() {
		BeforeEach(func() {
		})
		AfterEach(func() {
		})
		It("success the config validation", func() {
			var foo map[string]interface{}
			Expect(yaml.Unmarshal([]byte(test1), &foo)).To(Succeed())

			Expect(jsonschema.AddSchema(jsonschema.Config, schema1)).To(Succeed())
			Expect(jsonschema.Config.Validate(foo)).To(Succeed())
		})
		It("success the config validation with byte", func() {
			Expect(jsonschema.AddSchema(jsonschema.Config, schema1)).To(Succeed())
			Expect(jsonschema.Config.Validate([]byte(test2))).To(Succeed())
		})
		It("fail on schema does not exists", func() {
			var foo map[string]interface{}
			Expect(yaml.Unmarshal([]byte(test1), &foo)).To(Succeed())
			var c jsonschema.SchemaType = 1
			err := c.Validate(foo)
			Expect(err).To(HaveOccurred())
			e := err.(*errorx.Error)
			Expect(e).To(Equal(&errorx.Error{Message: "the specified schema type (1) is not supported"}))
		})
		It("fails on compilation", func() {
			var foo map[string]interface{}
			Expect(yaml.Unmarshal([]byte(test1), &foo)).To(Succeed())
			var st jsonschema.SchemaType = 2
			Expect(jsonschema.AddSchema(st, schema2)).To(Succeed())
			err := st.Validate(foo)
			Expect(err).To(HaveOccurred())
		})
	})
})
