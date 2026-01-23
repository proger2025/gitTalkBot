package parsing


type FuncS struct {
	Name string
	Signature string
	DocComment string
	PackageName string
	FilePath string
}

type MethodS struct {
	Name string
	Signature string
	DocComment string
	PackageName string
	FilePath string

}

type TypeS struct {
	Name string
	Kind string
	DocComment string
	PackageName string
	FilePath string
}


type ParseResult struct {
	Funcs []FuncS
	Methods []MethodS
	Types []TypeS
}


type Symbol struct {
	Id string
	Kind string
	PackageName string
	Signature string
	DocComment string
}
