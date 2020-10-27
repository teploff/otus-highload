from datetime import datetime


class User:
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

    def __repr__(self):
        return f'{self.email}\t{self.password}\t{self.name}\t{self.surname}\t{self.birthday}\t{self.sex}\t' \
               f'{self.city}\t{self.interests}'
