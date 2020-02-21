; ModuleID = '/home/camil.bancioiu/Work/Elrond/arwen-wasm-vm/test/contracts/num-with-fp/num-with-fp.c'
source_filename = "/home/camil.bancioiu/Work/Elrond/arwen-wasm-vm/test/contracts/num-with-fp/num-with-fp.c"
target datalayout = "e-m:e-p:32:32-i64:64-n32:64-S128"
target triple = "wasm32-unknown-unknown-wasm"

; Function Attrs: noinline nounwind optnone
define void @doSomething() #0 {
entry:
  %x = alloca i64, align 8
  %a = alloca float, align 4
  %q = alloca float, align 4
  %s = alloca i64, align 8
  store i64 6, i64* %x, align 8
  store float 1.000000e+00, float* %a, align 4
  %0 = load float, float* %a, align 4
  %add = fadd float %0, 0x3FD3333340000000
  store float %add, float* %a, align 4
  %1 = load i64, i64* %x, align 8
  %conv = uitofp i64 %1 to float
  %2 = load float, float* %a, align 4
  %mul = fmul float %conv, %2
  store float %mul, float* %q, align 4
  %3 = bitcast float* %q to i64*
  %4 = load i64, i64* %3, align 4
  store i64 %4, i64* %s, align 8
  %5 = load i64, i64* %s, align 8
  call void @int64finish(i64 %5)
  ret void
}

declare void @int64finish(i64) #1

; Function Attrs: noinline nounwind optnone
define void @init() #0 {
entry:
  ret void
}

; Function Attrs: noinline nounwind optnone
define void @_main() #0 {
entry:
  ret void
}

attributes #0 = { noinline nounwind optnone "correctly-rounded-divide-sqrt-fp-math"="false" "disable-tail-calls"="false" "less-precise-fpmad"="false" "min-legal-vector-width"="0" "no-frame-pointer-elim"="false" "no-infs-fp-math"="false" "no-jump-tables"="false" "no-nans-fp-math"="false" "no-signed-zeros-fp-math"="false" "no-trapping-math"="false" "stack-protector-buffer-size"="8" "unsafe-fp-math"="false" "use-soft-float"="false" }
attributes #1 = { "correctly-rounded-divide-sqrt-fp-math"="false" "disable-tail-calls"="false" "less-precise-fpmad"="false" "no-frame-pointer-elim"="false" "no-infs-fp-math"="false" "no-nans-fp-math"="false" "no-signed-zeros-fp-math"="false" "no-trapping-math"="false" "stack-protector-buffer-size"="8" "unsafe-fp-math"="false" "use-soft-float"="false" }

!llvm.module.flags = !{!0}
!llvm.ident = !{!1}

!0 = !{i32 1, !"wchar_size", i32 4}
!1 = !{!"clang version 9.0.0 (tags/RELEASE_900/final)"}
