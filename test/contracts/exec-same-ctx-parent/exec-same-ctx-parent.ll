; ModuleID = '/home/camil.bancioiu/Work/Elrond/arwen-wasm-vm/test/contracts/exec-same-ctx-parent/exec-same-ctx-parent.c'
source_filename = "/home/camil.bancioiu/Work/Elrond/arwen-wasm-vm/test/contracts/exec-same-ctx-parent/exec-same-ctx-parent.c"
target datalayout = "e-m:e-p:32:32-i64:64-n32:64-S128"
target triple = "wasm32-unknown-unknown-wasm"

@parentKeyA = global [33 x i8] c"parentKeyA......................\00", align 16
@parentDataA = global [12 x i8] c"parentDataA\00", align 1
@parentKeyB = global [33 x i8] c"parentKeyB......................\00", align 16
@parentDataB = global [12 x i8] c"parentDataB\00", align 1
@parentFinishA = global [14 x i8] c"parentFinishA\00", align 1
@parentFinishB = global [14 x i8] c"parentFinishB\00", align 1
@parentTransferReceiver = global [33 x i8] c"parentTransferReceiver..........\00", align 16
@parentTransferValue = global [32 x i8] c"\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00*", align 16
@parentTransferData = global [19 x i8] c"parentTransferData\00", align 16
@executeValue = global [32 x i8] c"\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00c", align 16
@executeArgumentsLengths = global [2 x i32] [i32 32, i32 6], align 4
@executeArgumentsData = global [39 x i8] c"asdfoottxxwlllllllllllwrraatttttqwerty\00", align 16
@__const.parentFunctionWrongCall.childAddress = private unnamed_addr constant [33 x i8] c"wrongSC.........................\00", align 16
@__const.parentFunctionChildCall.childAddress = private unnamed_addr constant [33 x i8] c"secondSC........................\00", align 16
@__const.parentFunctionChildCall.functionName = private unnamed_addr constant [14 x i8] c"childFunction\00", align 1
@__const.finishResult.message.1 = private unnamed_addr constant [7 x i8] c"failed\00", align 1
@__const.finishResult.message.2 = private unnamed_addr constant [15 x i8] c"unknown result\00", align 1

; Function Attrs: nounwind
define void @parentFunctionPrepare() local_unnamed_addr #0 {
entry:
  %call = tail call i32 @storageStore(i8* getelementptr inbounds ([33 x i8], [33 x i8]* @parentKeyA, i32 0, i32 0), i8* getelementptr inbounds ([12 x i8], [12 x i8]* @parentDataA, i32 0, i32 0), i32 11) #3
  %call1 = tail call i32 @storageStore(i8* getelementptr inbounds ([33 x i8], [33 x i8]* @parentKeyB, i32 0, i32 0), i8* getelementptr inbounds ([12 x i8], [12 x i8]* @parentDataB, i32 0, i32 0), i32 11) #3
  tail call void @finish(i8* getelementptr inbounds ([14 x i8], [14 x i8]* @parentFinishA, i32 0, i32 0), i32 13) #3
  tail call void @finish(i8* getelementptr inbounds ([14 x i8], [14 x i8]* @parentFinishB, i32 0, i32 0), i32 13) #3
  %call2 = tail call i32 @transferValue(i8* getelementptr inbounds ([33 x i8], [33 x i8]* @parentTransferReceiver, i32 0, i32 0), i8* getelementptr inbounds ([32 x i8], [32 x i8]* @parentTransferValue, i32 0, i32 0), i8* getelementptr inbounds ([19 x i8], [19 x i8]* @parentTransferData, i32 0, i32 0), i32 18) #3
  ret void
}

declare i32 @storageStore(i8*, i8*, i32) local_unnamed_addr #1

declare void @finish(i8*, i32) local_unnamed_addr #1

declare i32 @transferValue(i8*, i8*, i8*, i32) local_unnamed_addr #1

