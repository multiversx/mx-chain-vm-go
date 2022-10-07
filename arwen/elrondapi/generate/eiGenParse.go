package elrondapigenerate

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

const eiFunctionPrefix = "v1_5_"

func publicFuncName(originalFuncName string) string {
	// trim prefix
	trimmed := strings.TrimPrefix(originalFuncName, eiFunctionPrefix)
	// capitalize
	return strings.ToUpper(trimmed[0:1]) + trimmed[1:]
}

func extractEIFunctionArguments(decl *ast.FuncDecl) ([]*EIFunctionArg, error) {
	var arguments []*EIFunctionArg
	for index, param := range decl.Type.Params.List {
		if index == 0 {
			if len(param.Names) != 1 && param.Names[0].Name != "context" {
				return nil, fmt.Errorf("bad context argument in function %s", decl.Name.Name)
			}
		} else {
			for _, name := range param.Names {
				arguments = append(arguments, &EIFunctionArg{
					Name: name.String(),
					Type: fmt.Sprintf("%s", param.Type),
				})
			}
		}
	}
	return arguments, nil
}

func extractEIFunctionResult(decl *ast.FuncDecl) (*EIFunctionResult, error) {
	if decl.Type.Results == nil {
		return nil, nil
	}
	switch len(decl.Type.Results.List) {
	case 0:
		return nil, nil
	case 1:
		return &EIFunctionResult{
			Type: fmt.Sprintf("%s", decl.Type.Results.List[0].Type),
		}, nil
	default:
		return nil, fmt.Errorf("too many results in function %s, no more than 1 accepted", decl.Name.Name)
	}
}

func extractEIFunction(decl *ast.FuncDecl) (*EIFunction, error) {
	if decl.Recv != nil {
		return nil, errors.New("no receiver expected")
	}
	originalFunctionName := decl.Name.Name
	if !strings.HasPrefix(originalFunctionName, eiFunctionPrefix) {
		return nil, nil
	}
	arguments, err := extractEIFunctionArguments(decl)
	if err != nil {
		return nil, err
	}
	result, err := extractEIFunctionResult(decl)
	if err != nil {
		return nil, err
	}
	eiFunction := &EIFunction{
		OriginalName: originalFunctionName,
		PublicName:   publicFuncName(originalFunctionName),
		Arguments:    arguments,
		Result:       result,
	}

	return eiFunction, nil
}

func extractEIFunctions(f *ast.File) ([]*EIFunction, error) {
	var result []*EIFunction
	for _, d := range f.Decls {
		if funcDecl, ok := d.(*ast.FuncDecl); ok {
			eiFunc, err := extractEIFunction(funcDecl)
			if err != nil {
				return nil, err
			}
			if eiFunc != nil {
				result = append(result, eiFunc)
			}
		}
	}
	return result, nil
}

func ReadAndParseEIMetadata(fset *token.FileSet, pathToSources string, eiMetadata *EIMetadata) error {
	for fileName, fileMetadata := range eiMetadata.FileMap {
		f, err := parser.ParseFile(fset, pathToSources+fileName, nil, parser.AllErrors)
		if err != nil {
			return err
		}
		fileFunctions, err := extractEIFunctions(f)
		if err != nil {
			return err
		}
		fileMetadata.Functions = fileFunctions
		eiMetadata.AllFunctions = append(eiMetadata.AllFunctions, fileFunctions...)
	}
	return nil
}
