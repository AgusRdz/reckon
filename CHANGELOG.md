# Changelog

All notable changes to reckon are documented here.

## [0.3.3] - 2026-04-01

### Bug Fixes
- Prepend reckon to SessionStart so it runs before ctx restore
([187f73c](https://github.com/AgusRdz/reckon/commit/187f73c742c207c4f932d5b07469c8f338be516c))
## [0.3.2] - 2026-04-01

### Bug Fixes
- Use quoted forward-slash path format for hook command registration
([7367fff](https://github.com/AgusRdz/reckon/commit/7367fffefe34b6f05f28f8c7a94858ef7512094a))
## [0.3.1] - 2026-04-01

### Bug Fixes
- Use dedicated 'hook' subcommand for SessionStart, drop TTY detection
([de4e07a](https://github.com/AgusRdz/reckon/commit/de4e07a622a1588676405460cf563539dbd827c9))
## [0.3.0] - 2026-04-01

### Features
- Add reckon update command
([d5ad574](https://github.com/AgusRdz/reckon/commit/d5ad5742935e816a702941aa8a69edb711845491))
## [0.2.0] - 2026-04-01

### Bug Fixes
- Detect hook mode via stdout instead of stdin
([36a7b98](https://github.com/AgusRdz/reckon/commit/36a7b986ff523443d9cd9e6009b5a679c10a60d0))

### Documentation
- Add PowerShell install, auto-register hook in all install steps
([e92920c](https://github.com/AgusRdz/reckon/commit/e92920c9574d03209b141406820498c3efe5dd5c))

### Features
- Add install.sh and install.ps1 with auto hook registration
([711a528](https://github.com/AgusRdz/reckon/commit/711a52842e5a5eacde4db563ed3375dbdce4d1fd))
- Auto-add .codeindex to .gitignore on index build
([cf003e0](https://github.com/AgusRdz/reckon/commit/cf003e090208d8549c5726838b35fba6f5a2cc1d))

### Miscellaneous
- Add release signing public key
([a2ad12c](https://github.com/AgusRdz/reckon/commit/a2ad12ce7ba973a505d1d3d19e402ddfd9b46c23))
- Update release signing public key
([873f615](https://github.com/AgusRdz/reckon/commit/873f615eeb0b67726aa5915099905ee1e25ca3aa))
- Rename public key, ignore signing.pem
([affdbc3](https://github.com/AgusRdz/reckon/commit/affdbc33d8f5dc884b096be0759d1663173676ea))
- Install git-cliff in Docker, use it for release targets
([0c9f401](https://github.com/AgusRdz/reckon/commit/0c9f4019ecbee5a5d88a39e13950b97b258b69fa))
## [0.1.0] - 2026-04-01

### Documentation
- Add README with install and hook registration instructions
([58d3677](https://github.com/AgusRdz/reckon/commit/58d3677f8b1a5cc00315025b62593b975b29b379))

### Features
- Initial implementation of reckon symbol indexer
([8fdeaee](https://github.com/AgusRdz/reckon/commit/8fdeaee6e751d5768beb1c383002fc372da590e1))

### Miscellaneous
- Untrack CLAUDE.md and PLAN.md, add to .gitignore
([a6d2f56](https://github.com/AgusRdz/reckon/commit/a6d2f56ddaec723a4c819fef46fbd1a997492847))
- Remove unused check package
([9170425](https://github.com/AgusRdz/reckon/commit/9170425da08fdf4e6cc315b1027313704802a856))

### Testing
- Add unit tests for all packages
([1e746d4](https://github.com/AgusRdz/reckon/commit/1e746d46fcb1a3a74fec6fd40e7220da27cc7565))