; Function Attrs: nounwind
define void @parentFunctionWrongCall() local_unnamed_addr #0 {
entry:
  %message.i = alloca i64, align 8
  %message3.i = alloca [7 x i8], align 1
  %message9.i = alloca [15 x i8], align 1
  %childAddress = alloca [33 x i8], align 16
  %functionName = alloca [14 x i8], align 1
  %call.i = tail call i32 @storageStore(i8* getelementptr inbounds ([33 x i8], [33 x i8]* @parentKeyA, i32 0, i32 0), i8* getelementptr inbounds ([12 x i8], [12 x i8]* @parentDataA, i32 0, i32 0), i32 11) #3
  %call1.i = tail call i32 @storageStore(i8* getelementptr inbounds ([33 x i8], [33 x i8]* @parentKeyB, i32 0, i32 0), i8* getelementptr inbounds ([12 x i8], [12 x i8]* @parentDataB, i32 0, i32 0), i32 11) #3
  tail call void @finish(i8* getelementptr inbounds ([14 x i8], [14 x i8]* @parentFinishA, i32 0, i32 0), i32 13) #3
  tail call void @finish(i8* getelementptr inbounds ([14 x i8], [14 x i8]* @parentFinishB, i32 0, i32 0), i32 13) #3
  %call2.i = tail call i32 @transferValue(i8* getelementptr inbounds ([33 x i8], [33 x i8]* @parentTransferReceiver, i32 0, i32 0), i8* getelementptr inbounds ([32 x i8], [32 x i8]* @parentTransferValue, i32 0, i32 0), i8* getelementptr inbounds ([19 x i8], [19 x i8]* @parentTransferData, i32 0, i32 0), i32 18) #3
  %0 = getelementptr inbounds [33 x i8], [33 x i8]* %childAddress, i32 0, i32 0
  call void @llvm.lifetime.start.p0i8(i64 33, i8* nonnull %0) #3
  call void @llvm.memcpy.p0i8.p0i8.i32(i8* nonnull align 16 %0, i8* align 16 getelementptr inbounds ([33 x i8], [33 x i8]* @__const.parentFunctionWrongCall.childAddress, i32 0, i32 0), i32 33, i1 false)
  %1 = getelementptr inbounds [14 x i8], [14 x i8]* %functionName, i32 0, i32 0
  call void @llvm.lifetime.start.p0i8(i64 14, i8* nonnull %1) #3
  call void @llvm.memcpy.p0i8.p0i8.i32(i8* nonnull align 1 %1, i8* align 1 getelementptr inbounds ([14 x i8], [14 x i8]* @__const.parentFunctionChildCall.functionName, i32 0, i32 0), i32 14, i1 false)
  %call = call i32 @executeOnSameContext(i64 10000, i8* nonnull %0, i8* getelementptr inbounds ([32 x i8], [32 x i8]* @executeValue, i32 0, i32 0), i8* nonnull %1, i32 13, i32 2, i8* bitcast ([2 x i32]* @executeArgumentsLengths to i8*), i8* getelementptr inbounds ([39 x i8], [39 x i8]* @executeArgumentsData, i32 0, i32 0)) #3
  switch i32 %call, label %if.then8.i [
    i32 0, label %if.then.i
    i32 1, label %if.then2.i
  ]

if.then.i:                                        ; preds = %entry
  %2 = bitcast i64* %message.i to i8*
  call void @llvm.lifetime.start.p0i8(i64 8, i8* nonnull %2) #3
  store i64 32496501618079091, i64* %message.i, align 8
  call void @finish(i8* nonnull %2, i32 7) #3
  call void @llvm.lifetime.end.p0i8(i64 8, i8* nonnull %2) #3
  br label %finishResult.exit

if.then2.i:                                       ; preds = %entry
  %3 = getelementptr inbounds [7 x i8], [7 x i8]* %message3.i, i32 0, i32 0
  call void @llvm.lifetime.start.p0i8(i64 7, i8* nonnull %3) #3
  call void @llvm.memcpy.p0i8.p0i8.i32(i8* nonnull align 1 %3, i8* align 1 getelementptr inbounds ([7 x i8], [7 x i8]* @__const.finishResult.message.1, i32 0, i32 0), i32 7, i1 false) #3
  call void @finish(i8* nonnull %3, i32 6) #3
  call void @llvm.lifetime.end.p0i8(i64 7, i8* nonnull %3) #3
  br label %finishResult.exit

if.then8.i:                                       ; preds = %entry
  %4 = getelementptr inbounds [15 x i8], [15 x i8]* %message9.i, i32 0, i32 0
  call void @llvm.lifetime.start.p0i8(i64 15, i8* nonnull %4) #3
  call void @llvm.memcpy.p0i8.p0i8.i32(i8* nonnull align 1 %4, i8* align 1 getelementptr inbounds ([15 x i8], [15 x i8]* @__const.finishResult.message.2, i32 0, i32 0), i32 15, i1 false) #3
  call void @finish(i8* nonnull %4, i32 14) #3
  call void @llvm.lifetime.end.p0i8(i64 15, i8* nonnull %4) #3
  br label %finishResult.exit

finishResult.exit:                                ; preds = %if.then.i, %if.then2.i, %if.then8.i
  call void @llvm.lifetime.end.p0i8(i64 14, i8* nonnull %1) #3
  call void @llvm.lifetime.end.p0i8(i64 33, i8* nonnull %0) #3
  ret void
}

; Function Attrs: argmemonly nounwind
declare void @llvm.lifetime.start.p0i8(i64 immarg, i8* nocapture) #2

