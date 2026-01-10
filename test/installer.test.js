const assert = require("node:assert/strict");
const crypto = require("node:crypto");
const fs = require("node:fs");
const os = require("node:os");
const path = require("node:path");
const test = require("node:test");

const installer = require("../lib/installer");

test("resolvePlatform maps win32 to windows", () => {
  assert.equal(installer.resolvePlatform("win32"), "windows");
});

test("resolveArch maps x64 to amd64", () => {
  assert.equal(installer.resolveArch("x64"), "amd64");
});

test("getArchiveName builds Windows archive name", () => {
  assert.equal(
    installer.getArchiveName("1.2.3", "windows", "amd64"),
    "patchline_1.2.3_windows_amd64.zip",
  );
});

test("getDownloadUrl builds GitHub release URL", () => {
  assert.equal(
    installer.getDownloadUrl("1.2.3", "darwin", "arm64"),
    "https://github.com/AksharP5/Patchline/releases/download/v1.2.3/patchline_1.2.3_darwin_arm64.tar.gz",
  );
});

test("getBinaryName adds .exe on Windows", () => {
  assert.equal(installer.getBinaryName("windows"), "patchline.exe");
});

test("getInstalledBinaryName uses bin suffix", () => {
  assert.equal(installer.getInstalledBinaryName("linux"), "patchline-bin");
});

test("getCandidateBinaryPaths includes nested archive folder", () => {
  const candidates = installer.getCandidateBinaryPaths(
    path.join("/", "tmp"),
    "1.2.3",
    "linux",
    "amd64",
  );

  assert.deepEqual(candidates, [
    path.join("/", "tmp", "patchline"),
    path.join("/", "tmp", "patchline_1.2.3_linux_amd64", "patchline"),
  ]);
});

test("getChecksumsUrl builds release checksum URL", () => {
  assert.equal(
    installer.getChecksumsUrl("1.2.3"),
    "https://github.com/AksharP5/Patchline/releases/download/v1.2.3/checksums.txt",
  );
});

test("parseChecksums maps filenames to hashes", () => {
  const text = "abc123  file-one.tar.gz\nfff999  *file-two.zip\n";
  const checksums = installer.parseChecksums(text);

  assert.equal(checksums["file-one.tar.gz"], "abc123");
  assert.equal(checksums["file-two.zip"], "fff999");
});

test("verifyChecksum succeeds for matching hash", async () => {
  const root = fs.mkdtempSync(path.join(os.tmpdir(), "patchline-hash-"));
  const filePath = path.join(root, "test.txt");
  fs.writeFileSync(filePath, "hello");

  const hash = crypto.createHash("sha256").update("hello").digest("hex");
  const checksums = { "test.txt": hash };

  await installer.verifyChecksum(checksums, "test.txt", filePath);
});

test("verifyChecksum rejects mismatched hash", async () => {
  const root = fs.mkdtempSync(path.join(os.tmpdir(), "patchline-hash-"));
  const filePath = path.join(root, "test.txt");
  fs.writeFileSync(filePath, "hello");

  const checksums = { "test.txt": "deadbeef" };

  await assert.rejects(
    () => installer.verifyChecksum(checksums, "test.txt", filePath),
    /checksum mismatch/,
  );
});

test("shouldSkipDownload honors env flag", () => {
  assert.equal(installer.shouldSkipDownload({ PATCHLINE_SKIP_DOWNLOAD: "1" }), true);
  assert.equal(installer.shouldSkipDownload({ PATCHLINE_SKIP_DOWNLOAD: "true" }), true);
  assert.equal(installer.shouldSkipDownload({ PATCHLINE_SKIP_DOWNLOAD: "" }), false);
});
