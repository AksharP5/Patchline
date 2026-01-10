const fs = require("node:fs");
const fsPromises = require("node:fs/promises");
const https = require("node:https");
const os = require("node:os");
const path = require("node:path");
const { pipeline } = require("node:stream/promises");

const extractZip = require("extract-zip");
const tar = require("tar");

const {
  getArchiveExtension,
  getArchiveName,
  getBinaryName,
  getCandidateBinaryPaths,
  getChecksumsUrl,
  getDownloadUrl,
  getInstalledBinaryName,
  parseChecksums,
  resolveArch,
  resolvePlatform,
  shouldSkipDownload,
  verifyChecksum,
} = require("../lib/installer");

const packageJson = require("../package.json");

function request(url) {
  return new Promise((resolve, reject) => {
    const req = https.get(url, (res) => {
      if (res.statusCode >= 300 && res.statusCode < 400 && res.headers.location) {
        res.resume();
        resolve(request(res.headers.location));
        return;
      }

      if (res.statusCode !== 200) {
        res.resume();
        reject(new Error(`download failed: ${res.statusCode} ${res.statusMessage}`));
        return;
      }

      resolve(res);
    });

    req.on("error", reject);
  });
}

async function downloadFile(url, destination) {
  const response = await request(url);
  await pipeline(response, fs.createWriteStream(destination));
}

async function downloadText(url) {
  const response = await request(url);
  const chunks = [];
  for await (const chunk of response) {
    chunks.push(chunk);
  }
  return Buffer.concat(chunks).toString("utf8");
}

async function extractArchive(archivePath, extractDir, extension) {
  if (extension === "zip") {
    await extractZip(archivePath, { dir: extractDir });
    return;
  }
  await tar.x({ file: archivePath, cwd: extractDir });
}

async function findExistingPath(paths) {
  for (const candidate of paths) {
    try {
      await fsPromises.access(candidate);
      return candidate;
    } catch (error) {
      continue;
    }
  }
  return "";
}

async function install() {
  if (shouldSkipDownload(process.env)) {
    console.log("Skipping Patchline download because PATCHLINE_SKIP_DOWNLOAD is set.");
    return;
  }

  const version = packageJson.version;
  const goos = resolvePlatform(process.platform);
  const goarch = resolveArch(process.arch);
  const extension = getArchiveExtension(goos);
  const archiveName = getArchiveName(version, goos, goarch);
  const url = getDownloadUrl(version, goos, goarch);
  const checksumsUrl = getChecksumsUrl(version);

  const tempRoot = await fsPromises.mkdtemp(path.join(os.tmpdir(), "patchline-"));
  const archivePath = path.join(tempRoot, archiveName);
  const extractDir = path.join(tempRoot, "extract");

  await fsPromises.mkdir(extractDir, { recursive: true });
  await downloadFile(url, archivePath);

  const checksumsText = await downloadText(checksumsUrl);
  const checksums = parseChecksums(checksumsText);
  await verifyChecksum(checksums, archiveName, archivePath);

  await extractArchive(archivePath, extractDir, extension);

  const binaryName = getBinaryName(goos);
  const candidates = getCandidateBinaryPaths(extractDir, version, goos, goarch);
  const extractedBinaryPath = await findExistingPath(candidates);

  if (!extractedBinaryPath) {
    throw new Error(`binary ${binaryName} not found in archive`);
  }

  const installedBinaryName = getInstalledBinaryName(goos);
  const targetPath = path.join(__dirname, "..", "bin", installedBinaryName);

  await fsPromises.mkdir(path.dirname(targetPath), { recursive: true });
  await fsPromises.copyFile(extractedBinaryPath, targetPath);

  if (goos !== "windows") {
    await fsPromises.chmod(targetPath, 0o755);
  }
}

install().catch((error) => {
  console.error(error.message);
  process.exit(1);
});
