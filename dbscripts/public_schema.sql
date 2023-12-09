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
    scheduled_date text NULL,
    scheduled boolean NOT NULL,
    is_active boolean NOT NULL DEFAULT true,
    company_id text NOT NULL,

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
    dt real NULL,
    pressure real NULL,
    humidity real NULL,
    sunrise real NULL,
    sunset real NULL,
    speed real NULL,
    deg real NULL,
    clouds real NULL,
    rain real NULL,
    snow real NULL,
    icon text NULL,
    description text NULL,
    main text NULL,
    city_id text NULL,
    city_name text NULL,
    country text NULL,
    time_zone real NULL,
    population real NULL,
    latitude real NULL,
    longitude real NULL,
    temp_day real NULL,
    temp_min real NULL,
    temp_max real NULL,
    temp_night real NULL,
    temp_eve real NULL,
    temp_morn real NULL,
    feels_like_day real NULL,
    feels_like_night real NULL,
    feels_like_eve real NULL,
    feels_like_morn real NULL,


    CONSTRAINT weathers_pk PRIMARY KEY (id),
    CONSTRAINT fk_weathers_jobs_id FOREIGN KEY (job_id)
        REFERENCES jobs (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
);

