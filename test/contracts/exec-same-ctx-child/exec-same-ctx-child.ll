; ModuleID = '/var/work/Elrond/arwen-wasm-vm/test/contracts/exec-same-ctx-child/exec-same-ctx-child.c'
source_filename = "/var/work/Elrond/arwen-wasm-vm/test/contracts/exec-same-ctx-child/exec-same-ctx-child.c"
target datalayout = "e-m:e-p:32:32-i64:64-n32:64-S128"
target triple = "wasm32-unknown-unknown-wasm"

@bla = local_unnamed_addr global [32 x i8] c"bla bla bla bla bla bla bla bla\00", align 16
@dataA = global [20 x i8] zeroinitializer, align 16
@dataB = global [20 x i8] zeroinitializer, align 16
@parentKeyA = global [33 x i8] c"parentKeyA......................\00", align 16
@parentDataA = local_unnamed_addr global [12 x i8] c"parentDataA\00", align 1
@parentKeyB = global [33 x i8] c"parentKeyB......................\00", align 16
@parentDataB = local_unnamed_addr global [12 x i8] c"parentDataB\00", align 1
@childKey = global [33 x i8] c"childKey........................\00", align 16
@childData = global [10 x i8] c"childData\00", align 1
@childFinish = global [12 x i8] c"childFinish\00", align 1
@recipient = global [32 x i8] zeroinitializer, align 16
@value = global [32 x i8] c"\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00`", align 16
@__const.not_ok.msg = private unnamed_addr constant [7 x i8] c"not ok\00", align 1
@__const.childFunction.message = private unnamed_addr constant [26 x i8] c"wrong number of arguments\00", align 16
@__const.childFunction.err = private unnamed_addr constant [9 x i8] c"err lenA\00", align 1
@__const.childFunction.err.1 = private unnamed_addr constant [9 x i8] c"err lenB\00", align 1
@__const.childFunction.msg = private unnamed_addr constant [9 x i8] c"child ok\00", align 1

; Function Attrs: nounwind
define void @not_ok() local_unnamed_addr #0 {
entry:
  %msg = alloca [7 x i8], align 1
  %0 = getelementptr inbounds [7 x i8], [7 x i8]* %msg, i32 0, i32 0
  call void @llvm.lifetime.start.p0i8(i64 7, i8* nonnull %0) #4
  call void @llvm.memcpy.p0i8.p0i8.i32(i8* nonnull align 1 %0, i8* align 1 getelementptr inbounds ([7 x i8], [7 x i8]* @__const.not_ok.msg, i32 0, i32 0), i32 7, i1 false)
  call void @finish(i8* nonnull %0, i32 6) #4
  call void @llvm.lifetime.end.p0i8(i64 7, i8* nonnull %0) #4
  ret void
}

; Function Attrs: argmemonly nounwind
declare void @llvm.lifetime.start.p0i8(i64 immarg, i8* nocapture) #1

; Function Attrs: argmemonly nounwind
declare void @llvm.memcpy.p0i8.p0i8.i32(i8* nocapture writeonly, i8* nocapture readonly, i32, i1 immarg) #1

declare void @finish(i8*, i32) local_unnamed_addr #2

; Function Attrs: argmemonly nounwind
declare void @llvm.lifetime.end.p0i8(i64 immarg, i8* nocapture) #1

; Function Attrs: nounwind
define void @childFunction() local_unnamed_addr #0 {
entry:
  %msg.i = alloca [7 x i8], align 1
  %message = alloca [26 x i8], align 16
  %transferData = alloca [100 x i8], align 16
  %err = alloca [9 x i8], align 1
  %err17 = alloca [9 x i8], align 1
  %msg = alloca [9 x i8], align 1
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
  %call8 = call i32 @storageGetValueLength(i8* getelementptr inbounds ([33 x i8], [33 x i8]* @parentKeyA, i32 0, i32 0)) #4
  %cmp9 = icmp eq i32 %call8, 11
  br i1 %cmp9, label %if.end12, label %if.then10

