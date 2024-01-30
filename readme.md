 Second version of my language, this time static typed and compiled.
**Installation:**

 - So far the amd64  version works fine and
   arm64 version I haven't tested yet.
   
 - For the linux version to work you need to have llc and clang installed both can be installed with `apt-get`.
 - For windows you ned clang in %PATH% (i use mingw64 on msys64)


**Using:**

 - `-i [file]` - specify the input file.
 - `-o [file]` - specifies the output file.
 - `-ast` - generates a json representing the abstract syntax tree of the program.
 - `-ir` will preserve the llvm ir file of the program.
 - `-obj` will preserve the program object file.
 - `-help` to show help.

**Code:**

 - Type functions after the function name: func Foo()string{} 
 - You can use if and while for controling flow.

**Hello world:**
```go
func main(){
print("Hello world")
}