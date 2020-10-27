import pymysql.cursors
from settings import Configuration
from dummy import Dummy


def app():
    cfg = Configuration()
    cfg.parse()

    import time
    start_time = time.time()
    dummy = Dummy(cfg.fake_count)
    dummy.generate()
    dummy.make_snapshot(cfg.snapshots_path)
    print("--- %s seconds ---" % (time.time() - start_time))


if __name__ == "__main__":
    app()

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
#
#     fake = Faker()
#     interests = []
#     # print_hi(names.get_full_name(gender=male) + ' ' + male + ' ' + fake.password() + ' ' + fake.email())
#     print_hi(fake.email() + ' ' + fake.password() + ' ' + fake.first_name() + ' ' + fake.last_name() + ' ' +
#              fake.profile().get('sex') + ' ' + fake.date_of_birth(minimum_age=18, maximum_age=100).strftime("%m/%d/%Y")
#              + ' ' + fake.city() + ' ' + fake.text(max_nb_chars=100))
#
