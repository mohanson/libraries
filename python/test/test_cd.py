import os
import libraries.cd


def test_cd():
    cwd = os.getcwd()
    with libraries.cd.cd('libraries'):
        sub = os.getcwd()
        assert os.path.join(cwd, 'libraries') == sub
    assert os.getcwd() == cwd
