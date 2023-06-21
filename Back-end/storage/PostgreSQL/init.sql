DROP DATABASE IF EXISTS bachelordeploma;
DROP ROLE IF EXISTS bachelordeploma;
-- DROP ROLE IF EXISTS viktor;

CREATE ROLE bachelordeploma WITH PASSWORD 'bachelordeploma';
-- CREATE ROLE viktor WITH PASSWORD 'viktor';
ALTER ROLE bachelordeploma WITH LOGIN superuser;
-- ALTER ROLE viktor WITH LOGIN superuser;

CREATE DATABASE bachelordeploma
WITH OWNER = bachelordeploma
ENCODING = 'UTF8'
TABLESPACE = pg_default
CONNECTION LIMIT = -1;
