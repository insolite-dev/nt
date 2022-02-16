<p align="center">
<img src="https://user-images.githubusercontent.com/59066341/136813272-a90861f4-1e6c-4a83-9a3b-18e01b99de34.png" width="400px">
</p>

---

 <a href="https://github.com/anonistas/notya/blob/main/LICENSE">
  <img src="https://img.shields.io/badge/License-Apache-red.svg" alt="License: MIT"/>
 </a>
 <a href="https://discord.gg/CtStkzrHV3">
   <img src="https://img.shields.io/discord/914899238415130714?color=blue&label=Anon Community&logo=discord" alt="Anoncord" />
 </a>   


# Installation
See the [last release](https://github.com/anonistas/notya/releases/latest), where you can find binary files for your ecosystem

### Brew:
```
brew install --build-from-source notya
```

### Curl:
```
curl -sfL https://raw.githubusercontent.com/anonistas/notya/main/install.sh | sh
```

# Usage 
**Note**: _notya only available for local service (local machine database) for now._

### Help:
Run `notya help` or `notya -h` to see default [help.txt](https://github.com/anonistas/notya/wiki/help.txt). <br>
 
Use `notya [command] --help` for more information about a command.

### Init: 
Use `notya init` to initialize application. <br/>
*(It isn't must actually, whenever you call any command of notya, it checks initialization status and if it isn't initialized, initializes app automatically).*

### Settings (config):
The config file of notya is autogeneratable, it'd be genereted by `Init` functionality. <br>
**Refer to settings documentation for details - [Settings Wiki](https://github.com/anonistas/notya/wiki/Settings)**

### Available commands:
- **[See all notes](https://github.com/anonistas/notya/wiki/List)** - `notya list`
- **[View note](https://github.com/anonistas/notya/wiki/View)** - `notya view` or `notya view [name]`
- **[Create note](https://github.com/anonistas/notya/wiki/Create)** - `notya create` or `notya create [title]`
- **[Rename note](https://github.com/anonistas/notya/wiki/Rename)** - `notya rename` or `notya rename [name]`
- **[Edit note](https://github.com/anonistas/notya/wiki/Edit)** - `notya edit` or `notya edit [name]`
- **[Remove note](https://github.com/anonistas/notya/wiki/Remove)** - `notya remove` or `notya rm [name]`

# Contributing
For information regarding contributions, please refer to [CONTRIBUTING.md](https://github.com/anonistas/notya/blob/develop/CONTRIBUTING.md) file.
