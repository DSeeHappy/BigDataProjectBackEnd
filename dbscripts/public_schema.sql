-- Initial public schema relates to Library 0.x

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET client_min_messages = warning;
SET row_security = off;

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA pg_catalog;
CREATE EXTENSION IF NOT EXISTS pgcrypto;

SET search_path = public, pg_catalog;
SET default_tablespace = '';

-- jobs - static company id for now
CREATE TABLE jobs (
    id uuid NOT NULL DEFAULT uuid_generate_v1mc(),
    name text NOT NULL,
    address text NULL,
    city text NULL,
    state text NOT NULL,
    zip_code text NOT NULL,
    country text NULL,
    latitude text NULL,
    longitude text NULL,
    scheduled_date date NULL,
    scheduled boolean NOT NULL,
    is_active boolean NOT NULL DEFAULT true,
    company_id int NOT NULL DEFAULT 1,

    CONSTRAINT jobs_pk PRIMARY KEY (id)
);

CREATE INDEX jobs_city
ON jobs (city);

CREATE INDEX jobs_zip_code
ON jobs (zip_code);

-- weather
CREATE TABLE weathers (
    id uuid NOT NULL DEFAULT uuid_generate_v1mc(),
    job_id uuid NOT NULL,
    pressure text NULL,
    humidity text NULL,
    sunrise text NULL,
    sunset text NULL,
    speed text NULL,
    deg text NULL,
    clouds text NULL,
    rain text NULL,
    snow text NULL,
    icon text NULL,
    description text NULL,
    main text NULL,
    city_id text NULL,
    city_name text NULL,
    country text NULL,
    time_zone text NULL,
    population text NULL,
    latitude text NULL,
    longitude text NULL,
    temp_day text NULL,
    temp_min text NULL,
    temp_max text NULL,
    temp_night text NULL,
    temp_eve text NULL,
    temp_morn text NULL,
    feels_like_day text NULL,
    feels_like_night text NULL,
    feels_like_eve text NULL,
    feels_like_morn text NULL,


    CONSTRAINT weathers_pk PRIMARY KEY (id),
    CONSTRAINT fk_weathers_jobs_id FOREIGN KEY (job_id)
        REFERENCES jobs (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
);

