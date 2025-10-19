import libraries.objectdict


def test_objectdict():
    obj = libraries.objectdict.ObjectDict()
    obj.int = 1
    obj.str = 'Hello World'
    obj.map = {'k': 'v'}
    assert obj.int == 1
    assert obj.str == 'Hello World'
    assert obj.map.k == 'v'
    obj.map.k = 'w'
    assert obj.map.k == 'w'
