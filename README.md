# Headless Askpass

An askpass helper that lets you enter password in another terminal. Useful for automated sudo in headless sessions.

## Rationale

You have a program, script, or AI agent that wants to run some command over SSH, it requires sudo privilege. But you cannot enter password because there is no terminal to interact with.

For example:

```
$ ssh remote sudo tcpdump -U -i enp5s0 -w - 'not port 22' | wireshark -i -
sudo: a terminal is required to read the password; either use ssh's -t option or configure an askpass helper
sudo: a password is required
```

Existing askpass helpers mostly require a trusted graphical session to obtain passwords from. This program lets you use another terminal session instead.

## Usage

In a terminal, start a server by `headless-askpass -s`. On the same machine and user, subsequent calls to `sudo` without a terminal or `sudo -A` will prompt for password in that terminal. End the server with Ctrl+C when you are done.

## Installation

1. Download a binary from [releases](https://github.com/FiveYellowMice/headless-askpass/releases) or build one from source.
2. Put the binary to somewhere in your `$PATH`, like `/usr/local/bin/headless-askpass`. Make it executable.
3. Configure sudo to use this askpass helper. Choose one of:
    - Add to `/etc/sudoers` or a file in `/etc/suders.d`:
        ```
        Path askpass /usr/local/bin/headless-askpass
        ```
    - If `PermitUserEnvironment` is turned on in SSH, add to `~/.ssh/environment`:
        ```
        SUDO_ASKPASS=/usr/local/bin/headless-askpass
        ```
    - Add to `/etc/environment`:
        ```
        SUDO_ASKPASS=/usr/local/bin/headless-askpass
        ```

## Building

1. Install [Go](https://go.dev/).
2. `make`

## Similar Projects

- [SSH Askpass Helper](https://github.com/GlassOnTin/secure-askpass/): Use SSH key encryption to provide sudo passwords.
