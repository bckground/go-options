package test

import (
	"net/url"
	"time"
	time2 "time"
)

//go:generate go run .. -imports=time,net/url,time2=time -namespace=OptionNs -option=Option config
type config struct {
	myInt            int
	myIntWithDefault int `options:",1"`
	myRenamedInt     int `options:"yourInt"`

	// does something
	myDocumentedInt int
	myCommentedInt  int // for some reason

	// does something else
	myDocAndCommentInt int // for some other reason

	// takes a float
	myFloat            float64 // really a float
	myFloatWithDefault float64 `options:",1.23"`

	myString              string
	myStringWithDefault   string `options:",default string"`
	myStringWithoutOption string `options:"-"` // nolint:structcheck,unused // not expected to be used

	myFunc func() int

	myIntPointer *int

	myInterface interface{}

	// types requiring imports
	myURL       url.URL
	myDuration  time.Duration
	myDuration2 time2.Duration

	myStruct            struct{ a, b int }
	myStructWithDefault struct {
		a int `options:",1"`
	}
	myPointerToStruct         *struct{ a, b int }
	myStructWithVariadicSlice struct {
		a int
		b []int `options:"..."`
	}

	mySlice          []int  `options:"..."`
	myPointerToSlice *[]int `options:"..."`
	myRenamedSlice   []int  `options:"yourSlice..."`

	myPointerToInt        *int `options:"*"`
	myPointerToRenamedInt *int `options:"*yourIntWithPointer"`

	// ensure we can handle multiple tags
	WithJsonTagButNoOptions string `json:"-"`
	WithBothJsonAndOptions  string `json:"-" options:"gotBoth"`
}

//go:generate go run .. -func applyDifferent -option DifferentOption -new=false configWithDifferentApply
type configWithDifferentApply struct {
}

//go:generate go run .. -namespace=MyOptNs -option=MyOpt configWithDifferentPrefix
type configWithDifferentPrefix struct {
	myFloat float64
}

//go:generate go run .. -namespace=SuffixOptionNs -option=SuffixOption configWithSuffix
type configWithSuffix struct {
	myFloat float64
}

//go:generate go run .. -quote-default-strings=false -namespace=UnquotedOptionNs -option=UnquotedOption configWithUnquotedString
type configWithUnquotedString struct {
	myString string `options:",\"quoted\""`
}

//go:generate go run .. -cmp=false -namespace=NoCmpOptionNs -option=NoCmpOption configWithoutCmp
type configWithoutCmp struct {
	myInt int
}

//go:generate go run .. -stringer=false -namespace=NoStringerOptionNs -option=NoStringerOption configWithoutStringer
type configWithoutStringer struct {
	myInt int
}

//go:generate go run .. -noerror=false -namespace=NoErrorOptionNs -option=NoErrorOption configWithNoError
type configWithNoError struct {
	myInt int
}

//go:generate go run .. -public=true -namespace=PublicFuncOptionNs -option=PublicFuncOption configWithPublicNewFunc
type configWithPublicNewFunc struct {
	myInt int
}

//go:generate go run .. -build=testing -func applyBuild -namespace=BuildOptionNs -option=BuildOption configWithBuild
type configWithBuild struct {
	myInt int
}
