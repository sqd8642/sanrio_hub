CREATE TABLE ShowCharacters (
    id SERIAL PRIMARY KEY,
    show_id INT NOT NULL,
    character_id INT NOT NULL,
    FOREIGN KEY (show_id) REFERENCES Shows(id),
    FOREIGN KEY (character_id) REFERENCES Characters(id)
);