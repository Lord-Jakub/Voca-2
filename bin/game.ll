@str0 = global [10 x i8] c"Voca game\00"
@str1 = global [11 x i8] c"player.png\00"
@str2 = global [2 x i8] c"W\00"
@str3 = global [2 x i8] c"S\00"
@str4 = global [2 x i8] c"D\00"
@str5 = global [2 x i8] c"A\00"
@str6 = global [11 x i8] c"player.png\00"
@str7 = global [2 x i8] c"W\00"
@str8 = global [2 x i8] c"S\00"
@str9 = global [2 x i8] c"D\00"
@str10 = global [2 x i8] c"A\00"
@str11 = global [12 x i8] c"Closing app\00"

declare void @print(i8* %s)

declare i8* @append(i8* %s1, i8* %s2)

declare i32 @strlen(i8* %s)

declare i8* @IntToString(i32 %num)

declare i8* @Read()

declare float @StringToFloat(i8* %s)

declare i32 @StringToInt(i8* %s)

declare i8* @FloatToString(float %f)

declare i32 @FloatToInt(float %f)

declare void @delay(i32 %ms)

declare i32 @gInit()

declare i32 @gCreateWindow(i8* %title, i32 %w, i32 %h)

declare void @gSetColor(i32 %r, i32 %g, i32 %b, i32 %a)

declare void @gClear()

declare void @gClose()

declare void @gQuit()

declare void @gCreateFillRect(i32 %x, i32 %y, i32 %w, i32 %h)

declare void @gUpdate()

declare void @gCreateRect(i32 %x, i32 %y, i32 %w, i32 %h)

declare void @gCreateLine(i32 %x1, i32 %y1, i32 %x2, i32 %y2)

declare void @gCreateCircle(i32 %x, i32 %y, i32 %r)

declare void @gCreatePoint(i32 %x, i32 %y)

declare void @gCreatePicture(i8* %path, i32 %x, i32 %y, i32 %w, i32 %h)

declare i1 @gKeyPressed(i8* %keyName)

declare i1 @gMouseDown(i8* %button)

declare i32 @gMouseX()

declare i32 @gMouseY()

declare i1 @gIsRunning()

define i32 @graphics.Init() {
entry:
	%0 = call i32 @gInit()
	ret i32 %0
}

define void @graphics.CreateWindow(i8* %title, i32 %w, i32 %h) {
entry:
	%0 = alloca i8*
	store i8* %title, i8** %0
	%1 = alloca i32
	store i32 %w, i32* %1
	%2 = alloca i32
	store i32 %h, i32* %2
	%3 = load i8*, i8** %0
	%4 = load i32, i32* %1
	%5 = load i32, i32* %2
	%6 = call i32 @gCreateWindow(i8* %3, i32 %4, i32 %5)
	ret void
}

define void @graphics.SetColor(i32 %r, i32 %g, i32 %b, i32 %a) {
entry:
	%0 = alloca i32
	store i32 %r, i32* %0
	%1 = alloca i32
	store i32 %g, i32* %1
	%2 = alloca i32
	store i32 %b, i32* %2
	%3 = alloca i32
	store i32 %a, i32* %3
	%4 = load i32, i32* %0
	%5 = load i32, i32* %1
	%6 = load i32, i32* %2
	%7 = load i32, i32* %3
	call void @gSetColor(i32 %4, i32 %5, i32 %6, i32 %7)
	ret void
}

define void @graphics.Clear() {
entry:
	call void @gClear()
	ret void
}

define void @graphics.Close() {
entry:
	call void @gClose()
	ret void
}

define void @graphics.Quit() {
entry:
	call void @gQuit()
	ret void
}

define void @graphics.CreateFillRect(i32 %x, i32 %y, i32 %w, i32 %h) {
entry:
	%0 = alloca i32
	store i32 %x, i32* %0
	%1 = alloca i32
	store i32 %y, i32* %1
	%2 = alloca i32
	store i32 %w, i32* %2
	%3 = alloca i32
	store i32 %h, i32* %3
	%4 = load i32, i32* %0
	%5 = load i32, i32* %1
	%6 = load i32, i32* %2
	%7 = load i32, i32* %3
	call void @gCreateFillRect(i32 %4, i32 %5, i32 %6, i32 %7)
	ret void
}

define void @graphics.Update() {
entry:
	call void @gUpdate()
	ret void
}

define void @graphics.CreateRect(i32 %x, i32 %y, i32 %w, i32 %h) {
entry:
	%0 = alloca i32
	store i32 %x, i32* %0
	%1 = alloca i32
	store i32 %y, i32* %1
	%2 = alloca i32
	store i32 %w, i32* %2
	%3 = alloca i32
	store i32 %h, i32* %3
	%4 = load i32, i32* %0
	%5 = load i32, i32* %1
	%6 = load i32, i32* %2
	%7 = load i32, i32* %3
	call void @gCreateRect(i32 %4, i32 %5, i32 %6, i32 %7)
	ret void
}

