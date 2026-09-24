import { describe, it, expect } from 'vitest'
import { RANDOM_KEY_NAMES, getRandomKeyName } from '../randomKeyNames'

describe('randomKeyNames', () => {
  it('should have more than 200 preset names', () => {
    expect(RANDOM_KEY_NAMES.length).toBeGreaterThanOrEqual(200)
  })

  it('should include user-requested names like 卡卡罗特, 哪吒, 魔童, 小丸子', () => {
    expect(RANDOM_KEY_NAMES).toContain('卡卡罗特')
    expect(RANDOM_KEY_NAMES).toContain('哪吒')
    expect(RANDOM_KEY_NAMES).toContain('魔童')
    expect(RANDOM_KEY_NAMES).toContain('小丸子')
  })

  it('should return a non-empty string that exists in the preset array', () => {
    for (let i = 0; i < 50; i++) {
      const name = getRandomKeyName()
      expect(typeof name).toBe('string')
      expect(name.length).toBeGreaterThan(0)
      expect(RANDOM_KEY_NAMES).toContain(name as any)
    }
  })
})
