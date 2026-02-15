package parsing

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
)


func Ast(repoFiles []string) *ParseResult {
	res := &ParseResult{}

	for _, fileFrom := range repoFiles {
		f, m, t, err := fileParsing(fileFrom)
		if err != nil {
			log.Println(err)
			continue
		}

		res.Funcs = append(res.Funcs, f...)
		res.Methods = append(res.Methods, m...)
		res.Types = append(res.Types, t...)
	}

	return res
}

func fileParsing(fileFrom string) ([]FuncS, []MethodS, []TypeS, error) {
	var funcs []FuncS
	var methods []MethodS
	var types []TypeS
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, fileFrom, nil, parser.ParseComments)
	if err != nil {
    	log.Println(err)
    	return nil, nil, nil, err
	}

	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		
		case *ast.FuncDecl:
			if !ast.IsExported(x.Name.Name) {
				return true
			}

			docX := checkDocDecl(x)
			if x.Recv == nil {
				newStuctFuncs := FuncS {
					x.Name.Name,
					getSignature(fset, x),
					docX,
					node.Name.Name,
					fileFrom,
				}
				funcs = append(funcs, newStuctFuncs)
			} else  {
				newStuctMethod := MethodS {
					x.Name.Name,
					getSignature(fset, x),
					docX,
					node.Name.Name,
					fileFrom,
				}
				methods = append(methods, newStuctMethod)
			}
		case *ast.GenDecl:
    		if x.Tok != token.TYPE {
    		    return true
    		}
		
    		blockDoc := ""
    		if x.Doc != nil {
    		    blockDoc = x.Doc.Text()
    		}
		
    		for _, spec := range x.Specs {
    		    ts, ok := spec.(*ast.TypeSpec)
    		    if !ok {
    		        continue
    		    }
			
    		    if !ast.IsExported(ts.Name.Name) {
    		        continue
    		    }
			
    		    docX := blockDoc
    		    if ts.Doc != nil {
    		        docX = ts.Doc.Text()
    		    }
			
    		    kind := "alias"
    		    switch ts.Type.(type) {
    		    case *ast.StructType:
    		        kind = "struct"
    		    case *ast.InterfaceType:
    		        kind = "interface"
    		    }
			
    		    newStructTypes := TypeS{
    		        ts.Name.Name,   
    		        kind,          
    		        docX,          
    		        node.Name.Name, 
    		        fileFrom,      
    		    }
			
    		    types = append(types, newStructTypes)
    		}
		}
		return true
	})

	return funcs, methods, types, err

}


func getSignature(fset *token.FileSet, decl *ast.FuncDecl) string {
	var typeFunc bytes.Buffer
	printer.Fprint(&typeFunc, fset, decl.Type)

	if decl.Recv == nil {
		return "func " + decl.Name.Name + typeFunc.String()
	} else {
		var receiverStr bytes.Buffer
		printer.Fprint(&receiverStr, fset, decl.Recv)
		return "func " + receiverStr.String() + " " + decl.Name.Name + typeFunc.String()
	}

}

func checkDocDecl(decl *ast.FuncDecl) string {
	if decl.Doc == nil {
		return ""
	} else {
		return decl.Doc.Text()
	}
}


