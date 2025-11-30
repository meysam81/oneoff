import { createGunzip } from "node:zlib";
import { pipeline } from "node:stream/promises";
import { Writable } from "node:stream";
import log from "./logger";

const GITHUB_REPO = "meysam81/oneoff";
const GITHUB_API_URL = `https://api.github.com/repos/${GITHUB_REPO}/releases/latest`;

// Fallback values if GitHub API fails
const FALLBACK_VERSION = "v1.0";
const FALLBACK_BINARY_SIZE = "~7MB";

interface GitHubRelease {
  tag_name: string;
  assets: Array<{
    name: string;
    browser_download_url: string;
    size: number;
  }>;
}

interface GitHubData {
  version: string; // e.g., "v1.0"
  fullVersion: string; // e.g., "v1.0.2"
  binarySize: string; // e.g., "~15MB"
  downloadUrl: string; // URL for linux_amd64 tar.gz
}

/**
 * Extract major.minor version from full semver tag
 * e.g., "v1.0.2" -> "v1.0"
 */
function extractMajorMinor(tag: string): string {
  const match = tag.match(/^v?(\d+)\.(\d+)/);
  if (match) {
    return `v${match[1]}.${match[2]}`;
  }
  return FALLBACK_VERSION;
}

/**
 * Format bytes to human readable size
 */
function formatBytes(bytes: number): string {
  const mb = bytes / (1024 * 1024);
  if (mb >= 1) {
    return `~${Math.round(mb)}MB`;
  }
  const kb = bytes / 1024;
  return `~${Math.round(kb)}KB`;
}

/**
 * Parse tar header to extract file size
 * TAR header format: size is at offset 124, 12 bytes, octal string
 */
function parseTarHeaderSize(header: Buffer): number {
  // File size is at offset 124, 12 bytes in octal
  const sizeField = header.subarray(124, 136).toString("utf8").trim();
  // Remove null bytes and parse as octal
  const cleanSize = sizeField.replace(/\0/g, "").trim();
  return parseInt(cleanSize, 8) || 0;
}

/**
 * Download tar.gz and extract the binary size without fully extracting
 */
async function getBinarySize(
  downloadUrl: string,
  token?: string,
): Promise<number> {
  const headers: Record<string, string> = {
    Accept: "application/octet-stream",
    "User-Agent": "OneOff-Landing-Page",
  };

  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }

  const response = await fetch(downloadUrl, {
    headers,
    redirect: "follow",
  });

  if (!response.ok) {
    throw new Error(`Failed to download: ${response.status}`);
  }

  const arrayBuffer = await response.arrayBuffer();
  const compressedBuffer = Buffer.from(arrayBuffer);

  // Decompress gzip
  const chunks: Buffer[] = [];
  const gunzip = createGunzip();

  const collectChunks = new Writable({
    write(chunk, _encoding, callback) {
      chunks.push(chunk);
      callback();
    },
  });

  // Create a readable stream from buffer
  const { Readable } = await import("node:stream");
  const readable = Readable.from(compressedBuffer);

  await pipeline(readable, gunzip, collectChunks);

  const tarBuffer = Buffer.concat(chunks);

  // Parse tar to find the main binary (largest file, or file named "oneoff")
  let maxSize = 0;
  let offset = 0;

  while (offset < tarBuffer.length) {
    const header = tarBuffer.subarray(offset, offset + 512);

    // Check for end of archive (two zero blocks)
    if (header.every((b) => b === 0)) {
      break;
    }

    // Get filename (first 100 bytes)
    const filename = header.subarray(0, 100).toString("utf8").split("\0")[0];

    // Get file size
    const fileSize = parseTarHeaderSize(header);

    // Check if this is the oneoff binary
    if (filename === "oneoff" || filename.endsWith("/oneoff")) {
      return fileSize;
    }

    // Track largest file as fallback
    if (fileSize > maxSize) {
      maxSize = fileSize;
    }

    // Move to next header (512 byte header + file content padded to 512)
    const contentBlocks = Math.ceil(fileSize / 512);
    offset += 512 + contentBlocks * 512;
  }

  return maxSize;
}

/**
 * Fetch GitHub release data at build time
 * Uses GITHUB_TOKEN env var for authentication if available
 */
export async function fetchGitHubData(): Promise<GitHubData> {
  const token = process.env.GITHUB_TOKEN;

  try {
    const headers: Record<string, string> = {
      Accept: "application/vnd.github.v3+json",
      "User-Agent": "OneOff-Landing-Page",
    };

    if (token) {
      headers["Authorization"] = `Bearer ${token}`;
    }

    log.info("Fetching GitHub release data...");

    const response = await fetch(GITHUB_API_URL, { headers });

    if (!response.ok) {
      throw new Error(`GitHub API returned ${response.status}`);
    }

    const release: GitHubRelease = await response.json();
    const fullVersion = release.tag_name;
    const version = extractMajorMinor(fullVersion);

    log.info(`Found version: ${fullVersion} -> ${version}`);

    // Find linux_amd64 tar.gz asset
    const linuxAsset = release.assets.find(
      (a) => a.name === "oneoff_linux_amd64.tar.gz",
    );

    if (!linuxAsset) {
      throw new Error("Linux amd64 asset not found");
    }

    const downloadUrl = linuxAsset.browser_download_url;

    // Get actual binary size by downloading and decompressing
    log.info("Downloading and measuring binary size...");
    const binaryBytes = await getBinarySize(downloadUrl, token);
    const binarySize = formatBytes(binaryBytes);

    log.info(`Binary size: ${binarySize} (${binaryBytes} bytes)`);

    return {
      version,
      fullVersion,
      binarySize,
      downloadUrl,
    };
  } catch (error) {
    log.warn(`Failed to fetch GitHub data: ${error}`);
    log.warn("Using fallback values");

    return {
      version: FALLBACK_VERSION,
      fullVersion: FALLBACK_VERSION,
      binarySize: FALLBACK_BINARY_SIZE,
      downloadUrl: `https://github.com/${GITHUB_REPO}/releases/latest/download/oneoff_linux_amd64.tar.gz`,
    };
  }
}

// Cache the result to avoid multiple API calls during build
let cachedData: GitHubData | null = null;

export async function getGitHubData(): Promise<GitHubData> {
  if (!cachedData) {
    cachedData = await fetchGitHubData();
  }
  return cachedData;
}
