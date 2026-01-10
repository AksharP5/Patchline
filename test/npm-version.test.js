const assert = require("node:assert/strict");
const fs = require("node:fs");
const os = require("node:os");
const path = require("node:path");
const test = require("node:test");

const { normalizeTag, updatePackageVersion } = require("../lib/npm-version");

test("normalizeTag strips leading v", () => {
  assert.equal(normalizeTag("v1.2.3"), "1.2.3");
});

test("normalizeTag keeps plain version", () => {
  assert.equal(normalizeTag("2.0.0"), "2.0.0");
});

test("updatePackageVersion rewrites package.json version", () => {
  const root = fs.mkdtempSync(path.join(os.tmpdir(), "patchline-npm-"));
  const packagePath = path.join(root, "package.json");
  fs.writeFileSync(packagePath, JSON.stringify({ name: "patchline", version: "0.1.0" }));

  const updated = updatePackageVersion(packagePath, "v3.4.5");

  const contents = JSON.parse(fs.readFileSync(packagePath, "utf8"));
  assert.equal(updated, "3.4.5");
  assert.equal(contents.version, "3.4.5");
});
