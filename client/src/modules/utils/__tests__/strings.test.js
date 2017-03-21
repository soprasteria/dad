import { containsWithoutAccents } from '../strings';

describe('utility method to check if a string contains another string, case and accent insensitive', () => {
  describe('an empty string', () => {
    it('should contain the empty string', () => {
      const needle = '';
      const haystack = '';
      expect(containsWithoutAccents(haystack, needle)).toBe(true);
    });
  }),
  describe('a string without accents', () => {
    it('should contain itself', () => {
      const needle = 'Lorem ipsum dolor sit amet';
      const haystack = 'Lorem ipsum dolor sit amet';
      expect(containsWithoutAccents(haystack, needle)).toBe(true);
    }),
    it('should contain itself, case insensitively', () => {
      const needle = 'Lorem ipsum dolor sit amet';
      const haystack = 'lorem ipsum dolor sit amet';
      expect(containsWithoutAccents(haystack, needle)).toBe(true);
    });
  }),
  describe('a string with accents', () => {
    it('should contain itself', () => {
      const needle = 'Ante vulputaté éros convallis etiam';
      const haystack = 'Ante vulputaté éros convallis etiam';
      expect(containsWithoutAccents(haystack, needle)).toBe(true);
    }),
    it('should contain itself without accents and case insensitively', () => {
      const needle = 'ante vulputate eros convallis etiam';
      const haystack = 'Ante vulputaté éros convallis etiam';
      expect(containsWithoutAccents(haystack, needle)).toBe(true);
    }),
    it('should contain a word from the string without accents and case insensitively', () => {
      const needle = 'vulputate';
      const haystack = 'Ante vulputaté éros convallis etiam';
      expect(containsWithoutAccents(haystack, needle)).toBe(true);
    });
  });
});
