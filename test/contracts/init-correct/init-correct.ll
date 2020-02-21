; ModuleID = '/home/camil.bancioiu/Work/Elrond/arwen-wasm-vm/test/contracts/init-correct/init-correct.c'
source_filename = "/home/camil.bancioiu/Work/Elrond/arwen-wasm-vm/test/contracts/init-correct/init-correct.c"
target datalayout = "e-m:e-p:32:32-i64:64-n32:64-S128"
target triple = "wasm32-unknown-unknown-wasm"

@__const.init.msg = private unnamed_addr constant [16 x i8] c"init successful\00", align 16

; Function Attrs: nounwind
define void @init() local_unnamed_addr #0 {
entry:
  %msg = alloca [16 x i8], align 16
  %0 = getelementptr inbounds [16 x i8], [16 x i8]* %msg, i32 0, i32 0
  call void @llvm.lifetime.start.p0i8(i64 16, i8* nonnull %0) #3
  call void @llvm.memcpy.p0i8.p0i8.i32(i8* nonnull align 16 %0, i8* align 16 getelementptr inbounds ([16 x i8], [16 x i8]* @__const.init.msg, i32 0, i32 0), i32 16, i1 false)
  call void @finish(i8* nonnull %0, i32 15) #3
  call void @llvm.lifetime.end.p0i8(i64 16, i8* nonnull %0) #3
  ret void
}

; Function Attrs: argmemonly nounwind
declare void @llvm.lifetime.start.p0i8(i64 immarg, i8* nocapture) #1

; Function Attrs: argmemonly nounwind
declare void @llvm.memcpy.p0i8.p0i8.i32(i8* nocapture writeonly, i8* nocapture readonly, i32, i1 immarg) #1

declare void @finish(i8*, i32) local_unnamed_addr #2

; Function Attrs: argmemonly nounwind
declare void @llvm.lifetime.end.p0i8(i64 immarg, i8* nocapture) #1

attributes #0 = { nounwind "correctly-rounded-divide-sqrt-fp-math"="false" "disable-tail-calls"="false" "less-precise-fpmad"="false" "min-legal-vector-width"="0" "no-frame-pointer-elim"="false" "no-infs-fp-math"="false" "no-jump-tables"="false" "no-nans-fp-math"="false" "no-signed-zeros-fp-math"="false" "no-trapping-math"="false" "stack-protector-buffer-size"="8" "unsafe-fp-math"="false" "use-soft-float"="false" }
attributes #1 = { argmemonly nounwind }
attributes #2 = { "correctly-rounded-divide-sqrt-fp-math"="false" "disable-tail-calls"="false" "less-precise-fpmad"="false" "no-frame-pointer-elim"="false" "no-infs-fp-math"="false" "no-nans-fp-math"="false" "no-signed-zeros-fp-math"="false" "no-trapping-math"="false" "stack-protector-buffer-size"="8" "unsafe-fp-math"="false" "use-soft-float"="false" }
attributes #3 = { nounwind }

!llvm.module.flags = !{!0}
!llvm.ident = !{!1}

!0 = !{i32 1, !"wchar_size", i32 4}
!1 = !{!"clang version 9.0.0 (tags/RELEASE_900/final)"}
