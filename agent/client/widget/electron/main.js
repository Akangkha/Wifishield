import { app, BrowserWindow, Tray, Menu, nativeImage } from "electron";
import { join } from "path";

  const iconPath = join(__dirname, "icon.png");

let tray = null;
let win = null;

function createWindow() {
  win = new BrowserWindow({
    width: 420,
    height: 300,
    resizable: false,
    frame: false,        
    title: "NetShield",
    icon: iconPath,
    autoHideMenuBar: true,
    alwaysOnTop: true,  
    webPreferences: {
      nodeIntegration: false,
      contextIsolation: true,
    },
  });

  const url ="http://localhost:3000"
 

  win.loadURL(url);
  win.once("ready-to-show", () => {
    win.show();
  });
}

function createTray() {
  const iconPath = join(__dirname, "icon.png");
  const trayIcon = nativeImage.createFromPath(iconPath);

  tray = new Tray(trayIcon);
  tray.setToolTip("NetShield – Wi-Fi Guardian");

  const contextMenu = Menu.buildFromTemplate([
    {
      label: "Show",
      click: () => {
        if (!win) return;
        win.show();
        win.focus();
      },
    },
    {
      label: "Hide",
      click: () => {
        if (!win) return;
        win.hide();
      },
    },
    { type: "separator" },
    {
      label: "Quit",
      click: () => {
        app.isQuitting = true;
        app.quit();
      },
    },
  ]);

  tray.setContextMenu(contextMenu);


  tray.on("click", () => {
    if (!win) return;
    if (win.isVisible()) {
      win.hide();
    } else {
      win.show();
      win.focus();
    }
  });
}

app.whenReady().then(() => {
  console.log("[electron] app ready, creating window + tray");
  createWindow();
  createTray();

  app.on("activate", () => {
    if (BrowserWindow.getAllWindows().length === 0) {
      createWindow();
    }
  });
});

app.on("window-all-closed", () => {
  if (process.platform !== "darwin") {
    app.quit();
  }
});
