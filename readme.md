 Second version of my language, this time static typed and compiled.
**Installation:**

 - So far the amd64 linux version works fine (windows doesn't work yet) and
   arm64 version I haven't tested yet.
   
 - For the linux version to work you need to have llc and clang installed both can be installed with `apt-get`.

**Using:**

 - `-i [file]` - specify the input file.
 - `-o [file]` - specifies the output file.
 - `-ast` - generates a json representing the abstract syntax tree of the program.
 - `-ir` will preserve the llvm ir file of the program.
 - `-obj` will preserve the program object file.

**Code:**

 - Type functions after the function name: func Foo()string{} 
 - The only flowcontrol so far is if (and even that is limited)

**Hello world:**
```go
func main(){
print("Hello world")
}