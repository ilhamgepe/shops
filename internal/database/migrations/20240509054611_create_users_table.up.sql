CREATE TABLE users(
    id BIGINT NOT NULL AUTO_INCREMENT,
    first_name varchar(100) NOT NULL,
    last_name varchar(100),
    email varchar(255) NOT NULL,
    password varchar(255) NOT NULL,
    role ENUM('users','admin') DEFAULT 'users',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP ON UPDATE CURRENT_TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (email),
    PRIMARY KEY (id)
);