if.then10:                                        ; preds = %if.end
  %2 = getelementptr inbounds [9 x i8], [9 x i8]* %err, i32 0, i32 0
  call void @llvm.lifetime.start.p0i8(i64 9, i8* nonnull %2) #4
  call void @llvm.memcpy.p0i8.p0i8.i32(i8* nonnull align 1 %2, i8* align 1 getelementptr inbounds ([9 x i8], [9 x i8]* @__const.childFunction.err, i32 0, i32 0), i32 9, i1 false)
  call void @finish(i8* nonnull %2, i32 8) #4
  %conv = sext i32 %call8 to i64
  call void @int64finish(i64 %conv) #4
  %3 = getelementptr inbounds [7 x i8], [7 x i8]* %msg.i, i32 0, i32 0
  call void @llvm.lifetime.start.p0i8(i64 7, i8* nonnull %3) #4
  call void @llvm.memcpy.p0i8.p0i8.i32(i8* nonnull align 1 %3, i8* align 1 getelementptr inbounds ([7 x i8], [7 x i8]* @__const.not_ok.msg, i32 0, i32 0), i32 7, i1 false) #4
  call void @finish(i8* nonnull %3, i32 6) #4
  call void @llvm.lifetime.end.p0i8(i64 7, i8* nonnull %3) #4
  call void @llvm.lifetime.end.p0i8(i64 9, i8* nonnull %2) #4
  br label %cleanup65

if.end12:                                         ; preds = %if.end
  %call13 = call i32 @storageGetValueLength(i8* getelementptr inbounds ([33 x i8], [33 x i8]* @parentKeyB, i32 0, i32 0)) #4
  %cmp14 = icmp eq i32 %call13, 11
  br i1 %cmp14, label %if.end19, label %if.then16

if.then16:                                        ; preds = %if.end12
  %4 = getelementptr inbounds [9 x i8], [9 x i8]* %err17, i32 0, i32 0
  call void @llvm.lifetime.start.p0i8(i64 9, i8* nonnull %4) #4
  call void @llvm.memcpy.p0i8.p0i8.i32(i8* nonnull align 1 %4, i8* align 1 getelementptr inbounds ([9 x i8], [9 x i8]* @__const.childFunction.err.1, i32 0, i32 0), i32 9, i1 false)
  call void @finish(i8* nonnull %4, i32 8) #4
  %5 = getelementptr inbounds [7 x i8], [7 x i8]* %msg.i, i32 0, i32 0
  call void @llvm.lifetime.start.p0i8(i64 7, i8* nonnull %5) #4
  call void @llvm.memcpy.p0i8.p0i8.i32(i8* nonnull align 1 %5, i8* align 1 getelementptr inbounds ([7 x i8], [7 x i8]* @__const.not_ok.msg, i32 0, i32 0), i32 7, i1 false) #4
  call void @finish(i8* nonnull %5, i32 6) #4
  call void @llvm.lifetime.end.p0i8(i64 7, i8* nonnull %5) #4
  call void @llvm.lifetime.end.p0i8(i64 9, i8* nonnull %4) #4
  br label %cleanup65

