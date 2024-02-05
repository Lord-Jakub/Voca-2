#include <string.h>
// function declarations
extern unsigned int FloatToInt(float var0);
extern unsigned char* FloatToString(float var0);
extern unsigned char* IntToString(unsigned int var0);
extern unsigned char* Read(void);
extern float StringToFloat(unsigned char* var0);
extern unsigned int StringToInt(unsigned char* var0);
extern unsigned char* append(unsigned char* var0, unsigned char* var1);
extern void gClear(void);
extern void gClose(void);
extern void gCreateCircle(unsigned int var0, unsigned int var1, unsigned int var2);
extern void gCreateFillRect(unsigned int var0, unsigned int var1, unsigned int var2, unsigned int var3);
extern void gCreateLine(unsigned int var0, unsigned int var1, unsigned int var2, unsigned int var3);
extern void gCreatePicture(unsigned char* var0, unsigned int var1, unsigned int var2, unsigned int var3, unsigned int var4);
extern void gCreatePoint(unsigned int var0, unsigned int var1);
extern void gCreateRect(unsigned int var0, unsigned int var1, unsigned int var2, unsigned int var3);
extern unsigned int gCreateWindow(unsigned char* var0, unsigned int var1, unsigned int var2);
extern unsigned int gInit(void);
extern unsigned int gIsRunning(void);
extern unsigned int gKeyPressed(unsigned char* var0);
extern unsigned int gMouseDown(unsigned char* var0);
extern unsigned int gMouseX(void);
extern unsigned int gMouseY(void);
extern void gQuit(void);
extern void gSetColor(unsigned int var0, unsigned int var1, unsigned int var2, unsigned int var3);
extern void gUpdate(void);
void graphics_Clear(void);
void graphics_Close(void);
void graphics_CreateCircle(unsigned int var0, unsigned int var1, unsigned int var2);
void graphics_CreateFillRect(unsigned int var0, unsigned int var1, unsigned int var2, unsigned int var3);
void graphics_CreateLine(unsigned int var0, unsigned int var1, unsigned int var2, unsigned int var3);
void graphics_CreatePicture(unsigned char* var0, unsigned int var1, unsigned int var2, unsigned int var3, unsigned int var4);
void graphics_CreatePoint(unsigned int var0, unsigned int var1);
void graphics_CreateRect(unsigned int var0, unsigned int var1, unsigned int var2, unsigned int var3);
unsigned int graphics_CreateWindow(unsigned char* var0, unsigned int var1, unsigned int var2);
unsigned int graphics_Init(void);
unsigned int graphics_IsRunning(void);
unsigned int graphics_KeyPressed(unsigned char* var0);
unsigned int graphics_MouseDown(unsigned char* var0);
unsigned int graphics_MouseX(void);
unsigned int graphics_MouseY(void);
void graphics_Quit(void);
void graphics_SetColor(unsigned int var0, unsigned int var1, unsigned int var2, unsigned int var3);
void graphics_Update(void);
int main(int argc, char** argv);
extern void print(unsigned char* var0);
//extern unsigned int strlen(unsigned char* var0);

void graphics_Clear(void){
    block0:
    gClear();
    return;
}

void graphics_Close(void){
    block0:
    gClose();
    return;
}

void graphics_CreateCircle(unsigned int var0, unsigned int var1, unsigned int var2){
    unsigned int var3;
    unsigned int var4;
    unsigned int var5;
    block0:
    var3 = var0;
    var4 = var1;
    var5 = var2;
    gCreateCircle(var3, var4, var5);
    return;
}

void graphics_CreateFillRect(unsigned int var0, unsigned int var1, unsigned int var2, unsigned int var3){
    unsigned int var4;
    unsigned int var5;
    unsigned int var6;
    unsigned int var7;
    block0:
    var4 = var0;
    var5 = var1;
    var6 = var2;
    var7 = var3;
    gCreateFillRect(var4, var5, var6, var7);
    return;
}

void graphics_CreateLine(unsigned int var0, unsigned int var1, unsigned int var2, unsigned int var3){
    unsigned int var4;
    unsigned int var5;
    unsigned int var6;
    unsigned int var7;
    block0:
    var4 = var0;
    var5 = var1;
    var6 = var2;
    var7 = var3;
    gCreateLine(var4, var5, var6, var7);
    return;
}

