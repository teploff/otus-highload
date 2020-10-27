import uuid

from faker import Faker

from entity import User


class Dummy:
    """

    """

    def __init__(self, count: int):
        self.count = count
        self.users = []

    def generate(self):
        """

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

    def make_snapshot(self, path: str):
        """

        :return:
        """
        file_name = uuid.uuid4()

        with open(path + str(file_name) + '.txt', 'w') as f:
            for user in self.users:
                f.write(str(user))
                f.write('\n')
