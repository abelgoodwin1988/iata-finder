/**
 * Script for creating/recreating the table structure for the data
 */

-- Drop table's sequence in reverse order of creation.

DROP TABLE IF EXISTS airline;
DROP TABLE IF EXISTS airport;

-- Create tables sequence

CREATE TABLE airport (
     id                     INTEGER
    ,name                   VARCHAR
    ,city                   VARCHAR
    ,country                VARCHAR
    ,iata                   CHAR(3)
    ,icao                   CHAR(4)
    ,latitude               NUMERIC
    ,longitude              NUMERIC
    ,altitude               NUMERIC
    ,timezone               VARCHAR
    ,daylight_savings_time  VARCHAR
    ,tz                     VARCHAR
    ,type                   VARCHAR
    ,source                 VARCHAR
);

CREATE TABLE airline (
     id                     INTEGER
    ,name                   VARCHAR
    ,alias                  VARCHAR
    ,iata                   VARCHAR
    ,icao                   VARCHAR
    ,callsign               VARCHAR
    ,country                VARCHAR
    ,active                 VARCHAR
);