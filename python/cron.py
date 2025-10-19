import time
import typing


def cron(e: int, d: int) -> typing.Generator[None, None, None]:
    # Cron schedules tasks at regular intervals. This function receives two parameters: the elapsed time (duration)
    # between consecutive task executions (period), and the initial delay before the first task execution. This delay
    # is applied after aligning to the next period boundary.
    assert d < e
    while True:
        n = int(time.time())
        s = e + d - n % e
        time.sleep(s)
        yield
