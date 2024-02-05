#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <time.h>

//For some reason, the compiler is complaining about this function not being defined - so this must be done.
/*void ___chkstk_ms(void *target_location);
void __chkstk(void *target_location) {
    ___chkstk_ms(target_location);
}*/
char* IntToString(int num) {
    // Určíme maximální délku řetězce na převod
    int len = snprintf(NULL, 0, "%d", num);
    
    // Alokujeme paměť pro řetězec
    char* str = (char*)malloc(len + 1);
    
    // Převod integeru na řetězec
    snprintf(str, len + 1, "%d", num);
    
    return str;
}

char* FloatToString(float num) {
    char* str = (char*)malloc(20 * sizeof(char)); // allocate memory for string
    sprintf(str, "%f", num); // convert float to string
    return str;
}

int FloatToInt(float num) {
    return (int)num; // implicit conversion from float to int
}

int StringToInt(char* str) {
    return atoi(str); // convert string to int
}

float StringToFloat(char* str) {
    return atof(str); // convert string to float
}
char* Read() {
    char* input = (char*)malloc(100 * sizeof(char)); // allocate memory for input
    if (input == NULL) {
        exit(1);
    }
    scanf("%s", input); // read input from command line
    return input;
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
void delay(int milliseconds) {
    clock_t end_wait;
    end_wait = clock() + milliseconds * (CLOCKS_PER_SEC / 1000);
    while (clock() < end_wait) {}
}


