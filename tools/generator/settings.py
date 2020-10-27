import argparse
from pathlib import Path


class Configuration:
    """
    Main configuration.
    Configuration is based on program, arguments: how many users should be generated and path where their will persist
    in files.
    """
    def __init__(self):
        self.fake_count = None
        self.snapshots_path = False

    def parse(self) -> None:
        """
        Parsing program args.
        :return:
        """
        parser = argparse.ArgumentParser()
        parser.add_argument("-c", "--count", required=True, help="count generated dummy users", type=int)
        parser.add_argument("-path", "--snapshot_path", required=True, help="path to the snapshots directory", type=str)

        args = vars(parser.parse_args())

        self.fake_count = args['count']
        self.snapshots_path = Path(args['snapshot_path'])
