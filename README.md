# phishflood

<p align="center">
  <p align="center">
    <a href="https://github.com/andpalmier/phishflood/blob/main/LICENSE"><img alt="Software License" src="https://img.shields.io/badge/license-GPL3-brightgreen.svg?style=flat-square"></a>
    <a href="https://goreportcard.com/report/github.com/andpalmier/phishflood"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/andpalmier/phishflood?style=flat-square"></a>
    <a href="https://twitter.com/intent/follow?screen_name=andpalmier"><img src="https://img.shields.io/twitter/follow/andpalmier?style=social&logo=twitter" alt="follow on Twitter"></a>
  </p>
</p>


This is a proof of concept to pollute with fake data the credentials stolen with a phishing kit. An old version of this project was discussed [in this blog post](https://andpalmier.github.io/posts/flooding-phishing-kits/), I decided to use this repository to add some features.

## Usage

After downloading the repository, build the project:

```
$ make phishflood
```

This will create a folder `build` and an executable `phishflood`. You can then run the executable with the following flags:

- `-goroutines int`: number of goRoutines (default 10).
- `-proxies string`: one or multiple proxies; specify the schema (http default) and port, and use ',' as a separator.
- `-ua string`: User Agent to be used, using Chrome on iPhone by default.
- `-url string`: domain name or url, if schema is not specified, https is assumed.

### Todo

- [x] Improve code organization
- [x] Add custom User Agent flag
- [x] Remove colly dependency
- [ ] Add compatibility with known phishing kits
- [ ] Use [faker](https://github.com/bxcodec/faker) for data generation
