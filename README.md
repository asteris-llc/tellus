# Tellus

Tellus is a tool to coordinate [terraform](https://terraform.io) deployments
across your team.

<!-- markdown-toc start - Don't edit this section. Run M-x markdown-toc/generate-toc again -->
**Table of Contents**

- [Tellus](#tellus)
    - [The Name](#the-name)
    - [Storage](#storage)
    - [Installation](#installation)
    - [Commands](#commands)
        - [Server](#server)
        - [Client](#client)
    - [Configuration](#configuration)

<!-- markdown-toc end -->

## The Name

Tellus, in Latin, means "earth". This project is focused on providing a safe
home for you and your team while you're terraforming other planets (or servers,
whatever.)

## Storage

By default Tellus uses a disk-backed store. Ideally it'll support a
highly-available store in the not-to-distant future. For module storage, Tellus
uses [Vault](https://vaultproject.io/) (so you can either store your
authentication credentials inline safe in the knowledge that they're encrypted
on the server, or you can use environment variables or tfvars like normal.)

## Installation

`go get github.com/asteris-llc/tellus`

## Commands

### Server

`tellus server`

### Client

`tellus init https://tellus.yourdomain.com/your-project` in a directory with
your `.tf` files will set up Terraform with the appropriate remote syncing
settings. Afterwards, `tellus push` will do you just fine, and will push both
your Terraform state files and module config off to the server.

To pull down the files on a new machine, run the init again, and then hit
`tellus pull`.

## Configuration

It'll be in here when I write it! (AKA TODO)