define void @graphics.CreateLine(i32 %x1, i32 %y1, i32 %x2, i32 %y2) {
entry:
	%0 = alloca i32
	store i32 %x1, i32* %0
	%1 = alloca i32
	store i32 %y1, i32* %1
	%2 = alloca i32
	store i32 %x2, i32* %2
	%3 = alloca i32
	store i32 %y2, i32* %3
	%4 = load i32, i32* %0
	%5 = load i32, i32* %1
	%6 = load i32, i32* %2
	%7 = load i32, i32* %3
	call void @gCreateLine(i32 %4, i32 %5, i32 %6, i32 %7)
	ret void
}

define void @graphics.CreateCircle(i32 %x, i32 %y, i32 %r) {
entry:
	%0 = alloca i32
	store i32 %x, i32* %0
	%1 = alloca i32
	store i32 %y, i32* %1
	%2 = alloca i32
	store i32 %r, i32* %2
	%3 = load i32, i32* %0
	%4 = load i32, i32* %1
	%5 = load i32, i32* %2
	call void @gCreateCircle(i32 %3, i32 %4, i32 %5)
	ret void
}

define void @graphics.CreatePoint(i32 %x, i32 %y) {
entry:
	%0 = alloca i32
	store i32 %x, i32* %0
	%1 = alloca i32
	store i32 %y, i32* %1
	%2 = load i32, i32* %0
	%3 = load i32, i32* %1
	call void @gCreatePoint(i32 %2, i32 %3)
	ret void
}

define void @graphics.CreatePicture(i8* %path, i32 %x, i32 %y, i32 %w, i32 %h) {
entry:
	%0 = alloca i8*
	store i8* %path, i8** %0
	%1 = alloca i32
	store i32 %x, i32* %1
	%2 = alloca i32
	store i32 %y, i32* %2
	%3 = alloca i32
	store i32 %w, i32* %3
	%4 = alloca i32
	store i32 %h, i32* %4
	%5 = load i8*, i8** %0
	%6 = load i32, i32* %1
	%7 = load i32, i32* %2
	%8 = load i32, i32* %3
	%9 = load i32, i32* %4
	call void @gCreatePicture(i8* %5, i32 %6, i32 %7, i32 %8, i32 %9)
	ret void
}

define i1 @graphics.KeyPressed(i8* %keyName) {
entry:
	%0 = alloca i8*
	store i8* %keyName, i8** %0
	%1 = load i8*, i8** %0
	%2 = call i1 @gKeyPressed(i8* %1)
	ret i1 %2
}

define i1 @graphics.MouseDown(i8* %button) {
entry:
	%0 = alloca i8*
	store i8* %button, i8** %0
	%1 = load i8*, i8** %0
	%2 = call i1 @gMouseDown(i8* %1)
	ret i1 %2
}

define i32 @graphics.MouseX() {
entry:
	%0 = call i32 @gMouseX()
	ret i32 %0
}

define i32 @graphics.MouseY() {
entry:
	%0 = call i32 @gMouseY()
	ret i32 %0
}

define i1 @graphics.IsRunning() {
entry:
	%0 = call i1 @gIsRunning()
	ret i1 %0
}

