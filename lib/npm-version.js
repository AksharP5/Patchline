const fs = require("node:fs");

const semverPattern = /^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-[0-9A-Za-z.-]+)?(?:\+[0-9A-Za-z.-]+)?$/;

function normalizeTag(tag) {
  if (!tag) {
    throw new Error("version tag is required");
  }
  const trimmed = tag.trim();
  const version = trimmed.startsWith("v") ? trimmed.slice(1) : trimmed;
  if (!version) {
    throw new Error("version tag is required");
  }
  if (!semverPattern.test(version)) {
    throw new Error(`invalid semver: ${version}`);
  }
  return version;
}

function updatePackageVersion(packagePath, tag) {
  const version = normalizeTag(tag);
  const raw = fs.readFileSync(packagePath, "utf8");
  const pkg = JSON.parse(raw);
  pkg.version = version;
  fs.writeFileSync(packagePath, `${JSON.stringify(pkg, null, 2)}\n`);
  return version;
}

module.exports = {
  normalizeTag,
  updatePackageVersion,
};