if.end19:                                         ; preds = %if.end12
  %call20 = call i32 @storageLoad(i8* getelementptr inbounds ([33 x i8], [33 x i8]* @parentKeyA, i32 0, i32 0), i8* getelementptr inbounds ([20 x i8], [20 x i8]* @dataA, i32 0, i32 0)) #4
  %call22 = call i32 @storageLoad(i8* getelementptr inbounds ([33 x i8], [33 x i8]* @parentKeyB, i32 0, i32 0), i8* getelementptr inbounds ([20 x i8], [20 x i8]* @dataB, i32 0, i32 0)) #4
  call void @finish(i8* getelementptr inbounds ([20 x i8], [20 x i8]* @dataA, i32 0, i32 0), i32 11) #4
  call void @finish(i8* nonnull getelementptr inbounds ([20 x i8], [20 x i8]* @dataA, i32 0, i32 0), i32 1) #4
  call void @finish(i8* nonnull getelementptr inbounds ([20 x i8], [20 x i8]* @dataA, i32 0, i32 1), i32 1) #4
  call void @finish(i8* nonnull getelementptr inbounds ([20 x i8], [20 x i8]* @dataA, i32 0, i32 2), i32 1) #4
  call void @finish(i8* nonnull getelementptr inbounds ([20 x i8], [20 x i8]* @dataA, i32 0, i32 3), i32 1) #4
  call void @finish(i8* nonnull getelementptr inbounds ([20 x i8], [20 x i8]* @dataA, i32 0, i32 4), i32 1) #4
  call void @finish(i8* nonnull getelementptr inbounds ([20 x i8], [20 x i8]* @dataA, i32 0, i32 5), i32 1) #4
  call void @finish(i8* nonnull getelementptr inbounds ([20 x i8], [20 x i8]* @dataA, i32 0, i32 6), i32 1) #4
  call void @finish(i8* nonnull getelementptr inbounds ([20 x i8], [20 x i8]* @dataA, i32 0, i32 7), i32 1) #4
  call void @finish(i8* nonnull getelementptr inbounds ([20 x i8], [20 x i8]* @dataA, i32 0, i32 8), i32 1) #4
  call void @finish(i8* nonnull getelementptr inbounds ([20 x i8], [20 x i8]* @dataA, i32 0, i32 9), i32 1) #4
  call void @finish(i8* nonnull getelementptr inbounds ([20 x i8], [20 x i8]* @dataA, i32 0, i32 10), i32 1) #4
  call void @finish(i8* getelementptr inbounds ([20 x i8], [20 x i8]* @dataB, i32 0, i32 0), i32 11) #4
  call void @finish(i8* nonnull getelementptr inbounds ([20 x i8], [20 x i8]* @dataB, i32 0, i32 0), i32 1) #4
  call void @finish(i8* nonnull getelementptr inbounds ([20 x i8], [20 x i8]* @dataB, i32 0, i32 1), i32 1) #4
  call void @finish(i8* nonnull getelementptr inbounds ([20 x i8], [20 x i8]* @dataB, i32 0, i32 2), i32 1) #4
  call void @finish(i8* nonnull getelementptr inbounds ([20 x i8], [20 x i8]* @dataB, i32 0, i32 3), i32 1) #4
  call void @finish(i8* nonnull getelementptr inbounds ([20 x i8], [20 x i8]* @dataB, i32 0, i32 4), i32 1) #4
  call void @finish(i8* nonnull getelementptr inbounds ([20 x i8], [20 x i8]* @dataB, i32 0, i32 5), i32 1) #4
  call void @finish(i8* nonnull getelementptr inbounds ([20 x i8], [20 x i8]* @dataB, i32 0, i32 6), i32 1) #4
  call void @finish(i8* nonnull getelementptr inbounds ([20 x i8], [20 x i8]* @dataB, i32 0, i32 7), i32 1) #4
  call void @finish(i8* nonnull getelementptr inbounds ([20 x i8], [20 x i8]* @dataB, i32 0, i32 8), i32 1) #4
  call void @finish(i8* nonnull getelementptr inbounds ([20 x i8], [20 x i8]* @dataB, i32 0, i32 9), i32 1) #4
  call void @finish(i8* nonnull getelementptr inbounds ([20 x i8], [20 x i8]* @dataB, i32 0, i32 10), i32 1) #4
  %6 = load i8, i8* getelementptr inbounds ([20 x i8], [20 x i8]* @dataA, i32 0, i32 0), align 16, !tbaa !2
  %7 = load i8, i8* getelementptr inbounds ([12 x i8], [12 x i8]* @parentDataA, i32 0, i32 0), align 1, !tbaa !2
  %cmp45 = icmp eq i8 %6, %7
  br i1 %cmp45, label %if.end48, label %cleanup65

for.cond37:                                       ; preds = %if.end48
  %8 = load i8, i8* getelementptr inbounds ([20 x i8], [20 x i8]* @dataA, i32 0, i32 1), align 1, !tbaa !2
  %9 = load i8, i8* getelementptr inbounds ([12 x i8], [12 x i8]* @parentDataA, i32 0, i32 1), align 1, !tbaa !2
  %cmp45.1 = icmp eq i8 %8, %9
  br i1 %cmp45.1, label %if.end48.1, label %cleanup65

if.end48:                                         ; preds = %if.end19
  %10 = load i8, i8* getelementptr inbounds ([20 x i8], [20 x i8]* @dataB, i32 0, i32 0), align 16, !tbaa !2
  %11 = load i8, i8* getelementptr inbounds ([12 x i8], [12 x i8]* @parentDataB, i32 0, i32 0), align 1, !tbaa !2
  %cmp53 = icmp eq i8 %10, %11
  br i1 %cmp53, label %for.cond37, label %cleanup65

cleanup65:                                        ; preds = %if.end19, %if.end48, %for.cond37, %if.end48.1, %for.cond37.1, %if.end48.2, %for.cond37.2, %if.end48.3, %for.cond37.3, %if.end48.4, %for.cond37.4, %if.end48.5, %for.cond37.5, %if.end48.6, %for.cond37.6, %if.end48.7, %for.cond37.7, %if.end48.8, %for.cond37.8, %if.end48.9, %for.cond37.9, %if.end48.10, %if.then16, %for.cond37.10, %if.then10
  call void @llvm.lifetime.end.p0i8(i64 100, i8* nonnull %1) #4
  ret void

