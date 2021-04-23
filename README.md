<div align="center">
    <img src="misc/logo-with-text.png?raw=true?" alt="herald-logo" width="250">
    <h3>announce your samples</h3>
    <hr>
    <a href="https://github.com/will-rowe/herald/actions/workflows/tests.yml"><img src="https://github.com/will-rowe/herald/actions/workflows/tests.yml/badge.svg" alt="actions"></a>
    <a href='http://herald-docs.readthedocs.io/en/latest/?badge=latest'><img src='https://readthedocs.org/projects/herald/badge/?version=latest' alt='Documentation Status' /></a>
    <a href="https://goreportcard.com/report/github.com/will-rowe/herald"><img src="https://goreportcard.com/badge/github.com/will-rowe/herald" alt="reportcard"></a>
    <a href="https://github.com/will-rowe/herald/blob/master/LICENSE"><img src="https://img.shields.io/badge/license-MIT-orange.svg" alt="License"></a>
</div>

---

```
this is a work in progress and currently unstable
```

## Overview

The basic idea is to announce samples to APIs and track responses.

You create a sample via the **Herald** app. Once you have a sample on record, you can tag it with processes (e.g. sequence it, analyse it, upload it...). You can then tell **Herald** to announce the sample to the tagged processes, which it will monitor and update the sample record accordingly.

When announcing a sample, **Herald** will:

- check the sample record
- evaluate the tagged processes and create an execution order
- formulate the correct [gRPC]() messages and submit them to the process APIs
- wait for responses, update the sample record and submit the next message

## Installation

### Use a release

**Herald** is packaged as a desktop application (using [lorca](https://github.com/zserge/lorca)). Just download a [release]() for your platform.

> note: lorca apps require [Chrome/Chromium](https://www.google.com/chrome/) >= 70 to be installed on your system.

### Building from source

You will need the [Go tool chain](https://golang.org/) (**Herald** tested with v1.16) to build from source.

```sh
git clone https://github.com/will-rowe/herald
cd herald
make all
```

## Documentation

Docs are available via [read the docs](http://herald-docs.readthedocs.io/en/latest/?badge=latest) and are being written during development.
