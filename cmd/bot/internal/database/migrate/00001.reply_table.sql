CREATE TABLE IF NOT EXISTS replies(
    rid     INT AUTO_INCREMENT,
    reply   VARCHAR(255) NOT NULL ,
    keyword INT NOT NULL,
    PRIMARY KEY (rid),
    FOREIGN KEY (keyword) REFERENCES keywords(kid)
)