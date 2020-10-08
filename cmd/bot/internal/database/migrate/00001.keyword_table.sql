CREATE TABLE IF NOT EXISTS Keywords (
    kid INT AUTO_INCREMENT,
    keywords VARCHAR(255) NOT NULL,
    count int DEFAULT 0,
    PRIMARY KEY (kid)
)