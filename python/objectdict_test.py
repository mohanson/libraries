import objectdict


def test_objectdict():
    obj = objectdict.ObjectDict()
    obj.int = 1
    obj.str = 'Hello World'
    obj.map = objectdict.ObjectDict({'k': 'v'})
    assert obj.int == 1
    assert obj.str == 'Hello World'
    assert obj.map.k == 'v'
    obj.map.k = 'w'
    assert obj.map.k == 'w'
