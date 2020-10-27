from file_reader import Reader
from settings import Configuration


def app():
    cfg = Configuration()
    cfg.parse()

    reader = Reader(cfg.snapshots_path)
    reader.do()
    pass


if __name__ == '__main__':
    app()

# See PyCharm help at https://www.jetbrains.com/help/pycharm/
# # Press the green button in the gutter to run the script.
# if __name__ == '__main__':
#     # Connect to the database
#     connection = pymysql.connect(host='localhost',
#                                  user='user',
#                                  password='passwd',
#                                  db='db',
#                                  charset='utf8mb4',
#                                  cursorclass=pymysql.cursors.DictCursor)
#     try:
#         with connection.cursor() as cursor:
#             # Create a new record
#             sql = "INSERT INTO `users` (`email`, `password`) VALUES (%s, %s)"
#             cursor.execute(sql, ('webmaster@python.org', 'very-secret'))
#
#         # connection is not autocommit by default. So you must commit to save
#         # your changes.
#         connection.commit()
#     finally:
#         connection.close()