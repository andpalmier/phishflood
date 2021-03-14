# phishflood

<p align="center">
  <img alt="phishflood" src="https://github.com/andpalmier/phishflood/blob/main/img/phishflood.png?raw=true" />
  <p align="center">
    <a href="https://github.com/andpalmier/phishflood/blob/main/LICENSE"><img alt="Software License" src="https://img.shields.io/badge/license-GPL3-brightgreen.svg?style=flat-square"></a>
    <a href="https://goreportcard.com/report/github.com/andpalmier/phishflood"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/andpalmier/phishflood?style=flat-square"></a>
    <a href="https://twitter.com/intent/follow?screen_name=andpalmier"><img src="https://img.shields.io/twitter/follow/andpalmier?style=social&logo=twitter" alt="follow on Twitter"></a>
  </p>
</p>


This is a proof of concept to pollute phishing kits with fake data. An old version of this project was discussed [in this blog post](https://andpalmier.github.io/posts/flooding-phishing-kits/), I decided to use this repository to add some features.

**PLEASE NOTE: At the moment,** `phishflood` **is compatible only with the phishing kits mentioned in the blog post, and with some others which follows the same structure. This is due to a naif approach for the data generation (it's a PoC 😅).**

## Usage

After downloading the repository, build the project:

```
$ make phishflood
```

This will create a folder `build` and an executable `phishflood`. You can then run the executable with the following flags:

- `-dmax (int)`: maximum delay between consecutive requests, in seconds (default 3600).
- `-dmin (int)`: minimun delay between consecutive requests, in seconds (default 10).
- `-goroutines (int)`: number of goRoutines (default 10).
- `-proxies (string)`: one or multiple proxies; specify the schema (http default) and port, and use ',' as a separator.
- `-seed (int64)`: seed used for random data generation, random if not specified.
- `-url (string)`: domain name or url, if schema is not specified, https is assumed.
- `-ua (string)`: User Agent to be used, using Chrome on iPhone by default.

### Todo

- [x] Improve code organization
- [x] Add custom User Agent flag
- [x] Remove colly dependency
- [x] Started [gofakeit](https://github.com/brianvoe/gofakeit) integration for fake data generation
- [ ] Add compatibility with known phishing kits