if.end48.1:                                       ; preds = %for.cond37
  %12 = load i8, i8* getelementptr inbounds ([20 x i8], [20 x i8]* @dataB, i32 0, i32 1), align 1, !tbaa !2
  %13 = load i8, i8* getelementptr inbounds ([12 x i8], [12 x i8]* @parentDataB, i32 0, i32 1), align 1, !tbaa !2
  %cmp53.1 = icmp eq i8 %12, %13
  br i1 %cmp53.1, label %for.cond37.1, label %cleanup65

for.cond37.1:                                     ; preds = %if.end48.1
  %14 = load i8, i8* getelementptr inbounds ([20 x i8], [20 x i8]* @dataA, i32 0, i32 2), align 2, !tbaa !2
  %15 = load i8, i8* getelementptr inbounds ([12 x i8], [12 x i8]* @parentDataA, i32 0, i32 2), align 1, !tbaa !2
  %cmp45.2 = icmp eq i8 %14, %15
  br i1 %cmp45.2, label %if.end48.2, label %cleanup65

if.end48.2:                                       ; preds = %for.cond37.1
  %16 = load i8, i8* getelementptr inbounds ([20 x i8], [20 x i8]* @dataB, i32 0, i32 2), align 2, !tbaa !2
  %17 = load i8, i8* getelementptr inbounds ([12 x i8], [12 x i8]* @parentDataB, i32 0, i32 2), align 1, !tbaa !2
  %cmp53.2 = icmp eq i8 %16, %17
  br i1 %cmp53.2, label %for.cond37.2, label %cleanup65

for.cond37.2:                                     ; preds = %if.end48.2
  %18 = load i8, i8* getelementptr inbounds ([20 x i8], [20 x i8]* @dataA, i32 0, i32 3), align 1, !tbaa !2
  %19 = load i8, i8* getelementptr inbounds ([12 x i8], [12 x i8]* @parentDataA, i32 0, i32 3), align 1, !tbaa !2
  %cmp45.3 = icmp eq i8 %18, %19
  br i1 %cmp45.3, label %if.end48.3, label %cleanup65

if.end48.3:                                       ; preds = %for.cond37.2
  %20 = load i8, i8* getelementptr inbounds ([20 x i8], [20 x i8]* @dataB, i32 0, i32 3), align 1, !tbaa !2
  %21 = load i8, i8* getelementptr inbounds ([12 x i8], [12 x i8]* @parentDataB, i32 0, i32 3), align 1, !tbaa !2
  %cmp53.3 = icmp eq i8 %20, %21
  br i1 %cmp53.3, label %for.cond37.3, label %cleanup65

for.cond37.3:                                     ; preds = %if.end48.3
  %22 = load i8, i8* getelementptr inbounds ([20 x i8], [20 x i8]* @dataA, i32 0, i32 4), align 4, !tbaa !2
  %23 = load i8, i8* getelementptr inbounds ([12 x i8], [12 x i8]* @parentDataA, i32 0, i32 4), align 1, !tbaa !2
  %cmp45.4 = icmp eq i8 %22, %23
  br i1 %cmp45.4, label %if.end48.4, label %cleanup65

if.end48.4:                                       ; preds = %for.cond37.3
  %24 = load i8, i8* getelementptr inbounds ([20 x i8], [20 x i8]* @dataB, i32 0, i32 4), align 4, !tbaa !2
  %25 = load i8, i8* getelementptr inbounds ([12 x i8], [12 x i8]* @parentDataB, i32 0, i32 4), align 1, !tbaa !2
  %cmp53.4 = icmp eq i8 %24, %25
  br i1 %cmp53.4, label %for.cond37.4, label %cleanup65

for.cond37.4:                                     ; preds = %if.end48.4
  %26 = load i8, i8* getelementptr inbounds ([20 x i8], [20 x i8]* @dataA, i32 0, i32 5), align 1, !tbaa !2
  %27 = load i8, i8* getelementptr inbounds ([12 x i8], [12 x i8]* @parentDataA, i32 0, i32 5), align 1, !tbaa !2
  %cmp45.5 = icmp eq i8 %26, %27
  br i1 %cmp45.5, label %if.end48.5, label %cleanup65

