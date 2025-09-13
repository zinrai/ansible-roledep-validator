# Ansible Role Dependency Validator

`ansible-roledep-validator` is a command-line tool that helps identify missing roles in your Ansible project. It analyzes playbooks, detects required roles (including dependencies defined in meta/main.yml files), and compares them with the roles available in your roles directory.

## Installation

Build the binary:

```
$ go build
```

## Usage

Basic usage with default roles directory:

```bash
$ ansible-roledep-validator playbook.yml
```

Specify a custom roles directory:

```bash
$ ansible-roledep-validator -playbook playbook.yml -roles /path/to/roles
```

Verbose output with additional details:

```bash
$ ansible-roledep-validator -playbook playbook.yml -verbose
```

## Output

By default, the tool produces minimal output, listing only the missing roles:

```
role1
role2
role3
```

With the `-verbose` flag, the tool provides detailed information about the validation process:

```
=== Ansible Role Dependency Validator ===
Playbook: playbook.yml
Roles directory: roles

Current working directory: /path/to/project

...detailed logging...

=== VALIDATION RESULTS ===
Missing roles:
- role1
- role2
- role3
```

## How It Works

1. Reads the specified playbook file to extract directly referenced roles
2. Scans the roles directory to identify available roles (including symlinks)
3. For each role, reads its meta/main.yml file to identify dependencies
4. Recursively resolves all dependencies until the complete set of required roles is identified
5. Compares the required roles with available roles to identify missing ones

## License

This project is licensed under the [MIT License](./LICENSE).
