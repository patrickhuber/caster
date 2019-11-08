# caster
filesystem templating in go

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

Valid go templates can be used in file names. File names are first url decided and then run through the template engine.
