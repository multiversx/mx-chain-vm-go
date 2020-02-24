; ModuleID = '/home/camil.bancioiu/Work/Elrond/arwen-wasm-vm/test/contracts/exec-same-ctx-child/exec-same-ctx-child.c'
source_filename = "/home/camil.bancioiu/Work/Elrond/arwen-wasm-vm/test/contracts/exec-same-ctx-child/exec-same-ctx-child.c"
target datalayout = "e-m:e-p:32:32-i64:64-n32:64-S128"
target triple = "wasm32-unknown-unknown-wasm"

@childKey = global [33 x i8] c"childKey........................\00", align 16
@childData = global [10 x i8] c"childData\00", align 1
@childFinish = global [12 x i8] c"childFinish\00", align 1
@recipient = global [32 x i8] zeroinitializer, align 16
@value = global [32 x i8] c"\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00`", align 16
@__const.childFunction.message = private unnamed_addr constant [26 x i8] c"wrong number of arguments\00", align 16

; Function Attrs: nounwind
define void @childFunction() local_unnamed_addr #0 {
entry:
  %message = alloca [26 x i8], align 16
  %transferData = alloca [100 x i8], align 16
  %call = tail call i32 bitcast (i32 (...)* @getNumArguments to i32 ()*)() #4
  %cmp = icmp eq i32 %call, 2
  br i1 %cmp, label %if.end, label %if.then

if.then:                                          ; preds = %entry
  %0 = getelementptr inbounds [26 x i8], [26 x i8]* %message, i32 0, i32 0
  call void @llvm.lifetime.start.p0i8(i64 26, i8* nonnull %0) #4
  call void @llvm.memcpy.p0i8.p0i8.i32(i8* nonnull align 16 %0, i8* align 16 getelementptr inbounds ([26 x i8], [26 x i8]* @__const.childFunction.message, i32 0, i32 0), i32 26, i1 false)
  call void @signalError(i8* nonnull %0, i32 25) #4
  call void @llvm.lifetime.end.p0i8(i64 26, i8* nonnull %0) #4
  br label %if.end

if.end:                                           ; preds = %entry, %if.then
  %call1 = call i32 @getArgument(i32 0, i8* getelementptr inbounds ([32 x i8], [32 x i8]* @recipient, i32 0, i32 0)) #4
  %1 = getelementptr inbounds [100 x i8], [100 x i8]* %transferData, i32 0, i32 0
  call void @llvm.lifetime.start.p0i8(i64 100, i8* nonnull %1) #4
  %call3 = call i32 @getArgument(i32 1, i8* nonnull %1) #4
  %call4 = call i32 @getArgumentLength(i32 1) #4
  %call6 = call i32 @transferValue(i8* getelementptr inbounds ([32 x i8], [32 x i8]* @recipient, i32 0, i32 0), i8* getelementptr inbounds ([32 x i8], [32 x i8]* @value, i32 0, i32 0), i8* nonnull %1, i32 %call4) #4
  %call7 = call i32 @storageStore(i8* getelementptr inbounds ([33 x i8], [33 x i8]* @childKey, i32 0, i32 0), i8* getelementptr inbounds ([10 x i8], [10 x i8]* @childData, i32 0, i32 0), i32 9) #4
  call void @finish(i8* getelementptr inbounds ([12 x i8], [12 x i8]* @childFinish, i32 0, i32 0), i32 11) #4
  call void @llvm.lifetime.end.p0i8(i64 100, i8* nonnull %1) #4
  ret void
}

; Function Attrs: argmemonly nounwind
declare void @llvm.lifetime.start.p0i8(i64 immarg, i8* nocapture) #1

declare i32 @getNumArguments(...) local_unnamed_addr #2

; Function Attrs: argmemonly nounwind
declare void @llvm.memcpy.p0i8.p0i8.i32(i8* nocapture writeonly, i8* nocapture readonly, i32, i1 immarg) #1

declare void @signalError(i8*, i32) local_unnamed_addr #3

; Function Attrs: argmemonly nounwind
declare void @llvm.lifetime.end.p0i8(i64 immarg, i8* nocapture) #1

declare i32 @getArgument(i32, i8*) local_unnamed_addr #3

declare i32 @getArgumentLength(i32) local_unnamed_addr #3

declare i32 @transferValue(i8*, i8*, i8*, i32) local_unnamed_addr #3

declare i32 @storageStore(i8*, i8*, i32) local_unnamed_addr #3

declare void @finish(i8*, i32) local_unnamed_addr #3

attributes #0 = { nounwind "correctly-rounded-divide-sqrt-fp-math"="false" "disable-tail-calls"="false" "less-precise-fpmad"="false" "min-legal-vector-width"="0" "no-frame-pointer-elim"="false" "no-infs-fp-math"="false" "no-jump-tables"="false" "no-nans-fp-math"="false" "no-signed-zeros-fp-math"="false" "no-trapping-math"="false" "stack-protector-buffer-size"="8" "unsafe-fp-math"="false" "use-soft-float"="false" }
attributes #1 = { argmemonly nounwind }
attributes #2 = { "correctly-rounded-divide-sqrt-fp-math"="false" "disable-tail-calls"="false" "less-precise-fpmad"="false" "no-frame-pointer-elim"="false" "no-infs-fp-math"="false" "no-nans-fp-math"="false" "no-prototype" "no-signed-zeros-fp-math"="false" "no-trapping-math"="false" "stack-protector-buffer-size"="8" "unsafe-fp-math"="false" "use-soft-float"="false" }
attributes #3 = { "correctly-rounded-divide-sqrt-fp-math"="false" "disable-tail-calls"="false" "less-precise-fpmad"="false" "no-frame-pointer-elim"="false" "no-infs-fp-math"="false" "no-nans-fp-math"="false" "no-signed-zeros-fp-math"="false" "no-trapping-math"="false" "stack-protector-buffer-size"="8" "unsafe-fp-math"="false" "use-soft-float"="false" }
attributes #4 = { nounwind }

!llvm.module.flags = !{!0}
!llvm.ident = !{!1}

!0 = !{i32 1, !"wchar_size", i32 4}
!1 = !{!"clang version 9.0.0 (tags/RELEASE_900/final)"}
