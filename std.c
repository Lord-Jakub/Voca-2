#include <stdio.h>
#include <string.h>
#include <stdlib.h>


char* append(char* str1, char* str2) {
    char* result = malloc(strlen(str1) + strlen(str2) + 1);  
    strcpy(result, str1);  
    strcat(result, str2);  
    return result;
}


void print(const char text[]) {
    printf("%s\n", text);
}