if.end48.5:                                       ; preds = %for.cond37.4
  %28 = load i8, i8* getelementptr inbounds ([20 x i8], [20 x i8]* @dataB, i32 0, i32 5), align 1, !tbaa !2
  %29 = load i8, i8* getelementptr inbounds ([12 x i8], [12 x i8]* @parentDataB, i32 0, i32 5), align 1, !tbaa !2
  %cmp53.5 = icmp eq i8 %28, %29
  br i1 %cmp53.5, label %for.cond37.5, label %cleanup65

for.cond37.5:                                     ; preds = %if.end48.5
  %30 = load i8, i8* getelementptr inbounds ([20 x i8], [20 x i8]* @dataA, i32 0, i32 6), align 2, !tbaa !2
  %31 = load i8, i8* getelementptr inbounds ([12 x i8], [12 x i8]* @parentDataA, i32 0, i32 6), align 1, !tbaa !2
  %cmp45.6 = icmp eq i8 %30, %31
  br i1 %cmp45.6, label %if.end48.6, label %cleanup65

if.end48.6:                                       ; preds = %for.cond37.5
  %32 = load i8, i8* getelementptr inbounds ([20 x i8], [20 x i8]* @dataB, i32 0, i32 6), align 2, !tbaa !2
  %33 = load i8, i8* getelementptr inbounds ([12 x i8], [12 x i8]* @parentDataB, i32 0, i32 6), align 1, !tbaa !2
  %cmp53.6 = icmp eq i8 %32, %33
  br i1 %cmp53.6, label %for.cond37.6, label %cleanup65

for.cond37.6:                                     ; preds = %if.end48.6
  %34 = load i8, i8* getelementptr inbounds ([20 x i8], [20 x i8]* @dataA, i32 0, i32 7), align 1, !tbaa !2
  %35 = load i8, i8* getelementptr inbounds ([12 x i8], [12 x i8]* @parentDataA, i32 0, i32 7), align 1, !tbaa !2
  %cmp45.7 = icmp eq i8 %34, %35
  br i1 %cmp45.7, label %if.end48.7, label %cleanup65

if.end48.7:                                       ; preds = %for.cond37.6
  %36 = load i8, i8* getelementptr inbounds ([20 x i8], [20 x i8]* @dataB, i32 0, i32 7), align 1, !tbaa !2
  %37 = load i8, i8* getelementptr inbounds ([12 x i8], [12 x i8]* @parentDataB, i32 0, i32 7), align 1, !tbaa !2
  %cmp53.7 = icmp eq i8 %36, %37
  br i1 %cmp53.7, label %for.cond37.7, label %cleanup65

for.cond37.7:                                     ; preds = %if.end48.7
  %38 = load i8, i8* getelementptr inbounds ([20 x i8], [20 x i8]* @dataA, i32 0, i32 8), align 8, !tbaa !2
  %39 = load i8, i8* getelementptr inbounds ([12 x i8], [12 x i8]* @parentDataA, i32 0, i32 8), align 1, !tbaa !2
  %cmp45.8 = icmp eq i8 %38, %39
  br i1 %cmp45.8, label %if.end48.8, label %cleanup65

if.end48.8:                                       ; preds = %for.cond37.7
  %40 = load i8, i8* getelementptr inbounds ([20 x i8], [20 x i8]* @dataB, i32 0, i32 8), align 8, !tbaa !2
  %41 = load i8, i8* getelementptr inbounds ([12 x i8], [12 x i8]* @parentDataB, i32 0, i32 8), align 1, !tbaa !2
  %cmp53.8 = icmp eq i8 %40, %41
  br i1 %cmp53.8, label %for.cond37.8, label %cleanup65

for.cond37.8:                                     ; preds = %if.end48.8
  %42 = load i8, i8* getelementptr inbounds ([20 x i8], [20 x i8]* @dataA, i32 0, i32 9), align 1, !tbaa !2
  %43 = load i8, i8* getelementptr inbounds ([12 x i8], [12 x i8]* @parentDataA, i32 0, i32 9), align 1, !tbaa !2
  %cmp45.9 = icmp eq i8 %42, %43
  br i1 %cmp45.9, label %if.end48.9, label %cleanup65

if.end48.9:                                       ; preds = %for.cond37.8
  %44 = load i8, i8* getelementptr inbounds ([20 x i8], [20 x i8]* @dataB, i32 0, i32 9), align 1, !tbaa !2
  %45 = load i8, i8* getelementptr inbounds ([12 x i8], [12 x i8]* @parentDataB, i32 0, i32 9), align 1, !tbaa !2
  %cmp53.9 = icmp eq i8 %44, %45
  br i1 %cmp53.9, label %for.cond37.9, label %cleanup65

