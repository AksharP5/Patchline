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

module.exports = {
  resolvePlatform,
  resolveArch,
  getArchiveExtension,
  getArchiveName,
  getDownloadUrl,
  getBinaryName,
};
