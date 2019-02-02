# DB Structure

The purpose of this script is to create the structure that the database will require for
holding the iata information. The scripts will follow a simple pattern of tear-down, 
creation, and value insertion. With this in mind, if any data manipulation or etl 
happens on these source tables, it -will- be lost when running this script UNLESS
the steps for manipulation are written into a migration-style scripter and ran hereafter.