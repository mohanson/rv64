#include <stdint.h>

int main() {
  uint32_t a  = 0;
  a = a + 3;
  a = a * 3;
  a = a << 3;
  a = a >> 3;
  a = a & 0xffff;
  a = a | 0x2;
  a = a / 2;
  if (a != 5) {
    return 1;
  }
  uint64_t b = 7;
  b = b + 3;
  b = b * 17;
  b = b << 3;
  b = b >> 3;
  b = b & 0xffff;
  b = b | 0x5;
  b = b / 2;
  if (b != 87) {
    return 1;
  }
  return 0;
}
