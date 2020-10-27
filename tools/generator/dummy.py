from pathlib import Path
import uuid

from faker import Faker

from entity import User


class Dummy:
    """
    Dummy users' generator.
    """

    def __init__(self, count: int, path: Path):
        """
        count - user count, which should be generated
        path - path to directory, where generated users should be stored
        """
        self.count = count
        self.path = path
        self.users = []

    def generate(self) -> None:
        """
        Generate <count> users using faker library.
        :return:
        """

        for _ in range(self.count):
            fake = Faker()
            self.users.append(User(
                fake.email(),
                fake.password(),
                fake.first_name(),
                fake.last_name(),
                fake.date_of_birth(minimum_age=18, maximum_age=100),
                'male' if fake.profile().get('sex') == 'M' else 'female',
                fake.city(),
                fake.text(max_nb_chars=100)
            ))

    def make_snapshot(self) -> None:
        """
        Persist users to *.txt file.
        Name of file gets from uuid.
        """
        file_name = str(uuid.uuid4()) + ".txt"

        with open(self.path / file_name, 'w') as f:
            for user in self.users:
                f.write(str(user))
                f.write('\n')
