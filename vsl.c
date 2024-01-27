#include <stdio.h>
#include <string.h>
#include <stdlib.h>

//For some reason, the compiler is complaining about this function not being defined - so this must be done.
void ___chkstk_ms(void *target_location);
void __chkstk(void *target_location) {
    ___chkstk_ms(target_location);
}

char* append(char* str1, char* str2) {
    char* result = malloc(strlen(str1) + strlen(str2) + 1);  
    strcpy(result, str1);  
    strcat(result, str2);  
    return result;
}


void print(const char text[]) {
    printf("%s\n", text);
}


