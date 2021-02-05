# Dormouse

Dormouse is a stupidly simple tool, that builds and runs a simple CLI that wraps
existing tools and scripts based on a YAML configuration file. Originally
designed to work together with shell aliases, but can be used for all sorts of
things.

# Usage

Put a YAML file anywhere on your system, and point dormouse to it to execute
whatever command or script you have configured in the file. In the config, you
can specify either `exec` to execute a binary file with a list of arguments, or
`script`, to evaluate a shell script. Additional command line arguments are
either appended to the exec argument list, or passed as parameters to the shell
script. For example, the following configurations would be equivalent:

```yaml
---
script: |
  echo "$@"

---
exec: ["/bin/sh", "-c", "echo $@", "--"]

---
exec: ["echo"]

```

All would produce the same result:

```bash
$ dormouse config.yaml foo bar
foo bar
```

You can also require option flags and positional arguments. The `exec` and
`script` fields will be templated, and you can insert the option and argument
values via the `option` and `arg` template functions. For example:

```yaml
---
arguments:
  - name: firstname

options:
  - name: time
    short: t
    default: day

exec: ["echo", "Good {{ option `time` }}, {{ arg `firstname` }}!"]
```

```bash
$ dormouse config.yaml Steve
Good day, Steve!

$ dormouse config.yaml Steve --time=morning
Good morning, Steve!

$ dormouse config.yaml Steve -t eveninng
Good evening, Steve!
```

You can create a whole command tree by adding subcommands. Each command config
can have a `subcommands` field, which is a map from command names to a new
command config. Which in turn may have it's own subcommands and so on.

```yaml
---
subcommands:
  everyone:
    exec: ["echo", "Hello everyone!"]
```

```bash
$ dormouse config.yaml everyone
Hello everyone!
```

To wrap things up, we now save our config somewhere useful and set up a shell
alias for it, and suddenly we have our very own CLI, without a single line of code:

```bash
$ alias greet="dormouse ~/.dormouse/greet.yaml"

$ greet Steve
Good day, Steve!

$ greet everyone
Hello everyone!
```

## License & Authors

```text
Copyright 2021 Ole Claussen

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
