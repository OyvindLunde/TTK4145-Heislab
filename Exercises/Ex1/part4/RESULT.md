The variable i gives different values each run. This is because you can split the increment/decrement line
into two lines: read and write. E.g. the value i can be read and saved to memory, then the program can
switch to the other thread, where it reads the value of i again and increment/decrements it. However, when it
returns to the previous thread the value in the memory is different, meaning when it increments/decrements 
the value will be different from what it should be.