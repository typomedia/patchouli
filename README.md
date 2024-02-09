# Patchouli - Stay secure

[![Go Report Card](https://goreportcard.com/badge/github.com/typomedia/patchouli)](https://goreportcard.com/report/github.com/typomedia/patchouli)
[![Go Reference](https://pkg.go.dev/badge/github.com/typomedia/patchouli.svg)](https://pkg.go.dev/github.com/typomedia/patchouli)
[![GitHub release](https://img.shields.io/github/release/typomedia/patchouli.svg)](https://github.com/typomedia/patchouli/releases/latest)
[![GitHub license](https://img.shields.io/github/license/typomedia/patchouli.svg)](https://github.com/typomedia/patchouli/blob/master/LICENSE)

Patchouli is a lightweight patch management planner for **operating systems**. It comes with an intuitive web interface.

All data will be stored in a single `patchouli.boltdb` file in the current working directory.

## Run

    make run

## Build

    make build

## Cross compile

    make compile

## Technology

- [Go](https://golang.org/)
- [Fiber](https://gofiber.io/)
- [htmx](https://htmx.org/)
- [Pico](https://picocss.com/)
- [bbolt](https://github.com/etcd-io/bbolt)

## Todo

- [ ] add delete functionality
- [ ] add a login page
- [ ] add a `toml` config file
- [ ] email notifications
- [ ] protect api used by the frontend
- [ ] optimize json/csv export
- [ ] add json/csv import api
- [ ] refactor some quirky code
- [ ] write some tests...

## Systemd

You can use the provided [patchouli.service](patchouli.service) file to run `patchouli` as a daemon. 
If not running in a container, change the `User` and `Group` for security reasons.

    sudo mkdir /var/patchouli
    sudo cp patchouli.service /etc/systemd/system/patchouli.service
    sudo systemctl daemon-reload
    sudo systemctl enable --now patchouli
    sudo systemctl status patchouli

## License

Patchouli is licensed under the [GNU General Public License v3.0](LICENSE).

---
Copyright © 2024 Typomedia Foundation. All rights reserved.