for.cond37.9:                                     ; preds = %if.end48.9
  %46 = load i8, i8* getelementptr inbounds ([20 x i8], [20 x i8]* @dataA, i32 0, i32 10), align 2, !tbaa !2
  %47 = load i8, i8* getelementptr inbounds ([12 x i8], [12 x i8]* @parentDataA, i32 0, i32 10), align 1, !tbaa !2
  %cmp45.10 = icmp eq i8 %46, %47
  br i1 %cmp45.10, label %if.end48.10, label %cleanup65

if.end48.10:                                      ; preds = %for.cond37.9
  %48 = load i8, i8* getelementptr inbounds ([20 x i8], [20 x i8]* @dataB, i32 0, i32 10), align 2, !tbaa !2
  %49 = load i8, i8* getelementptr inbounds ([12 x i8], [12 x i8]* @parentDataB, i32 0, i32 10), align 1, !tbaa !2
  %cmp53.10 = icmp eq i8 %48, %49
  br i1 %cmp53.10, label %for.cond37.10, label %cleanup65

for.cond37.10:                                    ; preds = %if.end48.10
  %50 = getelementptr inbounds [9 x i8], [9 x i8]* %msg, i32 0, i32 0
  call void @llvm.lifetime.start.p0i8(i64 9, i8* nonnull %50) #4
  call void @llvm.memcpy.p0i8.p0i8.i32(i8* nonnull align 1 %50, i8* align 1 getelementptr inbounds ([9 x i8], [9 x i8]* @__const.childFunction.msg, i32 0, i32 0), i32 9, i1 false)
  call void @finish(i8* nonnull %50, i32 8) #4
  call void @llvm.lifetime.end.p0i8(i64 9, i8* nonnull %50) #4
  br label %cleanup65
}

declare i32 @getNumArguments(...) local_unnamed_addr #3

declare void @signalError(i8*, i32) local_unnamed_addr #2

declare i32 @getArgument(i32, i8*) local_unnamed_addr #2

declare i32 @getArgumentLength(i32) local_unnamed_addr #2

declare i32 @transferValue(i8*, i8*, i8*, i32) local_unnamed_addr #2

declare i32 @storageStore(i8*, i8*, i32) local_unnamed_addr #2

declare i32 @storageGetValueLength(i8*) local_unnamed_addr #2

declare void @int64finish(i64) local_unnamed_addr #2

declare i32 @storageLoad(i8*, i8*) local_unnamed_addr #2

attributes #0 = { nounwind "correctly-rounded-divide-sqrt-fp-math"="false" "disable-tail-calls"="false" "less-precise-fpmad"="false" "min-legal-vector-width"="0" "no-frame-pointer-elim"="false" "no-infs-fp-math"="false" "no-jump-tables"="false" "no-nans-fp-math"="false" "no-signed-zeros-fp-math"="false" "no-trapping-math"="false" "stack-protector-buffer-size"="8" "unsafe-fp-math"="false" "use-soft-float"="false" }
attributes #1 = { argmemonly nounwind }
attributes #2 = { "correctly-rounded-divide-sqrt-fp-math"="false" "disable-tail-calls"="false" "less-precise-fpmad"="false" "no-frame-pointer-elim"="false" "no-infs-fp-math"="false" "no-nans-fp-math"="false" "no-signed-zeros-fp-math"="false" "no-trapping-math"="false" "stack-protector-buffer-size"="8" "unsafe-fp-math"="false" "use-soft-float"="false" }
attributes #3 = { "correctly-rounded-divide-sqrt-fp-math"="false" "disable-tail-calls"="false" "less-precise-fpmad"="false" "no-frame-pointer-elim"="false" "no-infs-fp-math"="false" "no-nans-fp-math"="false" "no-prototype" "no-signed-zeros-fp-math"="false" "no-trapping-math"="false" "stack-protector-buffer-size"="8" "unsafe-fp-math"="false" "use-soft-float"="false" }
attributes #4 = { nounwind }

!llvm.module.flags = !{!0}
!llvm.ident = !{!1}

!0 = !{i32 1, !"wchar_size", i32 4}
!1 = !{!"clang version 9.0.0 (tags/RELEASE_900/final)"}
!2 = !{!3, !3, i64 0}
!3 = !{!"omnipotent char", !4, i64 0}
!4 = !{!"Simple C/C++ TBAA"}
