from dummy import Dummy
from settings import Configuration


def app():
    cfg = Configuration()
    cfg.parse()

    dummy = Dummy(cfg.fake_count, cfg.snapshots_path)
    dummy.generate()
    dummy.make_snapshot()


if __name__ == "__main__":
    app()
