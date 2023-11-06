CREATE TABLE IF NOT EXISTS tokens (
    id serial NOT NULL,
    user_id integer NOT NULL,
    token_value VARCHAR(255) NOT NULL,
    expiration_date timestamp NOT NULL,
    created_at timestamp DEFAULT now(),
    CONSTRAINT pk_tokens PRIMARY KEY(id),
    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id)
);

CREATE OR REPLACE FUNCTION activate_user_session()
RETURNS TRIGGER AS $$
BEGIN
    -- Cambiar el estado de sesión del usuario a activo
    UPDATE users SET session_active = true WHERE id = NEW.user_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION delete_expired_tokens()
RETURNS VOID AS $$
BEGIN
    -- Eliminar los tokens expirados
    DELETE FROM tokens WHERE expiration_date < now();

    -- Cambiar el estado de sesión de los usuarios correspondientes a falso
    UPDATE users SET session_active = false WHERE id IN (
        SELECT user_id FROM tokens WHERE expiration_date < now()
    );
END;
$$ LANGUAGE plpgsql;

-- Crear un trigger para ejecutar activate_user_session después de crear un token
CREATE TRIGGER activate_user_session_trigger
AFTER INSERT ON tokens
FOR EACH ROW
EXECUTE FUNCTION activate_user_session();

-- Crear una tarea cron para ejecutar delete_expired_tokens cada cierto tiempo
CREATE OR REPLACE FUNCTION run_delete_expired_tokens()
RETURNS VOID AS $$
BEGIN
    PERFORM delete_expired_tokens();
END;
$$ LANGUAGE plpgsql;

-- -- Configurar la tarea cron para que se ejecute cada 5 minutos
-- CREATE EVENT TRIGGER delete_expired_tokens_cron
-- ON SCHEDULE EVERY '5 minutes'
-- DO
-- $$
-- BEGIN
--     PERFORM run_delete_expired_tokens();
-- END;
-- $$;
