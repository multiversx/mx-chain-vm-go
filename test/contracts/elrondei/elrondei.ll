; ModuleID = '/home/camil.bancioiu/Work/Elrond/arwen-wasm-vm/test/contracts/elrondei/elrondei.c'
source_filename = "/home/camil.bancioiu/Work/Elrond/arwen-wasm-vm/test/contracts/elrondei/elrondei.c"
target datalayout = "e-m:e-p:32:32-i64:64-n32:64-S128"
target triple = "wasm32-unknown-unknown-wasm"

@msg_ok = global [3 x i8] c"ok\00", align 1
@msg_not_ok = global [7 x i8] c"not ok\00", align 1
@msg_unexpected = global [11 x i8] c"unexpected\00", align 1
@value = global [32 x i8] zeroinitializer, align 16

; Function Attrs: nounwind
define void @test_getCallValue_1byte() local_unnamed_addr #0 {
entry:
  %length = alloca i32, align 4
  %0 = bitcast i32* %length to i8*
  call void @llvm.lifetime.start.p0i8(i64 4, i8* nonnull %0) #3
  %call = tail call i32 @getCallValue(i8* getelementptr inbounds ([32 x i8], [32 x i8]* @value, i32 0, i32 0)) #3
  store i32 %call, i32* %length, align 4, !tbaa !2
  %cmp = icmp eq i32 %call, 1
  br i1 %cmp, label %if.end, label %if.then

if.then:                                          ; preds = %entry
  tail call void @signalError(i8* getelementptr inbounds ([11 x i8], [11 x i8]* @msg_unexpected, i32 0, i32 0), i32 10) #3
  br label %if.end

if.end:                                           ; preds = %entry, %if.then
  %1 = load i8, i8* getelementptr inbounds ([32 x i8], [32 x i8]* @value, i32 0, i32 0), align 16, !tbaa !6
  %cmp1 = icmp eq i8 %1, 64
  br i1 %cmp1, label %if.then3, label %if.else

if.then3:                                         ; preds = %if.end
  tail call void @finish(i8* getelementptr inbounds ([3 x i8], [3 x i8]* @msg_ok, i32 0, i32 0), i32 2) #3
  br label %if.end4

if.else:                                          ; preds = %if.end
  tail call void @finish(i8* getelementptr inbounds ([7 x i8], [7 x i8]* @msg_not_ok, i32 0, i32 0), i32 6) #3
  br label %if.end4

if.end4:                                          ; preds = %if.else, %if.then3
  call void @finish(i8* nonnull %0, i32 4) #3
  %2 = load i32, i32* %length, align 4, !tbaa !2
  call void @finish(i8* getelementptr inbounds ([32 x i8], [32 x i8]* @value, i32 0, i32 0), i32 %2) #3
  call void @llvm.lifetime.end.p0i8(i64 4, i8* nonnull %0) #3
  ret void
}

; Function Attrs: argmemonly nounwind
declare void @llvm.lifetime.start.p0i8(i64 immarg, i8* nocapture) #1

declare i32 @getCallValue(i8*) local_unnamed_addr #2

declare void @signalError(i8*, i32) local_unnamed_addr #2

declare void @finish(i8*, i32) local_unnamed_addr #2

; Function Attrs: argmemonly nounwind
declare void @llvm.lifetime.end.p0i8(i64 immarg, i8* nocapture) #1

; Function Attrs: nounwind
define void @test_getCallValue_4bytes() local_unnamed_addr #0 {
entry:
  %length = alloca i32, align 4
  %0 = bitcast i32* %length to i8*
  call void @llvm.lifetime.start.p0i8(i64 4, i8* nonnull %0) #3
  %call = tail call i32 @getCallValue(i8* getelementptr inbounds ([32 x i8], [32 x i8]* @value, i32 0, i32 0)) #3
  store i32 %call, i32* %length, align 4, !tbaa !2
  %cmp = icmp eq i32 %call, 4
  br i1 %cmp, label %if.end, label %if.then

