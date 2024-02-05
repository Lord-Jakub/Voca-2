declare void @print(i8* %s)

declare i8* @append(i8* %s1, i8* %s2)

declare i32 @strlen(i8* %s)

declare i8* @IntToString(i32 %num)

declare i8* @Read()

declare float @StringToFloat(i8* %s)

declare i32 @StringToInt(i8* %s)

declare i8* @FloatToString(float %f)

declare i32 @FloatToInt(float %f)

define i32 @math.Add(i32 %x, i32 %y) {
entry:
	%0 = alloca i32
	store i32 %x, i32* %0
	%1 = alloca i32
	store i32 %y, i32* %1
	%2 = load i32, i32* %0
	%3 = load i32, i32* %1
	%4 = add i32 %2, %3
	ret i32 %4
}

define void @main() {
entry:
	%i = alloca i32
	%0 = add i32 5, 5
	%1 = mul i32 %0, 3
	%2 = add i32 %1, 5
	%3 = sub i32 15, 5
	%4 = sub i32 %3, 5
	%5 = sub i32 %4, 5
	%6 = add i32 %2, %5
	store i32 %6, i32* %i
	%f = alloca float
	%7 = fmul float 5.5, 0x4015333320000000
	%8 = fmul float 5.5, 5.5
	%9 = fadd float %7, %8
	store float %9, float* %f
	%10 = load float, float* %f
	%11 = call i8* @FloatToString(float %10)
	call void @print(i8* %11)
	%12 = load i32, i32* %i
	%13 = icmp sgt i32 %12, 0
	br i1 %13, label %loop1, label %after1

after1:
	%b = alloca i1
	store i1 true, i1* %b
	%s = alloca i8*
	%14 = alloca [5 x i8]
	store [5 x i8] c"Ahoj\00", [5 x i8]* %14
	%15 = getelementptr [5 x i8], [5 x i8]* %14, i32 0, i32 0
	store i8* %15, i8** %s
	%s2 = alloca i8*
	%16 = alloca [7 x i8]
	store [7 x i8] c" Svete\00", [7 x i8]* %16
	%17 = getelementptr [7 x i8], [7 x i8]* %16, i32 0, i32 0
	store i8* %17, i8** %s2
	%s3 = alloca i8*
	%18 = load i8*, i8** %s
	%19 = load i8*, i8** %s2
	%20 = call i8* @append(i8* %18, i8* %19)
	store i8* %20, i8** %s3
	%21 = load i8*, i8** %s3
	call void @print(i8* %21)
	%22 = add i32 5, 5
	%23 = icmp slt i32 %22, 11
	br i1 %23, label %true2, label %false2

loop1:
	%24 = load i32, i32* %i
	%25 = icmp sgt i32 %24, 0
	br i1 %25, label %loop_body1, label %after1

loop_body1:
	%26 = load i32, i32* %i
	%27 = call i8* @IntToString(i32 %26)
	call void @print(i8* %27)
	%28 = load i32, i32* %i
	%29 = sub i32 %28, 1
	store i32 %29, i32* %i
	%30 = load i32, i32* %i
	%31 = icmp sgt i32 %30, 0
	br i1 %31, label %loop1, label %after1

true2:
	%32 = alloca [4 x i8]
	store [4 x i8] c"Ano\00", [4 x i8]* %32
	%33 = getelementptr [4 x i8], [4 x i8]* %32, i32 0, i32 0
	call void @print(i8* %33)
	br label %after2

false2:
	%34 = alloca [3 x i8]
	store [3 x i8] c"Ne\00", [3 x i8]* %34
	%35 = getelementptr [3 x i8], [3 x i8]* %34, i32 0, i32 0
	call void @print(i8* %35)
	br label %after2

after2:
	%36 = icmp eq i32 5, 5
	br i1 %36, label %true3, label %false3

true3:
	%37 = alloca [5 x i8]
	store [5 x i8] c"Ahoj\00", [5 x i8]* %37
	%38 = getelementptr [5 x i8], [5 x i8]* %37, i32 0, i32 0
	call void @print(i8* %38)
	br label %after3

false3:
	br label %after3

after3:
	%39 = load i1, i1* %b
	%40 = xor i1 %39, true
	br i1 %40, label %true4, label %false4

true4:
	%41 = alloca [6 x i8]
	store [6 x i8] c"false\00", [6 x i8]* %41
	%42 = getelementptr [6 x i8], [6 x i8]* %41, i32 0, i32 0
	call void @print(i8* %42)
	br label %after4

false4:
	%43 = alloca [5 x i8]
	store [5 x i8] c"true\00", [5 x i8]* %43
	%44 = getelementptr [5 x i8], [5 x i8]* %43, i32 0, i32 0
	call void @print(i8* %44)
	br label %after4

after4:
	%45 = call i32 @math.Add(i32 5, i32 5)
	%46 = call i8* @IntToString(i32 %45)
	call void @print(i8* %46)
	%47 = alloca [13 x i8]
	store [13 x i8] c"Zadej neco: \00", [13 x i8]* %47
	%48 = getelementptr [13 x i8], [13 x i8]* %47, i32 0, i32 0
	call void @print(i8* %48)
	%49 = alloca [12 x i8]
	store [12 x i8] c"Zadal jsi: \00", [12 x i8]* %49
	%50 = getelementptr [12 x i8], [12 x i8]* %49, i32 0, i32 0
	%51 = call i8* @Read()
	%52 = call i8* @append(i8* %50, i8* %51)
	call void @print(i8* %52)
	ret void
}
