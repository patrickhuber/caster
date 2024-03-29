# caster

A file and directory templating cli

# getting started

Caster is distributed as a single binary. You can download the latest version from the [github releases page](https://github.com/patrickhuber/caster/releases)

You can either place caster in a system path or add the location to your system path.

You can also install caster with `go install`

```bash
git clone https://github.com/patrickhuber/caster
go install -C caster/cmd/caster
```

Once installed, you can run the hello world example by typing 

```bash
caster init
caster apply
```

This creates a .caster.yml file in the current directory. The apply command will run the template in the current directory.

## creating a template

A template is just a directory with a .caster.yml or .caster.json file at its root that describes the template. 

The .caster file has no required fields

See the hello world example [here](https://github.com/patrickhuber/caster/tree/main/examples/simple) for a sample hello world caster example.