void graphics_CreatePicture(unsigned char* var0, unsigned int var1, unsigned int var2, unsigned int var3, unsigned int var4){
    unsigned char* var5;
    unsigned int var6;
    unsigned int var7;
    unsigned int var8;
    unsigned int var9;
    block0:
    var5 = var0;
    var6 = var1;
    var7 = var2;
    var8 = var3;
    var9 = var4;
    gCreatePicture(var5, var6, var7, var8, var9);
    return;
}

void graphics_CreatePoint(unsigned int var0, unsigned int var1){
    unsigned int var2;
    unsigned int var3;
    block0:
    var2 = var0;
    var3 = var1;
    gCreatePoint(var2, var3);
    return;
}

void graphics_CreateRect(unsigned int var0, unsigned int var1, unsigned int var2, unsigned int var3){
    unsigned int var4;
    unsigned int var5;
    unsigned int var6;
    unsigned int var7;
    block0:
    var4 = var0;
    var5 = var1;
    var6 = var2;
    var7 = var3;
    gCreateRect(var4, var5, var6, var7);
    return;
}

unsigned int graphics_CreateWindow(unsigned char* var0, unsigned int var1, unsigned int var2){
    unsigned char* var3;
    unsigned int var4;
    unsigned int var5;
    block0:
    var3 = var0;
    var4 = var1;
    var5 = var2;
    
    return gCreateWindow(var3, var4, var5);
}

unsigned int graphics_Init(void){
    block0:
    return gInit();
}

unsigned int graphics_IsRunning(void){
    block0:
    return gIsRunning();
}

unsigned int graphics_KeyPressed(unsigned char* var0){
    unsigned char* var1;
    block0:
    var1 = var0;
    return gKeyPressed(var1);
}

unsigned int graphics_MouseDown(unsigned char* var0){
    unsigned char* var1;
    block0:
    var1 = var0;
    return gMouseDown(var1);
}

unsigned int graphics_MouseX(void){
    block0:
    return gMouseX();
}

unsigned int graphics_MouseY(void){
    block0:
    return gMouseY();
}

void graphics_Quit(void){
    block0:
    gQuit();
    return;
}

void graphics_SetColor(unsigned int var0, unsigned int var1, unsigned int var2, unsigned int var3){
    unsigned int var4;
    unsigned int var5;
    unsigned int var6;
    unsigned int var7;
    block0:
    var4 = var0;
    var5 = var1;
    var6 = var2;
    var7 = var3;
    gSetColor(var4, var5, var6, var7);
    return;
}

void graphics_Update(void){
    block0:
    gUpdate();
    return;
}

int main(int argc, char** argv){
    unsigned int var2;
    unsigned char var3[14];
    unsigned char var4[5];
    unsigned int var5;
    block0:
    var2 = graphics_Init();
    if (var2 != 0){
        unsigned char s[23];
        unsigned char temp_var1[] = "Error initializing window";
        memcpy(s, temp_var1, sizeof(temp_var1));
        print(s);
    }
    unsigned char temp_var3[] = {86, 111, 99, 97, 32, 103, 114, 97, 112, 104, 105, 99, 115, 0};
    memcpy(var3, temp_var3, sizeof(temp_var3));

    unsigned int i = graphics_CreateWindow(&(var3[0]), 800, 600);
    if (i != 0){
        unsigned char s[23];
        unsigned char temp_var6[] = {69, 114, 114, 111, 114, 32, 99, 114, 101, 97, 116, 105, 110, 103, 32, 119, 105, 110, 100, 111, 119, 0,};
        memcpy(s, temp_var6, sizeof(temp_var6));
        print(s);
    }
    if (graphics_IsRunning()) {
        goto block2;
    } else {
        goto block1;
    }
    block2:
    if (graphics_IsRunning()) {
        goto block6;
    } else {
        goto block1;
    }
    block1:
    graphics_Close();
    return 0;
    block6:
    graphics_Update();
    if (graphics_IsRunning()) {
        goto block2;
    } else {
        goto block1;
    }
    block3:
    graphics_SetColor(255, 255, 255, 255);
    graphics_Clear();
    graphics_SetColor(255, 0, 0, 255);
    graphics_CreateFillRect(100, 100, 200, 200);
    unsigned char temp_var4[] = {108,101,102,116,0,};
    memcpy(var4, temp_var4, sizeof(temp_var4));
    if (graphics_MouseDown(&(var4[0]))) {
        graphics_SetColor(0, 255, 0, 255);
        var5 = graphics_MouseX();
        graphics_CreatePoint(var5, graphics_MouseY());
        goto block6;
    } else {
        goto block6;
    }
    return 0;
}