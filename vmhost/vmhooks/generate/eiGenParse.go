package vmhooksgenerate

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

func processType(ty ast.Expr) (EIType, error) {
	switch v := ty.(type) {
	case *ast.Ident:
		switch v.Name {
		case "int32":
			return EITypeInt32, nil
		case "int64":
			return EITypeInt64, nil
		}
	case *ast.SelectorExpr:
		if module, ok := v.X.(*ast.Ident); ok && module.Name == "executor" {
			switch v.Sel.Name {
			case "MemPtr":
				return EITypeMemPtr, nil
			case "MemLength":
				return EITypeMemLength, nil
			}
		}
	}

	return EITypeInvalid, fmt.Errorf("invalid EI type: %s", ty)
}

func extractEIFunctionArguments(decl *ast.FuncDecl) ([]*EIFunctionArg, error) {
	var arguments []*EIFunctionArg
	for _, param := range decl.Type.Params.List {
		for _, name := range param.Names {
			eiType, err := processType(param.Type)
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, &EIFunctionArg{
				Name: name.String(),
				Type: eiType,
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
		eiType, err := processType(decl.Type.Results.List[0].Type)
		if err != nil {
			return nil, err
		}
		return &EIFunctionResult{
			Type: eiType,
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
	err := validateReceiver(decl)
	if err != nil {
		return nil, fmt.Errorf("invalid method %s, %w", decl.Name.Name, err)
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
		Name:      decl.Name.Name,
		Arguments: arguments,
		Result:    result,
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

// ReadAndParseEIMetadata will read and parse EI metadata
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
