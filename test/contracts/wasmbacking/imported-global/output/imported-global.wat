(module
  (type $singleParam (func (param i64)))
	(type $void (func))
	(global $importedG (import "env" "importedG") i64)
  (import "env" "int64finish" (func $int64finish (type $singleParam)))
	(func $main (type $singleParam)
		global.get $importedG
		call $int64finish
	)
  (memory (;0;) 2)
  (export "memory" (memory 0))
  (export "get_imported_global" (func $main))
)
