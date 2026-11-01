#!/usr/bin/env node

const { execFileSync } = require("node:child_process")

const platformMap = {
	darwin: "darwin",
	linux: "linux",
	win32: "windows",
}

const archMap = {
	x64: "amd64",
	arm64: "arm64",
}

const platform = platformMap[process.platform]
const arch = archMap[process.arch]

if (!platform || !arch) {
	console.error(
		`Unsupported platform/architecture: ${process.platform}/${process.arch}`,
	)
	process.exit(1)
}

const packageName = `@codeshelldev/goplater-${platform}-${arch}`
const binaryName = platform === "windows" ? "goplater.exe" : "goplater"

let binaryPath

try {
	binaryPath = require.resolve(`${packageName}/bin/${binaryName}`)
} catch {
	console.error(
		`Could not find the GoPlater binary for ${platform}/${arch}.\n` +
			`Expected package: ${packageName}\n` +
			"Try reinstalling @codeshelldev/goplater.",
	)
	process.exit(1)
}

try {
	execFileSync(binaryPath, process.argv.slice(2), {
		stdio: "inherit",
	})
} catch (error) {
	if (error.status !== undefined) {
		process.exit(error.status)
	}

	throw error
}
