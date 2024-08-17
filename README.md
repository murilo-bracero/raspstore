# RASPSTORAGE

##  ⚠️ [DISCLAIMER] ⚠️

This is a pet project of mine and it is still in development, so it is not intended to be used in production. Having said that, please feel free to suggest new features or fix some bugs.

At this moment, raspstore requires an external IDP to manage users, like (Keycloak)[https://www.keycloak.org/] or (Firebase Authentication)[https://firebase.google.com/docs/auth].

## DESCRIPTION

Rasptorage is a storage application designed for personal private clouds (e.g. Home Servers).

It is built to require minimal maintenance and configuration, ensuring a simple "install-and-use" experience.

Rasptorage features a built-in UI for convenient file upload and download operations. Additionally, it allows users to manage storage settings, both on an individual user level and for the entire server.

## Software Architecture

The service itself is written in Go, the UI in Typescript using the Svelt framework and the mobile app is built using Flutter.

## Requirements

- An IDP that supports JWK (e.g. Keycloak, Firebase Authentication, Azure AD B2C, etc)
- SO: Linux
- Memory: TBD
- Space: TBD

## Installing

### Docker

TBD

### Manually

TBD

### Build from source

TBD

