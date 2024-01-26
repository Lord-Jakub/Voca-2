declare void @print(i8* %s)

declare i8* @append(i8* %s1, i8* %s2)

define void @main() {
entry:
	%i = alloca i32
	%0 = mul i32 5, 3
	%1 = add i32 %0, 5
	%2 = add i32 5, %1
	store i32 %2, i32* %i
	%s = alloca i8*
	%3 = alloca [5 x i8]
	store [5 x i8] c"Ahoj\00", [5 x i8]* %3
	%4 = getelementptr [5 x i8], [5 x i8]* %3, i32 0, i32 0
	store i8* %4, i8** %s
	%s2 = alloca i8*
	%5 = alloca [10 x i8]
	store [10 x i8] c" Sv\C3\84\C2\9Bte\00", [10 x i8]* %5
	%6 = getelementptr [10 x i8], [10 x i8]* %5, i32 0, i32 0
	store i8* %6, i8** %s2
	%s3 = alloca i8*
	%7 = load i8*, i8** %s
	%8 = load i8*, i8** %s2
	%9 = call i8* @append(i8* %7, i8* %8)
	store i8* %9, i8** %s3
	%10 = load i8*, i8** %s3
	call void @print(i8* %10)
	%11 = add i32 5, 5
	%12 = icmp slt i32 %11, 11
	br i1 %12, label %true1, label %false1

true1:
	%13 = alloca [4 x i8]
	store [4 x i8] c"Ano\00", [4 x i8]* %13
	%14 = getelementptr [4 x i8], [4 x i8]* %13, i32 0, i32 0
	call void @print(i8* %14)
	br label %after1

false1:
	%15 = alloca [3 x i8]
	store [3 x i8] c"Ne\00", [3 x i8]* %15
	%16 = getelementptr [3 x i8], [3 x i8]* %15, i32 0, i32 0
	call void @print(i8* %16)
	br label %after1

after1:
	%17 = icmp eq i32 5, 5
	br i1 %17, label %true2, label %false2

true2:
	%18 = alloca [5 x i8]
	store [5 x i8] c"Ahoj\00", [5 x i8]* %18
	%19 = getelementptr [5 x i8], [5 x i8]* %18, i32 0, i32 0
	call void @print(i8* %19)
	br label %after2

false2:
	br label %after2

after2:
	ret void
}
