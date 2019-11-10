# caster

filesystem templating in go

# getting started

Caster is distributed as a single binary. You can download the latest version from the [github releases page](https://github.com/patrickhuber/caster/releases)

You can either place caster in a system path or add the location to your system path.

## creating a template

A template is just a directory with a .caster file at its root that describes the template. 

The .caster file has no required fields

# commands

## show

```bash
caster show git://github.com/patrickhuber/caster/master/examples/simple
```

```yaml
name: simple
sample: "caster cast -l regions.yml"
files:
- name: data.yml
```

## scaffold

```bash
caster scaffold git://github.com/patrickhuber/caster/master/examples/simple
```

```
.
└── data.yml
```

## cast

# templates

## file 

### types

**Caster Metadata**

Allows you to set folder scoped variables. Caster files are not placed in output.

**Caster Molds**

Molds are templates. When a mold is cast the .mold extension is removed and the contents of the file are evaluated through the template engine.

**Regular files**

Regular files are copied verbatim to the target.

### extensions

The following naming conventions are used for files

* **Caster metadata** : \*.caster
* **Caster Molds** : .\*.mold
* **Regular files**: \*.\*

### names

Valid go templates can be used in file and flder names. 

Each line produced by the file or folder template will generate a file. When using ranges, you can specify a backtick ``n` character to force a newline and generate a new file. If you do not specify a newline the entire range of values will be used in the name of a single item. 

You can also use html escape characters to encode quotes and other characters. This is helpful when you need a special character in the file or folder template but that character is not allowed by posix or windows. 

| character | escape   | 
| --------- | -------- |
| `&`       | `&amp;`  |
| `<`       | `&lt;`   |
| `>`       | `&gt;`   |
| `"`       | `&quot;` |
| `'`       | `&apos;` |

given data file:

```yaml
regions:
- name: eastus
- name: westus
```

And a folder named:

```
{{ range .regions }}{{ .name }}`n
```

Will produce output:

```
.
├── eastus
└── westus
```

With the same data file, removing the ``n` produces a folder named:

```
.
└── eastuswestus
```


Empty file names or folder names will not produce any files or folders.

#### child foldes and the `{{end}}` block

If you do not specify a `{{end}}` block one is implicity added after all child items are processed. This allows you to drill down on nested data without needing explicit template variables.

given the data file

```yaml
parents:
- name: one
  children:
  - a
  - b
- name: two
  - c
  - d
```

With an `{{end}}` block, the scope of the current variable is at the root

```
.
└── {{ range .parents }}{{ .name }}{{end}}`n
    └── {{ range .parents }}{{ .name }}`n
```

generates output:

```
.
├── one
│   ├── one
│   └── two
└── two
    ├── one
    └── two
```

Without an `{{end}}` block, one is impliclity added to after all of the children of the root folder are traversed.

```
.
└── {{ range .parents }}{{ .name }}`n
    └── {{ range .children }}{{ . }}`n
```

generates output:

```
.
├── one
│   ├── a
│   └── b
└── two
    ├── c
    └── d
```

# data

Data files are files that contain json or yml data for the generator. The root elements are loaded as a variables. 

To create a data file, name the file `.caster.yml` or `.caster.json` and place it in the directory you want to generate. Users can also supply data files to the `cast` command with the -l flag.

## Examples

### Create a string variable named `iaas` at the root

```yaml
iaas: "azure"
```

```json
{"iaas": "azure"}
```

### Merging and Scoping

#### Merging

When you specify both a `.caster.json` and `.caster.yml` file they are merged into a single variable list. 

If variable names collide, the last file loaded will overwrite the values of the previous file loaded.

```json
{"name": "world"}
```

```yaml
name: hello
```

results in

```
name = "world"
```

User supplied values will always overwrite template supplied values. 

#### Scoping

If a `.caster.yml` or `.caster.json` file is specified in a directory, the values it provides are avaialble only in the current directory or child directories. 

If a variable is defined in child directory with the same name as a variable defined in the parent directory, the variable in the child directory will mask the value of the variable in the parent directory.