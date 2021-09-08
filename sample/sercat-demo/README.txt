
## Configuration file

Upon startup, the `sercat` command looks for a `.sercat` file in the current or
any readable parent directory. It loads all of them and merges all the top level
sections, giving override priority to the file in the current directory or the
closest parent.

The configuration file is using the YAML format. Here is an example of
configuration.

```YAML
default:
  port: /dev/tty.usbserial-*
  baudrate: 230400
  format: 8N1

/dev/tty.usbserial-*:
  baudrate: 230400
  format: 8N1
```

Each top level key is either:
- the name of a serial port
- the partial name of a group of serial port, suffixed with a `*`.
- an alias name for a specific port
- the special alias `default`, used if no port is specified on the commandline.

Each configuration section supports 3 options
- `port`: the name of the port to use, which must be specified only when
  defining the configuration for an aliases.
- `baudrate`: the default baudrate to use unless specified on the command-line.
  If the baudrate is not defined at all, neither in the active configuration nor
  on the command-line, it defaults to 115200.
- `format`: the default format to use unless specified on the command-line. If
  the format is not defined at all, neither in the active configuration nor on
  the command-line, it defaults to 8N1.
