# IATA-FINDER

Special thanks to openflights for making this high quality data readily available for use in multiple formats.

This repository attempts to follow golang [common project structure](https://github.com/golang-standards/project-layout)git c.

This repository contains:

1. A gRPC proto spec for IATA-FINDER
2. GCP Deployement Specs
3. TODO: ETL auto update scripts & GCP Deployment Specs

## Table of Contents

- [Background](#background)
- [Install](#install)
- [Usage](#usage)

## Background

I was having trouble finding an API that easily made available airport & airline codes with associated information for free. I thought I should use this chance to make my own project exploring some development concepts I have not yet had the chance to touch, such as API development.

## Install

## Usage

### Postgres Deployment
Find the `gcp_postgres` folder in the deployments main dir. In there, add a connection.config.json file following the specs in the connectino.config.sample.json.
At current the DB is not using SSL, but I plan to add in SSL as soon as I figure out how to config `pgx` to accept ssl certs for it's connection statement.
With assurance that you have the correct airlines.csv and airports.csv in the `/assets/` folder, you can run the db.go file and it will copy the csv contents
to your specificed db. Note that the copy happens to the postgres public schema, this can easily be changed by specifing the schema in the copy statements.