define void @main(i32 %argc, i8** %argv) {
entry:
	%init = alloca i32
	%0 = call i32 @graphics.Init()
	store i32 %0, i32* %init
	store [10 x i8] c"Voca game\00", [10 x i8]* @str0
	%1 = getelementptr [10 x i8], [10 x i8]* @str0, i32 0, i32 0
	call void @graphics.CreateWindow(i8* %1, i32 800, i32 600)
	%x = alloca i32
	store i32 10, i32* %x
	%y = alloca i32
	store i32 10, i32* %y
	%2 = call i1 @graphics.IsRunning()
	br i1 %2, label %loop1, label %after1

after1:
	store [12 x i8] c"Closing app\00", [12 x i8]* @str11
	%3 = getelementptr [12 x i8], [12 x i8]* @str11, i32 0, i32 0
	call void @print(i8* %3)
	call void @graphics.Close()
	ret void

loop1:
	%4 = call i1 @graphics.IsRunning()
	br i1 %4, label %loop_body1, label %after1

loop_body1:
	call void @graphics.SetColor(i32 255, i32 255, i32 255, i32 255)
	call void @graphics.Clear()
	call void @graphics.SetColor(i32 255, i32 0, i32 0, i32 255)
	store [11 x i8] c"player.png\00", [11 x i8]* @str1
	%5 = getelementptr [11 x i8], [11 x i8]* @str1, i32 0, i32 0
	%6 = load i32, i32* %x
	%7 = load i32, i32* %y
	call void @graphics.CreatePicture(i8* %5, i32 %6, i32 %7, i32 64, i32 64)
	store [2 x i8] c"W\00", [2 x i8]* @str2
	%8 = getelementptr [2 x i8], [2 x i8]* @str2, i32 0, i32 0
	%9 = call i1 @graphics.KeyPressed(i8* %8)
	br i1 %9, label %true2, label %false2

true2:
	%10 = load i32, i32* %y
	%11 = sub i32 %10, 2
	store i32 %11, i32* %y
	br label %after2

false2:
	br label %after2

after2:
	store [2 x i8] c"S\00", [2 x i8]* @str3
	%12 = getelementptr [2 x i8], [2 x i8]* @str3, i32 0, i32 0
	%13 = call i1 @graphics.KeyPressed(i8* %12)
	br i1 %13, label %true3, label %false3

true3:
	%14 = load i32, i32* %y
	%15 = add i32 %14, 2
	store i32 %15, i32* %y
	br label %after3

false3:
	br label %after3

after3:
	store [2 x i8] c"D\00", [2 x i8]* @str4
	%16 = getelementptr [2 x i8], [2 x i8]* @str4, i32 0, i32 0
	%17 = call i1 @graphics.KeyPressed(i8* %16)
	br i1 %17, label %true4, label %false4

true4:
	%18 = load i32, i32* %x
	%19 = add i32 %18, 2
	store i32 %19, i32* %x
	br label %after4

false4:
	br label %after4

after4:
	store [2 x i8] c"A\00", [2 x i8]* @str5
	%20 = getelementptr [2 x i8], [2 x i8]* @str5, i32 0, i32 0
	%21 = call i1 @graphics.KeyPressed(i8* %20)
	br i1 %21, label %true5, label %false5

true5:
	%22 = load i32, i32* %x
	%23 = sub i32 %22, 2
	store i32 %23, i32* %x
	br label %after5

false5:
	br label %after5

after5:
	call void @graphics.Update()
	%24 = sdiv i32 1000, 60
	call void @delay(i32 %24)
	call void @graphics.SetColor(i32 255, i32 255, i32 255, i32 255)
	call void @graphics.Clear()
	call void @graphics.SetColor(i32 255, i32 0, i32 0, i32 255)
	store [11 x i8] c"player.png\00", [11 x i8]* @str6
	%25 = getelementptr [11 x i8], [11 x i8]* @str6, i32 0, i32 0
	%26 = load i32, i32* %x
	%27 = load i32, i32* %y
	call void @graphics.CreatePicture(i8* %25, i32 %26, i32 %27, i32 64, i32 64)
	store [2 x i8] c"W\00", [2 x i8]* @str7
	%28 = getelementptr [2 x i8], [2 x i8]* @str7, i32 0, i32 0
	%29 = call i1 @graphics.KeyPressed(i8* %28)
	br i1 %29, label %true6, label %false6

true6:
	%30 = load i32, i32* %y
	%31 = sub i32 %30, 2
	store i32 %31, i32* %y
	br label %after6

false6:
	br label %after6

after6:
	store [2 x i8] c"S\00", [2 x i8]* @str8
	%32 = getelementptr [2 x i8], [2 x i8]* @str8, i32 0, i32 0
	%33 = call i1 @graphics.KeyPressed(i8* %32)
	br i1 %33, label %true7, label %false7

true7:
	%34 = load i32, i32* %y
	%35 = add i32 %34, 2
	store i32 %35, i32* %y
	br label %after7

false7:
	br label %after7

after7:
	store [2 x i8] c"D\00", [2 x i8]* @str9
	%36 = getelementptr [2 x i8], [2 x i8]* @str9, i32 0, i32 0
	%37 = call i1 @graphics.KeyPressed(i8* %36)
	br i1 %37, label %true8, label %false8

true8:
	%38 = load i32, i32* %x
	%39 = add i32 %38, 2
	store i32 %39, i32* %x
	br label %after8

false8:
	br label %after8

after8:
	store [2 x i8] c"A\00", [2 x i8]* @str10
	%40 = getelementptr [2 x i8], [2 x i8]* @str10, i32 0, i32 0
	%41 = call i1 @graphics.KeyPressed(i8* %40)
	br i1 %41, label %true9, label %false9

true9:
	%42 = load i32, i32* %x
	%43 = sub i32 %42, 2
	store i32 %43, i32* %x
	br label %after9

false9:
	br label %after9

after9:
	call void @graphics.Update()
	%44 = sdiv i32 1000, 60
	call void @delay(i32 %44)
	%45 = call i1 @graphics.IsRunning()
	br i1 %45, label %loop1, label %after1
}
