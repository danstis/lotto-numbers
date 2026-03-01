// Custom Node.js test reporter that outputs SonarQube Generic Test Execution XML format.
// See: https://docs.sonarsource.com/sonarqube-server/latest/analyzing-source-code/test-coverage/generic-test-data/

'use strict'

const path = require('node:path')
const cwd = process.cwd()

function relativePath(file) {
  return path.relative(cwd, file).replace(/\\/g, '/')
}

function escapeXml(str) {
  return String(str)
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
}

module.exports = async function* reporter(source) {
  const byFile = new Map()

  for await (const event of source) {
    if (event.type !== 'test:pass' && event.type !== 'test:fail') continue
    if (event.data.nesting !== 0) continue

    const file = relativePath(event.data.file || 'unknown')
    if (!byFile.has(file)) byFile.set(file, [])

    const duration = Math.round(event.data.details?.duration_ms ?? 0)

    if (event.type === 'test:pass') {
      byFile.get(file).push({ name: event.data.name, duration, failed: false })
    } else {
      const err = event.data.details?.error
      byFile.get(file).push({
        name: event.data.name,
        duration,
        failed: true,
        message: err?.message ?? 'Test failed',
        stack: err?.stack ?? '',
      })
    }
  }

  yield '<?xml version="1.0" encoding="UTF-8"?>\n'
  yield '<testExecutions version="1">\n'
  for (const [file, cases] of byFile) {
    yield `  <file path="${escapeXml(file)}">\n`
    for (const tc of cases) {
      if (tc.failed) {
        yield `    <testCase name="${escapeXml(tc.name)}" duration="${tc.duration}">\n`
        yield `      <failure message="${escapeXml(tc.message)}">${escapeXml(tc.stack)}</failure>\n`
        yield `    </testCase>\n`
      } else {
        yield `    <testCase name="${escapeXml(tc.name)}" duration="${tc.duration}"/>\n`
      }
    }
    yield `  </file>\n`
  }
  yield '</testExecutions>\n'
}