; Function Attrs: argmemonly nounwind
declare void @llvm.memcpy.p0i8.p0i8.i32(i8* nocapture writeonly, i8* nocapture readonly, i32, i1 immarg) #2

declare i32 @executeOnSameContext(i64, i8*, i8*, i8*, i32, i32, i8*, i8*) local_unnamed_addr #1

; Function Attrs: nounwind
define void @finishResult(i32 %result) local_unnamed_addr #0 {
entry:
  %message = alloca i64, align 8
  %message3 = alloca [7 x i8], align 1
  %message9 = alloca [15 x i8], align 1
  switch i32 %result, label %if.then8 [
    i32 0, label %if.then
    i32 1, label %if.then2
  ]

if.then:                                          ; preds = %entry
  %0 = bitcast i64* %message to i8*
  call void @llvm.lifetime.start.p0i8(i64 8, i8* nonnull %0) #3
  store i64 32496501618079091, i64* %message, align 8
  call void @finish(i8* nonnull %0, i32 7) #3
  call void @llvm.lifetime.end.p0i8(i64 8, i8* nonnull %0) #3
  br label %if.end11

if.then2:                                         ; preds = %entry
  %1 = getelementptr inbounds [7 x i8], [7 x i8]* %message3, i32 0, i32 0
  call void @llvm.lifetime.start.p0i8(i64 7, i8* nonnull %1) #3
  call void @llvm.memcpy.p0i8.p0i8.i32(i8* nonnull align 1 %1, i8* align 1 getelementptr inbounds ([7 x i8], [7 x i8]* @__const.finishResult.message.1, i32 0, i32 0), i32 7, i1 false)
  call void @finish(i8* nonnull %1, i32 6) #3
  call void @llvm.lifetime.end.p0i8(i64 7, i8* nonnull %1) #3
  br label %if.end11

if.then8:                                         ; preds = %entry
  %2 = getelementptr inbounds [15 x i8], [15 x i8]* %message9, i32 0, i32 0
  call void @llvm.lifetime.start.p0i8(i64 15, i8* nonnull %2) #3
  call void @llvm.memcpy.p0i8.p0i8.i32(i8* nonnull align 1 %2, i8* align 1 getelementptr inbounds ([15 x i8], [15 x i8]* @__const.finishResult.message.2, i32 0, i32 0), i32 15, i1 false)
  call void @finish(i8* nonnull %2, i32 14) #3
  call void @llvm.lifetime.end.p0i8(i64 15, i8* nonnull %2) #3
  br label %if.end11

if.end11:                                         ; preds = %if.then2, %if.then, %if.then8
  ret void
}

; Function Attrs: argmemonly nounwind
declare void @llvm.lifetime.end.p0i8(i64 immarg, i8* nocapture) #2

; Function Attrs: nounwind
define void @parentFunctionChildCall() local_unnamed_addr #0 {
entry:
  %message.i = alloca i64, align 8
  %message3.i = alloca [7 x i8], align 1
  %message9.i = alloca [15 x i8], align 1
  %childAddress = alloca [33 x i8], align 16
  %functionName = alloca [14 x i8], align 1
  %call.i = tail call i32 @storageStore(i8* getelementptr inbounds ([33 x i8], [33 x i8]* @parentKeyA, i32 0, i32 0), i8* getelementptr inbounds ([12 x i8], [12 x i8]* @parentDataA, i32 0, i32 0), i32 11) #3
  %call1.i = tail call i32 @storageStore(i8* getelementptr inbounds ([33 x i8], [33 x i8]* @parentKeyB, i32 0, i32 0), i8* getelementptr inbounds ([12 x i8], [12 x i8]* @parentDataB, i32 0, i32 0), i32 11) #3
  tail call void @finish(i8* getelementptr inbounds ([14 x i8], [14 x i8]* @parentFinishA, i32 0, i32 0), i32 13) #3
  tail call void @finish(i8* getelementptr inbounds ([14 x i8], [14 x i8]* @parentFinishB, i32 0, i32 0), i32 13) #3
  %call2.i = tail call i32 @transferValue(i8* getelementptr inbounds ([33 x i8], [33 x i8]* @parentTransferReceiver, i32 0, i32 0), i8* getelementptr inbounds ([32 x i8], [32 x i8]* @parentTransferValue, i32 0, i32 0), i8* getelementptr inbounds ([19 x i8], [19 x i8]* @parentTransferData, i32 0, i32 0), i32 18) #3
  %0 = getelementptr inbounds [33 x i8], [33 x i8]* %childAddress, i32 0, i32 0
  call void @llvm.lifetime.start.p0i8(i64 33, i8* nonnull %0) #3
  call void @llvm.memcpy.p0i8.p0i8.i32(i8* nonnull align 16 %0, i8* align 16 getelementptr inbounds ([33 x i8], [33 x i8]* @__const.parentFunctionChildCall.childAddress, i32 0, i32 0), i32 33, i1 false)
  %1 = getelementptr inbounds [14 x i8], [14 x i8]* %functionName, i32 0, i32 0
  call void @llvm.lifetime.start.p0i8(i64 14, i8* nonnull %1) #3
  call void @llvm.memcpy.p0i8.p0i8.i32(i8* nonnull align 1 %1, i8* align 1 getelementptr inbounds ([14 x i8], [14 x i8]* @__const.parentFunctionChildCall.functionName, i32 0, i32 0), i32 14, i1 false)
  %call = call i32 @executeOnSameContext(i64 20000, i8* nonnull %0, i8* getelementptr inbounds ([32 x i8], [32 x i8]* @executeValue, i32 0, i32 0), i8* nonnull %1, i32 13, i32 2, i8* bitcast ([2 x i32]* @executeArgumentsLengths to i8*), i8* getelementptr inbounds ([39 x i8], [39 x i8]* @executeArgumentsData, i32 0, i32 0)) #3
  switch i32 %call, label %if.then8.i [
    i32 0, label %if.then.i
    i32 1, label %if.then2.i
  ]

