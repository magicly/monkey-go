let fac = n => {
  if (n < 1) {
    return 1;
  }
  return n * fac(n - 1);
};

console.log(fac(10));

function fib(n) {
  if (n == 1) {
    return 0;
  }
  if (n == 2) {
    return 1;
  }
  return fib(n - 1) + fib(n - 2);
}

console.log(fib(35));
