; ModuleID = '/var/work/Elrond/arwen-wasm-vm/test/contracts/breakpoint/breakpoint.c'
source_filename = "/var/work/Elrond/arwen-wasm-vm/test/contracts/breakpoint/breakpoint.c"
target datalayout = "e-m:e-p:32:32-i64:64-n32:64-S128"
target triple = "wasm32-unknown-unknown-wasm"

@__const.testFunc.msg = private unnamed_addr constant [10 x i8] c"exit here\00", align 1
@__const.testFunc.msg2 = private unnamed_addr constant [11 x i8] c"exit later\00", align 1

; Function Attrs: nounwind
define void @testFunc() local_unnamed_addr #0 {
entry:
  %msg = alloca [10 x i8], align 1
  %msg2 = alloca [11 x i8], align 1
  %call = tail call i64 @int64getArgument(i32 0) #4
  %cmp = icmp eq i64 %call, 1
  br i1 %cmp, label %if.then, label %if.else

if.then:                                          ; preds = %entry
  %0 = getelementptr inbounds [10 x i8], [10 x i8]* %msg, i32 0, i32 0
  call void @llvm.lifetime.start.p0i8(i64 10, i8* nonnull %0) #4
  call void @llvm.memcpy.p0i8.p0i8.i32(i8* nonnull align 1 %0, i8* align 1 getelementptr inbounds ([10 x i8], [10 x i8]* @__const.testFunc.msg, i32 0, i32 0), i32 10, i1 false)
  call void @signalError(i8* nonnull %0, i32 9) #4
  %1 = getelementptr inbounds [11 x i8], [11 x i8]* %msg2, i32 0, i32 0
  call void @llvm.lifetime.start.p0i8(i64 11, i8* nonnull %1) #4
  call void @llvm.memcpy.p0i8.p0i8.i32(i8* nonnull align 1 %1, i8* align 1 getelementptr inbounds ([11 x i8], [11 x i8]* @__const.testFunc.msg2, i32 0, i32 0), i32 11, i1 false)
  call void @signalError(i8* nonnull %1, i32 10) #4
  call void @llvm.lifetime.end.p0i8(i64 11, i8* nonnull %1) #4
  call void @llvm.lifetime.end.p0i8(i64 10, i8* nonnull %0) #4
  br label %if.end

if.else:                                          ; preds = %entry
  tail call void @int64finish(i64 100) #4
  br label %if.end

if.end:                                           ; preds = %if.else, %if.then
  ret void
}

; Function Attrs: argmemonly nounwind
declare void @llvm.lifetime.start.p0i8(i64 immarg, i8* nocapture) #1

declare i64 @int64getArgument(i32) local_unnamed_addr #2

; Function Attrs: argmemonly nounwind
declare void @llvm.memcpy.p0i8.p0i8.i32(i8* nocapture writeonly, i8* nocapture readonly, i32, i1 immarg) #1

declare void @signalError(i8*, i32) local_unnamed_addr #2

; Function Attrs: argmemonly nounwind
declare void @llvm.lifetime.end.p0i8(i64 immarg, i8* nocapture) #1

declare void @int64finish(i64) local_unnamed_addr #2

; Function Attrs: norecurse nounwind readnone
define void @init() local_unnamed_addr #3 {
entry:
  ret void
}

; Function Attrs: norecurse nounwind readnone
define void @_main() local_unnamed_addr #3 {
entry:
  ret void
}

attributes #0 = { nounwind "correctly-rounded-divide-sqrt-fp-math"="false" "disable-tail-calls"="false" "less-precise-fpmad"="false" "min-legal-vector-width"="0" "no-frame-pointer-elim"="false" "no-infs-fp-math"="false" "no-jump-tables"="false" "no-nans-fp-math"="false" "no-signed-zeros-fp-math"="false" "no-trapping-math"="false" "stack-protector-buffer-size"="8" "unsafe-fp-math"="false" "use-soft-float"="false" }
attributes #1 = { argmemonly nounwind }
attributes #2 = { "correctly-rounded-divide-sqrt-fp-math"="false" "disable-tail-calls"="false" "less-precise-fpmad"="false" "no-frame-pointer-elim"="false" "no-infs-fp-math"="false" "no-nans-fp-math"="false" "no-signed-zeros-fp-math"="false" "no-trapping-math"="false" "stack-protector-buffer-size"="8" "unsafe-fp-math"="false" "use-soft-float"="false" }
attributes #3 = { norecurse nounwind readnone "correctly-rounded-divide-sqrt-fp-math"="false" "disable-tail-calls"="false" "less-precise-fpmad"="false" "min-legal-vector-width"="0" "no-frame-pointer-elim"="false" "no-infs-fp-math"="false" "no-jump-tables"="false" "no-nans-fp-math"="false" "no-signed-zeros-fp-math"="false" "no-trapping-math"="false" "stack-protector-buffer-size"="8" "unsafe-fp-math"="false" "use-soft-float"="false" }
attributes #4 = { nounwind }

!llvm.module.flags = !{!0}
!llvm.ident = !{!1}

!0 = !{i32 1, !"wchar_size", i32 4}
!1 = !{!"clang version 9.0.0 (tags/RELEASE_900/final)"}
