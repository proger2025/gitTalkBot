package parsing



func BuildSymbol(res *ParseResult) []Symbol {
	var newSymb []Symbol

	for _, v := range res.Funcs {
		newSymbStruct := Symbol {
			v.PackageName + "." + v.Name,
			"func",
			v.PackageName,
			v.Signature,
			v.DocComment,
		}

		newSymb = append(newSymb, newSymbStruct)
	}

	for _, v := range res.Methods {
		newSymbStruct := Symbol {
			v.PackageName + "." + v.Name ,
			"method",
			v.PackageName,
			v.Signature,
			v.DocComment,
		}

		newSymb = append(newSymb, newSymbStruct)
	}
	
	return newSymb

}