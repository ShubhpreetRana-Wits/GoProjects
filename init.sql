DO
$$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'your_database_name') THEN
        CREATE DATABASE your_database_name OWNER your_db_user;
    END IF;
END
$$;