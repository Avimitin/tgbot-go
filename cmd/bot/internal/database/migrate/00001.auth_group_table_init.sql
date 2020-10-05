CREATE TABLE IF NOT EXISTS AuthGroups (
    ID              int AUTO_INCREMENT,
    GroupID         int NOT NULL,
    GroupUsername   VARCHAR(50) NOT NULL,
    PRIMARY KEY (ID)
)CHARSET=utf8