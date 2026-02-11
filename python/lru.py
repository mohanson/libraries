import collections
import threading
import typing


class Lru:
    # Lru cache. It is safe for concurrent access.

    def __init__(self, size: int) -> None:
        assert size > 0
        self.data = collections.OrderedDict()
        self.lock = threading.Lock()
        self.size = size

    def __contains__(self, key: str) -> bool:
        with self.lock:
            return key in self.data

    def __delitem__(self, key: str) -> None:
        with self.lock:
            del self.data[key]

    def __getitem__(self, key: str) -> typing.Any:
        with self.lock:
            return self.data[key]

    def __len__(self) -> int:
        with self.lock:
            return len(self.data)

    def __setitem__(self, key: str, value: typing.Any) -> None:
        with self.lock:
            if len(self.data) >= self.size:
                self.data.popitem(False)
            self.data[key] = value

    def get(self, key: typing.Any, default=None) -> typing.Any:
        # Return the value for key if key is in the dictionary, else default.
        with self.lock:
            if key not in self.data:
                return default
            self.data.move_to_end(key)
            return self.data[key]

    def pop(self, key: typing.Any, default=None) -> typing.Any:
        # Remove specified key and return the corresponding value. If the key is not found, return the default if given;
        # otherwise, raise a KeyError.
        with self.lock:
            return self.data.pop(key, default)
