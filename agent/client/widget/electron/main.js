import { app, BrowserWindow, Tray, Menu, nativeImage } from "electron";
import path from "path"; // distinct import
import { fileURLToPath } from "url";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);


const iconPath = path.join(__dirname, "icon.ico");

let tray = null;
let win = null;

function createWindow() {
  win = new BrowserWindow({
    width: 420,
    height: 300,
    resizable: false,
    frame: false,
    title: "NetShield",
    icon: iconPath, // using the variable defined at top
    autoHideMenuBar: true,
    alwaysOnTop: true,
    webPreferences: {
      nodeIntegration: false,
      contextIsolation: true,
      webSecurity: false,
    },
  });


  const isDev = !app.isPackaged;

  if (isDev) {

    console.log("Running in DEV mode: Loading localhost");
    win.loadURL("http://localhost:3000");
  } else {
  
    const indexPath = path.join(__dirname, "../out/index.html");

    console.log("Running in PROD mode: Loading file from", indexPath);
    win.loadFile(indexPath);
  }



  win.once("ready-to-show", () => {
    win.show();
  });
}

function createTray() {

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
