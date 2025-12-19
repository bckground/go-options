package test

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestOptions(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Options Suite")
}

// Make sure we don't try to generate this
type OptionMyRenamedInt int // nolint:structcheck,unused // just exists so we would conflict with it

type OptionSetMyInt123 struct{}

func (o OptionSetMyInt123) apply(c *config) error {
	c.myInt = int(123)
	return nil
}

type OptionMakeError struct{}

func (o OptionMakeError) apply(c *config) error {
	return errors.New("bad news")
}

var _ = Describe("Generating options", func() {
	cfg := config{}

	It("generates options to set config value", func() {
		myInt := 456
		err := applyConfigOptions(&cfg,
			OptionNs.MyInt(123),
			OptionNs.MyFloat(4.56),
			OptionNs.MyString("my-string"),
			OptionNs.MyIntPointer(&myInt),
			OptionNs.MyInterface(789),
			OptionNs.MyFunc(func() int { return 0 }),
		)
		Ω(err).ShouldNot(HaveOccurred())
		Ω(cfg.myInt).Should(Equal(123))
		Ω(cfg.myFloat).Should(Equal(4.56))
		Ω(cfg.myString).Should(Equal("my-string"))
		Ω(cfg.myIntPointer).Should(Equal(&myInt))
		Ω(cfg.myInterface).Should(Equal(789))
	})

	It("generates an new function create a config", func() {
		cfg, err := newConfig(OptionNs.MyInt(123))
		Ω(err).ShouldNot(HaveOccurred())
		Ω(cfg.myInt).Should(Equal(123))
	})

	It("sets default values", func() {
		err := applyConfigOptions(&cfg)
		Ω(err).ShouldNot(HaveOccurred())
		Ω(cfg.myIntWithDefault).Should(Equal(1))
		Ω(cfg.myStringWithDefault).Should(Equal("default string"))
		Ω(cfg.myFloatWithDefault).Should(Equal(1.23))
	})

	It("compares using standard equality", func() {
		Ω(OptionNs.MyInt(1)).Should(Equal(OptionNs.MyInt(1)))
	})

	It("generates a String method", func() {
		Ω(fmt.Sprintf("%v", OptionNs.MyInt(1))).Should(Equal("MyInt: 1"))
		Ω(fmt.Sprintf("%v", OptionNs.MyString("abc"))).Should(Equal("MyString: abc"))
	})

	It("returns errors", func() {
		err := applyConfigOptions(&cfg, OptionMakeError{})
		Ω(err).Should(MatchError("bad news"))
	})

	It("allows option constructor to be renamed", func() {
		err := applyConfigOptions(&cfg, OptionNs.YourInt(1))
		Ω(err).ShouldNot(HaveOccurred())
		Ω(cfg.myRenamedInt).To(Equal(1))
	})

	Describe("custom options", func() {
		It("can be extended with custom options", func() {
			err := applyConfigOptions(&cfg, OptionSetMyInt123{})
			Ω(err).ShouldNot(HaveOccurred())
			Ω(cfg.myInt).Should(Equal(123))
		})

		It("returns error from custom options", func() {
			err := applyConfigOptions(&cfg, OptionMakeError{})
			Ω(err).Should(MatchError("bad news"))
		})
	})

	Describe("imports", func() {
		It("works with imported types", func() {
			err := applyConfigOptions(&cfg, OptionNs.MyDuration(time.Second))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(cfg.myDuration).Should(Equal(time.Second))
		})

		It("works with aliased imports", func() {
			err := applyConfigOptions(&cfg, OptionNs.MyDuration2(time.Second))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(cfg.myDuration).Should(Equal(time.Second))
		})

		It("works with nested packages", func() {
			myURL, err := url.Parse("http://example.com")
			Ω(err).ShouldNot(HaveOccurred())
			err = applyConfigOptions(&cfg, OptionNs.MyURL(*myURL))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(cfg.myURL).Should(Equal(*myURL))
		})
	})

	Describe("pointers", func() {
		It("can store a pointer to let us know if a value was set", func() {
			err := applyConfigOptions(&cfg)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(cfg.myPointerToInt).Should(BeNil())

			err = applyConfigOptions(&cfg, OptionNs.MyPointerToInt(1))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(cfg.myPointerToInt).ShouldNot(BeNil())
			Ω(*cfg.myPointerToInt).Should(Equal(1))
		})

		It("generates a String method", func() {
			Ω(fmt.Sprintf("%v", OptionNs.MyPointerToInt(1))).Should(Equal("MyPointerToInt: 1"))
		})
	})

	Describe("nested structs", func() {
		It("generates a constructor", func() {
			err := applyConfigOptions(&cfg, OptionNs.MyStruct(1, 2))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(cfg.myStruct.a).Should(Equal(1))
			Ω(cfg.myStruct.b).Should(Equal(2))
		})

		It("allows default values", func() {
			err := applyConfigOptions(&cfg)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(cfg.myStructWithDefault.a).Should(Equal(1))
		})

		It("defaults pointer structs to nil", func() {
			err := applyConfigOptions(&cfg)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(cfg.myPointerToStruct).Should(BeNil())

			err = applyConfigOptions(&cfg, OptionNs.MyPointerToStruct(1, 2))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(cfg.myPointerToStruct).ShouldNot(BeNil())
			Ω(cfg.myPointerToStruct.a).Should(Equal(1))
			Ω(cfg.myPointerToStruct.b).Should(Equal(2))
		})

		It("allows variadic arguments within a slice", func() {
			err := applyConfigOptions(&cfg, OptionNs.MyStructWithVariadicSlice(1, 1, 2))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(cfg.myStructWithDefault.a).Should(Equal(1))
		})

		It("generates a String method", func() {
			Ω(fmt.Sprintf("%v", OptionNs.MyStructWithVariadicSlice(1, 2))).Should(
				Equal("MyStructWithVariadicSlice: {a:1 b:[2]}"))
		})

		It("allows variadic arguments to be compared with cmp", func() {
			cmp.Equal(
				OptionNs.MyStructWithVariadicSlice(1, 1, 2),
				OptionNs.MyStructWithVariadicSlice(1, 1, 2))
		})
	})

	Describe("variadic slices", func() {
		It("creates a variadic constructor", func() {
			err := applyConfigOptions(&cfg, OptionNs.MySlice(1, 2))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(cfg.mySlice).Should(ConsistOf(1, 2))
		})

		It("allows them to be optional", func() {
			err := applyConfigOptions(&cfg)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(cfg.myPointerToSlice).Should(BeNil())

			err = applyConfigOptions(&cfg, OptionNs.MyPointerToSlice(1, 2))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(cfg.myPointerToSlice).ShouldNot(BeNil())
			Ω(*cfg.myPointerToSlice).Should(ConsistOf(1, 2))
		})

		It("allows them to be renamed", func() {
			err := applyConfigOptions(&cfg, OptionNs.YourSlice(1, 2))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(cfg.myRenamedSlice).ShouldNot(BeNil())
			Ω(cfg.myRenamedSlice).Should(ConsistOf(1, 2))
		})

		It("generates a String method", func() {
			Ω(fmt.Sprintf("%v", OptionNs.MySlice(1, 2))).Should(Equal("MySlice: [1 2]"))
		})

		It("allows them to the to be compared with cmp", func() {
			Ω(cmp.Equal(
				OptionNs.MySlice(1, 2),
				OptionNs.MySlice(1, 2))).Should(BeTrue())
		})
	})
})

