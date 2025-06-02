<div align="center">
    <img width="120" height="120" hspace="10" alt="ZDS Logo" src="https://codeberg.org/Zusty/ZDS/media/branch/main/img/zds.png"/>
    <h1>GCT</h1>
    A Git Commit Tool.
<br/>
<a href="https://codeberg.org/Zillowe/ZFPL">
<img alt="ZFPL-1.0" src="https://codeberg.org/Zillowe/ZFPL/raw/branch/main/badges/1-0/dark.svg"/>
</a>
</div>

<hr/>

## Installation

You can either build it from source or install it using installer scripts

### Scripts

To install GCT, you need to run this command:

```sh
# For Linux/macOS
curl -fsSL https://zusty.codeberg.page/GCT/@main/app/install.sh | bash
# For Windows
powershell -c "irm https://zusty.codeberg.page/GCT/@main/app/install.ps1|iex"
```

### Build

To build GCT from source you need to have [`go`](https://go.dev) installed.

then run this command to build it:

```sh
# For Linux/macOS
./build/build.sh
# For Windows
./build/build.ps1
```

To build it for all run this:

```sh
# For Linux/macOS
./build/build-all.sh
# For Windows
./build/build-all.ps1
```

## Usage

GCT is a command-line tool. Here are the available commands:

```sh
gct <command> [arguments...]
```

**Available Commands:**

* `gct version`: Show GCT version information.
* `gct about`: Display details and information about GCT.
* `gct update`: Check for and apply updates to GCT itself.
* `gct commit`: Create a new git commit interactively.
* `gct help`: Show the help message.

**Flags:**

* `gct -v`, `gct --version`: Show GCT version information.

For more details on a specific command, run `gct <command> --help`.

## Documentation

To get started with GCT please refer to the [GCT Wiki](https://codeberg.org/Zusty/GCT/wiki).

## Footer

GCT is developed by Zusty < Zillowe Foundation, part of the Zillowe Development Suite (ZDS)

### License

GCT is licensed under the [ZFPL](https://codeberg.org/Zillowe/ZFPL)-1.0 (Zillowe Foundation Public License, Version 1.0).
