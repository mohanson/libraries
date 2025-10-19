import collections
import contextlib
import json
import os.path
import typing
import threading


class MemDriver:
    # MemDriver cares to store data on memory, this means that MemDriver is fast. Since there is no expiration
    # mechanism, be careful that it might eats up all your memory.

    def __init__(self) -> None:
        self.data = {}

    def __contains__(self, key: str) -> bool:
        return key in self.data

    def __delitem__(self, key: str) -> None:
        del self.data[key]

    def __getitem__(self, key: str) -> bytearray:
        return self.data[key]

    def __setitem__(self, key: str, value: bytearray) -> None:
        self.data[key] = value

    def get(self, key: str, default=None) -> bytearray:
        return self.data.get(key, default)

    def pop(self, key: str, default=None) -> bytearray:
        return self.data.pop(key, default)


class DocDriver:
    # DocDriver use the OS's file system to manage data. In general, any high frequency operation is not recommended
    # unless you have an enough reason.

    def __init__(self, root: str) -> None:
        self.root = root
        if not os.path.exists(self.root):
            os.makedirs(self.root)

    def __contains__(self, key: str) -> bool:
        return os.path.exists(os.path.join(self.root, key))

    def __delitem__(self, key: str) -> None:
        os.remove(os.path.join(self.root, key))

    def __getitem__(self, key: str) -> bytearray:
        with open(os.path.join(self.root, key), 'rb') as f:
            return bytearray(f.read())

    def __setitem__(self, key: str, value: bytearray) -> None:
        with open(os.path.join(self.root, key), 'wb') as f:
            f.write(value)

    def get(self, key: str, default=None) -> bytearray:
        with contextlib.suppress(Exception):
            return self[key]
        return default

    def pop(self, key: str, default=None) -> bytearray:
        with contextlib.suppress(Exception):
            value = self[key]
            del self[key]
            return value
        return default


class LruDriver:
    # LruDriver implemention. In computing, cache algorithms (also frequently called cache replacement algorithms or
    # cache replacement policies) are optimizing instructions, or algorithms, that a computer program or a
    # hardware-maintained structure can utilize in order to manage a cache of information stored on the computer.
    # Caching improves performance by keeping recent or often-used data items in a memory locations that are faster or
    # computationally cheaper to access than normal memory stores. When the cache is full, the algorithm must choose
    # which items to discard to make room for the new ones.

    def __init__(self, size: int) -> None:
        self.data = collections.OrderedDict()
        self.size = size

    def __contains__(self, key: str) -> bool:
        return key in self.data

    def __delitem__(self, key: str) -> None:
        del self.data[key]

    def __getitem__(self, key: str) -> bytearray:
        self.data.move_to_end(key)
        return self.data[key]

    def __setitem__(self, key: str, value: bytearray) -> None:
        if len(self.data) >= self.size:
            self.data.popitem(False)
        self.data[key] = value

    def get(self, key: str, default=None) -> bytearray:
        with contextlib.suppress(Exception):
            return self[key]
        return default

    def pop(self, key: str, default=None) -> bytearray:
        return self.data.pop(key, default)


class MapDriver:
    # MapDriver is based on DocDriver and use LruDriver to provide caching at its interface layer. The size of
    # LruDriver is always 1024.

    def __init__(self, root: str) -> None:
        self.doc_driver = DocDriver(root)
        self.lru_driver = LruDriver(1024)

    def __contains__(self, key: str) -> bool:
        return key in self.lru_driver or key in self.doc_driver

    def __delitem__(self, key: str) -> None:
        with contextlib.suppress(Exception):
            del self.lru_driver[key]
        del self.doc_driver[key]

    def __getitem__(self, key: str) -> bytearray:
        with contextlib.suppress(KeyError):
            return self.lru_driver[key]
        value = self.doc_driver[key]
        self.lru_driver[key] = value
        return value

    def __setitem__(self, key: str, value: bytearray) -> None:
        self.lru_driver[key] = value
        self.doc_driver[key] = value

    def get(self, key: str, default=None) -> bytearray:
        with contextlib.suppress(KeyError):
            return self[key]
        return default

    def pop(self, key: str, default=None) -> bytearray:
        with contextlib.suppress(Exception):
            value = self[key]
            del self[key]
            return value
        return default


class Emerge:
    # Emerge is a actuator of the given drive. Do not worry, Is's concurrency-safety.

    def __init__(self, driver: MemDriver | DocDriver | LruDriver | MapDriver) -> None:
        self.driver = driver
        self.lock = threading.Lock()

    def __contains__(self, key: str) -> bool:
        with self.lock:
            return key in self.driver

    def __delitem__(self, key: str) -> None:
        with self.lock:
            del self.driver[key]

    def __getitem__(self, key: str) -> typing.Any:
        with self.lock:
            return json.loads(self.driver[key])

    def __setitem__(self, key: str, value: typing.Any) -> None:
        with self.lock:
            self.driver[key] = json.dumps(value).encode()

    def get(self, key: str, default=None) -> typing.Any:
        with self.lock:
            with contextlib.suppress(Exception):
                return self[key]
            return default

    def pop(self, key: str, default=None) -> None:
        with self.lock:
            with contextlib.suppress(Exception):
                value = self[key]
                del self[key]
                return value
            return default
