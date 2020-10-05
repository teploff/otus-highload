-- +goose Up
CREATE TABLE "user"
(
    id            VARCHAR(32)                                                     NOT NULL DEFAULT (uuid()) PRIMARY KEY,
    login         VARCHAR(100)                                                    NOT NULL UNIQUE
        CONSTRAINT login_empty_field CHECK (login != ''),
    password      TEXT                                                            NOT NULL UNIQUE
        CONSTRAINT password_empty_field CHECK (password != ''),
    name          VARCHAR(100)                                                    NOT NULL
        CONSTRAINT name_empty_field CHECK (name != ''),
    surname       VARCHAR(100)                                                    NOT NULL
        CONSTRAINT surname_empty_field CHECK (surname != ''),
    access_token  VARCHAR(350)             DEFAULT '',
    refresh_token VARCHAR(350)             DEFAULT '',
    birthday      TIMESTAMP WITH time ZONE DEFAULT timezone('utc' :: TEXT, now()) NOT NULL,
    sex           VARCHAR(10)                                                     NOT NULL
        CONSTRAINT sex_empty_field CHECK (sex != '' and (sex == 'male' or sex == 'female')),
    city          VARCHAR(50)                                                     NOT NULL
        CONSTRAINT city_empty_field CHECK (city != ''),
    interests     TEXT                                                            NOT NULL
        CONSTRAINT interests_empty_field CHECK (interests != ''),
    create_time   TIMESTAMP WITH time ZONE DEFAULT timezone('utc' :: TEXT, now()) NOT NULL,
    update_time   TIMESTAMP WITH TIME ZONE                                        NULL
);

CREATE INDEX idx ON "user" (id);
CREATE INDEX loginx ON "user" (login);
CREATE INDEX create_timex ON "user" (create_time);
CREATE INDEX update_timex ON "user" (update_time);

-- +goose Down
DROP TABLE "user";