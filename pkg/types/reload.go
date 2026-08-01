package types

type Reload struct {
	FileName   string
	DataStruct any
	Validate   func(Utils) bool
	IfOkay     func(Utils)
	IfFail     func(Utils)
}
