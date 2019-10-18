# IATA-FINDER

Special thanks to openflights for making this high quality data readily available for use in multiple formats.

This repository contains:

1. A gRPC proto spec & implementation for iata-finder
2. A dataservice which fetches a .csv of airports and airline data from open flights, sanitizes and stores in a poor-mans memory cache

## Table of Contents

- [Background](#background)
- [Install](#install)
- [Usage](#usage)
- [Todo](#todo)

## Background

I was having trouble finding an API that easily made available airport & airline codes with associated information for free. I thought I should use this chance to make my own project exploring some development concepts I have not yet had the chance to touch, such as API development.

## Install

## Todo

- Implement struct for managing the db connection in the ETL script.
- Organize the execution of the copy protocol in a better way
