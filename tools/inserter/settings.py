import argparse
from pathlib import Path

from omegaconf import OmegaConf


class Configuration:
    """
    Main configuration.
    Configuration is based on program, arguments: how many users should be generated and path where their will persist
    in files;
    """
    def __init__(self):
        self.storage = None
        self.snapshots_path = None
        self.batch_size = None

    def parse(self) -> None:
        """
        Parsing program args.
        :return:
        """
        parser = argparse.ArgumentParser()
        parser.add_argument("-cfg", "--configuration", required=True, help="path to the config file", type=str)
        parser.add_argument("-path", "--snapshot_path", required=True, help="directory path where snapshots are stored",
                            type=str)
        parser.add_argument("-size", "--batch_size", required=True, help="size of batching to insert into db", type=int)

        args = vars(parser.parse_args())

        conf = OmegaConf.load(args['configuration'])

        self.storage = conf.get('storage')
        self.snapshots_path = Path(args['snapshot_path'])
        self.batch_size = args['batch_size']
