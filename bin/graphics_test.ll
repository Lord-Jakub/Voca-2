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
	%1 = alloca [14 x i8]
	store [14 x i8] c"Voca graphics\00", [14 x i8]* %1
	%2 = getelementptr [14 x i8], [14 x i8]* %1, i32 0, i32 0
	call void @graphics.CreateWindow(i8* %2, i32 800, i32 600)
	%i = alloca i32
	store i32 0, i32* %i
	%3 = call i1 @graphics.IsRunning()
	br i1 %3, label %loop1, label %after1

after1:
	call void @graphics.Close()
	ret void

loop1:
	%4 = call i1 @graphics.IsRunning()
	br i1 %4, label %loop_body1, label %after1

loop_body1:
	call void @graphics.SetColor(i32 255, i32 255, i32 255, i32 255)
	call void @graphics.Clear()
	call void @graphics.SetColor(i32 255, i32 0, i32 0, i32 255)
	%5 = load i32, i32* %i
	call void @graphics.CreateFillRect(i32 %5, i32 100, i32 200, i32 200)
	%6 = alloca [5 x i8]
	store [5 x i8] c"left\00", [5 x i8]* %6
	%7 = getelementptr [5 x i8], [5 x i8]* %6, i32 0, i32 0
	%8 = call i1 @graphics.MouseDown(i8* %7)
	br i1 %8, label %true2, label %false2

true2:
	%9 = load i32, i32* %i
	%10 = add i32 %9, 1
	store i32 %10, i32* %i
	call void @graphics.SetColor(i32 0, i32 255, i32 0, i32 255)
	%11 = call i32 @graphics.MouseX()
	%12 = call i32 @graphics.MouseY()
	call void @graphics.CreatePoint(i32 %11, i32 %12)
	br label %after2

false2:
	br label %after2

after2:
	call void @graphics.Update()
	%13 = sdiv i32 1000, 60
	call void @delay(i32 %13)
	call void @graphics.SetColor(i32 255, i32 255, i32 255, i32 255)
	call void @graphics.Clear()
	call void @graphics.SetColor(i32 255, i32 0, i32 0, i32 255)
	%14 = load i32, i32* %i
	call void @graphics.CreateFillRect(i32 %14, i32 100, i32 200, i32 200)
	%15 = alloca [5 x i8]
	store [5 x i8] c"left\00", [5 x i8]* %15
	%16 = getelementptr [5 x i8], [5 x i8]* %15, i32 0, i32 0
	%17 = call i1 @graphics.MouseDown(i8* %16)
	br i1 %17, label %true3, label %false3

true3:
	%18 = load i32, i32* %i
	%19 = add i32 %18, 1
	store i32 %19, i32* %i
	call void @graphics.SetColor(i32 0, i32 255, i32 0, i32 255)
	%20 = call i32 @graphics.MouseX()
	%21 = call i32 @graphics.MouseY()
	call void @graphics.CreatePoint(i32 %20, i32 %21)
	br label %after3

false3:
	br label %after3

after3:
	call void @graphics.Update()
	%22 = sdiv i32 1000, 60
	call void @delay(i32 %22)
	%23 = call i1 @graphics.IsRunning()
	br i1 %23, label %loop1, label %after1
}
