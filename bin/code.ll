declare void @print(i8* %s)

declare i8* @append(i8* %s1, i8* %s2)

declare i32 @strlen(i8* %s)

declare i8* @IntToString(i32 %num)

define void @main() {
entry:
	%i = alloca i32
	%0 = mul i32 5, 3
	%1 = add i32 %0, 5
	%2 = add i32 5, %1
	store i32 %2, i32* %i
	%3 = load i32, i32* %i
	%4 = icmp sgt i32 %3, 0
	br i1 %4, label %loop1, label %after1

after1:
	%b = alloca i1
	store i1 true, i1* %b
	%s = alloca i8*
	%5 = alloca [5 x i8]
	store [5 x i8] c"Ahoj\00", [5 x i8]* %5
	%6 = getelementptr [5 x i8], [5 x i8]* %5, i32 0, i32 0
	store i8* %6, i8** %s
	%s2 = alloca i8*
	%7 = alloca [7 x i8]
	store [7 x i8] c" Svete\00", [7 x i8]* %7
	%8 = getelementptr [7 x i8], [7 x i8]* %7, i32 0, i32 0
	store i8* %8, i8** %s2
	%s3 = alloca i8*
	%9 = load i8*, i8** %s
	%10 = load i8*, i8** %s2
	%11 = call i8* @append(i8* %9, i8* %10)
	store i8* %11, i8** %s3
	%12 = load i8*, i8** %s3
	call void @print(i8* %12)
	%13 = add i32 5, 5
	%14 = icmp slt i32 %13, 11
	br i1 %14, label %true2, label %false2

loop1:
	%15 = load i32, i32* %i
	%16 = icmp sgt i32 %15, 0
	br i1 %16, label %loop_body1, label %after1

loop_body1:
	%17 = load i32, i32* %i
	%18 = call i8* @IntToString(i32 %17)
	call void @print(i8* %18)
	%19 = load i32, i32* %i
	%20 = sub i32 %19, 1
	store i32 %20, i32* %i
	%21 = load i32, i32* %i
	%22 = icmp sgt i32 %21, 0
	br i1 %22, label %loop1, label %after1

true2:
	%23 = alloca [4 x i8]
	store [4 x i8] c"Ano\00", [4 x i8]* %23
	%24 = getelementptr [4 x i8], [4 x i8]* %23, i32 0, i32 0
	call void @print(i8* %24)
	br label %after2

false2:
	%25 = alloca [3 x i8]
	store [3 x i8] c"Ne\00", [3 x i8]* %25
	%26 = getelementptr [3 x i8], [3 x i8]* %25, i32 0, i32 0
	call void @print(i8* %26)
	br label %after2

after2:
	%27 = icmp eq i32 5, 5
	br i1 %27, label %true3, label %false3

true3:
	%28 = alloca [5 x i8]
	store [5 x i8] c"Ahoj\00", [5 x i8]* %28
	%29 = getelementptr [5 x i8], [5 x i8]* %28, i32 0, i32 0
	call void @print(i8* %29)
	br label %after3

false3:
	br label %after3

after3:
	%30 = load i1, i1* %b
	%31 = xor i1 %30, true
	br i1 %31, label %true4, label %false4

true4:
	%32 = alloca [6 x i8]
	store [6 x i8] c"false\00", [6 x i8]* %32
	%33 = getelementptr [6 x i8], [6 x i8]* %32, i32 0, i32 0
	call void @print(i8* %33)
	br label %after4

false4:
	%34 = alloca [5 x i8]
	store [5 x i8] c"true\00", [5 x i8]* %34
	%35 = getelementptr [5 x i8], [5 x i8]* %34, i32 0, i32 0
	call void @print(i8* %35)
	br label %after4

after4:
	ret void
}
