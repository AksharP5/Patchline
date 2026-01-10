const assert = require("node:assert/strict");
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
