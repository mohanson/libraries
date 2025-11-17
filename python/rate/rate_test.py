import rate


def test_rate():
    limits = rate.Limits(10, 1)
    for _ in range(10):
        limits.peek(1)
    for _ in range(10):
        limits.wait(1)
