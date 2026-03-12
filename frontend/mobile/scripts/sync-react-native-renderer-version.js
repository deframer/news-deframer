const fs = require('fs');
const path = require('path');

const projectRoot = path.resolve(__dirname, '..');
const reactPackagePath = path.join(projectRoot, 'node_modules', 'react', 'package.json');
const reactNativePackagePath = path.join(projectRoot, 'node_modules', 'react-native', 'package.json');
const rendererDir = path.join(projectRoot, 'node_modules', 'react-native', 'Libraries', 'Renderer', 'implementations');

function readJson(filePath) {
  return JSON.parse(fs.readFileSync(filePath, 'utf8'));
}

function writeJson(filePath, value) {
  fs.writeFileSync(filePath, `${JSON.stringify(value, null, 2)}\n`);
}

function replaceAll(filePath, fromVersion, toVersion) {
  const source = fs.readFileSync(filePath, 'utf8');
  if (!source.includes(fromVersion)) {
    return false;
  }

  const updated = source.split(fromVersion).join(toVersion);
  fs.writeFileSync(filePath, updated);
  return true;
}

function main() {
  if (!fs.existsSync(reactPackagePath) || !fs.existsSync(reactNativePackagePath) || !fs.existsSync(rendererDir)) {
    console.log('[sync-react-native-renderer-version] Skipping; React Native install not present yet.');
    return;
  }

  const reactPackage = readJson(reactPackagePath);
  const reactNativePackage = readJson(reactNativePackagePath);
  const reactVersion = reactPackage.version;
  const expectedRange = reactNativePackage.peerDependencies && reactNativePackage.peerDependencies.react;
  const embeddedVersion = typeof expectedRange === 'string' ? expectedRange.replace(/^\^/, '') : null;

  if (!reactVersion || !embeddedVersion || reactVersion === embeddedVersion) {
    console.log('[sync-react-native-renderer-version] Nothing to patch.');
    return;
  }

  const rendererFiles = fs
    .readdirSync(rendererDir)
    .filter((name) => /^React(NativeRenderer|Fabric)-(dev|prod|profiling)\.js$/.test(name))
    .map((name) => path.join(rendererDir, name));

  let updatedFiles = 0;
  for (const filePath of rendererFiles) {
    if (replaceAll(filePath, embeddedVersion, reactVersion)) {
      updatedFiles += 1;
    }
  }

  reactNativePackage.peerDependencies = {
    ...reactNativePackage.peerDependencies,
    react: `^${reactVersion}`,
  };
  writeJson(reactNativePackagePath, reactNativePackage);

  console.log(
    `[sync-react-native-renderer-version] Patched embedded renderer version ${embeddedVersion} -> ${reactVersion} in ${updatedFiles} files.`
  );
}

main();
