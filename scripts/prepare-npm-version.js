const path = require("node:path");

const { updatePackageVersion } = require("../lib/npm-version");

const tag = process.env.NPM_VERSION || process.env.GITHUB_REF_NAME;
if (!tag) {
  console.error("NPM_VERSION or GITHUB_REF_NAME is required to set the npm version.");
  process.exit(1);
}

try {
  const version = updatePackageVersion(path.join(__dirname, "..", "package.json"), tag);
  console.log(`Prepared npm package version ${version}.`);
} catch (error) {
  console.error(error.message);
  process.exit(1);
}
