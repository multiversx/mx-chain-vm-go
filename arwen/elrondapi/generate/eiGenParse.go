package elrondapigenerate

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

func lowerInitial(name string) string {
	return strings.ToLower(name[0:1]) + name[1:]
}

func upperInitial(name string) string {
	return strings.ToUpper(name[0:1]) + name[1:]
}

func extractEIFunctionArguments(decl *ast.FuncDecl) ([]*EIFunctionArg, error) {
	var arguments []*EIFunctionArg
	for _, param := range decl.Type.Params.List {
		for _, name := range param.Names {
			arguments = append(arguments, &EIFunctionArg{
				Name: name.String(),
				Type: fmt.Sprintf("%s", param.Type),
			})
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

// isInterfaceMethod looks in the comments to determine if to take method into consideration or not
func isEIInterfaceMethod(decl *ast.FuncDecl) bool {
	if decl.Doc == nil {
		return false
	}

	text := decl.Doc.Text()
	// TODO: maybe also validate that doc is well-formed
	return strings.Contains(text, "@autogenerate(VMHooks)")
}

func validateReceiver(decl *ast.FuncDecl) error {
	if decl.Recv == nil {
		return errors.New("receiver expected")
	}
	if len(decl.Recv.List) != 1 {
		return errors.New("single receiver expected")
	}
	field := decl.Recv.List[0]
	if len(field.Names) != 1 {
		return errors.New("single receiver field name expected")
	}
	if field.Names[0].String() != "context" {
		return errors.New("method receiver should be named 'context'")
	}
	return nil
}

func extractEIFunction(decl *ast.FuncDecl) (*EIFunction, error) {
	interfaceFunction := isEIInterfaceMethod(decl)
	if !interfaceFunction {
		return nil, nil
	}
	originalFunctionName := decl.Name.Name
	err := validateReceiver(decl)
	if err != nil {
		return nil, fmt.Errorf("invalid method %s, %w", originalFunctionName, err)
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
		OriginalName:    originalFunctionName,
		LowerCaseName:   lowerInitial(originalFunctionName),
		CapitalizedName: upperInitial(originalFunctionName),
		Arguments:       arguments,
		Result:          result,
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
	for _, group := range eiMetadata.Groups {
		f, err := parser.ParseFile(fset, pathToSources+group.SourcePath, nil, parser.ParseComments)
		if err != nil {
			return err
		}
		fileFunctions, err := extractEIFunctions(f)
		if err != nil {
			return err
		}
		group.Functions = fileFunctions
		eiMetadata.AllFunctions = append(eiMetadata.AllFunctions, fileFunctions...)
	}
	return nil
}
