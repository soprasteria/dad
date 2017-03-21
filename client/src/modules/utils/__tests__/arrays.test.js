import { sortby } from '../arrays';

describe('comparator to sort an array by key', () => {
  describe('with an empty array', () => {
    it('should return an empty array', () => {
      const emptyArray = [];
      const expected = [];
      const actual = emptyArray.sort(sortby(''));
      expect(actual).toEqual(expected);
    });
  }),
  describe('with an unsorted array of objects', () => {
    describe('with an existing key', () => {
      it('should sort the array by key', () => {
        const unsortedArray = [
          { key: 'banana' },
          { key: 'apple' },
          { key: 'orange' },
          { key: 'kiwi' }
        ];
        const expected = [
          { key: 'apple' },
          { key: 'banana' },
          { key: 'kiwi' },
          { key: 'orange' }
        ];
        const actual = unsortedArray.sort(sortby('key'));
        expect(actual).toEqual(expected);
      });
    }),
    describe('with a non-existing key', () => {
      it('should throw an error', () => {
        const unsortedArray = [
          { key: 'banana' },
          { key: 'apple' },
          { key: 'orange' },
          { key: 'kiwi' }
        ];
        const actual = () => unsortedArray.sort(sortby('wrongKey'));
        expect(actual).toThrowError(TypeError);
      });
    });
  });
});
