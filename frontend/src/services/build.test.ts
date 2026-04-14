import { describe, it, expect } from 'vitest'
import * as fs from 'fs'
import * as path from 'path'

describe('빌드 및 배포 설정 검증', () => {
  it('shouldHaveDockerfileForBackend', () => {
    const dockerfilePath = path.resolve(__dirname, '../../../backend/Dockerfile')
    expect(fs.existsSync(dockerfilePath)).toBe(true)
  })

  it('shouldHaveDockerComposeFile', () => {
    const composePath = path.resolve(__dirname, '../../../docker-compose.yml')
    expect(fs.existsSync(composePath)).toBe(true)
  })

  it('shouldHaveConfigYamlExample', () => {
    const configPath = path.resolve(__dirname, '../../../backend/config.yaml.example')
    expect(fs.existsSync(configPath)).toBe(true)
  })

  it('shouldHaveVitestConfig', () => {
    const vitestPath = path.resolve(__dirname, '../../vitest.config.ts')
    expect(fs.existsSync(vitestPath)).toBe(true)
  })

  it('shouldHaveNextConfig', () => {
    const nextPath = path.resolve(__dirname, '../../next.config.ts')
    expect(fs.existsSync(nextPath)).toBe(true)
  })

  it('shouldHaveGoMod', () => {
    const goModPath = path.resolve(__dirname, '../../../backend/go.mod')
    expect(fs.existsSync(goModPath)).toBe(true)
  })

  it('shouldHavePackageJson', () => {
    const pkgPath = path.resolve(__dirname, '../../package.json')
    expect(fs.existsSync(pkgPath)).toBe(true)
  })
})
