# Build Instructions

## Development Build
```bash
wails build
```

## Production Build with Installer
```bash
wails build --target windows/amd64 --nsis
```

## Directory Structure After Installation
```
C:\Program Files\YourCompany\Monitor Profile Manager\
├── monitor-profile-manager-wails.exe
├── tools\
│   ├── svcl\
│   │   ├── svcl.exe
│   │   └── audio.csv
│   └── multimonitortool\
│       ├── MultiMonitorTool.exe
│       └── MultiMonitorTool.cfg
└── [other app files]
```

## Path Resolution
The application now uses dynamic path resolution:
- **Development**: Uses relative paths from project root (`tools/svcl/svcl.exe`)
- **Production**: Uses paths relative to executable directory (`exeDir/tools/svcl/svcl.exe`)

This ensures the tools are properly bundled with the installer and accessible at runtime.
