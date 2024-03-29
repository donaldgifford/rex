# rex

Rex is an ADR cli tool that helps create and manage ADRs in a project. Goals of this tool is to have a simple cli that can create markdown ADR's as well as generate HTML from the markdown ADR files to host using GitHub pages.

### Why another ADR tool?

All of them are uniquely built to satisfy the needs of how they want to use ADRs. I built this one so that I can build and manage them the way I want to. :D

### Primary Features

2 Options when using the cli tool:

1. Create markdown ADR's in the root of the project
2. Create HTML from markdown ADR's to use for hosting on GitHub pages.

Everything is configurable using the the `.rex.yaml` config file.

This repo dogfoods using `rex` with all features enabled. The Github Page generated can be found here: [Rex Github Page](https://donaldgifford.github.io/rex)

### Install

`go install github.com/donaldgifford/rex@latest` or whatever version you want to use.
