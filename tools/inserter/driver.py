from typing import List

import pymysql

from entity import User


class MySQLDriver:
    def __init__(self, settings: dict):
        self.connection = pymysql.connect(
            host=settings['host'],
            port=settings['port'],
            user=settings['user'],
            passwd=settings['password'],
            db=settings['db'],
            charset=settings['charset'],
            cursorclass=pymysql.cursors.DictCursor)
        self.written_users = 0

    def insert(self, users_set: List[List[User]]):
        try:
            with self.connection.cursor() as cursor:
                for users in users_set:
                    for user in users:
                        sql = '''
                        INSERT 
                            INTO user (email, password, name, surname, birthday, sex, city, interests)
                        VALUES 
                            ( %s, %s, %s, %s, %s, %s, %s, %s)'''
                        cursor.execute(sql, (
                            user.email,
                            user.password,
                            user.name,
                            user.surname,
                            user.birthday,
                            user.sex,
                            user.city,
                            user.interests))
                    self.connection.commit()
                    self.written_users += len(users)
        finally:
            self.connection.close()
