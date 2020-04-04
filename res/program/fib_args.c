#include <stdlib.h>

int fib(int n) {
    if (n <= 1) {
        return n;
    }
    return fib(n - 1) + fib(n - 2);
}

int main(int argc, char *argv[]) {
    int n = atoi(argv[1]);
    int s = atoi(argv[2]);
    if (fib(n) == s) {
        return 0;
    }
    return 1;
}
