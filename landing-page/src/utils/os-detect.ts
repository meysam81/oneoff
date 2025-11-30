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
}

export var platformList: PlatformInfo[] = [
  { id: "linux-amd64", name: "Linux", arch: "x86_64" },
  { id: "linux-arm64", name: "Linux", arch: "ARM64" },
  { id: "darwin-arm64", name: "macOS", arch: "Apple Silicon" },
  { id: "darwin-amd64", name: "macOS", arch: "Intel" },
  { id: "windows-amd64", name: "Windows", arch: "x86_64" },
];

export var allPlatforms: Platform[] = [
  "linux-amd64",
  "linux-arm64",
  "darwin-amd64",
  "darwin-arm64",
  "windows-amd64",
];

function getDownloadUrl(platform: Platform, version: string): string {
  var baseUrl = "https://github.com/meysam81/oneoff/releases/download";
  var files: Record<Platform, string> = {
    "linux-amd64": "oneoff_linux_amd64.tar.gz",
    "linux-arm64": "oneoff_linux_arm64.tar.gz",
    "darwin-amd64": "oneoff_darwin_amd64.tar.gz",
    "darwin-arm64": "oneoff_darwin_arm64.tar.gz",
    "windows-amd64": "oneoff_windows_amd64.zip",
  };
  return baseUrl + "/" + version + "/" + files[platform];
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
      '" `\n  -OutFile "oneoff.zip"\nExpand-Archive -Path "oneoff.zip" -DestinationPath "."'
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
      '" `\n  -OutFile "oneoff.zip"\nExpand-Archive -Path "oneoff.zip" -DestinationPath "."\n.\\oneoff.exe'
    );
  }
  return "curl -fsSL " + url + " | tar xz\n./oneoff";
}

export function getRunCommand(platform: Platform): string {
  return platform === "windows-amd64" ? ".\\oneoff.exe" : "./oneoff";
}

export function getOpenCommand(platform: Platform): string {
  if (platform === "windows-amd64")
    return 'Start-Process "http://localhost:8080"';
  if (platform.startsWith("darwin")) return "open http://localhost:8080";
  return "xdg-open http://localhost:8080";
}

export function getTerminalTitle(platform: Platform): string {
  return platform === "windows-amd64" ? "PowerShell" : "Terminal";
}

export function getShortInstallCommand(platform: Platform): string {
  if (platform === "windows-amd64")
    return "Invoke-WebRequest ... | Expand-Archive";
  return "curl -fsSL github.com/.../oneoff.tar.gz | tar xz";
}

export function getPromptChar(platform: Platform): string {
  return platform === "windows-amd64" ? ">" : "$";
}

export function getPlatformDisplayName(platform: Platform): string {
  var names: Record<Platform, string> = {
    "linux-amd64": "Linux",
    "linux-arm64": "Linux ARM",
    "darwin-amd64": "macOS",
    "darwin-arm64": "macOS",
    "windows-amd64": "Windows",
  };
  return names[platform] || "Linux";
}

function buildPlatformMap<T>(fn: (p: Platform) => T): Record<Platform, T> {
  var result: Record<string, T> = {};
  allPlatforms.forEach(function (p) {
    result[p] = fn(p);
  });
  return result as Record<Platform, T>;
}

export function getAllDownloadCommands(
  version: string,
): Record<Platform, string> {
  return buildPlatformMap(function (p) {
    return getDownloadCommand(p, version);
  });
}

export function getAllInstallCommands(
  version: string,
): Record<Platform, string> {
  return buildPlatformMap(function (p) {
    return getInstallCommand(p, version);
  });
}

export function getAllRunCommands(): Record<Platform, string> {
  return buildPlatformMap(getRunCommand);
}

export function getAllOpenCommands(): Record<Platform, string> {
  return buildPlatformMap(getOpenCommand);
}

export function getAllTerminalTitles(): Record<Platform, string> {
  return buildPlatformMap(getTerminalTitle);
}

export function getAllShortCommands(): Record<Platform, string> {
  return buildPlatformMap(getShortInstallCommand);
}

export function getAllPromptChars(): Record<Platform, string> {
  return buildPlatformMap(getPromptChar);
}

export function getAllDisplayNames(): Record<Platform, string> {
  return buildPlatformMap(getPlatformDisplayName);
}
