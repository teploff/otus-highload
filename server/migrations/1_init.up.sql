CREATE TABLE IF NOT EXISTS user (
    id            VARCHAR(100)  NOT NULL DEFAULT (uuid()) PRIMARY KEY,
    email         VARCHAR(100) NOT NULL UNIQUE CONSTRAINT email_empty_field CHECK (email != ''),
    password      TEXT         NOT NULL CONSTRAINT password_empty_field CHECK (password != ''),
    name          VARCHAR(100) NOT NULL CONSTRAINT name_empty_field CHECK (name != ''),
    surname       VARCHAR(100) NOT NULL CONSTRAINT surname_empty_field CHECK (surname != ''),
    sex           VARCHAR(10)  NOT NULL CONSTRAINT sex_empty_field CHECK (sex != '' and (sex = 'male' or sex = 'female')),
    birthday      DATE         NOT NULL ,
    city          VARCHAR(50)  NOT NULL CONSTRAINT city_empty_field CHECK (city != ''),
    interests     TEXT         NOT NULL CONSTRAINT interests_empty_field CHECK (interests != ''),
    access_token  VARCHAR(350) NULL ,
    refresh_token VARCHAR(350) NULL ,
    create_time   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time   TIMESTAMP    NULL
) ENGINE=InnoDB;