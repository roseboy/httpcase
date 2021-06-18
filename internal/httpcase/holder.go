package httpcase

type Holder struct {
	Function    *Function
	TestCase    *TestCase
	TestContext *TestContext
}

var TestHolder = &Holder{}