if.then.i:                                        ; preds = %entry
  %2 = bitcast i64* %message.i to i8*
  call void @llvm.lifetime.start.p0i8(i64 8, i8* nonnull %2) #3
  store i64 32496501618079091, i64* %message.i, align 8
  call void @finish(i8* nonnull %2, i32 7) #3
  call void @llvm.lifetime.end.p0i8(i64 8, i8* nonnull %2) #3
  br label %finishResult.exit

if.then2.i:                                       ; preds = %entry
  %3 = getelementptr inbounds [7 x i8], [7 x i8]* %message3.i, i32 0, i32 0
  call void @llvm.lifetime.start.p0i8(i64 7, i8* nonnull %3) #3
  call void @llvm.memcpy.p0i8.p0i8.i32(i8* nonnull align 1 %3, i8* align 1 getelementptr inbounds ([7 x i8], [7 x i8]* @__const.finishResult.message.1, i32 0, i32 0), i32 7, i1 false) #3
  call void @finish(i8* nonnull %3, i32 6) #3
  call void @llvm.lifetime.end.p0i8(i64 7, i8* nonnull %3) #3
  br label %finishResult.exit

if.then8.i:                                       ; preds = %entry
  %4 = getelementptr inbounds [15 x i8], [15 x i8]* %message9.i, i32 0, i32 0
  call void @llvm.lifetime.start.p0i8(i64 15, i8* nonnull %4) #3
  call void @llvm.memcpy.p0i8.p0i8.i32(i8* nonnull align 1 %4, i8* align 1 getelementptr inbounds ([15 x i8], [15 x i8]* @__const.finishResult.message.2, i32 0, i32 0), i32 15, i1 false) #3
  call void @finish(i8* nonnull %4, i32 14) #3
  call void @llvm.lifetime.end.p0i8(i64 15, i8* nonnull %4) #3
  br label %finishResult.exit

finishResult.exit:                                ; preds = %if.then.i, %if.then2.i, %if.then8.i
  call void @llvm.lifetime.end.p0i8(i64 14, i8* nonnull %1) #3
  call void @llvm.lifetime.end.p0i8(i64 33, i8* nonnull %0) #3
  ret void
}

attributes #0 = { nounwind "correctly-rounded-divide-sqrt-fp-math"="false" "disable-tail-calls"="false" "less-precise-fpmad"="false" "min-legal-vector-width"="0" "no-frame-pointer-elim"="false" "no-infs-fp-math"="false" "no-jump-tables"="false" "no-nans-fp-math"="false" "no-signed-zeros-fp-math"="false" "no-trapping-math"="false" "stack-protector-buffer-size"="8" "unsafe-fp-math"="false" "use-soft-float"="false" }
attributes #1 = { "correctly-rounded-divide-sqrt-fp-math"="false" "disable-tail-calls"="false" "less-precise-fpmad"="false" "no-frame-pointer-elim"="false" "no-infs-fp-math"="false" "no-nans-fp-math"="false" "no-signed-zeros-fp-math"="false" "no-trapping-math"="false" "stack-protector-buffer-size"="8" "unsafe-fp-math"="false" "use-soft-float"="false" }
attributes #2 = { argmemonly nounwind }
attributes #3 = { nounwind }

!llvm.module.flags = !{!0}
!llvm.ident = !{!1}

!0 = !{i32 1, !"wchar_size", i32 4}
!1 = !{!"clang version 9.0.0 (tags/RELEASE_900/final)"}
