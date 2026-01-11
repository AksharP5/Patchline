const crypto = require("node:crypto");
const fs = require("node:fs");
const path = require("node:path");

const platformMap = {
  darwin: "darwin",
  linux: "linux",
  win32: "windows",
};

const archMap = {
  x64: "amd64",
  arm64: "arm64",
};

function resolvePlatform(platform) {
  const resolved = platformMap[platform];
  if (!resolved) {
    throw new Error(`unsupported platform: ${platform}`);
  }
  return resolved;
}

function resolveArch(arch) {
  const resolved = archMap[arch];
  if (!resolved) {
    throw new Error(`unsupported architecture: ${arch}`);
  }
  return resolved;
}

function getArchiveExtension(goos) {
  if (goos === "windows") {
    return "zip";
  }
  return "tar.gz";
}

function getArchiveName(version, goos, goarch) {
  const ext = getArchiveExtension(goos);
  return `patchline_${version}_${goos}_${goarch}.${ext}`;
}

function getDownloadUrl(version, goos, goarch) {
  const archiveName = getArchiveName(version, goos, goarch);
  return `https://github.com/AksharP5/Patchline/releases/download/v${version}/${archiveName}`;
}

function getChecksumsUrl(version) {
  return `https://github.com/AksharP5/Patchline/releases/download/v${version}/checksums.txt`;
}

function getBinaryName(goos) {
  if (goos === "windows") {
    return "patchline.exe";
  }
  return "patchline";
}

function getInstalledBinaryName(goos) {
  if (goos === "windows") {
    return "patchline-bin.exe";
  }
  return "patchline-bin";
}

function getCandidateBinaryPaths(rootDir, version, goos, goarch) {
  const binaryName = getBinaryName(goos);
  return [
    path.join(rootDir, binaryName),
    path.join(rootDir, `patchline_${version}_${goos}_${goarch}`, binaryName),
  ];
}

function parseChecksums(text) {
  const result = {};
  for (const line of text.split(/\r?\n/)) {
    const trimmed = line.trim();
    if (!trimmed) {
      continue;
    }
    const parts = trimmed.split(/\s+/);
    if (parts.length < 2) {
      continue;
    }
    const hash = parts[0];
    const filename = parts[1].replace(/^\*?\.\//, "").replace(/^\*/, "");
    result[filename] = hash;
  }
  return result;
}

function calculateSha256(filePath) {
  return new Promise((resolve, reject) => {
    const hash = crypto.createHash("sha256");
    const stream = fs.createReadStream(filePath);

    stream.on("data", (chunk) => {
      hash.update(chunk);
    });
    stream.on("end", () => resolve(hash.digest("hex")));
    stream.on("error", reject);
  });
}

async function verifyChecksum(checksums, filename, filePath) {
  const expected = checksums[filename];
  if (!expected) {
    throw new Error(`checksum not found for ${filename}`);
  }
  const actual = await calculateSha256(filePath);
  if (actual !== expected) {
    throw new Error(`checksum mismatch for ${filename}`);
  }
}

function shouldSkipDownload(env) {
  if (env.CI === "true" || env.GITHUB_ACTIONS === "true") {
    return true;
  }
  const flag = env.PATCHLINE_SKIP_DOWNLOAD;
  if (!flag) {
    return false;
  }
  const normalized = flag.toLowerCase();
  return normalized === "1" || normalized === "true";
}

module.exports = {
  resolvePlatform,
  resolveArch,
  getArchiveExtension,
  getArchiveName,
  getDownloadUrl,
  getChecksumsUrl,
  getBinaryName,
  getInstalledBinaryName,
  getCandidateBinaryPaths,
  parseChecksums,
  verifyChecksum,
  shouldSkipDownload,
};
