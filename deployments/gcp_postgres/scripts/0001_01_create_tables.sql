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
    ,latitude               NUMBER
    ,longitude              NUMBER
    ,altitude               NUMBER
    ,timezone               VARCHAR
    ,daylight_savings_time  CHAR(1)
    ,timezone               VARCHAR
    ,type                   VARCHAR
    ,source                 VARCHAR
);

CREATE TABLE airline (
     id                     INTEGER
    ,name                   VARCHAR
    ,iata                   CHAR(2)
    ,icao                   CHAR(3)
    ,callsign               VARCHAR
    ,country                VARCHAR
    ,active                 CHAR(1)
);