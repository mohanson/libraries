import functools
import traceback
import typing


def wrap(func: typing.Callable) -> typing.Callable:
    # Catch-all exception handling.
    @functools.wraps(func)
    def wrapper(*args, **kwargs):
        try:
            return func(*args, **kwargs)
        except Exception:
            traceback.print_exc()
            return None
    return wrapper
