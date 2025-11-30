export type Platform =
  | "linux-amd64"
  | "linux-arm64"
  | "darwin-amd64"
  | "darwin-arm64"
  | "windows-amd64";

export interface PlatformInfo {
  id: Platform;
  name: string;
  arch: string;
  extension: string;
  extractCommand: string;
}

export var platforms: Record<Platform, PlatformInfo> = {
  "linux-amd64": {
    id: "linux-amd64",
    name: "Linux",
    arch: "x86_64",
    extension: "tar.gz",
    extractCommand: "tar xz",
  },
  "linux-arm64": {
    id: "linux-arm64",
    name: "Linux",
    arch: "ARM64",
    extension: "tar.gz",
    extractCommand: "tar xz",
  },
  "darwin-amd64": {
    id: "darwin-amd64",
    name: "macOS",
    arch: "Intel",
    extension: "tar.gz",
    extractCommand: "tar xz",
  },
  "darwin-arm64": {
    id: "darwin-arm64",
    name: "macOS",
    arch: "Apple Silicon",
    extension: "tar.gz",
    extractCommand: "tar xz",
  },
  "windows-amd64": {
    id: "windows-amd64",
    name: "Windows",
    arch: "x86_64",
    extension: "zip",
    extractCommand: "Expand-Archive",
  },
};

export var platformList: PlatformInfo[] = [
  platforms["linux-amd64"],
  platforms["linux-arm64"],
  platforms["darwin-arm64"],
  platforms["darwin-amd64"],
  platforms["windows-amd64"],
];

export function detectPlatform(): Platform {
  if (typeof navigator === "undefined") {
    return "linux-amd64";
  }

  var userAgent = navigator.userAgent.toLowerCase();
  var platform = (navigator.platform || "").toLowerCase();

  var isWindows = platform.includes("win") || userAgent.includes("windows");
  var isMac = platform.includes("mac") || userAgent.includes("macintosh");
  var isLinux = platform.includes("linux") || userAgent.includes("linux");

  var isArm =
    userAgent.includes("arm") ||
    userAgent.includes("aarch64") ||
    (navigator as any).userAgentData?.architecture === "arm";

  if (isWindows) {
    return "windows-amd64";
  }

  if (isMac) {
    return isArm ? "darwin-arm64" : "darwin-amd64";
  }

  if (isLinux) {
    return isArm ? "linux-arm64" : "linux-amd64";
  }

  return "linux-amd64";
}

export function getDownloadUrl(platform: Platform, version: string): string {
  var baseUrl = "https://github.com/meysam81/oneoff/releases/download";
  var platformMap: Record<Platform, string> = {
    "linux-amd64": "oneoff_linux_amd64.tar.gz",
    "linux-arm64": "oneoff_linux_arm64.tar.gz",
    "darwin-amd64": "oneoff_darwin_amd64.tar.gz",
    "darwin-arm64": "oneoff_darwin_arm64.tar.gz",
    "windows-amd64": "oneoff_windows_amd64.zip",
  };
  return baseUrl + "/" + version + "/" + platformMap[platform];
}

export function getDownloadCommand(
  platform: Platform,
  version: string,
): string {
  var url = getDownloadUrl(platform, version);

  if (platform === "windows-amd64") {
    return (
      'Invoke-WebRequest -Uri "' +
      url +
      '" -OutFile "oneoff.zip"; Expand-Archive -Path "oneoff.zip" -DestinationPath "."'
    );
  }

  return "curl -fsSL " + url + " | tar xz";
}

export function getInstallCommand(platform: Platform, version: string): string {
  var url = getDownloadUrl(platform, version);

  if (platform === "windows-amd64") {
    return (
      'Invoke-WebRequest -Uri "' +
      url +
      '" -OutFile "oneoff.zip"; Expand-Archive -Path "oneoff.zip" -DestinationPath "."; .\\oneoff.exe'
    );
  }

  return "curl -fsSL " + url + " | tar xz && ./oneoff";
}

export function getRunCommand(platform: Platform): string {
  if (platform === "windows-amd64") {
    return ".\\oneoff.exe";
  }
  return "./oneoff";
}

export function getOpenCommand(platform: Platform): string {
  if (platform === "windows-amd64") {
    return 'Start-Process "http://localhost:8080"';
  }
  if (platform.startsWith("darwin")) {
    return "open http://localhost:8080";
  }
  return "xdg-open http://localhost:8080";
}

export function getTerminalTitle(platform: Platform): string {
  if (platform === "windows-amd64") {
    return "PowerShell";
  }
  return "Terminal";
}

export function getShortInstallCommand(platform: Platform): string {
  if (platform === "windows-amd64") {
    return "Invoke-WebRequest ... | Expand-Archive";
  }
  return "curl -fsSL github.com/.../oneoff.tar.gz | tar xz";
}

export function generatePlatformCommandsJson(version: string): string {
  var commands: Record<Platform, string> = {
    "linux-amd64": getInstallCommand("linux-amd64", version),
    "linux-arm64": getInstallCommand("linux-arm64", version),
    "darwin-amd64": getInstallCommand("darwin-amd64", version),
    "darwin-arm64": getInstallCommand("darwin-arm64", version),
    "windows-amd64": getInstallCommand("windows-amd64", version),
  };
  return JSON.stringify(commands);
}
