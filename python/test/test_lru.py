import libraries.lru


def test_lru_append():
    c = libraries.lru.Lru(4)
    c[1] = 1
    c[2] = 2
    c[3] = 3
    c[4] = 4
    c[5] = 5
    assert not 1 in c
    assert c.get(5) == 5


def test_lru_change():
    c = libraries.lru.Lru(4)
    c[1] = 1
    c[2] = 2
    c[3] = 3
    c[4] = 4
    c[1] = 5
    assert c.get(1) == 5


def test_lru_delete():
    c = libraries.lru.Lru(4)
    c[1] = 1
    c[2] = 2
    c[3] = 3
    c[4] = 4
    c.pop(2)
    assert not 2 in c
    assert len(c) == 3


def test_lru_size():
    c = libraries.lru.Lru(4)
    assert len(c) == 0
    c[1] = 1
    assert len(c) == 1
    c[2] = 2
    assert len(c) == 2
    c[3] = 3
    assert len(c) == 3
    c[4] = 4
    assert len(c) == 4
    c[5] = 5
    assert len(c) == 4
