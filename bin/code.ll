declare void @print(i8* %s)

declare i8* @append(i8* %s1, i8* %s2)

declare i32 @strlen(i8* %s)

declare i8* @IntToString(i32 %num)

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
	%7 = load i32, i32* %i
	%8 = icmp sgt i32 %7, 0
	br i1 %8, label %loop1, label %after1

after1:
	%b = alloca i1
	store i1 true, i1* %b
	%s = alloca i8*
	%9 = alloca [5 x i8]
	store [5 x i8] c"Ahoj\00", [5 x i8]* %9
	%10 = getelementptr [5 x i8], [5 x i8]* %9, i32 0, i32 0
	store i8* %10, i8** %s
	%s2 = alloca i8*
	%11 = alloca [7 x i8]
	store [7 x i8] c" Svete\00", [7 x i8]* %11
	%12 = getelementptr [7 x i8], [7 x i8]* %11, i32 0, i32 0
	store i8* %12, i8** %s2
	%s3 = alloca i8*
	%13 = load i8*, i8** %s
	%14 = load i8*, i8** %s2
	%15 = call i8* @append(i8* %13, i8* %14)
	store i8* %15, i8** %s3
	%16 = load i8*, i8** %s3
	call void @print(i8* %16)
	%17 = add i32 5, 5
	%18 = icmp slt i32 %17, 11
	br i1 %18, label %true2, label %false2

loop1:
	%19 = load i32, i32* %i
	%20 = icmp sgt i32 %19, 0
	br i1 %20, label %loop_body1, label %after1

loop_body1:
	%21 = load i32, i32* %i
	%22 = call i8* @IntToString(i32 %21)
	call void @print(i8* %22)
	%23 = load i32, i32* %i
	%24 = sub i32 %23, 1
	store i32 %24, i32* %i
	%25 = load i32, i32* %i
	%26 = icmp sgt i32 %25, 0
	br i1 %26, label %loop1, label %after1

true2:
	%27 = alloca [4 x i8]
	store [4 x i8] c"Ano\00", [4 x i8]* %27
	%28 = getelementptr [4 x i8], [4 x i8]* %27, i32 0, i32 0
	call void @print(i8* %28)
	br label %after2

false2:
	%29 = alloca [3 x i8]
	store [3 x i8] c"Ne\00", [3 x i8]* %29
	%30 = getelementptr [3 x i8], [3 x i8]* %29, i32 0, i32 0
	call void @print(i8* %30)
	br label %after2

after2:
	%31 = icmp eq i32 5, 5
	br i1 %31, label %true3, label %false3

true3:
	%32 = alloca [5 x i8]
	store [5 x i8] c"Ahoj\00", [5 x i8]* %32
	%33 = getelementptr [5 x i8], [5 x i8]* %32, i32 0, i32 0
	call void @print(i8* %33)
	br label %after3

false3:
	br label %after3

after3:
	%34 = load i1, i1* %b
	%35 = xor i1 %34, true
	br i1 %35, label %true4, label %false4

true4:
	%36 = alloca [6 x i8]
	store [6 x i8] c"false\00", [6 x i8]* %36
	%37 = getelementptr [6 x i8], [6 x i8]* %36, i32 0, i32 0
	call void @print(i8* %37)
	br label %after4

false4:
	%38 = alloca [5 x i8]
	store [5 x i8] c"true\00", [5 x i8]* %38
	%39 = getelementptr [5 x i8], [5 x i8]* %38, i32 0, i32 0
	call void @print(i8* %39)
	br label %after4

after4:
	%40 = call i32 @math.Add(i32 5, i32 5)
	%41 = call i8* @IntToString(i32 %40)
	call void @print(i8* %41)
	ret void
}
