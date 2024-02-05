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
	%1 = alloca [10 x i8]
	store [10 x i8] c"Voca game\00", [10 x i8]* %1
	%2 = getelementptr [10 x i8], [10 x i8]* %1, i32 0, i32 0
	call void @graphics.CreateWindow(i8* %2, i32 800, i32 600)
	%x = alloca i32
	store i32 10, i32* %x
	%y = alloca i32
	store i32 10, i32* %y
	%3 = call i1 @graphics.IsRunning()
	br i1 %3, label %loop1, label %after1

after1:
	%4 = alloca [12 x i8]
	store [12 x i8] c"Closing app\00", [12 x i8]* %4
	%5 = getelementptr [12 x i8], [12 x i8]* %4, i32 0, i32 0
	call void @print(i8* %5)
	call void @graphics.Close()
	ret void

loop1:
	%6 = call i1 @graphics.IsRunning()
	br i1 %6, label %loop_body1, label %after1

loop_body1:
	call void @graphics.SetColor(i32 255, i32 255, i32 255, i32 255)
	call void @graphics.Clear()
	call void @graphics.SetColor(i32 255, i32 0, i32 0, i32 255)
	%7 = alloca [11 x i8]
	store [11 x i8] c"player.png\00", [11 x i8]* %7
	%8 = getelementptr [11 x i8], [11 x i8]* %7, i32 0, i32 0
	%9 = load i32, i32* %x
	%10 = load i32, i32* %y
	call void @graphics.CreatePicture(i8* %8, i32 %9, i32 %10, i32 64, i32 64)
	%11 = alloca [2 x i8]
	store [2 x i8] c"W\00", [2 x i8]* %11
	%12 = getelementptr [2 x i8], [2 x i8]* %11, i32 0, i32 0
	%13 = call i1 @graphics.KeyPressed(i8* %12)
	br i1 %13, label %true2, label %false2

true2:
	%14 = load i32, i32* %y
	%15 = sub i32 %14, 2
	store i32 %15, i32* %y
	br label %after2

false2:
	br label %after2

after2:
	%16 = alloca [2 x i8]
	store [2 x i8] c"S\00", [2 x i8]* %16
	%17 = getelementptr [2 x i8], [2 x i8]* %16, i32 0, i32 0
	%18 = call i1 @graphics.KeyPressed(i8* %17)
	br i1 %18, label %true3, label %false3

true3:
	%19 = load i32, i32* %y
	%20 = add i32 %19, 2
	store i32 %20, i32* %y
	br label %after3

false3:
	br label %after3

after3:
	%21 = alloca [2 x i8]
	store [2 x i8] c"D\00", [2 x i8]* %21
	%22 = getelementptr [2 x i8], [2 x i8]* %21, i32 0, i32 0
	%23 = call i1 @graphics.KeyPressed(i8* %22)
	br i1 %23, label %true4, label %false4

true4:
	%24 = load i32, i32* %x
	%25 = add i32 %24, 2
	store i32 %25, i32* %x
	br label %after4

false4:
	br label %after4

after4:
	%26 = alloca [2 x i8]
	store [2 x i8] c"A\00", [2 x i8]* %26
	%27 = getelementptr [2 x i8], [2 x i8]* %26, i32 0, i32 0
	%28 = call i1 @graphics.KeyPressed(i8* %27)
	br i1 %28, label %true5, label %false5

true5:
	%29 = load i32, i32* %x
	%30 = sub i32 %29, 2
	store i32 %30, i32* %x
	br label %after5

false5:
	br label %after5

after5:
	call void @graphics.Update()
	%31 = sdiv i32 1000, 60
	call void @delay(i32 %31)
	call void @graphics.SetColor(i32 255, i32 255, i32 255, i32 255)
	call void @graphics.Clear()
	call void @graphics.SetColor(i32 255, i32 0, i32 0, i32 255)
	%32 = alloca [11 x i8]
	store [11 x i8] c"player.png\00", [11 x i8]* %32
	%33 = getelementptr [11 x i8], [11 x i8]* %32, i32 0, i32 0
	%34 = load i32, i32* %x
	%35 = load i32, i32* %y
	call void @graphics.CreatePicture(i8* %33, i32 %34, i32 %35, i32 64, i32 64)
	%36 = alloca [2 x i8]
	store [2 x i8] c"W\00", [2 x i8]* %36
	%37 = getelementptr [2 x i8], [2 x i8]* %36, i32 0, i32 0
	%38 = call i1 @graphics.KeyPressed(i8* %37)
	br i1 %38, label %true6, label %false6

true6:
	%39 = load i32, i32* %y
	%40 = sub i32 %39, 2
	store i32 %40, i32* %y
	br label %after6

false6:
	br label %after6

after6:
	%41 = alloca [2 x i8]
	store [2 x i8] c"S\00", [2 x i8]* %41
	%42 = getelementptr [2 x i8], [2 x i8]* %41, i32 0, i32 0
	%43 = call i1 @graphics.KeyPressed(i8* %42)
	br i1 %43, label %true7, label %false7

true7:
	%44 = load i32, i32* %y
	%45 = add i32 %44, 2
	store i32 %45, i32* %y
	br label %after7

false7:
	br label %after7

after7:
	%46 = alloca [2 x i8]
	store [2 x i8] c"D\00", [2 x i8]* %46
	%47 = getelementptr [2 x i8], [2 x i8]* %46, i32 0, i32 0
	%48 = call i1 @graphics.KeyPressed(i8* %47)
	br i1 %48, label %true8, label %false8

true8:
	%49 = load i32, i32* %x
	%50 = add i32 %49, 2
	store i32 %50, i32* %x
	br label %after8

false8:
	br label %after8

after8:
	%51 = alloca [2 x i8]
	store [2 x i8] c"A\00", [2 x i8]* %51
	%52 = getelementptr [2 x i8], [2 x i8]* %51, i32 0, i32 0
	%53 = call i1 @graphics.KeyPressed(i8* %52)
	br i1 %53, label %true9, label %false9

true9:
	%54 = load i32, i32* %x
	%55 = sub i32 %54, 2
	store i32 %55, i32* %x
	br label %after9

false9:
	br label %after9

after9:
	call void @graphics.Update()
	%56 = sdiv i32 1000, 60
	call void @delay(i32 %56)
	%57 = call i1 @graphics.IsRunning()
	br i1 %57, label %loop1, label %after1
}
