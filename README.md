# flowcast

A small workflow automation tool for running steps against a set of remote
Linux hosts over SSH. Define an inventory of hosts and a list of steps in a
YAML file and flowcast runs every step against every host in parallel.

## Usage

```
go run cmd/main.go [-k ssh_key] <workflow.yaml>
```

or, once built:

```
flowcast [-k ssh_key] <workflow.yaml>
```

### Flags

| Flag | Shorthand | Default | Description |
|---|---|---|---|
| `-ssh_key` | `-k` | `id_ed25519` | Name of the private key to use for SSH auth. If the value contains a `/`, it's used as a path as-is. Otherwise it's looked up in `~/.ssh/`. |

The workflow file is a required positional argument and must come **after**
the flags, e.g. `flowcast -k my_key workflow.yaml`. 

Exit code is `0` only if every host completed its steps without error.

## Workflow file format

A workflow file has two top-level sections: `inventory` and `steps`.

```yaml
inventory:
  - host: 10.0.0.5
    user: root
    port: 22

  - host: 10.0.0.6
    user: root

steps:
  - name: echo
    value: "hello"
```

### `inventory`

A list of hosts to run the workflow against. Each host runs the full list
of steps independently and concurrently. A failure on one host does not
stop the others.

| Field | Required | Description |
|---|---|---|
| `host` | yes | IP address or hostname |
| `user` | yes | SSH user |
| `port` | no | defaults to `22` |

### `steps`

A list of steps, executed in order for each host.

#### `echo`

Logs a value. Mostly useful for testing and for setting a literal value via
`save_as`.

```yaml
- name: echo
  value: "hello world"
```

#### `shell`

Runs a shell command on the remote host over an SSH session.

```yaml
- name: shell
  cmd: "ls -la /var/log"
```

Its result (usable via `save_as`, see below) is a map with three keys:
`stdout`, `stderr`, and `exitCode`. A non-zero exit code does **not** fail
the step or the workflow on its own — check `result.exitCode` explicitly
if you need to react to it.

#### `set_var`

Does nothing on its own; it's a no-op step used purely as a place to attach
a `save_as` block, typically to set an initial variable.

```yaml
- name: set_var
  save_as:
    count: 1
```

#### `http`

Currently a stub — it logs the intended request but does not actually make
one.

```yaml
- name: http
  url: https://example.com
  method: GET
```

#### `if`

Conditional branching. `condition` is evaluated as an
[expr](https://github.com/expr-lang/expr) expression and must evaluate to a
boolean. Runs the `then` steps if true, `else` if false. Both `then` and
`else` are lists of ordinary steps, and either can be omitted.

```yaml
- name: if
  condition: "result.exitCode == 0"
  then:
    - name: shell
      cmd: "echo ok"
  else:
    - name: shell
      cmd: "echo failed"
```

#### `for`

A conditional loop. Before each iteration, `condition` is evaluated the
same way as `if`. While it's true, the `do` steps run, then the condition
is checked again. Stops as soon as `condition` evaluates to false.

```yaml
- name: for
  condition: "count < 3"
  do:
    - name: shell
      cmd: "echo counting"
      save_as:
        count: "count + 1"
```

There is currently no iteration limit — a condition that never becomes
false will loop forever.

### Variables (`save_as`)

Every step can optionally have a `save_as` block, mapping a variable name
to an [expr](https://github.com/expr-lang/expr) expression. After the step
runs, each expression is evaluated and the result is stored under that
name in a variables map, shared by the whole workflow run for that host.

The expression is evaluated against an environment containing:

- every variable saved so far (referenced by bare name, e.g. `count`)
- `result`, the value the step just produced (only set while `save_as` for
  that step is being evaluated)

```yaml
- name: shell
  cmd: "df -h / | tail -1"
  save_as:
    disk_line: result.stdout
    ok: "result.exitCode == 0"
    limit: 90          # a literal constant works too, not just a field lookup
```

Variables set this way are visible to every later step on that host,
including inside `if`/`for` conditions and `{{ }}` templating (see below).
They are **not** shared across hosts — each host runs with its own,
independent variable map.

### Templating (`{{ }}`)

Any string field in a step (except `condition`, `do`, `then`, and `else`,
where you can use the bare names of the saved variables) can contain 
one or more `{{ expression }}` placeholders. Before the step runs, 
each placeholder is evaluated as an expr expression against the current 
variables and substituted into the string.

```yaml
- name: shell
  cmd: "echo iteration {{ count }} >> log.txt"
```

A placeholder referencing an undefined variable, or one that fails to
parse, is an error and stops the workflow for that host — it does not
silently insert an empty string.

## Testing against a local target

`Dockerfile` builds a lightweight Alpine container running `sshd`, useful
for testing without a real remote machine. It copies in a **public** key
(`key.pub`) as the only authorized key and runs as root.
The key needs to be stored in the same directory as the Dockerfile. 

```
docker build -t flowcast-test .
docker run -d -p 2222:22 container-1 
docker run -d -p 2223:22 container-1
```

Then point a workflow's inventory at `127.0.0.1` on those ports (see
`test_workflow.yaml` for an example), and make sure `key.pub` is the public
half of whichever private key you pass via `-k`.

## Known limitations (v1)

- Host key verification is disabled (`InsecureIgnoreHostKey`) — fine for
  testing against a local container.
- Passphrase-protected private keys are not supported.
- Steps are parsed lazily, as they run — a typo or invalid step further
  down a workflow will only surface after earlier steps have already run
  on every host, rather than being caught up front.
- The `http` step is not implemented.
- `for` loops have no iteration cap and no timeout; a bad condition can
  loop forever.
- Individual shell commands have no timeout; a hung remote command blocks
  that host's workflow indefinitely.
