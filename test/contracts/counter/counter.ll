; ModuleID = '/home/camil.bancioiu/Work/Elrond/arwen-wasm-vm/test/contracts/counter/counter.c'
source_filename = "/home/camil.bancioiu/Work/Elrond/arwen-wasm-vm/test/contracts/counter/counter.c"
target datalayout = "e-m:e-p:32:32-i64:64-n32:64-S128"
target triple = "wasm32-unknown-unknown-wasm"

@counterKey = global <{ [9 x i8], [23 x i8] }> <{ [9 x i8] c"mycounter", [23 x i8] zeroinitializer }>, align 16

; Function Attrs: nounwind
define void @init() local_unnamed_addr #0 {
entry:
  %call = tail call i32 @int64storageStore(i8* getelementptr inbounds (<{ [9 x i8], [23 x i8] }>, <{ [9 x i8], [23 x i8] }>* @counterKey, i32 0, i32 0, i32 0), i64 1) #2
  ret void
}

declare i32 @int64storageStore(i8*, i64) local_unnamed_addr #1

; Function Attrs: nounwind
define void @increment() local_unnamed_addr #0 {
entry:
  %call = tail call i64 @int64storageLoad(i8* getelementptr inbounds (<{ [9 x i8], [23 x i8] }>, <{ [9 x i8], [23 x i8] }>* @counterKey, i32 0, i32 0, i32 0)) #2
  %inc = add i64 %call, 1
  %call1 = tail call i32 @int64storageStore(i8* getelementptr inbounds (<{ [9 x i8], [23 x i8] }>, <{ [9 x i8], [23 x i8] }>* @counterKey, i32 0, i32 0, i32 0), i64 %inc) #2
  tail call void @int64finish(i64 %inc) #2
  ret void
}

declare i64 @int64storageLoad(i8*) local_unnamed_addr #1

declare void @int64finish(i64) local_unnamed_addr #1

; Function Attrs: nounwind
define void @decrement() local_unnamed_addr #0 {
entry:
  %call = tail call i64 @int64storageLoad(i8* getelementptr inbounds (<{ [9 x i8], [23 x i8] }>, <{ [9 x i8], [23 x i8] }>* @counterKey, i32 0, i32 0, i32 0)) #2
  %dec = add i64 %call, -1
  %call1 = tail call i32 @int64storageStore(i8* getelementptr inbounds (<{ [9 x i8], [23 x i8] }>, <{ [9 x i8], [23 x i8] }>* @counterKey, i32 0, i32 0, i32 0), i64 %dec) #2
  tail call void @int64finish(i64 %dec) #2
  ret void
}

; Function Attrs: nounwind
define void @get() local_unnamed_addr #0 {
entry:
  %call = tail call i64 @int64storageLoad(i8* getelementptr inbounds (<{ [9 x i8], [23 x i8] }>, <{ [9 x i8], [23 x i8] }>* @counterKey, i32 0, i32 0, i32 0)) #2
  tail call void @int64finish(i64 %call) #2
  ret void
}

attributes #0 = { nounwind "correctly-rounded-divide-sqrt-fp-math"="false" "disable-tail-calls"="false" "less-precise-fpmad"="false" "min-legal-vector-width"="0" "no-frame-pointer-elim"="false" "no-infs-fp-math"="false" "no-jump-tables"="false" "no-nans-fp-math"="false" "no-signed-zeros-fp-math"="false" "no-trapping-math"="false" "stack-protector-buffer-size"="8" "unsafe-fp-math"="false" "use-soft-float"="false" }
attributes #1 = { "correctly-rounded-divide-sqrt-fp-math"="false" "disable-tail-calls"="false" "less-precise-fpmad"="false" "no-frame-pointer-elim"="false" "no-infs-fp-math"="false" "no-nans-fp-math"="false" "no-signed-zeros-fp-math"="false" "no-trapping-math"="false" "stack-protector-buffer-size"="8" "unsafe-fp-math"="false" "use-soft-float"="false" }
attributes #2 = { nounwind }

!llvm.module.flags = !{!0}
!llvm.ident = !{!1}

!0 = !{i32 1, !"wchar_size", i32 4}
!1 = !{!"clang version 9.0.0 (tags/RELEASE_900/final)"}
