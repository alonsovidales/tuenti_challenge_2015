"""
Using python to work with huge numbers, the best languaje for this proposal :)
"""

from collections import defaultdict

NUMBERS_FILE_PATH = "/home/avidales/tuenti_contest/numbers.txt"

def isPrime(n):
    if n == 2 or n == 3:
        return True

    if n % 2 == 0 or n % 3 == 0:
        return False

    i = 5
    w = 2
    while i * i <= n:
        if n % i == 0:
            return False

        i += w
        w = 6 - w

    return True

def getFirstprimes(num_of_primes):
    primes = []
    i = 2
    while len(primes) < num_of_primes:
        while not isPrime(i):
            i += 1

        primes.append(i)
        i += 1

    return primes

def get_primes_by_line(file_name):
    primes = getFirstprimes(25)
    primes_by_num = []

    with open(file_name) as f:
        for line in f:
            n = long(line)

            #print line
            primes_count = defaultdict(int)

            for prime in primes:
                while n % prime == 0:
                    n = n / prime
                    primes_count[prime] += 1

            primes_by_num.append(primes_count)

    return primes_by_num

def get_max_primes(from_to, primes_by_line):
    primes_count = defaultdict(int)
    for i in range(from_to[0], from_to[1]):
        for prime, total in primes_by_line[i].items():
            primes_count[prime] += total

    max_total = max(primes_count.values())
    to_return = []
    for prime, total in primes_count.items():
        if total == max_total:
            to_return.append(prime)

    sorted(to_return)

    return max_total, to_return

"""
main problem execution
"""
primes_by_line = get_primes_by_line(NUMBERS_FILE_PATH)

for problem in xrange(int(raw_input())):
    max_times, primes = get_max_primes(map(int, raw_input().split()), primes_by_line)

    print "%d %s" % (max_times, ' '.join(map(str, primes)))
