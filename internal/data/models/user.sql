CREATE TABLE IF NOT EXISTS users (
    id serial NOT NULL,
    name VARCHAR(150) NOT NULL,
    password varchar(120) NOT NULL UNIQUE CHECK (password <> '' ),
    addres varchar(500) NOT NULL, 
    email VARCHAR(150) NOT NULL UNIQUE CHECK (email <> '' AND email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'),
    phone VARCHAR(20) NOT NULL UNIQUE CHECK (phone <> '' AND phone ~* '^\+[0-9]+$'),
    created_at timestamp DEFAULT now(),
    updated_at timestamp NOT NULL,
    CONSTRAINT pk_users PRIMARY KEY(id)
);
