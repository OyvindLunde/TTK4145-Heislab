# Python 3.3.3 and 2.7.6
# python fo.py

from threading import Thread
import threading

# Potentially useful thing:
#   In Python you "import" a global variable, instead of "export"ing it when you declare it
#   (This is probably an effort to make you feel bad about typing the word "global")
i = 0
lock = threading.Lock()

def incrementingFunction():
    global i
    for j in range(0, 1000000):
        lock.acquire()
        i += 1
        lock.release()
    # TODO: increment i 1_000_000 times

def decrementingFunction():
    global i
    for j in range(0, 1000000):
        lock.acquire()
        i -= 1
        lock.release(
    # TODO: decrement i 1_000_000 times


def main():
    global i

    incrementing = Thread(target = incrementingFunction, args = (),)
    incrementing.start()
    decrementing = Thread(target = decrementingFunction, args = (),)
    decrementing.start()
    
    # TODO: Start both threads
    
    incrementing.join()
    decrementing.join()
    
    print("The magic number is %d" % (i))


main()
