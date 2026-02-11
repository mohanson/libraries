import hashlib
import os
import random
import time

import gool


def once() -> None:
    for _ in range(1024):
        hasher = hashlib.sha256(random.randbytes(256))
        hasher.digest()


def main_loop() -> int:
    stim = time.time()
    etim = stim
    cnts = 0
    for _ in range(1 << 32):
        once()
        cnts += 1
        etim = time.time()
        if etim - stim >= 4.0:
            break
    return int(float(cnts * 1024) / (etim - stim))


def main_gool() -> int:
    stim = time.time()
    etim = stim
    cnts = 0
    grun = gool.cpu()
    for _ in range(1 << 32):
        grun.call(once)
        cnts += 1
        etim = time.time()
        if etim - stim >= 4.0:
            break
    grun.wait()
    etim = time.time()
    return int(float(cnts * 1024) / (etim - stim))


print('main:', os.cpu_count(), 'logical cpus usable by the current process')
print('main: sha256 by loop')
print('main: sha256 by loop rate', main_loop())
print('main: sha256 by gool')
print('main: sha256 by gool rate', main_gool())
