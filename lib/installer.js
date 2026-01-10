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

module.exports = {
  resolvePlatform,
  resolveArch,
  getArchiveExtension,
  getArchiveName,
  getDownloadUrl,
  getBinaryName,
  getInstalledBinaryName,
  getCandidateBinaryPaths,
};
