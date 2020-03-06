; ModuleID = '/var/work/Elrond/arwen-wasm-vm/test/contracts/exec-same-ctx-simple-parent/exec-same-ctx-simple-parent.c'
source_filename = "/var/work/Elrond/arwen-wasm-vm/test/contracts/exec-same-ctx-simple-parent/exec-same-ctx-simple-parent.c"
target datalayout = "e-m:e-p:32:32-i64:64-n32:64-S128"
target triple = "wasm32-unknown-unknown-wasm"

@executeValue = global [32 x i8] c"\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00c", align 16
@__const.parentFunctionChildCall.childAddress = private unnamed_addr constant [33 x i8] c"secondSC........................\00", align 16
@__const.parentFunctionChildCall.functionName = private unnamed_addr constant [14 x i8] c"childFunction\00", align 1
@__const.parentFunctionChildCall.msg = private unnamed_addr constant [7 x i8] c"parent\00", align 1

; Function Attrs: nounwind
define void @parentFunctionChildCall() local_unnamed_addr #0 {
entry:
  %childAddress = alloca [33 x i8], align 16
  %functionName = alloca [14 x i8], align 1
  %msg = alloca [7 x i8], align 1
  %0 = getelementptr inbounds [33 x i8], [33 x i8]* %childAddress, i32 0, i32 0
  call void @llvm.lifetime.start.p0i8(i64 33, i8* nonnull %0) #3
  call void @llvm.memcpy.p0i8.p0i8.i32(i8* nonnull align 16 %0, i8* align 16 getelementptr inbounds ([33 x i8], [33 x i8]* @__const.parentFunctionChildCall.childAddress, i32 0, i32 0), i32 33, i1 false)
  %1 = getelementptr inbounds [14 x i8], [14 x i8]* %functionName, i32 0, i32 0
  call void @llvm.lifetime.start.p0i8(i64 14, i8* nonnull %1) #3
  call void @llvm.memcpy.p0i8.p0i8.i32(i8* nonnull align 1 %1, i8* align 1 getelementptr inbounds ([14 x i8], [14 x i8]* @__const.parentFunctionChildCall.functionName, i32 0, i32 0), i32 14, i1 false)
  %call = call i32 @executeOnSameContext(i64 200000, i8* nonnull %0, i8* getelementptr inbounds ([32 x i8], [32 x i8]* @executeValue, i32 0, i32 0), i8* nonnull %1, i32 13, i32 0, i8* null, i8* null) #3
  %conv = sext i32 %call to i64
  call void @int64finish(i64 %conv) #3
  %call4 = call i32 @executeOnSameContext(i64 200000, i8* nonnull %0, i8* getelementptr inbounds ([32 x i8], [32 x i8]* @executeValue, i32 0, i32 0), i8* nonnull %1, i32 13, i32 0, i8* null, i8* null) #3
  %conv5 = sext i32 %call4 to i64
  call void @int64finish(i64 %conv5) #3
  %2 = getelementptr inbounds [7 x i8], [7 x i8]* %msg, i32 0, i32 0
  call void @llvm.lifetime.start.p0i8(i64 7, i8* nonnull %2) #3
  call void @llvm.memcpy.p0i8.p0i8.i32(i8* nonnull align 1 %2, i8* align 1 getelementptr inbounds ([7 x i8], [7 x i8]* @__const.parentFunctionChildCall.msg, i32 0, i32 0), i32 7, i1 false)
  call void @finish(i8* nonnull %2, i32 6) #3
  call void @llvm.lifetime.end.p0i8(i64 7, i8* nonnull %2) #3
  call void @llvm.lifetime.end.p0i8(i64 14, i8* nonnull %1) #3
  call void @llvm.lifetime.end.p0i8(i64 33, i8* nonnull %0) #3
  ret void
}

; Function Attrs: argmemonly nounwind
declare void @llvm.lifetime.start.p0i8(i64 immarg, i8* nocapture) #1

; Function Attrs: argmemonly nounwind
declare void @llvm.memcpy.p0i8.p0i8.i32(i8* nocapture writeonly, i8* nocapture readonly, i32, i1 immarg) #1

declare i32 @executeOnSameContext(i64, i8*, i8*, i8*, i32, i32, i8*, i8*) local_unnamed_addr #2

declare void @int64finish(i64) local_unnamed_addr #2

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
