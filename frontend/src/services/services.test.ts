import { describe, it, expect } from 'vitest'
import { authApi } from './auth'
import { pluginsApi } from './plugins'
import { reviewsApi } from './reviews'
import { adminApi } from './admin'

describe('API 서비스 모듈', () => {
  it('shouldExportAuthApiMethods', () => {
    expect(authApi.register).toBeDefined()
    expect(authApi.login).toBeDefined()
    expect(authApi.refresh).toBeDefined()
    expect(authApi.getMe).toBeDefined()
  })

  it('shouldExportPluginsApiMethods', () => {
    expect(pluginsApi.getList).toBeDefined()
    expect(pluginsApi.getById).toBeDefined()
    expect(pluginsApi.create).toBeDefined()
    expect(pluginsApi.update).toBeDefined()
    expect(pluginsApi.delete).toBeDefined()
    expect(pluginsApi.install).toBeDefined()
    expect(pluginsApi.uninstall).toBeDefined()
    expect(pluginsApi.toggleActive).toBeDefined()
    expect(pluginsApi.getMyInstallations).toBeDefined()
  })

  it('shouldExportReviewsApiMethods', () => {
    expect(reviewsApi.getByPluginId).toBeDefined()
    expect(reviewsApi.create).toBeDefined()
    expect(reviewsApi.update).toBeDefined()
    expect(reviewsApi.delete).toBeDefined()
  })

  it('shouldExportAdminApiMethods', () => {
    expect(adminApi.getPendingPlugins).toBeDefined()
    expect(adminApi.approve).toBeDefined()
    expect(adminApi.reject).toBeDefined()
    expect(adminApi.hide).toBeDefined()
  })
})
