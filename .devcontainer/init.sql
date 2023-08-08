CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email character varying(100) NOT NULL UNIQUE,
    name character varying(30) NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);

INSERT INTO users (id, email, name) VALUES
    (1, 'tim@devlet.me', 'Tim Devlet')
;

CREATE TABLE IF NOT EXISTS events (
    id SERIAL PRIMARY KEY,
    name character varying(100) NOT NULL,
    payload jsonb NOT NULL DEFAULT '{}',
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);