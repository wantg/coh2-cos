// Modules to control application life and create native browser window
const { app, BrowserWindow, Menu, shell } = require('electron');
const { exec, spawn } = require('child_process');
const path = require('path');
const fs = require('fs');

const config = {
  url: 'http://127.0.0.1:25416/match',
  width: 1280,
  height: 1080,
};

Menu.setApplicationMenu(null);
app.commandLine.appendSwitch('ignore-certificate-errors', 'true');

let mainWindow = null;
const gotTheLock = app.requestSingleInstanceLock();
let tunnel = null;

function launchTunnel(tunnelPath) {
  tunnel = spawn(tunnelPath, ['-c', path.join(process.env.PORTABLE_EXECUTABLE_DIR, 'config.yml')]);
  tunnel.stdout.on('data', (data) => {
    console.log(`stdout: ${data}`);
  });
  tunnel.stderr.on('data', (data) => {
    console.error(`stderr: ${data}`);
  });
  tunnel.on('close', (code) => {
    console.log(`child process exited with code ${code}`);
  });
}

function init() {
  const tunnelPathList = ['dist/hq/hq.exe', '../../MacOS/hq', '../../hq/hq.exe'];
  for (const _p of tunnelPathList) {
    const p = path.join(__dirname, _p);
    console.log(p);
    if (fs.existsSync(p)) {
      if (fs.lstatSync(p).isFile()) {
        tunnelPath = p;
        console.log('tunnelPath is', tunnelPath);
        break;
      }
    }
  }
  if (!tunnelPath) {
    app.quit();
    return;
  }
  launchTunnel(tunnelPath);
}

function killTunnel() {
  if(tunnel){
    tunnel.kill(1);
  }
}

function createWindow() {
  // Create the browser window.
  mainWindow = new BrowserWindow({
    width: config.width,
    height: config.height,
    webPreferences: {
      preload: path.join(__dirname, 'preload.js'),
      devTools: false,
    },
  });

  // and load the index.html of the app.
  // mainWindow.loadFile('index.html')
  mainWindow.loadURL(config.url);

  mainWindow.webContents.on('new-window', function (e, url) {
    // mainWindow.loadURL(url);
    shell.openExternal(url);
    e.preventDefault();
  });

  // Open the DevTools.
  // mainWindow.webContents.openDevTools()
}

if (!gotTheLock) {
  killTunnel();
  app.quit();
  return;
}

app.on('second-instance', (event, commandLine, workingDirectory) => {
  // Someone tried to run a second instance, we should focus our window.
  if (mainWindow) {
    if (mainWindow.isMinimized()) mainWindow.restore();
    mainWindow.focus();
  }
});

// This method will be called when Electron has finished
// initialization and is ready to create browser windows.
// Some APIs can only be used after this event occurs.
app.whenReady().then(() => {
  init();
  createWindow();

  app.on('activate', function () {
    // On macOS it's common to re-create a window in the app when the
    // dock icon is clicked and there are no other windows open.
    if (BrowserWindow.getAllWindows().length === 0) createWindow();
  });
});

// Quit when all windows are closed, except on macOS. There, it's common
// for applications and their menu bar to stay active until the user quits
// explicitly with Cmd + Q.
app.on('window-all-closed', function () {
  if (true || process.platform !== 'darwin') {
    app.quit();
  }
});

app.on('before-quit', function () {
  killTunnel();
  console.log('before-quit');
});

app.on('will-quit', function () {
  killTunnel();
  console.log('will-quit');
});

app.on('quit', function () {
  killTunnel();
  console.log('quit');
});

// In this file you can include the rest of your app's specific main process
// code. You can also put them in separate files and require them here.
