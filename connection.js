const fs = require('fs').promises;
const { readFileSync } = require('fs');
const { exec } = require('child_process');
const path = require('path');
const os = require('os');
const util = require('util');

const execAsync = util.promisify(exec);
const thresholdAccuracy = 90;

async function processNotebook(base64String, originalFileName) {
  const decodedDir = path.join(os.homedir(), 'decoded_notebooks');
  const outputDir = path.join(os.homedir(), 'output');

  const decodedPath = path.join(decodedDir, originalFileName);
  const outputName = originalFileName.replace(/\.ipynb$/, '_output.txt');
  const outputPath = path.join(outputDir, outputName);

  await fs.mkdir(decodedDir, { recursive: true });
  await fs.mkdir(outputDir, { recursive: true });

  const decoded = Buffer.from(base64String, 'base64').toString('utf-8');
  await fs.writeFile(decodedPath, decoded, 'utf-8');

  const batPath = path.resolve(__dirname, 'run_go_compiler.bat');
  const command = `"${batPath}" "${decodedPath}"`;

  console.log(originalFileName);

  let stdout;
  try {
    const { stdout: rawOutput } = await execAsync(command);
    stdout = rawOutput;
  } catch (err) {
    console.error("[processNotebook] ❌ Failed to run Go compiler via .bat:", err);
    return "Failure";
  }

  let accuracy = 0;
  try {
    const parsed = JSON.parse(stdout.trim().split("\n").pop());
    accuracy = parsed.accuracy || 0;
  } catch (err) {
    console.error("[processNotebook] ❌ Failed to parse accuracy from stdout:", err);
    return "Failure";
  }

  return accuracy >= thresholdAccuracy ? "Success" : "Failure";
}

module.exports = { processNotebook };
