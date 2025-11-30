import { ref, computed, onMounted } from "vue";

export type Platform =
  | "linux-amd64"
  | "linux-arm64"
  | "darwin-amd64"
  | "darwin-arm64"
  | "windows-amd64";

var platformDisplayNames: Record<Platform, string> = {
  "linux-amd64": "Linux",
  "linux-arm64": "Linux ARM",
  "darwin-amd64": "macOS",
  "darwin-arm64": "macOS",
  "windows-amd64": "Windows",
};

export function usePlatform() {
  var platform = ref<Platform>("linux-amd64");

  function detectPlatform(): Platform {
    if (typeof navigator === "undefined") return "linux-amd64";
    var userAgent = navigator.userAgent.toLowerCase();
    var p = (navigator.platform || "").toLowerCase();
    var isWindows = p.includes("win") || userAgent.includes("windows");
    var isMac = p.includes("mac") || userAgent.includes("macintosh");
    var isArm =
      userAgent.includes("arm") ||
      userAgent.includes("aarch64") ||
      ((navigator as any).userAgentData &&
        (navigator as any).userAgentData.architecture === "arm");
    if (isWindows) return "windows-amd64";
    if (isMac) return isArm ? "darwin-arm64" : "darwin-amd64";
    if (isArm) return "linux-arm64";
    return "linux-amd64";
  }

  var displayName = computed(function () {
    return platformDisplayNames[platform.value] || "Linux";
  });

  var promptChar = computed(function () {
    return platform.value === "windows-amd64" ? ">" : "$";
  });

  var terminalTitle = computed(function () {
    return platform.value === "windows-amd64" ? "PowerShell" : "Terminal";
  });

  var runCommand = computed(function () {
    return platform.value === "windows-amd64" ? ".\\oneoff.exe" : "./oneoff";
  });

  var openCommand = computed(function () {
    if (platform.value === "windows-amd64")
      return 'Start-Process "http://localhost:8080"';
    if (platform.value.startsWith("darwin"))
      return "open http://localhost:8080";
    return "xdg-open http://localhost:8080";
  });

  onMounted(function () {
    platform.value = detectPlatform();
  });

  return {
    platform,
    displayName,
    promptChar,
    terminalTitle,
    runCommand,
    openCommand,
  };
}