if.then:                                          ; preds = %entry
  tail call void @signalError(i8* getelementptr inbounds ([11 x i8], [11 x i8]* @msg_unexpected, i32 0, i32 0), i32 10) #3
  br label %if.end

if.end:                                           ; preds = %entry, %if.then
  %1 = load i8, i8* getelementptr inbounds ([32 x i8], [32 x i8]* @value, i32 0, i32 0), align 16, !tbaa !6
  %cmp1 = icmp eq i8 %1, 64
  %conv2 = zext i1 %cmp1 to i32
  %2 = load i8, i8* getelementptr inbounds ([32 x i8], [32 x i8]* @value, i32 0, i32 1), align 1, !tbaa !6
  %cmp4 = icmp eq i8 %2, 12
  %conv5 = zext i1 %cmp4 to i32
  %add6 = add nuw nsw i32 %conv5, %conv2
  %3 = load i8, i8* getelementptr inbounds ([32 x i8], [32 x i8]* @value, i32 0, i32 2), align 2, !tbaa !6
  %cmp8 = icmp eq i8 %3, 16
  %conv9 = zext i1 %cmp8 to i32
  %add10 = add nuw nsw i32 %add6, %conv9
  %4 = load i8, i8* getelementptr inbounds ([32 x i8], [32 x i8]* @value, i32 0, i32 3), align 1, !tbaa !6
  %cmp12 = icmp eq i8 %4, 99
  %conv13 = zext i1 %cmp12 to i32
  %add14 = add nuw nsw i32 %add10, %conv13
  %cmp15 = icmp eq i32 %add14, 4
  br i1 %cmp15, label %if.then17, label %if.else

if.then17:                                        ; preds = %if.end
  tail call void @finish(i8* getelementptr inbounds ([3 x i8], [3 x i8]* @msg_ok, i32 0, i32 0), i32 2) #3
  br label %if.end18

if.else:                                          ; preds = %if.end
  tail call void @finish(i8* getelementptr inbounds ([7 x i8], [7 x i8]* @msg_not_ok, i32 0, i32 0), i32 6) #3
  br label %if.end18

if.end18:                                         ; preds = %if.else, %if.then17
  call void @finish(i8* nonnull %0, i32 4) #3
  %5 = load i32, i32* %length, align 4, !tbaa !2
  call void @finish(i8* getelementptr inbounds ([32 x i8], [32 x i8]* @value, i32 0, i32 0), i32 %5) #3
  call void @llvm.lifetime.end.p0i8(i64 4, i8* nonnull %0) #3
  ret void
}

attributes #0 = { nounwind "correctly-rounded-divide-sqrt-fp-math"="false" "disable-tail-calls"="false" "less-precise-fpmad"="false" "min-legal-vector-width"="0" "no-frame-pointer-elim"="false" "no-infs-fp-math"="false" "no-jump-tables"="false" "no-nans-fp-math"="false" "no-signed-zeros-fp-math"="false" "no-trapping-math"="false" "stack-protector-buffer-size"="8" "unsafe-fp-math"="false" "use-soft-float"="false" }
attributes #1 = { argmemonly nounwind }
attributes #2 = { "correctly-rounded-divide-sqrt-fp-math"="false" "disable-tail-calls"="false" "less-precise-fpmad"="false" "no-frame-pointer-elim"="false" "no-infs-fp-math"="false" "no-nans-fp-math"="false" "no-signed-zeros-fp-math"="false" "no-trapping-math"="false" "stack-protector-buffer-size"="8" "unsafe-fp-math"="false" "use-soft-float"="false" }
attributes #3 = { nounwind }

!llvm.module.flags = !{!0}
!llvm.ident = !{!1}

!0 = !{i32 1, !"wchar_size", i32 4}
!1 = !{!"clang version 9.0.0 (tags/RELEASE_900/final)"}
!2 = !{!3, !3, i64 0}
!3 = !{!"int", !4, i64 0}
!4 = !{!"omnipotent char", !5, i64 0}
!5 = !{!"Simple C/C++ TBAA"}
!6 = !{!4, !4, i64 0}
