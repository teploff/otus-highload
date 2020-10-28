from driver import MySQLDriver
from file_reader import Reader
from settings import Configuration
from utils import divide_sequence_into_chunks


def app():
    cfg = Configuration()
    cfg.parse()

    reader = Reader(cfg.snapshots_path)
    reader.do()

    repository = MySQLDriver(cfg.storage)
    repository.insert(list(divide_sequence_into_chunks(reader.users, cfg.batch_size)))


if __name__ == '__main__':
    app()
