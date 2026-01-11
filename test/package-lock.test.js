const assert = require("node:assert/strict");
const fs = require("node:fs");
const path = require("node:path");
const test = require("node:test");

const lockPath = path.join(__dirname, "..", "package-lock.json");

test("package-lock.json exists for npm ci", () => {
  assert.ok(fs.existsSync(lockPath), "package-lock.json is missing");
});
