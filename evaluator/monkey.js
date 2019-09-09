let fac = fn (n) {
  if (n < 1) {
    return 1;
  }
  return n * fac(n - 1)
}

puts(fac(10));


let fib = fn(n) {
  if (n == 1) {
    return 0;
  }
  if (n == 2) {
    return 1;
  }
  return fib(n - 1) + fib(n - 2);
}

puts(fib(35));