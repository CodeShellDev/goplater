#!/usr/bin/env node
const { execFileSync } = require("child_process")
const fs = require("fs")
const path = require("path")

const platformMap = { darwin: "darwin", linux: "linux", win32: "windows" }
const archMap = { x64: "amd64", arm64: "arm64" }

const plat = platformMap[process.platform]
const arch = archMap[process.arch]

if (!plat || !arch) {
	console.error(
		`Unsupported platform/arch: ${process.platform}/${process.arch}`,
	)
	process.exit(1)
}

const distDir = path.resolve(__dirname)
const files = fs.readdirSync(distDir)

let binaryFile = files.find((f) => f.includes(`${plat}-${arch}`))
if (!binaryFile) {
	console.error(`No binary found for ${plat}-${arch} in ${distDir}`)
	process.exit(1)
}

if (plat === "windows" && !binaryFile.endsWith(".exe")) {
	binaryFile += ".exe"
}

const binaryPath = path.join(distDir, binaryFile)
execFileSync(binaryPath, process.argv.slice(2), { stdio: "inherit" })
