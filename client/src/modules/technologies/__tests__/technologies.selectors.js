import { flattenTechnologies } from '../technologies.selectors';

describe('flattenTechnologies', () => {
  describe('with an empty object', () => {
    it('should return an empty array', () => {
      const expected = [];
      const actual = flattenTechnologies({});
      expect(actual).toEqual(expected);
    });
  });

  describe('with a single element array', () => {
    it('should format the technologies object received from the server as a flat array', () => {
      const technologies = {
        '5901c690ec732920805784c5': {
          id: '5901c690ec732920805784c5',
          name: 'Java'
        }
      };
      const expected = ['Java'];
      const actual = flattenTechnologies(technologies);
      expect(actual).toEqual(expected);
    });
  });

  describe('with a multiple element array', () => {
    it('should format the technologies object array from the server as a flat array', () => {
      const technologies = {
        'ObjectId_1': {
          id: 'ObjectId_1',
          name: 'Java'
        },
        'ObjectId_2': {
          id: 'ObjectId_2',
          name: '.NET'
        },
        'ObjectId_3': {
          id: 'ObjectId_3',
          name: 'Cobol'
        }
      };
      const expected = [
        'Java',
        '.NET',
        'Cobol'
      ];
      const actual = flattenTechnologies(technologies);
      expect(actual).toEqual(expected);
    });
  });
});
