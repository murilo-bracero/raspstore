# RASPSTORAGE

## DESCRIPTION

A cloud-like storage app for raspberry pi cluster, with automatic installation, low-maintenance and low-config, with user management, private folders and files, encryption at rest, light and resilient.

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

- v0.6:
    - Native Web Application

- v0.7:
    - Native Mobile App

- v0.8:
    - Architecture plugins

-v0.9:
    - Plugins
    - Auto infrastructure deployment using scripts and terraform

## Software Architecture

- Programming Languages: 
    - Go (Backend)
    - JS (Web Frontend)
    - Dart (Mobile)

## Recommended Architecture

- Databases:
    - MongoDB (User Data Storage)
    - Firebase (User Credentials Storage)
    - Stackdriver (History and Logging)
    - Cassandra (Metadata fast recovery)

- Infrastructure:
    - Application is hosted in Raspberry Pi infrastructure
    - MongoDB hosted at Atlas
    - Firebase is a Cloud Service, hosted at Google Cloud
    - Stackdriver is a managed logging system hosted at Google Cloud
    - Cassandra self-hosted or hosted at some Cloud environment

Disclaimer: I highly recommend that you use cloud environments if you are looking for low-config and low-management, but there is no problem to host all these thing in your own infrastructure if you have a well-suited raspberry pi cluster

## Plugins

- Datastore (everything)
- DynamoDB (everything)
- Firestore (everything)
- Cloud Pub/Sub (History and Logging)
- Kafka (History and Logging)
- RabbitMQ (History and Logging)
