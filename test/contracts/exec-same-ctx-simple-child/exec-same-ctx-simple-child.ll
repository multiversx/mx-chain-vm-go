; ModuleID = '/var/work/Elrond/arwen-wasm-vm/test/contracts/exec-same-ctx-simple-child/exec-same-ctx-simple-child.c'
source_filename = "/var/work/Elrond/arwen-wasm-vm/test/contracts/exec-same-ctx-simple-child/exec-same-ctx-simple-child.c"
target datalayout = "e-m:e-p:32:32-i64:64-n32:64-S128"
target triple = "wasm32-unknown-unknown-wasm"

@dataLen = local_unnamed_addr constant i32 1000, align 4
@data = local_unnamed_addr global [1000 x i8] zeroinitializer, align 16
@__const.childFunction.msg = private unnamed_addr constant [6 x i8] c"child\00", align 1

; Function Attrs: nounwind
define void @childFunction() local_unnamed_addr #0 {
entry:
  %msg = alloca [6 x i8], align 1
  %0 = getelementptr inbounds [6 x i8], [6 x i8]* %msg, i32 0, i32 0
  call void @llvm.lifetime.start.p0i8(i64 6, i8* nonnull %0) #3
  call void @llvm.memcpy.p0i8.p0i8.i32(i8* nonnull align 1 %0, i8* align 1 getelementptr inbounds ([6 x i8], [6 x i8]* @__const.childFunction.msg, i32 0, i32 0), i32 6, i1 false)
  call void @finish(i8* nonnull %0, i32 5) #3
  br label %for.body

for.body:                                         ; preds = %for.body, %entry
  %i.020 = phi i32 [ 0, %entry ], [ %inc, %for.body ]
  %conv = trunc i32 %i.020 to i8
  %arrayidx = getelementptr inbounds [1000 x i8], [1000 x i8]* @data, i32 0, i32 %i.020
  store i8 %conv, i8* %arrayidx, align 1, !tbaa !2
  %inc = add nuw nsw i32 %i.020, 1
  %exitcond21 = icmp eq i32 %inc, 1000
  br i1 %exitcond21, label %for.body6, label %for.body

for.cond.cleanup5:                                ; preds = %for.body6
  call void @llvm.lifetime.end.p0i8(i64 6, i8* nonnull %0) #3
  ret void

for.body6:                                        ; preds = %for.body, %for.body6
  %i1.019 = phi i32 [ %inc10, %for.body6 ], [ 1, %for.body ]
  %sub = add nsw i32 %i1.019, -1
  %arrayidx7 = getelementptr inbounds [1000 x i8], [1000 x i8]* @data, i32 0, i32 %sub
  %1 = load i8, i8* %arrayidx7, align 1, !tbaa !2
  %conv8 = zext i8 %1 to i64
  call void @int64finish(i64 %conv8) #3
  %inc10 = add nuw nsw i32 %i1.019, 1
  %exitcond = icmp eq i32 %inc10, 1001
  br i1 %exitcond, label %for.cond.cleanup5, label %for.body6
}

; Function Attrs: argmemonly nounwind
declare void @llvm.lifetime.start.p0i8(i64 immarg, i8* nocapture) #1

; Function Attrs: argmemonly nounwind
declare void @llvm.memcpy.p0i8.p0i8.i32(i8* nocapture writeonly, i8* nocapture readonly, i32, i1 immarg) #1

declare void @finish(i8*, i32) local_unnamed_addr #2

; Function Attrs: argmemonly nounwind
declare void @llvm.lifetime.end.p0i8(i64 immarg, i8* nocapture) #1

declare void @int64finish(i64) local_unnamed_addr #2

attributes #0 = { nounwind "correctly-rounded-divide-sqrt-fp-math"="false" "disable-tail-calls"="false" "less-precise-fpmad"="false" "min-legal-vector-width"="0" "no-frame-pointer-elim"="false" "no-infs-fp-math"="false" "no-jump-tables"="false" "no-nans-fp-math"="false" "no-signed-zeros-fp-math"="false" "no-trapping-math"="false" "stack-protector-buffer-size"="8" "unsafe-fp-math"="false" "use-soft-float"="false" }
attributes #1 = { argmemonly nounwind }
attributes #2 = { "correctly-rounded-divide-sqrt-fp-math"="false" "disable-tail-calls"="false" "less-precise-fpmad"="false" "no-frame-pointer-elim"="false" "no-infs-fp-math"="false" "no-nans-fp-math"="false" "no-signed-zeros-fp-math"="false" "no-trapping-math"="false" "stack-protector-buffer-size"="8" "unsafe-fp-math"="false" "use-soft-float"="false" }
attributes #3 = { nounwind }

!llvm.module.flags = !{!0}
!llvm.ident = !{!1}

!0 = !{i32 1, !"wchar_size", i32 4}
!1 = !{!"clang version 9.0.0 (tags/RELEASE_900/final)"}
!2 = !{!3, !3, i64 0}
!3 = !{!"omnipotent char", !4, i64 0}
!4 = !{!"Simple C/C++ TBAA"}
