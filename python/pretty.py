import io
import os
import sys
import typing


class Progress:
    # Progress represents a progress bar in the terminal.

    def __init__(self) -> None:
        self.chardev = os.isatty(sys.stdout.fileno())
        self.current = 0.0

    def print(self, percent: float) -> None:
        # Update the progress bar to the specified percent (0 to 1).
        if percent > 1:
            raise ValueError('pretty: the percent cannot be greater than 1')
        if percent < self.current:
            raise ValueError('pretty: the percent cannot be decreased')
        if percent not in (0, 1) and percent - self.current < 0.01:
            return
        if percent == 1 and percent == self.current:
            return
        if percent == 0 and self.chardev:
            sys.stderr.write('\x1b7')
            sys.stderr.flush()
        if percent != 0 and self.chardev:
            sys.stderr.write('\x1b8')
            sys.stderr.flush()

        self.current = percent
        cap = int(percent * 44)
        buf = list('[                                             ] 000%')
        for index in range(1, cap + 1):
            buf[index] = '='
        buf[1 + cap] = '>'
        number = f'{int(percent * 100):3d}'
        buf[48] = number[0]
        buf[49] = number[1]
        buf[50] = number[2]
        print('pretty:', ''.join(buf))


class ProgressWriter(io.Writer):
    # ProgressWriter updates a progress bar as data is written.

    def __init__(self, w: io.Writer, n: int) -> None:
        self.progress = Progress()
        self.progress.print(0)
        self.w = w
        self.m = 0
        self.n = n

    def write(self, data: typing.Any) -> int:
        a = self.w.write(data)
        self.m += a
        self.progress.print(self.m / self.n)
        return a


class Table:
    # Table represents a table structure with a head and body.

    def __init__(self) -> None:
        self.conf: list[str] = []
        self.head: list[str] = []
        self.body: list[list[str]] = []

    def print(self) -> None:
        conf = list(self.conf)
        while len(conf) < len(self.head):
            conf.append('<')

        size = [len(c) for c in self.head]
        for r in self.body:
            for i, c in enumerate(r):
                size[i] = max(size[i], len(c))

        line = [''] * len(self.head)
        for i, c in enumerate(self.head):
            w = size[i]
            if conf[i] == '>':
                line[i] = c.rjust(w)
            else:
                line[i] = c.ljust(w)
        print('pretty:', ' '.join(line))

        for i, w in enumerate(size):
            line[i] = '-' * w
        print('pretty:', '-'.join(line))

        for r in self.body:
            for i, c in enumerate(r):
                w = size[i]
                if conf[i] == '>':
                    line[i] = c.rjust(w)
                else:
                    line[i] = c.ljust(w)
            print('pretty:', ' '.join(line))


class Tree:
    # Tree represents a node in a tree structure.

    def __init__(self, name: str) -> None:
        self.name = name
        self.leaf: list[Tree] = []

    def inner(self, prefix: str) -> None:
        for index, elem in enumerate(self.leaf):
            is_last = index == len(self.leaf) - 1
            branch = '└── ' if is_last else '├── '
            print('pretty:', prefix + branch + elem.name)
            if elem.leaf:
                middle = '    ' if is_last else '│   '
                elem.inner(prefix + middle)

    def print(self) -> None:
        print('pretty:', self.name)
        self.inner('')
