package types

type Reload struct {
	FileName   string
	DataStruct any
	IfOkay     func(Utils)
	IfFail     func(Utils)
}
