# rex

Rex is an ADR cli tool that helps create and manage ADRs in a project. Goals of this tool is to have a simple cli that can do that as well as generate HTML from the markdown ADR files to host using something like GitHub pages.

### Why another ADR tool?
All of them are uniquely built to satisfy the needs of how they want to use ADRs. I built this one so that I can build and mmanage them the way I want to. :D


## Using Rex
get the binary for your machine. Goto the root level of your projects directory and run the command:
```
rex generate config
```

That will create the default config file. Now you can run the `init` command to initialize your project.
```
rex init
```

This will create 2 directories:
- `docs/adr`
- `templates`

`docs/adr` is where your ADR's will be created in.
`templates` is where the defualt ADR template is and where you can make changes to it.

To create a new ADR, run:
```
rex create -t "My ADR"
```
`-t` is a flag to specify the title of your ADR.


### Todos 
- [ ] Auto update a README.md in the ADR dir with a list of all ADR's and links
- [ ] Use some library to take ADR markdown and generate HTML to push to a Github Pages site.
- [ ] Create a simple web server to host the generate HTML files locally.
