# RASPSTORAGE

## DESCRIPTION

A cloud-like storage app for raspberry pi clusters, with automatic installation, low-maintenance and low-config, with user management, private folders and files, encryption at rest, light and resilient.

## Goals

- v0.1:
    - Clean Architecture
    - User Management (Authentication)
    - Shared File Storage

- v0.2:
    - User Management (Authorization)
    - Per-User File Storage
    - Secret Folders

- v0.3:
    - Encryption at rest
    - Auto deleted files
    - Stoage and files metadata

- v0.4:
    - Audit system
    - File sharing between users

- v0.5:
    - Backup plans
    - Auto maintenance windows
    - File integrity checkup
    - File integrity report
    - Environment security

## Architecture

- Programming Languages: 
    - Go
    - JS
    - Dart

- Frameworks:
    - Go:
        - MUX (Rest API)

    - JS:
        - ReactJS

    - Dart:
        - Flutter

- Databases:
    - MongoDB (User Data Storage)
    - Firebase (User Credentials Storage)
    - Cassandra (History and file metadata fast recovery)

- Infrastructure:
    - Application is hosted in Raspberry Pi infrastructure
    - MongoDB is hosted in Atlas
    - Firebase is a Clud Service, hosted at Google Cloud
    - Cassandra is a cloud agnostic database, hosted in whatever it suited better, including in your own raspberry pi infrastructure
