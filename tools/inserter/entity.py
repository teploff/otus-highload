from datetime import datetime


class User:
    """
    User entity.
    """
    def __init__(self, email: str, password: str, name: str, surname: str, birthday: datetime, sex: str, city: str,
                 interests: str):
        self.email = email
        self.password = password
        self.name = name
        self.surname = surname
        self.birthday = birthday
        self.sex = sex
        self.city = city
        self.interests = interests
