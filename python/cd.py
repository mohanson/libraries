import contextlib
import os
import typing


@contextlib.contextmanager
def cd(dst: str) -> typing.Generator[None, typing.Any, None]:
    # Context manager to temporarily change the current working directory.
    cwd = os.getcwd()
    os.chdir(dst)
    yield
    os.chdir(cwd)