var _ = Describe("Customizing the apply function name", func() {
	cfg := configWithDifferentApply{}

	It("uses the provided function name", func() {
		err := applyDifferent(&cfg)
		Ω(err).ShouldNot(HaveOccurred())
	})
})

var _ = Describe("Customizing the option prefix", func() {
	It("creates options with the custom prefix", func() {
		_, err := newConfigWithDifferentPrefix(MyOptNs.MyFloat(1.23))
		Ω(err).ShouldNot(HaveOccurred())
	})
})

var _ = Describe("Customizing the option suffix", func() {
	It("creates options with the custom prefix", func() {
		_, err := newConfigWithSuffix(SuffixOptionNs.MyFloat(1.23))
		Ω(err).ShouldNot(HaveOccurred())
	})
})

var _ = Describe("not quoting strings by default", func() {
	It("requires them to be quoted", func() {
		cfg, err := newConfigWithUnquotedString()
		Ω(err).ShouldNot(HaveOccurred())
		Ω(cfg.myString).Should(Equal("quoted"))
	})
})

var _ = Describe("Disabling cmp", func() {
	It("prevents options from implementing Equal", func() {
		_, equalsFound := reflect.TypeOf(NoCmpOptionNs.MyInt(1)).MethodByName("Equal")
		Ω(equalsFound).Should(BeFalse())
	})
})

var _ = Describe("Disabling stringer", func() {
	It("prevents options from implementing String", func() {
		Ω(fmt.Sprintf("%v", NoStringerOptionNs.MyInt(1))).Should(Equal("{1}"))
	})
})

var _ = Describe("Public new function", func() {
	It("Makes the new confic cunction public", func() {
		cfg, err := NewConfigWithPublicNewFunc(PublicFuncOptionNs.MyInt(10))
		Ω(err).ShouldNot(HaveOccurred())
		Ω(cfg.myInt).Should(Equal(10))
	})
})
