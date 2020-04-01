int fib(int n) {
    if (n <= 1) {
        return n;
    }
    return fib(n - 1) + fib(n - 2);
}

int main() {
    if (fib(10) == 55) {
        return 0;
    }
    return 1;
}
