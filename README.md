# ansible-roledep-validator

`ansible-roledep-validator` is a command-line tool that helps identify missing roles in your Ansible project. It analyzes playbooks, detects required roles (including dependencies defined in meta/main.yml files), and compares them with the roles available in your roles directory.

## Usage

```bash
$ ansible-roledep-validator -playbook playbook.yml -roles /path/to/roles
```

### Flags

- `-playbook`: Path to the Ansible playbook YAML file (required)
- `-roles`: Path to the Ansible roles directory (default: `roles`)

## Output

Missing role names go to stdout, one per line, sorted. No output means every required role is present.

```
role1
role2
role3
```

### Exit status

- `0`: every required role is present
- `1`: missing roles were printed to stdout
- `2`: an error occurred

## How roles are resolved

A role is a directory, or a symlink to one, directly under the roles directory. Both a play's `roles:` and a role's `dependencies:` name it in one of two forms:

```yaml
roles:
  - role1
  - role: role2
    tags: tag1
```

Any other entry is an error rather than being skipped: a mapping without a `role` key, a path, or a name containing `{{`. Skipping one would report that nothing is missing while a required role is absent.

## Limitations

Roles reached through `include_role` or `import_role` are not analyzed, because they are named inside tasks and only play-level `roles:` and `meta/main.yml` are read. A role required only that way is therefore not reported when it is missing.

`import_playbook` is not followed. Pass the imported playbook itself to check the roles it declares.

Only `meta/main.yml` is read for dependencies; `meta/main.yaml` is not.

## License

This project is licensed under the [MIT License](./LICENSE).
