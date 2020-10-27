import argparse
from omegaconf import OmegaConf


class Configuration:
    def __init__(self):
        self.storage = None
        self.fake_count = None
        self.snapshots_path = False

    def parse(self):
        """
        Parsing YAML configuration with environment variables.
        :return:
        """
        parser = argparse.ArgumentParser()
        parser.add_argument("-cfg", "--configuration", required=True, help="path to the config file", type=str)
        parser.add_argument("-c", "--count", required=True, help="count generated dummy users", type=int)
        parser.add_argument("-path", "--snapshot_path", required=True, help="path to the snapshots directory", type=str)

        args = vars(parser.parse_args())

        conf = OmegaConf.load(args['configuration'])

        self.storage = conf.get("storage")
        self.fake_count = args['count']
        self.snapshots_path = args['snapshot_path']
