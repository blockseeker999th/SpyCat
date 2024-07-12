
CREATE TABLE spy_cats (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    years_of_experience INT NOT NULL CHECK (years_of_experience >= 0),
    breed VARCHAR(255) NOT NULL,
    salary NUMERIC(10, 2) NOT NULL
);

CREATE TABLE missions (
    id SERIAL PRIMARY KEY,
    spy_cat_id INT NOT NULL,
    status VARCHAR(50) NOT NULL CHECK (status IN ('in_progress', 'completed')),
    FOREIGN KEY (spy_cat_id) REFERENCES spy_cats (id)
);

CREATE TABLE targets (
    id SERIAL PRIMARY KEY,
    mission_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    country VARCHAR(255) NOT NULL,
    notes TEXT,
    completed BOOLEAN NOT NULL,
    FOREIGN KEY (mission_id) REFERENCES missions (id)
);

CREATE INDEX idx_spy_cats_name ON spy_cats (name);
CREATE INDEX idx_missions_status ON missions (status);
CREATE INDEX idx_targets_completed ON targets (name);
