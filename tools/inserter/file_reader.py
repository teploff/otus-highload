from pathlib import Path

from entity import User


class Reader:
    """

    """
    def __init__(self, path: Path):
        """
        :param path: path where all snapshots of users are stored.
        """
        self.path = path
        self._users = []

    def do(self) -> None:
        """

        :return:
        """
        for item in self.path.iterdir():
            if item.is_file() and item.suffix == '.txt':
                self._read_file(item)

    def _read_file(self, file_path: Path):
        """

        :param file_path:
        :return:
        """
        with open(file_path, 'r') as f:
            line = f.readline().rstrip().split('\t')

            # TODO: fix empty string
            while line and line != ['']:
                self._users.append(User(*line))

                line = f.readline().rstrip().split('\t')
